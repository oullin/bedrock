package guidelines_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bedrock/packages/ai/boost/guidelines"
)

type mockGuidelinesAgent struct {
	path string
}

func (m *mockGuidelinesAgent) GuidelinesPath() string               { return m.path }
func (m *mockGuidelinesAgent) Frontmatter() bool                    { return false }
func (m *mockGuidelinesAgent) TransformGuidelines(md string) string { return md }

// GuidelineWriterTest::test_it_writes_guidelines_to_new_file
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

// GuidelineWriterTest::test_it_creates_directory_when_it_does_not_exist
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

// GuidelineWriterTest::test_it_throws_exception_when_directory_creation_fails
func TestGuidelineWriterDirectoryCreationFails(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	blocker := filepath.Join(tmp, "blocked")

	if err := os.WriteFile(blocker, []byte("file"), 0o644); err != nil {
		t.Fatalf("write blocker file: %v", err)
	}

	agent := &mockGuidelinesAgent{path: filepath.Join(blocker, "AGENTS.md")}
	w := guidelines.NewGuidelineWriter()

	if err := w.Write(agent, "content"); err == nil {
		t.Fatal("expected error when parent directory cannot be created")
	}
}
