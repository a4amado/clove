package redisPool

import (
	envConsts "clove/internals/consts/env"
	"os"

	goRedis "github.com/redis/go-redis/v9"
)

type RedisDB int

const (
	RedisStore     RedisDB = 0
	RedisHeartbeat RedisDB = 1
	RedisFanout    RedisDB = 2
)

// Client returns a Redis client configured for the specified RedisDB pool.
// 
// It selects the connection URL from the environment based on `pool` (using
// envConsts.REDIS_STORE_URL, envConsts.REDIS_FANOUT_URL, or
// envConsts.REDIS_HEARTBEAT_URL), parses that URL into options, and returns a
// new *goRedis.Client configured with those options. The function panics if the
// connection URL cannot be parsed.
func Client(pool RedisDB) *goRedis.Client {

	connString := ""
	switch pool {
	case RedisStore:
		connString = os.Getenv(string(envConsts.REDIS_STORE_URL))
	case RedisFanout:
		connString = os.Getenv(string(envConsts.REDIS_FANOUT_URL))
	case RedisHeartbeat:
		connString = os.Getenv(string(envConsts.REDIS_HEARTBEAT_URL))
	}
	opts, err := goRedis.ParseURL(connString)
	if err != nil {
		panic(err)
	}
	return goRedis.NewClient(opts)
}