package guidelines_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bedrock/packages/boost/guidelines"
	"github.com/bedrock/packages/config"
)

// TestGuidelineComposerEmptyConfig returns empty string when no guidelines exist.
func TestGuidelineComposerEmptyConfig(t *testing.T) {
	t.Parallel()

	repo := config.New(map[string]any{
		"boost.base_path": t.TempDir(),
	})
	cfg := guidelines.NewGuidelineConfig(repo)
	c := guidelines.NewGuidelineComposer(cfg)

	if got := c.Compose(); got != "" {
		t.Errorf("Compose() = %q, want empty string", got)
	}

	if used := c.Used(); len(used) != 0 {
		t.Errorf("Used() = %v, want empty", used)
	}
}

// TestGuidelineComposerCustomFiles discovers markdown files from the custom path.
func TestGuidelineComposerCustomFiles(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	customDir := filepath.Join(tmp, ".ai", "guidelines")

	if err := os.MkdirAll(customDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	content := "# My Guidelines\nDo the right thing."

	if err := os.WriteFile(filepath.Join(customDir, "project.md"), []byte(content), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	repo := config.New(map[string]any{
		"boost.base_path": tmp,
	})
	cfg := guidelines.NewGuidelineConfig(repo)
	c := guidelines.NewGuidelineComposer(cfg)

	result := c.Compose()

	if result == "" {
		t.Error("Compose() should not be empty when custom guideline files exist")
	}

	used := c.Used()
	found := false

	for _, u := range used {
		if u == "custom:project" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Used() = %v; expected 'custom:project' entry", used)
	}
}

// TestComposeGuidelinesHelper verifies the static ComposeGuidelines function.
func TestComposeGuidelinesHelper(t *testing.T) {
	t.Parallel()

	gs := []guidelines.Guideline{
		{Key: "a", Content: "# A\nContent A"},
		{Key: "b", Content: "# B\nContent B"},
	}

	result := guidelines.ComposeGuidelines(gs)

	if result == "" {
		t.Error("ComposeGuidelines should produce non-empty output")
	}

	// Both sections must appear.
	for _, want := range []string{"# A", "# B", "Content A", "Content B"} {
		if !containsString(result, want) {
			t.Errorf("ComposeGuidelines result missing %q", want)
		}
	}
}

func containsString(haystack, needle string) bool {
	return len(haystack) >= len(needle) &&
		(haystack == needle ||
			(len(needle) > 0 && indexString(haystack, needle) >= 0))
}

func indexString(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
