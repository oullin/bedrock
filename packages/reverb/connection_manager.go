package reverb

import "sync"

// ConnectionManager is a thread-safe registry of live WebSocket connections,
// keyed by appID and then socketID.
type ConnectionManager struct {
	mu    sync.RWMutex
	conns map[string]map[string]*Conn // appID → socketID → *Conn
}

// NewConnectionManager constructs an empty ConnectionManager.
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		conns: make(map[string]map[string]*Conn),
	}
}

// Add registers the connection under its appID and socketID. It is safe to
// call Add for a connection that is already registered; the entry is replaced.
func (m *ConnectionManager) Add(conn *Conn) {
	m.mu.Lock()

	defer m.mu.Unlock()

	appID := conn.AppID()

	if m.conns[appID] == nil {
		m.conns[appID] = make(map[string]*Conn)
	}

	m.conns[appID][conn.SocketID()] = conn
}

// Remove deletes the connection identified by appID and socketID from the
// registry. It is a no-op if the connection does not exist.
func (m *ConnectionManager) Remove(appID, socketID string) {
	m.mu.Lock()

	defer m.mu.Unlock()

	if app, ok := m.conns[appID]; ok {
		delete(app, socketID)
	}
}

// Find returns the connection for the given appID and socketID. The boolean
// reports whether the connection was found.
func (m *ConnectionManager) Find(appID, socketID string) (*Conn, bool) {
	m.mu.RLock()

	defer m.mu.RUnlock()

	app, ok := m.conns[appID]

	if !ok {
		return nil, false
	}

	conn, ok := app[socketID]

	return conn, ok
}

// All returns a snapshot of every live connection for the given appID.
func (m *ConnectionManager) All(appID string) []*Conn {
	m.mu.RLock()

	defer m.mu.RUnlock()

	app, ok := m.conns[appID]

	if !ok {
		return nil
	}

	conns := make([]*Conn, 0, len(app))

	for _, conn := range app {
		conns = append(conns, conn)
	}

	return conns
}

// Count returns the number of live connections for the given appID.
func (m *ConnectionManager) Count(appID string) int {
	m.mu.RLock()

	defer m.mu.RUnlock()

	return len(m.conns[appID])
}
