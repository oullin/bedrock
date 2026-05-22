package debugbar

// WatcherFactory is a function that constructs a Watcher and attaches it to
// the given DebugBar instance. Callers use this to wire concrete watcher
// types in the service-provider config.
type WatcherFactory func(t *DebugBar, options map[string]any) Watcher

// WatcherConfig holds per-watcher configuration: whether it is enabled and
// any driver-specific options.
type WatcherConfig struct {
	Enabled bool
	Options map[string]any
	Factory WatcherFactory
}

// Config is the typed configuration consumed by DebugBarServiceProvider.
type Config struct {
	// Enabled is the master switch. When false the DebugBar instance is
	// created but recording is never started.
	Enabled bool

	// Watchers is the ordered list of watcher configurations. Each entry
	// must supply a Factory function.
	Watchers []WatcherConfig
}

// DebugBarServiceProvider registers the DebugBar instance with the
// application and boots all configured watchers. It follows the bedrock
// service-provider pattern (Register + Boot).
//
// The app parameter is typed as any to avoid a hard dependency on the
// container package in the debugbar package itself. Concrete callers that
// use the bedrock container may pass *container.Container.
type DebugBarServiceProvider struct {
	app      any
	cfg      Config
	opts     []Option
	scope    *DebugBar
	watchers []Watcher
}

// NewDebugBarServiceProvider constructs the provider. opts are forwarded to
// New() when creating the DebugBar instance.
func NewDebugBarServiceProvider(app any, cfg Config, opts ...Option) *DebugBarServiceProvider {
	return &DebugBarServiceProvider{
		app:  app,
		cfg:  cfg,
		opts: opts,
	}
}

// Register creates the DebugBar instance and all configured watchers. It does
// not start recording or call Watcher.Register; that happens in Boot.
func (p *DebugBarServiceProvider) Register() {
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
func (p *DebugBarServiceProvider) Boot() {
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
func (p *DebugBarServiceProvider) Provides() []string {
	return []string{"debugbar"}
}

// DebugBar returns the DebugBar instance created by Register. Returns nil
// before Register is called.
func (p *DebugBarServiceProvider) DebugBar() *DebugBar {
	return p.scope
}

// Watchers returns the slice of registered Watcher instances.
func (p *DebugBarServiceProvider) Watchers() []Watcher {
	return p.watchers
}
