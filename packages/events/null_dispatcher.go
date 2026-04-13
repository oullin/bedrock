package events

import (
	"context"

	cevents "github.com/bedrock/packages/contracts/events"
)

// NullDispatcher wraps an EventDispatcher and suppresses all event firing while
// still forwarding listener registrations. This is useful in testing scenarios
// where events should be captured but not actually dispatched.
type NullDispatcher struct {
	dispatcher *EventDispatcher
}

// NewNullDispatcher creates a NullDispatcher wrapping the given dispatcher.
func NewNullDispatcher(dispatcher *EventDispatcher) *NullDispatcher {
	return &NullDispatcher{dispatcher: dispatcher}
}

// compile-time interface check.
var _ cevents.Dispatcher = (*NullDispatcher)(nil)

// Listen delegates to the wrapped dispatcher so listeners are registered.
func (n *NullDispatcher) Listen(events any, listeners ...Listener) {
	n.dispatcher.Listen(events, listeners...)
}

// HasListeners delegates to the wrapped dispatcher.
func (n *NullDispatcher) HasListeners(event any) bool {
	return n.dispatcher.HasListeners(event)
}

// HasWildcardListeners delegates to the wrapped dispatcher.
func (n *NullDispatcher) HasWildcardListeners(event any) bool {
	return n.dispatcher.HasWildcardListeners(event)
}

// Subscribe delegates to the wrapped dispatcher.
func (n *NullDispatcher) Subscribe(subscriber Subscriber) {
	n.dispatcher.Subscribe(subscriber)
}

// Dispatch is a no-op. Events are not fired.
func (n *NullDispatcher) Dispatch(_ context.Context, _ any) ([]any, error) {
	return nil, nil
}

// Until is a no-op. Events are not fired.
func (n *NullDispatcher) Until(_ context.Context, _ any) (any, error) {
	return nil, nil
}

// Push is a no-op. Events are not queued.
func (n *NullDispatcher) Push(_ context.Context, _ any) {}

// Flush is a no-op. Pushed events are not dispatched.
func (n *NullDispatcher) Flush(_ context.Context, _ string) error {
	return nil
}

// Forget delegates to the wrapped dispatcher.
func (n *NullDispatcher) Forget(event any) {
	n.dispatcher.Forget(event)
}

// ForgetPushed delegates to the wrapped dispatcher.
func (n *NullDispatcher) ForgetPushed() {
	n.dispatcher.ForgetPushed()
}

// GetListeners delegates to the wrapped dispatcher.
func (n *NullDispatcher) GetListeners(event any) []Listener {
	return n.dispatcher.GetListeners(event)
}
