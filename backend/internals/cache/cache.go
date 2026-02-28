package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/valkey-io/valkey-go"
)

// Client wraps a valkey.Client and provides typed cache operations.
type Client struct {
	v valkey.Client
}

// New creates a new cache Client backed by the given valkey.Client.
func New(client valkey.Client) *Client {
	return &Client{v: client}
}

// Get retrieves and unmarshals data from cache into the provided destination.
func Get[T any](ctx context.Context, c *Client, key CacheKey, dest *T) error {
	cmd := c.v.B().Get().Key(string(key)).Build()
	res := c.v.Do(ctx, cmd)

	if err := res.Error(); err != nil {
		return fmt.Errorf("cache miss: %w", err)
	}

	byts, err := res.AsBytes()
	if err != nil {
		return fmt.Errorf("failed to read cache bytes: %w", err)
	}

	if err := json.Unmarshal(byts, dest); err != nil {
		return fmt.Errorf("failed to unmarshal cache value: %w", err)
	}

	return nil
}

// Set marshals and stores data in cache with the given TTL.
func Set[T any](ctx context.Context, c *Client, key CacheKey, value T, ttl time.Duration) error {
	byts, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value: %w", err)
	}

	cmd := c.v.B().Set().Key(string(key)).Value(string(byts)).Ex(ttl).Build()
	res := c.v.Do(ctx, cmd)

	if err := res.Error(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// Delete removes a key from cache.
func (c *Client) Delete(ctx context.Context, key string) error {
	cmd := c.v.B().Del().Key(key).Build()
	res := c.v.Do(ctx, cmd)

	if err := res.Error(); err != nil {
		return fmt.Errorf("failed to delete from cache: %w", err)
	}

	return nil
}
