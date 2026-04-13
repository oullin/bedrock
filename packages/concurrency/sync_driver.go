package concurrency

import (
	"context"
	"fmt"
)

// SyncDriver executes tasks sequentially. Useful for testing and debugging.
type SyncDriver struct{}

// NewSyncDriver creates a SyncDriver.
func NewSyncDriver() *SyncDriver {
	return &SyncDriver{}
}

// Run executes tasks sequentially and returns results ordered by task index.
// Stops on the first error or context cancellation.
func (d *SyncDriver) Run(ctx context.Context, tasks []Task) ([]any, error) {
	if len(tasks) == 0 {
		return nil, ErrNoTasks
	}

	results := make([]any, len(tasks))

	for i, task := range tasks {
		if err := ctx.Err(); err != nil {
			return results, err
		}

		val, err := d.safeCall(task)

		if err != nil {
			return results, err
		}

		results[i] = val
	}

	return results, nil
}

// Defer stores tasks for later execution.
func (d *SyncDriver) Defer(tasks []Task) Deferrable {
	return NewDeferredCallback(d, tasks)
}

// safeCall invokes the task, recovering from panics.
func (d *SyncDriver) safeCall(fn Task) (val any, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%w: %v", ErrTaskPanicked, r)
		}
	}()

	return fn()
}
