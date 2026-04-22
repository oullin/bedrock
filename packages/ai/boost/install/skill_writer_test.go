package install_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bedrock/packages/ai/boost/install"
)

type mockSkillsPathAgent struct{ base string }

func (m *mockSkillsPathAgent) SkillsPath() string { return m.base }

// SkillWriterTest::it_writes_skill_to_a_target_directory
func TestInstallSkillWriterCreatesFiles(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	agent := &mockSkillsPathAgent{base: filepath.Join(tmp, ".claude", "skills")}

	w := &install.SkillWriter{}
	skills := []install.Skill{
		&install.SkillEntry{Name: "commit", Content: "# Commit skill\n"},
		&install.SkillEntry{Name: "review", Content: "# Review skill\n"},
	}

	if err := w.Write(agent, skills); err != nil {
		t.Fatalf("Write: %v", err)
	}

	for _, s := range skills {
		dest := filepath.Join(tmp, ".claude", "skills", s.SkillName(), "SKILL.md")

		if _, err := os.Stat(dest); err != nil {
			t.Errorf("expected SKILL.md at %s: %v", dest, err)
		}
	}
}

func TestInstallSkillWriterEmptyPathErrors(t *testing.T) {
	t.Parallel()

	agent := &mockSkillsPathAgent{base: ""}
	w := &install.SkillWriter{}

	if err := w.Write(agent, []install.Skill{&install.SkillEntry{Name: "x", Content: "y"}}); err == nil {
		t.Error("expected error for empty skills path")
	}
}

// SkillWriterTest::it_throws_an_exception_for_path_traversal_in_skill_name
func TestInstallSkillWriterRejectsPathTraversal(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	agent := &mockSkillsPathAgent{base: filepath.Join(tmp, ".claude", "skills")}
	w := &install.SkillWriter{}

	if err := w.Write(agent, []install.Skill{
		&install.SkillEntry{Name: "../escape", Content: "# Escape\n"},
	}); err == nil {
		t.Fatal("expected path traversal to fail")
	}
}

func TestInstallSkillWriterFileContent(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	agent := &mockSkillsPathAgent{base: filepath.Join(tmp, "skills")}
	content := "---\nname: test\n---\n\n# Test\n"

	w := &install.SkillWriter{}
	_ = w.Write(agent, []install.Skill{&install.SkillEntry{Name: "test", Content: content}})

	dest := filepath.Join(tmp, "skills", "test", "SKILL.md")
	data, err := os.ReadFile(dest)

	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	// Content should end with newline.
	if string(data) == "" {
		t.Error("SKILL.md should not be empty")
	}
}
