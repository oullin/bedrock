package cache_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/cache"
)

func TestTaggedCachePutAndGet(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	ts := cache.NewTagSet(store, []string{"users"})
	tc := cache.NewTaggedCache(store, ts)
	ctx := context.Background()

	_ = tc.Put(ctx, "name", "Alice", time.Minute)

	v, err := tc.Get(ctx, "name")

	if err != nil {
		t.Fatal(err)
	}

	if v != "Alice" {
		t.Fatalf("expected 'Alice', got %v", v)
	}
}

func TestTaggedCacheMiss(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	ts := cache.NewTagSet(store, []string{"users"})
	tc := cache.NewTaggedCache(store, ts)

	_, err := tc.Get(context.Background(), "missing")

	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestTaggedCacheFlushInvalidates(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	ts := cache.NewTagSet(store, []string{"users"})
	tc := cache.NewTaggedCache(store, ts)
	ctx := context.Background()

	_ = tc.Put(ctx, "name", "Alice", time.Minute)
	_ = tc.Flush(ctx)

	_, err := tc.Get(ctx, "name")

	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatal("expected ErrNotFound after Flush")
	}
}

func TestTaggedCacheIsolation(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	ctx := context.Background()

	// Untagged key.
	_ = store.Put(ctx, "global", "data", time.Minute)

	// Tagged key.
	tc := cache.NewTaggedCache(store, cache.NewTagSet(store, []string{"scoped"}))
	_ = tc.Put(ctx, "local", "val", time.Minute)

	// Flush tagged.
	_ = tc.Flush(ctx)

	// Global key should survive.
	v, err := store.Get(ctx, "global")

	if err != nil {
		t.Fatal("global key should survive tagged flush")
	}

	if v != "data" {
		t.Fatalf("expected 'data', got %v", v)
	}
}

func TestTaggedCacheGetMany(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	ts := cache.NewTagSet(store, []string{"batch"})
	tc := cache.NewTaggedCache(store, ts)
	ctx := context.Background()

	_ = tc.Put(ctx, "a", 1, time.Minute)
	_ = tc.Put(ctx, "b", 2, time.Minute)

	got, err := tc.GetMany(ctx, []string{"a", "b", "c"})

	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestTaggedCachePutMany(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	ts := cache.NewTagSet(store, []string{"batch"})
	tc := cache.NewTaggedCache(store, ts)
	ctx := context.Background()

	_ = tc.PutMany(ctx, map[string]any{"x": 10, "y": 20}, time.Minute)

	v, _ := tc.Get(ctx, "x")

	if v != 10 {
		t.Fatalf("expected 10, got %v", v)
	}
}

func TestTaggedCacheAdd(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	ts := cache.NewTagSet(store, []string{"adds"})
	tc := cache.NewTaggedCache(store, ts)
	ctx := context.Background()

	ok, _ := tc.Add(ctx, "k", "v1", time.Minute)

	if !ok {
		t.Fatal("expected Add to succeed")
	}

	ok, _ = tc.Add(ctx, "k", "v2", time.Minute)

	if ok {
		t.Fatal("expected Add to fail for existing")
	}
}

func TestTaggedCacheForever(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	ts := cache.NewTagSet(store, []string{"forever"})
	tc := cache.NewTaggedCache(store, ts)
	ctx := context.Background()

	_ = tc.Forever(ctx, "k", "eternal")

	v, _ := tc.Get(ctx, "k")

	if v != "eternal" {
		t.Fatalf("expected 'eternal', got %v", v)
	}
}

func TestTaggedCacheIncrement(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	ts := cache.NewTagSet(store, []string{"counters"})
	tc := cache.NewTaggedCache(store, ts)
	ctx := context.Background()

	v, _ := tc.Increment(ctx, "hits", 5)

	if v != 5 {
		t.Fatalf("expected 5, got %d", v)
	}

	v, _ = tc.Decrement(ctx, "hits", 2)

	if v != 3 {
		t.Fatalf("expected 3, got %d", v)
	}
}

func TestTaggedCacheForget(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	ts := cache.NewTagSet(store, []string{"temp"})
	tc := cache.NewTaggedCache(store, ts)
	ctx := context.Background()

	_ = tc.Put(ctx, "k", "v", time.Minute)
	_ = tc.Forget(ctx, "k")

	_, err := tc.Get(ctx, "k")

	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatal("expected ErrNotFound after Forget")
	}
}

func TestTaggedCacheGetPrefix(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	ts := cache.NewTagSet(store, []string{"test"})
	tc := cache.NewTaggedCache(store, ts)

	// Should return the inner store's prefix.
	if tc.GetPrefix() != "" {
		t.Fatalf("expected empty prefix, got %q", tc.GetPrefix())
	}
}

func TestTaggedCacheMultipleTags(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	ctx := context.Background()

	tc1 := cache.NewTaggedCache(store, cache.NewTagSet(store, []string{"tag1", "tag2"}))
	tc2 := cache.NewTaggedCache(store, cache.NewTagSet(store, []string{"tag1"}))

	_ = tc1.Put(ctx, "shared", "data1", time.Minute)
	_ = tc2.Put(ctx, "shared", "data2", time.Minute)

	// Different tag sets create different namespaces.
	v1, _ := tc1.Get(ctx, "shared")
	v2, _ := tc2.Get(ctx, "shared")

	if v1 != "data1" || v2 != "data2" {
		t.Fatalf("expected separate values, got %v and %v", v1, v2)
	}
}
