package agents_test

import (
	"testing"

	"github.com/bedrock/packages/ai/boost/agents"
	"github.com/bedrock/packages/ai/boost/internal/platform"
)

func TestClaudeCodeIdentity(t *testing.T) {
	t.Parallel()

	a := agents.NewClaudeCode()

	if a.Name() != "claude_code" {
		t.Errorf("Name() = %q, want \"claude_code\"", a.Name())
	}

	if a.DisplayName() != "Claude Code" {
		t.Errorf("DisplayName() = %q, want \"Claude Code\"", a.DisplayName())
	}
}

func TestClaudeCodeDefaultPaths(t *testing.T) {
	t.Parallel()

	a := agents.NewClaudeCode()

	if got := a.McpConfigPath(); got != ".mcp.json" {
		t.Errorf("McpConfigPath() = %q, want \".mcp.json\"", got)
	}

	if got := a.GuidelinesPath(); got != "CLAUDE.md" {
		t.Errorf("GuidelinesPath() = %q, want \"CLAUDE.md\"", got)
	}

	if got := a.SkillsPath(); got != ".claude/skills" {
		t.Errorf("SkillsPath() = %q, want \".claude/skills\"", got)
	}
}

// TestClaudeCodePathOverrides verifies AgentOptions overrides take effect.
func TestClaudeCodePathOverrides(t *testing.T) {
	t.Parallel()

	a := agents.NewClaudeCode(agents.AgentOptions{
		McpConfigPath:  "custom/.mcp.json",
		GuidelinesPath: "custom/CLAUDE.md",
		SkillsPath:     "custom/skills",
	})

	if got := a.McpConfigPath(); got != "custom/.mcp.json" {
		t.Errorf("McpConfigPath() = %q, want \"custom/.mcp.json\"", got)
	}

	if got := a.GuidelinesPath(); got != "custom/CLAUDE.md" {
		t.Errorf("GuidelinesPath() = %q, want \"custom/CLAUDE.md\"", got)
	}

	if got := a.SkillsPath(); got != "custom/skills" {
		t.Errorf("SkillsPath() = %q, want \"custom/skills\"", got)
	}
}

func TestClaudeCodeMcpStrategy(t *testing.T) {
	t.Parallel()

	a := agents.NewClaudeCode()

	if got := a.McpInstallationStrategy(); got != platform.McpStrategyFile {
		t.Errorf("McpInstallationStrategy() = %v, want McpStrategyFile", got)
	}
}

func TestClaudeCodeUseAbsolutePath(t *testing.T) {
	t.Parallel()

	a := agents.NewClaudeCode()

	if a.UseAbsolutePathForMcp() {
		t.Error("UseAbsolutePathForMcp() should be false")
	}
}

// TestClaudeCodeGoBinaryPath verifies GoBinaryPath(false) returns "go".
func TestClaudeCodeGoBinaryPath(t *testing.T) {
	t.Parallel()

	a := agents.NewClaudeCode()

	if got := a.GoBinaryPath(false); got != "go" {
		t.Errorf("GoBinaryPath(false) = %q, want \"go\"", got)
	}
}

// TestClaudeCodeHttpMcpServerConfig verifies the HTTP MCP payload shape.
func TestClaudeCodeHttpMcpServerConfig(t *testing.T) {
	t.Parallel()

	a := agents.NewClaudeCode()
	cfg := a.HttpMcpServerConfig("http://localhost:9000")

	if cfg["type"] != "http" {
		t.Errorf("HttpMcpServerConfig type = %v, want \"http\"", cfg["type"])
	}

	if cfg["url"] != "http://localhost:9000" {
		t.Errorf("HttpMcpServerConfig url = %v, want \"http://localhost:9000\"", cfg["url"])
	}
}

// TestClaudeCodeDetectInProject verifies project detection.
func TestClaudeCodeDetectInProject(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()

	a := agents.NewClaudeCode()

	// No markers → false.
	if a.DetectInProject(tmp) {
		t.Error("DetectInProject should return false when no markers present")
	}
}
