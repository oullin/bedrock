package telescope

// WatcherFactory is a function that constructs a Watcher and attaches it to
// the given Telescope instance. Callers use this to wire concrete watcher
// types in the service-provider config.
type WatcherFactory func(t *Telescope, options map[string]any) Watcher

// WatcherConfig holds per-watcher configuration: whether it is enabled and
// any driver-specific options.
type WatcherConfig struct {
	Enabled bool
	Options map[string]any
	Factory WatcherFactory
}

// Config is the typed configuration consumed by TelescopeServiceProvider.
type Config struct {
	// Enabled is the master switch. When false the Telescope instance is
	// created but recording is never started.
	Enabled bool

	// Watchers is the ordered list of watcher configurations. Each entry
	// must supply a Factory function.
	Watchers []WatcherConfig
}

// TelescopeServiceProvider registers the Telescope instance with the
// application and boots all configured watchers. It follows the bedrock
// service-provider pattern (Register + Boot).
//
// The app parameter is typed as any to avoid a hard dependency on the
// container package in the telescope package itself. Concrete callers that
// use the bedrock container may pass *container.Container.
type TelescopeServiceProvider struct {
	app      any
	cfg      Config
	opts     []Option
	scope    *Telescope
	watchers []Watcher
}

// NewTelescopeServiceProvider constructs the provider. opts are forwarded to
// New() when creating the Telescope instance.
func NewTelescopeServiceProvider(app any, cfg Config, opts ...Option) *TelescopeServiceProvider {
	return &TelescopeServiceProvider{
		app:  app,
		cfg:  cfg,
		opts: opts,
	}
}

// Register creates the Telescope instance and all configured watchers. It does
// not start recording or call Watcher.Register; that happens in Boot.
func (p *TelescopeServiceProvider) Register() {
	p.scope = New(p.opts...)

	for _, wc := range p.cfg.Watchers {
		if !wc.Enabled || wc.Factory == nil {
			continue
		}

		w := wc.Factory(p.scope, wc.Options)
		p.watchers = append(p.watchers, w)
	}
}

// Boot starts recording (when enabled) and calls Register on every watcher.
func (p *TelescopeServiceProvider) Boot() {
	if p.scope == nil {
		return
	}

	if p.cfg.Enabled {
		p.scope.StartRecording()
	}

	for _, w := range p.watchers {
		w.Register(p.app) //nolint:errcheck
	}
}

// Provides returns the abstract key registered by this provider.
func (p *TelescopeServiceProvider) Provides() []string {
	return []string{"telescope"}
}

// Telescope returns the Telescope instance created by Register. Returns nil
// before Register is called.
func (p *TelescopeServiceProvider) Telescope() *Telescope {
	return p.scope
}

// Watchers returns the slice of registered Watcher instances.
func (p *TelescopeServiceProvider) Watchers() []Watcher {
	return p.watchers
}
