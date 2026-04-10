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

func TestConcurrencyLimiterAcquire(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	cl := cache.NewConcurrencyLimiter(store, "job", 3, time.Minute)
	ctx := context.Background()

	slot, err := cl.Acquire(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if slot == "" {
		t.Fatal("expected a slot")
	}

	_ = cl.Release(ctx, slot)
}

func TestConcurrencyLimiterMaxSlots(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	cl := cache.NewConcurrencyLimiter(store, "job", 2, time.Minute)
	ctx := context.Background()

	s1, _ := cl.Acquire(ctx)
	s2, _ := cl.Acquire(ctx)
	s3, _ := cl.Acquire(ctx)

	if s1 == "" || s2 == "" {
		t.Fatal("expected first two slots to succeed")
	}

	if s3 != "" {
		t.Fatal("expected third slot to fail (all occupied)")
	}

	_ = cl.Release(ctx, s1)

	s3, _ = cl.Acquire(ctx)

	if s3 == "" {
		t.Fatal("expected acquire after release")
	}
}

func TestConcurrencyLimiterBlock(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	cl := cache.NewConcurrencyLimiter(store, "job", 1, time.Minute)
	ctx := context.Background()

	executed := false

	err := cl.Block(ctx, time.Second, func() error {
		executed = true

		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if !executed {
		t.Fatal("expected callback to execute")
	}
}

func TestConcurrencyLimiterBlockTimeout(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	cl := cache.NewConcurrencyLimiter(store, "job", 1, time.Minute)
	ctx := context.Background()

	// Occupy the only slot.
	slot, _ := cl.Acquire(ctx)

	if slot == "" {
		t.Fatal("expected slot")
	}

	err := cl.Block(ctx, 100*time.Millisecond, func() error { return nil })

	if !errors.Is(err, cache.ErrLockTimeout) {
		t.Fatalf("expected ErrLockTimeout, got %v", err)
	}

	_ = cl.Release(ctx, slot)
}

func TestConcurrencyLimiterConcurrent(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	cl := cache.NewConcurrencyLimiter(store, "job", 3, time.Second)
	ctx := context.Background()

	var wg sync.WaitGroup

	var maxConcurrent atomic.Int32

	var current atomic.Int32

	for i := 0; i < 10; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			err := cl.Block(ctx, 5*time.Second, func() error {
				n := current.Add(1)

				for {
					old := maxConcurrent.Load()

					if n <= old || maxConcurrent.CompareAndSwap(old, n) {
						break
					}
				}

				time.Sleep(10 * time.Millisecond)
				current.Add(-1)

				return nil
			})

			if err != nil {
				t.Errorf("block failed: %v", err)
			}
		}()
	}

	wg.Wait()

	if maxConcurrent.Load() > 3 {
		t.Fatalf("max concurrent exceeded limit: %d > 3", maxConcurrent.Load())
	}
}

func TestConcurrencyLimiterContextCancel(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	cl := cache.NewConcurrencyLimiter(store, "job", 1, time.Minute)

	// Occupy slot.
	ctx := context.Background()
	slot, _ := cl.Acquire(ctx)

	if slot == "" {
		t.Fatal("expected slot")
	}

	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()

	err := cl.Block(cancelCtx, 5*time.Second, func() error { return nil })

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
