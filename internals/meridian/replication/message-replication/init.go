package MessageReplication

import (
	envConsts "clove/internals/consts/env"
	redisPool "clove/internals/data/redispool"
	"clove/internals/repository"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

type MessageReplication struct {
	conn               *redis.Client
	crossRegionWriters map[repository.Region]*kafka.Writer
	localRegion        repository.Region
	localKafkaWriter   *kafka.Writer
	localReaders       []*kafka.Reader
	meridianInitOnce   sync.Once
}

var MessageReplicationOnce = sync.Once{}
var replication *MessageReplication

func ReplicateMessage() *MessageReplication {

	MessageReplicationOnce.Do(func() {
		readers := make([]*kafka.Reader, 5)
		for i := range envConsts.KafkaNumReaders() {
			readers[i] = kafka.NewReader(kafka.ReaderConfig{
				Brokers:        []string{envConsts.KafkaBootstrap()},
				Topic:          fmt.Sprintf("%s-msg-replication", envConsts.Region()),
				GroupID:        fmt.Sprintf("%s-msg-replication-group", envConsts.Region()),
				QueueCapacity:  envConsts.KafkaReaderBufferSize(),
				CommitInterval: time.Duration(envConsts.KafkaCommitInterval()) * time.Second,
			})
		}
		replication = &MessageReplication{
			conn:         redisPool.Client(redisPool.RedisFanout),
			localReaders: readers,
			crossRegionWriters: map[repository.Region]*kafka.Writer{
				repository.RegionDk1: {
					Addr:                   kafka.TCP(envConsts.KafkaBootstrap()),
					Topic:                  fmt.Sprintf("%s-msg-replication", repository.RegionDk1),
					Balancer:               &kafka.RoundRobin{},
					MaxAttempts:            3,
					WriteTimeout:           3 * time.Second,
					AllowAutoTopicCreation: true,
					RequiredAcks:           kafka.RequireOne,
					Compression:            kafka.Gzip,
				},
			},
			localRegion: envConsts.Region(),
			localKafkaWriter: &kafka.Writer{
				Addr:                   kafka.TCP(envConsts.KafkaBootstrap()),
				Topic:                  fmt.Sprintf("%s-msg-replication", repository.RegionDk1),
				Balancer:               &kafka.RoundRobin{},
				MaxAttempts:            3,
				WriteTimeout:           3 * time.Second,
				AllowAutoTopicCreation: true,
				RequiredAcks:           kafka.RequireOne,
				Compression:            kafka.Gzip,
			},
		}
	})

	return replication
}
