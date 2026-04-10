package cache_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/cache"
)

func TestSessionStoreBasicGetPut(t *testing.T) {
	t.Parallel()

	session := newMockSession()
	s := cache.NewSessionStore(session, "")
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

func TestSessionStoreMissingKey(t *testing.T) {
	t.Parallel()

	session := newMockSession()
	s := cache.NewSessionStore(session, "")

	_, err := s.Get(context.Background(), "missing")

	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestSessionStoreExpiry(t *testing.T) {
	t.Parallel()

	session := newMockSession()
	clk := &fakeClock{now: time.Now()}
	s := cache.NewSessionStoreWithClock(session, "", clk)
	ctx := context.Background()

	_ = s.Put(ctx, "key", "val", 10*time.Second)
	clk.Advance(11 * time.Second)

	_, err := s.Get(ctx, "key")

	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after expiry, got %v", err)
	}
}

func TestSessionStoreNoExpiry(t *testing.T) {
	t.Parallel()

	session := newMockSession()
	clk := &fakeClock{now: time.Now()}
	s := cache.NewSessionStoreWithClock(session, "", clk)
	ctx := context.Background()

	_ = s.Put(ctx, "key", "val", 0)
	clk.Advance(365 * 24 * time.Hour)

	_, err := s.Get(ctx, "key")

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestSessionStoreAdd(t *testing.T) {
	t.Parallel()

	session := newMockSession()
	s := cache.NewSessionStore(session, "")
	ctx := context.Background()

	ok, _ := s.Add(ctx, "k", "v1", time.Minute)

	if !ok {
		t.Fatal("expected Add to succeed")
	}

	ok, _ = s.Add(ctx, "k", "v2", time.Minute)

	if ok {
		t.Fatal("expected Add to fail for existing key")
	}
}

func TestSessionStoreForever(t *testing.T) {
	t.Parallel()

	session := newMockSession()
	s := cache.NewSessionStore(session, "")
	ctx := context.Background()

	_ = s.Forever(ctx, "k", "eternal")

	v, err := s.Get(ctx, "k")

	if err != nil {
		t.Fatal(err)
	}

	if v != "eternal" {
		t.Fatalf("expected 'eternal', got %v", v)
	}
}

func TestSessionStoreIncrement(t *testing.T) {
	t.Parallel()

	session := newMockSession()
	s := cache.NewSessionStore(session, "")
	ctx := context.Background()

	v, _ := s.Increment(ctx, "counter", 3)

	if v != 3 {
		t.Fatalf("expected 3, got %d", v)
	}

	v, _ = s.Increment(ctx, "counter", 5)

	if v != 8 {
		t.Fatalf("expected 8, got %d", v)
	}
}

func TestSessionStoreDecrement(t *testing.T) {
	t.Parallel()

	session := newMockSession()
	s := cache.NewSessionStore(session, "")
	ctx := context.Background()

	_, _ = s.Increment(ctx, "k", 10)
	v, _ := s.Decrement(ctx, "k", 3)

	if v != 7 {
		t.Fatalf("expected 7, got %d", v)
	}
}

func TestSessionStoreForget(t *testing.T) {
	t.Parallel()

	session := newMockSession()
	s := cache.NewSessionStore(session, "")
	ctx := context.Background()

	_ = s.Put(ctx, "k", "v", 0)
	_ = s.Forget(ctx, "k")

	_, err := s.Get(ctx, "k")

	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatal("expected ErrNotFound after Forget")
	}
}

func TestSessionStoreFlush(t *testing.T) {
	t.Parallel()

	session := newMockSession()
	s := cache.NewSessionStore(session, "")
	ctx := context.Background()

	_ = s.Put(ctx, "a", 1, 0)
	_ = s.Put(ctx, "b", 2, 0)
	_ = s.Flush(ctx)

	_, err := s.Get(ctx, "a")

	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatal("expected ErrNotFound after Flush")
	}
}

func TestSessionStoreTouch(t *testing.T) {
	t.Parallel()

	session := newMockSession()
	clk := &fakeClock{now: time.Now()}
	s := cache.NewSessionStoreWithClock(session, "", clk)
	ctx := context.Background()

	_ = s.Put(ctx, "k", "v", 10*time.Second)

	ok, _ := s.Touch(ctx, "k", 30*time.Second)

	if !ok {
		t.Fatal("expected Touch to succeed")
	}

	clk.Advance(15 * time.Second)

	_, err := s.Get(ctx, "k")

	if err != nil {
		t.Fatal("expected key to survive after Touch")
	}
}

func TestSessionStorePrefix(t *testing.T) {
	t.Parallel()

	session := newMockSession()
	s := cache.NewSessionStore(session, "app")

	if s.GetPrefix() != "app" {
		t.Fatalf("expected 'app', got %q", s.GetPrefix())
	}
}

func TestSessionStoreGetMany(t *testing.T) {
	t.Parallel()

	session := newMockSession()
	s := cache.NewSessionStore(session, "")
	ctx := context.Background()

	_ = s.Put(ctx, "a", "x", 0)
	_ = s.Put(ctx, "b", "y", 0)

	got, _ := s.GetMany(ctx, []string{"a", "b", "c"})

	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestSessionStorePutMany(t *testing.T) {
	t.Parallel()

	session := newMockSession()
	s := cache.NewSessionStore(session, "")
	ctx := context.Background()

	_ = s.PutMany(ctx, map[string]any{"x": 10, "y": 20}, time.Minute)

	v, _ := s.Get(ctx, "x")

	if v != 10 {
		t.Fatalf("expected 10, got %v", v)
	}
}

func TestSessionStoreIncrementNonNumeric(t *testing.T) {
	t.Parallel()

	session := newMockSession()
	s := cache.NewSessionStore(session, "")
	ctx := context.Background()

	_ = s.Put(ctx, "k", "hello", time.Minute)

	_, err := s.Increment(ctx, "k", 1)

	if !errors.Is(err, cache.ErrInvalidValue) {
		t.Fatalf("expected ErrInvalidValue, got %v", err)
	}
}
