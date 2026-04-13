package concurrency

import (
	"context"

	"github.com/bedrock/packages/contracts/concurrency"
)

// Task is a unit of concurrent work.
type Task = concurrency.Task

// Driver defines the concurrency backend contract.
type Driver = concurrency.Driver

// Deferrable holds deferred tasks and executes them on demand.
type Deferrable = concurrency.Deferrable

// Ensure compile-time interface satisfaction for concrete types.
var (
	_ Driver     = (*GoroutineDriver)(nil)
	_ Driver     = (*SyncDriver)(nil)
	_ Deferrable = (*DeferredCallback)(nil)
)

// Run is a convenience function that executes tasks using the GoroutineDriver.
func Run(ctx context.Context, tasks []Task) ([]any, error) {
	return NewGoroutineDriver(0).Run(ctx, tasks)
}

// Defer is a convenience function that defers tasks using the GoroutineDriver.
func Defer(tasks []Task) Deferrable {
	return NewGoroutineDriver(0).Defer(tasks)
}
