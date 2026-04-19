package guidelines_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bedrock/packages/boost/guidelines"
)

type mockGuidelinesAgent struct {
	path string
}

func (m *mockGuidelinesAgent) GuidelinesPath() string               { return m.path }
func (m *mockGuidelinesAgent) Frontmatter() bool                    { return false }
func (m *mockGuidelinesAgent) TransformGuidelines(md string) string { return md }

// TestGuidelineWriterCreatesFile verifies the writer creates the file.
func TestGuidelineWriterCreatesFile(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	dest := filepath.Join(tmp, "CLAUDE.md")
	agent := &mockGuidelinesAgent{path: dest}

	w := guidelines.NewGuidelineWriter()
	content := "# My Guidelines\n\nDo the right thing."

	if err := w.Write(agent, content); err != nil {
		t.Fatalf("Write: %v", err)
	}

	data, err := os.ReadFile(dest)

	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	if string(data) != content {
		t.Errorf("file content = %q, want %q", string(data), content)
	}
}

// TestGuidelineWriterCreatesParentDirs verifies parent directories are created.
func TestGuidelineWriterCreatesParentDirs(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	dest := filepath.Join(tmp, "sub", "dir", "AGENTS.md")
	agent := &mockGuidelinesAgent{path: dest}

	w := guidelines.NewGuidelineWriter()

	if err := w.Write(agent, "content"); err != nil {
		t.Fatalf("Write: %v", err)
	}

	if _, err := os.Stat(dest); err != nil {
		t.Errorf("expected file at %s: %v", dest, err)
	}
}
