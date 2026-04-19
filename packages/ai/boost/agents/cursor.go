package agents

import (
	"os"
	"path/filepath"

	"github.com/bedrock/packages/ai/boost/internal/platform"
)

// Cursor implements boost.CodingAgent for the Cursor IDE.
// Mirrors Upstream\Boost\Install\Agents\Cursor.
type Cursor struct {
	BaseAgent
}

// NewCursor constructs a Cursor agent.
func NewCursor(opts ...AgentOptions) *Cursor {
	var o AgentOptions

	if len(opts) > 0 {
		o = opts[0]
	}

	return &Cursor{BaseAgent: NewBaseAgent(o)}
}

func (a *Cursor) Name() string        { return "cursor" }
func (a *Cursor) DisplayName() string { return "Cursor" }

func (a *Cursor) McpInstallationStrategy() platform.McpInstallationStrategy {
	return platform.McpStrategyFile
}

func (a *Cursor) McpConfigPath() string {
	return fallback(a.opts.McpConfigPath, ".cursor/mcp.json")
}

func (a *Cursor) GuidelinesPath() string {
	return fallback(a.opts.GuidelinesPath, "AGENTS.md")
}

func (a *Cursor) SkillsPath() string {
	return fallback(a.opts.SkillsPath, ".cursor/skills")
}

// HttpMcpServerConfig overrides the base to use npx mcp-remote instead of type:http.
func (a *Cursor) HttpMcpServerConfig(url string) map[string]any {
	return map[string]any{
		"command": "npx",
		"args":    []string{"-y", "mcp-remote", url},
	}
}

// DetectOnSystem checks platform-specific installation locations.
func (a *Cursor) DetectOnSystem(p platform.Platform) bool {
	switch p {
	case platform.Darwin:
		return existsOnDisk("/Applications/Cursor.app")
	case platform.Windows:
		return existsOnDisk(filepath.Join(os.Getenv("ProgramFiles"), "Cursor")) ||
			existsOnDisk(filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "Cursor"))
	default:
		return existsOnDisk("/opt/cursor") ||
			commandInPath("cursor") ||
			existsOnDisk("~/.local/bin/cursor")
	}
}

// DetectInProject checks for a .cursor directory.
func (a *Cursor) DetectInProject(basePath string) bool {
	return existsOnDisk(basePath + "/.cursor")
}

// InstallMcp overrides BaseAgent to use the resolved McpConfigPath.
func (a *Cursor) InstallMcp(key, command string, args []string, env map[string]string) (bool, error) {
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
func (a *Cursor) InstallHttpMcp(key, url string) (bool, error) {
	return writeJSONConfigEntry(
		a.McpConfigPath(),
		a.McpConfigKey(),
		key,
		a.HttpMcpServerConfig(url),
		a.DefaultMcpConfig(),
	)
}
