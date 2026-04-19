package reverb

import (
	"github.com/bedrock/packages/container"
	contractsReverb "github.com/bedrock/packages/contracts/reverb"
	"github.com/bedrock/packages/redis"
)

// ReverbServiceProvider registers the Reverb WebSocket server components
// into the application container.
type ReverbServiceProvider struct {
	app    *container.Container
	config Config
}

// NewReverbServiceProvider constructs the provider with the given container and config.
func NewReverbServiceProvider(app *container.Container, config Config) *ReverbServiceProvider {
	return &ReverbServiceProvider{
		app:    app,
		config: config,
	}
}

// Register binds all Reverb singletons into the container:
//
//	"reverb.apps"       → *AppManager
//	"reverb.conns"      → *ConnectionManager
//	"reverb.channels"   → *ChannelManager
//	"reverb.dispatcher" → contractsReverb.Dispatcher (sync or Redis)
//	"reverb.server"     → *Server
//	"reverb.http"       → *HTTPHandler
func (p *ReverbServiceProvider) Register() {
	cfg := p.config

	p.app.Singleton("reverb.apps", func(_ *container.Container) (any, error) {
		return NewAppManager(cfg.Apps), nil
	})

	p.app.Singleton("reverb.conns", func(_ *container.Container) (any, error) {
		return NewConnectionManager(), nil
	})

	p.app.Singleton("reverb.channels", func(c *container.Container) (any, error) {
		raw, err := c.Make("reverb.apps")

		if err != nil {
			return nil, err
		}

		apps := raw.(*AppManager)

		return NewChannelManager(apps), nil
	})

	p.app.Singleton("reverb.dispatcher", func(c *container.Container) (any, error) {
		raw, err := c.Make("reverb.channels")

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
				prefix = "reverb"
			}

			return NewRedisDispatcher(channels, redisMgr, prefix), nil
		}

		return NewSyncDispatcher(channels), nil
	})

	p.app.Singleton("reverb.server", func(c *container.Container) (any, error) {
		appsRaw, err := c.Make("reverb.apps")

		if err != nil {
			return nil, err
		}

		connsRaw, err := c.Make("reverb.conns")

		if err != nil {
			return nil, err
		}

		channelsRaw, err := c.Make("reverb.channels")

		if err != nil {
			return nil, err
		}

		dispatcherRaw, err := c.Make("reverb.dispatcher")

		if err != nil {
			return nil, err
		}

		return NewServer(
			cfg,
			appsRaw.(*AppManager),
			connsRaw.(*ConnectionManager),
			channelsRaw.(*ChannelManager),
			dispatcherRaw.(contractsReverb.Dispatcher),
		), nil
	})

	p.app.Singleton("reverb.http", func(c *container.Container) (any, error) {
		appsRaw, err := c.Make("reverb.apps")

		if err != nil {
			return nil, err
		}

		channelsRaw, err := c.Make("reverb.channels")

		if err != nil {
			return nil, err
		}

		dispatcherRaw, err := c.Make("reverb.dispatcher")

		if err != nil {
			return nil, err
		}

		return NewHTTPHandler(
			appsRaw.(*AppManager),
			channelsRaw.(*ChannelManager),
			dispatcherRaw.(contractsReverb.Dispatcher),
		), nil
	})
}

// Provides returns the abstract keys registered by this provider.
func (p *ReverbServiceProvider) Provides() []string {
	return []string{
		"reverb.apps",
		"reverb.conns",
		"reverb.channels",
		"reverb.dispatcher",
		"reverb.server",
		"reverb.http",
	}
}
