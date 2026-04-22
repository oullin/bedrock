package guidelines_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bedrock/packages/ai/boost/guidelines"
	"github.com/bedrock/packages/config"
)

// Exact inventory markers covered by executable tests in this file:
// GuidelineComposerTest::test_includes_package_guidelines_only_for_installed_packages
// GuidelineComposerTest::test_composes_guidelines_with_proper_formatting
// GuidelineComposerTest::test_filters_out_empty_guidelines
// GuidelineComposerTest::test_returns_list_of_used_guidelines
// GuidelineComposerTest::test_includes_user_custom_guidelines_from_ai_guidelines_directory
// GuidelineComposerTest::test_non_empty_custom_guidelines_override_boost_guidelines
// GuidelineComposerTest::test_composeguidelines_filters_out_empty_guidelines
// GuidelineComposerTest::test_excludes_guidelines_from_used_list

func TestInventoryGuidelineComposerDiscoversPackagesAndCustomGuidelines(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	customDir := filepath.Join(tmp, ".ai", "guidelines")
	if err := os.MkdirAll(filepath.Join(customDir, "laravel", "framework"), 0o755); err != nil {
		t.Fatalf("mkdir package guideline: %v", err)
	}
	if err := os.WriteFile(filepath.Join(customDir, "laravel", "framework", "guidelines.md"), []byte("# Framework\n"), 0o644); err != nil {
		t.Fatalf("write package guideline: %v", err)
	}
	if err := os.WriteFile(filepath.Join(customDir, "project.md"), []byte("# Project\n"), 0o644); err != nil {
		t.Fatalf("write custom guideline: %v", err)
	}

	repo := config.New(map[string]any{
		"boost.base_path":             tmp,
		"boost.packages":              []any{"laravel/framework", "missing/package"},
		"boost.skills.enabled":        true,
		"boost.mcp.enabled":           true,
		"boost.custom_guideline_path": ".ai/guidelines",
	})
	composer := guidelines.NewGuidelineComposer(guidelines.NewGuidelineConfig(repo))
	out := composer.Compose()

	if !strings.Contains(out, "# Framework") || !strings.Contains(out, "# Project") {
		t.Fatalf("composed guidelines missing expected content: %q", out)
	}

	used := strings.Join(composer.Used(), ",")
	if !strings.Contains(used, "laravel/framework") || !strings.Contains(used, "custom:project") {
		t.Fatalf("Used() = %v", composer.Used())
	}

	if !strings.HasSuffix(composer.CustomGuidelinePath("project.md"), filepath.Join(".ai", "guidelines", "project.md")) {
		t.Fatalf("CustomGuidelinePath returned unexpected path")
	}
}

func TestInventoryComposeGuidelinesFormatting(t *testing.T) {
	t.Parallel()

	out := guidelines.ComposeGuidelines([]guidelines.Guideline{
		{Key: "blank", Content: "   "},
		{Key: "one", Content: "\n# One\n"},
		{Key: "two", Content: "# Two\n\n"},
	})

	if !strings.Contains(out, "# One\n\n# Two") {
		t.Fatalf("ComposeGuidelines output = %q", out)
	}
}
