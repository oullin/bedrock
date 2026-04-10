package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/bedrock/packages/cache"
)

func TestNoLockAcquire(t *testing.T) {
	t.Parallel()

	l := &cache.NoLock{}
	ok, err := l.Acquire(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("expected always true")
	}
}

func TestNoLockRelease(t *testing.T) {
	t.Parallel()

	l := &cache.NoLock{}
	ok, err := l.Release(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("expected always true")
	}
}

func TestNoLockGet(t *testing.T) {
	t.Parallel()

	l := &cache.NoLock{}
	executed := false

	err := l.Get(context.Background(), func() error {
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

func TestNoLockBlock(t *testing.T) {
	t.Parallel()

	l := &cache.NoLock{}

	err := l.Block(context.Background(), time.Second)

	if err != nil {
		t.Fatal(err)
	}
}

func TestNoLockBlocked(t *testing.T) {
	t.Parallel()

	l := &cache.NoLock{}
	blocked, err := l.Blocked(context.Background())

	if err != nil {
		t.Fatal(err)
	}

	if blocked {
		t.Fatal("expected never blocked")
	}
}

func TestNoLockForceRelease(t *testing.T) {
	t.Parallel()

	l := &cache.NoLock{}

	err := l.ForceRelease(context.Background())

	if err != nil {
		t.Fatal(err)
	}
}
