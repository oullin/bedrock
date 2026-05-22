package websockets

import "context"

// Dispatcher routes events to channel subscribers.
// Implementations may be synchronous (in-process) or distributed (Redis pub/sub).
type Dispatcher interface {
	// Dispatch broadcasts the event to all subscribers of the named channel
	// within the given application. If Event.SocketID is set the sender is
	// excluded from receiving the broadcast.
	Dispatch(ctx context.Context, appID string, event Event) error

	// Subscribe starts listening for cross-server events for the given app.
	// For synchronous dispatchers this is a no-op.
	Subscribe(ctx context.Context, appID string) error
}
