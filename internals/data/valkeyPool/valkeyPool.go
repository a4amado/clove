package valkeyPool

import (
	envConsts "clove/internals/consts/env"
	"errors"
	"sync"

	"github.com/joho/godotenv"
	"github.com/valkey-io/valkey-go"
)

type ValkeyDB int

const (
	ValkeyStore     ValkeyDB = 0
	ValkeyHeartbeat ValkeyDB = 1
	ValkeyFanout    ValkeyDB = 2
)

// Client returns the singleton Valkey client for the specified ValkeyDB pool.
// Init() must be called before using Client(), or the returned client will be nil.
// Panics if an invalid pool type is provided.
var valkeyStoreConn valkey.Client
var valkeyStoreConnOnce = sync.Once{}
var valkeyFanoutConn valkey.Client
var valkeyFanoutConnOnce = sync.Once{}
var valkeyHeartbeatConn valkey.Client
var valkeyHeartbeatConnOnce = sync.Once{}

func Client(pool ValkeyDB) valkey.Client {
	wg := sync.WaitGroup{}
	wg.Go(func() {
		godotenv.Load()
		valkeyStoreConnOnce.Do(func() {
			opts, err := valkey.ParseURL(envConsts.RedisStoreURL())
			if err != nil {
				panic(err)
			}
			valkeyStoreConn, err = valkey.NewClient(opts)
			if err != nil {
				panic(err)
			}
		})
	})

	wg.Go(func() {

		valkeyFanoutConnOnce.Do(func() {
			opts, err := valkey.ParseURL(envConsts.RedisFanoutURL())
			if err != nil {
				panic(err)
			}
			valkeyFanoutConn, err = valkey.NewClient(opts)
			if err != nil {
				panic(err)
			}
		})
	})
	wg.Go(func() {

		valkeyHeartbeatConnOnce.Do(func() {
			opts, err := valkey.ParseURL(envConsts.RedisHeartbeatURL())
			if err != nil {
				panic(err)
			}
			valkeyHeartbeatConn, err = valkey.NewClient(opts)
			if err != nil {
				panic(err)
			}
		})
	})
	wg.Wait()
	switch pool {
	case ValkeyStore:
		return valkeyStoreConn
	case ValkeyFanout:
		return valkeyFanoutConn
	case ValkeyHeartbeat:
		return valkeyHeartbeatConn
	default:
		panic(errors.New("wrong valkey conn pool type"))
	}

}
