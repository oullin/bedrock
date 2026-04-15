package agents_test

import (
	"testing"

	"github.com/bedrock/packages/boost/agents"
	"github.com/bedrock/packages/boost/internal/platform"
)

func TestAmpIdentity(t *testing.T) {
	t.Parallel()

	a := agents.NewAmp()

	if a.Name() != "amp" {
		t.Errorf("Name() = %q, want \"amp\"", a.Name())
	}

	if a.DisplayName() != "Amp" {
		t.Errorf("DisplayName() = %q, want \"Amp\"", a.DisplayName())
	}
}

func TestAmpMcpStrategyNone(t *testing.T) {
	t.Parallel()

	a := agents.NewAmp()

	if got := a.McpInstallationStrategy(); got != platform.McpStrategyNone {
		t.Errorf("McpInstallationStrategy() = %v, want McpStrategyNone", got)
	}
}

func TestAmpDefaultMcpConfigPathEmpty(t *testing.T) {
	t.Parallel()

	a := agents.NewAmp()

	if got := a.McpConfigPath(); got != "" {
		t.Errorf("McpConfigPath() = %q, want empty string", got)
	}
}

func TestAmpGuidelinesPath(t *testing.T) {
	t.Parallel()

	a := agents.NewAmp()

	if got := a.GuidelinesPath(); got != "AGENTS.md" {
		t.Errorf("GuidelinesPath() = %q, want \"AGENTS.md\"", got)
	}
}

func TestAmpSkillsPath(t *testing.T) {
	t.Parallel()

	a := agents.NewAmp()

	if got := a.SkillsPath(); got != ".amp/skills" {
		t.Errorf("SkillsPath() = %q, want \".amp/skills\"", got)
	}
}
