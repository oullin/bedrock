package websockets

import (
	"context"
	"time"
)

// Connection represents a live WebSocket client connection.
// All methods must be safe for concurrent use.
type Connection interface {
	// SocketID returns the unique identifier for this connection.
	SocketID() string

	// AppID returns the application ID this connection belongs to.
	AppID() string

	// Send writes a raw message to the client.
	Send(ctx context.Context, msg []byte) error

	// Close terminates the WebSocket connection with a Pusher error code and reason.
	Close(ctx context.Context, code int, reason string) error

	// LastSeenAt returns the time the connection last sent any message.
	LastSeenAt() time.Time

	// Touch updates the LastSeenAt timestamp to now.
	Touch()

	// TouchMessage updates the last-message-at timestamp to now.
	TouchMessage()

	// TouchPong updates the last-pong-at timestamp to now.
	TouchPong()
}
