package agents_test

import (
	"testing"

	"github.com/bedrock/packages/boost/agents"
	"github.com/bedrock/packages/boost/internal/platform"
)

func TestJunieIdentity(t *testing.T) {
	t.Parallel()

	a := agents.NewJunie()

	if a.Name() != "junie" {
		t.Errorf("Name() = %q, want \"junie\"", a.Name())
	}

	if a.DisplayName() != "Junie" {
		t.Errorf("DisplayName() = %q, want \"Junie\"", a.DisplayName())
	}
}

func TestJunieDefaultPaths(t *testing.T) {
	t.Parallel()

	a := agents.NewJunie()

	if got := a.McpConfigPath(); got != ".junie/mcp.json" {
		t.Errorf("McpConfigPath() = %q, want \".junie/mcp.json\"", got)
	}

	if got := a.GuidelinesPath(); got != ".junie/guidelines.md" {
		t.Errorf("GuidelinesPath() = %q, want \".junie/guidelines.md\"", got)
	}

	if got := a.SkillsPath(); got != ".junie/skills" {
		t.Errorf("SkillsPath() = %q, want \".junie/skills\"", got)
	}
}

func TestJunieMcpStrategy(t *testing.T) {
	t.Parallel()

	a := agents.NewJunie()

	if got := a.McpInstallationStrategy(); got != platform.McpStrategyFile {
		t.Errorf("McpInstallationStrategy() = %v, want McpStrategyFile", got)
	}
}
