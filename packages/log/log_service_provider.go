package log

import (
	"github.com/bedrock/packages/config"
	"github.com/bedrock/packages/container"
)

// LogProviderConfig is the typed configuration consumed by LogServiceProvider.
// It is intentionally narrow — only the keys that the LogManager reads from
// the underlying config.Repository are exposed here so callers don't need to
// know about config.Repository at all.
type LogProviderConfig struct {
	// Default is the default channel name (e.g. "single", "stack").
	Default string

	// Channels maps channel name → channel options. Each channel's options
	// match the keys consumed by ParseChannelConfig (driver, level, path,
	// days, channels, etc.).
	Channels map[string]map[string]any
}

// LogServiceProvider registers the log manager into the container.
// It mirrors Illuminate\Log\LogServiceProvider.
type LogServiceProvider struct {
	app  *container.Container
	cfg  LogProviderConfig
	opts []ManagerOption
}

// NewLogServiceProvider constructs the provider from a typed config.
// opts are forwarded to NewManager (e.g. WithEventDispatcher).
func NewLogServiceProvider(app *container.Container, cfg LogProviderConfig, opts ...ManagerOption) *LogServiceProvider {
	return &LogServiceProvider{app: app, cfg: cfg, opts: opts}
}

// Register binds the log manager as a singleton under "log".
func (p *LogServiceProvider) Register() {
	p.app.Singleton("log", func(_ *container.Container) (any, error) {
		repo := p.cfg.toRepository()

		return NewManager(repo, p.opts...), nil
	})
}

// Provides returns the abstract keys registered by this provider.
func (p *LogServiceProvider) Provides() []string {
	return []string{"log"}
}

// toRepository builds a config.Repository in the shape the LogManager
// expects: nested under the "logging" key.
func (c LogProviderConfig) toRepository() *config.Repository {
	def := c.Default

	if def == "" {
		def = "single"
	}

	channels := c.Channels

	if channels == nil {
		channels = map[string]map[string]any{}
	}

	// config.Repository walks dot notation; nest the values so
	// "logging.default" and "logging.channels.<name>.<key>" resolve.
	logging := map[string]any{
		"default":  def,
		"channels": channels,
	}

	return config.New(map[string]any{"logging": logging})
}
