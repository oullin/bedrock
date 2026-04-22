package broadcasting

import "fmt"

// Manager resolves broadcaster connections by name.
type Manager struct {
	connections map[string]Broadcaster
}

// NewManager creates an empty broadcaster manager.
func NewManager() *Manager {
	return &Manager{connections: make(map[string]Broadcaster)}
}

// Extend registers a broadcaster connection. Use an empty name for default.
func (m *Manager) Extend(name string, broadcaster Broadcaster) *Manager {
	m.connections[name] = broadcaster

	return m
}

// Connection returns a broadcaster connection by name.
func (m *Manager) Connection(name string) (Broadcaster, error) {
	broadcaster, ok := m.connections[name]

	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrConnectionNotFound, name)
	}

	return broadcaster, nil
}
