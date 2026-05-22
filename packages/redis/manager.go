package redis

import (
	"fmt"
	"sync"
)

// DriverFactory constructs a Client from a ConnectionConfig. The Manager
// uses factories to lazily build Connections and to support Extend.
type DriverFactory func(cfg ConnectionConfig) (Client, error)

// Manager multiplexes named Connections and is the Go analogue of
// Ref: @bedrock/code-0282
// The zero value is not usable — use NewManager.
type Manager struct {
	mu            sync.RWMutex
	defaultDriver string
	configs       map[string]ConnectionConfig
	connections   map[string]*Connection
	drivers       map[string]DriverFactory
	events        *EventDispatcher
}

// NewManager returns a Manager seeded with the "default" driver (single-
// node via DialSingle), "cluster" (DialCluster) and "sentinel"
// (DialSentinel).
func NewManager(defaultConn string, configs map[string]ConnectionConfig) *Manager {
	m := &Manager{
		defaultDriver: defaultConn,
		configs:       make(map[string]ConnectionConfig, len(configs)),
		connections:   make(map[string]*Connection),
		drivers:       make(map[string]DriverFactory),
		events:        NewEventDispatcher(),
	}

	for k, v := range configs {
		m.configs[k] = v
	}

	m.drivers["default"] = DialSingle
	m.drivers["cluster"] = DialCluster
	m.drivers["sentinel"] = DialSentinel

	return m
}

// Extend registers a custom driver (parity with RedisManager::extend).
func (m *Manager) Extend(driver string, factory DriverFactory) {
	m.mu.Lock()
	m.drivers[driver] = factory
	m.mu.Unlock()
}

// Register inserts a ready-made Connection. Useful for tests and for
// injecting a pre-built fake.
func (m *Manager) Register(name string, conn *Connection) {
	m.mu.Lock()
	m.connections[name] = conn

	if m.events.Enabled() {
		conn.Events().Enable()
	}

	m.mu.Unlock()
}

// AddConfig registers a connection configuration.
func (m *Manager) AddConfig(name string, cfg ConnectionConfig) {
	m.mu.Lock()
	m.configs[name] = cfg
	m.mu.Unlock()
}

// Connection returns the named connection, building it on first use.
// An empty name resolves to the default.
func (m *Manager) Connection(name string) (*Connection, error) {
	if name == "" {
		name = m.defaultDriver
	}

	m.mu.RLock()
	conn, ok := m.connections[name]
	m.mu.RUnlock()

	if ok {
		return conn, nil
	}

	return m.Resolve(name)
}

// Resolve always builds a fresh connection (parity with ::resolve).
func (m *Manager) Resolve(name string) (*Connection, error) {
	m.mu.Lock()

	defer m.mu.Unlock()

	cfg, ok := m.configs[name]

	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrConnectionNotFound, name)
	}

	driver := driverFor(cfg)
	factory, ok := m.drivers[driver]

	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrDriverNotFound, driver)
	}

	client, err := factory(cfg)

	if err != nil {
		return nil, err
	}

	conn := NewConnection(name, client)
	conn.isCluster = cfg.Cluster != nil
	// Inherit manager-wide event enablement.
	if m.events.Enabled() {
		conn.Events().Enable()
		// Fan out dispatcher listeners to the connection.
		conn.Events().Listen(func(e CommandExecuted) { m.events.DispatchExecuted(e) })
	}

	m.connections[name] = conn

	return conn, nil
}

// Connections returns a snapshot of cached connections.
func (m *Manager) Connections() map[string]*Connection {
	m.mu.RLock()

	defer m.mu.RUnlock()

	out := make(map[string]*Connection, len(m.connections))

	for k, v := range m.connections {
		out[k] = v
	}

	return out
}

// Purge removes a cached connection so it will be rebuilt on next access.
func (m *Manager) Purge(name string) {
	m.mu.Lock()

	if c, ok := m.connections[name]; ok {
		_ = c.Close()
		delete(m.connections, name)
	}

	m.mu.Unlock()
}

// SetDriver changes the default connection name.
func (m *Manager) SetDriver(name string) {
	m.mu.Lock()
	m.defaultDriver = name
	m.mu.Unlock()
}

// DefaultConnection returns the default connection name.
func (m *Manager) DefaultConnection() string {
	m.mu.RLock()

	defer m.mu.RUnlock()

	return m.defaultDriver
}

// EnableEvents turns on event dispatching for the manager and every
// currently cached connection. Connections resolved later will inherit.
func (m *Manager) EnableEvents() {
	m.mu.Lock()
	m.events.Enable()

	for _, c := range m.connections {
		c.Events().Enable()
	}

	m.mu.Unlock()
}

// DisableEvents turns off event dispatching.
func (m *Manager) DisableEvents() {
	m.mu.Lock()
	m.events.Disable()

	for _, c := range m.connections {
		c.Events().Disable()
	}

	m.mu.Unlock()
}

// Listen registers a manager-wide listener. Each connection's events are
// forwarded to the manager dispatcher when resolved through Manager.
func (m *Manager) Listen(fn func(CommandExecuted)) {
	m.events.Listen(fn)
	m.EnableEvents()
}

func driverFor(cfg ConnectionConfig) string {
	switch {
	case cfg.Cluster != nil:
		return "cluster"
	case cfg.Sentinel != nil:
		return "sentinel"
	default:
		return "default"
	}
}
