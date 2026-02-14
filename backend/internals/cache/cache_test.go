package cache

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestFormatAppCacheKey(t *testing.T) {
	tests := []struct {
		name     string
		id       uuid.UUID
		expected string
	}{
		{
			name:     "formats app cache key correctly",
			id:       uuid.MustParse("12345678-1234-1234-1234-123456789012"),
			expected: "apps:12345678-1234-1234-1234-123456789012",
		},
		{
			name:     "formats another app cache key",
			id:       uuid.MustParse("abcdef00-0000-0000-0000-000000000000"),
			expected: "apps:abcdef00-0000-0000-0000-000000000000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatAppCacheKey(tt.id)
			assert.Equal(t, CacheKey(tt.expected), result)
		})
	}
}

func TestFormatKeyCacheKey(t *testing.T) {
	tests := []struct {
		name     string
		id       uuid.UUID
		keyValue string
		expected string
	}{
		{
			name:     "formats key cache key correctly",
			id:       uuid.MustParse("12345678-1234-1234-1234-123456789012"),
			keyValue: "api-key-123",
			expected: "apps:12345678-1234-1234-1234-123456789012:keys:api-key-123",
		},
		{
			name:     "handles special characters in key value",
			id:       uuid.MustParse("abcdef00-0000-0000-0000-000000000000"),
			keyValue: "key_with-special.chars",
			expected: "apps:abcdef00-0000-0000-0000-000000000000:keys:key_with-special.chars",
		},
		{
			name:     "handles empty key value",
			id:       uuid.MustParse("12345678-1234-1234-1234-123456789012"),
			keyValue: "",
			expected: "apps:12345678-1234-1234-1234-123456789012:keys:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatKeyCacheKey(tt.id, tt.keyValue)
			assert.Equal(t, CacheKey(tt.expected), result)
		})
	}
}

func TestCacheKeyType(t *testing.T) {
	key := CacheKey("test:key")

	// Verify it can be converted to string
	assert.Equal(t, "test:key", string(key))
}

func TestFormatAppCacheKey_NewUUIDs(t *testing.T) {
	// Test with random UUIDs to ensure formatting is consistent
	for i := 0; i < 5; i++ {
		id := uuid.New()
		result := FormatAppCacheKey(id)

		// Should start with "apps:"
		assert.Contains(t, string(result), "apps:")

		// Should contain the UUID string
		assert.Contains(t, string(result), id.String())
	}
}

func TestFormatKeyCacheKey_NewUUIDs(t *testing.T) {
	// Test with random UUIDs to ensure formatting is consistent
	for i := 0; i < 5; i++ {
		id := uuid.New()
		keyValue := "test-key-" + uuid.NewString()[:8]
		result := FormatKeyCacheKey(id, keyValue)

		// Should contain expected components
		assert.Contains(t, string(result), "apps:")
		assert.Contains(t, string(result), id.String())
		assert.Contains(t, string(result), ":keys:")
		assert.Contains(t, string(result), keyValue)
	}
}
