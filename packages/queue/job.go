package queue

import (
	"context"
	"time"
)

// Job represents a queued job instance.
type Job interface {
	// UUID returns the job's unique identifier.
	UUID() string
	// GetJobID returns the backend-specific job identifier.
	GetJobID() string
	// Payload returns the raw job payload.
	Payload() []byte
	// Fire executes the job.
	Fire(ctx context.Context) error
	// Release puts the job back on the queue after a delay.
	Release(delay time.Duration) error
	// Delete removes the job from the queue.
	Delete() error
	// Fail marks the job as failed.
	Fail(err error) error
	// MarkAsFailed is an alias for Fail.
	MarkAsFailed(err error) error
	// Attempts returns the number of times the job has been attempted.
	Attempts() int
	// MaxTries returns the maximum number of allowed attempts.
	MaxTries() int
	// MaxExceptions returns the maximum number of exceptions before failing.
	MaxExceptions() int
	// Timeout returns the job execution timeout.
	Timeout() time.Duration
	// Backoff returns per-attempt backoff durations.
	Backoff() []time.Duration
	// RetryUntil returns the time after which to stop retrying.
	RetryUntil() *time.Time
	// IsDeleted reports whether the job has been deleted.
	IsDeleted() bool
	// IsReleased reports whether the job has been released.
	IsReleased() bool
	// HasFailed reports whether the job has been marked as failed.
	HasFailed() bool
	// GetQueue returns the queue name.
	GetQueue() string
	// GetConnectionName returns the connection name.
	GetConnectionName() string
}

// Handler handles a job type.
type Handler interface {
	Handle(ctx context.Context, job Job) error
}

// HandlerFunc is a function that implements Handler.
type HandlerFunc func(ctx context.Context, job Job) error

// FailureHandler is the optional contract a Handler can implement to
// receive a callback when a job it was processing has been marked as
// failed. The driver or worker invokes Failed after Job.Fail has been
// called and before the JobFailed event is emitted.
//
// callback invoked by CallQueuedHandler::failed.
type FailureHandler interface {
	Failed(ctx context.Context, job Job, err error)
}

// Handle implements Handler.

// JobOptions configures job dispatch options.
type JobOptions struct {
	Queue                   string
	Connection              string
	Delay                   time.Duration
	BatchID                 string
	MaxTries                int
	MaxExceptions           int
	Timeout                 time.Duration
	Backoff                 []time.Duration
	RetryUntil              time.Time
	FailOnTimeout           bool
	UniqueFor               time.Duration
	DeleteWhenMissingModels bool
}

func (f HandlerFunc) Handle(ctx context.Context, job Job) error { return f(ctx, job) }

// WithoutDelay returns a copy of opts with its dispatch delay cleared.
func (o JobOptions) WithoutDelay() JobOptions {
	o.Delay = 0

	return o
}
