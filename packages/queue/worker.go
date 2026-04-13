package queue

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// WorkerOptions configures the Worker processing loop.
type WorkerOptions struct {
	// Sleep is how long to sleep when the queue is empty.
	Sleep time.Duration
	// MaxJobs is the maximum number of jobs to process before stopping (0 = unlimited).
	MaxJobs int
	// MaxTime is the maximum wall-clock time before stopping (0 = unlimited).
	MaxTime time.Duration
	// StopOnEmpty stops the worker when the queue is empty.
	StopOnEmpty bool
}

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
}

// NewWorker creates a Worker.
func NewWorker(q Queue, handler Handler, emitter EventEmitter, opts WorkerOptions) *Worker {
	if opts.Sleep == 0 {
		opts.Sleep = time.Second
	}

	return &Worker{queue: q, handler: handler, emitter: emitter, opts: opts}
}

// Run starts the daemon loop, processing jobs until a stop condition is met.
// It handles SIGTERM and SIGQUIT for graceful shutdown.
func (w *Worker) Run(ctx context.Context, queueName string) error {
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	defer signal.Stop(sigs)

	go func() {
		select {
		case <-sigs:
			cancel()
		case <-ctx.Done():
		}
	}()

	start := time.Now()
	processed := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if w.opts.MaxJobs > 0 && processed >= w.opts.MaxJobs {
			return nil
		}

		if w.opts.MaxTime > 0 && time.Since(start) >= w.opts.MaxTime {
			return nil
		}

		job, err := w.queue.Pop(ctx, queueName)

		if err == ErrNoJob {
			if w.opts.StopOnEmpty {
				return nil
			}

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(w.opts.Sleep):
			}

			continue
		}

		if err != nil {
			return err
		}

		w.processJob(ctx, job)
		processed++
	}
}

func (w *Worker) processJob(ctx context.Context, job Job) {
	w.emit(JobProcessing{ConnectionName: w.queue.ConnectionName(), Job: job})
	w.emit(JobAttempted{ConnectionName: w.queue.ConnectionName(), Job: job})

	jobCtx := ctx

	if timeout := job.Timeout(); timeout > 0 {
		var cancel context.CancelFunc

		jobCtx, cancel = context.WithTimeout(ctx, timeout)

		defer cancel()
	}

	if err := w.handler.Handle(jobCtx, job); err != nil {
		w.emit(JobExceptionOccurred{ConnectionName: w.queue.ConnectionName(), Job: job, Err: err})

		if w.shouldFail(job) {
			_ = job.Fail(err)
			w.emit(JobFailed{ConnectionName: w.queue.ConnectionName(), Job: job, Err: err})
		} else {
			backoff := w.backoffFor(job)
			_ = job.Release(backoff)
		}

		return
	}

	_ = job.Delete()
	w.emit(JobProcessed{ConnectionName: w.queue.ConnectionName(), Job: job})
}

// shouldFail reports whether the job has exhausted its retry options.
func (w *Worker) shouldFail(job Job) bool {
	if job.MaxTries() > 0 && job.Attempts() >= job.MaxTries() {
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

func (w *Worker) backoffFor(job Job) time.Duration {
	backoffs := job.Backoff()
	attempt := job.Attempts()

	if attempt > 0 && attempt <= len(backoffs) {
		return backoffs[attempt-1]
	}

	return 0
}

func (w *Worker) emit(event any) {
	if w.emitter != nil {
		w.emitter.Emit(event)
	}
}
