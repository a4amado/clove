package MessageReplication

import (
	"clove/internals/meridian/fanout"
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

// BridgeRabbitMQInternalDeliveryReplicatorToRedis starts a RabbitMQ consumer that continuously reads
// replication messages and saves them to the local Valkey instance.
// This function blocks until the context is cancelled.

func (mr *MessageReplication) BridgeRabbitMQInternalDeliveryReplicatorToRedis(ctx context.Context) {

	loopWg := sync.WaitGroup{}
	for _, consumerChannel := range mr.localConsumers {
		loopWg.Add(1)
		go func(ch *amqp.Channel) {
			messageBuffer := make(chan amqp.Delivery, 1000)
			childWg := sync.WaitGroup{}
			fanoutClient := fanout.Fanout()
			// Start consumer
			childWg.Go(func() {
				msgs, err := ch.Consume(
					mr.localQueueName, // queue
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

						select {
						case messageBuffer <- msg:
						case <-ctx.Done():
							msg.Nack(false, true)
							return
						}
					}
				}
			})

			// Start processor
			childWg.Go(func() {
				for msg := range messageBuffer {
					var appMsg InternalReplicatableDeliveryMsg
					if err := json.Unmarshal(msg.Body, &appMsg); err != nil {
						log.Printf("unmarshal error: %v", err)
						msg.Nack(false, false)
						continue
					}

					appID, err := uuid.Parse(appMsg.AppID.String())
					if err != nil {
						log.Printf("invalid app ID: %v", err)
						msg.Nack(false, false)
						continue
					}

					key := fanoutClient.FormatChannelKey(fanout.ChannelKey{
						AppID:     appID,
						ChannelID: appMsg.ChannelID,
					})

					if err := fanoutClient.Publish(ctx, key, appMsg.Payload); err != nil {
						log.Printf("valkey publish error: %v", err)
						msg.Nack(false, true)
						continue
					}

					msg.Ack(false)
				}
			})
			childWg.Wait()
		}(consumerChannel)
	}
	loopWg.Wait()
}
