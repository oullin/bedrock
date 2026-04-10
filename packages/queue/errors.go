package queue

import "errors"

var (
	// ErrNoJob is returned when the queue is empty.
	ErrNoJob = errors.New("queue: no job available")
	// ErrInvalidDriver is returned when no driver is registered for the given name.
	ErrInvalidDriver = errors.New("queue: invalid driver")
)
