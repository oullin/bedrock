package cache_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/cache"
)

func TestCacheLockAcquireAndRelease(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	l := cache.NewCacheLock(store, "res", "owner-1", time.Minute)
	ctx := context.Background()

	ok, err := l.Acquire(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("expected acquire to succeed")
	}

	ok, err = l.Release(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("expected release to succeed")
	}
}

func TestCacheLockOwnershipVerification(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	l1 := cache.NewCacheLock(store, "res", "owner-1", time.Minute)
	l2 := cache.NewCacheLock(store, "res", "owner-2", time.Minute)
	ctx := context.Background()

	l1.Acquire(ctx) //nolint:errcheck

	// Different owner cannot acquire.
	ok, _ := l2.Acquire(ctx)

	if ok {
		t.Fatal("expected second owner to fail acquiring")
	}

	// Different owner cannot release.
	ok, _ = l2.Release(ctx)

	if ok {
		t.Fatal("expected second owner to fail releasing")
	}
}

func TestCacheLockForceRelease(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	l1 := cache.NewCacheLock(store, "res", "owner-1", time.Minute)
	l2 := cache.NewCacheLock(store, "res", "owner-2", time.Minute)
	ctx := context.Background()

	l1.Acquire(ctx) //nolint:errcheck

	_ = l2.ForceRelease(ctx)

	ok, _ := l2.Acquire(ctx)

	if !ok {
		t.Fatal("expected acquire after force release")
	}
}

func TestCacheLockGet(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	l := cache.NewCacheLock(store, "res", "owner", time.Minute)
	ctx := context.Background()

	executed := false

	err := l.Get(ctx, func() error {
		executed = true

		return nil
	})

	if err != nil {
		t.Fatal(err)
	}

	if !executed {
		t.Fatal("expected callback to execute")
	}
}

func TestCacheLockGetFailsWhenHeld(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	l1 := cache.NewCacheLock(store, "res", "owner-1", time.Minute)
	l2 := cache.NewCacheLock(store, "res", "owner-2", time.Minute)
	ctx := context.Background()

	l1.Acquire(ctx) //nolint:errcheck

	err := l2.Get(ctx, func() error { return nil })

	if !errors.Is(err, cache.ErrLockTimeout) {
		t.Fatalf("expected ErrLockTimeout, got %v", err)
	}
}

func TestCacheLockBlocked(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	l1 := cache.NewCacheLock(store, "res", "owner-1", time.Minute)
	l2 := cache.NewCacheLock(store, "res", "owner-2", time.Minute)
	ctx := context.Background()

	l1.Acquire(ctx) //nolint:errcheck

	blocked, _ := l2.Blocked(ctx)

	if !blocked {
		t.Fatal("expected blocked")
	}
}

func TestCacheLockBlockTimeout(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	l1 := cache.NewCacheLock(store, "res", "owner-1", time.Minute)
	l2 := cache.NewCacheLock(store, "res", "owner-2", time.Minute)
	ctx := context.Background()

	l1.Acquire(ctx) //nolint:errcheck

	err := l2.Block(ctx, 100*time.Millisecond)

	if !errors.Is(err, cache.ErrLockTimeout) {
		t.Fatalf("expected ErrLockTimeout, got %v", err)
	}
}

func TestCacheLockOwner(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	l := cache.NewCacheLock(store, "res", "owner-42", time.Minute)

	if l.Owner() != "owner-42" {
		t.Fatalf("expected 'owner-42', got %q", l.Owner())
	}
}

func TestCacheLockIsOwnedByCurrentProcess(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	l := cache.NewCacheLock(store, "res", "owner-1", time.Minute)
	ctx := context.Background()

	owned, _ := l.IsOwnedByCurrentProcess(ctx)

	if owned {
		t.Fatal("expected false before acquire")
	}

	l.Acquire(ctx) //nolint:errcheck

	owned, _ = l.IsOwnedByCurrentProcess(ctx)

	if !owned {
		t.Fatal("expected true after acquire")
	}
}

func TestCacheLockIsOwnedBy(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	l := cache.NewCacheLock(store, "res", "owner-1", time.Minute)
	ctx := context.Background()

	l.Acquire(ctx) //nolint:errcheck

	owned, _ := l.IsOwnedBy(ctx, "owner-1")

	if !owned {
		t.Fatal("expected owned by owner-1")
	}

	owned, _ = l.IsOwnedBy(ctx, "owner-2")

	if owned {
		t.Fatal("expected not owned by owner-2")
	}
}

func TestCacheLockBetweenBlockedAttemptsSleepFor(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	l := cache.NewCacheLock(store, "res", "owner-1", time.Minute)

	result := l.BetweenBlockedAttemptsSleepFor(100)

	if result != l {
		t.Fatal("expected same lock returned for chaining")
	}
}

func TestCacheLockRestoreLockRelease(t *testing.T) {
	t.Parallel()

	store := cache.NewArrayStore()
	ctx := context.Background()

	lock := store.Lock("res", "owner-1", time.Minute)
	ok, _ := lock.Acquire(ctx)

	if !ok {
		t.Fatal("expected acquire")
	}

	restored := store.RestoreLock("res", "owner-1")
	released, _ := restored.Release(ctx)

	if !released {
		t.Fatal("expected restore+release to succeed")
	}
}
