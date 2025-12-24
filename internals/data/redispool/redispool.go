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
