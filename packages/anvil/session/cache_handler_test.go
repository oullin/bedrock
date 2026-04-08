package session

import (
	"context"
	"testing"

	"github.com/bedrock/packages/anvil/cache"
)

func TestCacheHandlerOpen(t *testing.T) {
	t.Parallel()

	h := NewCacheHandler(cache.New(), 10)

	if err := h.Open(context.Background(), "/tmp", "sess"); err != nil {
		t.Fatalf("Open should be a no-op: %v", err)
	}
}

func TestCacheHandlerClose(t *testing.T) {
	t.Parallel()

	h := NewCacheHandler(cache.New(), 10)

	if err := h.Close(context.Background()); err != nil {
		t.Fatalf("Close should be a no-op: %v", err)
	}
}

func TestCacheHandlerRead(t *testing.T) {
	t.Parallel()

	store := cache.New()
	h := NewCacheHandler(store, 10)
	ctx := context.Background()

	_ = store.Put(ctx, "session-id", `{"foo":"bar"}`, 0)

	data, err := h.Read(ctx, "session-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data != `{"foo":"bar"}` {
		t.Fatalf("want cached data, got %q", data)
	}
}

func TestCacheHandlerReadMissing(t *testing.T) {
	t.Parallel()

	h := NewCacheHandler(cache.New(), 10)

	data, err := h.Read(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data != "" {
		t.Fatalf("missing session should return empty string, got %q", data)
	}
}

func TestCacheHandlerWrite(t *testing.T) {
	t.Parallel()

	store := cache.New()
	h := NewCacheHandler(store, 10)
	ctx := context.Background()

	if err := h.Write(ctx, "session-id", `{"key":"value"}`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	v, err := store.Get(ctx, "session-id")
	if err != nil {
		t.Fatalf("data should be in cache: %v", err)
	}

	if v != `{"key":"value"}` {
		t.Fatalf("want written data in cache, got %q", v)
	}
}

func TestCacheHandlerDestroy(t *testing.T) {
	t.Parallel()

	store := cache.New()
	h := NewCacheHandler(store, 10)
	ctx := context.Background()

	_ = h.Write(ctx, "session-id", "data")

	if err := h.Destroy(ctx, "session-id"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if store.Has(ctx, "session-id") {
		t.Fatal("session should be removed from cache")
	}
}

func TestCacheHandlerGC(t *testing.T) {
	t.Parallel()

	h := NewCacheHandler(cache.New(), 10)

	if err := h.GC(context.Background(), 3600); err != nil {
		t.Fatalf("GC should be a no-op: %v", err)
	}
}

func TestCacheHandlerGetCache(t *testing.T) {
	t.Parallel()

	store := cache.New()
	h := NewCacheHandler(store, 10)

	if h.GetCache() != store {
		t.Fatal("GetCache should return the same cache instance")
	}
}
