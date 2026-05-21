package queue

// Ref: @bedrock/code-0266
// In the upstream framework this is an empty RuntimeException used as a marker type
// when a job calls $this->fail() without an underlying cause. Go's
// equivalent is a small error struct with a default message; callers
// use errors.As to detect a manual failure.
type ManuallyFailedError struct {
	Message string
}

// Error implements the error interface.
func (e *ManuallyFailedError) Error() string {
	if e.Message == "" {
		return "queue: job manually failed"
	}

	return e.Message
}

// NewManuallyFailedError constructs a manually-failed marker error.
func NewManuallyFailedError(message string) *ManuallyFailedError {
	return &ManuallyFailedError{Message: message}
}
