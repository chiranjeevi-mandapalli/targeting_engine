package cache

import (
	"context"
	"time"
)

// Cache defines the fundamental cache operations used by all domain-specific caches
type Cache interface {
	// Get retrieves a value from cache
	Get(ctx context.Context, key string) ([]byte, error)

	// Set stores a value in cache with expiration
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error

	// Delete removes a value from cache
	Delete(ctx context.Context, key string) error

	// MGet retrieves multiple values in one call (optional)
	MGet(ctx context.Context, keys ...string) ([]interface{}, error)

	// Ping checks cache connectivity (optional)
	Ping(ctx context.Context) error
}
