package agents

import "github.com/bedrock/packages/boost/internal/platform"

// Amp implements boost.CodingAgent for Amp (Anthropic's web-based assistant).
// Mirrors Upstream\Boost\Install\Agents\Amp.
type Amp struct {
	BaseAgent
}

// NewAmp constructs an Amp agent.
func NewAmp(opts ...AgentOptions) *Amp {
	var o AgentOptions
	if len(opts) > 0 {
		o = opts[0]
	}

	return &Amp{BaseAgent: NewBaseAgent(o)}
}

func (a *Amp) Name() string        { return "amp" }
func (a *Amp) DisplayName() string { return "Amp" }

func (a *Amp) McpInstallationStrategy() platform.McpInstallationStrategy {
	return platform.McpStrategyNone
}

func (a *Amp) McpConfigPath() string { return fallback(a.opts.McpConfigPath, "") }

func (a *Amp) GuidelinesPath() string {
	return fallback(a.opts.GuidelinesPath, "AGENTS.md")
}

func (a *Amp) SkillsPath() string {
	return fallback(a.opts.SkillsPath, ".amp/skills")
}

// DetectOnSystem checks for the "amp" binary in PATH.
func (a *Amp) DetectOnSystem(p platform.Platform) bool {
	switch p {
	case platform.Windows:
		return commandExists("cmd /c where amp 2>nul")
	default:
		return commandInPath("amp")
	}
}

// DetectInProject checks for an .amp directory.
func (a *Amp) DetectInProject(basePath string) bool {
	return existsOnDisk(basePath + "/.amp")
}
