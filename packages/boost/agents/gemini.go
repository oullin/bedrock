package agents

import "github.com/bedrock/packages/boost/internal/platform"

// Gemini implements boost.CodingAgent for Google Gemini CLI.
// Mirrors Laravel\Boost\Install\Agents\Gemini.
type Gemini struct {
	BaseAgent
}

// NewGemini constructs a Gemini agent.
func NewGemini(opts ...AgentOptions) *Gemini {
	var o AgentOptions

	if len(opts) > 0 {
		o = opts[0]
	}

	return &Gemini{BaseAgent: NewBaseAgent(o)}
}

func (a *Gemini) Name() string        { return "gemini" }
func (a *Gemini) DisplayName() string { return "Gemini" }

func (a *Gemini) McpInstallationStrategy() platform.McpInstallationStrategy {
	return platform.McpStrategyFile
}

func (a *Gemini) McpConfigPath() string {
	return fallback(a.opts.McpConfigPath, ".gemini/settings.json")
}

func (a *Gemini) GuidelinesPath() string {
	return fallback(a.opts.GuidelinesPath, "GEMINI.md")
}

func (a *Gemini) SkillsPath() string {
	return fallback(a.opts.SkillsPath, ".gemini/skills")
}

// McpConfigKey for Gemini uses "mcpServers".
func (a *Gemini) McpConfigKey() string { return "mcpServers" }

// DetectOnSystem checks for the "gemini" binary in PATH.
func (a *Gemini) DetectOnSystem(p platform.Platform) bool {
	switch p {
	case platform.Windows:
		return commandExists("cmd /c where gemini 2>nul")
	default:
		return commandInPath("gemini")
	}
}

// DetectInProject checks for a .gemini directory.
func (a *Gemini) DetectInProject(basePath string) bool {
	return existsOnDisk(basePath + "/.gemini")
}

// InstallMcp overrides BaseAgent to use the resolved McpConfigPath.
func (a *Gemini) InstallMcp(key, command string, args []string, env map[string]string) (bool, error) {
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
func (a *Gemini) InstallHttpMcp(key, url string) (bool, error) {
	return writeJSONConfigEntry(
		a.McpConfigPath(),
		a.McpConfigKey(),
		key,
		a.HttpMcpServerConfig(url),
		a.DefaultMcpConfig(),
	)
}
