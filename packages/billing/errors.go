package billing

import (
	"errors"
	"fmt"
	"strings"
)

// Sentinel errors.

// PaddleError wraps an error from the Paddle API.
type PaddleError struct {
	Code    string
	Message string
	Err     error
}

// ValidationError represents a single field validation failure.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors is a collection of validation errors.
type ValidationErrors []ValidationError

var (
	ErrNotFound          = errors.New("billing: not found")
	ErrBillableRequired  = errors.New("billing: billable is required")
	ErrAlreadySubscribed = errors.New("billing: already subscribed")
	ErrNotSubscribed     = errors.New("billing: not subscribed")
	ErrInvalidProvider   = errors.New("billing: invalid payment provider")
)

func (e *PaddleError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("paddle: %s (%s): %v", e.Message, e.Code, e.Err)
	}

	return fmt.Sprintf("paddle: %s (%s)", e.Message, e.Code)
}

func (e *PaddleError) Unwrap() error {
	return e.Err
}

// Error implements the error interface.
func (ve ValidationErrors) Error() string {
	msgs := make([]string, len(ve))

	for i, e := range ve {
		msgs[i] = e.Field + ": " + e.Message
	}

	return strings.Join(msgs, "; ")
}

// HasErrors reports whether there are any validation errors.
func (ve ValidationErrors) HasErrors() bool {
	return len(ve) > 0
}
