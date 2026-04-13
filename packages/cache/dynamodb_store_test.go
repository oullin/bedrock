package cache_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/cache"
)

func TestDynamoDbStoreBasicGetPut(t *testing.T) {
	t.Parallel()

	client := newMockDynamoClient()
	s := cache.NewDynamoDbStore(client, "cache", "")
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

func TestDynamoDbStoreMissingKey(t *testing.T) {
	t.Parallel()

	client := newMockDynamoClient()
	s := cache.NewDynamoDbStore(client, "cache", "")

	_, err := s.Get(context.Background(), "missing")

	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDynamoDbStorePrefix(t *testing.T) {
	t.Parallel()

	client := newMockDynamoClient()
	s := cache.NewDynamoDbStore(client, "cache", "app")
	ctx := context.Background()

	_ = s.Put(ctx, "key", "val", time.Minute)

	if _, ok := client.data["app:key"]; !ok {
		t.Fatal("expected prefixed key 'app:key'")
	}

	if s.GetPrefix() != "app" {
		t.Fatalf("expected 'app', got %q", s.GetPrefix())
	}
}

func TestDynamoDbStoreAdd(t *testing.T) {
	t.Parallel()

	client := newMockDynamoClient()
	s := cache.NewDynamoDbStore(client, "cache", "")
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

func TestDynamoDbStoreForever(t *testing.T) {
	t.Parallel()

	client := newMockDynamoClient()
	s := cache.NewDynamoDbStore(client, "cache", "")
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

func TestDynamoDbStoreForget(t *testing.T) {
	t.Parallel()

	client := newMockDynamoClient()
	s := cache.NewDynamoDbStore(client, "cache", "")
	ctx := context.Background()

	_ = s.Put(ctx, "k", "v", 0)
	_ = s.Forget(ctx, "k")

	_, err := s.Get(ctx, "k")

	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatal("expected ErrNotFound after Forget")
	}
}

func TestDynamoDbStoreFlush(t *testing.T) {
	t.Parallel()

	client := newMockDynamoClient()
	s := cache.NewDynamoDbStore(client, "cache", "")
	ctx := context.Background()

	_ = s.Put(ctx, "a", 1, 0)
	_ = s.Put(ctx, "b", 2, 0)
	_ = s.Flush(ctx)

	if len(client.data) != 0 {
		t.Fatalf("expected empty after flush, got %d", len(client.data))
	}
}

func TestDynamoDbStoreGetMany(t *testing.T) {
	t.Parallel()

	client := newMockDynamoClient()
	s := cache.NewDynamoDbStore(client, "cache", "")
	ctx := context.Background()

	_ = s.Put(ctx, "a", "x", 0)
	_ = s.Put(ctx, "b", "y", 0)

	got, _ := s.GetMany(ctx, []string{"a", "b", "c"})

	if len(got) != 2 {
		t.Fatalf("expected 2, got %d", len(got))
	}
}

func TestDynamoDbStoreLock(t *testing.T) {
	t.Parallel()

	client := newMockDynamoClient()
	s := cache.NewDynamoDbStore(client, "cache", "")

	l := s.Lock("res", "owner", time.Minute)

	if l == nil {
		t.Fatal("expected non-nil lock")
	}
}

func TestDynamoDbStoreSetPrefix(t *testing.T) {
	t.Parallel()

	client := newMockDynamoClient()
	s := cache.NewDynamoDbStore(client, "cache", "old")

	if s.GetPrefix() != "old" {
		t.Fatalf("expected 'old', got %q", s.GetPrefix())
	}

	s.SetPrefix("new")

	if s.GetPrefix() != "new" {
		t.Fatalf("expected 'new', got %q", s.GetPrefix())
	}
}

func TestDynamoDbStoreGetClient(t *testing.T) {
	t.Parallel()

	client := newMockDynamoClient()
	s := cache.NewDynamoDbStore(client, "cache", "")

	if s.GetClient() != client {
		t.Fatal("expected same client")
	}
}
