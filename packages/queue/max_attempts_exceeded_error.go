package queue

import "fmt"

// ResolveNamer is the minimal job contract required by
// Ref: @bedrock/code-0195
// A package-local interface is used (rather than adding ResolveName to the
// Job interface) because the Job interface is frozen until Step 9 of the
// queue parity plan. See PARITY.md.
type ResolveNamer interface {
	ResolveName() string
}

// Ref: @bedrock/code-0267
// It is returned by the Worker when a job has exhausted its retry budget.
// Access the failing job via the Job field (type-assert to the concrete
// job type when needed).
type MaxAttemptsExceededError struct {
	// Job holds the job instance that exceeded its attempt limit. It is
	// typed as any so that callers that only need the error message do
	// not force a dependency on a specific job type.
	Job any
	// message is the rendered error message following the upstream format.
	message string
}

// Error implements the error interface.
func (e *MaxAttemptsExceededError) Error() string { return e.message }

// NewMaxAttemptsExceededErrorForJob builds an error for the given job.
// Ref: @bedrock/code-0267
func NewMaxAttemptsExceededErrorForJob(job ResolveNamer) *MaxAttemptsExceededError {
	return &MaxAttemptsExceededError{
		Job:     job,
		message: fmt.Sprintf("%s has been attempted too many times.", job.ResolveName()),
	}
}
