package boost_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bedrock/packages/ai/boost"
)

// TestDetectorGetAgents verifies that GetAgents returns all 9 defaults.
func TestDetectorGetAgents(t *testing.T) {
	t.Parallel()

	m := boost.New()
	d := boost.NewDetector(m)
	got := d.GetAgents()

	if len(got) != 9 {
		t.Errorf("GetAgents: want 9, got %d", len(got))
	}
}

// TestDetectorDiscoverProjectInstalled verifies project detection by creating
// a temp directory with a CLAUDE.md file (triggering the claude_code agent).
func TestDetectorDiscoverProjectInstalled(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()

	// Write a CLAUDE.md file → should trigger claude_code
	if err := os.WriteFile(filepath.Join(tmp, "CLAUDE.md"), []byte("# guide"), 0644); err != nil {
		t.Fatalf("write CLAUDE.md: %v", err)
	}

	m := boost.New()
	d := boost.NewDetector(m)
	found := d.DiscoverProjectInstalledAgents(tmp)

	// At least claude_code should be detected.
	var names []string

	for _, a := range found {
		names = append(names, a.Name())
	}

	hasClaudeCode := false

	for _, n := range names {
		if n == "claude_code" {
			hasClaudeCode = true

			break
		}
	}

	if !hasClaudeCode {
		t.Errorf("expected claude_code in project-detected agents, got %v", names)
	}
}

// TestDetectorDiscoverSystemInstalled verifies that DiscoverSystemInstalledAgents
// returns a slice (possibly empty, since this runs in a CI-like env).
func TestDetectorDiscoverSystemInstalled(t *testing.T) {
	t.Parallel()

	m := boost.New()
	d := boost.NewDetector(m)
	got := d.DiscoverSystemInstalledAgents()

	// Just assert that the call doesn't panic and returns a slice.
	if got == nil {
		t.Error("DiscoverSystemInstalledAgents should return non-nil slice")
	}
}
