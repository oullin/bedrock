package scout

import (
	"fmt"
	"sync"

	contract "github.com/bedrock/packages/contracts/scout"
)

// DriverFactory creates an Engine from configuration.
type DriverFactory func(config map[string]any) (contract.Engine, error)

// EngineManager creates and manages named search engine instances.
// It follows the Manager/Driver pattern used by cache.Manager and
// database.Manager throughout the Bedrock framework.
type EngineManager struct {
	mu            sync.RWMutex
	engines       map[string]contract.Engine
	drivers       map[string]DriverFactory
	defaultDriver string
}

// NewEngineManager creates an empty EngineManager.
func NewEngineManager() *EngineManager {
	return &EngineManager{
		engines: make(map[string]contract.Engine),
		drivers: make(map[string]DriverFactory),
	}
}

// Register adds a named engine instance directly.
func (m *EngineManager) Register(name string, engine contract.Engine) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.engines[name] = engine
}

// Extend registers a custom driver factory under the given driver name.
func (m *EngineManager) Extend(driver string, factory DriverFactory) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.drivers[driver] = factory
}

// Engine returns a named engine, creating it via the registered factory
// if not yet instantiated. If no name is given, the default driver is used.
func (m *EngineManager) Engine(name ...string) (contract.Engine, error) {
	n := m.defaultDriver

	if len(name) > 0 && name[0] != "" {
		n = name[0]
	}

	if n == "" {
		return nil, ErrEngineNotConfigured
	}

	m.mu.RLock()
	e, ok := m.engines[n]
	m.mu.RUnlock()

	if ok {
		return e, nil
	}

	return nil, fmt.Errorf("%w: %q", ErrEngineNotConfigured, n)
}

// Build creates an engine from a registered driver factory without caching.
func (m *EngineManager) Build(driver string, config map[string]any) (contract.Engine, error) {
	m.mu.RLock()
	factory, ok := m.drivers[driver]
	m.mu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrDriverNotSupported, driver)
	}

	return factory(config)
}

// BuildAndRegister creates an engine from a driver factory and caches it.
func (m *EngineManager) BuildAndRegister(name, driver string, config map[string]any) (contract.Engine, error) {
	engine, err := m.Build(driver, config)

	if err != nil {
		return nil, err
	}

	m.Register(name, engine)

	return engine, nil
}

// GetDefaultDriver returns the default driver name.
func (m *EngineManager) GetDefaultDriver() string {
	return m.defaultDriver
}

// SetDefaultDriver sets the default driver name.
func (m *EngineManager) SetDefaultDriver(name string) {
	m.defaultDriver = name
}

// Driver returns the default driver's engine.
func (m *EngineManager) Driver() (contract.Engine, error) {
	return m.Engine(m.defaultDriver)
}

// Purge removes a cached engine instance.
func (m *EngineManager) Purge(name string) {
	m.mu.Lock()

	defer m.mu.Unlock()

	delete(m.engines, name)
}

// ForgetDriver is an alias for Purge.
func (m *EngineManager) ForgetDriver(name string) {
	m.Purge(name)
}

// GetEngines returns all registered engine names.
func (m *EngineManager) GetEngines() []string {
	m.mu.RLock()

	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.engines))

	for name := range m.engines {
		names = append(names, name)
	}

	return names
}

// GetDrivers returns all registered driver factory names.
func (m *EngineManager) GetDrivers() []string {
	m.mu.RLock()

	defer m.mu.RUnlock()

	names := make([]string, 0, len(m.drivers))

	for name := range m.drivers {
		names = append(names, name)
	}

	return names
}
