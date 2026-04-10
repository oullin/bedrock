package cache_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/cache"
)

func TestNullStoreGet(t *testing.T) {
	t.Parallel()

	s := cache.NewNullStore()
	_ = s.Put(context.Background(), "k", "v", time.Minute)

	_, err := s.Get(context.Background(), "k")

	if !errors.Is(err, cache.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestNullStoreGetMany(t *testing.T) {
	t.Parallel()

	s := cache.NewNullStore()
	got, err := s.GetMany(context.Background(), []string{"a", "b"})

	if err != nil || len(got) != 0 {
		t.Fatalf("expected empty map, got %v %v", got, err)
	}
}

func TestNullStoreAdd(t *testing.T) {
	t.Parallel()

	s := cache.NewNullStore()
	ok, err := s.Add(context.Background(), "k", "v", time.Minute)

	if err != nil || !ok {
		t.Fatalf("expected Add to succeed, got ok=%v err=%v", ok, err)
	}
}

func TestNullStoreLock(t *testing.T) {
	t.Parallel()

	s := cache.NewNullStore()
	l := s.Lock("r", "owner", time.Minute)

	ok, err := l.Acquire(context.Background())

	if err != nil || !ok {
		t.Fatalf("expected no-op lock acquire to succeed")
	}

	blocked, err := l.Blocked(context.Background())

	if err != nil || blocked {
		t.Fatalf("expected no-op lock to report not blocked")
	}
}
