package boost

import "github.com/bedrock/packages/boost/internal/platform"

// AgentsDetector discovers which coding agents are installed on the system or
// configured in a project directory.
// Mirrors Upstream\Boost\Install\AgentsDetector.
type AgentsDetector struct {
	manager *Manager
}

// NewDetector constructs an AgentsDetector backed by the given Manager.
func NewDetector(m *Manager) *AgentsDetector {
	return &AgentsDetector{manager: m}
}

// DiscoverSystemInstalledAgents returns all agents that DetectOnSystem returns
// true for the current platform.
func (d *AgentsDetector) DiscoverSystemInstalledAgents() []CodingAgent {
	p := platform.Current()
	agents := d.manager.GetAgents()
	found := make([]CodingAgent, 0, len(agents))

	for _, agent := range agents {
		if agent.DetectOnSystem(p) {
			found = append(found, agent)
		}
	}

	return found
}

// DiscoverProjectInstalledAgents returns all agents whose DetectInProject
// returns true for basePath.
func (d *AgentsDetector) DiscoverProjectInstalledAgents(basePath string) []CodingAgent {
	agents := d.manager.GetAgents()
	found := make([]CodingAgent, 0, len(agents))

	for _, agent := range agents {
		if agent.DetectInProject(basePath) {
			found = append(found, agent)
		}
	}

	return found
}

// GetAgents returns all registered agents as a flat slice.
func (d *AgentsDetector) GetAgents() []CodingAgent {
	agents := d.manager.GetAgents()
	result := make([]CodingAgent, 0, len(agents))

	for _, agent := range agents {
		result = append(result, agent)
	}

	return result
}
