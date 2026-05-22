package agents

import (
	"path/filepath"

	"github.com/bedrock/packages/ai/boost/internal/platform"
)

// Windsurf implements boost.CodingAgent for the Windsurf editor (Codeium).
// Net-new agent. Rules-only today —
// MCP installation will be wired in when Windsurf gains MCP support upstream.
type Windsurf struct {
	BaseAgent
}

// NewWindsurf constructs a Windsurf agent.
func NewWindsurf(opts ...AgentOptions) *Windsurf {
	var o AgentOptions

	if len(opts) > 0 {
		o = opts[0]
	}

	return &Windsurf{BaseAgent: NewBaseAgent(o)}
}

func (a *Windsurf) Name() string        { return "windsurf" }
func (a *Windsurf) DisplayName() string { return "Windsurf" }

func (a *Windsurf) McpInstallationStrategy() platform.McpInstallationStrategy {
	return platform.McpStrategyNone
}

func (a *Windsurf) McpConfigPath() string { return fallback(a.opts.McpConfigPath, "") }

func (a *Windsurf) GuidelinesPath() string {
	return fallback(a.opts.GuidelinesPath, ".windsurfrules")
}

func (a *Windsurf) SkillsPath() string {
	return fallback(a.opts.SkillsPath, ".windsurf/skills")
}

// DetectOnSystem checks for the "windsurf" binary in PATH.
func (a *Windsurf) DetectOnSystem(p platform.Platform) bool {
	switch p {
	case platform.Windows:
		return commandExists("cmd /c where windsurf 2>nul")
	default:
		return commandInPath("windsurf")
	}
}

// DetectInProject checks for a .windsurfrules file or a .windsurf directory.
func (a *Windsurf) DetectInProject(basePath string) bool {
	return existsOnDisk(filepath.Join(basePath, ".windsurfrules")) ||
		existsOnDisk(filepath.Join(basePath, ".windsurf"))
}
