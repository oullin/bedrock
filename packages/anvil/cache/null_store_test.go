package cache

import (
	"context"
	"errors"
	"testing"
	"time"
)

// --- testItemsCanNotBeCached ---

func TestNullStoreItemsCanNotBeCached(t *testing.T) {
	t.Parallel()

	s := NewNullStore()
	ctx := context.Background()

	_ = s.Put(ctx, "foo", "bar", 10*time.Minute)

	_, err := s.Get(ctx, "foo")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}

// --- testGetMultipleReturnsMultipleNulls ---

func TestNullStoreGetMultipleReturnsEmpty(t *testing.T) {
	t.Parallel()

	s := NewNullStore()

	result, err := s.GetMultiple(context.Background(), []string{"foo", "bar"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result) != 0 {
		t.Fatalf("want empty map, got %d entries", len(result))
	}
}

// --- testIncrementAndDecrementReturnFalse ---

func TestNullStoreIncrementAndDecrementReturnZero(t *testing.T) {
	t.Parallel()

	s := NewNullStore()
	ctx := context.Background()

	inc, err := s.Increment(ctx, "foo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if inc != 0 {
		t.Fatalf("want 0, got %d", inc)
	}

	dec, err := s.Decrement(ctx, "foo", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if dec != 0 {
		t.Fatalf("want 0, got %d", dec)
	}
}

// --- testTouchReturnsFalse ---

func TestNullStoreTouchReturnsFalse(t *testing.T) {
	t.Parallel()

	ok, err := NewNullStore().Touch(context.Background(), "foo", 30*time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ok {
		t.Fatal("Touch should return false on NullStore")
	}
}

// --- testHasReturnsFalse ---

func TestNullStoreHasReturnsFalse(t *testing.T) {
	t.Parallel()

	s := NewNullStore()
	ctx := context.Background()

	_ = s.Put(ctx, "key", "value", time.Minute)

	if s.Has(ctx, "key") {
		t.Fatal("Has should always return false on NullStore")
	}
}

// --- testMissingReturnsTrue ---

func TestNullStoreMissingReturnsTrue(t *testing.T) {
	t.Parallel()

	if !NewNullStore().Missing(context.Background(), "key") {
		t.Fatal("Missing should always return true on NullStore")
	}
}

// --- testAddAlwaysSucceeds ---

func TestNullStoreAddAlwaysSucceeds(t *testing.T) {
	t.Parallel()

	s := NewNullStore()
	ctx := context.Background()

	ok, err := s.Add(ctx, "key", "value", time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ok {
		t.Fatal("Add should always return true on NullStore")
	}

	// Adding the same key again should also succeed since nothing is stored.
	ok, err = s.Add(ctx, "key", "other", time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !ok {
		t.Fatal("Add should always return true on NullStore")
	}
}

// --- testRememberAlwaysCallsCallback ---

func TestNullStoreRememberAlwaysCallsCallback(t *testing.T) {
	t.Parallel()

	s := NewNullStore()
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

	// Calling again should invoke callback again (nothing cached).
	result, err = s.Remember(ctx, "key", time.Minute, func() (any, error) {
		return "recomputed", nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "recomputed" {
		t.Fatalf("want %q, got %v", "recomputed", result)
	}
}

// --- testFlushSucceeds ---

func TestNullStoreFlushSucceeds(t *testing.T) {
	t.Parallel()

	if err := NewNullStore().Flush(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
