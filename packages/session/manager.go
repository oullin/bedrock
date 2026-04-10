package session

import (
	"context"
	"fmt"
	"sync"
)

// DriverCreator is a factory function that creates a Handler.
type DriverCreator func(config map[string]any) (Handler, error)

// Manager creates and manages named session stores backed by registered drivers.
type Manager struct {
	mu      sync.RWMutex
	stores  map[string]*Store
	drivers map[string]DriverCreator
	config  map[string]map[string]any
	name    string // default session name / cookie name
}

// NewManager creates a Manager. name is the default session name.
func NewManager(name string) *Manager {
	return &Manager{
		stores:  make(map[string]*Store),
		drivers: make(map[string]DriverCreator),
		config:  make(map[string]map[string]any),
		name:    name,
	}
}

// Extend registers a custom driver factory under driverName.
func (m *Manager) Extend(driverName string, creator DriverCreator) *Manager {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.drivers[driverName] = creator

	return m
}

// SetDriverConfig stores the configuration map for a specific driver.
func (m *Manager) SetDriverConfig(driverName string, config map[string]any) *Manager {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.config[driverName] = config

	return m
}

// Driver creates (or returns a cached) Store for the given driver name.
func (m *Manager) Driver(ctx context.Context, driverName string) (*Store, error) {
	m.mu.Lock()

	defer m.mu.Unlock()

	if s, ok := m.stores[driverName]; ok {
		return s, nil
	}

	creator, ok := m.drivers[driverName]

	if !ok {
		return nil, fmt.Errorf("session: unsupported driver %q", driverName)
	}

	cfg := m.config[driverName]

	handler, err := creator(cfg)

	if err != nil {
		return nil, fmt.Errorf("session: create driver %q: %w", driverName, err)
	}

	if err = handler.Open(ctx, "", m.name); err != nil {
		return nil, fmt.Errorf("session: open driver %q: %w", driverName, err)
	}

	store := New(m.name, handler)
	m.stores[driverName] = store

	return store, nil
}
