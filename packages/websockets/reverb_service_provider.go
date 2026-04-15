package websockets

import (
	contractsWebSockets "github.com/bedrock/packages/contracts/websockets"
	"github.com/bedrock/packages/container"
	"github.com/bedrock/packages/redis"
)

// WebSocketsServiceProvider registers the WebSockets WebSocket server components
// into the application container.
type WebSocketsServiceProvider struct {
	app    *container.Container
	config Config
}

// NewWebSocketsServiceProvider constructs the provider with the given container and config.
func NewWebSocketsServiceProvider(app *container.Container, config Config) *WebSocketsServiceProvider {
	return &WebSocketsServiceProvider{
		app:    app,
		config: config,
	}
}

// Register binds all WebSockets singletons into the container:
//
//	"websockets.apps"       → *AppManager
//	"websockets.conns"      → *ConnectionManager
//	"websockets.channels"   → *ChannelManager
//	"websockets.dispatcher" → contractsWebSockets.Dispatcher (sync or Redis)
//	"websockets.server"     → *Server
//	"websockets.http"       → *HTTPHandler
func (p *WebSocketsServiceProvider) Register() {
	cfg := p.config

	p.app.Singleton("websockets.apps", func(_ *container.Container) (any, error) {
		return NewAppManager(cfg.Apps), nil
	})

	p.app.Singleton("websockets.conns", func(_ *container.Container) (any, error) {
		return NewConnectionManager(), nil
	})

	p.app.Singleton("websockets.channels", func(c *container.Container) (any, error) {
		raw, err := c.Make("websockets.apps")
		if err != nil {
			return nil, err
		}

		apps := raw.(*AppManager)

		return NewChannelManager(apps), nil
	})

	p.app.Singleton("websockets.dispatcher", func(c *container.Container) (any, error) {
		raw, err := c.Make("websockets.channels")
		if err != nil {
			return nil, err
		}

		channels := raw.(*ChannelManager)

		if cfg.Redis != nil {
			redisMgrRaw, err := c.Make("redis")
			if err != nil {
				return nil, err
			}

			redisMgr := redisMgrRaw.(*redis.Manager)
			prefix := cfg.Redis.Prefix
			if prefix == "" {
				prefix = "websockets"
			}

			return NewRedisDispatcher(channels, redisMgr, prefix), nil
		}

		return NewSyncDispatcher(channels), nil
	})

	p.app.Singleton("websockets.server", func(c *container.Container) (any, error) {
		appsRaw, err := c.Make("websockets.apps")
		if err != nil {
			return nil, err
		}

		connsRaw, err := c.Make("websockets.conns")
		if err != nil {
			return nil, err
		}

		channelsRaw, err := c.Make("websockets.channels")
		if err != nil {
			return nil, err
		}

		dispatcherRaw, err := c.Make("websockets.dispatcher")
		if err != nil {
			return nil, err
		}

		return NewServer(
			cfg,
			appsRaw.(*AppManager),
			connsRaw.(*ConnectionManager),
			channelsRaw.(*ChannelManager),
			dispatcherRaw.(contractsWebSockets.Dispatcher),
		), nil
	})

	p.app.Singleton("websockets.http", func(c *container.Container) (any, error) {
		appsRaw, err := c.Make("websockets.apps")
		if err != nil {
			return nil, err
		}

		channelsRaw, err := c.Make("websockets.channels")
		if err != nil {
			return nil, err
		}

		dispatcherRaw, err := c.Make("websockets.dispatcher")
		if err != nil {
			return nil, err
		}

		return NewHTTPHandler(
			appsRaw.(*AppManager),
			channelsRaw.(*ChannelManager),
			dispatcherRaw.(contractsWebSockets.Dispatcher),
		), nil
	})
}

// Provides returns the abstract keys registered by this provider.
func (p *WebSocketsServiceProvider) Provides() []string {
	return []string{
		"websockets.apps",
		"websockets.conns",
		"websockets.channels",
		"websockets.dispatcher",
		"websockets.server",
		"websockets.http",
	}
}
