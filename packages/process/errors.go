package process

import (
	"errors"
	"fmt"
)

// ErrProcessFailed is returned when a process exits unsuccessfully.

// ErrProcessTimedOut is returned when a process exceeds its timeout.

// ErrStrayProcess is returned when a command runs without a matching fake
// while stray processes are prevented.

// ErrSequenceEmpty is returned when a fake sequence has no result left.

// ProcessError wraps a failed process result.
type ProcessError struct {
	result *Result
	err    error
}

var (
	ErrProcessFailed = errors.New("process: failed")

	ErrProcessTimedOut = errors.New("process: timed out")

	ErrStrayProcess = errors.New("process: stray process")

	ErrSequenceEmpty = errors.New("process: fake sequence empty")
)

func (e *ProcessError) Error() string {
	if e == nil || e.result == nil {
		return ErrProcessFailed.Error()
	}

	return fmt.Sprintf("%s: %s exited with code %d", e.err, e.result.command.String(), e.result.exitCode)
}

func (e *ProcessError) Unwrap() error {
	if e == nil || e.err == nil {
		return ErrProcessFailed
	}

	return e.err
}

// Result returns the failed process result.
func (e *ProcessError) Result() *Result {
	if e == nil {
		return nil
	}

	return e.result
}
