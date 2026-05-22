package queue

import (
	"context"
	"errors"
)

// ExceptionPredicate decides whether err should fail job immediately.
// It gives Go callers the same "inspect the job object" escape hatch
// the upstream FailOnException middleware exposes through closures.
type ExceptionPredicate func(err error, job Job) bool

// FailOnException marks a job as failed when a handler returns a matching
// error. Sentinel errors are matched with errors.Is; predicates can inspect
// both the error and the job.
type FailOnException struct {
	errors     []error
	predicates []ExceptionPredicate
}

// NewFailOnException creates middleware that fails on the supplied sentinels.
func NewFailOnException(errs ...error) *FailOnException {
	out := &FailOnException{}

	for _, err := range errs {
		if err != nil {
			out.errors = append(out.errors, err)
		}
	}

	return out
}

// When appends an inspection predicate to the middleware.
func (m *FailOnException) When(predicate ExceptionPredicate) *FailOnException {
	if predicate != nil {
		m.predicates = append(m.predicates, predicate)
	}

	return m
}

// ShouldFail reports whether err should immediately fail job.
func (m *FailOnException) ShouldFail(err error, job Job) bool {
	if m == nil || err == nil {
		return false
	}

	for _, target := range m.errors {
		if errors.Is(err, target) {
			return true
		}
	}

	for _, predicate := range m.predicates {
		if predicate(err, job) {
			return true
		}
	}

	return false
}

// Handle runs next and fails job when the returned error matches.
func (m *FailOnException) Handle(ctx context.Context, job Job, next func(context.Context, Job) error) error {
	err := next(ctx, job)

	if m.ShouldFail(err, job) {
		_ = job.Fail(err)
	}

	return err
}
