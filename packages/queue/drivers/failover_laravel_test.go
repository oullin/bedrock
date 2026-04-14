package drivers_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
	"github.com/bedrock/packages/queue/drivers"
	"github.com/bedrock/packages/queue/events"
)

// Port of Framework\Tests\Queue\FailoverQueueTest.
//
// Upstream's FailoverQueue takes a QueueManager, an events Dispatcher,
// and a list of connection names; it resolves each connection on
// demand via $manager->connection(). The Go FailoverDriver is given
// its sub-queues directly at construction time (resolution happens in
// the application wiring, not the driver), so this port substitutes
// the two sub-queues with small recording fakes. The behaviour under
// assertion — fall over on Push exception, dispatch a single event,
// land on the next backend — matches the PHP test line-for-line.

// --- fakes ------------------------------------------------------------

// throwingPushQueue fails every Push with a fixed error. Every other
// method is a no-op.
type throwingPushQueue struct {
	err  error
	name string
}

// capturingPushQueue accepts every Push and records the payloads for
// later assertion.
type capturingPushQueue struct {
	mu       sync.Mutex
	pushes   [][]byte
	name     string
	pushID   string
	pushCall int
}

// failoverRecorder is a small EventEmitter that collects every event
// for post-hoc assertion.
type failoverRecorder struct {
	mu     sync.Mutex
	events []any
}

func (q *throwingPushQueue) Push(context.Context, string, []byte) (string, error) {
	return "", q.err
}

func (q *throwingPushQueue) PushDelayed(context.Context, string, []byte, time.Duration) (string, error) {
	return "", q.err
}

func (q *throwingPushQueue) PushMultiple(context.Context, string, [][]byte) ([]string, error) {
	return nil, q.err
}

func (q *throwingPushQueue) Pop(context.Context, string) (queue.Job, error) {
	return nil, queue.ErrNoJob
}

func (q *throwingPushQueue) Size(context.Context, string) (int64, error) { return 0, nil }

func (q *throwingPushQueue) PendingSize(context.Context, string) (int64, error) { return 0, nil }

func (q *throwingPushQueue) DelayedSize(context.Context, string) (int64, error) { return 0, nil }

func (q *throwingPushQueue) ReservedSize(context.Context, string) (int64, error) { return 0, nil }

func (q *throwingPushQueue) ConnectionName() string { return q.name }

func (q *capturingPushQueue) Push(_ context.Context, _ string, payload []byte) (string, error) {
	q.mu.Lock()

	defer q.mu.Unlock()

	q.pushes = append(q.pushes, payload)
	q.pushCall++

	return q.pushID, nil
}

func (q *capturingPushQueue) PushDelayed(ctx context.Context, name string, payload []byte, _ time.Duration) (string, error) {
	return q.Push(ctx, name, payload)
}

func (q *capturingPushQueue) PushMultiple(ctx context.Context, name string, payloads [][]byte) ([]string, error) {
	ids := make([]string, len(payloads))

	for i, p := range payloads {
		_, _ = q.Push(ctx, name, p)

		ids[i] = ""
	}

	return ids, nil
}

func (q *capturingPushQueue) Pop(context.Context, string) (queue.Job, error) {
	return nil, queue.ErrNoJob
}

func (q *capturingPushQueue) Size(context.Context, string) (int64, error) { return 0, nil }

func (q *capturingPushQueue) PendingSize(context.Context, string) (int64, error) { return 0, nil }

func (q *capturingPushQueue) DelayedSize(context.Context, string) (int64, error) { return 0, nil }

func (q *capturingPushQueue) ReservedSize(context.Context, string) (int64, error) { return 0, nil }

func (q *capturingPushQueue) ConnectionName() string { return q.name }

func (r *failoverRecorder) Emit(event any) {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.events = append(r.events, event)
}

// Port of Framework\Tests\Queue\FailoverQueueTest::test_push_fails_over_on_exception
func TestPushFailsOverOnException(t *testing.T) {
	t.Parallel()

	redis := &throwingPushQueue{err: errors.New("error"), name: "redis"}
	sync := &capturingPushQueue{name: "sync", pushID: "sync-id"}
	rec := &failoverRecorder{}

	failover := drivers.NewFailoverDriver("failover", redis, sync).SetEmitter(rec)

	id, err := failover.Push(context.Background(), "default", []byte("some-job"))

	if err != nil {
		t.Fatalf("Push: %v", err)
	}

	if id != "sync-id" {
		t.Errorf("Push returned id %q, want sync-id", id)
	}

	if sync.pushCall != 1 {
		t.Errorf("sync.pushCall: got %d, want 1", sync.pushCall)
	}

	rec.mu.Lock()

	defer rec.mu.Unlock()

	if len(rec.events) != 1 {
		t.Fatalf("events: got %d, want 1", len(rec.events))
	}

	ev, ok := rec.events[0].(events.QueueFailedOver)

	if !ok {
		t.Fatalf("expected QueueFailedOver event, got %T", rec.events[0])
	}

	if ev.From != "redis" || ev.To != "sync" {
		t.Errorf("QueueFailedOver: got From=%q To=%q, want redis/sync", ev.From, ev.To)
	}

	if ev.Err == nil || ev.Err.Error() != "error" {
		t.Errorf("QueueFailedOver: Err=%v, want 'error'", ev.Err)
	}
}
