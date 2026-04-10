package queue

import (
	"fmt"
	"sync"
)

// DriverCreator creates a Queue from a configuration map.
type DriverCreator func(config map[string]any) (Queue, error)

// Manager creates and manages named queue connections.
type Manager struct {
	mu       sync.RWMutex
	queues   map[string]Queue
	creators map[string]DriverCreator
	configs  map[string]map[string]any
}

// NewManager creates a Manager.
func NewManager() *Manager {
	return &Manager{
		queues:   make(map[string]Queue),
		creators: make(map[string]DriverCreator),
		configs:  make(map[string]map[string]any),
	}
}

// Register registers a named driver creator (e.g. "sync", "redis").
func (m *Manager) Register(driver string, creator DriverCreator) *Manager {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.creators[driver] = creator

	return m
}

// Extend is an alias for Register (matches Laravel API naming).
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

// Driver returns (or creates) the Queue for the given connection name.
func (m *Manager) Driver(connection string) (Queue, error) {
	m.mu.Lock()

	defer m.mu.Unlock()

	if q, ok := m.queues[connection]; ok {
		return q, nil
	}

	cfg := m.configs[connection]
	driver, _ := cfg["driver"].(string)

	creator, ok := m.creators[driver]

	if !ok {
		return nil, fmt.Errorf("%w: %q (driver: %q)", ErrInvalidDriver, connection, driver)
	}

	q, err := creator(cfg)

	if err != nil {
		return nil, fmt.Errorf("queue: create driver %q: %w", driver, err)
	}

	m.queues[connection] = q

	return q, nil
}
