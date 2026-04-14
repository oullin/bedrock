package queue

import "fmt"

// TimeoutExceededError is the Go port of
// Illuminate\Queue\TimeoutExceededException.
//
// In Laravel this class extends MaxAttemptsExceededException. Go has no
// inheritance so TimeoutExceededError embeds *MaxAttemptsExceededError by
// pointer; errors.As therefore matches both types, preserving the
// Laravel "is-a" relationship. The Error() method is promoted from the
// embedded pointer.
type TimeoutExceededError struct {
	*MaxAttemptsExceededError
}

// Unwrap exposes the embedded MaxAttemptsExceededError so errors.As and
// errors.Is walk through it, preserving Laravel's PHP "TimeoutExceededException
// extends MaxAttemptsExceededException" relationship.
func (e *TimeoutExceededError) Unwrap() error { return e.MaxAttemptsExceededError }

// NewTimeoutExceededErrorForJob builds an error for the given job.
// Mirrors Laravel's TimeoutExceededException::forJob static factory.
func NewTimeoutExceededErrorForJob(job ResolveNamer) *TimeoutExceededError {
	return &TimeoutExceededError{
		MaxAttemptsExceededError: &MaxAttemptsExceededError{
			Job:     job,
			message: fmt.Sprintf("%s has timed out.", job.ResolveName()),
		},
	}
}
