package replication

import (
	redisPool "clove/internals/data/redispool"
	"sync"

	"github.com/redis/go-redis/v9"
)

type Replication struct {
	conn *redis.Client
}

// RedisFanOutConn:   redisPool.Client(redisPool.RedisFanout),
// RedisHearbeatConn: redisPool.Client(redisPool.RedisHeartbeat),
var replicationOnce = sync.Once{}
var replication *Replication

func Replicate() *Replication {
	replicationOnce.Do(func() {
		replication = &Replication{conn: redisPool.Client(redisPool.RedisStore)}
	})
	Init()
	return replication
}
