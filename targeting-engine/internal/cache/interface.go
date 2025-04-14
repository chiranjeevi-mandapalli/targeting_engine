package cache

import (
	"context"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)

	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error

	Delete(ctx context.Context, key string) error

	MGet(ctx context.Context, keys ...string) ([]interface{}, error)

	Ping(ctx context.Context) error
}
