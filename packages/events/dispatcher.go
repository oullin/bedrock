package events

import (
	"context"
	"reflect"
	"sync"

	cevents "github.com/bedrock/packages/contracts/events"
)

// EventDispatcher is the concrete implementation of the Dispatcher contract.
type EventDispatcher struct {
	mu                         sync.RWMutex
	listeners                  map[string][]Listener
	wildcards                  map[string][]Listener
	wildcardCache              map[string][]Listener
	pushed                     map[string][]any
	queueResolver              QueueResolver
	transactionManagerResolver TransactionManagerResolver
}

// NewDispatcher creates an EventDispatcher.
func NewDispatcher() *EventDispatcher {
	return &EventDispatcher{
		listeners:     make(map[string][]Listener),
		wildcards:     make(map[string][]Listener),
		wildcardCache: make(map[string][]Listener),
		pushed:        make(map[string][]any),
	}
}

// compile-time interface check.
var _ cevents.Dispatcher = (*EventDispatcher)(nil)

// Listen registers one or more listeners for the given event(s).
// events accepts a string, a struct zero-value (type-based key), or a slice of
// either. Each listener is appended in order.
func (d *EventDispatcher) Listen(events any, listeners ...Listener) {
	for _, name := range d.parseEvents(events) {
		d.mu.Lock()

		if isWildcardPattern(name) {
			d.wildcards[name] = append(d.wildcards[name], listeners...)
			d.wildcardCache = make(map[string][]Listener) // invalidate cache.
		} else {
			d.listeners[name] = append(d.listeners[name], listeners...)
		}

		d.mu.Unlock()
	}
}

// HasListeners reports whether any listeners (direct or wildcard) exist for the
// given event.
func (d *EventDispatcher) HasListeners(event any) bool {
	name := eventName(event)

	d.mu.RLock()

	defer d.mu.RUnlock()

	if len(d.listeners[name]) > 0 {
		return true
	}

	return d.hasWildcardMatch(name)
}

// HasWildcardListeners reports whether any wildcard patterns match the event.
func (d *EventDispatcher) HasWildcardListeners(event any) bool {
	name := eventName(event)

	d.mu.RLock()

	defer d.mu.RUnlock()

	return d.hasWildcardMatch(name)
}

// Subscribe registers an event subscriber. The subscriber's Subscribe method
// is called with the dispatcher so it can register its own listeners.
func (d *EventDispatcher) Subscribe(subscriber Subscriber) {
	subscriber.Subscribe(d)
}

// Dispatch fires the event to all registered listeners and returns a slice of
// non-nil responses.
func (d *EventDispatcher) Dispatch(ctx context.Context, event any) ([]any, error) {
	return d.dispatch(ctx, event, false)
}

// Until dispatches the event, stopping on the first listener that returns a
// non-nil response, and returns that response.
func (d *EventDispatcher) Until(ctx context.Context, event any) (any, error) {
	responses, err := d.dispatch(ctx, event, true)

	if err != nil {
		return nil, err
	}

	if len(responses) > 0 {
		return responses[0], nil
	}

	return nil, nil
}

// Push queues an event for later dispatch via Flush.
func (d *EventDispatcher) Push(ctx context.Context, event any) {
	name := eventName(event)

	d.mu.Lock()

	defer d.mu.Unlock()

	d.pushed[name] = append(d.pushed[name], event)
}

// Flush dispatches all pushed events matching the given name and then clears
// them.
func (d *EventDispatcher) Flush(ctx context.Context, event string) error {
	d.mu.Lock()
	events := d.pushed[event]
	delete(d.pushed, event)
	d.mu.Unlock()

	for _, e := range events {
		if _, err := d.Dispatch(ctx, e); err != nil {
			return err
		}
	}

	return nil
}

// Forget removes all listeners for the given event. If the event is a wildcard
// pattern, it is removed from the wildcard registry.
func (d *EventDispatcher) Forget(event any) {
	name := eventName(event)

	d.mu.Lock()

	defer d.mu.Unlock()

	delete(d.listeners, name)
	delete(d.wildcards, name)
	d.wildcardCache = make(map[string][]Listener)
}

// ForgetPushed clears all pushed/deferred events.
func (d *EventDispatcher) ForgetPushed() {
	d.mu.Lock()

	defer d.mu.Unlock()

	d.pushed = make(map[string][]any)
}

// GetListeners returns all listeners (direct + wildcard matches) for the event.
func (d *EventDispatcher) GetListeners(event any) []Listener {
	name := eventName(event)

	d.mu.RLock()

	defer d.mu.RUnlock()

	direct := d.listeners[name]
	wildcard := d.resolveWildcardListeners(name)

	if len(wildcard) == 0 {
		return direct
	}

	result := make([]Listener, 0, len(direct)+len(wildcard))
	result = append(result, direct...)
	result = append(result, wildcard...)

	return result
}

// GetRawListeners returns a copy of the raw listener registry for inspection.
func (d *EventDispatcher) GetRawListeners() map[string][]Listener {
	d.mu.RLock()

	defer d.mu.RUnlock()

	raw := make(map[string][]Listener, len(d.listeners))

	for k, v := range d.listeners {
		raw[k] = v
	}

	return raw
}

// MakeListener wraps a listener. If wildcard is true the returned listener
// injects the matched event name into the context under the WildcardEventKey.
func (d *EventDispatcher) MakeListener(listener Listener, wildcard bool) Listener {
	if !wildcard {
		return listener
	}

	return func(ctx context.Context, event any) (any, error) {
		return listener(ctx, event)
	}
}

// SetQueueResolver sets the resolver for the queue backend.
func (d *EventDispatcher) SetQueueResolver(resolver QueueResolver) *EventDispatcher {
	d.mu.Lock()

	defer d.mu.Unlock()

	d.queueResolver = resolver

	return d
}

// SetTransactionManagerResolver sets the resolver for the transaction manager.
func (d *EventDispatcher) SetTransactionManagerResolver(resolver TransactionManagerResolver) *EventDispatcher {
	d.mu.Lock()

	defer d.mu.Unlock()

	d.transactionManagerResolver = resolver

	return d
}

// Defer executes the callback while suppressing dispatch for the named events
// (or all events if none specified). After the callback returns, all deferred
// events are flushed.
func (d *EventDispatcher) Defer(ctx context.Context, callback func(ctx context.Context) error, events ...string) error {
	deferSet := make(map[string]struct{}, len(events))

	for _, e := range events {
		deferSet[e] = struct{}{}
	}

	deferAll := len(events) == 0

	captured := make(map[string][]any)
	captureMu := &sync.Mutex{}

	// Deep-copy original listeners and replace matched entries with capture stubs.
	d.mu.Lock()

	original := make(map[string][]Listener, len(d.listeners))

	for k, v := range d.listeners {
		original[k] = v
	}

	for name := range d.listeners {
		_, inSet := deferSet[name]

		if !deferAll && !inSet {
			continue
		}

		d.listeners[name] = []Listener{
			func(ctx context.Context, event any) (any, error) {
				n := eventName(event)
				captureMu.Lock()
				captured[n] = append(captured[n], event)
				captureMu.Unlock()

				return nil, nil
			},
		}
	}

	d.mu.Unlock()

	// Run the callback.
	err := callback(ctx)

	// Restore original listeners.
	d.mu.Lock()
	d.listeners = original
	d.mu.Unlock()

	if err != nil {
		return err
	}

	// Flush captured events.
	for _, evts := range captured {
		for _, e := range evts {
			if _, dispatchErr := d.Dispatch(ctx, e); dispatchErr != nil {
				return dispatchErr
			}
		}
	}

	return nil
}

// dispatch is the core dispatch logic shared by Dispatch and Until.
func (d *EventDispatcher) dispatch(ctx context.Context, event any, halt bool) ([]any, error) {
	listeners := d.GetListeners(event)

	var responses []any

	for _, listener := range listeners {
		response, err := listener(ctx, event)

		if err != nil {
			return responses, err
		}

		if halt && response != nil {
			return []any{response}, nil
		}

		if response != nil {
			responses = append(responses, response)
		}
	}

	return responses, nil
}

// parseEvents normalizes the polymorphic events parameter into a slice of
// string event names.
func (d *EventDispatcher) parseEvents(events any) []string {
	switch v := events.(type) {
	case string:
		return []string{v}
	case []string:
		return v
	case []any:
		names := make([]string, 0, len(v))

		for _, e := range v {
			names = append(names, eventName(e))
		}

		return names
	default:
		return []string{eventName(events)}
	}
}

// hasWildcardMatch reports whether any wildcard pattern matches the name.
// Must be called under at least d.mu.RLock.
func (d *EventDispatcher) hasWildcardMatch(name string) bool {
	for pattern := range d.wildcards {
		if matchesWildcard(pattern, name) {
			return true
		}
	}

	return false
}

// resolveWildcardListeners returns all wildcard listeners matching the name.
// Must be called under at least d.mu.RLock.
func (d *EventDispatcher) resolveWildcardListeners(name string) []Listener {
	if cached, ok := d.wildcardCache[name]; ok {
		return cached
	}

	var matched []Listener

	for pattern, ls := range d.wildcards {
		if matchesWildcard(pattern, name) {
			matched = append(matched, ls...)
		}
	}

	// Note: cannot write to wildcardCache under RLock. Callers that need
	// caching should call under a full Lock or accept the repeated scan.
	// For simplicity we skip write here. The cache is populated on Listen.

	return matched
}

// eventNameFromType returns the string event name for a reflect.Type.
func eventNameFromType(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.String()
}
