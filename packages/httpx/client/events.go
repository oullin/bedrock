package client

import (
	"net/http"
	"sync"
)

// RequestSending is dispatched before an HTTP request is sent.
type RequestSending struct {
	Request *http.Request
}

// ResponseReceived is dispatched after an HTTP response is received.
type ResponseReceived struct {
	Request  *http.Request
	Response *Response
}

// ConnectionFailed is dispatched when a connection attempt fails.
type ConnectionFailed struct {
	Request *http.Request
	Err     error
}

// EventListener is a function that receives dispatched events.
type EventListener func(event any)

// EventDispatcher dispatches events to registered listeners.
type EventDispatcher struct {
	mu        sync.Mutex
	listeners []EventListener
}

// NewEventDispatcher creates a new EventDispatcher.
func NewEventDispatcher() *EventDispatcher {
	return &EventDispatcher{}
}

// Listen registers an event listener.
func (d *EventDispatcher) Listen(listener EventListener) {
	d.mu.Lock()

	defer d.mu.Unlock()

	d.listeners = append(d.listeners, listener)
}

// Dispatch sends an event to all registered listeners.
func (d *EventDispatcher) Dispatch(event any) {
	d.mu.Lock()

	listeners := make([]EventListener, len(d.listeners))
	copy(listeners, d.listeners)

	d.mu.Unlock()

	for _, listener := range listeners {
		listener(event)
	}
}
