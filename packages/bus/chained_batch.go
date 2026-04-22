package bus

import (
	"context"
)

// ChainedBatch is a queue-dispatchable wrapper around PendingBatch
// that enables batches inside job chains.
type ChainedBatch struct {
	Queueable
	Batchable

	jobs           []any
	name           string
	options        map[string]any
	catchCallbacks []FailureCallback
}

// NewChainedBatch creates a ChainedBatch from a PendingBatch,
// capturing its jobs, name, options, connection, queue, and catch callbacks.

// Handle dispatches the batch. If the batch completes successfully
// (is not cancelled), the remaining chain jobs are dispatched.

// When the batch succeeds, dispatch remaining chain jobs.

// Attach catch callbacks.

// ToPendingBatch reconstructs a PendingBatch from the ChainedBatch's stored config.

// PrepareNestedBatches converts any *PendingBatch items in a job slice
// into *ChainedBatch instances for queue dispatch.

// dispatcherContextKey is the context key for the QueueingDispatcher.
type dispatcherContextKey struct{}

func NewChainedBatch(pb *PendingBatch) *ChainedBatch {
	cb := &ChainedBatch{
		jobs:           pb.Jobs(),
		name:           pb.GetName(),
		options:        pb.Options(),
		catchCallbacks: pb.CatchCallbacks(),
	}

	cb.Connection = pb.Connection()
	cb.Queue = pb.Queue()

	return cb
}

func (cb *ChainedBatch) Handle(ctx context.Context) (any, error) {
	dispatcher, ok := ctx.Value(dispatcherContextKey{}).(QueueingDispatcher)

	if !ok {
		return nil, nil
	}

	pb := cb.ToPendingBatch(dispatcher)

	if len(cb.ChainJobs) > 0 {
		remaining := cb.ChainJobs
		pb.Then(func(ctx context.Context, batch *Batch) {
			if batch.Cancelled() {
				return
			}

			for _, job := range remaining {
				_ = dispatcher.DispatchToQueue(ctx, job)
			}
		})
	}

	for _, fn := range cb.catchCallbacks {
		pb.Catch(fn)
	}

	_, err := pb.Dispatch(ctx)

	return nil, err
}

func (cb *ChainedBatch) ToPendingBatch(dispatcher QueueingDispatcher) *PendingBatch {
	pb := NewPendingBatch(dispatcher, cb.jobs).Name(cb.name)

	if cb.Connection != "" {
		pb.OnConnection(cb.Connection)
	}

	if cb.Queue != "" {
		pb.OnQueue(cb.Queue)
	}

	if cb.options != nil {
		if v, ok := cb.options["allowFailures"]; ok && v == true {
			pb.AllowFailures()
		}
	}

	return pb
}

func PrepareNestedBatches(jobs []any) []any {
	result := make([]any, 0, len(jobs))

	for _, job := range jobs {
		if pb, ok := job.(*PendingBatch); ok {
			result = append(result, NewChainedBatch(pb))
		} else {
			result = append(result, job)
		}
	}

	return result
}

// WithDispatcher returns a context with the QueueingDispatcher set.
func WithDispatcher(ctx context.Context, d QueueingDispatcher) context.Context {
	return context.WithValue(ctx, dispatcherContextKey{}, d)
}
