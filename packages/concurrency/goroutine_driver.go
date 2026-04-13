package concurrency

import (
	"context"
	"fmt"
	"sync"
)

// GoroutineDriver executes tasks concurrently using goroutines.
// Set maxConcurrency > 0 to limit the number of simultaneous goroutines.
type GoroutineDriver struct {
	maxConcurrency int
}

// NewGoroutineDriver creates a GoroutineDriver.
// A maxConcurrency of 0 means unlimited concurrency.
func NewGoroutineDriver(maxConcurrency int) *GoroutineDriver {
	return &GoroutineDriver{maxConcurrency: maxConcurrency}
}

// Run executes all tasks concurrently and returns results ordered by task index.
// On the first task error, remaining tasks are cancelled via the context.
// Panics in tasks are recovered and returned as errors wrapping ErrTaskPanicked.
func (d *GoroutineDriver) Run(ctx context.Context, tasks []Task) ([]any, error) {
	if len(tasks) == 0 {
		return nil, ErrNoTasks
	}

	if err := ctx.Err(); err != nil {
		return make([]any, len(tasks)), err
	}

	results := make([]any, len(tasks))

	ctx, cancel := context.WithCancel(ctx)

	defer cancel()

	var (
		wg       sync.WaitGroup
		once     sync.Once
		firstErr error
		sem      chan struct{}
	)

	if d.maxConcurrency > 0 {
		sem = make(chan struct{}, d.maxConcurrency)
	}

	for i, task := range tasks {
		wg.Add(1)

		go func(idx int, fn Task) {
			defer wg.Done()

			if sem != nil {
				select {
				case sem <- struct{}{}:
					defer func() { <-sem }()
				case <-ctx.Done():
					once.Do(func() { firstErr = ctx.Err() })

					return
				}
			}

			if ctx.Err() != nil {
				return
			}

			val, err := d.safeCall(fn)

			if err != nil {
				once.Do(func() {
					firstErr = err
					cancel()
				})

				return
			}

			results[idx] = val
		}(i, task)
	}

	wg.Wait()

	if firstErr != nil {
		return results, firstErr
	}

	return results, nil
}

// Defer stores tasks for later execution.
func (d *GoroutineDriver) Defer(tasks []Task) Deferrable {
	return NewDeferredCallback(d, tasks)
}

// safeCall invokes the task, recovering from panics.
func (d *GoroutineDriver) safeCall(fn Task) (val any, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%w: %v", ErrTaskPanicked, r)
		}
	}()

	return fn()
}
