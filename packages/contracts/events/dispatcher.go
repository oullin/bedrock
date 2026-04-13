package events

import "context"

// Listener handles an event. The returned value is collected by Dispatch and
// used by Until to halt on the first non-nil response.
type Listener func(ctx context.Context, event any) (any, error)

// Subscriber registers one or more event-listener mappings on a dispatcher.
type Subscriber interface {
	Subscribe(dispatcher Dispatcher)
}

// Dispatcher dispatches domain events and manages listeners.
type Dispatcher interface {
	// Listen registers one or more listeners for the given event(s).
	// events accepts a string name, a struct zero-value, or a slice of either.
	Listen(events any, listeners ...Listener)
	// HasListeners reports whether the event has any registered listeners.
	HasListeners(event any) bool
	// HasWildcardListeners reports whether any wildcard patterns match the event.
	HasWildcardListeners(event any) bool
	// Subscribe registers an event subscriber.
	Subscribe(subscriber Subscriber)
	// Until dispatches the event and stops on the first non-nil response.
	Until(ctx context.Context, event any) (any, error)
	// Dispatch fires the event to all registered listeners.
	Dispatch(ctx context.Context, event any) ([]any, error)
	// Push queues an event for later dispatch via Flush.
	Push(ctx context.Context, event any)
	// Flush dispatches all pushed events matching the given name.
	Flush(ctx context.Context, event string) error
	// Forget removes all listeners for the given event.
	Forget(event any)
	// ForgetPushed clears all pushed/deferred events.
	ForgetPushed()
	// GetListeners returns all listeners for the given event (direct + wildcard).
	GetListeners(event any) []Listener
}
