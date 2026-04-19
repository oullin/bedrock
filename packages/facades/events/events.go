// Package events is the facade for the event dispatcher. It forwards common
// operations to the EventDispatcher bound under "events" in the global
// bedrock Application.
package events

import (
	"context"
	"sync"

	"github.com/bedrock/packages/container"
	eventspkg "github.com/bedrock/packages/events"
)

var (
	mu     sync.Mutex
	cached *eventspkg.EventDispatcher
)

// Dispatcher returns the event dispatcher from the global Application.
// Resolved once per process and cached.
func Dispatcher() *eventspkg.EventDispatcher {
	mu.Lock()

	defer mu.Unlock()

	if cached == nil {
		cached = container.Resolve[*eventspkg.EventDispatcher]("events")
	}

	return cached
}

// Reset clears the cached dispatcher. Tests must call this after reinstalling
// a different Application via container.SetApp.
func Reset() {
	mu.Lock()

	defer mu.Unlock()

	cached = nil
}

// Listen registers listeners for the given event(s).
func Listen(eventList any, listeners ...eventspkg.Listener) {
	Dispatcher().Listen(eventList, listeners...)
}

// Dispatch synchronously dispatches the event to all listeners.
func Dispatch(ctx context.Context, event any) ([]any, error) {
	return Dispatcher().Dispatch(ctx, event)
}

// Until dispatches the event and stops at the first non-nil response.
func Until(ctx context.Context, event any) (any, error) {
	return Dispatcher().Until(ctx, event)
}

// Forget removes all listeners for the given event.
func Forget(event any) {
	Dispatcher().Forget(event)
}
