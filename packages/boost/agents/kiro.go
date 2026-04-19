package agents

import "github.com/bedrock/packages/boost/internal/platform"

// Kiro implements boost.CodingAgent for AWS Kiro IDE.
// Mirrors Upstream\Boost\Install\Agents\Kiro.
type Kiro struct {
	BaseAgent
}

// NewKiro constructs a Kiro agent.
func NewKiro(opts ...AgentOptions) *Kiro {
	var o AgentOptions

	if len(opts) > 0 {
		o = opts[0]
	}

	return &Kiro{BaseAgent: NewBaseAgent(o)}
}

func (a *Kiro) Name() string        { return "kiro" }
func (a *Kiro) DisplayName() string { return "Kiro" }

func (a *Kiro) McpInstallationStrategy() platform.McpInstallationStrategy {
	return platform.McpStrategyFile
}

func (a *Kiro) McpConfigPath() string {
	return fallback(a.opts.McpConfigPath, ".kiro/mcp.json")
}

func (a *Kiro) GuidelinesPath() string {
	return fallback(a.opts.GuidelinesPath, ".kiro/steering/guidelines.md")
}

func (a *Kiro) SkillsPath() string {
	return fallback(a.opts.SkillsPath, ".kiro/skills")
}

// DetectOnSystem checks for the "kiro" binary or app.
func (a *Kiro) DetectOnSystem(p platform.Platform) bool {
	switch p {
	case platform.Darwin:
		return existsOnDisk("/Applications/Kiro.app") || commandInPath("kiro")
	case platform.Windows:
		return commandExists("cmd /c where kiro 2>nul")
	default:
		return commandInPath("kiro")
	}
}

// DetectInProject checks for a .kiro directory in the project.
func (a *Kiro) DetectInProject(basePath string) bool {
	return existsOnDisk(basePath + "/.kiro")
}

// InstallMcp overrides BaseAgent to use the resolved McpConfigPath.
func (a *Kiro) InstallMcp(key, command string, args []string, env map[string]string) (bool, error) {
	cmd, normalArgs := normalizeCommand(command, args)

	return writeJSONConfigEntry(
		a.McpConfigPath(),
		a.McpConfigKey(),
		key,
		a.McpServerConfig(cmd, normalArgs, env),
		a.DefaultMcpConfig(),
	)
}

// InstallHttpMcp overrides BaseAgent to use the resolved McpConfigPath.
func (a *Kiro) InstallHttpMcp(key, url string) (bool, error) {
	return writeJSONConfigEntry(
		a.McpConfigPath(),
		a.McpConfigKey(),
		key,
		a.HttpMcpServerConfig(url),
		a.DefaultMcpConfig(),
	)
}
