package cache_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/cache"
)

func TestRepositoryEventDispatching(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	dispatcher := &mockEventDispatcher{}
	r := cache.NewRepositoryWithEvents(store, "array", dispatcher)
	ctx := context.Background()

	// Put dispatches WritingKey + KeyWritten.
	_ = r.Put(ctx, "k", "v", time.Minute)

	if dispatcher.count("WritingKey") != 1 {
		t.Fatal("expected 1 WritingKey event")
	}

	if dispatcher.count("KeyWritten") != 1 {
		t.Fatal("expected 1 KeyWritten event")
	}

	// Get dispatches CacheHit.
	_ = r.Get(ctx, "k", nil)

	if dispatcher.count("CacheHit") != 1 {
		t.Fatal("expected 1 CacheHit event")
	}

	// Miss dispatches CacheMissed.
	_ = r.Get(ctx, "missing", nil)

	if dispatcher.count("CacheMissed") != 1 {
		t.Fatal("expected 1 CacheMissed event")
	}

	// Forget dispatches ForgettingKey + KeyForgotten.
	_ = r.Forget(ctx, "k")

	if dispatcher.count("ForgettingKey") != 1 {
		t.Fatal("expected 1 ForgettingKey event")
	}

	if dispatcher.count("KeyForgotten") != 1 {
		t.Fatal("expected 1 KeyForgotten event")
	}

	// Flush dispatches CacheFlushing + CacheFlushed.
	_ = r.Flush(ctx)

	if dispatcher.count("CacheFlushing") != 1 {
		t.Fatal("expected 1 CacheFlushing event")
	}

	if dispatcher.count("CacheFlushed") != 1 {
		t.Fatal("expected 1 CacheFlushed event")
	}
}

func TestRepositoryNoEventsWithoutDispatcher(t *testing.T) {
	t.Parallel()

	r := cache.NewRepository(cache.NewArrayStore())
	ctx := context.Background()

	// Should not panic.
	_ = r.Put(ctx, "k", "v", time.Minute)
	_ = r.Get(ctx, "k", nil)
	_ = r.Forget(ctx, "k")
	_ = r.Flush(ctx)
}

func TestRepositoryStringTypedGetter(t *testing.T) {
	t.Parallel()

	r := cache.NewRepository(cache.NewArrayStore())
	ctx := context.Background()

	_ = r.Put(ctx, "name", "Alice", time.Minute)

	v, err := r.String(ctx, "name")

	if err != nil {
		t.Fatal(err)
	}

	if v != "Alice" {
		t.Fatalf("expected 'Alice', got %q", v)
	}
}

func TestRepositoryStringMissing(t *testing.T) {
	t.Parallel()

	r := cache.NewRepository(cache.NewArrayStore())
	ctx := context.Background()

	_, err := r.String(ctx, "missing")

	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestRepositoryStringInvalidType(t *testing.T) {
	t.Parallel()

	r := cache.NewRepository(cache.NewArrayStore())
	ctx := context.Background()

	_ = r.Put(ctx, "k", map[string]any{"a": 1}, time.Minute)

	_, err := r.String(ctx, "k")

	if !errors.Is(err, cache.ErrInvalidValue) {
		t.Fatalf("expected ErrInvalidValue, got %v", err)
	}
}

func TestRepositoryInteger(t *testing.T) {
	t.Parallel()

	r := cache.NewRepository(cache.NewArrayStore())
	ctx := context.Background()

	_ = r.Put(ctx, "count", int64(42), time.Minute)

	v, err := r.Integer(ctx, "count")

	if err != nil {
		t.Fatal(err)
	}

	if v != 42 {
		t.Fatalf("expected 42, got %d", v)
	}
}

func TestRepositoryIntegerFromFloat(t *testing.T) {
	t.Parallel()

	r := cache.NewRepository(cache.NewArrayStore())
	ctx := context.Background()

	_ = r.Put(ctx, "k", float64(99), time.Minute)

	v, err := r.Integer(ctx, "k")

	if err != nil {
		t.Fatal(err)
	}

	if v != 99 {
		t.Fatalf("expected 99, got %d", v)
	}
}

func TestRepositoryIntegerInvalidType(t *testing.T) {
	t.Parallel()

	r := cache.NewRepository(cache.NewArrayStore())
	ctx := context.Background()

	_ = r.Put(ctx, "k", "not-a-number", time.Minute)

	_, err := r.Integer(ctx, "k")

	if !errors.Is(err, cache.ErrInvalidValue) {
		t.Fatalf("expected ErrInvalidValue, got %v", err)
	}
}

func TestRepositoryFloat(t *testing.T) {
	t.Parallel()

	r := cache.NewRepository(cache.NewArrayStore())
	ctx := context.Background()

	_ = r.Put(ctx, "pi", 3.14, time.Minute)

	v, err := r.Float(ctx, "pi")

	if err != nil {
		t.Fatal(err)
	}

	if v != 3.14 {
		t.Fatalf("expected 3.14, got %f", v)
	}
}

func TestRepositoryFloatFromInt(t *testing.T) {
	t.Parallel()

	r := cache.NewRepository(cache.NewArrayStore())
	ctx := context.Background()

	_ = r.Put(ctx, "k", int64(5), time.Minute)

	v, err := r.Float(ctx, "k")

	if err != nil {
		t.Fatal(err)
	}

	if v != 5.0 {
		t.Fatalf("expected 5.0, got %f", v)
	}
}

func TestRepositoryBoolean(t *testing.T) {
	t.Parallel()

	r := cache.NewRepository(cache.NewArrayStore())
	ctx := context.Background()

	_ = r.Put(ctx, "flag", true, time.Minute)

	v, err := r.Boolean(ctx, "flag")

	if err != nil {
		t.Fatal(err)
	}

	if !v {
		t.Fatal("expected true")
	}
}

func TestRepositoryBooleanInvalidType(t *testing.T) {
	t.Parallel()

	r := cache.NewRepository(cache.NewArrayStore())
	ctx := context.Background()

	_ = r.Put(ctx, "k", "yes", time.Minute)

	_, err := r.Boolean(ctx, "k")

	if !errors.Is(err, cache.ErrInvalidValue) {
		t.Fatalf("expected ErrInvalidValue, got %v", err)
	}
}

func TestRepositoryMap(t *testing.T) {
	t.Parallel()

	r := cache.NewRepository(cache.NewArrayStore())
	ctx := context.Background()

	data := map[string]any{"name": "Alice", "age": 30}
	_ = r.Put(ctx, "user", data, time.Minute)

	v, err := r.Map(ctx, "user")

	if err != nil {
		t.Fatal(err)
	}

	if v["name"] != "Alice" {
		t.Fatalf("expected 'Alice', got %v", v["name"])
	}
}

func TestRepositoryMapInvalidType(t *testing.T) {
	t.Parallel()

	r := cache.NewRepository(cache.NewArrayStore())
	ctx := context.Background()

	_ = r.Put(ctx, "k", "not-a-map", time.Minute)

	_, err := r.Map(ctx, "k")

	if !errors.Is(err, cache.ErrInvalidValue) {
		t.Fatalf("expected ErrInvalidValue, got %v", err)
	}
}

func TestRepositoryGetManyEvents(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	dispatcher := &mockEventDispatcher{}
	r := cache.NewRepositoryWithEvents(store, "array", dispatcher)
	ctx := context.Background()

	_ = r.Put(ctx, "a", 1, time.Minute)

	_, _ = r.GetMany(ctx, []string{"a", "b"})

	if dispatcher.count("CacheHit") != 1 {
		t.Fatal("expected 1 CacheHit for key 'a'")
	}

	if dispatcher.count("CacheMissed") != 1 {
		t.Fatal("expected 1 CacheMissed for key 'b'")
	}
}

func TestRepositoryForeverEvents(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	dispatcher := &mockEventDispatcher{}
	r := cache.NewRepositoryWithEvents(store, "array", dispatcher)
	ctx := context.Background()

	_ = r.Forever(ctx, "k", "v")

	if dispatcher.count("WritingKey") != 1 {
		t.Fatal("expected 1 WritingKey event")
	}

	if dispatcher.count("KeyWritten") != 1 {
		t.Fatal("expected 1 KeyWritten event")
	}
}

func TestRepositoryPutManyEvents(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	dispatcher := &mockEventDispatcher{}
	r := cache.NewRepositoryWithEvents(store, "array", dispatcher)
	ctx := context.Background()

	_ = r.PutMany(ctx, map[string]any{"a": 1, "b": 2}, time.Minute)

	if dispatcher.count("WritingKey") != 2 {
		t.Fatalf("expected 2 WritingKey events, got %d", dispatcher.count("WritingKey"))
	}

	if dispatcher.count("KeyWritten") != 2 {
		t.Fatalf("expected 2 KeyWritten events, got %d", dispatcher.count("KeyWritten"))
	}
}

func TestRepositoryIncrement(t *testing.T) {
	t.Parallel()

	r := cache.NewRepository(cache.NewArrayStore())
	ctx := context.Background()

	v, err := r.Increment(ctx, "counter", 5)

	if err != nil {
		t.Fatal(err)
	}

	if v != 5 {
		t.Fatalf("expected 5, got %d", v)
	}

	v, err = r.Decrement(ctx, "counter", 2)

	if err != nil {
		t.Fatal(err)
	}

	if v != 3 {
		t.Fatalf("expected 3, got %d", v)
	}
}

func TestRepositoryTouch(t *testing.T) {
	t.Parallel()

	clk := &fakeClock{now: time.Now()}
	store := cache.NewArrayStoreWithClock(clk)
	r := cache.NewRepository(store)
	ctx := context.Background()

	_ = r.Put(ctx, "k", "v", 10*time.Second)
	_, _ = r.Touch(ctx, "k", 30*time.Second)

	clk.Advance(15 * time.Second)

	if r.Missing(ctx, "k") {
		t.Fatal("expected key to survive after Touch")
	}
}

func TestRepositoryFlexibleFreshValue(t *testing.T) {
	t.Parallel()

	r := cache.NewRepository(cache.NewArrayStore())
	ctx := context.Background()

	calls := 0
	fn := func() (any, error) {
		calls++

		return "computed", nil
	}

	v1, _ := r.Flexible(ctx, "k", time.Minute, 5*time.Minute, fn)

	if v1 != "computed" {
		t.Fatalf("expected 'computed', got %v", v1)
	}

	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}

	// Second call within fresh window should return cached value.
	v2, _ := r.Flexible(ctx, "k", time.Minute, 5*time.Minute, fn)

	if v2 != "computed" {
		t.Fatalf("expected 'computed', got %v", v2)
	}

	// fn should not be called again within the fresh window.
	if calls != 1 {
		t.Fatalf("expected 1 call (fresh hit), got %d", calls)
	}
}

func TestRepositoryFlexibleExpiredValue(t *testing.T) {
	t.Parallel()

	clk := &fakeClock{now: time.Now()}
	store := cache.NewArrayStoreWithClock(clk)
	r := cache.NewRepository(store)
	ctx := context.Background()

	calls := 0
	fn := func() (any, error) {
		calls++

		return "v" + string(rune('0'+calls)), nil
	}

	// First call: compute and store.
	_, _ = r.Flexible(ctx, "k", time.Minute, 5*time.Minute, fn)

	// Advance past stale TTL: the key expires in the store.
	clk.Advance(10 * time.Minute)

	// Should call fn synchronously since the key is fully expired.
	v, _ := r.Flexible(ctx, "k", time.Minute, 5*time.Minute, fn)

	if calls != 2 {
		t.Fatalf("expected 2 calls after full expiry, got %d", calls)
	}

	if v != "v2" {
		t.Fatalf("expected 'v2', got %v", v)
	}
}

func TestRepositoryRememberWithEvents(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	dispatcher := &mockEventDispatcher{}
	r := cache.NewRepositoryWithEvents(store, "array", dispatcher)
	ctx := context.Background()

	// First call: miss + compute.
	_, _ = r.Remember(ctx, "k", time.Minute, func() (any, error) {
		return "val", nil
	})

	if dispatcher.count("CacheMissed") != 1 {
		t.Fatal("expected 1 CacheMissed")
	}

	// Second call: hit.
	_, _ = r.Remember(ctx, "k", time.Minute, func() (any, error) {
		return "val2", nil
	})

	if dispatcher.count("CacheHit") != 1 {
		t.Fatal("expected 1 CacheHit")
	}
}

func TestRepositoryStringFromNumeric(t *testing.T) {
	t.Parallel()

	r := cache.NewRepository(cache.NewArrayStore())
	ctx := context.Background()

	_ = r.Put(ctx, "k", int64(42), time.Minute)

	v, err := r.String(ctx, "k")

	if err != nil {
		t.Fatal(err)
	}

	if v != "42" {
		t.Fatalf("expected '42', got %q", v)
	}
}

func TestRepositoryStringFromBool(t *testing.T) {
	t.Parallel()

	r := cache.NewRepository(cache.NewArrayStore())
	ctx := context.Background()

	_ = r.Put(ctx, "k", true, time.Minute)

	v, err := r.String(ctx, "k")

	if err != nil {
		t.Fatal(err)
	}

	if v != "true" {
		t.Fatalf("expected 'true', got %q", v)
	}
}
