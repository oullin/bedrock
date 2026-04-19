package websockets

import "errors"

// Pusher WebSocket close codes.
const (
	// CodeInvalidMessage is sent when a message fails protocol validation.
	CodeInvalidMessage = 4200

	// CodeRateLimitExceeded is sent when a connection exceeds the message rate limit.
	CodeRateLimitExceeded = 4202

	// CodeConnectionLimitReached is sent when the per-app connection limit is hit.
	CodeConnectionLimitReached = 4004

	// CodeUnauthorized is sent on failed authentication or invalid origin.
	CodeUnauthorized = 4009

	// CodeConnectionStale is sent when a connection did not respond to a ping.
	CodeConnectionStale = 4201

	// CodeClientEventsDisabled is sent when client events are disabled for the app.
	CodeClientEventsDisabled = 4301
)

var (
	// ErrInvalidMessage is returned when a Pusher message fails protocol validation.
	ErrInvalidMessage = errors.New("websockets: invalid message format")

	// ErrRateLimitExceeded is returned when a connection exceeds the configured rate limit.
	ErrRateLimitExceeded = errors.New("websockets: rate limit exceeded")

	// ErrConnectionLimitReached is returned when the per-app connection limit is reached.
	ErrConnectionLimitReached = errors.New("websockets: connection limit reached")

	// ErrUnauthorized is returned when channel authentication fails.
	ErrUnauthorized = errors.New("websockets: unauthorized")

	// ErrInvalidOrigin is returned when the request origin is not in the allowed list.
	ErrInvalidOrigin = errors.New("websockets: invalid origin")

	// ErrAppNotFound is returned when no app matches the given ID or key.
	ErrAppNotFound = errors.New("websockets: app not found")

	// ErrChannelNotFound is returned when the requested channel does not exist.
	ErrChannelNotFound = errors.New("websockets: channel not found")

	// ErrConnectionNotFound is returned when the requested connection does not exist.
	ErrConnectionNotFound = errors.New("websockets: connection not found")

	// ErrClientEventsDisabled is returned when a client-event is rejected by policy.
	ErrClientEventsDisabled = errors.New("websockets: client events disabled")

	// ErrClientEventsNonMember is returned when a non-member tries to whisper.
	ErrClientEventsNonMember = errors.New("websockets: client event non-member")

	// ErrConnectionStale is returned when a connection has not responded to a ping.
	ErrConnectionStale = errors.New("websockets: connection stale")
)
