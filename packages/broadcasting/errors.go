package broadcasting

import "errors"

var (
	// ErrAccessDenied is returned when a user cannot authenticate for a channel.
	ErrAccessDenied = errors.New("broadcasting: access denied")

	// ErrUnknownChannelHandler is returned when a registered channel handler
	// cannot be converted into a callable Go handler.
	ErrUnknownChannelHandler = errors.New("broadcasting: unknown channel handler")

	// ErrBroadcast is returned when a broadcaster backend fails to send an event.
	ErrBroadcast = errors.New("broadcasting: broadcast failed")

	// ErrConnectionNotFound is returned when a factory cannot resolve a named
	// broadcaster connection.
	ErrConnectionNotFound = errors.New("broadcasting: connection not found")
)
