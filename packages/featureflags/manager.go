package featureflags

import (
	"context"
	"fmt"
	"sync"
)

// DriverFactory creates a raw Driver from a config map.
type DriverFactory func(config map[string]any) (Driver, error)

// Manager is the top-level coordinator. It manages named Decorator instances
// backed by registered driver factories. Two built-in factories are registered
// automatically: "array" and "database".
type Manager struct {
	mu            sync.RWMutex
	stores        map[string]*Decorator    // cached Decorator instances keyed by driver name
	factories     map[string]DriverFactory // driver factories keyed by driver name
	defaultDriver string
	dispatcher    EventDispatcher
	scopeResolver func(ctx context.Context) (any, error)
}

// NewManager creates a Manager with the given default driver name and no event
// dispatcher. The "array" and "database" built-in factories are registered.
func NewManager(defaultDriver string) *Manager {
	m := &Manager{
		stores:        make(map[string]*Decorator),
		factories:     make(map[string]DriverFactory),
		defaultDriver: defaultDriver,
	}

	m.registerBuiltins()

	return m
}

// NewManagerWithDispatcher creates a Manager that dispatches featureflags events via
// d. The "array" and "database" built-in factories are registered.
func NewManagerWithDispatcher(defaultDriver string, d EventDispatcher) *Manager {
	m := &Manager{
		stores:        make(map[string]*Decorator),
		factories:     make(map[string]DriverFactory),
		defaultDriver: defaultDriver,
		dispatcher:    d,
	}

	m.registerBuiltins()

	return m
}

// registerBuiltins wires the two built-in driver factories.
func (m *Manager) registerBuiltins() {
	// "array" factory
	m.factories["array"] = func(_ map[string]any) (Driver, error) {
		if m.dispatcher != nil {
			return NewArrayDriverWithDispatcher(m.dispatcher), nil
		}

		return NewArrayDriver(), nil
	}

	// "database" factory
	m.factories["database"] = func(config map[string]any) (Driver, error) {
		if config == nil {
			return nil, fmt.Errorf("featureflags: database driver requires a config map with a \"db\" key")
		}

		dbVal, ok := config["db"]

		if !ok {
			return nil, fmt.Errorf("featureflags: database driver config missing required key \"db\"")
		}

		db, ok := dbVal.(DBExecutor)

		if !ok {
			return nil, fmt.Errorf("featureflags: database driver config key \"db\" must implement DBExecutor")
		}

		table := "features"

		if t, ok := config["table"].(string); ok && t != "" {
			table = t
		}

		if m.dispatcher != nil {
			return NewDatabaseDriverWithDispatcher(db, table, m.dispatcher), nil
		}

		return NewDatabaseDriver(db, table), nil
	}
}

// Store returns the named Decorator, creating it via the registered factory if
// not already instantiated. Passing no name (or an empty string) uses the
// default driver.
func (m *Manager) Store(name ...string) (*Decorator, error) {
	driverName := m.defaultDriver

	if len(name) > 0 && name[0] != "" {
		driverName = name[0]
	}

	// Fast path: cache hit.
	m.mu.RLock()

	if dec, ok := m.stores[driverName]; ok {
		m.mu.RUnlock()

		return dec, nil
	}

	m.mu.RUnlock()

	// Slow path: create the driver and wrap it in a Decorator.
	m.mu.Lock()

	defer m.mu.Unlock()

	// Double-check after acquiring the write lock.
	if dec, ok := m.stores[driverName]; ok {
		return dec, nil
	}

	factory, ok := m.factories[driverName]

	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrDriverNotFound, driverName)
	}

	raw, err := factory(nil)

	if err != nil {
		return nil, err
	}

	var dec *Decorator

	if m.dispatcher != nil {
		dec = NewDecoratorWithDispatcher(raw, m.dispatcher)
	} else {
		dec = NewDecorator(raw)
	}

	m.stores[driverName] = dec

	return dec, nil
}

// Driver is an alias for Store (Upstream API parity).
func (m *Manager) Driver(name ...string) (*Decorator, error) {
	return m.Store(name...)
}

// Extend registers a custom driver factory under the given driver name.
// Calling Extend for a name that already exists overwrites the previous factory.
func (m *Manager) Extend(driver string, factory DriverFactory) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.factories[driver] = factory
}

// ResolveScopeUsing sets the default scope resolver used by context-based
// operations.
func (m *Manager) ResolveScopeUsing(resolver func(ctx context.Context) (any, error)) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.scopeResolver = resolver
}

// SerializeScope is a convenience wrapper around the package-level
// SerializeScope function.
func (m *Manager) SerializeScope(scope any) (string, error) {
	return SerializeScope(scope)
}

// FlushCache calls FlushCache on every cached Decorator.
func (m *Manager) FlushCache() {
	m.mu.RLock()
	// Snapshot the slice of decorators while holding the read lock so we can
	// call FlushCache (which acquires its own lock) outside ours.
	decs := make([]*Decorator, 0, len(m.stores))

	for _, dec := range m.stores {
		decs = append(decs, dec)
	}

	m.mu.RUnlock()

	for _, dec := range decs {
		dec.FlushCache()
	}
}

// SetDefaultDriver changes the default driver name.
func (m *Manager) SetDefaultDriver(name string) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.defaultDriver = name
}

// GetDefaultDriver returns the current default driver name.
func (m *Manager) GetDefaultDriver() string {
	m.mu.RLock()

	defer m.mu.RUnlock()

	return m.defaultDriver
}

// DefaultDecorator returns the Decorator for the current default driver.
// Used internally by ScopedFeatureInteraction.
func (m *Manager) DefaultDecorator() (*Decorator, error) {
	return m.Store()
}
