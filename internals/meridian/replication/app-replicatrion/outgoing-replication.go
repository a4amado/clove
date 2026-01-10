package AppReplication

import (
	"clove/internals/cache"
	repository "clove/internals/services/generatedRepo"
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ReplicatableAppMsg represents an app that can be replicated across regions
type ReplicatableAppMsg struct {
	repository.App
}

// PublishReplicatableAppMsgToRabbitMQ publishes an app replication message to RabbitMQ
// for distribution to other regions. It first attempts to save locally, and on
// failure sends to all regions to ensure eventual consistency.
func (ar *AppReplication) PublishReplicatableAppMsgToRabbitMQ(ctx context.Context, msg ReplicatableAppMsg) error {
	messageBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshaling message: %w", err)
	}

	var targetRegions []repository.Region

	// Try saving to local valkey first
	err = cache.Apps().Set(ctx, msg.App)
	if err != nil {
		// If local save fails, send to all regions including source
		// to ensure eventual consistency
		targetRegions = repository.AllRegionValues()
	}

	if len(targetRegions) == 0 {
		return nil
	}

	// Publish to each target region
	for _, region := range targetRegions {
		var ch *amqp.Channel
		if region == ar.region {
			ch = ar.localRabbitMQChannel
		} else {
			ch = ar.crossRegionChannels[region]
		}
		if ch == nil {
			continue
		}

		routingKey := fmt.Sprintf("%s.app-replication", region)
		err := ch.PublishWithContext(
			ctx,
			ar.exchangeName, // exchange
			routingKey,      // routing key
			false,           // mandatory
			false,           // immediate
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         messageBytes,
				DeliveryMode: amqp.Persistent,
				Timestamp:    time.Now(),
				MessageId:    fmt.Sprintf("%s", msg.ID.Bytes),
			},
		)
		if err != nil {
			return fmt.Errorf("publishing to region %s: %w", region, err)
		}
	}

	return nil
}
