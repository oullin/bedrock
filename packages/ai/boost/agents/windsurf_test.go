package agents_test

import (
	"testing"

	"github.com/bedrock/packages/ai/boost/agents"
	"github.com/bedrock/packages/ai/boost/internal/platform"
)

func TestWindsurfIdentity(t *testing.T) {
	t.Parallel()

	a := agents.NewWindsurf()

	if a.Name() != "windsurf" {
		t.Errorf("Name() = %q, want \"windsurf\"", a.Name())
	}

	if a.DisplayName() != "Windsurf" {
		t.Errorf("DisplayName() = %q, want \"Windsurf\"", a.DisplayName())
	}
}

func TestWindsurfMcpStrategyNone(t *testing.T) {
	t.Parallel()

	a := agents.NewWindsurf()

	if got := a.McpInstallationStrategy(); got != platform.McpStrategyNone {
		t.Errorf("McpInstallationStrategy() = %v, want McpStrategyNone", got)
	}
}

func TestWindsurfDefaultPaths(t *testing.T) {
	t.Parallel()

	a := agents.NewWindsurf()

	if got := a.McpConfigPath(); got != "" {
		t.Errorf("McpConfigPath() = %q, want empty string", got)
	}

	if got := a.GuidelinesPath(); got != ".windsurfrules" {
		t.Errorf("GuidelinesPath() = %q, want \".windsurfrules\"", got)
	}

	if got := a.SkillsPath(); got != ".windsurf/skills" {
		t.Errorf("SkillsPath() = %q, want \".windsurf/skills\"", got)
	}
}

func TestWindsurfDetectInProject(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	a := agents.NewWindsurf()

	if a.DetectInProject(tmp) {
		t.Error("DetectInProject should be false when no Windsurf markers present")
	}
}
