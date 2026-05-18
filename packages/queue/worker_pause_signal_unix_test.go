//go:build unix

package queue_test

import (
	"context"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
)

// signalingQueue is a minimal Queue that returns ErrNoJob from every
// Pop so the worker stays in its idle/sleep loop, giving the test a
// stable window in which to send SIGUSR2 / SIGCONT and assert on the
// dispatched events.
type signalingQueue struct{}

type signalEventRecorder struct {
	mu     sync.Mutex
	events []any
}

func (signalingQueue) Push(context.Context, string, []byte) (string, error) {
	return "", nil
}

func (signalingQueue) PushDelayed(context.Context, string, []byte, time.Duration) (string, error) {
	return "", nil
}

func (signalingQueue) PushMultiple(_ context.Context, _ string, p [][]byte) ([]string, error) {
	return make([]string, len(p)), nil
}

func (signalingQueue) Pop(context.Context, string) (queue.Job, error)      { return nil, queue.ErrNoJob }
func (signalingQueue) Size(context.Context, string) (int64, error)         { return 0, nil }
func (signalingQueue) PendingSize(context.Context, string) (int64, error)  { return 0, nil }
func (signalingQueue) DelayedSize(context.Context, string) (int64, error)  { return 0, nil }
func (signalingQueue) ReservedSize(context.Context, string) (int64, error) { return 0, nil }
func (signalingQueue) ConnectionName() string                              { return "test-connection" }

func (r *signalEventRecorder) Emit(event any) {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.events = append(r.events, event)
}

func (r *signalEventRecorder) snapshot() []any {
	r.mu.Lock()

	defer r.mu.Unlock()

	out := make([]any, len(r.events))

	copy(out, r.events)

	return out
}

// waitFor polls predicate up to timeout, returning true once it
// returns true. Used to bridge the gap between sending a signal and
// the worker goroutine observing it.
func waitFor(timeout time.Duration, predicate func() bool) bool {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		if predicate() {
			return true
		}

		time.Sleep(10 * time.Millisecond)
	}

	return predicate()
}

func containsEvent[T any](events []any) (T, bool) {
	var zero T

	for _, e := range events {
		if typed, ok := e.(T); ok {
			return typed, true
		}
	}

	return zero, false
}

// TestWorkerEmitsPauseResumeEventsOnSignals exercises the SIGUSR2 →
// WorkerPausing and SIGCONT → WorkerResuming wiring end-to-end. Unix
// build tag only — these signals do not exist on Windows.
func TestWorkerEmitsPauseResumeEventsOnSignals(t *testing.T) {
	if os.Getenv("BEDROCK_SKIP_SIGNAL_TESTS") != "" {
		t.Skip("BEDROCK_SKIP_SIGNAL_TESTS set")
	}

	rec := &signalEventRecorder{}
	w := queue.NewWorker(signalingQueue{}, queue.HandlerFunc(func(context.Context, queue.Job) error { return nil }), rec, queue.WorkerOptions{
		Name:  "test-worker",
		Sleep: 20 * time.Millisecond,
	})

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	done := make(chan struct{})

	go func() {
		_ = w.Run(ctx, "default")

		close(done)
	}()

	// Wait for the WorkerStarting event so we know the signal handlers
	// have been installed before we send anything.
	if !waitFor(2*time.Second, func() bool {
		_, ok := containsEvent[queue.WorkerStarting](rec.snapshot())

		return ok
	}) {
		t.Fatalf("WorkerStarting was never emitted; signal handlers may not be wired")
	}

	if err := syscall.Kill(os.Getpid(), syscall.SIGUSR2); err != nil {
		t.Fatalf("kill SIGUSR2: %v", err)
	}

	if !waitFor(2*time.Second, func() bool {
		_, ok := containsEvent[queue.WorkerPausing](rec.snapshot())

		return ok
	}) {
		t.Fatalf("WorkerPausing was never emitted after SIGUSR2; events: %+v", rec.snapshot())
	}

	pausing, _ := containsEvent[queue.WorkerPausing](rec.snapshot())

	if pausing.ConnectionName != "test-connection" || pausing.Queue != "default" || pausing.WorkerName != "test-worker" {
		t.Errorf("WorkerPausing fields: got %+v", pausing)
	}

	if err := syscall.Kill(os.Getpid(), syscall.SIGCONT); err != nil {
		t.Fatalf("kill SIGCONT: %v", err)
	}

	if !waitFor(2*time.Second, func() bool {
		_, ok := containsEvent[queue.WorkerResuming](rec.snapshot())

		return ok
	}) {
		t.Fatalf("WorkerResuming was never emitted after SIGCONT")
	}

	resuming, _ := containsEvent[queue.WorkerResuming](rec.snapshot())

	if resuming.ConnectionName != "test-connection" || resuming.Queue != "default" || resuming.WorkerName != "test-worker" {
		t.Errorf("WorkerResuming fields: got %+v", resuming)
	}

	cancel()

	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatalf("worker did not exit after cancel")
	}
}
