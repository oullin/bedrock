package agents_test

import (
	"testing"

	"github.com/bedrock/packages/ai/boost/agents"
	"github.com/bedrock/packages/ai/boost/internal/platform"
)

func TestCodexIdentity(t *testing.T) {
	t.Parallel()

	a := agents.NewCodex()

	if a.Name() != "codex" {
		t.Errorf("Name() = %q, want \"codex\"", a.Name())
	}

	if a.DisplayName() != "Codex" {
		t.Errorf("DisplayName() = %q, want \"Codex\"", a.DisplayName())
	}
}

func TestCodexDefaultPaths(t *testing.T) {
	t.Parallel()

	a := agents.NewCodex()

	if got := a.McpConfigPath(); got != "codex.json" {
		t.Errorf("McpConfigPath() = %q, want \"codex.json\"", got)
	}

	if got := a.GuidelinesPath(); got != "AGENTS.md" {
		t.Errorf("GuidelinesPath() = %q, want \"AGENTS.md\"", got)
	}
}

func TestCodexMcpStrategy(t *testing.T) {
	t.Parallel()

	a := agents.NewCodex()

	if got := a.McpInstallationStrategy(); got != platform.McpStrategyFile {
		t.Errorf("McpInstallationStrategy() = %v, want McpStrategyFile", got)
	}
}
