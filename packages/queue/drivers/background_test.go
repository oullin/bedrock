package drivers_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
	"github.com/bedrock/packages/queue/drivers"
)

// Use a non-existent command so spawn is a harmless no-op.

// Let goroutine start.

// NullDriver implements both contracts (returning nil/nil), so the
// wrapper passes through the nil result rather than ErrNotSupported.

// stubInspector is a queue.Queue that satisfies QueueNamer and
// JobInspector with canned return values. It lets us exercise the
// happy-path delegation in wrappers (Background, Deferred, Failover).
type stubInspector struct {
	connection string
	names      []string
	pending    []queue.InspectedJob
	delayed    []queue.InspectedJob
	reserved   []queue.InspectedJob
	err        error
}

// noInspectorInner satisfies only the base Queue contract so the
// wrappers fall to the ErrNotSupported branch.
type noInspectorInner struct{ connection string }

func TestBackgroundDriverPushDelegatesToInner(t *testing.T) {
	t.Parallel()

	client := newMockRedisClient()
	inner := drivers.NewRedisDriver(client, "redis")

	drv := drivers.NewBackgroundDriver("true", nil, inner, "bg")

	_, err := drv.Push(context.Background(), "default", []byte("payload"))

	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(10 * time.Millisecond)

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

func TestBackgroundDriverInspectionDelegatesErrNotSupported(t *testing.T) {
	t.Parallel()

	inner := drivers.NewNullDriver("null")
	drv := drivers.NewBackgroundDriver("true", nil, inner, "bg")
	ctx := context.Background()

	names, err := drv.QueueNames(ctx)

	if err != nil || names != nil {
		t.Errorf("QueueNames passthrough: got (%v, %v)", names, err)
	}

	pending, err := drv.PendingJobs(ctx, "default")

	if err != nil || pending != nil {
		t.Errorf("PendingJobs passthrough: got (%v, %v)", pending, err)
	}

	delayed, err := drv.DelayedJobs(ctx, "default")

	if err != nil || delayed != nil {
		t.Errorf("DelayedJobs passthrough: got (%v, %v)", delayed, err)
	}

	reserved, err := drv.ReservedJobs(ctx, "default")

	if err != nil || reserved != nil {
		t.Errorf("ReservedJobs passthrough: got (%v, %v)", reserved, err)
	}
}

func (s *stubInspector) Push(_ context.Context, _ string, _ []byte) (string, error) {
	return "", nil
}

func (s *stubInspector) PushDelayed(_ context.Context, _ string, _ []byte, _ time.Duration) (string, error) {
	return "", nil
}

func (s *stubInspector) PushMultiple(_ context.Context, _ string, payloads [][]byte) ([]string, error) {
	return make([]string, len(payloads)), nil
}

func (s *stubInspector) Pop(_ context.Context, _ string) (queue.Job, error) {
	return nil, queue.ErrNoJob
}

func (s *stubInspector) Size(_ context.Context, _ string) (int64, error)         { return 0, nil }
func (s *stubInspector) PendingSize(_ context.Context, _ string) (int64, error)  { return 0, nil }
func (s *stubInspector) DelayedSize(_ context.Context, _ string) (int64, error)  { return 0, nil }
func (s *stubInspector) ReservedSize(_ context.Context, _ string) (int64, error) { return 0, nil }
func (s *stubInspector) ConnectionName() string                                  { return s.connection }

func (s *stubInspector) QueueNames(_ context.Context) ([]string, error) {
	return s.names, s.err
}

func (s *stubInspector) PendingJobs(_ context.Context, _ string) ([]queue.InspectedJob, error) {
	return s.pending, s.err
}

func (s *stubInspector) DelayedJobs(_ context.Context, _ string) ([]queue.InspectedJob, error) {
	return s.delayed, s.err
}

func (s *stubInspector) ReservedJobs(_ context.Context, _ string) ([]queue.InspectedJob, error) {
	return s.reserved, s.err
}

func TestBackgroundDriverInspectionPropagatesResults(t *testing.T) {
	t.Parallel()

	inner := &stubInspector{
		connection: "inner",
		names:      []string{"default", "emails"},
		pending:    []queue.InspectedJob{{ID: 1, Queue: "default"}},
		delayed:    []queue.InspectedJob{{ID: 2, Queue: "default"}},
		reserved:   []queue.InspectedJob{{ID: 3, Queue: "default"}},
	}

	drv := drivers.NewBackgroundDriver("true", nil, inner, "bg")
	ctx := context.Background()

	names, err := drv.QueueNames(ctx)

	if err != nil || len(names) != 2 || names[0] != "default" {
		t.Errorf("QueueNames: got (%v, %v)", names, err)
	}

	pending, err := drv.PendingJobs(ctx, "default")

	if err != nil || len(pending) != 1 || pending[0].ID != 1 {
		t.Errorf("PendingJobs: got (%v, %v)", pending, err)
	}

	delayed, err := drv.DelayedJobs(ctx, "default")

	if err != nil || len(delayed) != 1 || delayed[0].ID != 2 {
		t.Errorf("DelayedJobs: got (%v, %v)", delayed, err)
	}

	reserved, err := drv.ReservedJobs(ctx, "default")

	if err != nil || len(reserved) != 1 || reserved[0].ID != 3 {
		t.Errorf("ReservedJobs: got (%v, %v)", reserved, err)
	}
}

func (n *noInspectorInner) Push(_ context.Context, _ string, _ []byte) (string, error) {
	return "", nil
}

func (n *noInspectorInner) PushDelayed(_ context.Context, _ string, _ []byte, _ time.Duration) (string, error) {
	return "", nil
}

func (n *noInspectorInner) PushMultiple(_ context.Context, _ string, payloads [][]byte) ([]string, error) {
	return make([]string, len(payloads)), nil
}

func (n *noInspectorInner) Pop(_ context.Context, _ string) (queue.Job, error) {
	return nil, queue.ErrNoJob
}

func (n *noInspectorInner) Size(_ context.Context, _ string) (int64, error)         { return 0, nil }
func (n *noInspectorInner) PendingSize(_ context.Context, _ string) (int64, error)  { return 0, nil }
func (n *noInspectorInner) DelayedSize(_ context.Context, _ string) (int64, error)  { return 0, nil }
func (n *noInspectorInner) ReservedSize(_ context.Context, _ string) (int64, error) { return 0, nil }
func (n *noInspectorInner) ConnectionName() string                                  { return n.connection }

func TestBackgroundDriverInspectionMissingContractReturnsErrNotSupported(t *testing.T) {
	t.Parallel()

	inner := &noInspectorInner{connection: "bare"}
	drv := drivers.NewBackgroundDriver("true", nil, inner, "bg")
	ctx := context.Background()

	if _, err := drv.QueueNames(ctx); !errors.Is(err, queue.ErrNotSupported) {
		t.Errorf("QueueNames: want ErrNotSupported, got %v", err)
	}

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
