package events

import "errors"

var (
	// ErrInvalidEvent is returned when an event cannot be resolved to a name.
	ErrInvalidEvent = errors.New("events: invalid event type")
)
