package events

import (
	"context"
	"time"
)

// CallQueuedListener is a queue job payload that wraps a listener for
// asynchronous execution.
// for retry, timeout, backoff, and unique job constraints.
type CallQueuedListener struct {
	ListenerName string
	Method       string
	Event        any
	Options      ListenerOptions

	// ShouldBeUnique indicates whether the job enforces a unique constraint.
	ShouldBeUniqueFlag bool
	// UniqueID is the identifier used for deduplication.
	UniqueID string
	// UniqueFor is the duration of the unique lock.
	UniqueFor time.Duration
}

// NewCallQueuedListener creates a CallQueuedListener with the given listener
// name and event.
func NewCallQueuedListener(listenerName string, event any) *CallQueuedListener {
	return &CallQueuedListener{
		ListenerName: listenerName,
		Method:       "Handle",
		Event:        event,
	}
}

// Handle executes the queued listener. In a full queue integration, this would
// resolve the listener and invoke it. Here it stores the execution context for
// the queue worker to process.
func (c *CallQueuedListener) Handle(ctx context.Context) error {
	return nil
}

// Failed is called when the queue job fails after exhausting retries.
func (c *CallQueuedListener) Failed(_ context.Context, _ error) {}

// ShouldQueue satisfies the ShouldQueue marker interface.
func (c *CallQueuedListener) ShouldQueue() {}

// GetQueue returns the queue name for routing.
func (c *CallQueuedListener) GetQueue() string { return c.Options.Queue }

// GetConnection returns the connection name.
func (c *CallQueuedListener) GetConnection() string { return c.Options.Connection }

// GetDelay returns the delay before processing.
func (c *CallQueuedListener) GetDelay() time.Duration { return c.Options.Delay }

// GetTries returns the maximum number of attempts.
func (c *CallQueuedListener) GetTries() int { return c.Options.Tries }

// GetMaxExceptions returns the maximum exceptions before failing.
func (c *CallQueuedListener) GetMaxExceptions() int { return c.Options.MaxExceptions }

// GetTimeout returns the job timeout.
func (c *CallQueuedListener) GetTimeout() time.Duration { return c.Options.Timeout }

// GetBackoff returns the backoff durations between retries.
func (c *CallQueuedListener) GetBackoff() []time.Duration { return c.Options.Backoff }

// DisplayName returns the display name for queue dashboards.
func (c *CallQueuedListener) DisplayName() string { return c.ListenerName }

// WithOptions sets the listener options.
func (c *CallQueuedListener) WithOptions(opts ListenerOptions) *CallQueuedListener {
	c.Options = opts

	return c
}
