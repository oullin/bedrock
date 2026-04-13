package concurrency

import "errors"

var (
	// ErrTaskPanicked indicates a task recovered from a panic.
	ErrTaskPanicked = errors.New("concurrency: task panicked")

	// ErrInvalidDriver indicates the requested driver is not registered.
	ErrInvalidDriver = errors.New("concurrency: invalid driver")

	// ErrNoTasks indicates an empty task slice was provided.
	ErrNoTasks = errors.New("concurrency: no tasks provided")
)
