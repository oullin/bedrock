package reverb

// ClientEventsConfig controls the client-events (whisper) policy for an app.
type ClientEventsConfig struct {
	// Mode is "all", "members", or "none".
	//   "all"     – any subscriber may whisper on private/presence channels.
	//   "members" – only presence-channel members may whisper.
	//   "none"    – client events are disabled (returns error code 4301).
	Mode string
}

// AppConfig holds the per-application configuration.
type AppConfig struct {
	// ID is the unique application identifier.
	ID string

	// Key is the public key clients use when connecting.
	Key string

	// Secret is used for HMAC-SHA256 signature verification.
	Secret string

	// MaxConnections is the maximum simultaneous WebSocket connections.
	// Zero means unlimited.
	MaxConnections int

	// MaxMessageSize is the maximum WebSocket message size in bytes.
	// Zero means unlimited (nhooyr.io/websocket default applies).
	MaxMessageSize int64

	// PingInterval is the interval in seconds between server pings.
	// Defaults to 60.
	PingInterval int

	// ActivityTimeout is the seconds of inactivity before a ping is sent.
	// Defaults to 30.
	ActivityTimeout int

	// AllowedOrigins is the list of permitted request origin patterns.
	// An empty slice permits all origins.
	// Wildcards are supported: "*.example.com".
	AllowedOrigins []string

	// ClientEvents controls the client-events whisper policy.
	ClientEvents ClientEventsConfig
}

// RedisConfig holds the optional Redis pub/sub settings for horizontal scaling.
type RedisConfig struct {
	// Connection is the name of the Redis connection in the redis.Manager.
	Connection string

	// Prefix is prepended to all Redis pub/sub channel names.
	// Defaults to "reverb".
	Prefix string
}

// Config is the top-level Reverb server configuration.
type Config struct {
	// Apps is the list of configured applications.
	Apps []AppConfig

	// Redis holds the optional Redis scaling configuration.
	// When nil the synchronous (in-process) dispatcher is used.
	Redis *RedisConfig

	// Host is the address the WebSocket server listens on.
	Host string

	// Port is the port the WebSocket server listens on.
	Port int

	// MaxRequestSize is the maximum HTTP request body size in bytes.
	// Applies to the HTTP trigger API. Zero means 10 000.
	MaxRequestSize int64
}

// DefaultConfig returns a Config populated with sensible defaults that mirror
// the Laravel Reverb defaults.
func DefaultConfig() Config {
	return Config{
		Host:           "0.0.0.0",
		Port:           8080,
		MaxRequestSize: 10_000,
	}
}

// defaults applies zero-value defaults to an AppConfig.
func (c *AppConfig) defaults() {
	if c.PingInterval == 0 {
		c.PingInterval = 60
	}
	if c.ActivityTimeout == 0 {
		c.ActivityTimeout = 30
	}
	if c.ClientEvents.Mode == "" {
		c.ClientEvents.Mode = "none"
	}
}
