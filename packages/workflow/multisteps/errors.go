package multisteps

import "fmt"

// WorkflowError wraps a job failure with its attempt count and originating job name.
type WorkflowError struct {
	Job      string
	Attempts int
	Cause    error
}

// UnresolvedResponseError is returned when a Response Arg cannot find its field.
type UnresolvedResponseError struct {
	Job   string
	Field string
}

// CompileError reports static errors detected when building the DAG.
type CompileError struct {
	Reason string
}

func (e *WorkflowError) Error() string {
	return fmt.Sprintf("multisteps: job %q failed after %d attempt(s): %v", e.Job, e.Attempts, e.Cause)
}

func (e *WorkflowError) Unwrap() error { return e.Cause }

func (e *UnresolvedResponseError) Error() string {
	if e.Job == "" {
		return fmt.Sprintf("multisteps: response field %q not resolvable", e.Field)
	}

	return fmt.Sprintf("multisteps: response field %q on job %q not resolvable", e.Field, e.Job)
}

func (e *CompileError) Error() string {
	return "multisteps: compile error: " + e.Reason
}
