package AppReplication

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

type AppReplication struct {
	conn               *redis.Client
	crossRegionWriters map[repository.Region]*kafka.Writer
	localRegion        repository.Region
	localKafkaWriter   *kafka.Writer
	localReader        *kafka.Reader
	meridianInitOnce   sync.Once
}

var replicationOnce = sync.Once{}
var replication *AppReplication

func ReplicateApp() *AppReplication {
	replicationOnce.Do(func() {
		godotenv.Load()
		region := os.Getenv(string(envConsts.REGION))
		kafkaBootstrap := os.Getenv(string(envConsts.KAFKA_BOOTSTRAP))
		if kafkaBootstrap == "" {
			panic("KAFKA_BOOTSTRAP environment variable not set")
		}
		currentMachineRegion := repository.Region(region)
		if !currentMachineRegion.Valid() {
			panic(fmt.Sprintf("no valid region set in environment variable, got: %q", region))
		}
		replication = &AppReplication{
			localReader: kafka.NewReader(kafka.ReaderConfig{
				Brokers: []string{kafkaBootstrap},
				Topic:   fmt.Sprintf("%s-app-replication", currentMachineRegion),
				GroupID: fmt.Sprintf("%s-app-replication-group", currentMachineRegion),
			}),
			conn: redisPool.Client(redisPool.RedisStore),
			localKafkaWriter: &kafka.Writer{
				Addr:                   kafka.TCP(kafkaBootstrap),
				Topic:                  fmt.Sprintf("%s-app-replication", repository.RegionDk1),
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
