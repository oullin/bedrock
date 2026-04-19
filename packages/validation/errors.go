package validation

import "errors"

// ErrValidationFailed is returned by Validator.Validate() when any rule fails.
// Callers that need the full error list should use errors.As to unwrap a
// *ValidationException.

// ValidationException wraps the MessageBag produced by a failed Validator.
// It implements the error interface so it can be returned from Validate().
type ValidationException struct {
	Bag *MessageBag
}

var ErrValidationFailed = errors.New("validation: validation failed")

func (e *ValidationException) Error() string {
	return ErrValidationFailed.Error()
}

// Unwrap allows errors.Is(err, ErrValidationFailed) to work.
func (e *ValidationException) Unwrap() error {
	return ErrValidationFailed
}
