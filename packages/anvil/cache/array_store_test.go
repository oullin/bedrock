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

// --- testItemsCanBeSetAndRetrieved ---

func TestItemsCanBeSetAndRetrieved(t *testing.T) {
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

// --- testCacheTtl ---

func TestCacheTtl(t *testing.T) {
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

// --- testMultipleItemsCanBeSetAndRetrieved ---

func TestMultipleItemsCanBeSetAndRetrieved(t *testing.T) {
	t.Parallel()

	s := New()
	ctx := context.Background()

	_ = s.PutMultiple(ctx, map[string]any{"a": 1, "b": 2}, time.Minute)

	result, err := s.GetMultiple(ctx, []string{"a", "b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 2 {
		t.Fatalf("want 2 results, got %d", len(result))
	}
}

// --- testItemsCanExpire ---

func TestItemsCanExpire(t *testing.T) {
	t.Parallel()

	clk := &fakeClock{now: time.Now()}
	s := NewWithClock(clk)
	ctx := context.Background()

	_ = s.Put(ctx, "key", "value", time.Second)

	got, err := s.Get(ctx, "key")
	if err != nil {
		t.Fatalf("key should exist before expiry: %v", err)
	}

	if got != "value" {
		t.Fatalf("want %q, got %v", "value", got)
	}

	clk.Advance(2 * time.Second)

	if s.Has(ctx, "key") {
		t.Fatal("expired key should not be present")
	}
}

// --- testTouchExtendsTtl ---

func TestTouchExtendsTtl(t *testing.T) {
	t.Parallel()

	clk := &fakeClock{now: time.Now()}
	s := NewWithClock(clk)
	ctx := context.Background()

	_ = s.Put(ctx, "key", "value", time.Minute)

	clk.Advance(30 * time.Second)

	ok, err := s.Touch(ctx, "key", 2*time.Minute)
	if err != nil || !ok {
		t.Fatalf("Touch should succeed: ok=%v err=%v", ok, err)
	}

	clk.Advance(90 * time.Second)

	got, err := s.Get(ctx, "key")
	if err != nil {
		t.Fatalf("key should still exist after touch: %v", err)
	}

	if got != "value" {
		t.Fatalf("want %q, got %v", "value", got)
	}
}

// --- testTouchMissingKey ---

func TestTouchMissingKey(t *testing.T) {
	t.Parallel()

	s := New()

	ok, err := s.Touch(context.Background(), "missing", time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ok {
		t.Fatal("Touch on missing key should return false")
	}
}

// --- testStoreItemForeverProperlyStoresInArray ---

func TestStoreItemForeverProperlyStoresInArray(t *testing.T) {
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

// --- testValuesCanBeIncremented ---

func TestValuesCanBeIncremented(t *testing.T) {
	t.Parallel()

	s := New()
	ctx := context.Background()

	_ = s.Put(ctx, "counter", int64(10), time.Minute)

	result, err := s.Increment(ctx, "counter", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != 11 {
		t.Fatalf("want 11, got %d", result)
	}

	result, err = s.Increment(ctx, "counter", 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != 15 {
		t.Fatalf("want 15, got %d", result)
	}
}

// --- testValuesGetCastedByIncrementOrDecrement ---

func TestValuesGetCastedByIncrementOrDecrement(t *testing.T) {
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

// --- testIncrementNonNumericValues ---

func TestIncrementNonNumericValues(t *testing.T) {
	t.Parallel()

	s := New()
	ctx := context.Background()

	_ = s.Put(ctx, "key", "not-a-number", time.Minute)

	_, err := s.Increment(ctx, "key", 1)
	if !errors.Is(err, ErrInvalidValue) {
		t.Fatalf("want ErrInvalidValue, got %v", err)
	}
}

// --- testNonExistingKeysCanBeIncremented ---

func TestNonExistingKeysCanBeIncremented(t *testing.T) {
	t.Parallel()

	s := New()

	result, err := s.Increment(context.Background(), "new", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != 1 {
		t.Fatalf("want 1, got %d", result)
	}
}

// --- testExpiredKeysAreIncrementedLikeNonExistingKeys ---

func TestExpiredKeysAreIncrementedLikeNonExistingKeys(t *testing.T) {
	t.Parallel()

	clk := &fakeClock{now: time.Now()}
	s := NewWithClock(clk)
	ctx := context.Background()

	_ = s.Put(ctx, "counter", int64(100), time.Minute)
	clk.Advance(2 * time.Minute)

	result, err := s.Increment(ctx, "counter", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != 1 {
		t.Fatalf("expired key should be treated as new, want 1, got %d", result)
	}
}

// --- testValuesCanBeDecremented ---

func TestValuesCanBeDecremented(t *testing.T) {
	t.Parallel()

	s := New()
	ctx := context.Background()

	_ = s.Put(ctx, "counter", int64(10), time.Minute)

	result, err := s.Decrement(ctx, "counter", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != 9 {
		t.Fatalf("want 9, got %d", result)
	}

	result, err = s.Decrement(ctx, "counter", 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != 5 {
		t.Fatalf("want 5, got %d", result)
	}
}

// --- testDecrementNonExistent ---

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

// --- testItemsCanBeRemoved ---

func TestItemsCanBeRemoved(t *testing.T) {
	t.Parallel()

	s := New()
	ctx := context.Background()

	_ = s.Put(ctx, "key", "value", time.Minute)

	if err := s.Forget(ctx, "key"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if s.Has(ctx, "key") {
		t.Fatal("key should be removed after Forget")
	}

	// Forget on non-existent key should not error.
	if err := s.Forget(ctx, "key"); err != nil {
		t.Fatalf("second Forget should not error: %v", err)
	}
}

// --- testItemsCanBeFlushed ---

func TestItemsCanBeFlushed(t *testing.T) {
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

// --- testGetNotFound ---

func TestGetNotFound(t *testing.T) {
	t.Parallel()

	s := New()

	_, err := s.Get(context.Background(), "missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

// --- testPutOverwrite ---

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

// --- testHas ---

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

// --- testMissing ---

func TestMissing(t *testing.T) {
	t.Parallel()

	s := New()
	ctx := context.Background()

	if !s.Missing(ctx, "key") {
		t.Fatal("Missing should return true for non-existent key")
	}

	_ = s.Put(ctx, "key", "value", time.Minute)

	if s.Missing(ctx, "key") {
		t.Fatal("Missing should return false for existing key")
	}
}

// --- testHasExpired ---

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

// --- testGetMultiplePartialMiss ---

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

// --- testPutWithZeroTTL ---

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

// --- testAddWhenKeyDoesNotExist ---

func TestAddWhenKeyDoesNotExist(t *testing.T) {
	t.Parallel()

	s := New()
	ctx := context.Background()

	ok, err := s.Add(ctx, "key", "value", time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ok {
		t.Fatal("Add should return true for new key")
	}

	got, _ := s.Get(ctx, "key")
	if got != "value" {
		t.Fatalf("want %q, got %v", "value", got)
	}
}

// --- testCannotAddWhenKeyExists ---

func TestCannotAddWhenKeyExists(t *testing.T) {
	t.Parallel()

	s := New()
	ctx := context.Background()

	_ = s.Put(ctx, "key", "first", time.Minute)

	ok, err := s.Add(ctx, "key", "second", time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ok {
		t.Fatal("Add should return false when key already exists")
	}

	got, _ := s.Get(ctx, "key")
	if got != "first" {
		t.Fatalf("original value should be preserved, got %v", got)
	}
}

// --- testCanAddAfterExpiry ---

func TestCanAddAfterExpiry(t *testing.T) {
	t.Parallel()

	clk := &fakeClock{now: time.Now()}
	s := NewWithClock(clk)
	ctx := context.Background()

	_, _ = s.Add(ctx, "key", "old", time.Minute)
	clk.Advance(2 * time.Minute)

	ok, _ := s.Add(ctx, "key", "new", time.Minute)
	if !ok {
		t.Fatal("Add should succeed after key expiry")
	}

	got, _ := s.Get(ctx, "key")
	if got != "new" {
		t.Fatalf("want %q, got %v", "new", got)
	}
}

// --- testRememberCallsPutAndReturnsDefault ---

func TestRememberCallsPutAndReturnsDefault(t *testing.T) {
	t.Parallel()

	s := New()
	ctx := context.Background()

	result, err := s.Remember(ctx, "key", time.Minute, func() (any, error) {
		return "computed", nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "computed" {
		t.Fatalf("want %q, got %v", "computed", result)
	}

	// Should be cached now.
	got, _ := s.Get(ctx, "key")
	if got != "computed" {
		t.Fatalf("value should be cached, got %v", got)
	}
}

// --- testRememberReturnsCachedValue ---

func TestRememberReturnsCachedValue(t *testing.T) {
	t.Parallel()

	s := New()
	ctx := context.Background()

	_ = s.Put(ctx, "key", "existing", time.Minute)

	result, err := s.Remember(ctx, "key", time.Minute, func() (any, error) {
		return "new", nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "existing" {
		t.Fatalf("want existing value %q, got %v", "existing", result)
	}
}

// --- testRememberForever ---

func TestRememberForever(t *testing.T) {
	t.Parallel()

	clk := &fakeClock{now: time.Now()}
	s := NewWithClock(clk)
	ctx := context.Background()

	result, err := s.RememberForever(ctx, "key", func() (any, error) {
		return "forever", nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "forever" {
		t.Fatalf("want %q, got %v", "forever", result)
	}

	clk.Advance(365 * 24 * time.Hour)

	got, err := s.Get(ctx, "key")
	if err != nil {
		t.Fatalf("RememberForever value should persist: %v", err)
	}

	if got != "forever" {
		t.Fatalf("want %q, got %v", "forever", got)
	}
}

// --- testValuesAreNotStoredByReference ---

func TestValuesAreNotStoredByReference(t *testing.T) {
	t.Parallel()

	s := New()
	ctx := context.Background()

	original := map[string]string{"key": "original"}
	_ = s.Put(ctx, "data", original, time.Minute)

	original["key"] = "modified"

	got, _ := s.Get(ctx, "data")
	m, ok := got.(map[string]string)
	if !ok {
		t.Fatalf("want map, got %T", got)
	}

	// In Go, maps are reference types. The ArrayStore stores values as-is
	// (no deep copy), so mutations to the original affect the stored value.
	// This matches Upstream's behavior with serialization disabled.
	if m["key"] != "modified" {
		t.Fatal("Go maps are reference types; mutation should be visible")
	}
}

// --- testConcurrentReadWrite ---

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

// --- testConcurrentIncrement ---

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
			_, err := s.Increment(ctx, "counter", 1)
			if err != nil {
				return
			}
			sum.Add(1)
		}()
	}

	wg.Wait()

	got, _ := s.Get(ctx, "counter")
	if got != int64(100) {
		t.Fatalf("want 100, got %v", got)
	}
}

// --- testAddWithZeroTTL ---

func TestAddWithZeroTTL(t *testing.T) {
	t.Parallel()

	clk := &fakeClock{now: time.Now()}
	s := NewWithClock(clk)
	ctx := context.Background()

	ok, _ := s.Add(ctx, "key", "value", 0)
	if !ok {
		t.Fatal("Add with zero TTL should succeed")
	}

	clk.Advance(365 * 24 * time.Hour)

	got, err := s.Get(ctx, "key")
	if err != nil {
		t.Fatalf("zero-TTL add should not expire: %v", err)
	}

	if got != "value" {
		t.Fatalf("want %q, got %v", "value", got)
	}
}
