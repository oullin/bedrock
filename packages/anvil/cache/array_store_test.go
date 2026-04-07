package cache

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
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

func TestPutAndGet(t *testing.T) {
	t.Parallel()

	s := New()
	ctx := context.Background()

	if err := s.Put(ctx, "key", "value", time.Minute); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := s.Get(ctx, "key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "value" {
		t.Fatalf("want %q, got %v", "value", got)
	}
}

func TestGetNotFound(t *testing.T) {
	t.Parallel()

	s := New()

	_, err := s.Get(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

func TestGetExpired(t *testing.T) {
	t.Parallel()

	clk := &fakeClock{now: time.Now()}
	s := NewWithClock(clk)
	ctx := context.Background()

	_ = s.Put(ctx, "key", "value", time.Minute)

	clk.Advance(2 * time.Minute)

	_, err := s.Get(ctx, "key")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("want ErrNotFound for expired key, got %v", err)
	}
}

func TestPutOverwrite(t *testing.T) {
	t.Parallel()

	s := New()
	ctx := context.Background()

	_ = s.Put(ctx, "key", "first", time.Minute)
	_ = s.Put(ctx, "key", "second", time.Minute)

	got, _ := s.Get(ctx, "key")
	if got != "second" {
		t.Fatalf("want %q, got %v", "second", got)
	}
}

func TestForever(t *testing.T) {
	t.Parallel()

	clk := &fakeClock{now: time.Now()}
	s := NewWithClock(clk)
	ctx := context.Background()

	_ = s.Forever(ctx, "key", "value")

	clk.Advance(365 * 24 * time.Hour)

	got, err := s.Get(ctx, "key")
	if err != nil {
		t.Fatalf("Forever key should not expire: %v", err)
	}

	if got != "value" {
		t.Fatalf("want %q, got %v", "value", got)
	}
}

func TestHas(t *testing.T) {
	t.Parallel()

	s := New()
	ctx := context.Background()

	if s.Has(ctx, "key") {
		t.Fatal("should not have missing key")
	}

	_ = s.Put(ctx, "key", "value", time.Minute)

	if !s.Has(ctx, "key") {
		t.Fatal("should have existing key")
	}
}

func TestHasExpired(t *testing.T) {
	t.Parallel()

	clk := &fakeClock{now: time.Now()}
	s := NewWithClock(clk)
	ctx := context.Background()

	_ = s.Put(ctx, "key", "value", time.Minute)
	clk.Advance(2 * time.Minute)

	if s.Has(ctx, "key") {
		t.Fatal("expired key should not be reported as present")
	}
}

func TestForgetKey(t *testing.T) {
	t.Parallel()

	s := New()
	ctx := context.Background()

	_ = s.Put(ctx, "key", "value", time.Minute)
	_ = s.Forget(ctx, "key")

	if s.Has(ctx, "key") {
		t.Fatal("key should be removed after Forget")
	}
}

func TestFlushAll(t *testing.T) {
	t.Parallel()

	s := New()
	ctx := context.Background()

	_ = s.Put(ctx, "a", 1, time.Minute)
	_ = s.Put(ctx, "b", 2, time.Minute)
	_ = s.Flush(ctx)

	if s.Has(ctx, "a") || s.Has(ctx, "b") {
		t.Fatal("all keys should be removed after Flush")
	}
}

func TestIncrement(t *testing.T) {
	t.Parallel()

	s := New()
	ctx := context.Background()

	_ = s.Put(ctx, "counter", int64(10), time.Minute)

	result, err := s.Increment(ctx, "counter", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != 15 {
		t.Fatalf("want 15, got %d", result)
	}
}

func TestIncrementNonExistent(t *testing.T) {
	t.Parallel()

	s := New()

	result, err := s.Increment(context.Background(), "new", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != 3 {
		t.Fatalf("want 3, got %d", result)
	}
}

func TestIncrementNonNumeric(t *testing.T) {
	t.Parallel()

	s := New()
	ctx := context.Background()

	_ = s.Put(ctx, "key", "not-a-number", time.Minute)

	_, err := s.Increment(ctx, "key", 1)
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("want ErrInvalidValue, got %v", err)
	}
}

func TestDecrement(t *testing.T) {
	t.Parallel()

	s := New()
	ctx := context.Background()

	_ = s.Put(ctx, "counter", int64(10), time.Minute)

	result, err := s.Decrement(ctx, "counter", 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != 7 {
		t.Fatalf("want 7, got %d", result)
	}
}

func TestDecrementNonExistent(t *testing.T) {
	t.Parallel()

	s := New()

	result, err := s.Decrement(context.Background(), "new", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != -5 {
		t.Fatalf("want -5, got %d", result)
	}
}

func TestGetMultiple(t *testing.T) {
	t.Parallel()

	s := New()
	ctx := context.Background()

	_ = s.Put(ctx, "a", 1, time.Minute)
	_ = s.Put(ctx, "b", 2, time.Minute)

	result, err := s.GetMultiple(ctx, []string{"a", "b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("want 2 results, got %d", len(result))
	}
}

func TestGetMultiplePartialMiss(t *testing.T) {
	t.Parallel()

	s := New()
	ctx := context.Background()

	_ = s.Put(ctx, "a", 1, time.Minute)

	result, err := s.GetMultiple(ctx, []string{"a", "missing"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 1 {
		t.Fatalf("want 1 result, got %d", len(result))
	}

	if _, ok := result["missing"]; ok {
		t.Fatal("missing key should not be in result")
	}
}

func TestPutMultiple(t *testing.T) {
	t.Parallel()

	s := New()
	ctx := context.Background()

	_ = s.PutMultiple(ctx, map[string]any{"a": 1, "b": 2}, time.Minute)

	if !s.Has(ctx, "a") || !s.Has(ctx, "b") {
		t.Fatal("both keys should exist after PutMultiple")
	}
}

func TestPutWithZeroTTL(t *testing.T) {
	t.Parallel()

	clk := &fakeClock{now: time.Now()}
	s := NewWithClock(clk)
	ctx := context.Background()

	_ = s.Put(ctx, "key", "value", 0)

	clk.Advance(365 * 24 * time.Hour)

	got, err := s.Get(ctx, "key")
	if err != nil {
		t.Fatalf("zero-TTL key should not expire: %v", err)
	}

	if got != "value" {
		t.Fatalf("want %q, got %v", "value", got)
	}
}

func TestConcurrentReadWrite(t *testing.T) {
	t.Parallel()

	s := New()
	ctx := context.Background()

	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			_ = s.Put(ctx, "key", n, time.Minute)
		}(i)
	}

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = s.Get(ctx, "key")
		}()
	}

	wg.Wait()
}

func TestConcurrentIncrement(t *testing.T) {
	t.Parallel()

	s := New()
	ctx := context.Background()

	_ = s.Put(ctx, "counter", int64(0), time.Minute)

	var wg sync.WaitGroup
	var sum atomic.Int64

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := s.Increment(ctx, "counter", 1)
			if err != nil {
				return
			}
			_ = result
			sum.Add(1)
		}()
	}

	wg.Wait()

	got, _ := s.Get(ctx, "counter")
	if got != int64(100) {
		t.Fatalf("want 100, got %v", got)
	}
}

func TestIncrementWithVariousNumericTypes(t *testing.T) {
	t.Parallel()

	s := New()
	ctx := context.Background()

	cases := []struct {
		name  string
		value any
		want  int64
	}{
		{"int", int(5), 6},
		{"int32", int32(5), 6},
		{"float64", float64(5.7), 6},
		{"uint", uint(5), 6},
	}

	for _, tc := range cases {
		_ = s.Put(ctx, tc.name, tc.value, time.Minute)
		got, err := s.Increment(ctx, tc.name, 1)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", tc.name, err)
		}
		if got != tc.want {
			t.Fatalf("%s: want %d, got %d", tc.name, tc.want, got)
		}
	}
}
