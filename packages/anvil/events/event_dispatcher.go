package events

import (
	"context"
	"strings"
	"sync"
)

type queuedEvent struct {
	event   string
	payload any
}

// EventDispatcher is the concrete event dispatcher. It is safe for
// concurrent use.
type EventDispatcher struct {
	mu        sync.RWMutex
	listeners map[string][]Listener
	wildcards map[string][]Listener
	queued    []queuedEvent
}

// New creates an empty dispatcher.
func New() *EventDispatcher {
	return &EventDispatcher{
		listeners: make(map[string][]Listener),
		wildcards: make(map[string][]Listener),
	}
}

// Listen registers a listener for one or more events. Event names containing
// a trailing ".*" or bare "*" are treated as wildcard patterns.
func (d *EventDispatcher) Listen(events []string, listener Listener) {
	d.mu.Lock()
	defer d.mu.Unlock()

	for _, event := range events {
		if isWildcard(event) {
			d.wildcards[event] = append(d.wildcards[event], listener)
		} else {
			d.listeners[event] = append(d.listeners[event], listener)
		}
	}
}

// HasListeners reports whether any listeners (exact or wildcard) match event.
func (d *EventDispatcher) HasListeners(event string) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()

	if len(d.listeners[event]) > 0 {
		return true
	}

	for pattern := range d.wildcards {
		if matchesWildcard(pattern, event) {
			return true
		}
	}

	return false
}

// GetListeners returns all listeners that match event, exact first then
// wildcard, in registration order.
func (d *EventDispatcher) GetListeners(event string) []Listener {
	return d.getListeners(event)
}

// Subscribe calls subscriber.Subscribe with the dispatcher so the subscriber
// can register its own listeners.
func (d *EventDispatcher) Subscribe(subscriber Subscriber) {
	subscriber.Subscribe(d)
}

// Dispatch fires all matching listeners in registration order. It returns
// the first listener error encountered and short-circuits on context
// cancellation. If a listener returns false (bool), propagation stops.
func (d *EventDispatcher) Dispatch(ctx context.Context, event string, payload any) error {
	listeners := d.getListeners(event)

	for _, l := range listeners {
		if err := ctx.Err(); err != nil {
			return err
		}

		result, err := l(ctx, event, payload)
		if err != nil {
			return err
		}

		if b, ok := result.(bool); ok && !b {
			break
		}
	}

	return nil
}

// Until dispatches listeners until one returns a non-nil value. It returns
// that value, or nil if all listeners returned nil. Returning false (bool)
// also halts propagation and is returned.
func (d *EventDispatcher) Until(ctx context.Context, event string, payload any) (any, error) {
	listeners := d.getListeners(event)

	for _, l := range listeners {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		result, err := l(ctx, event, payload)
		if err != nil {
			return nil, err
		}

		if result != nil {
			return result, nil
		}
	}

	return nil, nil
}

// Push queues an event for later dispatch via FlushQueued.
func (d *EventDispatcher) Push(event string, payload any) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.queued = append(d.queued, queuedEvent{event: event, payload: payload})
}

// FlushQueued dispatches all queued events matching the given event name
// and removes them from the queue. Pass "*" to flush all queued events.
func (d *EventDispatcher) FlushQueued(ctx context.Context, event string) error {
	d.mu.Lock()
	var toDispatch []queuedEvent
	var remaining []queuedEvent

	for _, q := range d.queued {
		if event == "*" || q.event == event {
			toDispatch = append(toDispatch, q)
		} else {
			remaining = append(remaining, q)
		}
	}

	d.queued = remaining
	d.mu.Unlock()

	for _, q := range toDispatch {
		if err := d.Dispatch(ctx, q.event, q.payload); err != nil {
			return err
		}
	}

	return nil
}

// Forget removes all listeners for the given event (exact match only).
func (d *EventDispatcher) Forget(event string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if isWildcard(event) {
		delete(d.wildcards, event)
	} else {
		delete(d.listeners, event)
	}
}

// Flush removes all listeners, wildcard listeners, and queued events.
func (d *EventDispatcher) Flush() {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.listeners = make(map[string][]Listener)
	d.wildcards = make(map[string][]Listener)
	d.queued = nil
}

// getListeners returns exact listeners followed by matching wildcard listeners.
func (d *EventDispatcher) getListeners(event string) []Listener {
	d.mu.RLock()
	defer d.mu.RUnlock()

	exact := make([]Listener, len(d.listeners[event]))
	copy(exact, d.listeners[event])

	for pattern, wls := range d.wildcards {
		if matchesWildcard(pattern, event) {
			exact = append(exact, wls...)
		}
	}

	return exact
}

// isWildcard reports whether event contains a wildcard segment.
func isWildcard(event string) bool {
	return strings.Contains(event, "*")
}

// matchesWildcard checks whether event matches a wildcard pattern. A bare "*"
// matches everything. A pattern like "user.*" matches any event that starts
// with "user." (one or more trailing segments).
func matchesWildcard(pattern, event string) bool {
	if pattern == "*" {
		return true
	}

	prefix := strings.TrimSuffix(pattern, "*")

	return strings.HasPrefix(event, prefix)
}
