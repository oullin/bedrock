package agents_test

import (
	"testing"

	"github.com/bedrock/packages/boost/agents"
	"github.com/bedrock/packages/boost/internal/platform"
)

func TestGeminiIdentity(t *testing.T) {
	t.Parallel()

	a := agents.NewGemini()

	if a.Name() != "gemini" {
		t.Errorf("Name() = %q, want \"gemini\"", a.Name())
	}

	if a.DisplayName() != "Gemini" {
		t.Errorf("DisplayName() = %q, want \"Gemini\"", a.DisplayName())
	}
}

func TestGeminiDefaultPaths(t *testing.T) {
	t.Parallel()

	a := agents.NewGemini()

	if got := a.McpConfigPath(); got != ".gemini/settings.json" {
		t.Errorf("McpConfigPath() = %q, want \".gemini/settings.json\"", got)
	}

	if got := a.GuidelinesPath(); got != "GEMINI.md" {
		t.Errorf("GuidelinesPath() = %q, want \"GEMINI.md\"", got)
	}

	if got := a.SkillsPath(); got != ".gemini/skills" {
		t.Errorf("SkillsPath() = %q, want \".gemini/skills\"", got)
	}
}

func TestGeminiMcpStrategy(t *testing.T) {
	t.Parallel()

	a := agents.NewGemini()

	if got := a.McpInstallationStrategy(); got != platform.McpStrategyFile {
		t.Errorf("McpInstallationStrategy() = %v, want McpStrategyFile", got)
	}
}
