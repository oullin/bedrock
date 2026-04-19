package agents

import "github.com/bedrock/packages/ai/boost/internal/platform"

// Copilot implements boost.CodingAgent for GitHub Copilot (VS Code).
// Mirrors Laravel\Boost\Install\Agents\Copilot.
type Copilot struct {
	BaseAgent
}

// NewCopilot constructs a Copilot agent.
func NewCopilot(opts ...AgentOptions) *Copilot {
	var o AgentOptions

	if len(opts) > 0 {
		o = opts[0]
	}

	return &Copilot{BaseAgent: NewBaseAgent(o)}
}

func (a *Copilot) Name() string        { return "copilot" }
func (a *Copilot) DisplayName() string { return "GitHub Copilot" }

func (a *Copilot) McpInstallationStrategy() platform.McpInstallationStrategy {
	return platform.McpStrategyFile
}

func (a *Copilot) McpConfigPath() string {
	return fallback(a.opts.McpConfigPath, ".vscode/mcp.json")
}

func (a *Copilot) GuidelinesPath() string {
	return fallback(a.opts.GuidelinesPath, ".github/copilot-instructions.md")
}

func (a *Copilot) SkillsPath() string {
	return fallback(a.opts.SkillsPath, ".github/copilot-skills")
}

// DetectOnSystem checks for the "code" binary (VS Code) in PATH.
func (a *Copilot) DetectOnSystem(p platform.Platform) bool {
	switch p {
	case platform.Windows:
		return commandExists("cmd /c where code 2>nul")
	default:
		return commandInPath("code")
	}
}

// DetectInProject checks for a .github/copilot-instructions.md file.
func (a *Copilot) DetectInProject(basePath string) bool {
	return existsOnDisk(basePath + "/.github/copilot-instructions.md")
}

// InstallMcp overrides BaseAgent to use the resolved McpConfigPath.
func (a *Copilot) InstallMcp(key, command string, args []string, env map[string]string) (bool, error) {
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
func (a *Copilot) InstallHttpMcp(key, url string) (bool, error) {
	return writeJSONConfigEntry(
		a.McpConfigPath(),
		a.McpConfigKey(),
		key,
		a.HttpMcpServerConfig(url),
		a.DefaultMcpConfig(),
	)
}
