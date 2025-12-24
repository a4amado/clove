package meridian

import (
	envConsts "clove/internals/consts/env"
	"clove/internals/repository"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/segmentio/kafka-go"
)

// App here to the apps that the user own
type ReplicatableAppMsg struct {
	repository.App
	// Copied it so You can add additional metadata fields here if needed
}

func (m *ReplicatableAppMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(ReplicatableAppMsg{App: m.App})
}

var regionKafkaWriters = map[repository.Region]*kafka.Writer{
	repository.RegionDk1: {
		Addr:                   kafka.TCP(os.Getenv(string(envConsts.KAFKA_BOOTSTRAP))),
		Topic:                  fmt.Sprintf("%s-app-replication", repository.RegionDk1),
		Balancer:               &kafka.RoundRobin{},
		MaxAttempts:            3,
		WriteTimeout:           10 * time.Second,
		AllowAutoTopicCreation: true,
		RequiredAcks:           kafka.RequireOne,
		Compression:            kafka.Gzip,
	},
}

// getCurrentMachineRegion returns the region for the current machine as a repository.Region.
// In non-production environments it defaults the configured region to "dk1". It panics if the derived region is not valid.
func getCurrentMachineRegion() repository.Region {

	region := os.Getenv(string(envConsts.REGION))
	parsedRegion := repository.Region(region)
	// if no region is set, panic
	if !parsedRegion.Valid() {
		panic("no valid region set in environment variable")
	}
	return parsedRegion
}

var ReplicatableAppMsgReader = kafka.NewReader(kafka.ReaderConfig{
	Brokers: []string{os.Getenv(string(envConsts.KAFKA_BOOTSTRAP))},
	Topic:   fmt.Sprintf("%s-app-replication", getCurrentMachineRegion()),
	GroupID: fmt.Sprintf("%s-app-replication-group", getCurrentMachineRegion()),
})

func (c *Meridian) PublishReplicatableAppMsgToKafka(ctx context.Context, msg ReplicatableAppMsg) error {
	messageBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	var targetRegions []repository.Region = []repository.Region{}
	// 1. try saving the app to local redis first
	err = c.SaveApp(ctx, msg.App)
	// 2. if fails, send it via kafka to all regions (including source region), to ensure eventual consistency
	// 3. bc i don't want to deal with retry here
	if err != nil {
		targetRegions = repository.AllRegionValues()
	} else {
		targetRegions = slices.DeleteFunc(repository.AllRegionValues(), func(r repository.Region) bool {
			return getCurrentMachineRegion() == r
		})

	}

	if len(targetRegions) == 0 {
		return nil
	}
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
			return err
		}

	}
	return nil

}

// StartKafkaConsumer starts a kafka consumer that listens for app replication messages
// and saves them to the local redis instance
func (c *Meridian) StartKafkaToRedisBridge(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			kafkaMessage, err := ReplicatableAppMsgReader.ReadMessage(ctx)
			if err != nil {
				continue
			}
			var ReplicatableAppMsg ReplicatableAppMsg
			if err := json.Unmarshal(kafkaMessage.Value, &ReplicatableAppMsg); err != nil {
				continue
			}
			err = c.SaveApp(ctx, ReplicatableAppMsg.App)
			if err != nil {
				continue
			}

			err = ReplicatableAppMsgReader.CommitMessages(ctx, kafkaMessage)
			if err != nil {
				continue
			}

		}
	}
}
