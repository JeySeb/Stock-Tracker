package cache

import (
	"context"
	"reflect"
	"sync"
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
	mu   sync.RWMutex
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
	c.mu.RLock()
	entry, exists := c.data[key]
	c.mu.RUnlock()

	if !exists || time.Now().After(entry.expiresAt) {
		return ErrCacheMiss
	}

	// Use reflection to properly handle type assignment
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr {
		return ErrCacheMiss
	}

	destElem := destValue.Elem()

	// Handle nil values specially
	if entry.value == nil {
		if destElem.CanSet() {
			destElem.Set(reflect.Zero(destElem.Type()))
			return nil
		}
		return ErrCacheMiss
	}

	// For interface{} destinations, assign directly
	if destElem.Type() == reflect.TypeOf((*interface{})(nil)).Elem() {
		destElem.Set(reflect.ValueOf(entry.value))
		return nil
	}

	// For typed destinations, attempt type conversion
	entryValue := reflect.ValueOf(entry.value)
	if entryValue.Type().AssignableTo(destElem.Type()) {
		destElem.Set(entryValue)
		return nil
	}

	// If types don't match, return cache miss
	return ErrCacheMiss
}

// Set stores a value in the cache
func (c *InMemoryCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	c.mu.Lock()
	c.data[key] = cacheEntry{
		value:     value,
		expiresAt: time.Now().Add(ttl),
	}
	c.mu.Unlock()
	return nil
}

// Delete removes a value from the cache
func (c *InMemoryCache) Delete(ctx context.Context, key string) error {
	c.mu.Lock()
	delete(c.data, key)
	c.mu.Unlock()
	return nil
}

// Clear removes all values from the cache
func (c *InMemoryCache) Clear(ctx context.Context) error {
	c.mu.Lock()
	c.data = make(map[string]cacheEntry)
	c.mu.Unlock()
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
