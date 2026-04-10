package session

import "errors"

var (
	// ErrAlreadyStarted is returned when Start is called on an already-started session.
	ErrAlreadyStarted = errors.New("session: already started")
	// ErrInvalidID is returned when a session ID fails validation.
	ErrInvalidID = errors.New("session: invalid id")
	// ErrNotStarted is returned when an operation requires a started session.
	ErrNotStarted = errors.New("session: not started")
)
