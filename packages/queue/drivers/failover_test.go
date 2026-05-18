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

func TestFailoverDriverQueueNamesUnionAndDedupe(t *testing.T) {
	t.Parallel()

	d1 := &stubInspector{connection: "d1", names: []string{"a", "b"}}
	d2 := &stubInspector{connection: "d2", names: []string{"b", "c"}}
	drv := drivers.NewFailoverDriver("failover", d1, d2)

	names, err := drv.QueueNames(context.Background())
	if err != nil {
		t.Fatalf("QueueNames: %v", err)
	}

	if len(names) != 3 {
		t.Fatalf("got %d names, want 3: %v", len(names), names)
	}

	// First-seen order preserved across drivers.
	if names[0] != "a" || names[1] != "b" || names[2] != "c" {
		t.Errorf("ordering: got %v, want [a b c]", names)
	}
}

func TestFailoverDriverQueueNamesSkipsDriversWithoutContract(t *testing.T) {
	t.Parallel()

	bare := &noInspectorInner{connection: "bare"}
	d2 := &stubInspector{connection: "d2", names: []string{"x"}}
	drv := drivers.NewFailoverDriver("failover", bare, d2)

	names, err := drv.QueueNames(context.Background())
	if err != nil {
		t.Fatalf("QueueNames: %v", err)
	}

	if len(names) != 1 || names[0] != "x" {
		t.Errorf("got %v, want [x]", names)
	}
}

func TestFailoverDriverInspectionFallsThroughOnError(t *testing.T) {
	t.Parallel()

	d1 := &stubInspector{connection: "d1", err: errors.New("d1 broken")}
	d2 := &stubInspector{
		connection: "d2",
		pending:    []queue.InspectedJob{{ID: 99, Queue: "default"}},
		delayed:    []queue.InspectedJob{{ID: 100, Queue: "default"}},
		reserved:   []queue.InspectedJob{{ID: 101, Queue: "default"}},
	}
	drv := drivers.NewFailoverDriver("failover", d1, d2)
	ctx := context.Background()

	pending, err := drv.PendingJobs(ctx, "default")
	if err != nil || len(pending) != 1 || pending[0].ID != 99 {
		t.Errorf("PendingJobs fallback: got (%v, %v)", pending, err)
	}

	delayed, err := drv.DelayedJobs(ctx, "default")
	if err != nil || len(delayed) != 1 || delayed[0].ID != 100 {
		t.Errorf("DelayedJobs fallback: got (%v, %v)", delayed, err)
	}

	reserved, err := drv.ReservedJobs(ctx, "default")
	if err != nil || len(reserved) != 1 || reserved[0].ID != 101 {
		t.Errorf("ReservedJobs fallback: got (%v, %v)", reserved, err)
	}
}

func TestFailoverDriverInspectionAllBareReturnsErrNotSupported(t *testing.T) {
	t.Parallel()

	drv := drivers.NewFailoverDriver("failover",
		&noInspectorInner{connection: "a"},
		&noInspectorInner{connection: "b"},
	)
	ctx := context.Background()

	if _, err := drv.PendingJobs(ctx, "default"); !errors.Is(err, queue.ErrNotSupported) {
		t.Errorf("PendingJobs: want ErrNotSupported, got %v", err)
	}

	if _, err := drv.DelayedJobs(ctx, "default"); !errors.Is(err, queue.ErrNotSupported) {
		t.Errorf("DelayedJobs: want ErrNotSupported, got %v", err)
	}

	if _, err := drv.ReservedJobs(ctx, "default"); !errors.Is(err, queue.ErrNotSupported) {
		t.Errorf("ReservedJobs: want ErrNotSupported, got %v", err)
	}
}
