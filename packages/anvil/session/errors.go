package session

import "errors"

var (
	ErrNotFound       = errors.New("session: not found")
	ErrInvalidID      = errors.New("session: invalid session ID")
	ErrAlreadyStarted = errors.New("session: already started")
)
