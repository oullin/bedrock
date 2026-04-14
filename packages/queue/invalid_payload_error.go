package queue

// InvalidPayloadError is the Go port of
// Framework\Queue\InvalidPayloadException.
//
// It is returned by Payload marshalling/unmarshalling helpers when a
// payload cannot be encoded or decoded. The Value field carries the
// original input that failed to decode (matching Upstream's $value field)
// so callers can inspect it without re-parsing.
type InvalidPayloadError struct {
	Message string
	Value   any
}

// Error implements the error interface.
func (e *InvalidPayloadError) Error() string {
	if e.Message == "" {
		return "queue: invalid payload"
	}

	return e.Message
}

// NewInvalidPayloadError constructs an error with an explicit message and
// the offending value. Mirrors Upstream's constructor signature
// (message, value).
func NewInvalidPayloadError(message string, value any) *InvalidPayloadError {
	return &InvalidPayloadError{Message: message, Value: value}
}
