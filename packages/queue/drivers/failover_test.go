package drivers_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
	"github.com/bedrock/packages/queue/drivers"
)

func TestFailoverDriverPushUsesFirstSuccessful(t *testing.T) {
	t.Parallel()

	d1 := drivers.NewNullDriver("d1")
	d2 := drivers.NewNullDriver("d2")
	drv := drivers.NewFailoverDriver("failover", d1, d2)

	_, err := drv.Push(context.Background(), "q", []byte("payload"))

	if err != nil {
		t.Fatal(err)
	}
}

func TestFailoverDriverPushFallsBackOnError(t *testing.T) {
	t.Parallel()

	handler := queue.HandlerFunc(func(_ context.Context, _ queue.Job) error {
		return errors.New("fail")
	})

	failing := drivers.NewSyncDriver("d1", handler)

	ok := drivers.NewNullDriver("d2")

	drv := drivers.NewFailoverDriver("failover", failing, ok)

	_, err := drv.Push(context.Background(), "q", []byte("payload"))

	if err != nil {
		t.Fatal(err)
	}
}

func TestFailoverDriverPushAllFail(t *testing.T) {
	t.Parallel()

	handler := queue.HandlerFunc(func(_ context.Context, _ queue.Job) error {
		return errors.New("fail")
	})

	d1 := drivers.NewSyncDriver("d1", handler)
	d2 := drivers.NewSyncDriver("d2", handler)
	drv := drivers.NewFailoverDriver("failover", d1, d2)

	_, err := drv.Push(context.Background(), "q", []byte("payload"))

	if err == nil {
		t.Fatal("expected error when all drivers fail")
	}
}

func TestFailoverDriverPushDelayedFallback(t *testing.T) {
	t.Parallel()

	d1 := drivers.NewNullDriver("d1")
	d2 := drivers.NewNullDriver("d2")
	drv := drivers.NewFailoverDriver("failover", d1, d2)

	_, err := drv.PushDelayed(context.Background(), "q", []byte("payload"), 5*time.Second)

	if err != nil {
		t.Fatal(err)
	}
}

func TestFailoverDriverPushMultipleFallback(t *testing.T) {
	t.Parallel()

	d1 := drivers.NewNullDriver("d1")
	d2 := drivers.NewNullDriver("d2")
	drv := drivers.NewFailoverDriver("failover", d1, d2)

	ids, err := drv.PushMultiple(context.Background(), "q", [][]byte{[]byte("a"), []byte("b")})

	if err != nil {
		t.Fatal(err)
	}

	if len(ids) != 2 {
		t.Errorf("expected 2 ids, got %d", len(ids))
	}
}

func TestFailoverDriverPopTriesEach(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	redis := drivers.NewRedisDriver(client, "redis")
	null := drivers.NewNullDriver("null")

	_, _ = redis.Push(context.Background(), "q", []byte("from-redis"))

	drv := drivers.NewFailoverDriver("failover", null, redis)

	job, err := drv.Pop(context.Background(), "q")

	if err != nil {
		t.Fatal(err)
	}

	if string(job.Payload()) != "from-redis" {
		t.Errorf("expected 'from-redis', got %q", job.Payload())
	}
}

func TestFailoverDriverPopAllEmpty(t *testing.T) {
	t.Parallel()

	d1 := drivers.NewNullDriver("d1")
	d2 := drivers.NewNullDriver("d2")
	drv := drivers.NewFailoverDriver("failover", d1, d2)

	_, err := drv.Pop(context.Background(), "q")

	if !errors.Is(err, queue.ErrNoJob) {
		t.Fatalf("expected ErrNoJob, got %v", err)
	}
}

func TestFailoverDriverSizeFallback(t *testing.T) {
	t.Parallel()

	d1 := drivers.NewNullDriver("d1")
	d2 := drivers.NewNullDriver("d2")
	drv := drivers.NewFailoverDriver("failover", d1, d2)

	n, err := drv.Size(context.Background(), "q")

	if err != nil {
		t.Fatal(err)
	}

	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
}

func TestFailoverDriverConnectionName(t *testing.T) {
	t.Parallel()

	drv := drivers.NewFailoverDriver("my-failover")

	if drv.ConnectionName() != "my-failover" {
		t.Errorf("expected 'my-failover', got %q", drv.ConnectionName())
	}
}

func TestFailoverDriverPendingSizeFallback(t *testing.T) {
	t.Parallel()

	d1 := drivers.NewNullDriver("d1")
	drv := drivers.NewFailoverDriver("failover", d1)

	n, err := drv.PendingSize(context.Background(), "q")

	if err != nil {
		t.Fatal(err)
	}

	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
}

func TestFailoverDriverDelayedSizeFallback(t *testing.T) {
	t.Parallel()

	d1 := drivers.NewNullDriver("d1")
	drv := drivers.NewFailoverDriver("failover", d1)

	n, err := drv.DelayedSize(context.Background(), "q")

	if err != nil {
		t.Fatal(err)
	}

	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
}

func TestFailoverDriverReservedSizeFallback(t *testing.T) {
	t.Parallel()

	d1 := drivers.NewNullDriver("d1")
	drv := drivers.NewFailoverDriver("failover", d1)

	n, err := drv.ReservedSize(context.Background(), "q")

	if err != nil {
		t.Fatal(err)
	}

	if n != 0 {
		t.Errorf("expected 0, got %d", n)
	}
}
