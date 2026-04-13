package bus

import (
	"context"
	"fmt"
)

// PendingChain is a fluent builder for dispatching a chain of sequential jobs.
type PendingChain struct {
	jobs           []any
	connection     string
	queue          string
	catchCallbacks []func(ctx context.Context, err error)
	dispatcher     Dispatcher
}

// NewPendingChain creates a PendingChain.
func NewPendingChain(dispatcher Dispatcher, jobs []any) *PendingChain {
	return &PendingChain{dispatcher: dispatcher, jobs: jobs}
}

// OnConnection sets the queue connection for the chain.
func (c *PendingChain) OnConnection(connection string) *PendingChain {
	c.connection = connection

	return c
}

// OnQueue sets the queue name for the chain.
func (c *PendingChain) OnQueue(queue string) *PendingChain {
	c.queue = queue

	return c
}

// Catch registers a callback invoked when a chained job fails.
func (c *PendingChain) Catch(fn func(ctx context.Context, err error)) *PendingChain {
	c.catchCallbacks = append(c.catchCallbacks, fn)

	return c
}

// Dispatch dispatches the first job in the chain, attaching the remaining
// jobs as the chain on the first job's Queueable.
func (c *PendingChain) Dispatch(ctx context.Context) (any, error) {
	if len(c.jobs) == 0 {
		return nil, fmt.Errorf("bus: cannot dispatch an empty chain")
	}

	first := c.prepareFirstJob()

	return c.dispatcher.Dispatch(ctx, first)
}

// DispatchAfterResponse dispatches the chain after the response is sent.
func (c *PendingChain) DispatchAfterResponse(ctx context.Context) error {
	if len(c.jobs) == 0 {
		return fmt.Errorf("bus: cannot dispatch an empty chain")
	}

	first := c.prepareFirstJob()

	return c.dispatcher.DispatchAfterResponse(ctx, first)
}

func (c *PendingChain) prepareFirstJob() any {
	first := c.jobs[0]
	remaining := c.jobs[1:]

	// If the first job has a Queueable, configure it.
	if q, ok := first.(interface{ Chain(jobs ...any) *Queueable }); ok {
		q.Chain(remaining...)
	}

	if c.connection != "" {
		if q, ok := first.(interface{ OnConnection(string) *Queueable }); ok {
			q.OnConnection(c.connection)
		}
	}

	if c.queue != "" {
		if q, ok := first.(interface{ OnQueue(string) *Queueable }); ok {
			q.OnQueue(c.queue)
		}
	}

	// Wire catch callbacks onto the first job.
	for _, fn := range c.catchCallbacks {
		if q, ok := first.(interface {
			OnChainCatch(func(context.Context, error)) *Queueable
		}); ok {
			q.OnChainCatch(fn)
		}
	}

	return first
}
