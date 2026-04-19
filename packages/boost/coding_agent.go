package boost

import "github.com/bedrock/packages/boost/internal/platform"

// Platform is a type alias re-exporting internal/platform.Platform so callers
// only need to import the top-level boost package.
type Platform = platform.Platform

// McpInstallationStrategy determines how an agent writes its MCP server config.
// Type alias for internal/platform.McpInstallationStrategy.
type McpInstallationStrategy = platform.McpInstallationStrategy

// McpStrategyFile writes a JSON (or TOML) config file.

// McpStrategyShell runs a shell command to install the MCP server.

// McpStrategyNone indicates no MCP support for this agent.

// CodingAgent is the base interface every IDE coding-assistant agent must satisfy.
// Named CodingAgent (not Agent) to avoid collision with contracts/ai.Agent, which
// models LLM agents. Mirrors the abstract Agent class in laravel/boost.
type CodingAgent interface {
	// Name returns the canonical snake_case key (e.g. "claude_code").
	Name() string

	// DisplayName returns the human-readable name (e.g. "Claude Code").
	DisplayName() string

	// McpConfigPath returns the path where the MCP JSON config is written.
	// Returns "" when the agent has no file-based config.
	McpConfigPath() string

	// McpConfigKey returns the top-level JSON key that holds server entries.
	// Defaults to "mcpServers".
	McpConfigKey() string

	// ShellMcpCommand returns a shell command template for shell-strategy agents.
	ShellMcpCommand() string

	// DefaultMcpConfig returns the skeleton JSON written when the config file
	// does not yet exist. Usually an empty map.
	DefaultMcpConfig() map[string]any

	// Frontmatter reports whether this agent's guidelines file requires YAML
	// frontmatter.
	Frontmatter() bool

	// McpInstallationStrategy returns how the agent installs its MCP server.
	McpInstallationStrategy() McpInstallationStrategy

	// DetectOnSystem reports whether this agent is installed on the host OS.
	DetectOnSystem(p Platform) bool

	// DetectInProject reports whether this agent is configured in basePath.
	DetectInProject(basePath string) bool

	// InstallMcp writes an MCP server entry for key into the agent's config.
	InstallMcp(key, command string, args []string, env map[string]string) (bool, error)

	// InstallHttpMcp writes an HTTP-transport MCP server entry for key.
	InstallHttpMcp(key, url string) (bool, error)

	// HttpMcpServerConfig returns the config payload for an HTTP MCP server.
	HttpMcpServerConfig(url string) map[string]any

	// McpServerConfig returns the config payload for a stdio MCP server.
	McpServerConfig(command string, args []string, env map[string]string) map[string]any

	// UseAbsolutePathForMcp reports whether binary paths must be absolute.
	UseAbsolutePathForMcp() bool

	// GoBinaryPath returns the path to the Go runtime binary.
	// Adapts PHP's getPhpPath().
	GoBinaryPath(forceAbsolute bool) string

	// EntryPointPath returns the path to the application entry-point.
	// Adapts PHP's getArtisanPath(); defaults to "main.go".
	EntryPointPath(forceAbsolute bool) string

	// TransformGuidelines post-processes generated guidelines markdown.
	TransformGuidelines(markdown string) string
}

// SupportsGuidelines is satisfied by agents that accept an AI guidelines file.
// Mirrors Laravel\Boost\Contracts\SupportsGuidelines.
type SupportsGuidelines interface {
	CodingAgent
	GuidelinesPath() string
}

// SupportsMcp is satisfied by agents that can install an MCP server.
// Mirrors Laravel\Boost\Contracts\SupportsMcp.
type SupportsMcp interface {
	CodingAgent
	UseAbsolutePathForMcp() bool
}

// SupportsSkills is satisfied by agents that accept SKILL.md files.
// Mirrors Laravel\Boost\Contracts\SupportsSkills.
type SupportsSkills interface {
	CodingAgent
	SkillsPath() string
}

const (
	McpStrategyFile = platform.McpStrategyFile

	McpStrategyShell = platform.McpStrategyShell

	McpStrategyNone = platform.McpStrategyNone
)
