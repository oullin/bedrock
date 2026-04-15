package guidelines_test

import (
	"testing"

	"github.com/bedrock/packages/boost/guidelines"
	"github.com/bedrock/packages/config"
)

// TestGuidelineConfigDefaults mirrors GuidelineConfigTest::test_default_values.
func TestGuidelineConfigDefaults(t *testing.T) {
	t.Parallel()

	repo := config.New(map[string]any{})
	cfg := guidelines.NewGuidelineConfig(repo)

	if got := cfg.BasePath(); got != "." {
		t.Errorf("BasePath() = %q, want \".\"", got)
	}

	if got := cfg.CustomPath(); got != ".ai/guidelines" {
		t.Errorf("CustomPath() = %q, want \".ai/guidelines\"", got)
	}

	if pkgs := cfg.Packages(); len(pkgs) != 0 {
		t.Errorf("Packages() = %v, want empty slice", pkgs)
	}

	if cfg.HasSkillsEnabled() {
		t.Error("HasSkillsEnabled() should default to false")
	}

	if cfg.HasMcpEnabled() {
		t.Error("HasMcpEnabled() should default to false")
	}
}

// TestGuidelineConfigFromRepository verifies values are read from config.Repository.
func TestGuidelineConfigFromRepository(t *testing.T) {
	t.Parallel()

	repo := config.New(map[string]any{
		"boost.base_path":               "/app",
		"boost.custom_guideline_path":   "/app/.ai/guidelines",
		"boost.packages":                []any{"laravel/boost", "myorg/mypackage"},
		"boost.skills.enabled":          true,
		"boost.mcp.enabled":             true,
	})
	cfg := guidelines.NewGuidelineConfig(repo)

	if got := cfg.BasePath(); got != "/app" {
		t.Errorf("BasePath() = %q, want \"/app\"", got)
	}

	if got := cfg.CustomPath(); got != "/app/.ai/guidelines" {
		t.Errorf("CustomPath() = %q, want \"/app/.ai/guidelines\"", got)
	}

	pkgs := cfg.Packages()
	if len(pkgs) != 2 || pkgs[0] != "laravel/boost" || pkgs[1] != "myorg/mypackage" {
		t.Errorf("Packages() = %v, want [laravel/boost myorg/mypackage]", pkgs)
	}

	if !cfg.HasSkillsEnabled() {
		t.Error("HasSkillsEnabled() should be true")
	}

	if !cfg.HasMcpEnabled() {
		t.Error("HasMcpEnabled() should be true")
	}
}
