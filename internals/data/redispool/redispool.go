package redisPool

import (
	envConsts "clove/internals/consts/env"
	"errors"
	"sync"

	"github.com/joho/godotenv"
	goRedis "github.com/redis/go-redis/v9"
)

type RedisDB int

const (
	RedisStore     RedisDB = 0
	RedisHeartbeat RedisDB = 1
	RedisFanout    RedisDB = 2
)

// Client returns the singleton Redis client for the specified RedisDB pool.
// Init() must be called before using Client(), or the returned client will be nil.
// Panics if an invalid pool type is provided.
var redisStoreConn *goRedis.Client
var redisStoreConnOnce = sync.Once{}
var redisFanoutConn *goRedis.Client
var redisFanoutConnOnce = sync.Once{}
var redisHeartbeatConn *goRedis.Client
var redisHeartbeatConnOnce = sync.Once{}

func Init() {
	wg := sync.WaitGroup{}
	wg.Go(func() {
		godotenv.Load()
		redisStoreConnOnce.Do(func() {
			opts, err := goRedis.ParseURL(envConsts.RedisStoreURL())
			if err != nil {
				panic(err)
			}
			redisStoreConn = goRedis.NewClient(opts)
		})
	})

	wg.Go(func() {

		redisFanoutConnOnce.Do(func() {
			opts, err := goRedis.ParseURL(envConsts.RedisFanoutURL())
			if err != nil {
				panic(err)
			}
			redisFanoutConn = goRedis.NewClient(opts)
		})
	})
	wg.Go(func() {

		redisHeartbeatConnOnce.Do(func() {
			opts, err := goRedis.ParseURL(envConsts.RedisHeartbeatURL())
			if err != nil {
				panic(err)
			}
			redisHeartbeatConn = goRedis.NewClient(opts)
		})
	})
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
