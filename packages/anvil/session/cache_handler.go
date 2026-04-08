package session

import (
	"context"
	"time"

	"github.com/bedrock/packages/anvil/cache"
)

// CacheHandler stores sessions in a cache backend. It is safe for
// concurrent use when the underlying cache store is.
type CacheHandler struct {
	cache   cache.Store
	minutes int
}

// NewCacheHandler creates a cache-backed session handler. Sessions expire
// after minutes of inactivity.
func NewCacheHandler(store cache.Store, minutes int) *CacheHandler {
	return &CacheHandler{
		cache:   store,
		minutes: minutes,
	}
}

func (h *CacheHandler) Open(_ context.Context, _ string, _ string) error { return nil }
func (h *CacheHandler) Close(_ context.Context) error                    { return nil }
func (h *CacheHandler) GC(_ context.Context, _ int) error                { return nil }

// Read returns the session data from the cache, or an empty string if the
// session does not exist.
func (h *CacheHandler) Read(ctx context.Context, id string) (string, error) {
	v, err := h.cache.Get(ctx, id)
	if err != nil {
		return "", nil
	}

	s, _ := v.(string)

	return s, nil
}

// Write stores session data in the cache with the configured TTL.
func (h *CacheHandler) Write(ctx context.Context, id string, data string) error {
	return h.cache.Put(ctx, id, data, time.Duration(h.minutes)*time.Minute)
}

// Destroy removes session data from the cache.
func (h *CacheHandler) Destroy(ctx context.Context, id string) error {
	return h.cache.Forget(ctx, id)
}

// GetCache returns the underlying cache store.
func (h *CacheHandler) GetCache() cache.Store {
	return h.cache
}
