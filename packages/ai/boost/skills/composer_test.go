package skills_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bedrock/packages/ai/boost/guidelines"
	"github.com/bedrock/packages/ai/boost/skills"
	"github.com/bedrock/packages/config"
)

// TestSkillComposerBoostSkills verifies BoostSkills returns an empty slice.
func TestSkillComposerBoostSkills(t *testing.T) {
	t.Parallel()

	cfg := guidelines.NewGuidelineConfig(config.New(map[string]any{}))
	c := skills.NewSkillComposer(cfg)

	if got := c.BoostSkills(); len(got) != 0 {
		t.Errorf("BoostSkills() = %v, want empty", got)
	}
}

// TestSkillComposerUserSkills discovers SKILL.md files from the user skills dir.
func TestSkillComposerUserSkills(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	skillDir := filepath.Join(tmp, ".ai", "skills", "my-skill")

	if err := os.MkdirAll(skillDir, 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	content := "---\nname: my-skill\ndescription: A test skill\n---\n\n# My Skill\n\nDoes things.\n"

	if err := os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte(content), 0644); err != nil {
		t.Fatalf("write SKILL.md: %v", err)
	}

	cfg := guidelines.NewGuidelineConfig(config.New(map[string]any{
		"boost.base_path": tmp,
	}))

	c := skills.NewSkillComposer(cfg)
	found := c.UserSkills()

	if len(found) != 1 {
		t.Fatalf("UserSkills() returned %d skills, want 1", len(found))
	}

	if found[0].Name != "my-skill" {
		t.Errorf("skill Name = %q, want \"my-skill\"", found[0].Name)
	}

	if found[0].Description != "A test skill" {
		t.Errorf("skill Description = %q, want \"A test skill\"", found[0].Description)
	}
}

// TestSkillComposerParseSkillFrontmatter verifies frontmatter extraction.
func TestSkillComposerParseSkillFrontmatter(t *testing.T) {
	t.Parallel()

	cfg := guidelines.NewGuidelineConfig(config.New(map[string]any{}))
	c := skills.NewSkillComposer(cfg)

	content := "---\nname: test\ndescription: Test description\nauthor: tester\n---\n\n# Body\n"
	fm := c.ParseSkillFrontmatter(content)

	cases := map[string]string{
		"name":        "test",
		"description": "Test description",
		"author":      "tester",
	}

	for k, want := range cases {
		if got := fm[k]; got != want {
			t.Errorf("frontmatter[%q] = %q, want %q", k, got, want)
		}
	}
}
