package redis

import (
	"sync"
	"time"
)

// CommandExecuted mirrors Framework\Redis\Events\CommandExecuted. It is
// dispatched once per Connection.Command call when events are enabled.
type CommandExecuted struct {
	Command        string
	Parameters     []any
	Time           time.Duration
	ConnectionName string
}

// CommandFailed mirrors Framework\Redis\Events\CommandFailed.
type CommandFailed struct {
	Command        string
	Parameters     []any
	Exception      error
	ConnectionName string
}

// TimeMs returns the command duration in milliseconds (float64), matching
// Upstream's $time property.

// EventDispatcher fans out CommandExecuted events to registered listeners.
// A zero value is ready to use.
type EventDispatcher struct {
	mu               sync.RWMutex
	enabled          bool
	listeners        []func(CommandExecuted)
	failureListeners []func(CommandFailed)
}

func (e CommandExecuted) TimeMs() float64 {
	return float64(e.Time.Microseconds()) / 1000.0
}

// NewEventDispatcher returns a dispatcher with events disabled.
func NewEventDispatcher() *EventDispatcher { return &EventDispatcher{} }

// Enable turns event dispatching on.
func (d *EventDispatcher) Enable() {
	d.mu.Lock()
	d.enabled = true
	d.mu.Unlock()
}

// Disable turns event dispatching off. Registered listeners are kept.
func (d *EventDispatcher) Disable() {
	d.mu.Lock()
	d.enabled = false
	d.mu.Unlock()
}

// Enabled reports whether events are currently dispatched.
func (d *EventDispatcher) Enabled() bool {
	d.mu.RLock()

	defer d.mu.RUnlock()

	return d.enabled
}

// Listen registers a listener.
func (d *EventDispatcher) Listen(fn func(CommandExecuted)) {
	d.mu.Lock()
	d.listeners = append(d.listeners, fn)
	d.mu.Unlock()
}

// ListenForFailures registers a failure listener.
func (d *EventDispatcher) ListenForFailures(fn func(CommandFailed)) {
	d.mu.Lock()
	d.failureListeners = append(d.failureListeners, fn)
	d.mu.Unlock()
}

// DispatchExecuted fires the event to all command listeners when enabled.
func (d *EventDispatcher) DispatchExecuted(e CommandExecuted) {
	d.mu.RLock()

	if !d.enabled || len(d.listeners) == 0 {
		d.mu.RUnlock()

		return
	}

	ls := make([]func(CommandExecuted), len(d.listeners))
	copy(ls, d.listeners)
	d.mu.RUnlock()

	for _, l := range ls {
		l(e)
	}
}

// DispatchFailed fires the event to all failure listeners when enabled.
func (d *EventDispatcher) DispatchFailed(e CommandFailed) {
	d.mu.RLock()

	if !d.enabled || len(d.failureListeners) == 0 {
		d.mu.RUnlock()

		return
	}

	ls := make([]func(CommandFailed), len(d.failureListeners))
	copy(ls, d.failureListeners)
	d.mu.RUnlock()

	for _, l := range ls {
		l(e)
	}
}
