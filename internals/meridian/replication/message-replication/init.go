package MessageReplication

import (
	envConsts "clove/internals/consts/env"
	redisPool "clove/internals/data/redispool"
	"clove/internals/repository"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type MessageReplication struct {
	conn               *redis.Client
	crossRegionWriters map[repository.Region]*kafka.Writer
	localRegion        repository.Region
	localKafkaWriter   *kafka.Writer
	localReader        *kafka.Reader
	meridianInitOnce   sync.Once
}

var MessageReplicationOnce = sync.Once{}
var replication *MessageReplication

func ReplicateMessage() *MessageReplication {

	godotenv.Load()
	region := os.Getenv(string(envConsts.REGION))
	kafkaBootstrap := os.Getenv(string(envConsts.KAFKA_BOOTSTRAP))
	currentMachineRegion := repository.Region(region)
	if !currentMachineRegion.Valid() {
		panic(fmt.Sprintf("no valid region set in environment variable, got: %q", region))
	}
	if kafkaBootstrap == "" {

		panic("KAFKA_BOOTSTRAP environment variable not set")
	}
	MessageReplicationOnce.Do(func() {
		replication = &MessageReplication{
			conn: redisPool.Client(redisPool.RedisFanout),
			localReader: kafka.NewReader(kafka.ReaderConfig{
				Brokers:        []string{kafkaBootstrap},
				Topic:          fmt.Sprintf("%s-msg-replication", currentMachineRegion),
				GroupID:        fmt.Sprintf("%s-msg-replication-group", currentMachineRegion),
				QueueCapacity:  1000,
				CommitInterval: 10 * time.Second,
			}),
			crossRegionWriters: map[repository.Region]*kafka.Writer{
				repository.RegionDk1: {
					Addr:                   kafka.TCP(kafkaBootstrap),
					Topic:                  fmt.Sprintf("%s-msg-replication", repository.RegionDk1),
					Balancer:               &kafka.RoundRobin{},
					MaxAttempts:            3,
					WriteTimeout:           10 * time.Second,
					AllowAutoTopicCreation: true,
					RequiredAcks:           kafka.RequireOne,
					Compression:            kafka.Gzip,
				},
			},
			localRegion: repository.RegionDk1,
			localKafkaWriter: &kafka.Writer{
				Addr:                   kafka.TCP(kafkaBootstrap),
				Topic:                  fmt.Sprintf("%s-msg-replication", repository.RegionDk1),
				Balancer:               &kafka.RoundRobin{},
				MaxAttempts:            3,
				WriteTimeout:           10 * time.Second,
				AllowAutoTopicCreation: true,
				RequiredAcks:           kafka.RequireOne,
				Compression:            kafka.Gzip,
			},
		}
	})

	return replication
}
