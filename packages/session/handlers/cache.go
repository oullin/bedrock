package handlers

import (
	"context"
	"time"
)

// CacheStore is the minimal cache interface required by CacheBasedHandler.
type CacheStore interface {
	Get(ctx context.Context, key string) (string, error)
	Put(ctx context.Context, key, value string, ttlSeconds int) error
	Forget(ctx context.Context, key string) error
}

// CacheBasedHandler stores sessions in a cache backend (Redis, Memcached,
// DynamoDB, APC, etc.). The caller injects the concrete cache store.
type CacheBasedHandler struct {
	cache   CacheStore
	minutes int
}

// NewCacheBasedHandler creates a CacheBasedHandler. minutes is the session TTL.
func NewCacheBasedHandler(cache CacheStore, minutes int) *CacheBasedHandler {
	return &CacheBasedHandler{cache: cache, minutes: minutes}
}

func (h *CacheBasedHandler) Open(_ context.Context, _, _ string) error { return nil }

func (h *CacheBasedHandler) Close(_ context.Context) error { return nil }

func (h *CacheBasedHandler) Read(ctx context.Context, id string) (string, error) {
	val, err := h.cache.Get(ctx, h.key(id))

	if err != nil {
		return "", nil // Cache miss is not an error for sessions.
	}

	return val, nil
}

func (h *CacheBasedHandler) Write(ctx context.Context, id, data string) error {
	ttl := int((time.Duration(h.minutes) * time.Minute).Seconds())

	return h.cache.Put(ctx, h.key(id), data, ttl)
}

func (h *CacheBasedHandler) Destroy(ctx context.Context, id string) error {
	return h.cache.Forget(ctx, h.key(id))
}

func (h *CacheBasedHandler) GC(_ context.Context, _ int) error {
	// Cache drivers manage their own TTL-based expiry; no explicit GC needed.
	return nil
}

// GetCache exposes the underlying cache store for tests and integration code.
func (h *CacheBasedHandler) GetCache() CacheStore {
	return h.cache
}

func (h *CacheBasedHandler) key(id string) string {
	return "session:" + id
}
