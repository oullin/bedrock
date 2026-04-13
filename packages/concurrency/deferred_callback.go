package concurrency

import (
	"context"
	"sync"
)

// DeferredCallback stores tasks for later execution.
// Call Flush to execute all stored tasks via the associated driver.
type DeferredCallback struct {
	mu      sync.Mutex
	tasks   []Task
	driver  Driver
	flushed bool
}

// NewDeferredCallback creates a DeferredCallback bound to the given driver.
func NewDeferredCallback(driver Driver, tasks []Task) *DeferredCallback {
	cp := make([]Task, len(tasks))
	copy(cp, tasks)

	return &DeferredCallback{
		driver: driver,
		tasks:  cp,
	}
}

// Flush executes all pending tasks via the driver and returns the results.
// Subsequent calls are no-ops that return nil, nil.
func (d *DeferredCallback) Flush(ctx context.Context) ([]any, error) {
	d.mu.Lock()

	defer d.mu.Unlock()

	if d.flushed || len(d.tasks) == 0 {
		return nil, nil
	}

	d.flushed = true

	tasks := d.tasks
	d.tasks = nil

	return d.driver.Run(ctx, tasks)
}

// Pending reports whether there are unflushed tasks.
func (d *DeferredCallback) Pending() bool {
	d.mu.Lock()

	defer d.mu.Unlock()

	return !d.flushed && len(d.tasks) > 0
}

// Count returns the number of pending tasks.
func (d *DeferredCallback) Count() int {
	d.mu.Lock()

	defer d.mu.Unlock()

	if d.flushed {
		return 0
	}

	return len(d.tasks)
}
