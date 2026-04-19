package echo

// Callback is the function type invoked when a channel event fires.
// The data parameter carries the event payload.
type Callback func(data any)

// Channel defines the public API for a broadcast channel subscription.
// All implementations must be safe for concurrent use.
type Channel interface {
	// Listen registers a callback for the given event name.
	Listen(event string, callback Callback) Channel

	// StopListening removes a specific callback for the event.
	// If callback is nil, all callbacks for the event are removed.
	StopListening(event string, callback Callback) Channel

	// ListenToAll registers a callback for every event on this channel.
	ListenToAll(callback Callback) Channel

	// StopListeningToAll removes the given callback from "listen-all" listeners.
	// If callback is nil, all "listen-all" callbacks are removed.
	StopListeningToAll(callback Callback) Channel

	// Subscribed registers a callback invoked when the subscription is confirmed.
	Subscribed(callback Callback) Channel

	// Error registers a callback invoked on subscription errors.
	Error(callback Callback) Channel

	// On registers a callback for a raw (un-namespaced) event name.
	On(event string, callback Callback) Channel

	// Leave unsubscribes from this channel.
	Leave()
}

// PrivateChannel extends Channel with private-channel capabilities.
type PrivateChannel interface {
	Channel

	// Whisper sends a client event to other channel members.
	Whisper(event string, data any) PrivateChannel
}

// EncryptedPrivateChannel extends PrivateChannel for end-to-end encrypted channels.
type EncryptedPrivateChannel interface {
	PrivateChannel
}

// PresenceChannel extends Channel with presence-awareness.
type PresenceChannel interface {
	Channel

	// Here registers a callback invoked with the initial member list on join.
	Here(callback func(members []any)) PresenceChannel

	// Joining registers a callback invoked when a new member joins.
	Joining(callback func(member any)) PresenceChannel

	// Leaving registers a callback invoked when a member leaves.
	Leaving(callback func(member any)) PresenceChannel

	// Whisper sends a client event to other channel members.
	Whisper(event string, data any) PresenceChannel
}
