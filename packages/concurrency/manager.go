package concurrency

import (
	"fmt"
	"sync"
)

// DriverCreator creates a Driver from a configuration map.
type DriverCreator func(config map[string]any) (Driver, error)

// Manager creates and manages named concurrency driver instances.
type Manager struct {
	mu                sync.RWMutex
	drivers           map[string]Driver
	creators          map[string]DriverCreator
	configs           map[string]map[string]any
	defaultConnection string
}

// NewManager creates a Manager.
func NewManager() *Manager {
	return &Manager{
		drivers:  make(map[string]Driver),
		creators: make(map[string]DriverCreator),
		configs:  make(map[string]map[string]any),
	}
}

// Register registers a named driver creator (e.g. "goroutine", "sync").
func (m *Manager) Register(driver string, creator DriverCreator) *Manager {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.creators[driver] = creator

	return m
}

// Extend is an alias for Register (matches upstream API naming).
func (m *Manager) Extend(driver string, creator DriverCreator) *Manager {
	return m.Register(driver, creator)
}

// SetConfig stores the configuration for a named connection.
func (m *Manager) SetConfig(connection string, config map[string]any) *Manager {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.configs[connection] = config

	return m
}

// Driver returns (or creates) the Driver for the given connection name.
func (m *Manager) Driver(connection string) (Driver, error) {
	m.mu.Lock()

	defer m.mu.Unlock()

	if d, ok := m.drivers[connection]; ok {
		return d, nil
	}

	cfg := m.configs[connection]
	driver, _ := cfg["driver"].(string)

	creator, ok := m.creators[driver]

	if !ok {
		return nil, fmt.Errorf("%w: %q (driver: %q)", ErrInvalidDriver, connection, driver)
	}

	d, err := creator(cfg)

	if err != nil {
		return nil, fmt.Errorf("concurrency: create driver %q: %w", driver, err)
	}

	m.drivers[connection] = d

	return d, nil
}

// Connection is an alias for Driver (matches upstream naming).
func (m *Manager) Connection(connection string) (Driver, error) {
	return m.Driver(connection)
}

// GetDefaultConnection returns the default connection name.
func (m *Manager) GetDefaultConnection() string {
	m.mu.RLock()

	defer m.mu.RUnlock()

	return m.defaultConnection
}

// SetDefaultConnection sets the default connection name.
func (m *Manager) SetDefaultConnection(connection string) *Manager {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.defaultConnection = connection

	return m
}

// Purge removes a cached driver instance, forcing re-creation on next use.
func (m *Manager) Purge(connection string) {
	m.mu.Lock()

	defer m.mu.Unlock()

	delete(m.drivers, connection)
}

// ForgetDriver removes a registered driver creator.
func (m *Manager) ForgetDriver(driver string) {
	m.mu.Lock()

	defer m.mu.Unlock()

	delete(m.creators, driver)
}
