package cache

import (
	"context"
	"time"
)

// Cache defines the interface for caching operations
type Cache interface {
	Get(ctx context.Context, key string, dest interface{}) error
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Clear(ctx context.Context) error
}

// InMemoryCache is a simple in-memory cache implementation
// TODO: Replace with Redis implementation for production
type InMemoryCache struct {
	data map[string]cacheEntry
}

type cacheEntry struct {
	value     interface{}
	expiresAt time.Time
}

// NewInMemoryCache creates a new in-memory cache instance
func NewInMemoryCache() Cache {
	return &InMemoryCache{
		data: make(map[string]cacheEntry),
	}
}

// Get retrieves a value from the cache
func (c *InMemoryCache) Get(ctx context.Context, key string, dest interface{}) error {
	entry, exists := c.data[key]
	if !exists || time.Now().After(entry.expiresAt) {
		return ErrCacheMiss
	}

	// Simple assignment - in a real implementation this would need proper deserialization
	switch v := dest.(type) {
	case *interface{}:
		*v = entry.value
	default:
		// For now, just return cache miss for complex types
		return ErrCacheMiss
	}

	return nil
}

// Set stores a value in the cache
func (c *InMemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.data[key] = cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
	return nil
}

// Delete removes a value from the cache
func (c *InMemoryCache) Delete(ctx context.Context, key string) error {
	delete(c.data, key)
	return nil
}

// Clear removes all values from the cache
func (c *InMemoryCache) Clear(ctx context.Context) error {
	c.data = make(map[string]cacheEntry)
	return nil
}

// ErrCacheMiss is returned when a key is not found in the cache
var ErrCacheMiss = &CacheError{Message: "cache miss"}

// CacheError represents a cache-related error
type CacheError struct {
	Message string
}

func (e *CacheError) Error() string {
	return e.Message
}
