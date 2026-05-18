package agents

import "github.com/bedrock/packages/ai/boost/internal/platform"

// Aider implements boost.CodingAgent for the Aider CLI assistant.
// Rules-only agent: MCP installation is not yet supported.
type Aider struct {
	BaseAgent
}

// NewAider constructs an Aider agent.
func NewAider(opts ...AgentOptions) *Aider {
	var o AgentOptions

	if len(opts) > 0 {
		o = opts[0]
	}

	return &Aider{BaseAgent: NewBaseAgent(o)}
}

func (a *Aider) Name() string        { return "aider" }
func (a *Aider) DisplayName() string { return "Aider" }

func (a *Aider) McpInstallationStrategy() platform.McpInstallationStrategy {
	return platform.McpStrategyNone
}

func (a *Aider) McpConfigPath() string { return fallback(a.opts.McpConfigPath, "") }

func (a *Aider) GuidelinesPath() string {
	return fallback(a.opts.GuidelinesPath, ".aider.conf.yml")
}

func (a *Aider) SkillsPath() string {
	return fallback(a.opts.SkillsPath, ".aider/skills")
}

// DetectOnSystem checks for the "aider" binary in PATH.
func (a *Aider) DetectOnSystem(p platform.Platform) bool {
	switch p {
	case platform.Windows:
		return commandExists("cmd /c where aider 2>nul")
	default:
		return commandInPath("aider")
	}
}

// DetectInProject checks for an .aider.conf.yml file.
func (a *Aider) DetectInProject(basePath string) bool {
	return existsOnDisk(basePath + "/.aider.conf.yml")
}
