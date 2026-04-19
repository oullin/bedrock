package echo

import (
	"reflect"
	"sync"
)

// entry pairs a Callback with its identity pointer so it can be removed later.
type entry struct {
	ptr uintptr
	fn  Callback
}

// DispatchChannel is a Channel implementation with real listener storage and
// dispatch. It is the Go equivalent of SocketIoChannel's listener management
// layer, and is used as the base for real connector backends as well as for
// testing listener behaviour directly.
//
// Callback identity for StopListening is determined via
// reflect.ValueOf(cb).Pointer(). This works correctly when the same named
// variable is passed to both Listen and StopListening — the common pattern
// in tests and application code. It does not distinguish between distinct
// closures compiled from the same source-level literal.
type DispatchChannel struct {
	mu        sync.RWMutex
	formatter *EventFormatter
	listeners map[string][]entry // keyed by formatted event name
	allCBs    []entry            // ListenToAll callbacks
}

// compile-time interface check.
var _ Channel = (*DispatchChannel)(nil)

// NewDispatchChannel creates a DispatchChannel that formats event names using
// the provided EventFormatter.
func NewDispatchChannel(formatter *EventFormatter) *DispatchChannel {
	return &DispatchChannel{
		formatter: formatter,
		listeners: make(map[string][]entry),
	}
}

// Listen registers callback for the given event. The event name is normalized
// via the EventFormatter before storage.
func (c *DispatchChannel) Listen(event string, callback Callback) Channel {
	key := c.formatter.Format(event)

	c.mu.Lock()
	c.listeners[key] = append(c.listeners[key], entry{ptr: ptrOf(callback), fn: callback})
	c.mu.Unlock()

	return c
}

// StopListening removes callback from the listener list for event. If callback
// is nil, all listeners for the event are removed.
func (c *DispatchChannel) StopListening(event string, callback Callback) Channel {
	key := c.formatter.Format(event)

	c.mu.Lock()

	defer c.mu.Unlock()

	if callback == nil {
		delete(c.listeners, key)

		return c
	}

	p := ptrOf(callback)
	existing := c.listeners[key]
	kept := existing[:0]

	for _, e := range existing {
		if e.ptr != p {
			kept = append(kept, e)
		}
	}

	if len(kept) == 0 {
		delete(c.listeners, key)
	} else {
		c.listeners[key] = kept
	}

	return c
}

// ListenToAll registers callback to be invoked for every event dispatched on
// this channel.
func (c *DispatchChannel) ListenToAll(callback Callback) Channel {
	c.mu.Lock()
	c.allCBs = append(c.allCBs, entry{ptr: ptrOf(callback), fn: callback})
	c.mu.Unlock()

	return c
}

// StopListeningToAll removes the given callback from the "listen-all" list.
// If callback is nil, all "listen-all" callbacks are removed.
func (c *DispatchChannel) StopListeningToAll(callback Callback) Channel {
	c.mu.Lock()

	defer c.mu.Unlock()

	if callback == nil {
		c.allCBs = nil

		return c
	}

	p := ptrOf(callback)
	kept := c.allCBs[:0]

	for _, e := range c.allCBs {
		if e.ptr != p {
			kept = append(kept, e)
		}
	}

	c.allCBs = kept

	return c
}

// Subscribed registers callback to be invoked when the subscription is
// confirmed by the server. Stored under the reserved key "__subscribed__".
func (c *DispatchChannel) Subscribed(callback Callback) Channel {
	return c.Listen("__subscribed__", callback)
}

// Error registers callback to be invoked on subscription errors.
// Stored under the reserved key "__error__".
func (c *DispatchChannel) Error(callback Callback) Channel {
	return c.Listen("__error__", callback)
}

// On registers callback for a raw (already-formatted) event name, bypassing
// the EventFormatter. This mirrors the SocketIoChannel.on() method used for
// low-level transport events.
func (c *DispatchChannel) On(event string, callback Callback) Channel {
	c.mu.Lock()
	c.listeners[event] = append(c.listeners[event], entry{ptr: ptrOf(callback), fn: callback})
	c.mu.Unlock()

	return c
}

// Leave is a no-op on DispatchChannel; real connectors handle unsubscription.
func (c *DispatchChannel) Leave() {}

// Dispatch fires all registered callbacks for the given event. The event name
// is normalized via the EventFormatter before lookup. "Listen-all" callbacks
// are also invoked.
//
// This method is the primary test helper — it replaces the role of
// socket.emit() in the TypeScript SocketIoChannel test suite.
func (c *DispatchChannel) Dispatch(event string, data any) {
	key := c.formatter.Format(event)

	c.mu.RLock()
	entries := make([]entry, len(c.listeners[key]))
	copy(entries, c.listeners[key])
	all := make([]entry, len(c.allCBs))
	copy(all, c.allCBs)
	c.mu.RUnlock()

	for _, e := range entries {
		e.fn(data)
	}

	for _, e := range all {
		e.fn(data)
	}
}

// ptrOf returns a uintptr identifying the function value cb.
// Returns 0 for a nil callback.
func ptrOf(cb Callback) uintptr {
	if cb == nil {
		return 0
	}

	return reflect.ValueOf(cb).Pointer()
}
