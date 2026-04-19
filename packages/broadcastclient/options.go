package broadcastclient

// Options configures an BroadcastClient instance.
//
// Broadcaster selects the transport backend:
//   - "pusher"    — PusherConnector
//   - "websockets"    — PusherConnector (WebSockets variant)
//   - "socket.io" — SocketIOConnector
//   - "null"      — NullConnector (no-op, useful for tests)
//
// When Connector is non-nil it is used directly and Broadcaster is ignored.
// This is the Go equivalent of passing a custom connector constructor in TS.
type Options struct {
	// Broadcaster names the built-in transport. Ignored when Connector is set.
	Broadcaster string

	// Connector is a pre-built connector instance that takes precedence over
	// Broadcaster when non-nil.
	Connector Connector

	// Namespace is prepended to event names by EventFormatter.
	// Example: "App.Events" (the default in Upstream BroadcastClient).
	// Set to empty string "" to disable namespacing, equivalent to
	// namespace: false in the TypeScript library.
	Namespace string

	// Host is the broadcaster host address.
	Host string

	// Key is the Pusher application key.
	Key string

	// Auth holds authentication configuration for private/presence channels.
	Auth AuthOptions
}

// AuthOptions carries authentication configuration for channel subscriptions.
type AuthOptions struct {
	// Endpoint is the URL used to authenticate channel subscriptions.
	// Defaults to "/broadcasting/auth".
	Endpoint string

	// Headers are added to every authentication request.
	Headers map[string]string
}
