package events

import "context"

// Listener is a callback invoked when an event is dispatched. The return
// value is only meaningful for Until: a non-nil return halts propagation.
type Listener func(ctx context.Context, event string, payload any) (any, error)

// Dispatcher defines the event dispatching contract.
type Dispatcher interface {
	Listen(events []string, listener Listener)
	HasListeners(event string) bool
	GetListeners(event string) []Listener
	Subscribe(subscriber Subscriber)
	Dispatch(ctx context.Context, event string, payload any) error
	Until(ctx context.Context, event string, payload any) (any, error)
	Push(event string, payload any)
	FlushQueued(ctx context.Context, event string) error
	Forget(event string)
	Flush()
}

// Subscriber registers its own listeners with a dispatcher.
type Subscriber interface {
	Subscribe(dispatcher Dispatcher)
}
