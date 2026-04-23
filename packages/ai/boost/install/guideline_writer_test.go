package install_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bedrock/packages/ai/boost/install"
)

type mockGLAgent struct{ path string }

func (m *mockGLAgent) GuidelinesPath() string { return m.path }

// GuidelineWriterTest::test_it_writes_guidelines_to_new_file
func TestInstallGuidelineWriterCreatesFile(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	dest := filepath.Join(tmp, "CLAUDE.md")
	w := install.NewGuidelineWriter()

	if err := w.Write(&mockGLAgent{path: dest}, "# Guide\n"); err != nil {
		t.Fatalf("Write: %v", err)
	}

	data, err := os.ReadFile(dest)

	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	if string(data) != "# Guide\n" {
		t.Errorf("file content = %q, want \"# Guide\\n\"", string(data))
	}
}

// GuidelineWriterTest::test_it_throws_exception_when_directory_creation_fails
func TestInstallGuidelineWriterDirectoryCreationFails(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	blocker := filepath.Join(tmp, "blocked")

	if err := os.WriteFile(blocker, []byte("file"), 0o644); err != nil {
		t.Fatalf("write blocker file: %v", err)
	}

	w := install.NewGuidelineWriter()

	if err := w.Write(&mockGLAgent{path: filepath.Join(blocker, "AGENTS.md")}, "content"); err == nil {
		t.Fatal("expected error when parent directory cannot be created")
	}
}

// GuidelineWriterTest::test_it_creates_directory_when_it_does_not_exist
func TestInstallGuidelineWriterCreatesParentDirs(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	dest := filepath.Join(tmp, "deep", "nested", "AGENTS.md")
	w := install.NewGuidelineWriter()

	if err := w.Write(&mockGLAgent{path: dest}, "content"); err != nil {
		t.Fatalf("Write: %v", err)
	}

	if _, err := os.Stat(dest); err != nil {
		t.Errorf("expected file at %s: %v", dest, err)
	}
}
