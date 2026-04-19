package websockets

import contractsWebSockets "github.com/bedrock/packages/contracts/websockets"

// App holds the runtime configuration for a single WebSockets application.
// It implements contracts/websockets.App.
type App struct {
	id              string
	key             string
	secret          string
	maxConnections  int
	maxMessageSize  int64
	pingInterval    int
	activityTimeout int
	allowedOrigins  []string
	clientEvents    ClientEventsConfig
}

var _ contractsWebSockets.App = (*App)(nil)

// NewApp constructs an App from the given AppConfig, applying zero-value
// defaults before storing the configuration.
func NewApp(cfg AppConfig) *App {
	cfg.defaults()

	return &App{
		id:              cfg.ID,
		key:             cfg.Key,
		secret:          cfg.Secret,
		maxConnections:  cfg.MaxConnections,
		maxMessageSize:  cfg.MaxMessageSize,
		pingInterval:    cfg.PingInterval,
		activityTimeout: cfg.ActivityTimeout,
		allowedOrigins:  cfg.AllowedOrigins,
		clientEvents:    cfg.ClientEvents,
	}
}

// ID returns the unique application identifier.
func (a *App) ID() string { return a.id }

// Key returns the public application key used by clients to connect.
func (a *App) Key() string { return a.key }

// Secret returns the secret used for HMAC-SHA256 signature verification.
func (a *App) Secret() string { return a.secret }

// MaxConnections returns the maximum simultaneous connections allowed.
// Zero means unlimited.
func (a *App) MaxConnections() int { return a.maxConnections }

// MaxMessageSize returns the maximum WebSocket message size in bytes.
// Zero means unlimited.
func (a *App) MaxMessageSize() int64 { return a.maxMessageSize }

// PingInterval returns the interval in seconds between server pings.
func (a *App) PingInterval() int { return a.pingInterval }

// ActivityTimeout returns the seconds of inactivity before a ping is sent.
func (a *App) ActivityTimeout() int { return a.activityTimeout }

// AllowedOrigins returns the list of allowed origin patterns.
// An empty slice means all origins are allowed.
func (a *App) AllowedOrigins() []string { return a.allowedOrigins }

// ClientEventsMode returns the client events policy: "all", "members", or "none".
func (a *App) ClientEventsMode() string { return a.clientEvents.Mode }
