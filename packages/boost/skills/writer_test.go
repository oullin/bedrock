package skills_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bedrock/packages/boost/skills"
)

type mockSkillsAgent struct{ base string }

func (m *mockSkillsAgent) SkillsPath() string { return m.base }

// TestSkillWriterCreatesFiles verifies SKILL.md files are created.
func TestSkillWriterCreatesFiles(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	agent := &mockSkillsAgent{base: filepath.Join(tmp, ".claude", "skills")}

	w := skills.NewSkillWriter()
	ss := []skills.Skill{
		{Name: "commit", Content: "---\nname: commit\n---\n\n# Commit\n"},
		{Name: "review", Content: "---\nname: review\n---\n\n# Review\n"},
	}

	if err := w.Write(agent, ss); err != nil {
		t.Fatalf("Write: %v", err)
	}

	for _, s := range ss {
		dest := filepath.Join(tmp, ".claude", "skills", s.Name, "SKILL.md")
		data, err := os.ReadFile(dest)

		if err != nil {
			t.Errorf("expected SKILL.md at %s: %v", dest, err)
			continue
		}

		if string(data) != s.Content {
			t.Errorf("SKILL.md content = %q, want %q", string(data), s.Content)
		}
	}
}

// TestSkillWriterCreatesParentDirs verifies that nested skill dirs are created.
func TestSkillWriterCreatesParentDirs(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	// Non-existent nested path.
	agent := &mockSkillsAgent{base: filepath.Join(tmp, "a", "b", "c", "skills")}

	w := skills.NewSkillWriter()
	ss := []skills.Skill{{Name: "test", Content: "content"}}

	if err := w.Write(agent, ss); err != nil {
		t.Fatalf("Write: %v", err)
	}

	dest := filepath.Join(tmp, "a", "b", "c", "skills", "test", "SKILL.md")

	if _, err := os.Stat(dest); err != nil {
		t.Errorf("expected SKILL.md at %s: %v", dest, err)
	}
}
