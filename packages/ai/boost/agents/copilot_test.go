package agents_test

import (
	"testing"

	"github.com/bedrock/packages/ai/boost/agents"
	"github.com/bedrock/packages/ai/boost/internal/platform"
)

func TestCopilotIdentity(t *testing.T) {
	t.Parallel()

	a := agents.NewCopilot()

	if a.Name() != "copilot" {
		t.Errorf("Name() = %q, want \"copilot\"", a.Name())
	}

	if a.DisplayName() != "GitHub Copilot" {
		t.Errorf("DisplayName() = %q, want \"GitHub Copilot\"", a.DisplayName())
	}
}

func TestCopilotDefaultPaths(t *testing.T) {
	t.Parallel()

	a := agents.NewCopilot()

	if got := a.McpConfigPath(); got != ".vscode/mcp.json" {
		t.Errorf("McpConfigPath() = %q, want \".vscode/mcp.json\"", got)
	}

	if got := a.GuidelinesPath(); got != ".github/copilot-instructions.md" {
		t.Errorf("GuidelinesPath() = %q, want \".github/copilot-instructions.md\"", got)
	}

	if got := a.SkillsPath(); got != ".github/copilot-skills" {
		t.Errorf("SkillsPath() = %q, want \".github/copilot-skills\"", got)
	}
}

func TestCopilotMcpStrategy(t *testing.T) {
	t.Parallel()

	a := agents.NewCopilot()

	if got := a.McpInstallationStrategy(); got != platform.McpStrategyFile {
		t.Errorf("McpInstallationStrategy() = %v, want McpStrategyFile", got)
	}
}
