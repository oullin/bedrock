package cache_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/bedrock/packages/cache"
)

func TestMemoizedStoreGetCachesResult(t *testing.T) {
	t.Parallel()

	spy := newSpyStore()
	s := cache.NewMemoizedStore(spy)
	ctx := context.Background()

	_ = spy.inner.Put(ctx, "k", "v", time.Minute)

	// First Get: delegates to inner.
	v1, _ := s.Get(ctx, "k")

	if v1 != "v" {
		t.Fatalf("expected 'v', got %v", v1)
	}

	// Second Get: should return memoized, not hit inner.
	v2, _ := s.Get(ctx, "k")

	if v2 != "v" {
		t.Fatalf("expected 'v', got %v", v2)
	}

	if spy.callCount("Get") != 1 {
		t.Fatalf("expected 1 inner Get call, got %d", spy.callCount("Get"))
	}
}

func TestMemoizedStoreGetMissCached(t *testing.T) {
	t.Parallel()

	spy := newSpyStore()
	s := cache.NewMemoizedStore(spy)
	ctx := context.Background()

	// First Get: miss, delegates to inner.
	_, err := s.Get(ctx, "missing")

	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}

	// Second Get: should return cached miss without hitting inner.
	_, _ = s.Get(ctx, "missing")

	if spy.callCount("Get") != 1 {
		t.Fatalf("expected 1 inner Get call for cached miss, got %d", spy.callCount("Get"))
	}
}

func TestMemoizedStorePutInvalidatesMemo(t *testing.T) {
	t.Parallel()

	spy := newSpyStore()
	s := cache.NewMemoizedStore(spy)
	ctx := context.Background()

	_ = s.Put(ctx, "k", "v1", time.Minute)

	v, _ := s.Get(ctx, "k")

	if v != "v1" {
		t.Fatalf("expected 'v1', got %v", v)
	}

	// Put updates memo.
	_ = s.Put(ctx, "k", "v2", time.Minute)

	v, _ = s.Get(ctx, "k")

	if v != "v2" {
		t.Fatalf("expected 'v2' after Put, got %v", v)
	}
}

func TestMemoizedStoreForgetInvalidatesMemo(t *testing.T) {
	t.Parallel()

	spy := newSpyStore()
	s := cache.NewMemoizedStore(spy)
	ctx := context.Background()

	_ = s.Put(ctx, "k", "v", time.Minute)
	_, _ = s.Get(ctx, "k")

	_ = s.Forget(ctx, "k")

	// After Forget, memo should be cleared, so next Get hits inner.
	_, err := s.Get(ctx, "k")

	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after Forget, got %v", err)
	}
}

func TestMemoizedStoreFlushClearsMemo(t *testing.T) {
	t.Parallel()

	spy := newSpyStore()
	s := cache.NewMemoizedStore(spy)
	ctx := context.Background()

	_ = s.Put(ctx, "a", 1, 0)
	_ = s.Put(ctx, "b", 2, 0)
	_, _ = s.Get(ctx, "a")

	_ = s.Flush(ctx)

	_, err := s.Get(ctx, "a")

	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatal("expected ErrNotFound after Flush")
	}
}

func TestMemoizedStoreClear(t *testing.T) {
	t.Parallel()

	spy := newSpyStore()
	s := cache.NewMemoizedStore(spy)
	ctx := context.Background()

	_ = spy.inner.Put(ctx, "k", "v", time.Minute)
	_, _ = s.Get(ctx, "k")

	s.Clear()

	// After Clear, next Get should hit inner again.
	_, _ = s.Get(ctx, "k")

	if spy.callCount("Get") != 2 {
		t.Fatalf("expected 2 inner Get calls after Clear, got %d", spy.callCount("Get"))
	}
}

func TestMemoizedStoreIncrementUpdatesMemo(t *testing.T) {
	t.Parallel()

	spy := newSpyStore()
	s := cache.NewMemoizedStore(spy)
	ctx := context.Background()

	v, _ := s.Increment(ctx, "counter", 5)

	if v != 5 {
		t.Fatalf("expected 5, got %d", v)
	}

	// Get should return memoized value.
	got, _ := s.Get(ctx, "counter")

	if got != int64(5) {
		t.Fatalf("expected int64(5), got %v (%T)", got, got)
	}

	// Should not hit inner store for Get.
	if spy.callCount("Get") != 0 {
		t.Fatalf("expected 0 inner Get calls, got %d", spy.callCount("Get"))
	}
}

func TestMemoizedStoreGetMany(t *testing.T) {
	t.Parallel()

	spy := newSpyStore()
	s := cache.NewMemoizedStore(spy)
	ctx := context.Background()

	_ = spy.inner.Put(ctx, "a", 1, 0)
	_ = spy.inner.Put(ctx, "b", 2, 0)

	// Prime memo for "a".
	_, _ = s.Get(ctx, "a")

	// GetMany should use memo for "a" and fetch "b" from inner.
	got, _ := s.GetMany(ctx, []string{"a", "b"})

	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}

	if spy.callCount("GetMany") != 1 {
		t.Fatalf("expected 1 inner GetMany call, got %d", spy.callCount("GetMany"))
	}
}

func TestMemoizedStoreConcurrentAccess(t *testing.T) {
	t.Parallel()

	spy := newSpyStore()
	s := cache.NewMemoizedStore(spy)
	ctx := context.Background()

	_ = spy.inner.Put(ctx, "k", "v", time.Minute)

	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			s.Get(ctx, "k") //nolint:errcheck
		}()
	}

	wg.Wait()
}

func TestMemoizedStoreGetPrefix(t *testing.T) {
	t.Parallel()

	spy := newSpyStore()
	s := cache.NewMemoizedStore(spy)

	if s.GetPrefix() != "" {
		t.Fatalf("expected empty prefix, got %q", s.GetPrefix())
	}
}

func TestMemoizedStoreInner(t *testing.T) {
	t.Parallel()

	spy := newSpyStore()
	s := cache.NewMemoizedStore(spy)

	if s.Inner() != spy {
		t.Fatal("expected Inner to return the spy store")
	}
}
