package cache

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CacheKey string

const (
	TTLApp        = 10 * time.Minute
	TTLRegions    = 5 * time.Minute
	TTLKey        = 24 * time.Hour
	TTLCredential = 1 * time.Minute
)

func FormatAppCacheKey(id uuid.UUID) CacheKey {
	return CacheKey(fmt.Sprintf("apps:%s", id.String()))
}

func FormatAppRegionsCacheKey(appID uuid.UUID) CacheKey {
	return CacheKey(fmt.Sprintf("apps:%s:regions", appID.String()))
}

func FormatKeyCacheKey(appID uuid.UUID, keyValue string) CacheKey {
	return CacheKey(fmt.Sprintf("apps:%s:keys:%s", appID.String(), keyValue))
}

func FormatCredentialCacheKey(token string) CacheKey {
	return CacheKey(fmt.Sprintf("credentials:%s", token))
}
