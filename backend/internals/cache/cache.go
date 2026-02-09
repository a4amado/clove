package cache

import (
	"fmt"

	"github.com/google/uuid"
)

type CacheKey string

func FormatAppCacheKey(id uuid.UUID) CacheKey {
	return CacheKey(fmt.Sprintf("apps:%s", id.String()))
}
func FormatKeyCacheKey(id uuid.UUID, keyValue string) CacheKey {
	return CacheKey(fmt.Sprintf("apps:%s:keys:%s", id.String(), keyValue))
}
