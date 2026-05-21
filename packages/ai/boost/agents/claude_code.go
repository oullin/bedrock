package agents

import "github.com/bedrock/packages/ai/boost/internal/platform"

// ClaudeCode implements boost.CodingAgent for the Claude Code IDE.
// It also satisfies SupportsGuidelines, SupportsMcp, and SupportsSkills.
type ClaudeCode struct {
	BaseAgent
}

// NewClaudeCode constructs a ClaudeCode agent. Pass AgentOptions to override
// the default paths).
func NewClaudeCode(opts ...AgentOptions) *ClaudeCode {
	var o AgentOptions

	if len(opts) > 0 {
		o = opts[0]
	}

	return &ClaudeCode{BaseAgent: NewBaseAgent(o)}
}

func (a *ClaudeCode) Name() string        { return "claude_code" }
func (a *ClaudeCode) DisplayName() string { return "Claude Code" }

// McpInstallationStrategy uses file-based installation (.mcp.json).
func (a *ClaudeCode) McpInstallationStrategy() platform.McpInstallationStrategy {
	return platform.McpStrategyFile
}

func (a *ClaudeCode) McpConfigPath() string {
	return fallback(a.opts.McpConfigPath, ".mcp.json")
}

func (a *ClaudeCode) GuidelinesPath() string {
	return fallback(a.opts.GuidelinesPath, "CLAUDE.md")
}

func (a *ClaudeCode) SkillsPath() string {
	return fallback(a.opts.SkillsPath, ".claude/skills")
}

// DetectOnSystem checks for the "claude" binary in PATH.
func (a *ClaudeCode) DetectOnSystem(p platform.Platform) bool {
	switch p {
	case platform.Windows:
		return commandExists("cmd /c where claude 2>nul")
	default:
		return commandInPath("claude")
	}
}

// DetectInProject checks for a .claude directory or CLAUDE.md file.
func (a *ClaudeCode) DetectInProject(basePath string) bool {
	return existsOnDisk(basePath+"/.claude") ||
		existsOnDisk(basePath+"/CLAUDE.md")
}

// InstallMcp overrides BaseAgent so it uses the resolved McpConfigPath.
func (a *ClaudeCode) InstallMcp(key, command string, args []string, env map[string]string) (bool, error) {
	cmd, normalArgs := normalizeCommand(command, args)

	return writeJSONConfigEntry(
		a.McpConfigPath(),
		a.McpConfigKey(),
		key,
		a.McpServerConfig(cmd, normalArgs, env),
		a.DefaultMcpConfig(),
	)
}

// InstallHttpMcp overrides BaseAgent so it uses the resolved McpConfigPath.
func (a *ClaudeCode) InstallHttpMcp(key, url string) (bool, error) {
	return writeJSONConfigEntry(
		a.McpConfigPath(),
		a.McpConfigKey(),
		key,
		a.HttpMcpServerConfig(url),
		a.DefaultMcpConfig(),
	)
}
