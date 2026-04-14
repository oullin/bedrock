package queue_test

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
)

// Partial port of Framework\Tests\Queue\QueueWorkerTest (8 of 25).
//
// The Upstream tests drive a Worker that takes a QueueManager and a
// connection/queue-name pair, supports comma-separated priority queues,
// installs an exception handler, and exposes runNextJob / daemon.
// The Go worker currently takes a single Queue instance; a full
// rewrite is tracked as Step 8 proper. This file ports the subset of
// tests that can be expressed against the additive surface landed in
// Step 8a:
//
//   - RunNextJob (single-iteration primitive)
//   - MemoryExceeded method
//   - WorkerStarting / WorkerStopping / JobPopping / JobPopped events
//   - SleepFunc hook + SleptFor accessor
//
// Deferred until the full rewrite:
//   - priority queues (comma-separated queue names)
//   - exception handler integration
//   - maintenance-mode hook
//   - daemon status code (WorkerStopReason)
//   - every *PickedJobUsingCustomCallbacks / *StopsWithLostConnectionReason test

// --- fixtures ---------------------------------------------------------

// workerUpstreamJob is a recording Job fixture — fires or errors on
// demand, and tracks delete/release/fail state. It is the Go stand-in
// for Upstream's WorkerFakeJob and supports the full mutable attempts /
// maxTries / backoff / retryUntil surface the retry-semantics tests
// exercise.
type workerUpstreamJob struct {
	fired         atomic.Bool
	attempts      int
	maxTries      int
	maxExceptions int
	backoffs      []time.Duration
	retryUntil    *time.Time
	deleted       bool
	released      bool
	releaseDelay  time.Duration
	failedWith    error
	handler       func(j *workerUpstreamJob) error
	connection    string
	queueName     string
}

// Fail mirrors Upstream's WorkerFakeJob::fail($e): stores the error in
// failedWith and marks the job as deleted.

// workerFakeQueue is a deterministic test queue: it returns the
// pre-staged jobs on Pop and then ErrNoJob forever. The Go stand-in
// for a single Upstream WorkerFakeManager/WorkerFakeQueue pair.
type workerFakeQueue struct {
	mu         sync.Mutex
	jobs       []queue.Job
	popped     int
	connection string
}

// workerEventRecorder collects emitted events for post-hoc assertions.
type workerEventRecorder struct {
	mu     sync.Mutex
	events []any
}

// --- ports ------------------------------------------------------------

// Port of Framework\Tests\Queue\QueueWorkerTest::testWorkerMemoryExceededWhenMemoryIsZero

// Port of Framework\Tests\Queue\QueueWorkerTest::testWorkerMemoryExceededWhenMemoryGreaterThanZero

// The Go test binary itself has allocated well over 1 MiB by the
// time this runs, so a 1 MiB cap must register as exceeded — the
// same logic as Upstream's test: memory_get_usage() > 1.

// Port of Framework\Tests\Queue\QueueWorkerTest::testWorkerMemoryExceededWhenMemoryIsNegative

// Port of Framework\Tests\Queue\QueueWorkerTest::testJobCanBeFired

// Port of Framework\Tests\Queue\QueueWorkerTest::testJobPoppingEvent

// Assert exactly one JobPopping event carrying the connection name.

// Port of Framework\Tests\Queue\QueueWorkerTest::testWorkerCanWorkUntilQueueIsEmpty

// Port of Framework\Tests\Queue\QueueWorkerTest::testWorkerSleepsWhenQueueIsEmpty

// no jobs

// Override SleepFunc so the test does not actually wait 5 seconds.
// The Go adaptation of Upstream's $worker->sleptFor override.

// SleptFor accessor reflects the same value.

// Port of Framework\Tests\Queue\QueueWorkerTest::testWorkerStartingIsDispatched

// --- recording exception reporter ------------------------------------

type workerReporter struct {
	mu     sync.Mutex
	errors []error
}

// newWorkerRun runs a single RunNextJob against the supplied job and
// returns the populated recorder + reporter + worker for assertions.

// passthroughHandler is the worker-facing Handler shim that calls the
// job's own Fire method. Matches how Upstream's runNextJob dispatches
// to the job — the Go worker's Handler takes the job and forwards to
// its embedded handler.

// --- retry/fail-semantics ports ---------------------------------------

// Port of Framework\Tests\Queue\QueueWorkerTest::testJobIsReleasedOnException

// Port of Framework\Tests\Queue\QueueWorkerTest::testExceptionIsNotReportedIfReportJobExceptionsIsDisabled

// The JobExceptionOccurred event is still dispatched — only the
// reporter path is suppressed.

// Port of Framework\Tests\Queue\QueueWorkerTest::testJobIsNotReleasedIfItHasExceededMaxAttempts

// Simulate the "pop incremented attempts" behaviour.

// Port of Framework\Tests\Queue\QueueWorkerTest::testJobIsNotReleasedIfItHasExpired

// Upstream uses Carbon::setTestNow to advance past retryUntil. Go
// has no equivalent; set retryUntil to a past timestamp instead.

// Port of Framework\Tests\Queue\QueueWorkerTest::testJobIsFailedIfItHasAlreadyExceededMaxAttempts

// Pre-fire check should have short-circuited — handler never ran.

// Port of Framework\Tests\Queue\QueueWorkerTest::testJobIsFailedIfItHasAlreadyExpired

// Port of Framework\Tests\Queue\QueueWorkerTest::testJobBasedMaxRetries
//
// Upstream asserts deleted==false on this success path, but that's a
// PHP artifact of CallQueuedHandler owning the delete call. The Go
// worker no longer auto-deletes either, so the assertion is preserved.

// opts.MaxTries=1 should be ignored because job.MaxTries()=10 wins.

// Port of Framework\Tests\Queue\QueueWorkerTest::testJobBasedFailedDelay

// opts.Backoff=3s should be ignored because job.Backoff() is set.

// Port of Framework\Tests\Queue\QueueWorkerTest::testJobDoesNotFireIfDeleted

// Delete the job before the worker sees it.

// Upstream dispatches JobProcessed exactly once for a pre-deleted job.

// --- multi-queue and pop-error fixtures ------------------------------

// workerMultiQueueFake is a workerFakeQueue variant that tracks jobs
// per queue name, so priority-queue tests can stage jobs across
// multiple buckets.
type workerMultiQueueFake struct {
	mu         sync.Mutex
	buckets    map[string][]queue.Job
	connection string
}

// brokenPopQueue errors on every Pop with a fixed error. Used by the
// exception-reporter port.
type brokenPopQueue struct {
	err        error
	connection string
}

func (j *workerUpstreamJob) UUID() string              { return "" }
func (j *workerUpstreamJob) GetJobID() string          { return "" }
func (j *workerUpstreamJob) Payload() []byte           { return nil }
func (j *workerUpstreamJob) Attempts() int             { return j.attempts }
func (j *workerUpstreamJob) MaxTries() int             { return j.maxTries }
func (j *workerUpstreamJob) MaxExceptions() int        { return j.maxExceptions }
func (j *workerUpstreamJob) Timeout() time.Duration    { return 0 }
func (j *workerUpstreamJob) Backoff() []time.Duration  { return j.backoffs }
func (j *workerUpstreamJob) RetryUntil() *time.Time    { return j.retryUntil }
func (j *workerUpstreamJob) GetQueue() string          { return j.queueName }
func (j *workerUpstreamJob) GetConnectionName() string { return j.connection }
func (j *workerUpstreamJob) IsDeleted() bool           { return j.deleted }
func (j *workerUpstreamJob) IsReleased() bool          { return j.released }
func (j *workerUpstreamJob) HasFailed() bool           { return j.failedWith != nil }

func (j *workerUpstreamJob) Fire(context.Context) error {
	j.fired.Store(true)

	if j.handler != nil {
		return j.handler(j)
	}

	return nil
}

func (j *workerUpstreamJob) Delete() error {
	j.deleted = true

	return nil
}

func (j *workerUpstreamJob) Release(d time.Duration) error {
	j.released = true
	j.releaseDelay = d

	return nil
}

func (j *workerUpstreamJob) Fail(err error) error {
	j.failedWith = err
	j.deleted = true

	return nil
}

func (j *workerUpstreamJob) MarkAsFailed(err error) error { return j.Fail(err) }

func (q *workerFakeQueue) Push(context.Context, string, []byte) (string, error) {
	return "", nil
}

func (q *workerFakeQueue) PushDelayed(context.Context, string, []byte, time.Duration) (string, error) {
	return "", nil
}

func (q *workerFakeQueue) PushMultiple(context.Context, string, [][]byte) ([]string, error) {
	return nil, nil
}

func (q *workerFakeQueue) Pop(context.Context, string) (queue.Job, error) {
	q.mu.Lock()

	defer q.mu.Unlock()

	if q.popped >= len(q.jobs) {
		return nil, queue.ErrNoJob
	}

	job := q.jobs[q.popped]
	q.popped++

	return job, nil
}

func (q *workerFakeQueue) Size(context.Context, string) (int64, error) {
	q.mu.Lock()

	defer q.mu.Unlock()

	return int64(len(q.jobs) - q.popped), nil
}

func (q *workerFakeQueue) PendingSize(ctx context.Context, name string) (int64, error) {
	return q.Size(ctx, name)
}

func (q *workerFakeQueue) DelayedSize(context.Context, string) (int64, error)  { return 0, nil }
func (q *workerFakeQueue) ReservedSize(context.Context, string) (int64, error) { return 0, nil }
func (q *workerFakeQueue) ConnectionName() string                              { return q.connection }

func (r *workerEventRecorder) Emit(event any) {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.events = append(r.events, event)
}

func (r *workerEventRecorder) countByType(pred func(any) bool) int {
	r.mu.Lock()

	defer r.mu.Unlock()

	n := 0

	for _, e := range r.events {
		if pred(e) {
			n++
		}
	}

	return n
}

func TestWorkerMemoryExceededWhenMemoryIsZero(t *testing.T) {
	t.Parallel()

	w := queue.NewWorker(&workerFakeQueue{connection: "default"}, nil, nil, queue.WorkerOptions{})

	if w.MemoryExceeded(0) {
		t.Error("MemoryExceeded(0): got true, want false")
	}
}

func TestWorkerMemoryExceededWhenMemoryGreaterThanZero(t *testing.T) {
	t.Parallel()

	w := queue.NewWorker(&workerFakeQueue{connection: "default"}, nil, nil, queue.WorkerOptions{})

	if !w.MemoryExceeded(1) {
		t.Error("MemoryExceeded(1): got false, want true")
	}
}

func TestWorkerMemoryExceededWhenMemoryIsNegative(t *testing.T) {
	t.Parallel()

	w := queue.NewWorker(&workerFakeQueue{connection: "default"}, nil, nil, queue.WorkerOptions{})

	if w.MemoryExceeded(-1) {
		t.Error("MemoryExceeded(-1): got true, want false")
	}
}

func TestJobCanBeFired(t *testing.T) {
	t.Parallel()

	job := &workerUpstreamJob{connection: "default", queueName: "queue"}
	q := &workerFakeQueue{connection: "default", jobs: []queue.Job{job}}
	rec := &workerEventRecorder{}

	handler := queue.HandlerFunc(func(_ context.Context, j queue.Job) error {
		return j.Fire(context.Background())
	})

	w := queue.NewWorker(q, handler, rec, queue.WorkerOptions{})

	if err := w.RunNextJob(context.Background(), "queue"); err != nil {
		t.Fatalf("RunNextJob: %v", err)
	}

	if !job.fired.Load() {
		t.Error("expected job.Fire to have been called")
	}

	if got := rec.countByType(func(e any) bool { _, ok := e.(queue.JobPopping); return ok }); got != 1 {
		t.Errorf("JobPopping: got %d, want 1", got)
	}

	if got := rec.countByType(func(e any) bool { _, ok := e.(queue.JobPopped); return ok }); got != 1 {
		t.Errorf("JobPopped: got %d, want 1", got)
	}

	if got := rec.countByType(func(e any) bool { _, ok := e.(queue.JobProcessing); return ok }); got != 1 {
		t.Errorf("JobProcessing: got %d, want 1", got)
	}

	if got := rec.countByType(func(e any) bool { _, ok := e.(queue.JobProcessed); return ok }); got != 1 {
		t.Errorf("JobProcessed: got %d, want 1", got)
	}
}

func TestJobPoppingEvent(t *testing.T) {
	t.Parallel()

	job := &workerUpstreamJob{connection: "default", queueName: "queue"}
	q := &workerFakeQueue{connection: "default", jobs: []queue.Job{job}}
	rec := &workerEventRecorder{}

	handler := queue.HandlerFunc(func(_ context.Context, j queue.Job) error {
		return j.Fire(context.Background())
	})

	w := queue.NewWorker(q, handler, rec, queue.WorkerOptions{})

	if err := w.RunNextJob(context.Background(), "queue"); err != nil {
		t.Fatalf("RunNextJob: %v", err)
	}

	var popping *queue.JobPopping

	rec.mu.Lock()

	for _, e := range rec.events {
		if p, ok := e.(queue.JobPopping); ok {
			pCopy := p
			popping = &pCopy

			break
		}
	}

	rec.mu.Unlock()

	if popping == nil {
		t.Fatal("expected a JobPopping event")
	}

	if popping.ConnectionName != "default" {
		t.Errorf("JobPopping.ConnectionName: got %q, want default", popping.ConnectionName)
	}
}

func TestWorkerCanWorkUntilQueueIsEmpty(t *testing.T) {
	t.Parallel()

	j1 := &workerUpstreamJob{connection: "default", queueName: "queue"}
	j2 := &workerUpstreamJob{connection: "default", queueName: "queue"}
	q := &workerFakeQueue{connection: "default", jobs: []queue.Job{j1, j2}}
	rec := &workerEventRecorder{}

	handler := queue.HandlerFunc(func(_ context.Context, j queue.Job) error {
		return j.Fire(context.Background())
	})

	w := queue.NewWorker(q, handler, rec, queue.WorkerOptions{StopOnEmpty: true, Sleep: time.Millisecond})

	if err := w.Run(context.Background(), "queue"); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if !j1.fired.Load() || !j2.fired.Load() {
		t.Errorf("expected both jobs fired, got j1=%v j2=%v", j1.fired.Load(), j2.fired.Load())
	}

	if got := rec.countByType(func(e any) bool { _, ok := e.(queue.JobProcessing); return ok }); got != 2 {
		t.Errorf("JobProcessing: got %d, want 2", got)
	}

	if got := rec.countByType(func(e any) bool { _, ok := e.(queue.JobProcessed); return ok }); got != 2 {
		t.Errorf("JobProcessed: got %d, want 2", got)
	}
}

func TestWorkerSleepsWhenQueueIsEmpty(t *testing.T) {
	t.Parallel()

	q := &workerFakeQueue{connection: "default"}
	w := queue.NewWorker(q, nil, nil, queue.WorkerOptions{Sleep: 5 * time.Second})

	var sleptFor time.Duration

	w.SleepFunc = func(_ context.Context, d time.Duration) { sleptFor = d }

	if err := w.RunNextJob(context.Background(), "queue"); err != nil {
		t.Fatalf("RunNextJob: %v", err)
	}

	if sleptFor != 5*time.Second {
		t.Errorf("sleptFor: got %s, want 5s", sleptFor)
	}

	if w.SleptFor() != 5*time.Second {
		t.Errorf("SleptFor(): got %s, want 5s", w.SleptFor())
	}
}

func TestWorkerStartingIsDispatched(t *testing.T) {
	t.Parallel()

	q := &workerFakeQueue{connection: "default"}
	rec := &workerEventRecorder{}

	w := queue.NewWorker(q, nil, rec, queue.WorkerOptions{
		Name:        "test-worker",
		Sleep:       time.Millisecond,
		StopOnEmpty: true,
	})

	if err := w.Run(context.Background(), "queue"); err != nil {
		t.Fatalf("Run: %v", err)
	}

	var starting *queue.WorkerStarting

	rec.mu.Lock()

	for _, e := range rec.events {
		if s, ok := e.(queue.WorkerStarting); ok {
			sCopy := s
			starting = &sCopy

			break
		}
	}

	rec.mu.Unlock()

	if starting == nil {
		t.Fatal("expected a WorkerStarting event")
	}

	if starting.ConnectionName != "default" {
		t.Errorf("WorkerStarting.ConnectionName: got %q, want default", starting.ConnectionName)
	}

	if starting.Queue != "queue" {
		t.Errorf("WorkerStarting.Queue: got %q, want queue", starting.Queue)
	}

	if starting.WorkerName != "test-worker" {
		t.Errorf("WorkerStarting.WorkerName: got %q, want test-worker", starting.WorkerName)
	}
}

func (r *workerReporter) ReportException(err error) {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.errors = append(r.errors, err)
}

func (r *workerReporter) count() int {
	r.mu.Lock()

	defer r.mu.Unlock()

	return len(r.errors)
}

func newWorkerRun(
	t *testing.T,
	job *workerUpstreamJob,
	handler queue.Handler,
	opts queue.WorkerOptions,
) (*queue.Worker, *workerEventRecorder, *workerReporter) {
	t.Helper()

	q := &workerFakeQueue{connection: "default", jobs: []queue.Job{job}}
	rec := &workerEventRecorder{}
	rep := &workerReporter{}

	w := queue.NewWorker(q, handler, rec, opts)
	w.ExceptionReporter = rep

	if err := w.RunNextJob(context.Background(), "queue"); err != nil {
		t.Fatalf("RunNextJob: %v", err)
	}

	return w, rec, rep
}

func passthroughHandler() queue.Handler {
	return queue.HandlerFunc(func(ctx context.Context, j queue.Job) error {
		return j.Fire(ctx)
	})
}

func TestJobIsReleasedOnException(t *testing.T) {
	t.Parallel()

	theErr := errors.New("boom")

	job := &workerUpstreamJob{
		connection: "default",
		queueName:  "queue",
		handler: func(j *workerUpstreamJob) error {
			return theErr
		},
	}

	_, rec, rep := newWorkerRun(t, job, passthroughHandler(), queue.WorkerOptions{Backoff: 10 * time.Second})

	if !job.released {
		t.Error("expected job to be released")
	}

	if job.releaseDelay != 10*time.Second {
		t.Errorf("releaseDelay: got %s, want 10s", job.releaseDelay)
	}

	if job.deleted {
		t.Error("expected job NOT to be deleted on release path")
	}

	if rep.count() != 1 {
		t.Errorf("reporter calls: got %d, want 1", rep.count())
	}

	if n := rec.countByType(func(e any) bool { _, ok := e.(queue.JobExceptionOccurred); return ok }); n != 1 {
		t.Errorf("JobExceptionOccurred: got %d, want 1", n)
	}

	if n := rec.countByType(func(e any) bool { _, ok := e.(queue.JobProcessed); return ok }); n != 0 {
		t.Errorf("JobProcessed: got %d, want 0", n)
	}
}

func TestExceptionIsNotReportedIfReportJobExceptionsIsDisabled(t *testing.T) {
	t.Parallel()

	theErr := errors.New("boom")

	job := &workerUpstreamJob{
		connection: "default",
		queueName:  "queue",
		handler: func(*workerUpstreamJob) error {
			return theErr
		},
	}

	q := &workerFakeQueue{connection: "default", jobs: []queue.Job{job}}
	rec := &workerEventRecorder{}
	rep := &workerReporter{}

	w := queue.NewWorker(q, passthroughHandler(), rec, queue.WorkerOptions{Backoff: 10 * time.Second})
	w.ExceptionReporter = rep
	w.ReportJobExceptions = false

	if err := w.RunNextJob(context.Background(), "queue"); err != nil {
		t.Fatalf("RunNextJob: %v", err)
	}

	if rep.count() != 0 {
		t.Errorf("reporter calls: got %d, want 0 (reporting disabled)", rep.count())
	}

	if n := rec.countByType(func(e any) bool { _, ok := e.(queue.JobExceptionOccurred); return ok }); n != 1 {
		t.Errorf("JobExceptionOccurred: got %d, want 1", n)
	}
}

func TestJobIsNotReleasedIfItHasExceededMaxAttempts(t *testing.T) {
	t.Parallel()

	theErr := errors.New("boom")

	job := &workerUpstreamJob{
		connection: "default",
		queueName:  "queue",
		attempts:   1,
		handler: func(j *workerUpstreamJob) error {

			j.attempts++

			return theErr
		},
	}

	_, rec, rep := newWorkerRun(t, job, passthroughHandler(), queue.WorkerOptions{MaxTries: 1})

	if job.released {
		t.Error("expected job NOT to be released (exhausted)")
	}

	if !job.deleted {
		t.Error("expected job to be deleted via Fail")
	}

	if job.failedWith == nil || job.failedWith.Error() != "boom" {
		t.Errorf("failedWith: got %v, want 'boom'", job.failedWith)
	}

	if rep.count() != 1 {
		t.Errorf("reporter calls: got %d, want 1", rep.count())
	}

	if n := rec.countByType(func(e any) bool { _, ok := e.(queue.JobExceptionOccurred); return ok }); n != 1 {
		t.Errorf("JobExceptionOccurred: got %d, want 1", n)
	}

	if n := rec.countByType(func(e any) bool { _, ok := e.(queue.JobProcessed); return ok }); n != 0 {
		t.Errorf("JobProcessed: got %d, want 0", n)
	}
}

func TestJobIsNotReleasedIfItHasExpired(t *testing.T) {
	t.Parallel()

	theErr := errors.New("boom")

	past := time.Now().Add(-time.Second)

	job := &workerUpstreamJob{
		connection: "default",
		queueName:  "queue",
		retryUntil: &past,
		handler: func(j *workerUpstreamJob) error {
			j.attempts++

			return theErr
		},
	}

	_, _, _ = newWorkerRun(t, job, passthroughHandler(), queue.WorkerOptions{})

	if job.released {
		t.Error("expected job NOT to be released (expired)")
	}

	if !job.deleted {
		t.Error("expected job to be deleted via Fail")
	}

	if job.failedWith == nil {
		t.Error("expected failedWith to be set")
	}
}

func TestJobIsFailedIfItHasAlreadyExceededMaxAttempts(t *testing.T) {
	t.Parallel()

	job := &workerUpstreamJob{
		connection: "default",
		queueName:  "queue",
		attempts:   2,
		handler: func(j *workerUpstreamJob) error {
			j.attempts++

			return nil
		},
	}

	_, rec, rep := newWorkerRun(t, job, passthroughHandler(), queue.WorkerOptions{MaxTries: 1})

	if job.fired.Load() {
		t.Error("expected handler NOT to fire (pre-fire exhaustion)")
	}

	if !job.deleted {
		t.Error("expected job to be deleted via Fail")
	}

	var maxErr *queue.MaxAttemptsExceededError

	if !errors.As(job.failedWith, &maxErr) {
		t.Errorf("failedWith: got %T (%v), want *MaxAttemptsExceededError", job.failedWith, job.failedWith)
	}

	if rep.count() != 1 {
		t.Errorf("reporter calls: got %d, want 1", rep.count())
	}

	if n := rec.countByType(func(e any) bool { _, ok := e.(queue.JobExceptionOccurred); return ok }); n != 1 {
		t.Errorf("JobExceptionOccurred: got %d, want 1", n)
	}
}

func TestJobIsFailedIfItHasAlreadyExpired(t *testing.T) {
	t.Parallel()

	past := time.Now().Add(-time.Second)

	job := &workerUpstreamJob{
		connection: "default",
		queueName:  "queue",
		attempts:   1,
		retryUntil: &past,
		handler: func(j *workerUpstreamJob) error {
			j.attempts++

			return nil
		},
	}

	_, _, _ = newWorkerRun(t, job, passthroughHandler(), queue.WorkerOptions{})

	if job.fired.Load() {
		t.Error("expected handler NOT to fire (pre-fire expired)")
	}

	if !job.deleted {
		t.Error("expected job to be deleted via Fail")
	}

	var maxErr *queue.MaxAttemptsExceededError

	if !errors.As(job.failedWith, &maxErr) {
		t.Errorf("failedWith: got %T, want *MaxAttemptsExceededError", job.failedWith)
	}
}

func TestJobBasedMaxRetries(t *testing.T) {
	t.Parallel()

	job := &workerUpstreamJob{
		connection: "default",
		queueName:  "queue",
		attempts:   2,
		maxTries:   10,
		handler: func(j *workerUpstreamJob) error {
			j.attempts++

			return nil
		},
	}

	_, _, _ = newWorkerRun(t, job, passthroughHandler(), queue.WorkerOptions{MaxTries: 1})

	if !job.fired.Load() {
		t.Error("expected handler to fire (job's own maxTries overrides opts)")
	}

	if job.failedWith != nil {
		t.Errorf("failedWith: got %v, want nil", job.failedWith)
	}

	if job.released {
		t.Error("expected job NOT to be released (success path)")
	}

	if job.deleted {
		t.Error("expected job NOT to be deleted (worker no longer auto-deletes)")
	}
}

func TestJobBasedFailedDelay(t *testing.T) {
	t.Parallel()

	job := &workerUpstreamJob{
		connection: "default",
		queueName:  "queue",
		attempts:   1,
		backoffs:   []time.Duration{10 * time.Second},
		handler: func(j *workerUpstreamJob) error {
			return errors.New("something went wrong")
		},
	}

	_, _, _ = newWorkerRun(t, job, passthroughHandler(), queue.WorkerOptions{
		Backoff:  3 * time.Second,
		MaxTries: 0,
	})

	if !job.released {
		t.Error("expected job to be released")
	}

	if job.releaseDelay != 10*time.Second {
		t.Errorf("releaseDelay: got %s, want 10s", job.releaseDelay)
	}
}

func TestJobDoesNotFireIfDeleted(t *testing.T) {
	t.Parallel()

	job := &workerUpstreamJob{
		connection: "default",
		queueName:  "queue",
		handler: func(*workerUpstreamJob) error {
			return nil
		},
	}

	_ = job.Delete()

	_, rec, _ := newWorkerRun(t, job, passthroughHandler(), queue.WorkerOptions{})

	if job.fired.Load() {
		t.Error("expected handler NOT to fire (pre-deleted)")
	}

	if !job.IsDeleted() {
		t.Error("expected job to remain deleted")
	}

	if job.HasFailed() {
		t.Error("expected job not failed")
	}

	if job.IsReleased() {
		t.Error("expected job not released")
	}

	if n := rec.countByType(func(e any) bool { _, ok := e.(queue.JobProcessed); return ok }); n != 1 {
		t.Errorf("JobProcessed: got %d, want 1", n)
	}
}

func newWorkerMultiQueueFake(connection string, buckets map[string][]queue.Job) *workerMultiQueueFake {
	return &workerMultiQueueFake{connection: connection, buckets: buckets}
}

func (q *workerMultiQueueFake) Push(context.Context, string, []byte) (string, error) {
	return "", nil
}

func (q *workerMultiQueueFake) PushDelayed(context.Context, string, []byte, time.Duration) (string, error) {
	return "", nil
}

func (q *workerMultiQueueFake) PushMultiple(context.Context, string, [][]byte) ([]string, error) {
	return nil, nil
}

func (q *workerMultiQueueFake) Pop(_ context.Context, queueName string) (queue.Job, error) {
	q.mu.Lock()

	defer q.mu.Unlock()

	bucket := q.buckets[queueName]

	if len(bucket) == 0 {
		return nil, queue.ErrNoJob
	}

	job := bucket[0]
	q.buckets[queueName] = bucket[1:]

	return job, nil
}

func (q *workerMultiQueueFake) Size(context.Context, string) (int64, error)         { return 0, nil }
func (q *workerMultiQueueFake) PendingSize(context.Context, string) (int64, error)  { return 0, nil }
func (q *workerMultiQueueFake) DelayedSize(context.Context, string) (int64, error)  { return 0, nil }
func (q *workerMultiQueueFake) ReservedSize(context.Context, string) (int64, error) { return 0, nil }
func (q *workerMultiQueueFake) ConnectionName() string                              { return q.connection }

func (q *brokenPopQueue) Push(context.Context, string, []byte) (string, error) {
	return "", nil
}

func (q *brokenPopQueue) PushDelayed(context.Context, string, []byte, time.Duration) (string, error) {
	return "", nil
}

func (q *brokenPopQueue) PushMultiple(context.Context, string, [][]byte) ([]string, error) {
	return nil, nil
}

func (q *brokenPopQueue) Pop(context.Context, string) (queue.Job, error) { return nil, q.err }
func (q *brokenPopQueue) Size(context.Context, string) (int64, error)    { return 0, nil }
func (q *brokenPopQueue) PendingSize(context.Context, string) (int64, error) {
	return 0, nil
}
func (q *brokenPopQueue) DelayedSize(context.Context, string) (int64, error)  { return 0, nil }
func (q *brokenPopQueue) ReservedSize(context.Context, string) (int64, error) { return 0, nil }
func (q *brokenPopQueue) ConnectionName() string                              { return q.connection }

// --- Step 8c ports ---------------------------------------------------

// Port of Framework\Tests\Queue\QueueWorkerTest::testJobCanBeFiredBasedOnPriority
func TestJobCanBeFiredBasedOnPriority(t *testing.T) {
	t.Parallel()

	high1 := &workerUpstreamJob{connection: "default", queueName: "high"}
	high2 := &workerUpstreamJob{connection: "default", queueName: "high"}
	low := &workerUpstreamJob{connection: "default", queueName: "low"}

	q := newWorkerMultiQueueFake("default", map[string][]queue.Job{
		"high": {high1, high2},
		"low":  {low},
	})

	handler := passthroughHandler()
	rec := &workerEventRecorder{}

	w := queue.NewWorker(q, handler, rec, queue.WorkerOptions{})

	// 1st RunNextJob: high has 2 jobs → pop high1.
	if err := w.RunNextJob(context.Background(), "high,low"); err != nil {
		t.Fatalf("RunNextJob 1: %v", err)
	}

	if !high1.fired.Load() {
		t.Error("high1: expected fired")
	}

	if high2.fired.Load() || low.fired.Load() {
		t.Error("only high1 should have fired on iteration 1")
	}

	// 2nd RunNextJob: high still has 1 job → pop high2.
	if err := w.RunNextJob(context.Background(), "high,low"); err != nil {
		t.Fatalf("RunNextJob 2: %v", err)
	}

	if !high2.fired.Load() {
		t.Error("high2: expected fired")
	}

	if low.fired.Load() {
		t.Error("low should not fire yet")
	}

	// 3rd RunNextJob: high empty, fall through to low.
	if err := w.RunNextJob(context.Background(), "high,low"); err != nil {
		t.Fatalf("RunNextJob 3: %v", err)
	}

	if !low.fired.Load() {
		t.Error("low: expected fired")
	}
}

// Port of Framework\Tests\Queue\QueueWorkerTest::testExceptionIsReportedIfConnectionThrowsExceptionOnJobPop
func TestExceptionIsReportedIfConnectionThrowsExceptionOnJobPop(t *testing.T) {
	t.Parallel()

	theErr := errors.New("runtime failure")

	q := &brokenPopQueue{err: theErr, connection: "default"}
	rep := &workerReporter{}

	w := queue.NewWorker(q, passthroughHandler(), nil, queue.WorkerOptions{Sleep: time.Millisecond})
	w.ExceptionReporter = rep
	// Short-circuit the sleep so the test does not stall on the
	// empty-queue backoff after the reported error.
	w.SleepFunc = func(context.Context, time.Duration) {}

	if err := w.RunNextJob(context.Background(), "queue"); err != nil {
		t.Fatalf("RunNextJob: %v", err)
	}

	if rep.count() != 1 {
		t.Errorf("reporter calls: got %d, want 1", rep.count())
	}

	rep.mu.Lock()

	if !errors.Is(rep.errors[0], theErr) {
		t.Errorf("reported error: got %v, want %v", rep.errors[0], theErr)
	}

	rep.mu.Unlock()
}

// Port of Framework\Tests\Queue\QueueWorkerTest::testJobRunsIfAppIsNotInMaintenanceMode
//
// Upstream advances a pre-canned maintenance-flag list [false, true]
// and then throws a LoopBreakerException. The Go port uses a bounded
// MaxTime to break out of the loop and asserts the same observable
// outcome: the first iteration (not maintenance) processes the first
// job, subsequent iterations (maintenance true) never process the
// second.
func TestJobRunsIfAppIsNotInMaintenanceMode(t *testing.T) {
	t.Parallel()

	j1 := &workerUpstreamJob{
		connection: "default",
		queueName:  "queue",
		handler: func(j *workerUpstreamJob) error {
			j.attempts++

			return nil
		},
	}

	j2 := &workerUpstreamJob{
		connection: "default",
		queueName:  "queue",
		handler: func(j *workerUpstreamJob) error {
			j.attempts++

			return nil
		},
	}

	q := &workerFakeQueue{connection: "default", jobs: []queue.Job{j1, j2}}

	// maintenance: false on first call, true forever after.
	var calls int

	w := queue.NewWorker(q, passthroughHandler(), nil, queue.WorkerOptions{
		Sleep:   time.Millisecond,
		MaxTime: 80 * time.Millisecond,
	})
	w.MaintenanceMode = func() bool {
		calls++

		return calls > 1
	}

	if err := w.Run(context.Background(), "queue"); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if j1.attempts != 1 {
		t.Errorf("j1.attempts: got %d, want 1", j1.attempts)
	}

	if j2.attempts != 0 {
		t.Errorf("j2.attempts: got %d, want 0 (maintenance blocked its pop)", j2.attempts)
	}
}

// Port of Framework\Tests\Queue\QueueWorkerTest::testWorkerStopsWhenMemoryExceeded
func TestWorkerStopsWhenMemoryExceeded(t *testing.T) {
	t.Parallel()

	j1 := &workerUpstreamJob{connection: "default", queueName: "queue"}
	j2 := &workerUpstreamJob{connection: "default", queueName: "queue"}

	q := &workerFakeQueue{connection: "default", jobs: []queue.Job{j1, j2}}
	rec := &workerEventRecorder{}

	w := queue.NewWorker(q, passthroughHandler(), rec, queue.WorkerOptions{
		// 1 MiB is always exceeded by the Go test binary's Sys value,
		// so the memory cap trips after the first processed job.
		MemoryLimitMiB: 1,
		Sleep:          time.Millisecond,
	})

	if err := w.Run(context.Background(), "queue"); err != nil {
		t.Fatalf("Run: %v", err)
	}

	if !j1.fired.Load() {
		t.Error("j1: expected fired")
	}

	if j2.fired.Load() {
		t.Error("j2: expected NOT fired (memory cap stopped worker)")
	}

	if w.LastStopReason() != queue.WorkerStopReasonMemoryLimitReached {
		t.Errorf("LastStopReason: got %d, want %d (MemoryLimitReached)", w.LastStopReason(), queue.WorkerStopReasonMemoryLimitReached)
	}

	// Sanity: the status value matches Upstream's 12.
	if int(queue.WorkerStopReasonMemoryLimitReached) != 12 {
		t.Errorf("WorkerStopReasonMemoryLimitReached: got %d, want 12", queue.WorkerStopReasonMemoryLimitReached)
	}
}

// Port of Framework\Tests\Queue\QueueWorkerTest::testJobReleasedEvent
//
// Asserts a JobReleasedAfterException is dispatched when a thrown job
// is released for retry. The retry/fail path is already covered by
// TestJobIsReleasedOnException; this test focuses the assertion on
// the dedicated event.
func TestJobReleasedEvent(t *testing.T) {
	t.Parallel()

	job := &workerUpstreamJob{
		connection: "default",
		queueName:  "queue",
		handler: func(*workerUpstreamJob) error {
			return errors.New("boom")
		},
	}

	_, rec, _ := newWorkerRun(t, job, passthroughHandler(), queue.WorkerOptions{Backoff: time.Second})

	if n := rec.countByType(func(e any) bool { _, ok := e.(queue.JobReleasedAfterException); return ok }); n != 1 {
		t.Errorf("JobReleasedAfterException: got %d, want 1", n)
	}

	if !job.released {
		t.Error("expected job to be released")
	}
}

// Port of Framework\Tests\Queue\QueueWorkerTest::testWorkerStoppingIsDispatched
func TestWorkerStoppingIsDispatched(t *testing.T) {
	t.Parallel()

	q := &workerFakeQueue{connection: "default"}
	rec := &workerEventRecorder{}

	w := queue.NewWorker(q, nil, rec, queue.WorkerOptions{
		Name:        "stop-worker",
		Sleep:       time.Millisecond,
		StopOnEmpty: true,
	})

	if err := w.Run(context.Background(), "queue"); err != nil {
		t.Fatalf("Run: %v", err)
	}

	var stopping *queue.WorkerStopping

	rec.mu.Lock()

	for _, e := range rec.events {
		if s, ok := e.(queue.WorkerStopping); ok {
			sCopy := s
			stopping = &sCopy

			break
		}
	}

	rec.mu.Unlock()

	if stopping == nil {
		t.Fatal("expected a WorkerStopping event")
	}

	if stopping.WorkerName != "stop-worker" {
		t.Errorf("WorkerStopping.WorkerName: got %q, want stop-worker", stopping.WorkerName)
	}
}
