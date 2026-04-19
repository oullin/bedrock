package database

import (
	"context"
	"fmt"
	"sync"

	dbcontract "github.com/bedrock/packages/contracts/database"
	cevents "github.com/bedrock/packages/contracts/events"
	dbevents "github.com/bedrock/packages/database/events"
)

// ConnectorFactory creates a *Connection from configuration.
type ConnectorFactory func(config ConnectionConfig) (*Connection, error)

// Manager creates and manages named database connections. It follows the
// Manager/Driver pattern used by the cache, queue, and redis packages.
type Manager struct {
	mu                sync.RWMutex
	connections       map[string]*Connection
	configs           map[string]ConnectionConfig
	connectorFactory  map[string]ConnectorFactory
	defaultConnection string
	events            cevents.Dispatcher
	reconnector       func(*Connection) (*Connection, error)
}

var _ dbcontract.ConnectionResolver = (*Manager)(nil)

// NewManager creates a new database Manager.
func NewManager() *Manager {
	return &Manager{
		connections:      make(map[string]*Connection),
		configs:          make(map[string]ConnectionConfig),
		connectorFactory: make(map[string]ConnectorFactory),
	}
}

// AddConnection registers a connection configuration.
func (m *Manager) AddConnection(name string, config ConnectionConfig) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.configs[name] = config
}

// Connection returns the named connection, creating it if necessary.
// If no name is given, the default connection is used.
func (m *Manager) Connection(_ context.Context, name ...string) (dbcontract.Connection, error) {
	n := m.defaultConnection

	if len(name) > 0 && name[0] != "" {
		n = name[0]
	}

	m.mu.RLock()
	conn, ok := m.connections[n]
	m.mu.RUnlock()

	if ok {
		return conn, nil
	}

	return m.makeConnection(n)
}

// Reconnect purges the cached connection and creates a new one.
func (m *Manager) Reconnect(name string) (*Connection, error) {
	m.Purge(name)
	conn, err := m.Connection(context.Background(), name)

	if err != nil {
		return nil, err
	}

	return conn.(*Connection), nil
}

// Purge closes and removes a cached connection.
func (m *Manager) Purge(name string) {
	m.mu.Lock()
	conn, ok := m.connections[name]

	if ok {
		_ = conn.Disconnect()
		delete(m.connections, name)
	}

	m.mu.Unlock()
}

// Disconnect closes all connections.
func (m *Manager) Disconnect() {
	m.mu.Lock()

	defer m.mu.Unlock()

	for name, conn := range m.connections {
		_ = conn.Disconnect()
		delete(m.connections, name)
	}
}

// Extend registers a custom connector factory for the given driver.
func (m *Manager) Extend(driver string, factory ConnectorFactory) {
	m.mu.Lock()

	defer m.mu.Unlock()

	m.connectorFactory[driver] = factory
}

// GetDefaultConnection returns the default connection name.
func (m *Manager) GetDefaultConnection() string {
	return m.defaultConnection
}

// SetDefaultConnection sets the default connection name.
func (m *Manager) SetDefaultConnection(name string) {
	m.defaultConnection = name
}

// SetEventDispatcher sets the event dispatcher for all new connections.
func (m *Manager) SetEventDispatcher(events cevents.Dispatcher) {
	m.events = events
}

// SetReconnector registers a function to reconnect a lost connection.
func (m *Manager) SetReconnector(fn func(*Connection) (*Connection, error)) {
	m.reconnector = fn
}

// GetConnections returns all active connections.
func (m *Manager) GetConnections() map[string]*Connection {
	m.mu.RLock()

	defer m.mu.RUnlock()

	conns := make(map[string]*Connection, len(m.connections))

	for k, v := range m.connections {
		conns[k] = v
	}

	return conns
}

// SupportedDrivers returns the list of supported driver names.
func (m *Manager) SupportedDrivers() []string {
	return []string{"mysql", "mariadb", "pgsql", "sqlite"}
}

func (m *Manager) makeConnection(name string) (*Connection, error) {
	m.mu.Lock()

	defer m.mu.Unlock()

	// Double-check after acquiring write lock.
	if conn, ok := m.connections[name]; ok {
		return conn, nil
	}

	config, ok := m.configs[name]

	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrConnectionNotConfigured, name)
	}

	factory, ok := m.connectorFactory[config.Driver]

	if !ok {
		return nil, fmt.Errorf("%w: %q", ErrDriverNotSupported, config.Driver)
	}

	conn, err := factory(config)

	if err != nil {
		return nil, err
	}

	conn.name = name
	conn.SetDriverName(config.Driver)

	if m.events != nil {
		conn.SetEventDispatcher(m.events)
		m.events.Dispatch(context.Background(), dbevents.ConnectionEstablished{
			ConnectionName: name,
			Driver:         config.Driver,
		})
	}

	m.connections[name] = conn

	return conn, nil
}
