package events

import "errors"

var (
	ErrListenerNotFound = errors.New("events: listener not found")
	ErrInvalidListener  = errors.New("events: invalid listener")
)
