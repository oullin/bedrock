package cache_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bedrock/packages/cache"
)

type fakeClock struct {
	mu  sync.Mutex
	now time.Time
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.now
}

func (c *fakeClock) Advance(d time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.now = c.now.Add(d)
}

func TestArrayStoreBasicGet(t *testing.T) {
	t.Parallel()

	s := cache.NewArrayStore()
	ctx := context.Background()

	_ = s.Put(ctx, "key", "value", time.Minute)

	got, err := s.Get(ctx, "key")
	if err != nil {
		t.Fatal(err)
	}

	if got != "value" {
		t.Fatalf("expected 'value', got %v", got)
	}
}

func TestArrayStoreMissingKey(t *testing.T) {
	t.Parallel()

	s := cache.NewArrayStore()
	_, err := s.Get(context.Background(), "missing")

	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestArrayStoreExpiry(t *testing.T) {
	t.Parallel()

	clk := &fakeClock{now: time.Now()}
	s := cache.NewArrayStoreWithClock(clk)
	ctx := context.Background()

	_ = s.Put(ctx, "key", "val", 10*time.Second)
	clk.Advance(11 * time.Second)

	_, err := s.Get(ctx, "key")
	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after expiry, got %v", err)
	}
}

func TestArrayStoreNoExpiry(t *testing.T) {
	t.Parallel()

	clk := &fakeClock{now: time.Now()}
	s := cache.NewArrayStoreWithClock(clk)
	ctx := context.Background()

	_ = s.Put(ctx, "key", "val", 0)
	clk.Advance(365 * 24 * time.Hour)

	_, err := s.Get(ctx, "key")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestArrayStoreAdd(t *testing.T) {
	t.Parallel()

	s := cache.NewArrayStore()
	ctx := context.Background()

	ok, _ := s.Add(ctx, "k", "v1", time.Minute)
	if !ok {
		t.Fatal("expected Add to succeed")
	}

	ok, _ = s.Add(ctx, "k", "v2", time.Minute)
	if ok {
		t.Fatal("expected Add to fail on existing key")
	}

	got, _ := s.Get(ctx, "k")
	if got != "v1" {
		t.Fatalf("expected original value, got %v", got)
	}
}

func TestArrayStoreIncrement(t *testing.T) {
	t.Parallel()

	s := cache.NewArrayStore()
	ctx := context.Background()

	v, _ := s.Increment(ctx, "counter", 1)
	if v != 1 {
		t.Fatalf("expected 1, got %d", v)
	}

	v, _ = s.Increment(ctx, "counter", 5)
	if v != 6 {
		t.Fatalf("expected 6, got %d", v)
	}

	v, _ = s.Decrement(ctx, "counter", 2)
	if v != 4 {
		t.Fatalf("expected 4, got %d", v)
	}
}

func TestArrayStoreForget(t *testing.T) {
	t.Parallel()

	s := cache.NewArrayStore()
	ctx := context.Background()

	_ = s.Put(ctx, "key", "val", 0)
	_ = s.Forget(ctx, "key")

	if _, err := s.Get(ctx, "key"); !errors.Is(err, cache.ErrNotFound) {
		t.Fatal("expected ErrNotFound after Forget")
	}
}

func TestArrayStoreFlush(t *testing.T) {
	t.Parallel()

	s := cache.NewArrayStore()
	ctx := context.Background()

	_ = s.Put(ctx, "a", 1, 0)
	_ = s.Put(ctx, "b", 2, 0)
	_ = s.Flush(ctx)

	for _, k := range []string{"a", "b"} {
		if _, err := s.Get(ctx, k); !errors.Is(err, cache.ErrNotFound) {
			t.Fatalf("expected ErrNotFound for %q after Flush", k)
		}
	}
}

func TestArrayStoreConcurrentAccess(t *testing.T) {
	t.Parallel()

	s := cache.NewArrayStore()
	ctx := context.Background()

	var wg sync.WaitGroup
	var failures atomic.Int32

	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			k := "key"
			_ = s.Put(ctx, k, i, 0)

			if _, err := s.Get(ctx, k); err != nil {
				failures.Add(1)
			}
		}(i)
	}

	wg.Wait()

	if failures.Load() > 0 {
		t.Fatalf("concurrent access failures: %d", failures.Load())
	}
}

func TestArrayStoreLock(t *testing.T) {
	t.Parallel()

	s := cache.NewArrayStore()
	ctx := context.Background()

	l1 := s.Lock("resource", "owner-1", time.Minute)
	l2 := s.Lock("resource", "owner-2", time.Minute)

	ok, _ := l1.Acquire(ctx)
	if !ok {
		t.Fatal("expected owner-1 to acquire lock")
	}

	ok, _ = l2.Acquire(ctx)
	if ok {
		t.Fatal("expected owner-2 to fail acquiring held lock")
	}

	l1.Release(ctx) //nolint:errcheck

	ok, _ = l2.Acquire(ctx)
	if !ok {
		t.Fatal("expected owner-2 to acquire after release")
	}
}

func TestArrayStoreGetMany(t *testing.T) {
	t.Parallel()

	s := cache.NewArrayStore()
	ctx := context.Background()

	_ = s.Put(ctx, "a", 1, 0)
	_ = s.Put(ctx, "b", 2, 0)

	got, err := s.GetMany(ctx, []string{"a", "b", "missing"})
	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(got))
	}
}
