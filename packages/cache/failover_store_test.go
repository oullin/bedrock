package cache_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/cache"
)

func TestFailoverStoreGetFromPrimary(t *testing.T) {
	t.Parallel()

	primary := cache.NewArrayStore()
	secondary := cache.NewArrayStore()
	s := cache.NewFailoverStore(primary, secondary)
	ctx := context.Background()

	_ = primary.Put(ctx, "k", "primary", 0)
	_ = secondary.Put(ctx, "k", "secondary", 0)

	v, err := s.Get(ctx, "k")

	if err != nil {
		t.Fatal(err)
	}

	if v != "primary" {
		t.Fatalf("expected 'primary', got %v", v)
	}
}

func TestFailoverStoreGetFallback(t *testing.T) {
	t.Parallel()

	primary := cache.NewArrayStore()
	secondary := cache.NewArrayStore()
	s := cache.NewFailoverStore(primary, secondary)
	ctx := context.Background()

	// Only set in secondary.
	_ = secondary.Put(ctx, "k", "fallback", 0)

	v, err := s.Get(ctx, "k")

	if err != nil {
		t.Fatal(err)
	}

	if v != "fallback" {
		t.Fatalf("expected 'fallback', got %v", v)
	}
}

func TestFailoverStoreGetAllMiss(t *testing.T) {
	t.Parallel()

	s := cache.NewFailoverStore(cache.NewArrayStore(), cache.NewArrayStore())

	_, err := s.Get(context.Background(), "missing")

	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestFailoverStorePutWritesToAll(t *testing.T) {
	t.Parallel()

	primary := cache.NewArrayStore()
	secondary := cache.NewArrayStore()
	s := cache.NewFailoverStore(primary, secondary)
	ctx := context.Background()

	_ = s.Put(ctx, "k", "value", time.Minute)

	// Both stores should have the value.
	v1, _ := primary.Get(ctx, "k")
	v2, _ := secondary.Get(ctx, "k")

	if v1 != "value" || v2 != "value" {
		t.Fatalf("expected both to have 'value', got %v and %v", v1, v2)
	}
}

func TestFailoverStoreForget(t *testing.T) {
	t.Parallel()

	primary := cache.NewArrayStore()
	secondary := cache.NewArrayStore()
	s := cache.NewFailoverStore(primary, secondary)
	ctx := context.Background()

	_ = s.Put(ctx, "k", "v", 0)
	_ = s.Forget(ctx, "k")

	_, err1 := primary.Get(ctx, "k")
	_, err2 := secondary.Get(ctx, "k")

	if !errors.Is(err1, cache.ErrNotFound) || !errors.Is(err2, cache.ErrNotFound) {
		t.Fatal("expected both to be empty after Forget")
	}
}

func TestFailoverStoreFlush(t *testing.T) {
	t.Parallel()

	primary := cache.NewArrayStore()
	secondary := cache.NewArrayStore()
	s := cache.NewFailoverStore(primary, secondary)
	ctx := context.Background()

	_ = s.Put(ctx, "a", 1, 0)
	_ = s.Flush(ctx)

	_, err1 := primary.Get(ctx, "a")
	_, err2 := secondary.Get(ctx, "a")

	if !errors.Is(err1, cache.ErrNotFound) || !errors.Is(err2, cache.ErrNotFound) {
		t.Fatal("expected both flushed")
	}
}

func TestFailoverStoreIncrement(t *testing.T) {
	t.Parallel()

	s := cache.NewFailoverStore(cache.NewArrayStore(), cache.NewArrayStore())
	ctx := context.Background()

	v, _ := s.Increment(ctx, "k", 5)

	if v != 5 {
		t.Fatalf("expected 5, got %d", v)
	}
}

func TestFailoverStoreTouch(t *testing.T) {
	t.Parallel()

	primary := cache.NewArrayStore()
	secondary := cache.NewArrayStore()
	s := cache.NewFailoverStore(primary, secondary)
	ctx := context.Background()

	_ = s.Put(ctx, "k", "v", time.Minute)

	ok, _ := s.Touch(ctx, "k", time.Hour)

	if !ok {
		t.Fatal("expected Touch to succeed")
	}
}

func TestFailoverStoreAdd(t *testing.T) {
	t.Parallel()

	s := cache.NewFailoverStore(cache.NewArrayStore())
	ctx := context.Background()

	ok, _ := s.Add(ctx, "k", "v", time.Minute)

	if !ok {
		t.Fatal("expected Add to succeed")
	}
}

func TestFailoverStoreGetPrefix(t *testing.T) {
	t.Parallel()

	s := cache.NewFailoverStore(cache.NewArrayStore())

	if s.GetPrefix() != "" {
		t.Fatalf("expected empty prefix, got %q", s.GetPrefix())
	}
}
