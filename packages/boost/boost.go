package boost

import (
	"fmt"
	"sync"

	"github.com/bedrock/packages/boost/agents"
)

// defaultAgentFactories maps the canonical keys for the nine built-in agents to
// their factory functions. This mirrors the $agents array in BoostManager.php.

// Manager is the central boost registry that tracks registered coding agents.
// It is registered in the container under the key "boost".
// Mirrors BoostManager in upstream/boost.
type Manager struct {
	mu     sync.RWMutex
	agents map[string]CodingAgent
}

var defaultAgentFactories = map[string]func() CodingAgent{
	"amp":         func() CodingAgent { return agents.NewAmp() },
	"junie":       func() CodingAgent { return agents.NewJunie() },
	"cursor":      func() CodingAgent { return agents.NewCursor() },
	"claude_code": func() CodingAgent { return agents.NewClaudeCode() },
	"codex":       func() CodingAgent { return agents.NewCodex() },
	"copilot":     func() CodingAgent { return agents.NewCopilot() },
	"kiro":        func() CodingAgent { return agents.NewKiro() },
	"opencode":    func() CodingAgent { return agents.NewOpenCode() },
	"gemini":      func() CodingAgent { return agents.NewGemini() },
}

// New returns a Manager pre-loaded with all nine default agents.
func New() *Manager {
	m := &Manager{
		agents: make(map[string]CodingAgent, len(defaultAgentFactories)),
	}

	for key, factory := range defaultAgentFactories {
		m.agents[key] = factory()
	}

	return m
}

// RegisterAgent registers agent under key.
// Returns a wrapped ErrAgentAlreadyRegistered if the key is already taken.
func (m *Manager) RegisterAgent(key string, agent CodingAgent) error {
	m.mu.Lock()

	defer m.mu.Unlock()

	if _, exists := m.agents[key]; exists {
		return fmt.Errorf("%w: %q", ErrAgentAlreadyRegistered, key)
	}

	m.agents[key] = agent

	return nil
}

// GetAgents returns a shallow copy of the registered agents map.
// Callers receive a snapshot; subsequent registrations are not reflected.
func (m *Manager) GetAgents() map[string]CodingAgent {
	m.mu.RLock()

	defer m.mu.RUnlock()

	snapshot := make(map[string]CodingAgent, len(m.agents))

	for k, v := range m.agents {
		snapshot[k] = v
	}

	return snapshot
}
