package agents

import "github.com/bedrock/packages/boost/internal/platform"

// OpenCode implements boost.CodingAgent for the OpenCode CLI agent.
// Mirrors Upstream\Boost\Install\Agents\OpenCode.
type OpenCode struct {
	BaseAgent
}

// NewOpenCode constructs an OpenCode agent.
func NewOpenCode(opts ...AgentOptions) *OpenCode {
	var o AgentOptions

	if len(opts) > 0 {
		o = opts[0]
	}

	return &OpenCode{BaseAgent: NewBaseAgent(o)}
}

func (a *OpenCode) Name() string        { return "opencode" }
func (a *OpenCode) DisplayName() string { return "OpenCode" }

func (a *OpenCode) McpInstallationStrategy() platform.McpInstallationStrategy {
	return platform.McpStrategyFile
}

func (a *OpenCode) McpConfigPath() string {
	return fallback(a.opts.McpConfigPath, "opencode.json")
}

func (a *OpenCode) GuidelinesPath() string {
	return fallback(a.opts.GuidelinesPath, "AGENTS.md")
}

func (a *OpenCode) SkillsPath() string {
	return fallback(a.opts.SkillsPath, ".opencode/skills")
}

// DetectOnSystem checks for the "opencode" binary in PATH.
func (a *OpenCode) DetectOnSystem(p platform.Platform) bool {
	switch p {
	case platform.Windows:
		return commandExists("cmd /c where opencode 2>nul")
	default:
		return commandInPath("opencode")
	}
}

// DetectInProject checks for an opencode.json file.
func (a *OpenCode) DetectInProject(basePath string) bool {
	return existsOnDisk(basePath + "/opencode.json")
}

// InstallMcp overrides BaseAgent to use the resolved McpConfigPath.
func (a *OpenCode) InstallMcp(key, command string, args []string, env map[string]string) (bool, error) {
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
func (a *OpenCode) InstallHttpMcp(key, url string) (bool, error) {
	return writeJSONConfigEntry(
		a.McpConfigPath(),
		a.McpConfigKey(),
		key,
		a.HttpMcpServerConfig(url),
		a.DefaultMcpConfig(),
	)
}
