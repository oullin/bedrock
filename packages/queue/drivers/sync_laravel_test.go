package drivers_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/bedrock/packages/queue"
	"github.com/bedrock/packages/queue/drivers"
)

// Ports of Framework\Tests\Queue\QueueSyncQueueTest (partial: 4 of 8).
//
// The four after-commit tests require the Step 14 transaction-manager
// integration and are deferred. The four ported here exercise the
// observable Upstream behaviours that land in Step 7:
//
//   - Push fires the handler synchronously with the job + data visible.
//   - A thrown handler triggers the failure path: Fail, the handler's
//     Failed callback, and the 4-event lifecycle (JobProcessing,
//     JobExceptionOccurred, JobFailed, JobAttempted).
//   - The failed handler can read the payload via job.Payload() and
//     observes data injected by a createPayloadUsing hook.
//   - A successful handler can read the same payload and throw with a
//     value plucked from it — proving createPayloadUsing runs end to end.

// --- fakes ------------------------------------------------------------

// syncEventRecorder collects every emitted event for post-hoc assertion.
type syncEventRecorder struct {
	mu     sync.Mutex
	events []any
}

// --- handlers ---------------------------------------------------------

// recordingSyncHandler is the Go stand-in for SyncQueueTestHandler —
// records the (job, data) it was invoked with.
type recordingSyncHandler struct {
	mu     sync.Mutex
	called atomic.Bool
	job    queue.Job
	data   map[string]any
}

// failingSyncHandlerWithFailedCallback is the Go stand-in for
// FailingSyncQueueTestHandler — throws on Handle and sets a flag in
// its Failed callback.
type failingSyncHandlerWithFailedCallback struct {
	failedCalled atomic.Bool
}

// payloadReadingFailedHandler is the Go stand-in for FailingSyncQueueJob —
// throws on Handle and reads the payload from its Failed callback.
type payloadReadingFailedHandler struct {
	mu       sync.Mutex
	extra    string
	failedAt atomic.Bool
}

// payloadMessageHandler is the Go stand-in for SyncQueueJob — reads the
// payload and returns an error whose message is the extra value, so the
// test can assert the hook data flowed all the way through.
type payloadMessageHandler struct{}

func (r *syncEventRecorder) Emit(event any) {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.events = append(r.events, event)
}

func (r *syncEventRecorder) count() int {
	r.mu.Lock()

	defer r.mu.Unlock()

	return len(r.events)
}

func (h *recordingSyncHandler) Handle(_ context.Context, job queue.Job) error {
	h.mu.Lock()

	defer h.mu.Unlock()

	h.called.Store(true)
	h.job = job

	if p, err := queue.UnmarshalPayload(job.Payload()); err == nil {
		h.data = p.Data
	}

	return nil
}

func (h *failingSyncHandlerWithFailedCallback) Handle(_ context.Context, _ queue.Job) error {
	return errors.New("whoops")
}

func (h *failingSyncHandlerWithFailedCallback) Failed(_ context.Context, _ queue.Job, _ error) {
	h.failedCalled.Store(true)
}

func (h *payloadReadingFailedHandler) Handle(_ context.Context, _ queue.Job) error {
	return errors.New("logic")
}

func (h *payloadReadingFailedHandler) Failed(_ context.Context, job queue.Job, _ error) {
	h.mu.Lock()

	defer h.mu.Unlock()

	h.failedAt.Store(true)

	p, err := queue.UnmarshalPayload(job.Payload())

	if err != nil {
		return
	}

	if v, ok := p.Data["extra"].(string); ok {
		h.extra = v
	}
}

func (payloadMessageHandler) Handle(_ context.Context, job queue.Job) error {
	p, err := queue.UnmarshalPayload(job.Payload())

	if err != nil {
		return err
	}

	if v, ok := p.Data["extra"].(string); ok {
		return errors.New(v)
	}

	return errors.New("missing extra")
}

// helper: build a payload for the given "job name".
func buildSyncPayload(t *testing.T, connection, queueName, jobName string, data map[string]any) []byte {
	t.Helper()

	_, raw, err := queue.CreatePayloadFor(connection, queueName, jobName, data, queue.JobOptions{})

	if err != nil {
		t.Fatalf("CreatePayloadFor: %v", err)
	}

	return raw
}

// countByType returns how many events of type T the recorder holds.
func countByType[T any](r *syncEventRecorder) int {
	r.mu.Lock()

	defer r.mu.Unlock()

	n := 0

	for _, ev := range r.events {
		if _, ok := ev.(T); ok {
			n++
		}
	}

	return n
}

// --- ports ------------------------------------------------------------

// Port of Framework\Tests\Queue\QueueSyncQueueTest::testPushShouldFireJobInstantly
func TestPushShouldFireJobInstantly(t *testing.T) {
	t.Parallel()

	handler := &recordingSyncHandler{}
	drv := drivers.NewSyncDriver("sync", handler)

	payload := buildSyncPayload(t, "sync", "default", "SyncQueueTestHandler", map[string]any{"foo": "bar"})

	if _, err := drv.Push(context.Background(), "default", payload); err != nil {
		t.Fatalf("Push: %v", err)
	}

	if !handler.called.Load() {
		t.Fatal("handler was not invoked")
	}

	if handler.job == nil {
		t.Error("expected job instance to be captured")
	}

	if got, _ := handler.data["foo"].(string); got != "bar" {
		t.Errorf("data[foo]: got %v, want bar", handler.data["foo"])
	}
}

// Port of Framework\Tests\Queue\QueueSyncQueueTest::testFailedJobGetsHandledWhenAnExceptionIsThrown
func TestFailedJobGetsHandledWhenAnExceptionIsThrown(t *testing.T) {
	t.Parallel()

	handler := &failingSyncHandlerWithFailedCallback{}
	rec := &syncEventRecorder{}
	drv := drivers.NewSyncDriver("sync", handler).SetEmitter(rec)

	payload := buildSyncPayload(t, "sync", "default", "FailingSyncQueueTestHandler", map[string]any{"foo": "bar"})

	_, err := drv.Push(context.Background(), "default", payload)

	if err == nil || err.Error() != "whoops" {
		t.Fatalf("Push err: got %v, want 'whoops'", err)
	}

	if !handler.failedCalled.Load() {
		t.Error("handler.Failed was not invoked")
	}

	// Upstream dispatches exactly 4 events on the failure path:
	// JobProcessing, JobExceptionOccurred, JobFailed, JobAttempted.
	if got := rec.count(); got != 4 {
		t.Errorf("event count: got %d, want 4", got)
	}

	if n := countByType[queue.JobProcessing](rec); n != 1 {
		t.Errorf("JobProcessing count: got %d, want 1", n)
	}

	if n := countByType[queue.JobExceptionOccurred](rec); n != 1 {
		t.Errorf("JobExceptionOccurred count: got %d, want 1", n)
	}

	if n := countByType[queue.JobFailed](rec); n != 1 {
		t.Errorf("JobFailed count: got %d, want 1", n)
	}

	if n := countByType[queue.JobAttempted](rec); n != 1 {
		t.Errorf("JobAttempted count: got %d, want 1", n)
	}
}

// Port of Framework\Tests\Queue\QueueSyncQueueTest::testFailedJobHasAccessToJobInstance
//
// Not t.Parallel: mutates the global payload-hook list.
func TestFailedJobHasAccessToJobInstance(t *testing.T) {
	queue.ClearPayloadHooks()

	defer queue.ClearPayloadHooks()

	queue.CreatePayloadUsing(func(_, _ string, p *queue.Payload) {
		if p.Data == nil {
			p.Data = map[string]any{}
		}

		p.Data["extra"] = "extraValue"
	})

	handler := &payloadReadingFailedHandler{}
	drv := drivers.NewSyncDriver("sync", handler)

	payload := buildSyncPayload(t, "sync", "default", "FailingSyncQueueJob", nil)

	_, err := drv.Push(context.Background(), "default", payload)

	if err == nil || err.Error() != "logic" {
		t.Fatalf("Push err: got %v, want 'logic'", err)
	}

	if !handler.failedAt.Load() {
		t.Fatal("handler.Failed was not invoked")
	}

	handler.mu.Lock()

	defer handler.mu.Unlock()

	if handler.extra != "extraValue" {
		t.Errorf("handler.extra: got %q, want extraValue", handler.extra)
	}
}

// Port of Framework\Tests\Queue\QueueSyncQueueTest::testCreatesPayloadObject
//
// Not t.Parallel: mutates the global payload-hook list.
func TestCreatesPayloadObject(t *testing.T) {
	queue.ClearPayloadHooks()

	defer queue.ClearPayloadHooks()

	queue.CreatePayloadUsing(func(_, _ string, p *queue.Payload) {
		if p.Data == nil {
			p.Data = map[string]any{}
		}

		p.Data["extra"] = "extraValue"
	})

	drv := drivers.NewSyncDriver("sync", payloadMessageHandler{})

	payload := buildSyncPayload(t, "sync", "default", "SyncQueueJob", nil)

	_, err := drv.Push(context.Background(), "default", payload)

	if err == nil {
		t.Fatal("Push: expected error from handler")
	}

	if err.Error() != "extraValue" {
		t.Errorf("Push err: got %q, want extraValue", err.Error())
	}
}
