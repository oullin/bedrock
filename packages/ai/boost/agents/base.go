// Package agents provides the nine built-in IDE coding-agent implementations.
// Each agent embeds BaseAgent for shared behaviour and overrides only what differs.
package agents

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/bedrock/packages/ai/boost/internal/boosterr"
	"github.com/bedrock/packages/ai/boost/internal/jsonconfig"
	"github.com/bedrock/packages/ai/boost/internal/platform"
)

// Platform re-exports the shared platform type for agent implementations.
type Platform = platform.Platform

// AgentOptions holds per-agent configuration overrides that replace the
// config('boost.agents.*') lookup pattern from laravel/boost.
type AgentOptions struct {
	McpConfigPath  string
	GuidelinesPath string
	SkillsPath     string
	GoBinary       string // default: "go"
	EntryPoint     string // default: "main.go"
}

// BaseAgent provides the shared implementation for all nine concrete agents.
// Embed this struct in every concrete agent.
type BaseAgent struct {
	opts AgentOptions
}

// NewBaseAgent constructs a BaseAgent from the provided options.
func NewBaseAgent(opts AgentOptions) BaseAgent {
	return BaseAgent{opts: opts}
}

// ---- CodingAgent shared implementations ---------------------------------

// McpConfigKey returns "mcpServers", the default top-level JSON key.
func (b *BaseAgent) McpConfigKey() string { return "mcpServers" }

// ShellMcpCommand returns "" — most agents use file-based installation.
func (b *BaseAgent) ShellMcpCommand() string { return "" }

// DefaultMcpConfig returns an empty map; the config file starts as {}.
func (b *BaseAgent) DefaultMcpConfig() map[string]any { return map[string]any{} }

// Frontmatter returns false; agents that require frontmatter override this.
func (b *BaseAgent) Frontmatter() bool { return false }

// UseAbsolutePathForMcp returns false by default.
func (b *BaseAgent) UseAbsolutePathForMcp() bool { return false }

// TransformGuidelines is the identity transform; subclasses may override.
func (b *BaseAgent) TransformGuidelines(markdown string) string { return markdown }

// GoBinaryPath returns the path to the Go runtime binary.
// When forceAbsolute is true or UseAbsolutePathForMcp() returns true, the
// binary is resolved via exec.LookPath. Adapts PHP's getPhpPath().
func (b *BaseAgent) GoBinaryPath(forceAbsolute bool) string {
	binary := b.opts.GoBinary

	if binary == "" {
		binary = "go"
	}

	if binary != "go" {
		return binary
	}

	if forceAbsolute || b.UseAbsolutePathForMcp() {
		if p, err := exec.LookPath(binary); err == nil {
			return p
		}

		if gr := runtime.GOROOT(); gr != "" {
			return filepath.Join(gr, "bin", "go")
		}
	}

	return binary
}

// EntryPointPath returns the application entry-point path.
// Adapts PHP's getArtisanPath(); defaults to "main.go".
func (b *BaseAgent) EntryPointPath(forceAbsolute bool) string {
	ep := b.opts.EntryPoint

	if ep == "" {
		ep = "main.go"
	}

	if forceAbsolute || b.UseAbsolutePathForMcp() {
		if abs, err := filepath.Abs(ep); err == nil {
			return abs
		}
	}

	return ep
}

// HttpMcpServerConfig returns the standard HTTP MCP payload.
// Cursor overrides this to return an npx mcp-remote command.
func (b *BaseAgent) HttpMcpServerConfig(url string) map[string]any {
	return map[string]any{
		"type": "http",
		"url":  url,
	}
}

// McpServerConfig returns the stdio MCP server payload.
func (b *BaseAgent) McpServerConfig(command string, args []string, env map[string]string) map[string]any {
	return map[string]any{
		"command": command,
		"args":    args,
		"env":     env,
	}
}

// InstallMcp writes a stdio MCP server entry. Concrete agents that have a
// config path call writeJSONConfigEntry via their own InstallMcp override.
func (b *BaseAgent) InstallMcp(key, command string, args []string, env map[string]string) (bool, error) {
	if b.opts.McpConfigPath == "" {
		return false, ErrNoMcpConfigPath
	}

	cmd, normalArgs := normalizeCommand(command, args)

	return writeJSONConfigEntry(
		b.opts.McpConfigPath,
		b.McpConfigKey(),
		key,
		b.McpServerConfig(cmd, normalArgs, env),
		b.DefaultMcpConfig(),
	)
}

// InstallHttpMcp writes an HTTP-transport MCP server entry.
func (b *BaseAgent) InstallHttpMcp(key, url string) (bool, error) {
	if b.opts.McpConfigPath == "" {
		return false, ErrNoMcpConfigPath
	}

	return writeJSONConfigEntry(
		b.opts.McpConfigPath,
		b.McpConfigKey(),
		key,
		b.HttpMcpServerConfig(url),
		b.DefaultMcpConfig(),
	)
}

// ---- Utilities ----------------------------------------------------------

// normalizeCommand splits a space-separated command string into command + args,
// but never splits absolute paths (which may contain spaces on macOS).
// Mirrors CommandNormalizer::normalize in laravel/boost.
func normalizeCommand(command string, extraArgs []string) (string, []string) {
	if filepath.IsAbs(command) || (len(command) > 2 && command[1] == ':') {
		return command, extraArgs
	}

	parts := strings.Fields(command)

	if len(parts) <= 1 {
		return command, extraArgs
	}

	combined := make([]string, 0, len(parts)-1+len(extraArgs))
	combined = append(combined, parts[1:]...)
	combined = append(combined, extraArgs...)

	return parts[0], combined
}

// writeJSONConfigEntry reads the config file, merges in the new server entry
// under configKey/serverKey, and writes it back atomically with file locking.
func writeJSONConfigEntry(
	configPath string,
	configKey string,
	serverKey string,
	serverConfig map[string]any,
	skeleton map[string]any,
) (bool, error) {
	_, err := jsonconfig.WriteEntry(configPath, configKey, serverKey, serverConfig, skeleton)

	if err != nil {
		return false, err
	}

	// This caller treats an already-existing entry as success.
	return true, nil
}

// ErrNoMcpConfigPath is returned when an agent's config path is empty.
var ErrNoMcpConfigPath = boosterr.ErrNoMcpConfigPath

// fallback returns override when non-empty, otherwise returns def.
func fallback(override, def string) string {
	if override != "" {
		return override
	}

	return def
}

// existsOnDisk reports whether path exists on the filesystem.
func existsOnDisk(path string) bool {
	_, err := os.Stat(expandHomeMarker(path))

	return err == nil
}

func expandHomeMarker(path string) string {
	if path == "~" {
		home, err := os.UserHomeDir()

		if err == nil && home != "" {
			return home
		}

		return path
	}

	if !strings.HasPrefix(path, "~/") {
		return path
	}

	home, err := os.UserHomeDir()

	if err != nil || home == "" {
		return path
	}

	return filepath.Join(home, filepath.FromSlash(path[2:]))
}

// commandInPath reports whether the binary name is in PATH.
func commandInPath(name string) bool {
	_, err := exec.LookPath(name)

	return err == nil
}

// commandExists runs the shell detection command and returns true on exit 0.
func commandExists(command string) bool {
	parts := strings.Fields(command)

	if len(parts) == 0 {
		return false
	}

	cmd := exec.Command(parts[0], parts[1:]...) //nolint:gosec

	return cmd.Run() == nil
}

// Ensure Platform is accessible from agent files without re-importing internal.
var _ Platform = platform.Darwin
