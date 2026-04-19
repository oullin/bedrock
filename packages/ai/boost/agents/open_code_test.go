package agents_test

import (
	"testing"

	"github.com/bedrock/packages/ai/boost/agents"
	"github.com/bedrock/packages/ai/boost/internal/platform"
)

func TestOpenCodeIdentity(t *testing.T) {
	t.Parallel()

	a := agents.NewOpenCode()

	if a.Name() != "opencode" {
		t.Errorf("Name() = %q, want \"opencode\"", a.Name())
	}

	if a.DisplayName() != "OpenCode" {
		t.Errorf("DisplayName() = %q, want \"OpenCode\"", a.DisplayName())
	}
}

func TestOpenCodeDefaultPaths(t *testing.T) {
	t.Parallel()

	a := agents.NewOpenCode()

	if got := a.McpConfigPath(); got != "opencode.json" {
		t.Errorf("McpConfigPath() = %q, want \"opencode.json\"", got)
	}

	if got := a.GuidelinesPath(); got != "AGENTS.md" {
		t.Errorf("GuidelinesPath() = %q, want \"AGENTS.md\"", got)
	}

	if got := a.SkillsPath(); got != ".opencode/skills" {
		t.Errorf("SkillsPath() = %q, want \".opencode/skills\"", got)
	}
}

func TestOpenCodeMcpStrategy(t *testing.T) {
	t.Parallel()

	a := agents.NewOpenCode()

	if got := a.McpInstallationStrategy(); got != platform.McpStrategyFile {
		t.Errorf("McpInstallationStrategy() = %v, want McpStrategyFile", got)
	}
}
