package cache

import (
	"clove/internals/data/valkeyPool"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type ApiKey struct{}

var constApiKeys = ApiKey{}

const (
	apiKeyTTL = 24 * time.Hour
)

func (*AppCache) Keys() *ApiKey {
	return &constApiKeys
}
func (*ApiKey) formatAppApiKey(appId uuid.UUID, keyId uuid.UUID) string {
	return fmt.Sprintf("app:%s:key:%s", appId.String(), keyId.String())
}

func (a *ApiKey) Get(ctx context.Context, appId uuid.UUID, keyId uuid.UUID) (*string, error) {
	var valkeyClient = valkeyPool.Client(valkeyPool.ValkeyStore)

	valkeyCMD := valkeyClient.
		B().
		Get().
		Key(a.formatAppApiKey(appId, keyId)).
		Build()

	res := valkeyClient.Do(ctx, valkeyCMD)

	// Handle Valkey errors
	if err := res.Error(); err != nil {
		return nil, fmt.Errorf("failed to get API key from cache: %w", err)
	}

	// Get the raw bytes from the response
	data, err := res.AsBytes()
	if err != nil {
		return nil, fmt.Errorf("failed to read cache response: %w", err)
	}

	key := string(data)
	return &key, nil
}

func (a *ApiKey) Set(ctx context.Context, appId uuid.UUID, keyId uuid.UUID, apiKey string) error {
	var valkeyClient = valkeyPool.Client(valkeyPool.ValkeyStore)

	valkeyCMD := valkeyClient.
		B().
		Set().
		Key(a.formatAppApiKey(appId, keyId)).
		Value(apiKey).
		Ex(apiKeyTTL).
		Build()

	res := valkeyClient.Do(ctx, valkeyCMD)

	if err := res.Error(); err != nil {
		return fmt.Errorf("failed to set API key in cache: %w", err)
	}

	return nil
}

func (a *ApiKey) Delete(ctx context.Context, appId uuid.UUID, keyId uuid.UUID) error {
	var valkeyClient = valkeyPool.Client(valkeyPool.ValkeyStore)

	valkeyCMD := valkeyClient.
		B().
		Del().
		Key(a.formatAppApiKey(appId, keyId)).
		Build()

	res := valkeyClient.Do(ctx, valkeyCMD)

	if err := res.Error(); err != nil {
		return fmt.Errorf("failed to delete API key from cache: %w", err)
	}

	return nil
}
