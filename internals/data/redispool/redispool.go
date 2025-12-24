package redisPool

import (
	envConsts "clove/internals/consts/env"
	"errors"
	"os"
	"sync"

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

var redisStoreConn *goRedis.Client
var redisStoreConnOnce = sync.Once{}
var redisFanoutConn *goRedis.Client
var redisFanoutConnOnce = sync.Once{}
var redisHeartbeatConn *goRedis.Client
var redisHeartbeatConnOnce = sync.Once{}

func Init() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		connString := os.Getenv(string(envConsts.REDIS_STORE_URL))
		redisStoreConnOnce.Do(func() {
			opts, err := goRedis.ParseURL(connString)
			if err != nil {
				panic(err)
			}
			redisStoreConn = goRedis.NewClient(opts)
		})
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		connString := os.Getenv(string(envConsts.REDIS_FANOUT_URL))
		redisFanoutConnOnce.Do(func() {
			opts, err := goRedis.ParseURL(connString)
			if err != nil {
				panic(err)
			}
			redisFanoutConn = goRedis.NewClient(opts)
		})
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		connString := os.Getenv(string(envConsts.REDIS_HEARTBEAT_URL))
		redisHeartbeatConnOnce.Do(func() {
			opts, err := goRedis.ParseURL(connString)
			if err != nil {
				panic(err)
			}
			redisHeartbeatConn = goRedis.NewClient(opts)
		})
	}()
	wg.Wait()
}
func Client(pool RedisDB) *goRedis.Client {

	switch pool {
	case RedisStore:
		return redisStoreConn
	case RedisFanout:
		return redisFanoutConn
	case RedisHeartbeat:
		return redisHeartbeatConn
	default:
		panic(errors.New("wrong redis conn pool type"))
	}

}
