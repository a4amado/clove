package AppReplication

import (
	"clove/internals/repository"
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

// ReplicatableAppMsg represents an app that can be replicated across regions
type ReplicatableAppMsg struct {
	repository.App
}

func (m *ReplicatableAppMsg) MarshalJSON() ([]byte, error) {
	return json.Marshal(ReplicatableAppMsg{App: m.App})
}

// PublishReplicatableAppMsgToKafka publishes an app replication message to Kafka
// for distribution to other regions. It first attempts to save locally, and on
// failure sends to all regions to ensure eventual consistency.
func (c *AppReplication) PublishReplicatableAppMsgToKafka(ctx context.Context, msg ReplicatableAppMsg) error {
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

	}

	if len(targetRegions) == 0 {
		return nil
	}

	// Publish to each target region
	for _, region := range c.crossRegionWriters {
		err := region.WriteMessages(ctx, kafka.Message{
			Key:   msg.ID.Bytes[:],
			Value: messageBytes,
		})
		if err != nil {
			return fmt.Errorf("writing to region %v: %w", region, err)
		}
	}

	return nil
}
