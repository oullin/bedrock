package drivers_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
	"github.com/bedrock/packages/queue/drivers"
)

func TestBackgroundDriverPushDelegatesToInner(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	inner := drivers.NewRedisDriver(client, "redis")

	// Use a non-existent command so spawn is a harmless no-op.
	drv := drivers.NewBackgroundDriver("true", nil, inner, "bg")

	_, err := drv.Push(context.Background(), "default", []byte("payload"))

	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Millisecond) // Let goroutine start.

	n, _ := inner.Size(context.Background(), "default")

	if n != 1 {
		t.Errorf("expected 1 in inner queue, got %d", n)
	}
}

func TestBackgroundDriverPushDelayedDelegates(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	inner := drivers.NewRedisDriver(client, "redis")
	drv := drivers.NewBackgroundDriver("true", nil, inner, "bg")

	_, err := drv.PushDelayed(context.Background(), "default", []byte("payload"), time.Millisecond)

	if err != nil {
		t.Fatal(err)
	}

	n, _ := inner.DelayedSize(context.Background(), "default")

	if n != 1 {
		t.Errorf("expected 1 delayed in inner, got %d", n)
	}
}

func TestBackgroundDriverPopDelegates(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	inner := drivers.NewRedisDriver(client, "redis")
	drv := drivers.NewBackgroundDriver("true", nil, inner, "bg")

	_, _ = inner.Push(context.Background(), "default", []byte("payload"))

	job, err := drv.Pop(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if string(job.Payload()) != "payload" {
		t.Errorf("expected 'payload', got %q", job.Payload())
	}
}

func TestBackgroundDriverPopEmptyDelegates(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	inner := drivers.NewRedisDriver(client, "redis")
	drv := drivers.NewBackgroundDriver("true", nil, inner, "bg")

	_, err := drv.Pop(context.Background(), "default")

	if !errors.Is(err, queue.ErrNoJob) {
		t.Fatalf("expected ErrNoJob, got %v", err)
	}
}

func TestBackgroundDriverSizeDelegates(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	inner := drivers.NewRedisDriver(client, "redis")
	drv := drivers.NewBackgroundDriver("true", nil, inner, "bg")

	_, _ = inner.Push(context.Background(), "default", []byte("a"))
	_, _ = inner.Push(context.Background(), "default", []byte("b"))

	n, err := drv.Size(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if n != 2 {
		t.Errorf("expected 2, got %d", n)
	}
}

func TestBackgroundDriverConnectionName(t *testing.T) {
	t.Parallel()

	inner := drivers.NewNullDriver("null")
	drv := drivers.NewBackgroundDriver("true", nil, inner, "my-bg")

	if drv.ConnectionName() != "my-bg" {
		t.Errorf("expected 'my-bg', got %q", drv.ConnectionName())
	}
}

func TestBackgroundDriverPendingSizeDelegates(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	inner := drivers.NewRedisDriver(client, "redis")
	drv := drivers.NewBackgroundDriver("true", nil, inner, "bg")

	n, err := drv.PendingSize(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
}

func TestBackgroundDriverDelayedSizeDelegates(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	inner := drivers.NewRedisDriver(client, "redis")
	drv := drivers.NewBackgroundDriver("true", nil, inner, "bg")

	n, err := drv.DelayedSize(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
}

func TestBackgroundDriverReservedSizeDelegates(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	inner := drivers.NewRedisDriver(client, "redis")
	drv := drivers.NewBackgroundDriver("true", nil, inner, "bg")

	n, err := drv.ReservedSize(context.Background(), "default")

	if err != nil {
		t.Fatal(err)
	}

	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
}
