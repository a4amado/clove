package AppReplication

import (
	"clove/internals/cache"
	"context"
	"encoding/json"
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

// BridgeRabbitMQAppReplicatorToRedis starts a RabbitMQ consumer that continuously reads app
// replication messages and saves them to the local Valkey instance.
// This function blocks until the context is cancelled.
func (ar *AppReplication) BridgeRabbitMQAppReplicatorToRedis(ctx context.Context) {

	wg := sync.WaitGroup{}
	for _, consumerChannel := range ar.localConsumers {
		wg.Add(1)
		messagePipe := make(chan ReplicatableAppMsg, 5)
		go func(ch *amqp.Channel) {
			defer wg.Done()
			localWG := sync.WaitGroup{}

			// Start consumer
			localWG.Go(func() {
				msgs, err := ch.Consume(
					ar.localQueueName, // queue
					"",                // consumer tag (empty = auto-generate)
					false,             // auto-ack
					false,             // exclusive
					false,             // no-local
					false,             // no-wait
					nil,               // args
				)
				if err != nil {
					log.Printf("Failed to register consumer: %v", err)
					return
				}

				for {
					select {
					case <-ctx.Done():
						return
					case msg, ok := <-msgs:
						if !ok {
							return
						}
						var appMsg ReplicatableAppMsg
						if err := json.Unmarshal(msg.Body, &appMsg); err != nil {
							log.Printf("unmarshal error: %v", err)
							msg.Nack(false, false)
							continue
						}
						select {
						case messagePipe <- appMsg:
							msg.Ack(false)
						case <-ctx.Done():
							msg.Nack(false, true)
							return
						default:
							msg.Nack(false, true)
						}
					}
				}
			})

			// Start processor
			localWG.Go(func() {
				for {
					select {
					case msg, ok := <-messagePipe:
						if !ok {
							continue
						}
						err := cache.Apps().Set(ctx, msg.App)
						if err != nil {
							log.Printf("save app error: %v", err)
							return
						}
					case <-ctx.Done():
						return
					default:
					}
				}
			})

			localWG.Wait()
		}(consumerChannel)
	}

	wg.Wait()
}
