package cache

import (
	"clove/internals/data/valkeyPool"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type ApiKey struct{}

var constApiKeys = ApiKey{}

const (
	apiKeyTTL = 24 * time.Hour
)

// Get retrieves and unmarshals data from cache into the provided destination
func Get[T any](ctx context.Context, key CacheKey, dest *T) error {
	var valkeyClient = valkeyPool.Client(valkeyPool.ValkeyStore)

	valkeyCMD := valkeyClient.
		B().
		Get().
		Key(string(key)).
		Build()

	res := valkeyClient.Do(ctx, valkeyCMD)

	if err := res.Error(); err != nil {
		return fmt.Errorf("failed to get from cache: %w", err)
	}

	byts, err := res.AsBytes()
	if err != nil {
		return fmt.Errorf("failed to read response bytes: %w", err)
	}

	if err := json.Unmarshal(byts, dest); err != nil {
		return fmt.Errorf("failed to unmarshal: %w", err)
	}

	return nil
}

// Set marshals and stores data in cache with the given TTL
func Set[T any](ctx context.Context, key CacheKey, value T, ttl time.Duration) error {
	var valkeyClient = valkeyPool.Client(valkeyPool.ValkeyStore)

	byts, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal: %w", err)
	}

	valkeyCMD := valkeyClient.
		B().
		Set().
		Key(string(key)).
		Value(string(byts)).
		Ex(ttl).
		Build()

	res := valkeyClient.Do(ctx, valkeyCMD)

	if err := res.Error(); err != nil {
		return fmt.Errorf("failed to set in cache: %w", err)
	}

	return nil
}

// Delete removes a key from cache
func Delete(ctx context.Context, key string) error {
	var valkeyClient = valkeyPool.Client(valkeyPool.ValkeyStore)

	valkeyCMD := valkeyClient.
		B().
		Del().
		Key(key).
		Build()

	res := valkeyClient.Do(ctx, valkeyCMD)

	if err := res.Error(); err != nil {
		return fmt.Errorf("failed to delete from cache: %w", err)
	}

	return nil
}
