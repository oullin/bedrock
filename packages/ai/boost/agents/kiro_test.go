package agents_test

import (
	"testing"

	"github.com/bedrock/packages/ai/boost/agents"
	"github.com/bedrock/packages/ai/boost/internal/platform"
)

func TestKiroIdentity(t *testing.T) {
	t.Parallel()

	a := agents.NewKiro()

	if a.Name() != "kiro" {
		t.Errorf("Name() = %q, want \"kiro\"", a.Name())
	}

	if a.DisplayName() != "Kiro" {
		t.Errorf("DisplayName() = %q, want \"Kiro\"", a.DisplayName())
	}
}

func TestKiroDefaultPaths(t *testing.T) {
	t.Parallel()

	a := agents.NewKiro()

	if got := a.McpConfigPath(); got != ".kiro/mcp.json" {
		t.Errorf("McpConfigPath() = %q, want \".kiro/mcp.json\"", got)
	}

	if got := a.GuidelinesPath(); got != ".kiro/steering/guidelines.md" {
		t.Errorf("GuidelinesPath() = %q, want \".kiro/steering/guidelines.md\"", got)
	}

	if got := a.SkillsPath(); got != ".kiro/skills" {
		t.Errorf("SkillsPath() = %q, want \".kiro/skills\"", got)
	}
}

func TestKiroMcpStrategy(t *testing.T) {
	t.Parallel()

	a := agents.NewKiro()

	if got := a.McpInstallationStrategy(); got != platform.McpStrategyFile {
		t.Errorf("McpInstallationStrategy() = %v, want McpStrategyFile", got)
	}
}
