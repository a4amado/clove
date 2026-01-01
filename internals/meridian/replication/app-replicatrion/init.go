package AppReplication

import (
	envConsts "clove/internals/consts/env"
	"clove/internals/data/valkeyPool"
	repository "clove/internals/services/generatedRepo"
	"fmt"
	"log"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/valkey-io/valkey-go"
)

type AppReplication struct {
	conn                 valkey.Client
	crossRegionChannels  map[repository.Region]*amqp.Channel
	localRegion          repository.Region
	localRabbitMQChannel *amqp.Channel
	localRabbitMQConn    *amqp.Connection
	localConsumers       []*amqp.Channel
	region               repository.Region
	meridianInitOnce     sync.Once
	exchangeName         string
	localQueueName       string
}

var replicationOnce = sync.Once{}
var replicationInstance *AppReplication

func ReplicateApp() *AppReplication {
	replicationOnce.Do(func() {
		rabbitMQURL := envConsts.RabbitMQURL()
		conn, err := amqp.Dial(rabbitMQURL)
		if err != nil {
			log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		}

		// Main channel for publishing
		mainChannel, err := conn.Channel()
		if err != nil {
			log.Fatalf("Failed to open channel: %v", err)
		}

		// Declare topic exchange for app replication
		exchangeName := "app-replication-exchange"
		err = mainChannel.ExchangeDeclare(
			exchangeName, // name
			"direct",     // type
			true,         // durable
			false,        // auto-deleted
			false,        // internal
			false,        // no-wait
			nil,          // arguments
		)
		if err != nil {
			log.Fatalf("Failed to declare exchange: %v", err)
		}

		// Create local queue
		localRegion := envConsts.Region()
		localQueueName := fmt.Sprintf("%s-app-replication-queue", localRegion)
		_, err = mainChannel.QueueDeclare(
			localQueueName, // name
			true,           // durable
			false,          // delete when unused
			false,          // exclusive
			false,          // no-wait
			nil,            // arguments
		)
		if err != nil {
			log.Fatalf("Failed to declare queue: %v", err)
		}

		// Bind local queue to exchange with routing key
		routingKey := fmt.Sprintf("%s.app-replication", localRegion)
		err = mainChannel.QueueBind(
			localQueueName, // queue name
			routingKey,     // routing key
			exchangeName,   // exchange
			false,
			nil,
		)
		if err != nil {
			log.Fatalf("Failed to bind queue: %v", err)
		}

		// Create consumer channels
		consumers := make([]*amqp.Channel, envConsts.RabbitMQNumReaders())
		for i := range consumers {
			ch, err := conn.Channel()
			if err != nil {
				log.Fatalf("Failed to open consumer channel: %v", err)
			}
			err = ch.Qos(
				envConsts.RabbitMQPrefetchCount(), // prefetch count
				0,                                 // prefetch size
				false,                             // global
			)
			if err != nil {
				log.Fatalf("Failed to set QoS: %v", err)
			}
			consumers[i] = ch
		}

		// Create channels for cross-region publishing
		crossRegionChannels := make(map[repository.Region]*amqp.Channel)
		for _, region := range repository.AllRegionValues() {
			if region == localRegion {
				continue
			}
			ch, err := conn.Channel()
			if err != nil {
				log.Fatalf("Failed to open cross-region channel: %v", err)
			}
			crossRegionChannels[region] = ch
		}

		replicationInstance = &AppReplication{
			localConsumers:       consumers,
			conn:                 valkeyPool.Client(valkeyPool.ValkeyStore),
			localRabbitMQChannel: mainChannel,
			localRabbitMQConn:    conn,
			region:               localRegion,
			crossRegionChannels:  crossRegionChannels,
			exchangeName:         exchangeName,
			localQueueName:       localQueueName,
		}
	})

	return replicationInstance
}
