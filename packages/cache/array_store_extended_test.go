package cache_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/cache"
)

func TestArrayStorePutMany(t *testing.T) {
	t.Parallel()

	s := cache.NewArrayStore()
	ctx := context.Background()

	err := s.PutMany(ctx, map[string]any{"a": 1, "b": 2, "c": 3}, time.Minute)

	if err != nil {
		t.Fatal(err)
	}

	got, err := s.GetMany(ctx, []string{"a", "b", "c"})

	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(got))
	}
}

func TestArrayStoreTouch(t *testing.T) {
	t.Parallel()

	clk := &fakeClock{now: time.Now()}
	s := cache.NewArrayStoreWithClock(clk)
	ctx := context.Background()

	_ = s.Put(ctx, "k", "v", 10*time.Second)

	ok, err := s.Touch(ctx, "k", 30*time.Second)

	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("expected Touch to succeed")
	}

	clk.Advance(15 * time.Second)

	_, err = s.Get(ctx, "k")

	if err != nil {
		t.Fatal("expected key to survive after Touch extended TTL")
	}
}

func TestArrayStoreTouchMissing(t *testing.T) {
	t.Parallel()

	s := cache.NewArrayStore()
	ctx := context.Background()

	ok, _ := s.Touch(ctx, "missing", time.Minute)

	if ok {
		t.Fatal("expected Touch to return false for missing key")
	}
}

func TestArrayStoreForever(t *testing.T) {
	t.Parallel()

	clk := &fakeClock{now: time.Now()}
	s := cache.NewArrayStoreWithClock(clk)
	ctx := context.Background()

	_ = s.Forever(ctx, "k", "eternal")
	clk.Advance(100 * 365 * 24 * time.Hour)

	v, err := s.Get(ctx, "k")

	if err != nil {
		t.Fatal("expected Forever key to never expire")
	}

	if v != "eternal" {
		t.Fatalf("expected 'eternal', got %v", v)
	}
}

func TestArrayStoreGetPrefix(t *testing.T) {
	t.Parallel()

	s := cache.NewArrayStore()

	if s.GetPrefix() != "" {
		t.Fatalf("expected empty prefix, got %q", s.GetPrefix())
	}
}

func TestArrayStoreAddAfterExpiry(t *testing.T) {
	t.Parallel()

	clk := &fakeClock{now: time.Now()}
	s := cache.NewArrayStoreWithClock(clk)
	ctx := context.Background()

	_, _ = s.Add(ctx, "k", "v1", 5*time.Second)
	clk.Advance(10 * time.Second)

	ok, _ := s.Add(ctx, "k", "v2", time.Minute)

	if !ok {
		t.Fatal("expected Add to succeed after expiry")
	}

	v, _ := s.Get(ctx, "k")

	if v != "v2" {
		t.Fatalf("expected 'v2', got %v", v)
	}
}

func TestArrayStoreIncrementNonNumeric(t *testing.T) {
	t.Parallel()

	s := cache.NewArrayStore()
	ctx := context.Background()

	_ = s.Put(ctx, "k", "not-a-number", time.Minute)

	_, err := s.Increment(ctx, "k", 1)

	if !errors.Is(err, cache.ErrInvalidValue) {
		t.Fatalf("expected ErrInvalidValue, got %v", err)
	}
}

func TestArrayStoreDecrementBelowZero(t *testing.T) {
	t.Parallel()

	s := cache.NewArrayStore()
	ctx := context.Background()

	v, _ := s.Increment(ctx, "k", 2)

	if v != 2 {
		t.Fatalf("expected 2, got %d", v)
	}

	v, _ = s.Decrement(ctx, "k", 5)

	if v != -3 {
		t.Fatalf("expected -3, got %d", v)
	}
}

func TestArrayStoreLockExpiry(t *testing.T) {
	t.Parallel()

	clk := &fakeClock{now: time.Now()}
	s := cache.NewArrayStoreWithClock(clk)
	ctx := context.Background()

	l1 := s.Lock("res", "owner-1", 5*time.Second)
	ok, _ := l1.Acquire(ctx)

	if !ok {
		t.Fatal("expected acquire to succeed")
	}

	clk.Advance(10 * time.Second)

	l2 := s.Lock("res", "owner-2", 5*time.Second)
	ok, _ = l2.Acquire(ctx)

	if !ok {
		t.Fatal("expected acquire to succeed after lock expired")
	}
}

func TestArrayStoreLockGet(t *testing.T) {
	t.Parallel()

	s := cache.NewArrayStore()
	ctx := context.Background()
	l := s.Lock("res", "owner", time.Minute)

	executed := false

	err := l.Get(ctx, func() error {
		executed = true

		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if !executed {
		t.Fatal("expected callback to be executed")
	}
}

func TestArrayStoreLockBlocked(t *testing.T) {
	t.Parallel()

	s := cache.NewArrayStore()
	ctx := context.Background()

	l1 := s.Lock("res", "owner-1", time.Minute)
	l1.Acquire(ctx) //nolint:errcheck

	l2 := s.Lock("res", "owner-2", time.Minute)
	blocked, _ := l2.Blocked(ctx)

	if !blocked {
		t.Fatal("expected lock to report as blocked")
	}
}

func TestArrayStoreLockForceRelease(t *testing.T) {
	t.Parallel()

	s := cache.NewArrayStore()
	ctx := context.Background()

	l1 := s.Lock("res", "owner-1", time.Minute)
	l1.Acquire(ctx) //nolint:errcheck

	l2 := s.Lock("res", "owner-2", time.Minute)
	_ = l2.ForceRelease(ctx)

	ok, _ := l2.Acquire(ctx)

	if !ok {
		t.Fatal("expected acquire after force release")
	}
}

func TestArrayStoreTagsBasic(t *testing.T) {
	t.Parallel()

	s := cache.NewArrayStore()
	ctx := context.Background()

	tagged := s.Tags("users")
	_ = tagged.Put(ctx, "name", "Alice", time.Minute)

	v, err := tagged.Get(ctx, "name")

	if err != nil {
		t.Fatal(err)
	}

	if v != "Alice" {
		t.Fatalf("expected 'Alice', got %v", v)
	}
}

func TestArrayStoreTagsFlush(t *testing.T) {
	t.Parallel()

	s := cache.NewArrayStore()
	ctx := context.Background()

	tagged := s.Tags("users")
	_ = tagged.Put(ctx, "name", "Alice", time.Minute)

	_ = tagged.Flush(ctx)

	// After flush, the tagged namespace changes so old keys are inaccessible.
	_, err := tagged.Get(ctx, "name")

	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after flush, got %v", err)
	}
}

func TestArrayStoreTagsIsolation(t *testing.T) {
	t.Parallel()

	s := cache.NewArrayStore()
	ctx := context.Background()

	_ = s.Put(ctx, "global", "data", time.Minute)
	tagged := s.Tags("scoped")
	_ = tagged.Put(ctx, "local", "data", time.Minute)

	// Global key should survive tagged flush.
	_ = tagged.Flush(ctx)

	v, err := s.Get(ctx, "global")

	if err != nil {
		t.Fatal("expected global key to survive tagged flush")
	}

	if v != "data" {
		t.Fatalf("expected 'data', got %v", v)
	}
}

func TestArrayStoreConcurrentLocks(t *testing.T) {
	t.Parallel()

	s := cache.NewArrayStore()
	ctx := context.Background()

	acquired := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(id int) {
			l := s.Lock("shared", "owner-"+string(rune('0'+id)), time.Minute)
			ok, _ := l.Acquire(ctx)
			acquired <- ok
		}(i)
	}

	successes := 0

	for i := 0; i < 10; i++ {
		if <-acquired {
			successes++
		}
	}

	if successes != 1 {
		t.Fatalf("expected exactly 1 lock acquisition, got %d", successes)
	}
}

func TestArrayStoreAll(t *testing.T) {
	t.Parallel()

	clock := &fakeClock{now: time.Now()}
	s := cache.NewArrayStoreWithClock(clock)
	ctx := context.Background()

	_ = s.Put(ctx, "a", "val-a", time.Minute)
	_ = s.Put(ctx, "b", "val-b", time.Minute)
	_ = s.Put(ctx, "expired", "gone", 100*time.Millisecond)

	clock.Advance(time.Second)

	all := s.All()

	if len(all) != 2 {
		t.Fatalf("expected 2 items, got %d", len(all))
	}

	if all["a"] != "val-a" || all["b"] != "val-b" {
		t.Fatalf("unexpected values: %v", all)
	}
}

func TestArrayStoreSetPrefix(t *testing.T) {
	t.Parallel()

	s := cache.NewArrayStore()

	if s.GetPrefix() != "" {
		t.Fatal("expected empty prefix")
	}

	s.SetPrefix("test")

	if s.GetPrefix() != "test" {
		t.Fatalf("expected 'test', got %q", s.GetPrefix())
	}
}
