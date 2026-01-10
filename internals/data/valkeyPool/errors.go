package valkeyPool

import "errors"

type RedisError error

var (
	ErrCacheMiss RedisError = errors.New("err_cache_miss")
)
