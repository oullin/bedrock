package cache_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/cache"
)

func newRepo() *cache.Repository {
	return cache.NewRepository(cache.NewArrayStore())
}

func TestRepositoryHas(t *testing.T) {
	t.Parallel()

	r := newRepo()
	ctx := context.Background()

	if r.Has(ctx, "k") {
		t.Fatal("expected Has to be false for missing key")
	}

	_ = r.Put(ctx, "k", "v", time.Minute)

	if !r.Has(ctx, "k") {
		t.Fatal("expected Has to be true after Put")
	}
}

func TestRepositoryMissing(t *testing.T) {
	t.Parallel()

	r := newRepo()
	ctx := context.Background()

	if !r.Missing(ctx, "k") {
		t.Fatal("expected Missing for absent key")
	}
}

func TestRepositoryGet(t *testing.T) {
	t.Parallel()

	r := newRepo()
	ctx := context.Background()

	got := r.Get(ctx, "missing", "default")

	if got != "default" {
		t.Fatalf("expected default, got %v", got)
	}

	_ = r.Put(ctx, "k", "v", time.Minute)
	got = r.Get(ctx, "k", "default")

	if got != "v" {
		t.Fatalf("expected 'v', got %v", got)
	}
}

func TestRepositoryPull(t *testing.T) {
	t.Parallel()

	r := newRepo()
	ctx := context.Background()

	_ = r.Put(ctx, "k", "v", time.Minute)

	got := r.Pull(ctx, "k", nil)

	if got != "v" {
		t.Fatalf("expected 'v', got %v", got)
	}

	if r.Has(ctx, "k") {
		t.Fatal("expected key to be deleted after Pull")
	}
}

func TestRepositoryRemember(t *testing.T) {
	t.Parallel()

	r := newRepo()
	ctx := context.Background()

	calls := 0
	fn := func() (any, error) {
		calls++

		return "computed", nil
	}

	v1, _ := r.Remember(ctx, "k", time.Minute, fn)
	v2, _ := r.Remember(ctx, "k", time.Minute, fn)

	if v1 != "computed" || v2 != "computed" {
		t.Fatalf("unexpected values: %v %v", v1, v2)
	}

	if calls != 1 {
		t.Fatalf("expected fn to be called once, got %d", calls)
	}
}

func TestRepositoryRememberForever(t *testing.T) {
	t.Parallel()

	r := newRepo()
	ctx := context.Background()

	calls := 0
	fn := func() (any, error) { calls++; return "forever", nil }

	_, _ = r.RememberForever(ctx, "k", fn)
	_, _ = r.RememberForever(ctx, "k", fn)

	if calls != 1 {
		t.Fatalf("expected fn called once, got %d", calls)
	}
}

func TestRepositoryRememberPropagatesError(t *testing.T) {
	t.Parallel()

	r := newRepo()
	ctx := context.Background()

	sentinel := errors.New("oops")

	_, err := r.Remember(ctx, "k", time.Minute, func() (any, error) {
		return nil, sentinel
	})

	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestRepositoryLock(t *testing.T) {
	t.Parallel()

	r := cache.NewRepository(cache.NewArrayStore())
	l := r.Lock("res", "o1", time.Minute)

	if l == nil {
		t.Fatal("expected non-nil lock from ArrayStore")
	}

	ctx := context.Background()

	ok, _ := l.Acquire(ctx)

	if !ok {
		t.Fatal("expected acquire to succeed")
	}
}

func TestRepositoryLockNilForNullStore(t *testing.T) {
	t.Parallel()

	// NullStore does implement Locker, so lock is non-nil.
	r := cache.NewRepository(cache.NewNullStore())
	l := r.Lock("res", "o", time.Minute)

	if l == nil {
		t.Fatal("expected non-nil lock from NullStore")
	}
}
