package cache

import (
	"fmt"
	"sync"
)

// DriverFactory creates a Store from configuration.
type DriverFactory func(config map[string]any) (Store, error)

// Manager creates and manages named cache stores.
type Manager struct {
	mu       sync.RWMutex
	stores   map[string]Store
	drivers  map[string]DriverFactory
}

// NewManager creates an empty Manager.
func NewManager() *Manager {
	return &Manager{
		stores:  make(map[string]Store),
		drivers: make(map[string]DriverFactory),
	}
}

// Register adds a named store instance directly.
func (m *Manager) Register(name string, store Store) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.stores[name] = store
}

// Extend registers a custom driver factory under the given driver name.
func (m *Manager) Extend(driver string, factory DriverFactory) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.drivers[driver] = factory
}

// Store returns a named store, creating it via the registered factory if not
// yet instantiated.
func (m *Manager) Store(name string) (Store, error) {
	m.mu.RLock()
	s, ok := m.stores[name]
	m.mu.RUnlock()

	if ok {
		return s, nil
	}

	return nil, fmt.Errorf("cache: store %q is not registered", name)
}

// Repository wraps the named store in a Repository.
func (m *Manager) Repository(name string) (*Repository, error) {
	s, err := m.Store(name)
	if err != nil {
		return nil, err
	}

	return NewRepository(s), nil
}
