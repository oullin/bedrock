package agents_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bedrock/packages/ai/boost/agents"
	"github.com/bedrock/packages/ai/boost/internal/platform"
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

func TestJunieDetectOnSystemDarwinExpandsHome(t *testing.T) {
	tmp := t.TempDir()

	t.Setenv("HOME", tmp)

	a := agents.NewJunie()

	if a.DetectOnSystem(platform.Darwin) {
		t.Error("DetectOnSystem(platform.Darwin) should be false when JetBrains dir is absent")
	}

	jetbrainsDir := filepath.Join(tmp, "Library", "Application Support", "JetBrains")

	if err := os.MkdirAll(jetbrainsDir, 0755); err != nil {
		t.Fatalf("MkdirAll(JetBrains): %v", err)
	}

	if !a.DetectOnSystem(platform.Darwin) {
		t.Error("DetectOnSystem(platform.Darwin) should be true when JetBrains dir exists under HOME")
	}
}

func TestJunieDetectOnSystemLinuxExpandsHome(t *testing.T) {
	tmp := t.TempDir()

	t.Setenv("HOME", tmp)
	t.Setenv("PATH", "")

	a := agents.NewJunie()

	if a.DetectOnSystem(platform.Linux) {
		t.Error("DetectOnSystem(platform.Linux) should be false when ~/.config/junie is absent and junie is not in PATH")
	}

	junieDir := filepath.Join(tmp, ".config", "junie")

	if err := os.MkdirAll(junieDir, 0755); err != nil {
		t.Fatalf("MkdirAll(junie): %v", err)
	}

	if !a.DetectOnSystem(platform.Linux) {
		t.Error("DetectOnSystem(platform.Linux) should be true when ~/.config/junie exists under HOME")
	}
}
