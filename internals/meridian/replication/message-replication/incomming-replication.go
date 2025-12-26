package MessageReplication

import (
	"clove/internals/meridian/fanout"
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"
)

// StartKafkaToRedisBridge starts a Kafka consumer that continuously reads app
// replication messages and saves them to the local Valkey instance.
// This function blocks until the context is cancelled.

func (c *MessageReplication) BridgeKafkaInternalDelevieryReplicatorToRedis(ctx context.Context) {

	loopWg := sync.WaitGroup{}
	for _, reader := range c.localReaders {
		loopWg.Add(1)
		go func(reader *kafka.Reader) {
			buffer := make(chan kafka.Message, 1000)
			childWg := sync.WaitGroup{}
			fanoutClient := fanout.Fanout()
			// start local reader
			childWg.Go(func() {
				for {

					select {
					case <-ctx.Done():
						return
					default:
						msg, err := reader.ReadMessage(ctx)

						if err != nil {
							if ctx.Err() != nil {
								return
							}
							log.Printf("kafka read error: %v", err)
							continue
						}
						buffer <- msg
					}
				}
			})
			// start local writer
			childWg.Go(func() {
				for msg := range buffer {
					var appMsg InternalReplicatableDeliveryMsg
					if err := json.Unmarshal(msg.Value, &appMsg); err != nil {
						log.Printf("unmarshal error: %v", err)
						continue
					}

					key := fanoutClient.FormatChannelKey(fanout.ChannelKey{
						AppId:     appMsg.AppID,
						ChannelId: "test",
					})

					if err := fanoutClient.Publish(ctx, key, appMsg.Payload); err != nil {
						log.Printf("valkey publish error: %v", err)
						continue
					}

					if err := reader.CommitMessages(ctx, msg); err != nil {
						log.Printf("commit error: %v", err)
					}
				}
			})
			childWg.Wait()
		}(reader)
	}
	loopWg.Wait()

}
