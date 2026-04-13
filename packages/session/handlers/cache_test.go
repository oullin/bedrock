package handlers_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/bedrock/packages/session/handlers"
)

// fakeCacheStore is an in-memory implementation of CacheStore for testing.
type fakeCacheStore struct {
	mu   sync.Mutex
	data map[string]string
	ttls map[string]int
}

func newFakeCacheStore() *fakeCacheStore {
	return &fakeCacheStore{
		data: make(map[string]string),
		ttls: make(map[string]int),
	}
}

func (c *fakeCacheStore) Get(_ context.Context, key string) (string, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	v, ok := c.data[key]

	if !ok {
		return "", fmt.Errorf("key not found: %s", key)
	}

	return v, nil
}

func (c *fakeCacheStore) Put(_ context.Context, key, value string, ttlSeconds int) error {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.data[key] = value
	c.ttls[key] = ttlSeconds

	return nil
}

func (c *fakeCacheStore) Forget(_ context.Context, key string) error {
	c.mu.Lock()

	defer c.mu.Unlock()

	delete(c.data, key)
	delete(c.ttls, key)

	return nil
}

func TestCacheHandlerOpen(t *testing.T) {
	h := handlers.NewCacheBasedHandler(newFakeCacheStore(), 10)

	if err := h.Open(context.Background(), "", "test"); err != nil {
		t.Errorf("Open returned error: %v", err)
	}
}

func TestCacheHandlerClose(t *testing.T) {
	h := handlers.NewCacheBasedHandler(newFakeCacheStore(), 10)

	if err := h.Close(context.Background()); err != nil {
		t.Errorf("Close returned error: %v", err)
	}
}

func TestCacheHandlerReadMissing(t *testing.T) {
	h := handlers.NewCacheBasedHandler(newFakeCacheStore(), 10)
	ctx := context.Background()

	data, err := h.Read(ctx, "nonexistent")

	if err != nil {
		t.Errorf("Read returned error: %v", err)
	}

	if data != "" {
		t.Errorf("expected empty string, got %q", data)
	}
}

func TestCacheHandlerWriteAndRead(t *testing.T) {
	cache := newFakeCacheStore()
	h := handlers.NewCacheBasedHandler(cache, 10)
	ctx := context.Background()

	payload := `{"foo":"bar"}`

	if err := h.Write(ctx, "sess1", payload); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	data, err := h.Read(ctx, "sess1")

	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if data != payload {
		t.Errorf("expected %q, got %q", payload, data)
	}

	// Verify key prefix.
	if _, ok := cache.data["session:sess1"]; !ok {
		t.Error("expected cache key to be prefixed with 'session:'")
	}
}

func TestCacheHandlerDestroy(t *testing.T) {
	h := handlers.NewCacheBasedHandler(newFakeCacheStore(), 10)
	ctx := context.Background()

	_ = h.Write(ctx, "sess1", "data")

	if err := h.Destroy(ctx, "sess1"); err != nil {
		t.Fatalf("Destroy failed: %v", err)
	}

	data, _ := h.Read(ctx, "sess1")

	if data != "" {
		t.Error("session should be empty after Destroy")
	}
}

func TestCacheHandlerGC(t *testing.T) {
	h := handlers.NewCacheBasedHandler(newFakeCacheStore(), 10)

	if err := h.GC(context.Background(), 3600); err != nil {
		t.Errorf("GC returned error: %v", err)
	}
}

func TestCacheHandlerTTL(t *testing.T) {
	cache := newFakeCacheStore()
	h := handlers.NewCacheBasedHandler(cache, 10)
	ctx := context.Background()

	_ = h.Write(ctx, "sess1", "data")

	// 10 minutes = 600 seconds.
	ttl := cache.ttls["session:sess1"]

	if ttl != 600 {
		t.Errorf("expected TTL 600, got %d", ttl)
	}
}

func TestCacheHandlerDestroyNonExistent(t *testing.T) {
	h := handlers.NewCacheBasedHandler(newFakeCacheStore(), 10)

	if err := h.Destroy(context.Background(), "nope"); err != nil {
		t.Errorf("Destroy non-existent should not error: %v", err)
	}
}
