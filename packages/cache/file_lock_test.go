package cache_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/cache"
)

func TestFileLockAcquireAndRelease(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	l := cache.NewFileLock(dir, "res", "owner-1", time.Minute, nil)
	ctx := context.Background()

	ok, err := l.Acquire(ctx)

	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("expected acquire to succeed")
	}

	ok, _ = l.Release(ctx)

	if !ok {
		t.Fatal("expected release to succeed")
	}
}

func TestFileLockOwnershipVerification(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	l1 := cache.NewFileLock(dir, "res", "owner-1", time.Minute, nil)
	l2 := cache.NewFileLock(dir, "res", "owner-2", time.Minute, nil)
	ctx := context.Background()

	l1.Acquire(ctx) //nolint:errcheck

	ok, _ := l2.Acquire(ctx)

	if ok {
		t.Fatal("expected second owner to fail")
	}

	// Cannot release with wrong owner.
	ok, _ = l2.Release(ctx)

	if ok {
		t.Fatal("expected release to fail for wrong owner")
	}
}

func TestFileLockForceRelease(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	l1 := cache.NewFileLock(dir, "res", "owner-1", time.Minute, nil)
	l2 := cache.NewFileLock(dir, "res", "owner-2", time.Minute, nil)
	ctx := context.Background()

	l1.Acquire(ctx) //nolint:errcheck

	_ = l2.ForceRelease(ctx)

	ok, _ := l2.Acquire(ctx)

	if !ok {
		t.Fatal("expected acquire after force release")
	}
}

func TestFileLockGet(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	l := cache.NewFileLock(dir, "res", "owner", time.Minute, nil)
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

func TestFileLockBlocked(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	l1 := cache.NewFileLock(dir, "res", "owner-1", time.Minute, nil)
	l2 := cache.NewFileLock(dir, "res", "owner-2", time.Minute, nil)
	ctx := context.Background()

	l1.Acquire(ctx) //nolint:errcheck

	blocked, _ := l2.Blocked(ctx)

	if !blocked {
		t.Fatal("expected blocked")
	}
}

func TestFileLockBlockTimeout(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	l1 := cache.NewFileLock(dir, "res", "owner-1", time.Minute, nil)
	l2 := cache.NewFileLock(dir, "res", "owner-2", time.Minute, nil)
	ctx := context.Background()

	l1.Acquire(ctx) //nolint:errcheck

	err := l2.Block(ctx, 100*time.Millisecond)

	if !errors.Is(err, cache.ErrLockTimeout) {
		t.Fatalf("expected ErrLockTimeout, got %v", err)
	}
}

func TestFileLockExpiry(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	clk := &fakeClock{now: time.Now()}
	l1 := cache.NewFileLock(dir, "res", "owner-1", 5*time.Second, clk)
	ctx := context.Background()

	l1.Acquire(ctx) //nolint:errcheck

	clk.Advance(10 * time.Second)

	l2 := cache.NewFileLock(dir, "res", "owner-2", 5*time.Second, clk)
	ok, _ := l2.Acquire(ctx)

	if !ok {
		t.Fatal("expected acquire after lock expired")
	}
}

func TestFileLockReacquireBySameOwner(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	l := cache.NewFileLock(dir, "res", "owner", time.Minute, nil)
	ctx := context.Background()

	ok, _ := l.Acquire(ctx)

	if !ok {
		t.Fatal("expected first acquire to succeed")
	}

	ok, _ = l.Acquire(ctx)

	if !ok {
		t.Fatal("expected reacquire by same owner to succeed")
	}
}
