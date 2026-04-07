package cache

import (
	"context"
	"time"
)

// Clock reports wall-clock time. Inject a fake for deterministic tests.
type Clock interface {
	Now() time.Time
}

// Store defines the cache store contract.
type Store interface {
	Get(ctx context.Context, key string) (any, error)
	GetMultiple(ctx context.Context, keys []string) (map[string]any, error)
	Put(ctx context.Context, key string, value any, ttl time.Duration) error
	PutMultiple(ctx context.Context, values map[string]any, ttl time.Duration) error
	Forever(ctx context.Context, key string, value any) error
	Has(ctx context.Context, key string) bool
	Increment(ctx context.Context, key string, delta int64) (int64, error)
	Decrement(ctx context.Context, key string, delta int64) (int64, error)
	Forget(ctx context.Context, key string) error
	Flush(ctx context.Context) error
}
