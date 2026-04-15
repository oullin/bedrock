package websockets

import "sync"

// AppManager is a thread-safe registry of all configured WebSockets applications.
// Applications are indexed by both ID and Key for fast lookup.
type AppManager struct {
	mu    sync.RWMutex
	byID  map[string]*App
	byKey map[string]*App
}

// NewAppManager constructs an AppManager pre-populated from the given configs.
func NewAppManager(configs []AppConfig) *AppManager {
	m := &AppManager{
		byID:  make(map[string]*App, len(configs)),
		byKey: make(map[string]*App, len(configs)),
	}

	for _, cfg := range configs {
		app := NewApp(cfg)
		m.byID[app.ID()] = app
		m.byKey[app.Key()] = app
	}

	return m
}

// FindByID returns the App with the given ID, or ErrAppNotFound if none exists.
func (m *AppManager) FindByID(id string) (*App, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	app, ok := m.byID[id]
	if !ok {
		return nil, ErrAppNotFound
	}

	return app, nil
}

// FindByKey returns the App with the given key, or ErrAppNotFound if none exists.
func (m *AppManager) FindByKey(key string) (*App, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	app, ok := m.byKey[key]
	if !ok {
		return nil, ErrAppNotFound
	}

	return app, nil
}

// All returns a snapshot of every registered App.
func (m *AppManager) All() []*App {
	m.mu.RLock()
	defer m.mu.RUnlock()

	apps := make([]*App, 0, len(m.byID))
	for _, app := range m.byID {
		apps = append(apps, app)
	}

	return apps
}
