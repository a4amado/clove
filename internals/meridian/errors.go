package meridian

import "errors"

type MeridianError error

var (
	ErrCacheMiss MeridianError = errors.New("err_cache_miss")
)
