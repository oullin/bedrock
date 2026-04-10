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
	// Attempts returns the number of times the job has been attempted.
	Attempts() int
	// MaxTries returns the maximum number of allowed attempts.
	MaxTries() int
	// Timeout returns the job execution timeout.
	Timeout() time.Duration
	// Backoff returns per-attempt backoff durations.
	Backoff() []time.Duration
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

// Handle implements Handler.
func (f HandlerFunc) Handle(ctx context.Context, job Job) error { return f(ctx, job) }

// JobOptions configures job dispatch options.
type JobOptions struct {
	Queue                   string
	Connection              string
	Delay                   time.Duration
	MaxTries                int
	MaxExceptions           int
	Timeout                 time.Duration
	Backoff                 []time.Duration
	RetryUntil              time.Time
	FailOnTimeout           bool
	UniqueFor               time.Duration
	DeleteWhenMissingModels bool
}
