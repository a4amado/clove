package replication

import (
	envConsts "clove/internals/consts/env"
	"clove/internals/repository"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/segmentio/kafka-go"
)

// ReplicatableAppMsg represents an app that can be replicated across regions
type ReplicatableAppMsg struct {
	repository.App
}

func (m *ReplicatableAppMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(ReplicatableAppMsg{App: m.App})
}

// Package-level variables that need to be initialized after env vars are loaded
var (
	regionKafkaWriters       map[repository.Region]*kafka.Writer
	replicatableAppMsgReader *kafka.Reader
	currentMachineRegion     repository.Region
	meridianInitOnce         sync.Once
)

// Init initializes the meridian package with Kafka connections and region configuration.
// This must be called after environment variables are loaded.
// It is safe to call multiple times - initialization only happens once.
func Init() {
	meridianInitOnce.Do(func() {
		// Get and validate region
		godotenv.Load()
		region := os.Getenv(string(envConsts.REGION))
		kafkaBootstrap := os.Getenv(string(envConsts.KAFKA_BOOTSTRAP))
		if kafkaBootstrap == "" {
			panic("KAFKA_BOOTSTRAP environment variable not set")
		}
		currentMachineRegion = repository.Region(region)
		if !currentMachineRegion.Valid() {
			panic(fmt.Sprintf("no valid region set in environment variable, got: %q", region))
		}

		// Initialize Kafka writers for each region
		regionKafkaWriters = map[repository.Region]*kafka.Writer{
			repository.RegionDk1: {
				Addr:                   kafka.TCP(kafkaBootstrap),
				Topic:                  fmt.Sprintf("%s-app-replication", repository.RegionDk1),
				Balancer:               &kafka.RoundRobin{},
				MaxAttempts:            3,
				WriteTimeout:           10 * time.Second,
				AllowAutoTopicCreation: true,
				RequiredAcks:           kafka.RequireOne,
				Compression:            kafka.Gzip,
			},
			// Add other regions here as needed
		}

		// Initialize Kafka reader for this region
		replicatableAppMsgReader = kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{kafkaBootstrap},
			Topic:   fmt.Sprintf("%s-app-replication", currentMachineRegion),
			GroupID: fmt.Sprintf("%s-app-replication-group", currentMachineRegion),
		})
	})
}

// getCurrentMachineRegion returns the current machine's region.
// Panics if Init() has not been called.
func getCurrentMachineRegion() repository.Region {
	if !currentMachineRegion.Valid() {
		panic("meridian.Init() must be called before using getCurrentMachineRegion()")
	}
	return currentMachineRegion
}

// PublishReplicatableAppMsgToKafka publishes an app replication message to Kafka
// for distribution to other regions. It first attempts to save locally, and on
// failure sends to all regions to ensure eventual consistency.
func (c *Replication) PublishReplicatableAppMsgToKafka(ctx context.Context, msg ReplicatableAppMsg) error {
	messageBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshaling message: %w", err)
	}

	var targetRegions []repository.Region

	// Try saving to local redis first
	err = c.SaveApp(ctx, msg.App)
	if err != nil {
		// If local save fails, send to all regions including source
		// to ensure eventual consistency
		targetRegions = repository.AllRegionValues()
	} else {
		// Local save succeeded, only send to other regions
		targetRegions = slices.DeleteFunc(repository.AllRegionValues(), func(r repository.Region) bool {
			return getCurrentMachineRegion() == r
		})
	}

	if len(targetRegions) == 0 {
		return nil
	}

	// Publish to each target region
	for _, region := range targetRegions {
		writer, exists := regionKafkaWriters[region]
		if !exists {
			return fmt.Errorf("no kafka writer configured for region: %s", region)
		}

		err := writer.WriteMessages(ctx, kafka.Message{
			Key:   msg.ID.Bytes[:],
			Value: messageBytes,
		})
		if err != nil {
			return fmt.Errorf("writing to region %s: %w", region, err)
		}
	}

	return nil
}

// StartKafkaToRedisBridge starts a Kafka consumer that continuously reads app
// replication messages and saves them to the local Redis instance.
// This function blocks until the context is cancelled.
func (c *Replication) BridgeKafkaAppReplicatorToRedis(ctx context.Context) {
	if replicatableAppMsgReader == nil {
		panic("meridian.Init() must be called before starting Kafka bridge")
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			kafkaMessage, err := replicatableAppMsgReader.ReadMessage(ctx)
			if err != nil {
				// Log error in production, but continue processing
				continue
			}

			var msg ReplicatableAppMsg
			if err := json.Unmarshal(kafkaMessage.Value, &msg); err != nil {
				// Log unmarshal error in production
				continue
			}

			err = c.SaveApp(ctx, msg.App)
			if err != nil {
				// Log save error in production
				continue
			}

			err = replicatableAppMsgReader.CommitMessages(ctx, kafkaMessage)
			if err != nil {
				// Log commit error in production
				continue
			}
		}
	}
}
