package drivers

import (
	"context"
	"time"
)

// BaseJob provides default implementations of Job methods for embedding.
type BaseJob struct {
	id            string
	uuid          string
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
	releaseFunc   func(delay time.Duration) error
	deleteFunc    func() error
	failFunc      func(err error) error
	fireFunc      func(ctx context.Context) error
}

func (j *BaseJob) UUID() string              { return j.uuid }
func (j *BaseJob) GetJobID() string          { return j.id }
func (j *BaseJob) Payload() []byte           { return j.payload }
func (j *BaseJob) Attempts() int             { return j.attempts }
func (j *BaseJob) MaxTries() int             { return j.maxTries }
func (j *BaseJob) MaxExceptions() int        { return j.maxExceptions }
func (j *BaseJob) Timeout() time.Duration    { return j.timeout }
func (j *BaseJob) Backoff() []time.Duration  { return j.backoff }
func (j *BaseJob) RetryUntil() *time.Time    { return j.retryUntil }
func (j *BaseJob) IsDeleted() bool           { return j.deleted }
func (j *BaseJob) IsReleased() bool          { return j.released }
func (j *BaseJob) HasFailed() bool           { return j.failed }
func (j *BaseJob) GetQueue() string          { return j.queue }
func (j *BaseJob) GetConnectionName() string { return j.connection }

func (j *BaseJob) Fire(ctx context.Context) error {
	if j.fireFunc != nil {
		return j.fireFunc(ctx)
	}

	return nil
}

func (j *BaseJob) Release(delay time.Duration) error {
	j.released = true

	if j.releaseFunc != nil {
		return j.releaseFunc(delay)
	}

	return nil
}

func (j *BaseJob) Delete() error {
	j.deleted = true

	if j.deleteFunc != nil {
		return j.deleteFunc()
	}

	return nil
}

func (j *BaseJob) Fail(err error) error {
	j.failed = true

	if j.failFunc != nil {
		return j.failFunc(err)
	}

	return nil
}

func (j *BaseJob) MarkAsFailed(err error) error {
	return j.Fail(err)
}
