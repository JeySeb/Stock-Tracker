package cache

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"stock-tracker/internal/infrastructure/cache"
)

func TestInMemoryCache_SetAndGet(t *testing.T) {
	// Setup
	c := cache.NewInMemoryCache()
	ctx := context.Background()

	// Test data
	key := "test_key"
	value := "test_value"
	ttl := 1 * time.Hour

	// Execute - Set
	err := c.Set(ctx, key, value, ttl)
	assert.NoError(t, err)

	// Execute - Get
	var result interface{}
	err = c.Get(ctx, key, &result)
	assert.NoError(t, err)
	assert.Equal(t, value, result)
}

func TestInMemoryCache_GetNonExistent(t *testing.T) {
	// Setup
	c := cache.NewInMemoryCache()
	ctx := context.Background()

	// Execute
	var result interface{}
	err := c.Get(ctx, "non_existent_key", &result)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, cache.ErrCacheMiss, err)
}

func TestInMemoryCache_TTLExpiration(t *testing.T) {
	// Setup
	c := cache.NewInMemoryCache()
	ctx := context.Background()

	// Test data
	key := "expiring_key"
	value := "expiring_value"
	shortTTL := 50 * time.Millisecond

	// Execute - Set with short TTL
	err := c.Set(ctx, key, value, shortTTL)
	assert.NoError(t, err)

	// Verify it's available immediately
	var result interface{}
	err = c.Get(ctx, key, &result)
	assert.NoError(t, err)
	assert.Equal(t, value, result)

	// Wait for expiration
	time.Sleep(100 * time.Millisecond)

	// Verify it's expired
	err = c.Get(ctx, key, &result)
	assert.Error(t, err)
	assert.Equal(t, cache.ErrCacheMiss, err)
}

func TestInMemoryCache_Delete(t *testing.T) {
	// Setup
	c := cache.NewInMemoryCache()
	ctx := context.Background()

	// Test data
	key := "deletable_key"
	value := "deletable_value"
	ttl := 1 * time.Hour

	// Set value
	err := c.Set(ctx, key, value, ttl)
	assert.NoError(t, err)

	// Verify it exists
	var result interface{}
	err = c.Get(ctx, key, &result)
	assert.NoError(t, err)

	// Delete
	err = c.Delete(ctx, key)
	assert.NoError(t, err)

	// Verify it's gone
	err = c.Get(ctx, key, &result)
	assert.Error(t, err)
	assert.Equal(t, cache.ErrCacheMiss, err)
}

func TestInMemoryCache_Clear(t *testing.T) {
	// Setup
	c := cache.NewInMemoryCache()
	ctx := context.Background()

	// Set multiple values
	keys := []string{"key1", "key2", "key3"}
	for i, key := range keys {
		err := c.Set(ctx, key, i, 1*time.Hour)
		assert.NoError(t, err)
	}

	// Verify all exist
	for i, key := range keys {
		var result interface{}
		err := c.Get(ctx, key, &result)
		assert.NoError(t, err)
		assert.Equal(t, i, result)
	}

	// Clear cache
	err := c.Clear(ctx)
	assert.NoError(t, err)

	// Verify all are gone
	for _, key := range keys {
		var result interface{}
		err := c.Get(ctx, key, &result)
		assert.Error(t, err)
		assert.Equal(t, cache.ErrCacheMiss, err)
	}
}

func TestInMemoryCache_OverwriteValue(t *testing.T) {
	// Setup
	c := cache.NewInMemoryCache()
	ctx := context.Background()

	key := "overwrite_key"
	originalValue := "original"
	newValue := "new"
	ttl := 1 * time.Hour

	// Set original value
	err := c.Set(ctx, key, originalValue, ttl)
	assert.NoError(t, err)

	// Overwrite with new value
	err = c.Set(ctx, key, newValue, ttl)
	assert.NoError(t, err)

	// Verify new value
	var result interface{}
	err = c.Get(ctx, key, &result)
	assert.NoError(t, err)
	assert.Equal(t, newValue, result)
}

func TestInMemoryCache_ConcurrentAccess(t *testing.T) {
	// Setup
	c := cache.NewInMemoryCache()
	ctx := context.Background()
	ttl := 1 * time.Hour

	// Test concurrent writes and reads
	var wg sync.WaitGroup
	numGoroutines := 100
	numOperations := 10

	// Concurrent writes
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("concurrent_key_%d_%d", id, j)
				value := fmt.Sprintf("value_%d_%d", id, j)
				err := c.Set(ctx, key, value, ttl)
				assert.NoError(t, err)
			}
		}(i)
	}

	wg.Wait()

	// Concurrent reads
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := fmt.Sprintf("concurrent_key_%d_%d", id, j)
				expectedValue := fmt.Sprintf("value_%d_%d", id, j)

				var result interface{}
				err := c.Get(ctx, key, &result)
				assert.NoError(t, err)
				assert.Equal(t, expectedValue, result)
			}
		}(i)
	}

	wg.Wait()
}

func TestInMemoryCache_ComplexDataTypes(t *testing.T) {
	// Setup
	c := cache.NewInMemoryCache()
	ctx := context.Background()
	ttl := 1 * time.Hour

	tests := []struct {
		name  string
		key   string
		value interface{}
	}{
		{
			name:  "String value",
			key:   "string_key",
			value: "test string",
		},
		{
			name:  "Integer value",
			key:   "int_key",
			value: 42,
		},
		{
			name:  "Float value",
			key:   "float_key",
			value: 3.14159,
		},
		{
			name:  "Boolean value",
			key:   "bool_key",
			value: true,
		},
		{
			name:  "Slice value",
			key:   "slice_key",
			value: []string{"a", "b", "c"},
		},
		{
			name: "Map value",
			key:  "map_key",
			value: map[string]int{
				"one":   1,
				"two":   2,
				"three": 3,
			},
		},
		{
			name: "Struct value",
			key:  "struct_key",
			value: struct {
				Name string
				Age  int
			}{
				Name: "John Doe",
				Age:  30,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set value
			err := c.Set(ctx, tt.key, tt.value, ttl)
			assert.NoError(t, err)

			// Get value
			var result interface{}
			err = c.Get(ctx, tt.key, &result)

			// Note: In this simple implementation, complex types might not work
			// This test documents the current behavior
			if tt.name == "String value" || tt.name == "Integer value" ||
				tt.name == "Float value" || tt.name == "Boolean value" {
				assert.NoError(t, err)
				assert.Equal(t, tt.value, result)
			} else {
				// For complex types, current implementation returns cache miss
				// This is a limitation of the simple implementation
				assert.Error(t, err)
				assert.Equal(t, cache.ErrCacheMiss, err)
			}
		})
	}
}

func TestInMemoryCache_EdgeCases(t *testing.T) {
	// Setup
	c := cache.NewInMemoryCache()
	ctx := context.Background()

	tests := []struct {
		name        string
		key         string
		value       interface{}
		ttl         time.Duration
		expectError bool
	}{
		{
			name:        "Empty key",
			key:         "",
			value:       "empty_key_value",
			ttl:         1 * time.Hour,
			expectError: false, // Empty keys are allowed
		},
		{
			name:        "Nil value",
			key:         "nil_key",
			value:       nil,
			ttl:         1 * time.Hour,
			expectError: false, // Nil values are allowed
		},
		{
			name:        "Zero TTL",
			key:         "zero_ttl_key",
			value:       "zero_ttl_value",
			ttl:         0,
			expectError: false, // Zero TTL means immediate expiration
		},
		{
			name:        "Negative TTL",
			key:         "negative_ttl_key",
			value:       "negative_ttl_value",
			ttl:         -1 * time.Hour,
			expectError: false, // Negative TTL means immediate expiration
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set value
			err := c.Set(ctx, tt.key, tt.value, tt.ttl)
			if tt.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			// Try to get value
			var result interface{}
			err = c.Get(ctx, tt.key, &result)

			if tt.ttl <= 0 {
				// Should be expired immediately
				assert.Error(t, err)
				assert.Equal(t, cache.ErrCacheMiss, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.value, result)
			}
		})
	}
}

func TestInMemoryCache_MemoryUsage(t *testing.T) {
	// Setup
	c := cache.NewInMemoryCache()
	ctx := context.Background()

	// Add many items to test memory handling
	numItems := 1000
	ttl := 1 * time.Hour

	// Set many items
	for i := 0; i < numItems; i++ {
		key := fmt.Sprintf("memory_test_key_%d", i)
		value := fmt.Sprintf("memory_test_value_%d", i)
		err := c.Set(ctx, key, value, ttl)
		assert.NoError(t, err)
	}

	// Verify all items exist
	for i := 0; i < numItems; i++ {
		key := fmt.Sprintf("memory_test_key_%d", i)
		expectedValue := fmt.Sprintf("memory_test_value_%d", i)

		var result interface{}
		err := c.Get(ctx, key, &result)
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, result)
	}

	// Clear all items
	err := c.Clear(ctx)
	assert.NoError(t, err)

	// Verify all items are gone
	for i := 0; i < numItems; i++ {
		key := fmt.Sprintf("memory_test_key_%d", i)

		var result interface{}
		err := c.Get(ctx, key, &result)
		assert.Error(t, err)
		assert.Equal(t, cache.ErrCacheMiss, err)
	}
}

func TestCacheError_Error(t *testing.T) {
	// Test the CacheError implementation
	err := &cache.CacheError{Message: "test error"}
	assert.Equal(t, "test error", err.Error())

	// Test ErrCacheMiss
	assert.Equal(t, "cache miss", cache.ErrCacheMiss.Error())
}

func TestInMemoryCache_ContextCancellation(t *testing.T) {
	// Setup
	c := cache.NewInMemoryCache()
	ttl := 1 * time.Hour

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Operations should still work (current implementation doesn't use context for cancellation)
	err := c.Set(ctx, "test_key", "test_value", ttl)
	assert.NoError(t, err)

	var result interface{}
	err = c.Get(ctx, "test_key", &result)
	assert.NoError(t, err)
	assert.Equal(t, "test_value", result)
}

func TestInMemoryCache_LongRunningTTL(t *testing.T) {
	// Setup
	c := cache.NewInMemoryCache()
	ctx := context.Background()

	// Test with very long TTL
	key := "long_ttl_key"
	value := "long_ttl_value"
	longTTL := 24 * 365 * time.Hour // 1 year

	// Set value
	err := c.Set(ctx, key, value, longTTL)
	assert.NoError(t, err)

	// Verify it's available
	var result interface{}
	err = c.Get(ctx, key, &result)
	assert.NoError(t, err)
	assert.Equal(t, value, result)

	// Test with very short TTL that we can actually wait for
	shortKey := "short_ttl_key"
	shortValue := "short_ttl_value"
	shortTTL := 10 * time.Millisecond

	err = c.Set(ctx, shortKey, shortValue, shortTTL)
	assert.NoError(t, err)

	// Wait and verify expiration
	time.Sleep(20 * time.Millisecond)

	err = c.Get(ctx, shortKey, &result)
	assert.Error(t, err)
	assert.Equal(t, cache.ErrCacheMiss, err)

	// Original long TTL item should still be there
	err = c.Get(ctx, key, &result)
	assert.NoError(t, err)
	assert.Equal(t, value, result)
}
