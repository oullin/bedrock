package queue_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/bedrock/packages/queue"
	"github.com/bedrock/packages/queue/drivers"
)

// mockQueue wraps a real queue to control behavior in tests.
type mockQueue struct {
	queue.Queue
	jobs       []queue.Job
	mu         sync.Mutex
	popCount   int
	connection string
}

// mockJob implements queue.Job for worker tests.
type mockJob struct {
	payload       []byte
	queue         string
	connection    string
	attempts      int
	maxTries      int
	maxExceptions int
	timeout       time.Duration
	backoff       []time.Duration
	retryUntil    *time.Time
	released      bool
	deleted       bool
	failed        bool
	releaseDelay  time.Duration
}

// eventRecorder records emitted events.
type eventRecorder struct {
	mu     sync.Mutex
	events []any
}

// Use a real redis-backed queue with infinite jobs.

// With StopOnEmpty false, Sleep 0 (should default to 1s), we test by using context timeout.

// Just verify it doesn't panic and returns.

// No max tries limit.

// A queue that returns a non-ErrNoJob error on Pop.

// errorQueue always returns an error from Pop.
type errorQueue struct {
	err error
}

// newWorkerMockRedisClient creates a simple in-memory Redis client for worker tests.

type workerRedisClient struct {
	mu    sync.Mutex
	lists map[string][]string
}

func (q *mockQueue) Push(_ context.Context, queueName string, payload []byte) (string, error) {
	q.mu.Lock()

	defer q.mu.Unlock()

	q.jobs = append(q.jobs, &mockJob{
		payload:    payload,
		queue:      queueName,
		connection: q.connection,
	})

	return "", nil
}

func (q *mockQueue) Pop(_ context.Context, _ string) (queue.Job, error) {
	q.mu.Lock()

	defer q.mu.Unlock()

	if len(q.jobs) == 0 {
		return nil, queue.ErrNoJob
	}

	q.popCount++

	job := q.jobs[0]
	q.jobs = q.jobs[1:]

	return job, nil
}

func (q *mockQueue) ConnectionName() string { return q.connection }

func (j *mockJob) UUID() string                 { return "test-uuid" }
func (j *mockJob) GetJobID() string             { return "test-id" }
func (j *mockJob) Payload() []byte              { return j.payload }
func (j *mockJob) Fire(_ context.Context) error { return nil }
func (j *mockJob) Attempts() int                { return j.attempts }
func (j *mockJob) MaxTries() int                { return j.maxTries }
func (j *mockJob) MaxExceptions() int           { return j.maxExceptions }
func (j *mockJob) Timeout() time.Duration       { return j.timeout }
func (j *mockJob) Backoff() []time.Duration     { return j.backoff }
func (j *mockJob) RetryUntil() *time.Time       { return j.retryUntil }
func (j *mockJob) IsDeleted() bool              { return j.deleted }
func (j *mockJob) IsReleased() bool             { return j.released }
func (j *mockJob) HasFailed() bool              { return j.failed }
func (j *mockJob) GetQueue() string             { return j.queue }
func (j *mockJob) GetConnectionName() string    { return j.connection }

func (j *mockJob) Release(delay time.Duration) error {
	j.released = true
	j.releaseDelay = delay

	return nil
}

func (j *mockJob) Delete() error {
	j.deleted = true

	return nil
}

func (j *mockJob) Fail(_ error) error {
	j.failed = true

	return nil
}

func (j *mockJob) MarkAsFailed(err error) error {
	return j.Fail(err)
}

func (r *eventRecorder) Emit(event any) {
	r.mu.Lock()

	defer r.mu.Unlock()

	r.events = append(r.events, event)
}

func (r *eventRecorder) count(filter func(any) bool) int {
	r.mu.Lock()

	defer r.mu.Unlock()

	n := 0

	for _, e := range r.events {
		if filter(e) {
			n++
		}
	}

	return n
}

func TestWorkerProcessesJobAndEmitsEvents(t *testing.T) {
	t.Parallel()

	q := &mockQueue{connection: "test"}
	q.jobs = []queue.Job{&mockJob{payload: []byte("p"), queue: "q", connection: "test"}}

	recorder := &eventRecorder{}
	handled := 0

	handler := queue.HandlerFunc(func(_ context.Context, _ queue.Job) error {
		handled++

		return nil
	})

	w := queue.NewWorker(q, handler, recorder, queue.WorkerOptions{StopOnEmpty: true})
	_ = w.Run(context.Background(), "q")

	if handled != 1 {
		t.Errorf("expected 1 handled, got %d", handled)
	}

	processing := recorder.count(func(e any) bool { _, ok := e.(queue.JobProcessing); return ok })

	if processing != 1 {
		t.Errorf("expected 1 JobProcessing event, got %d", processing)
	}

	processed := recorder.count(func(e any) bool { _, ok := e.(queue.JobProcessed); return ok })

	if processed != 1 {
		t.Errorf("expected 1 JobProcessed event, got %d", processed)
	}

	attempted := recorder.count(func(e any) bool { _, ok := e.(queue.JobAttempted); return ok })

	if attempted != 1 {
		t.Errorf("expected 1 JobAttempted event, got %d", attempted)
	}
}

func TestWorkerFailedJobEmitsFailedEvent(t *testing.T) {
	t.Parallel()

	job := &mockJob{payload: []byte("p"), queue: "q", connection: "test", maxTries: 1, attempts: 1}
	q := &mockQueue{connection: "test", jobs: []queue.Job{job}}

	recorder := &eventRecorder{}

	handler := queue.HandlerFunc(func(_ context.Context, _ queue.Job) error {
		return errors.New("fail")
	})

	w := queue.NewWorker(q, handler, recorder, queue.WorkerOptions{StopOnEmpty: true})
	_ = w.Run(context.Background(), "q")

	failed := recorder.count(func(e any) bool { _, ok := e.(queue.JobFailed); return ok })

	if failed != 1 {
		t.Errorf("expected 1 JobFailed event, got %d", failed)
	}

	exception := recorder.count(func(e any) bool { _, ok := e.(queue.JobExceptionOccurred); return ok })

	if exception != 1 {
		t.Errorf("expected 1 JobExceptionOccurred event, got %d", exception)
	}

	if !job.failed {
		t.Error("expected job to be marked as failed")
	}
}

func TestWorkerReleasesJobOnError(t *testing.T) {
	t.Parallel()

	job := &mockJob{payload: []byte("p"), queue: "q", connection: "test", maxTries: 3, attempts: 1}
	q := &mockQueue{connection: "test", jobs: []queue.Job{job}}

	handler := queue.HandlerFunc(func(_ context.Context, _ queue.Job) error {
		return errors.New("retry")
	})

	w := queue.NewWorker(q, handler, nil, queue.WorkerOptions{StopOnEmpty: true})
	_ = w.Run(context.Background(), "q")

	if !job.released {
		t.Error("expected job to be released for retry")
	}

	if job.failed {
		t.Error("expected job NOT to be failed (retries remaining)")
	}
}

func TestWorkerBackoffStrategy(t *testing.T) {
	t.Parallel()

	job := &mockJob{
		payload:    []byte("p"),
		queue:      "q",
		connection: "test",
		maxTries:   5,
		attempts:   2,
		backoff:    []time.Duration{time.Second, 5 * time.Second, 10 * time.Second},
	}
	q := &mockQueue{connection: "test", jobs: []queue.Job{job}}

	handler := queue.HandlerFunc(func(_ context.Context, _ queue.Job) error {
		return errors.New("retry")
	})

	w := queue.NewWorker(q, handler, nil, queue.WorkerOptions{StopOnEmpty: true})
	_ = w.Run(context.Background(), "q")

	if job.releaseDelay != 5*time.Second {
		t.Errorf("expected 5s backoff for attempt 2, got %v", job.releaseDelay)
	}
}

func TestWorkerMaxJobsLimit(t *testing.T) {
	t.Parallel()

	q := &mockQueue{connection: "test"}

	for i := 0; i < 10; i++ {
		q.jobs = append(q.jobs, &mockJob{payload: []byte("p"), queue: "q", connection: "test"})
	}

	handled := 0

	handler := queue.HandlerFunc(func(_ context.Context, _ queue.Job) error {
		handled++

		return nil
	})

	w := queue.NewWorker(q, handler, nil, queue.WorkerOptions{MaxJobs: 3})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)

	defer cancel()

	_ = w.Run(ctx, "q")

	if handled != 3 {
		t.Errorf("expected 3 handled (MaxJobs), got %d", handled)
	}
}

func TestWorkerMaxTimeLimit(t *testing.T) {
	t.Parallel()

	client := newWorkerMockRedisClient()
	inner := drivers.NewRedisDriver(client, "redis")

	for i := 0; i < 100; i++ {
		_, _ = inner.Push(context.Background(), "q", []byte("p"))
	}

	handler := queue.HandlerFunc(func(_ context.Context, _ queue.Job) error {
		time.Sleep(10 * time.Millisecond)

		return nil
	})

	w := queue.NewWorker(inner, handler, nil, queue.WorkerOptions{MaxTime: 50 * time.Millisecond})

	start := time.Now()
	_ = w.Run(context.Background(), "q")
	elapsed := time.Since(start)

	if elapsed > 200*time.Millisecond {
		t.Errorf("expected to stop within ~50ms, took %v", elapsed)
	}
}

func TestWorkerDefaultSleep(t *testing.T) {
	t.Parallel()

	drv := drivers.NewNullDriver("null")

	handler := queue.HandlerFunc(func(_ context.Context, _ queue.Job) error { return nil })

	w := queue.NewWorker(drv, handler, nil, queue.WorkerOptions{})

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)

	defer cancel()

	_ = w.Run(ctx, "q")

}

func TestWorkerContextCancellation(t *testing.T) {
	t.Parallel()

	drv := drivers.NewNullDriver("null")

	handler := queue.HandlerFunc(func(_ context.Context, _ queue.Job) error { return nil })

	w := queue.NewWorker(drv, handler, nil, queue.WorkerOptions{Sleep: time.Millisecond})

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	err := w.Run(ctx, "q")

	if err == nil || !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
}

func TestWorkerDeletesSuccessfulJob(t *testing.T) {
	t.Parallel()

	job := &mockJob{payload: []byte("p"), queue: "q", connection: "test"}
	q := &mockQueue{connection: "test", jobs: []queue.Job{job}}

	handler := queue.HandlerFunc(func(_ context.Context, _ queue.Job) error { return nil })

	w := queue.NewWorker(q, handler, nil, queue.WorkerOptions{StopOnEmpty: true})
	_ = w.Run(context.Background(), "q")

	if !job.deleted {
		t.Error("expected job to be deleted after success")
	}
}

func TestWorkerNilEmitter(t *testing.T) {
	t.Parallel()

	job := &mockJob{payload: []byte("p"), queue: "q", connection: "test"}
	q := &mockQueue{connection: "test", jobs: []queue.Job{job}}

	handler := queue.HandlerFunc(func(_ context.Context, _ queue.Job) error { return nil })

	w := queue.NewWorker(q, handler, nil, queue.WorkerOptions{StopOnEmpty: true})

	err := w.Run(context.Background(), "q")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestWorkerRetryUntilExpired(t *testing.T) {
	t.Parallel()

	past := time.Now().Add(-time.Hour)
	job := &mockJob{
		payload:    []byte("p"),
		queue:      "q",
		connection: "test",
		maxTries:   0,
		attempts:   1,
		retryUntil: &past,
	}
	q := &mockQueue{connection: "test", jobs: []queue.Job{job}}

	handler := queue.HandlerFunc(func(_ context.Context, _ queue.Job) error {
		return errors.New("fail")
	})

	w := queue.NewWorker(q, handler, nil, queue.WorkerOptions{StopOnEmpty: true})
	_ = w.Run(context.Background(), "q")

	if !job.failed {
		t.Error("expected job to be failed when RetryUntil is in the past")
	}
}

func TestWorkerMaxExceptionsExceeded(t *testing.T) {
	t.Parallel()

	job := &mockJob{
		payload:       []byte("p"),
		queue:         "q",
		connection:    "test",
		maxTries:      0,
		maxExceptions: 2,
		attempts:      3,
	}
	q := &mockQueue{connection: "test", jobs: []queue.Job{job}}

	handler := queue.HandlerFunc(func(_ context.Context, _ queue.Job) error {
		return errors.New("fail")
	})

	w := queue.NewWorker(q, handler, nil, queue.WorkerOptions{StopOnEmpty: true})
	_ = w.Run(context.Background(), "q")

	if !job.failed {
		t.Error("expected job to be failed when MaxExceptions exceeded")
	}
}

func TestWorkerPopError(t *testing.T) {
	t.Parallel()

	q := &errorQueue{err: errors.New("connection lost")}

	handler := queue.HandlerFunc(func(_ context.Context, _ queue.Job) error { return nil })

	w := queue.NewWorker(q, handler, nil, queue.WorkerOptions{})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)

	defer cancel()

	err := w.Run(ctx, "q")

	if err == nil || err.Error() != "connection lost" {
		t.Errorf("expected 'connection lost', got %v", err)
	}
}

func (q *errorQueue) Push(_ context.Context, _ string, _ []byte) (string, error) { return "", nil }
func (q *errorQueue) PushDelayed(_ context.Context, _ string, _ []byte, _ time.Duration) (string, error) {
	return "", nil
}
func (q *errorQueue) PushMultiple(_ context.Context, _ string, _ [][]byte) ([]string, error) {
	return nil, nil
}
func (q *errorQueue) Pop(_ context.Context, _ string) (queue.Job, error)      { return nil, q.err }
func (q *errorQueue) Size(_ context.Context, _ string) (int64, error)         { return 0, nil }
func (q *errorQueue) PendingSize(_ context.Context, _ string) (int64, error)  { return 0, nil }
func (q *errorQueue) DelayedSize(_ context.Context, _ string) (int64, error)  { return 0, nil }
func (q *errorQueue) ReservedSize(_ context.Context, _ string) (int64, error) { return 0, nil }
func (q *errorQueue) ConnectionName() string                                  { return "error" }

func newWorkerMockRedisClient() *workerRedisClient {
	return &workerRedisClient{
		lists: make(map[string][]string),
	}
}

func (c *workerRedisClient) LPush(_ context.Context, key string, values ...any) error {
	c.mu.Lock()

	defer c.mu.Unlock()

	for _, v := range values {
		c.lists[key] = append([]string{v.(string)}, c.lists[key]...)
	}

	return nil
}

func (c *workerRedisClient) RPop(_ context.Context, key string) (string, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	list := c.lists[key]

	if len(list) == 0 {
		return "", errors.New("empty")
	}

	val := list[len(list)-1]
	c.lists[key] = list[:len(list)-1]

	return val, nil
}

func (c *workerRedisClient) ZAdd(_ context.Context, _ string, _ float64, _ string) error {
	return nil
}

func (c *workerRedisClient) ZRangeByScore(_ context.Context, _ string, _, _ float64) ([]string, error) {
	return nil, nil
}

func (c *workerRedisClient) ZRem(_ context.Context, _ string, _ ...any) error { return nil }

func (c *workerRedisClient) LLen(_ context.Context, key string) (int64, error) {
	c.mu.Lock()

	defer c.mu.Unlock()

	return int64(len(c.lists[key])), nil
}

func (c *workerRedisClient) ZCard(_ context.Context, _ string) (int64, error) { return 0, nil }
