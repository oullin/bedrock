package events

import "context"

// Dispatcher dispatches domain events.
type Dispatcher interface {
	Dispatch(ctx context.Context, event any) error
}
