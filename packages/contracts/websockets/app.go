package websockets

// App holds the configuration for a single WebSockets application.
// Each application has its own credentials and resource limits.
type App interface {
	// ID returns the unique application identifier.
	ID() string

	// Key returns the public application key used by clients to connect.
	Key() string

	// Secret returns the secret used for HMAC-SHA256 signature verification.
	Secret() string

	// MaxConnections returns the maximum simultaneous connections allowed.
	// Zero means unlimited.
	MaxConnections() int

	// MaxMessageSize returns the maximum WebSocket message size in bytes.
	// Zero means unlimited.
	MaxMessageSize() int64

	// PingInterval returns the interval in seconds between server pings.
	PingInterval() int

	// ActivityTimeout returns the seconds of inactivity before a ping is sent.
	ActivityTimeout() int

	// AllowedOrigins returns the list of allowed origin patterns.
	// An empty slice means all origins are allowed.
	// Supports wildcard patterns such as "*.example.com".
	AllowedOrigins() []string

	// ClientEventsMode returns the client events policy: "all", "members", or "none".
	ClientEventsMode() string
}
