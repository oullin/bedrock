package queue

import (
	"context"
	"runtime"
	"strings"
	"sync/atomic"
	"time"
)

// WorkerOptions configures the Worker processing loop.
type WorkerOptions struct {
	// Name identifies this worker in the WorkerStarting/WorkerStopping
	// events.
	Name string
	// Sleep is how long to sleep when the queue is empty.
	Sleep time.Duration
	// MaxJobs is the maximum number of jobs to process before stopping (0 = unlimited).
	MaxJobs int
	// MaxTime is the maximum wall-clock time before stopping (0 = unlimited).
	MaxTime time.Duration
	// StopOnEmpty stops the worker when the queue is empty.
	StopOnEmpty bool
	// MemoryLimitMiB is the RSS cap the worker self-enforces by polling
	// runtime.MemStats.Sys between iterations. A value of 0 disables
	// the cap.
	MemoryLimitMiB int64
	// Backoff is the default per-attempt delay used when a released job
	// has no Backoff() of its own.
	Backoff time.Duration
	// MaxTries is the default attempt cap used when a job's own
	// MaxTries() returns 0.
	MaxTries int
}

// ExceptionReporter is the optional contract the Worker uses to report
// exceptions that escape a job handler. It is the Go analogue of
// Ref: @bedrock/code-0193
type ExceptionReporter interface {
	ReportException(err error)
}

// WorkerPopCallback overrides how a worker pops from one queue name.
type WorkerPopCallback func(ctx context.Context, q Queue, queueName string) (Job, error)

// WorkerStopReason is the machine-readable status code returned from
// the daemon loop when it exits. The numeric values intentionally
// mirror the upstream WorkerStopReason enum so operators can treat the
// process exit codes identically.
type WorkerStopReason int

// WorkerStopReasonNone means the worker has not stopped yet or
// exited cleanly with no dedicated reason.

// WorkerStopReasonStopOnEmpty — the queue drained while
// StopOnEmpty was set, so the worker exited voluntarily.

// WorkerStopReasonMemoryLimitReached — MemoryExceeded tripped
// between iterations. Matches the upstream STATUS_MEMORY_LIMIT = 12.

// WorkerStopReasonLostConnection — pop repeatedly errored on a
// single connection. Matches the upstream STATUS_LOST_CONNECTION = 13.

// WorkerStopReasonMaxTimeExceeded — opts.MaxTime elapsed. Matches
// the upstream STATUS_OUT_OF_TIME = 14.

// WorkerStopReasonMaxJobsExceeded — opts.MaxJobs reached.

// EventEmitter receives worker lifecycle events.
type EventEmitter interface {
	Emit(event any)
}

// Worker processes jobs from a Queue in a daemon loop with graceful shutdown.
type Worker struct {
	queue   Queue
	handler Handler
	emitter EventEmitter
	opts    WorkerOptions
	// SleepFunc is called when the queue is empty. It defaults to a
	// context-aware time.After sleep. Tests override it to capture the
	// requested duration without blocking. Setting SleepFunc to nil
	// restores the default behaviour.
	SleepFunc func(ctx context.Context, d time.Duration)
	// sleptFor records the duration of the most recent sleep call. The
	// SleptFor accessor exposes it to tests — this is the Go analogue
	// of the upstream public $sleptFor property on the worker.
	sleptFor time.Duration
	// ExceptionReporter, if non-nil, receives a ReportException call
	// for every exception raised by a job handler (or by the pre-fire
	// max-attempts/retry-until check). May be nil.
	ExceptionReporter ExceptionReporter
	// ReportJobExceptions gates the reporter. Default is true; set
	// false to suppress reporter calls without suppressing the
	// JobExceptionOccurred event emission.
	// Worker::$reportJobExceptions flag.
	ReportJobExceptions bool
	// MaintenanceMode, if non-nil, is called by the Run loop before
	// each iteration. When it returns true, the worker sleeps for
	// opts. Sleep and skips the pop.
	// into $this->manager->isDownForMaintenance().
	MaintenanceMode func() bool
	// lastStopReason captures the reason the most recent Run call
	// exited. Read via LastStopReason().
	lastStopReason WorkerStopReason
	// popCallbacks holds per-queue pop selectors. The "*" entry applies
	// to queue names that do not have a more specific callback.
	popCallbacks map[string]WorkerPopCallback
	// LostConnectionDetector, when non-nil, decides whether a pop error
	// should stop the daemon with WorkerStopReasonLostConnection.
	LostConnectionDetector func(error) bool
	// paused is the atomic flag flipped by SIGUSR2 (pause) and SIGCONT
	// (resume) so the run loop can stall between iterations without
	// the platform-specific signal goroutine reaching back into the
	// loop's local state. Treated as a strict 0/1; access only via the
	// IsPaused/SetPaused helpers.
	paused atomic.Bool
}

// NewWorker creates a Worker.

// SleptFor returns the duration of the most recent sleep the worker
// performed. Tests use this to assert sleep behaviour without taking
// a timing dependency; production callers will usually ignore it.

// LastStopReason returns the reason the most recent Run call exited,
// or WorkerStopReasonNone if Run has not yet been called or exited
// via a context cancellation.

// parseQueueNames splits a comma-separated queue string into its
// individual names. Empty segments are dropped. A single-queue string
// returns a single-element slice.

// splitComma is a tiny helper that avoids pulling strings.Split into
// this file solely for a one-shot split. Keeps the worker self-contained.

// popFromQueues tries each queue name in order and returns the first
// job it finds. A real error (other than ErrNoJob) aborts the loop and
// is returned to the caller. If every queue is empty it returns
// (nil, ErrNoJob).
//
// This is the Go analogue of the upstream "high,low" priority-queue
// handling: the worker pops from the first bucket, falling through to
// the next only when the higher-priority bucket is empty.

// reportPopError dispatches an exception reporter call for a pop-time
// error when reporting is enabled.

// MemoryExceeded reports whether the worker's resident memory footprint
// in MiB has met or exceeded memoryLimitMiB. A zero or negative limit is
// treated as "no limit" and always returns false.
//
// The Go analogue of the upstream memory_get_usage(true) is
// runtime.MemStats.Sys — total bytes obtained from the OS, which
// includes the unused portion the allocator is holding. Using Sys
// keeps the observable behaviour close to PHP's "real" memory count.

// sleep blocks for d honouring context cancellation. Overridable via
// SleepFunc for tests. It also records d in sleptFor.

// RunNextJob performs a single iteration of the worker loop: it emits
// JobPopping, pops one job from the (comma-separated) queue names in
// priority order, and either processes it (emitting JobPopped + the
// rest of the lifecycle) or sleeps when every queue is empty.
//
// queueNames accepts either a single queue ("default") or a priority
// chain ("high,low") — higher-priority queues are tried first. On a
// non-ErrNoJob pop error, the error is dispatched through
// ExceptionReporter (if configured) and the worker sleeps; it never
// propagates back to the caller.

// Run starts the daemon loop, processing jobs until a stop condition is met.
// It handles SIGTERM and SIGQUIT for graceful shutdown.
//
// Run dispatches WorkerStarting before entering the loop and
// WorkerStopping on exit (both with the configured WorkerOptions.Name).
// Per-iteration events (JobPopping/JobPopped/JobProcessing/...) are
// emitted from the main loop.

// Between iterations, check the memory cap. the upstream daemon
// does this post-process so the current job is allowed to
// finish before the worker recycles itself.

// processJob runs the full lifecycle of a popped job:
//
//   - If the job was already deleted before the worker saw it, emit
// JobProcessed and return (skip Fire).
//   - Pre-fire exhaustion check: if attempts already exceeds the
//     effective max-tries OR retry-until has expired, fail with a
//     MaxAttemptsExceededError without calling the handler.
//   - Fire the handler (with a timeout context if set).
//   - On error: dispatch JobExceptionOccurred + report, then either
//     fail (attempts exhausted post-fire) or release with the
//     effective backoff. Release emits JobReleasedAfterException.
//   - On success: emit JobProcessed. No auto-delete — handlers that
//     need the job removed from the backend must call job.Delete
//     themselves (matching the upstream CallQueuedHandler contract).

// handleJobException runs the shared failure pipeline used by both
// the pre-fire exhaustion check and the post-handler error path.
//
// Sequence: JobExceptionOccurred → reporter (if enabled) →
// {JobFailed via job.Fail(err)} OR {JobReleasedAfterException via
// job.Release(backoff)}.

// markIfExhausted returns a MaxAttemptsExceededError when the job's
// pre-fire state indicates it has no retry budget left (strict >
// because attempts have not yet been incremented for the current
// run), and nil otherwise.

// shouldFail reports whether the job has exhausted its retry options
// and should therefore go through the fail path rather than release.
// Uses inclusive >= for the post-fire decision (attempts has just
// been incremented by the handler's run).

// effectiveMaxTries returns the max-tries budget to use for job:
// the job's own MaxTries wins over the WorkerOptions fallback.

// effectiveBackoff returns the release delay for the next retry of
// job. The job's own Backoff slice wins over WorkerOptions.Backoff.
// When the slice is shorter than attempts, the last element is used
// When the slice is shorter than attempts, the last element is reused.

// jobNameShim adapts a queue.Job to the ResolveNamer interface used
// by MaxAttemptsExceededError. The Job interface itself stays frozen
// (Step 1), so the shim provides the missing method via a fallback
// to GetQueue. A future Job contract change (Step 9) can drop the
// shim and have Job implement ResolveNamer directly.
type jobNameShim struct{ job Job }

const (
	WorkerStopReasonNone WorkerStopReason = 0

	WorkerStopReasonStopOnEmpty WorkerStopReason = 11

	WorkerStopReasonMemoryLimitReached WorkerStopReason = 12

	WorkerStopReasonLostConnection WorkerStopReason = 13

	WorkerStopReasonMaxTimeExceeded WorkerStopReason = 14

	WorkerStopReasonMaxJobsExceeded WorkerStopReason = 15
)

func NewWorker(q Queue, handler Handler, emitter EventEmitter, opts WorkerOptions) *Worker {
	if opts.Sleep == 0 {
		opts.Sleep = time.Second
	}

	return &Worker{
		queue:               q,
		handler:             handler,
		emitter:             emitter,
		opts:                opts,
		ReportJobExceptions: true,
	}
}

func (w *Worker) SleptFor() time.Duration { return w.sleptFor }

func (w *Worker) LastStopReason() WorkerStopReason { return w.lastStopReason }

func (w *Worker) PopUsing(queueName string, callback WorkerPopCallback) {
	if w.popCallbacks == nil {
		w.popCallbacks = make(map[string]WorkerPopCallback)
	}

	if callback == nil {
		delete(w.popCallbacks, queueName)

		return
	}

	w.popCallbacks[queueName] = callback
}

func parseQueueNames(raw string) []string {
	parts := splitComma(raw)
	out := parts[:0]

	for _, p := range parts {
		if p != "" {
			out = append(out, p)
		}
	}

	if len(out) == 0 {
		return []string{"default"}
	}

	return out
}

func splitComma(raw string) []string {
	out := make([]string, 0, 2)
	start := 0

	for i := 0; i < len(raw); i++ {
		if raw[i] == ',' {
			out = append(out, raw[start:i])
			start = i + 1
		}
	}

	out = append(out, raw[start:])

	return out
}

func (w *Worker) popFromQueues(ctx context.Context, queueNames []string) (Job, error) {
	for _, q := range queueNames {
		pop := w.popCallbackFor(q)
		job, err := pop(ctx, w.queue, q)

		if err == nil && job != nil {
			return job, nil
		}

		if err != nil && err != ErrNoJob {
			return nil, err
		}
	}

	return nil, ErrNoJob
}

func (w *Worker) popCallbackFor(queueName string) WorkerPopCallback {
	if w.popCallbacks != nil {
		if callback := w.popCallbacks[queueName]; callback != nil {
			return callback
		}

		if callback := w.popCallbacks["*"]; callback != nil {
			return callback
		}
	}

	return func(ctx context.Context, q Queue, name string) (Job, error) {
		return q.Pop(ctx, name)
	}
}

func (w *Worker) isLostConnection(err error) bool {
	if err == nil {
		return false
	}

	if w.LostConnectionDetector != nil {
		return w.LostConnectionDetector(err)
	}

	return strings.Contains(strings.ToLower(err.Error()), "lost connection")
}

func (w *Worker) reportPopError(err error) {
	if w.ReportJobExceptions && w.ExceptionReporter != nil {
		w.ExceptionReporter.ReportException(err)
	}
}

func (w *Worker) MemoryExceeded(memoryLimitMiB int64) bool {
	if memoryLimitMiB <= 0 {
		return false
	}

	var ms runtime.MemStats

	runtime.ReadMemStats(&ms)

	return ms.Sys >= uint64(memoryLimitMiB)*1024*1024
}

func (w *Worker) sleep(ctx context.Context, d time.Duration) {
	w.sleptFor = d

	if w.SleepFunc != nil {
		w.SleepFunc(ctx, d)

		return
	}

	select {
	case <-ctx.Done():
	case <-time.After(d):
	}
}

func (w *Worker) RunNextJob(ctx context.Context, queueNames string) error {
	queues := parseQueueNames(queueNames)

	w.emit(JobPopping{ConnectionName: w.queue.ConnectionName()})

	job, err := w.popFromQueues(ctx, queues)

	if err != nil && err != ErrNoJob {
		w.reportPopError(err)
		w.sleep(ctx, w.opts.Sleep)

		return nil
	}

	if job == nil {
		w.sleep(ctx, w.opts.Sleep)

		return nil
	}

	w.emit(JobPopped{ConnectionName: w.queue.ConnectionName(), Job: job})
	w.processJob(ctx, job)

	return nil
}

func (w *Worker) Run(ctx context.Context, queueName string) error {
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	stopSignals := w.installSignalHandlers(ctx, cancel, queueName)

	defer stopSignals()

	w.emit(WorkerStarting{
		ConnectionName: w.queue.ConnectionName(),
		Queue:          queueName,
		WorkerName:     w.opts.Name,
	})

	defer func() {
		w.emit(WorkerStopping{WorkerName: w.opts.Name, Status: int(w.lastStopReason)})
	}()

	queues := parseQueueNames(queueName)
	start := time.Now()
	processed := 0
	w.lastStopReason = WorkerStopReasonNone

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if w.paused.Load() {
			w.sleep(ctx, w.opts.Sleep)

			if ctx.Err() != nil {
				return ctx.Err()
			}

			continue
		}

		if w.opts.MaxJobs > 0 && processed >= w.opts.MaxJobs {
			w.lastStopReason = WorkerStopReasonMaxJobsExceeded

			return nil
		}

		if w.opts.MaxTime > 0 && time.Since(start) >= w.opts.MaxTime {
			w.lastStopReason = WorkerStopReasonMaxTimeExceeded

			return nil
		}

		if w.MaintenanceMode != nil && w.MaintenanceMode() {
			w.sleep(ctx, w.opts.Sleep)

			if ctx.Err() != nil {
				return ctx.Err()
			}

			continue
		}

		w.emit(JobPopping{ConnectionName: w.queue.ConnectionName()})

		job, err := w.popFromQueues(ctx, queues)

		if err != nil && err != ErrNoJob {
			w.reportPopError(err)

			if w.isLostConnection(err) {
				w.lastStopReason = WorkerStopReasonLostConnection

				return nil
			}

			w.sleep(ctx, w.opts.Sleep)

			if ctx.Err() != nil {
				return ctx.Err()
			}

			continue
		}

		if err == ErrNoJob || job == nil {
			if w.opts.StopOnEmpty {
				w.lastStopReason = WorkerStopReasonStopOnEmpty

				return nil
			}

			w.sleep(ctx, w.opts.Sleep)

			if ctx.Err() != nil {
				return ctx.Err()
			}

			continue
		}

		w.emit(JobPopped{ConnectionName: w.queue.ConnectionName(), Job: job})
		w.processJob(ctx, job)
		processed++

		if w.opts.MemoryLimitMiB > 0 && w.MemoryExceeded(w.opts.MemoryLimitMiB) {
			w.lastStopReason = WorkerStopReasonMemoryLimitReached

			return nil
		}
	}
}

func (w *Worker) processJob(ctx context.Context, job Job) {
	if job.IsDeleted() {
		w.emit(JobProcessed{ConnectionName: w.queue.ConnectionName(), Job: job})

		return
	}

	if err := w.markIfExhausted(job); err != nil {
		w.handleJobException(job, err)

		return
	}

	w.emit(JobProcessing{ConnectionName: w.queue.ConnectionName(), Job: job})
	w.emit(JobAttempted{ConnectionName: w.queue.ConnectionName(), Job: job})

	jobCtx := ctx

	if timeout := job.Timeout(); timeout > 0 {
		var cancel context.CancelFunc

		jobCtx, cancel = context.WithTimeout(ctx, timeout)

		defer cancel()
	}

	if err := w.handler.Handle(jobCtx, job); err != nil {
		w.handleJobException(job, err)

		return
	}

	w.emit(JobProcessed{ConnectionName: w.queue.ConnectionName(), Job: job})
}

func (w *Worker) handleJobException(job Job, err error) {
	w.emit(JobExceptionOccurred{ConnectionName: w.queue.ConnectionName(), Job: job, Err: err})

	if w.ReportJobExceptions && w.ExceptionReporter != nil {
		w.ExceptionReporter.ReportException(err)
	}

	if w.shouldFail(job) {
		_ = job.Fail(err)
		w.emit(JobFailed{ConnectionName: w.queue.ConnectionName(), Job: job, Err: err})

		return
	}

	_ = job.Release(w.effectiveBackoff(job))
	w.emit(JobReleasedAfterException{ConnectionName: w.queue.ConnectionName(), Job: job})
}

func (w *Worker) markIfExhausted(job Job) error {
	if max := w.effectiveMaxTries(job); max > 0 && job.Attempts() > max {
		return NewMaxAttemptsExceededErrorForJob(jobNameShim{job: job})
	}

	if retryUntil := job.RetryUntil(); retryUntil != nil && time.Now().After(*retryUntil) {
		return NewMaxAttemptsExceededErrorForJob(jobNameShim{job: job})
	}

	return nil
}

func (w *Worker) shouldFail(job Job) bool {
	if max := w.effectiveMaxTries(job); max > 0 && job.Attempts() >= max {
		return true
	}

	if retryUntil := job.RetryUntil(); retryUntil != nil && time.Now().After(*retryUntil) {
		return true
	}

	if job.MaxExceptions() > 0 && job.Attempts() >= job.MaxExceptions() {
		return true
	}

	return false
}

func (w *Worker) effectiveMaxTries(job Job) int {
	if n := job.MaxTries(); n > 0 {
		return n
	}

	return w.opts.MaxTries
}

func (w *Worker) effectiveBackoff(job Job) time.Duration {
	if b := job.Backoff(); len(b) > 0 {
		idx := job.Attempts() - 1

		if idx < 0 {
			idx = 0
		}

		if idx >= len(b) {
			idx = len(b) - 1
		}

		return b[idx]
	}

	return w.opts.Backoff
}

func (s jobNameShim) ResolveName() string {
	if nr, ok := s.job.(ResolveNamer); ok {
		return nr.ResolveName()
	}

	if q := s.job.GetQueue(); q != "" {
		return q
	}

	return "job"
}

func (w *Worker) emit(event any) {
	if w.emitter != nil {
		w.emitter.Emit(event)
	}
}
