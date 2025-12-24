// Package meridian is the core routing coordination service for Clove.
//
// Meridian maintains the global routing registry through a multi-tier storage architecture:
//   - PostgreSQL: Source of truth for persistent routing state
//   - Redis: Local cache for sub-millisecond route lookups
//   - Kafka: Event stream for global state propagation
//
// Architecture:
//
//	Producer: Publishes local routing changes to Kafka for global replication
//	Consumer: Ingests routing updates from Kafka and materializes them in local Redis
//	Query Interface: Serves routing lookups from Redis with PostgreSQL fallback
//
// This design ensures eventual consistency across distributed Clove instances while
// maintaining low-latency access to routing data.

package meridian

import (
	redisPool "clove/internals/data/redispool"
	"sync"

	"github.com/redis/go-redis/v9"
)

type Meridian struct {
	RedisStoreConn    *redis.Client
	RedisFanOutConn   *redis.Client
	RedisHearbeatConn *redis.Client
}

var meridianOnce *Meridian
var once = sync.Once{}

func Client() *Meridian {
	once.Do(func() {
		meridianOnce = &Meridian{
			RedisStoreConn:    redisPool.Client(redisPool.RedisStore),
			RedisFanOutConn:   redisPool.Client(redisPool.RedisFanout),
			RedisHearbeatConn: redisPool.Client(redisPool.RedisHeartbeat),
		}
	})
	return meridianOnce
}
