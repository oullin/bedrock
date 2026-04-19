package agents_test

import (
	"testing"

	"github.com/bedrock/packages/boost/agents"
	"github.com/bedrock/packages/boost/internal/platform"
)

func TestCursorIdentity(t *testing.T) {
	t.Parallel()

	a := agents.NewCursor()

	if a.Name() != "cursor" {
		t.Errorf("Name() = %q, want \"cursor\"", a.Name())
	}

	if a.DisplayName() != "Cursor" {
		t.Errorf("DisplayName() = %q, want \"Cursor\"", a.DisplayName())
	}
}

func TestCursorDefaultPaths(t *testing.T) {
	t.Parallel()

	a := agents.NewCursor()

	if got := a.McpConfigPath(); got != ".cursor/mcp.json" {
		t.Errorf("McpConfigPath() = %q, want \".cursor/mcp.json\"", got)
	}

	if got := a.GuidelinesPath(); got != "AGENTS.md" {
		t.Errorf("GuidelinesPath() = %q, want \"AGENTS.md\"", got)
	}

	if got := a.SkillsPath(); got != ".cursor/skills" {
		t.Errorf("SkillsPath() = %q, want \".cursor/skills\"", got)
	}
}

func TestCursorMcpStrategy(t *testing.T) {
	t.Parallel()

	a := agents.NewCursor()

	if got := a.McpInstallationStrategy(); got != platform.McpStrategyFile {
		t.Errorf("McpInstallationStrategy() = %v, want McpStrategyFile", got)
	}
}

// TestCursorHttpMcpServerConfig mirrors CursorTest::test_http_mcp_uses_npx_mcp_remote.
func TestCursorHttpMcpServerConfig(t *testing.T) {
	t.Parallel()

	a := agents.NewCursor()
	cfg := a.HttpMcpServerConfig("http://localhost:9000")

	if cfg["command"] != "npx" {
		t.Errorf("command = %v, want \"npx\"", cfg["command"])
	}

	args, ok := cfg["args"].([]string)

	if !ok {
		t.Fatalf("args type = %T, want []string", cfg["args"])
	}

	if len(args) != 3 || args[0] != "-y" || args[1] != "mcp-remote" || args[2] != "http://localhost:9000" {
		t.Errorf("args = %v, want [\"-y\", \"mcp-remote\", \"http://localhost:9000\"]", args)
	}
}

func TestCursorDetectInProject(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	a := agents.NewCursor()

	if a.DetectInProject(tmp) {
		t.Error("DetectInProject should be false when .cursor dir absent")
	}
}
