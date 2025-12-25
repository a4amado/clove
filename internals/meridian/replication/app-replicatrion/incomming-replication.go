package AppReplication

import (
	"context"
	"encoding/json"
)

// StartKafkaToRedisBridge starts a Kafka consumer that continuously reads app
// replication messages and saves them to the local Redis instance.
// This function blocks until the context is cancelled.
func (c *AppReplication) BridgeKafkaAppReplicatorToRedis(ctx context.Context) {

	for {
		select {
		case <-ctx.Done():
			return
		default:
			kafkaMessage, err := c.localReader.ReadMessage(ctx)
			if err != nil {
				panic(err)
				// Log error in production, but continue processing
				continue
			}

			var msg ReplicatableAppMsg
			if err := json.Unmarshal(kafkaMessage.Value, &msg); err != nil {
				// Log unmarshal error in production
				panic(err)
				continue
			}

			err = c.SaveApp(ctx, msg.App)
			if err != nil {
				// Log save error in production
				panic(err)
				continue
			}

			err = c.localReader.CommitMessages(ctx, kafkaMessage)
			if err != nil {
				// Log commit error in production
				panic(err)
				continue
			}
		}
	}
}
