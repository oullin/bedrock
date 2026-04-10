package cache_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/cache"
)

func TestRedisStoreBasicGetPut(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	s := cache.NewRedisStore(client, "")
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

func TestRedisStoreMissingKey(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	s := cache.NewRedisStore(client, "")

	_, err := s.Get(context.Background(), "missing")

	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestRedisStorePrefixed(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	s := cache.NewRedisStore(client, "app")
	ctx := context.Background()

	_ = s.Put(ctx, "key", "val", time.Minute)

	// The key should be stored with prefix.
	if _, ok := client.data["app:key"]; !ok {
		t.Fatal("expected key to be stored with prefix 'app:key'")
	}
}

func TestRedisStoreGetPrefix(t *testing.T) {
	t.Parallel()

	s := cache.NewRedisStore(newMockRedisClient(), "myprefix")

	if s.GetPrefix() != "myprefix" {
		t.Fatalf("expected 'myprefix', got %q", s.GetPrefix())
	}
}

func TestRedisStoreAdd(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	s := cache.NewRedisStore(client, "")
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

func TestRedisStoreForget(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	s := cache.NewRedisStore(client, "")
	ctx := context.Background()

	_ = s.Put(ctx, "k", "v", time.Minute)
	_ = s.Forget(ctx, "k")

	_, err := s.Get(ctx, "k")

	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatal("expected ErrNotFound after Forget")
	}
}

func TestRedisStoreFlush(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	s := cache.NewRedisStore(client, "")
	ctx := context.Background()

	_ = s.Put(ctx, "a", 1, 0)
	_ = s.Put(ctx, "b", 2, 0)
	_ = s.Flush(ctx)

	if len(client.data) != 0 {
		t.Fatalf("expected empty store after flush, got %d entries", len(client.data))
	}
}

func TestRedisStoreForever(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	s := cache.NewRedisStore(client, "")
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

func TestRedisStoreGetMany(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	s := cache.NewRedisStore(client, "")
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

func TestRedisStorePutMany(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	s := cache.NewRedisStore(client, "")
	ctx := context.Background()

	_ = s.PutMany(ctx, map[string]any{"a": 1, "b": 2}, time.Minute)

	v, _ := s.Get(ctx, "a")

	if v != float64(1) {
		t.Fatalf("expected 1 (as float64 from JSON), got %v (%T)", v, v)
	}
}

func TestRedisStoreTags(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	s := cache.NewRedisStore(client, "")
	ctx := context.Background()

	tagged := s.Tags("users")
	_ = tagged.Put(ctx, "name", "Bob", time.Minute)

	v, err := tagged.Get(ctx, "name")

	if err != nil {
		t.Fatal(err)
	}

	if v != "Bob" {
		t.Fatalf("expected 'Bob', got %v", v)
	}
}

func TestRedisStoreLock(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	s := cache.NewRedisStore(client, "")
	ctx := context.Background()

	l := s.Lock("res", "owner", time.Minute)

	ok, err := l.Acquire(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("expected acquire to succeed")
	}
}

func TestRedisStoreJSONSerialization(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	s := cache.NewRedisStore(client, "")
	ctx := context.Background()

	_ = s.Put(ctx, "num", 42, time.Minute)

	v, _ := s.Get(ctx, "num")

	// JSON unmarshals numbers as float64.
	if v != float64(42) {
		t.Fatalf("expected float64(42), got %v (%T)", v, v)
	}
}
