package agents

import "github.com/bedrock/packages/boost/internal/platform"

// Junie implements boost.CodingAgent for JetBrains Junie.
// Mirrors Laravel\Boost\Install\Agents\Junie.
type Junie struct {
	BaseAgent
}

// NewJunie constructs a Junie agent.
func NewJunie(opts ...AgentOptions) *Junie {
	var o AgentOptions

	if len(opts) > 0 {
		o = opts[0]
	}

	return &Junie{BaseAgent: NewBaseAgent(o)}
}

func (a *Junie) Name() string        { return "junie" }
func (a *Junie) DisplayName() string { return "Junie" }

func (a *Junie) McpInstallationStrategy() platform.McpInstallationStrategy {
	return platform.McpStrategyFile
}

func (a *Junie) McpConfigPath() string {
	return fallback(a.opts.McpConfigPath, ".junie/mcp.json")
}

func (a *Junie) GuidelinesPath() string {
	return fallback(a.opts.GuidelinesPath, ".junie/guidelines.md")
}

func (a *Junie) SkillsPath() string {
	return fallback(a.opts.SkillsPath, ".junie/skills")
}

// DetectOnSystem checks for a .junie directory on the system.
func (a *Junie) DetectOnSystem(p platform.Platform) bool {
	switch p {
	case platform.Darwin:
		return existsOnDisk("~/Library/Application Support/JetBrains")
	default:
		return commandInPath("junie") || existsOnDisk("~/.config/junie")
	}
}

// DetectInProject checks for a .junie directory in the project.
func (a *Junie) DetectInProject(basePath string) bool {
	return existsOnDisk(basePath + "/.junie")
}

// InstallMcp overrides BaseAgent to use the resolved McpConfigPath.
func (a *Junie) InstallMcp(key, command string, args []string, env map[string]string) (bool, error) {
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
func (a *Junie) InstallHttpMcp(key, url string) (bool, error) {
	return writeJSONConfigEntry(
		a.McpConfigPath(),
		a.McpConfigKey(),
		key,
		a.HttpMcpServerConfig(url),
		a.DefaultMcpConfig(),
	)
}
