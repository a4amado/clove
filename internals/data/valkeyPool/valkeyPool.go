package valkeyPool

import (
	envConsts "clove/internals/consts/env"
	"errors"
	"sync"

	"github.com/joho/godotenv"
	"github.com/valkey-io/valkey-go"
)

type RedisDB int

const (
	RedisStore     RedisDB = 0
	RedisHeartbeat RedisDB = 1
	RedisFanout    RedisDB = 2
)

// Client returns the singleton Valkey client for the specified RedisDB pool.
// Init() must be called before using Client(), or the returned client will be nil.
// Panics if an invalid pool type is provided.
var redisStoreConn valkey.Client
var redisStoreConnOnce = sync.Once{}
var redisFanoutConn valkey.Client
var redisFanoutConnOnce = sync.Once{}
var redisHeartbeatConn valkey.Client
var redisHeartbeatConnOnce = sync.Once{}

func Client(pool RedisDB) valkey.Client {
	wg := sync.WaitGroup{}
	wg.Go(func() {
		godotenv.Load()
		redisStoreConnOnce.Do(func() {
			opts, err := valkey.ParseURL(envConsts.RedisStoreURL())
			if err != nil {
				panic(err)
			}
			redisStoreConn, err = valkey.NewClient(opts)
			if err != nil {
				panic(err)
			}
		})
	})

	wg.Go(func() {

		redisFanoutConnOnce.Do(func() {
			opts, err := valkey.ParseURL(envConsts.RedisFanoutURL())
			if err != nil {
				panic(err)
			}
			redisFanoutConn, err = valkey.NewClient(opts)
			if err != nil {
				panic(err)
			}
		})
	})
	wg.Go(func() {

		redisHeartbeatConnOnce.Do(func() {
			opts, err := valkey.ParseURL(envConsts.RedisHeartbeatURL())
			if err != nil {
				panic(err)
			}
			redisHeartbeatConn, err = valkey.NewClient(opts)
			if err != nil {
				panic(err)
			}
		})
	})
	wg.Wait()
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
