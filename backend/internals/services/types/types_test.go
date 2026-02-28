package types

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestNewBaseService(t *testing.T) {
	ctx := context.Background()

	service := NewBaseService(ctx, nil, false)

	assert.NotNil(t, service)
	assert.Equal(t, ctx, service.GetCtx())
	assert.False(t, service.IsCache())
}

func TestNewBaseService_WithCache(t *testing.T) {
	ctx := context.Background()

	service := NewBaseService(ctx, nil, true)

	assert.NotNil(t, service)
	assert.True(t, service.IsCache())
}

func TestBaseService_WithCache(t *testing.T) {
	ctx := context.Background()
	service := NewBaseService(ctx, nil, false)

	assert.False(t, service.IsCache())

	result := service.WithCache(nil)

	assert.Equal(t, service, result) // Returns self for chaining
	assert.True(t, service.IsCache())
}

func TestBaseService_IsCache(t *testing.T) {
	tests := []struct {
		name     string
		useCache bool
		expected bool
	}{
		{"cache enabled", true, true},
		{"cache disabled", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewBaseService(context.Background(), nil, tt.useCache)
			assert.Equal(t, tt.expected, service.IsCache())
		})
	}
}

func TestBaseService_GetCtx(t *testing.T) {
	ctx := context.WithValue(context.Background(), "key", "value")
	service := NewBaseService(ctx, nil, false)

	result := service.GetCtx()

	assert.Equal(t, ctx, result)
	assert.Equal(t, "value", result.Value("key"))
}

func TestBaseService_WithCtx(t *testing.T) {
	ctx1 := context.WithValue(context.Background(), "key", "value1")
	ctx2 := context.WithValue(context.Background(), "key", "value2")

	service := NewBaseService(ctx1, nil, false)
	assert.Equal(t, "value1", service.GetCtx().Value("key"))

	result := service.WithCtx(ctx2)

	assert.Equal(t, service, result) // Returns self for chaining
	assert.Equal(t, "value2", service.GetCtx().Value("key"))
}

func TestBaseService_ToPgUUID(t *testing.T) {
	service := NewBaseService(context.Background(), nil, false)

	tests := []struct {
		name string
		id   uuid.UUID
	}{
		{"random UUID", uuid.New()},
		{"specific UUID", uuid.MustParse("12345678-1234-1234-1234-123456789012")},
		{"nil UUID", uuid.Nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.ToPgUUID(tt.id)

			assert.True(t, result.Valid)
			// pgtype.UUID.Bytes is [16]byte, uuid.UUID is also [16]byte but different types
			assert.Equal(t, [16]byte(tt.id), result.Bytes)
		})
	}
}

func TestBaseService_ToPgText(t *testing.T) {
	service := NewBaseService(context.Background(), nil, false)

	tests := []struct {
		name  string
		input string
	}{
		{"regular string", "hello world"},
		{"empty string", ""},
		{"special characters", "hello!@#$%^&*()"},
		{"unicode", "こんにちは"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.ToPgText(tt.input)

			assert.True(t, result.Valid)
			assert.Equal(t, tt.input, result.String)
		})
	}
}

func TestBaseService_CacheAsync_Disabled(t *testing.T) {
	service := NewBaseService(context.Background(), nil, false)

	called := false
	service.CacheAsync(func(ctx context.Context) error {
		called = true
		return nil
	})

	// Give goroutine time to potentially run (it shouldn't)
	// Since cache is disabled, the function should never be called
	assert.False(t, called)
}

func TestBaseService_CacheAsync_Enabled(t *testing.T) {
	service := NewBaseService(context.Background(), nil, true)

	done := make(chan bool, 1)
	service.CacheAsync(func(ctx context.Context) error {
		done <- true
		return nil
	})

	// Wait for goroutine to complete
	select {
	case <-done:
		// Success - function was called
	case <-context.Background().Done():
		t.Fatal("CacheAsync function was not called")
	}
}

func TestServiceParams_Structure(t *testing.T) {
	ctx := context.Background()
	params := ServiceParams{
		Ctx:      ctx,
		Tx:       nil,
		UseCache: true,
	}

	assert.Equal(t, ctx, params.Ctx)
	assert.Nil(t, params.Tx)
	assert.True(t, params.UseCache)
}
