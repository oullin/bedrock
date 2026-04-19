package install_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bedrock/packages/boost/install"
)

type mockGLAgent struct{ path string }

func (m *mockGLAgent) GuidelinesPath() string { return m.path }

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

func TestInstallGuidelineWriterEmptyPathErrors(t *testing.T) {
	t.Parallel()

	w := install.NewGuidelineWriter()

	if err := w.Write(&mockGLAgent{path: ""}, "content"); err == nil {
		t.Error("expected error for empty guidelines path")
	}
}

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
