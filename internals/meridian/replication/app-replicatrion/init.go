package AppReplication

import (
	envConsts "clove/internals/consts/env"
	"clove/internals/data/valkeyPool"
	repository "clove/internals/services/generatedRepo"
	"fmt"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/valkey-io/valkey-go"
)

type AppReplication struct {
	conn               valkey.Client
	crossRegionWriters map[repository.Region]*kafka.Writer
	localRegion        repository.Region
	localKafkaWriter   *kafka.Writer
	localReaders       []*kafka.Reader
	region             repository.Region
	meridianInitOnce   sync.Once
}

var replicationOnce = sync.Once{}
var replication *AppReplication

func ReplicateApp() *AppReplication {
	replicationOnce.Do(func() {
		readers := make([]*kafka.Reader, 5)
		for i := range envConsts.KafkaNumReaders() {
			readers[i] = kafka.NewReader(kafka.ReaderConfig{
				Brokers:        []string{envConsts.KafkaBootstrap()},
				Topic:          fmt.Sprintf("%s-app-replication", envConsts.Region()),
				GroupID:        fmt.Sprintf("%s-app-replication-group", envConsts.Region()),
				QueueCapacity:  envConsts.KafkaReaderBufferSize(),
				CommitInterval: time.Duration(envConsts.KafkaCommitInterval()) * time.Second,
			})
		}
		replication = &AppReplication{
			localReaders: readers,
			conn:         valkeyPool.Client(valkeyPool.RedisStore),
			localKafkaWriter: &kafka.Writer{
				Addr:                   kafka.TCP(envConsts.KafkaBootstrap()),
				Topic:                  fmt.Sprintf("%s-app-replication", envConsts.Region()),
				Balancer:               &kafka.RoundRobin{},
				MaxAttempts:            3,
				WriteTimeout:           10 * time.Second,
				AllowAutoTopicCreation: true,
				RequiredAcks:           kafka.RequireOne,
				Compression:            kafka.Gzip,
			},
			region: envConsts.Region(),
			crossRegionWriters: map[repository.Region]*kafka.Writer{
				repository.RegionDk1: {
					Addr:                   kafka.TCP(envConsts.KafkaBootstrap()),
					Topic:                  fmt.Sprintf("%s-app-replication", repository.RegionDk1),
					Balancer:               &kafka.RoundRobin{},
					MaxAttempts:            3,
					WriteTimeout:           10 * time.Second,
					AllowAutoTopicCreation: true,
					RequiredAcks:           kafka.RequireOne,
					Compression:            kafka.Gzip,
				},
			},
		}
	})

	return replication
}
