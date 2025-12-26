package MessageReplication

import (
	"clove/internals/meridian/fanout"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

// StartKafkaToRedisBridge starts a Kafka consumer that continuously reads app
// replication messages and saves them to the local Redis instance.
// This function blocks until the context is cancelled.
var let = 1

func (c *MessageReplication) BridgeKafkaInternalDelevieryReplicatorToRedis(ctx context.Context) {
	buffer := make(chan kafka.Message, 1000)

	fanoutClient := fanout.Fanout()
	wg := sync.WaitGroup{}

	wg.Go(func() {
		ticker := time.NewTicker(time.Second * 5)
		for {
			<-ticker.C
			fmt.Println(let)
		}
	})

	wg.Go(func() {
		for {

			select {
			case <-ctx.Done():
				return
			default:
				msg, err := c.localReader.ReadMessage(ctx)

				let += 1
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

	// Writer goroutine
	wg.Go(func() {
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

			res := fanoutClient.Publish(ctx, key, appMsg.Payload)
			if res.Err() != nil {
				log.Printf("redis publish error: %v", res.Err())
				continue
			}

			if err := c.localReader.CommitMessages(ctx, msg); err != nil {
				log.Printf("commit error: %v", err)
			}
		}
	})

	wg.Wait()
}
