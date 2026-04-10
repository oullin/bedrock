package cache_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/cache"
)

func TestDatabaseStoreBasicGetPut(t *testing.T) {
	t.Parallel()

	conn := newMockDBConnection()
	s := cache.NewDatabaseStore(conn, "cache", "")
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

func TestDatabaseStoreMissingKey(t *testing.T) {
	t.Parallel()

	conn := newMockDBConnection()
	s := cache.NewDatabaseStore(conn, "cache", "")

	_, err := s.Get(context.Background(), "missing")

	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDatabaseStorePrefix(t *testing.T) {
	t.Parallel()

	conn := newMockDBConnection()
	s := cache.NewDatabaseStore(conn, "cache", "app_")
	ctx := context.Background()

	_ = s.Put(ctx, "key", "val", time.Minute)

	// Key should be stored with prefix.
	if _, ok := conn.data["app_key"]; !ok {
		t.Fatal("expected key with prefix 'app_key'")
	}

	if s.GetPrefix() != "app_" {
		t.Fatalf("expected 'app_', got %q", s.GetPrefix())
	}
}

func TestDatabaseStoreAdd(t *testing.T) {
	t.Parallel()

	conn := newMockDBConnection()
	s := cache.NewDatabaseStore(conn, "cache", "")
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

func TestDatabaseStoreForever(t *testing.T) {
	t.Parallel()

	conn := newMockDBConnection()
	s := cache.NewDatabaseStore(conn, "cache", "")
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

func TestDatabaseStoreIncrement(t *testing.T) {
	t.Parallel()

	conn := newMockDBConnection()
	s := cache.NewDatabaseStore(conn, "cache", "")
	ctx := context.Background()

	v, _ := s.Increment(ctx, "counter", 5)

	if v != 5 {
		t.Fatalf("expected 5, got %d", v)
	}

	v, _ = s.Increment(ctx, "counter", 3)

	if v != 8 {
		t.Fatalf("expected 8, got %d", v)
	}
}

func TestDatabaseStoreDecrement(t *testing.T) {
	t.Parallel()

	conn := newMockDBConnection()
	s := cache.NewDatabaseStore(conn, "cache", "")
	ctx := context.Background()

	_, _ = s.Increment(ctx, "k", 10)
	v, _ := s.Decrement(ctx, "k", 3)

	if v != 7 {
		t.Fatalf("expected 7, got %d", v)
	}
}

func TestDatabaseStoreForget(t *testing.T) {
	t.Parallel()

	conn := newMockDBConnection()
	s := cache.NewDatabaseStore(conn, "cache", "")
	ctx := context.Background()

	_ = s.Put(ctx, "k", "v", 0)
	_ = s.Forget(ctx, "k")

	_, err := s.Get(ctx, "k")

	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatal("expected ErrNotFound after Forget")
	}
}

func TestDatabaseStoreFlush(t *testing.T) {
	t.Parallel()

	conn := newMockDBConnection()
	s := cache.NewDatabaseStore(conn, "cache", "")
	ctx := context.Background()

	_ = s.Put(ctx, "a", 1, 0)
	_ = s.Put(ctx, "b", 2, 0)
	_ = s.Flush(ctx)

	if len(conn.data) != 0 {
		t.Fatalf("expected empty after flush, got %d", len(conn.data))
	}
}

func TestDatabaseStoreGetMany(t *testing.T) {
	t.Parallel()

	conn := newMockDBConnection()
	s := cache.NewDatabaseStore(conn, "cache", "")
	ctx := context.Background()

	_ = s.Put(ctx, "a", "x", 0)
	_ = s.Put(ctx, "b", "y", 0)

	got, err := s.GetMany(ctx, []string{"a", "b", "c"})

	if err != nil {
		t.Fatal(err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestDatabaseStoreTouch(t *testing.T) {
	t.Parallel()

	conn := newMockDBConnection()
	s := cache.NewDatabaseStore(conn, "cache", "")
	ctx := context.Background()

	_ = s.Put(ctx, "k", "v", time.Minute)

	ok, _ := s.Touch(ctx, "k", time.Hour)

	if !ok {
		t.Fatal("expected Touch to succeed")
	}
}

func TestDatabaseStoreTouchMissing(t *testing.T) {
	t.Parallel()

	conn := newMockDBConnection()
	s := cache.NewDatabaseStore(conn, "cache", "")

	ok, _ := s.Touch(context.Background(), "missing", time.Minute)

	if ok {
		t.Fatal("expected Touch to return false for missing key")
	}
}

func TestDatabaseStoreTags(t *testing.T) {
	t.Parallel()

	conn := newMockDBConnection()
	s := cache.NewDatabaseStore(conn, "cache", "")
	ctx := context.Background()

	tagged := s.Tags("orders")
	_ = tagged.Put(ctx, "item", "widget", time.Minute)

	v, err := tagged.Get(ctx, "item")

	if err != nil {
		t.Fatal(err)
	}

	if v != "widget" {
		t.Fatalf("expected 'widget', got %v", v)
	}
}

func TestDatabaseStoreDefaultTable(t *testing.T) {
	t.Parallel()

	conn := newMockDBConnection()
	s := cache.NewDatabaseStore(conn, "", "")

	// Default table should be "cache", no panic.
	_ = s.Put(context.Background(), "k", "v", time.Minute)
}

func TestDatabaseStoreLock(t *testing.T) {
	t.Parallel()

	conn := newMockDBConnection()
	s := cache.NewDatabaseStore(conn, "cache", "")

	l := s.Lock("res", "owner", time.Minute)

	if l == nil {
		t.Fatal("expected non-nil lock")
	}
}
