package agents

import "github.com/bedrock/packages/ai/boost/internal/platform"

// Codex implements boost.CodingAgent for OpenAI Codex CLI.
type Codex struct {
	BaseAgent
}

// NewCodex constructs a Codex agent.
func NewCodex(opts ...AgentOptions) *Codex {
	var o AgentOptions

	if len(opts) > 0 {
		o = opts[0]
	}

	return &Codex{BaseAgent: NewBaseAgent(o)}
}

func (a *Codex) Name() string        { return "codex" }
func (a *Codex) DisplayName() string { return "Codex" }

func (a *Codex) McpInstallationStrategy() platform.McpInstallationStrategy {
	return platform.McpStrategyFile
}

func (a *Codex) McpConfigPath() string {
	return fallback(a.opts.McpConfigPath, "codex.json")
}

func (a *Codex) GuidelinesPath() string {
	return fallback(a.opts.GuidelinesPath, "AGENTS.md")
}

func (a *Codex) SkillsPath() string {
	return fallback(a.opts.SkillsPath, ".codex/skills")
}

// McpConfigKey for Codex uses "mcpServers".
func (a *Codex) McpConfigKey() string { return "mcpServers" }

// DetectOnSystem checks for the "codex" binary in PATH.
func (a *Codex) DetectOnSystem(p platform.Platform) bool {
	switch p {
	case platform.Windows:
		return commandExists("cmd /c where codex 2>nul")
	default:
		return commandInPath("codex")
	}
}

// DetectInProject checks for a codex.json file.
func (a *Codex) DetectInProject(basePath string) bool {
	return existsOnDisk(basePath + "/codex.json")
}

// InstallMcp overrides BaseAgent to use the resolved McpConfigPath.
func (a *Codex) InstallMcp(key, command string, args []string, env map[string]string) (bool, error) {
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
func (a *Codex) InstallHttpMcp(key, url string) (bool, error) {
	return writeJSONConfigEntry(
		a.McpConfigPath(),
		a.McpConfigKey(),
		key,
		a.HttpMcpServerConfig(url),
		a.DefaultMcpConfig(),
	)
}
