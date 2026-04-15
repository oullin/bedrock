package reverb

import (
	"sync"

	contractsReverb "github.com/bedrock/packages/contracts/reverb"
)

// ChannelManager is a thread-safe registry of channels grouped by application ID.
type ChannelManager struct {
	mu       sync.RWMutex
	channels map[string]map[string]contractsReverb.Channel // appID → channelName → Channel
	apps     *AppManager
}

// NewChannelManager constructs a ChannelManager backed by the given AppManager.
func NewChannelManager(apps *AppManager) *ChannelManager {
	return &ChannelManager{
		channels: make(map[string]map[string]contractsReverb.Channel),
		apps:     apps,
	}
}

// GetOrCreate returns the existing channel for (appID, channelName), or creates
// a new one via NewChannel if it does not exist yet.
// It returns ErrAppNotFound when no app with appID is registered.
func (m *ChannelManager) GetOrCreate(appID, channelName string) (contractsReverb.Channel, error) {
	// Fast path: channel already exists.
	m.mu.RLock()
	if appChans, ok := m.channels[appID]; ok {
		if ch, ok := appChans[channelName]; ok {
			m.mu.RUnlock()
			return ch, nil
		}
	}
	m.mu.RUnlock()

	// Resolve the app to build the channel.
	app, err := m.apps.FindByID(appID)
	if err != nil {
		return nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check after acquiring the write lock.
	if appChans, ok := m.channels[appID]; ok {
		if ch, ok := appChans[channelName]; ok {
			return ch, nil
		}
	}

	ch := NewChannel(channelName, app)

	if _, ok := m.channels[appID]; !ok {
		m.channels[appID] = make(map[string]contractsReverb.Channel)
	}

	m.channels[appID][channelName] = ch

	return ch, nil
}

// Get returns the channel for (appID, channelName) and reports whether it exists.
func (m *ChannelManager) Get(appID, channelName string) (contractsReverb.Channel, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	appChans, ok := m.channels[appID]
	if !ok {
		return nil, false
	}

	ch, ok := appChans[channelName]

	return ch, ok
}

// Remove deletes the channel identified by (appID, channelName).
func (m *ChannelManager) Remove(appID, channelName string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if appChans, ok := m.channels[appID]; ok {
		delete(appChans, channelName)
	}
}

// All returns a snapshot of every channel registered for appID.
func (m *ChannelManager) All(appID string) []contractsReverb.Channel {
	m.mu.RLock()
	defer m.mu.RUnlock()

	appChans, ok := m.channels[appID]
	if !ok {
		return nil
	}

	out := make([]contractsReverb.Channel, 0, len(appChans))
	for _, ch := range appChans {
		out = append(out, ch)
	}

	return out
}

// CleanupEmpty removes every channel under appID that has no subscribers.
func (m *ChannelManager) CleanupEmpty(appID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	appChans, ok := m.channels[appID]
	if !ok {
		return
	}

	for name, ch := range appChans {
		if len(ch.Connections()) == 0 {
			delete(appChans, name)
		}
	}
}
