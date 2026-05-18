package agents_test

import (
	"testing"

	"github.com/bedrock/packages/ai/boost/agents"
	"github.com/bedrock/packages/ai/boost/internal/platform"
)

func TestAiderIdentity(t *testing.T) {
	t.Parallel()

	a := agents.NewAider()

	if a.Name() != "aider" {
		t.Errorf("Name() = %q, want \"aider\"", a.Name())
	}

	if a.DisplayName() != "Aider" {
		t.Errorf("DisplayName() = %q, want \"Aider\"", a.DisplayName())
	}
}

func TestAiderMcpStrategyNone(t *testing.T) {
	t.Parallel()

	a := agents.NewAider()

	if got := a.McpInstallationStrategy(); got != platform.McpStrategyNone {
		t.Errorf("McpInstallationStrategy() = %v, want McpStrategyNone", got)
	}
}

func TestAiderDefaultPaths(t *testing.T) {
	t.Parallel()

	a := agents.NewAider()

	if got := a.McpConfigPath(); got != "" {
		t.Errorf("McpConfigPath() = %q, want empty string", got)
	}

	if got := a.GuidelinesPath(); got != ".aider.conf.yml" {
		t.Errorf("GuidelinesPath() = %q, want \".aider.conf.yml\"", got)
	}

	if got := a.SkillsPath(); got != ".aider/skills" {
		t.Errorf("SkillsPath() = %q, want \".aider/skills\"", got)
	}
}

func TestAiderDetectInProject(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	a := agents.NewAider()

	if a.DetectInProject(tmp) {
		t.Error("DetectInProject should be false when .aider.conf.yml is absent")
	}
}
