package bus

import (
	"context"
	"time"
)

// ShouldQueue is a marker interface for commands that should be dispatched to the queue.
type ShouldQueue interface {
	ShouldQueue()
}

// Handler handles a command/job.
type Handler func(ctx context.Context, command any) (any, error)

// Pipe is middleware in the command pipeline.
// It receives the command and a next function to call the remaining pipeline.
type Pipe func(ctx context.Context, command any, next Handler) (any, error)

// Dispatcher dispatches commands synchronously or to a queue.
type Dispatcher interface {
	// Dispatch sends a command through the pipeline and executes it.
	Dispatch(ctx context.Context, command any) (any, error)
	// DispatchSync executes the command synchronously, bypassing the queue.
	DispatchSync(ctx context.Context, command any) (any, error)
	// DispatchNow is an alias for DispatchSync.
	DispatchNow(ctx context.Context, command any) (any, error)
	// DispatchAfterResponse queues a command for execution after the response is sent.
	DispatchAfterResponse(ctx context.Context, command any) error
	// PipeThrough sets the middleware pipeline.
	PipeThrough(pipes ...Pipe) Dispatcher
	// Map registers a command→handler mapping.
	Map(command any, handler Handler) Dispatcher
	// HasCommandHandler reports whether a handler is registered for the command type.
	HasCommandHandler(command any) bool
	// GetCommandHandler returns the handler for the given command type.
	GetCommandHandler(command any) (Handler, bool)
	// Chain creates a PendingChain for sequential job execution.
	Chain(jobs []any) *PendingChain
}

// QueueingDispatcher extends Dispatcher with queue-based dispatch.
type QueueingDispatcher interface {
	Dispatcher
	// DispatchToQueue sends a command to the queue backend.
	DispatchToQueue(ctx context.Context, command any) error
	// FindBatch retrieves a batch by ID.
	FindBatch(ctx context.Context, id string) (*Batch, error)
	// Batch creates a PendingBatch for the given jobs.
	Batch(jobs []any) *PendingBatch
}

// BatchRepository persists batch state.
type BatchRepository interface {
	Get(ctx context.Context, id string) (*Batch, error)
	GetList(ctx context.Context, limit int, before string) ([]*Batch, error)
	Store(ctx context.Context, batch *Batch) error
	IncrementTotalJobs(ctx context.Context, id string, amount int) error
	DecrementPendingJobs(ctx context.Context, id string) (*UpdatedBatchJobCounts, error)
	IncrementFailedJobs(ctx context.Context, id string, failedJobID string) (*UpdatedBatchJobCounts, error)
	MarkAsFinished(ctx context.Context, id string) error
	Cancel(ctx context.Context, id string) error
	Delete(ctx context.Context, id string) error
	Transaction(ctx context.Context, fn func(BatchRepository) error) error
	RollBack(ctx context.Context) error
}

// PrunableBatchRepository extends BatchRepository with cleanup operations.
type PrunableBatchRepository interface {
	BatchRepository
	Prune(ctx context.Context, before time.Time) (int, error)
	PruneCancelled(ctx context.Context, before time.Time) (int, error)
	PruneUnfinished(ctx context.Context, before time.Time) (int, error)
}
