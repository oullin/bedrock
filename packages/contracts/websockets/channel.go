package websockets

import "context"

// Channel represents a broadcast channel and its subscriber set.
// All methods must be safe for concurrent use.
type Channel interface {
	// Name returns the channel name.
	Name() string

	// Subscribe adds the connection to this channel.
	// For private and presence channels it validates the HMAC auth token.
	// data is the channel_data JSON string used for presence channels.
	Subscribe(ctx context.Context, conn Connection, auth string, data string) error

	// Unsubscribe removes the connection from this channel.
	Unsubscribe(ctx context.Context, conn Connection)

	// Broadcast sends the event to all subscribers except the one identified by
	// except (the sender's socket ID). Pass nil to broadcast to everyone.
	Broadcast(ctx context.Context, event Event, except *string) error

	// BroadcastToAll sends the event to every subscriber without exclusions.
	BroadcastToAll(ctx context.Context, event Event) error

	// Connections returns a snapshot of all subscribed connections.
	Connections() []Connection

	// HasConnection reports whether the given socket ID is subscribed.
	HasConnection(socketID string) bool
}

// CacheableChannel extends Channel with last-event caching.
// Cache channels deliver the most recent event to new subscribers.
type CacheableChannel interface {
	Channel

	// LastEvent returns the most recently broadcast event, or nil if none.
	LastEvent() *Event

	// CacheEvent stores an event as the last broadcast event.
	CacheEvent(event Event)
}

// PresenceChanneler extends Channel with member presence tracking.
type PresenceChanneler interface {
	Channel

	// Members returns a deduplicated map of user_id → user_info for all subscribers.
	Members() map[string]any

	// MemberCount returns the number of unique users subscribed.
	MemberCount() int

	// MemberIDs returns the unique user IDs of all subscribers.
	MemberIDs() []string
}

// Event is the payload passed to Broadcast calls.
type Event struct {
	Event   string `json:"event"`
	Data    string `json:"data"`
	Channel string `json:"channel,omitempty"`
}
