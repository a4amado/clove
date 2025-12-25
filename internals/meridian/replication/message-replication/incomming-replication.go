package MessageReplication

import (
	"clove/internals/meridian/fanout"
	"context"
	"encoding/json"
)

// StartKafkaToRedisBridge starts a Kafka consumer that continuously reads app
// replication messages and saves them to the local Redis instance.
// This function blocks until the context is cancelled.
func (c *MessageReplication) BridgeKafkaInternalDelevieryReplicatorToRedis(ctx context.Context) {
	fanoutClient := fanout.Fanout()
	for {
		select {
		case <-ctx.Done():
			return
		default:
			kafkaMessage, err := c.localReader.ReadMessage(ctx)
			if err != nil {
				// Log error in production, but continue processing
				continue
			}

			var msg InternalReplicatableDeliveryMsg
			if err := json.Unmarshal(kafkaMessage.Value, &msg); err != nil {
				// Log unmarshal error in production
				continue
			}

			key := fanoutClient.FormatChannelKey(fanout.ChannelKey{
				AppId:     msg.AppID,
				ChannelId: "test",
			})
			res := fanoutClient.Publish(ctx, key, msg.Payload)
			if res.Err() != nil {
				panic(res.Err())
			}
			err = c.localReader.CommitMessages(ctx, kafkaMessage)
			if err != nil {
				// Log commit error in production
				continue
			}
		}
	}
}
