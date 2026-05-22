package queue

import (
	"errors"
	"fmt"
	"time"
)

// Ref: @bedrock/code-0261
// in their job-handler structs to get ergonomic access to the
// currently-processing Job (delete, release, fail) without having to
// re-implement the lifecycle plumbing.
//
// Typical use:
//
//	type SendEmailHandler struct {
//	    queue.InteractsWithQueue
//	    Mailer *Mailer
//	}
//
//	func (h *SendEmailHandler) Handle(ctx context.Context, job queue.Job) error {
//	    h.Job = job
//	    if err := h.Mailer.Send(...); err != nil {
//	        return h.Fail(err)
//	    }
//	    return nil
//	}
//
// The struct holds a Job field (set by the worker before invoking the
// handler) plus thin wrappers over Job.Release / Job.Delete / Job.Fail.
type InteractsWithQueue struct {
	Job Job
}

// Release re-queues the associated job after the given delay. Pass
// zero or omit the argument for an immediate release, matching
// the upstream $this->job->release($delay = 0) default.
//
// Returns nil when no job is attached.
func (i *InteractsWithQueue) Release(delay ...time.Duration) error {
	if i.Job == nil {
		return nil
	}

	var d time.Duration

	if len(delay) > 0 {
		d = delay[0]
	}

	return i.Job.Release(d)
}

// Delete removes the associated job from the queue. Returns nil when
// no job is attached.
func (i *InteractsWithQueue) Delete() error {
	if i.Job == nil {
		return nil
	}

	return i.Job.Delete()
}

// Fail marks the associated job as failed. reason is the same duck-
// typed argument the upstream accepts:
//
//   - nil          → ManuallyFailedError
//   - error        → passed through verbatim
//   - string       → wrapped with errors.New(reason)
//   - anything else → fmt.Errorf("%v", reason)
//
// Returns nil when no job is attached.
//
// Ref: @bedrock/code-0261
// including the "string becomes Exception" conversion tested by
// InteractsWithQueueTest::testCreatesAnExceptionFromString.
func (i *InteractsWithQueue) Fail(reason any) error {
	if i.Job == nil {
		return nil
	}

	return i.Job.Fail(coerceFailReason(reason))
}

// coerceFailReason converts a caller-supplied fail argument into the
// error value the Job.Fail contract expects.
func coerceFailReason(reason any) error {
	switch v := reason.(type) {
	case nil:
		return NewManuallyFailedError("")
	case error:
		return v
	case string:
		return errors.New(v)
	default:
		return fmt.Errorf("%v", v)
	}
}
