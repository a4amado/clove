package AppReplication

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/segmentio/kafka-go"
)

// StartKafkaToRedisBridge starts a Kafka consumer that continuously reads app
// replication messages and saves them to the local Valkey instance.
// This function blocks until the context is cancelled.
func (c *AppReplication) BridgeKafkaAppReplicatorToRedis(ctx context.Context) {

	wg := sync.WaitGroup{}
	for _, region := range c.localReaders {
		wg.Add(1)
		pipe := make(chan ReplicatableAppMsg, 5)
		go func(region *kafka.Reader) {
			defer wg.Done()
			localWG := sync.WaitGroup{}
			localWG.Go(func() {
				for {
					select {
					case <-ctx.Done():
						return
					default:
						kafkaMessage, err := region.ReadMessage(ctx)
						if err != nil {
							panic(err)
						}
						var msg ReplicatableAppMsg
						if err := json.Unmarshal(kafkaMessage.Value, &msg); err != nil {
							panic(err)
						}
						select {
						case pipe <- msg:
						default:
						}

					}
				}

			})
			localWG.Go(func() {

				for {
					select {
					case msg, ok := <-pipe:
						if !ok {
							continue
						}
						err := c.SaveApp(ctx, msg.App)
						if err != nil {
							return
						}
					case <-ctx.Done():
						return
					default:

					}
				}

			})

			localWG.Wait()
		}(region)
	}

	wg.Wait()
}
