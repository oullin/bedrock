package events

import (
	"context"
	"time"
)

// QueuedClosure wraps a Listener for queued execution with fluent configuration
// for connection, queue, delay, and failure handling.
type QueuedClosure struct {
	closure    Listener
	connection string
	queue      string
	delay      time.Duration
	catchFn    func(ctx context.Context, err error)
}

// Queueable creates a QueuedClosure from a listener. This is the Go equivalent
// of the upstream queueable() helper function.
func Queueable(listener Listener) *QueuedClosure {
	return &QueuedClosure{
		closure: listener,
	}
}

// OnConnection sets the queue connection name.
func (q *QueuedClosure) OnConnection(connection string) *QueuedClosure {
	q.connection = connection

	return q
}

// OnQueue sets the queue name.
func (q *QueuedClosure) OnQueue(queue string) *QueuedClosure {
	q.queue = queue

	return q
}

// WithDelay sets the delay before the job is processed.
func (q *QueuedClosure) WithDelay(delay time.Duration) *QueuedClosure {
	q.delay = delay

	return q
}

// Catch registers a failure handler that is called if the queued listener fails.
func (q *QueuedClosure) Catch(fn func(ctx context.Context, err error)) *QueuedClosure {
	q.catchFn = fn

	return q
}

// Resolve returns a Listener that, when called, invokes the wrapped closure.
// In a queue-integrated environment, the returned listener would create a queue
// job instead of executing inline.
func (q *QueuedClosure) Resolve() Listener {
	return q.closure
}

// GetConnection returns the configured queue connection.
func (q *QueuedClosure) GetConnection() string { return q.connection }

// GetQueue returns the configured queue name.
func (q *QueuedClosure) GetQueue() string { return q.queue }

// GetDelay returns the configured delay.
func (q *QueuedClosure) GetDelay() time.Duration { return q.delay }

// GetCatchFn returns the configured failure handler.
func (q *QueuedClosure) GetCatchFn() func(ctx context.Context, err error) { return q.catchFn }
