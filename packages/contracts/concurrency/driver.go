package concurrency

import "context"

// Task is a unit of concurrent work.
type Task func() (any, error)

// Driver defines the concurrency backend contract.
type Driver interface {
	// Run executes the given tasks concurrently and returns results ordered by task index.
	Run(ctx context.Context, tasks []Task) ([]any, error)

	// Defer stores tasks for later execution, returning a handle to flush them.
	Defer(tasks []Task) Deferrable
}

// Deferrable holds deferred tasks and executes them on demand.
type Deferrable interface {
	Flush(ctx context.Context) ([]any, error)
	Pending() bool
	Count() int
}
