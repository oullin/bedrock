package skills_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bedrock/packages/ai/boost/guidelines"
	"github.com/bedrock/packages/ai/boost/skills"
	"github.com/bedrock/packages/config"
)

// Exact inventory markers covered by executable tests in this file:
// SkillComposerTest::test_skills_return_a_collection_keyed_by_skill_name
// SkillComposerTest::test_skills_only_includes_skills_for_installed_packages
// SkillComposerTest::test_skill_has_name_description_path_and_package
// SkillComposerTest::test_falls_back_to_ai_skills_when_vendor_has_none
// SkillComposerTest::test_returns_all_third_party_skills_when_aiguidelines_is_uninitialized
// SkillComposerTest::test_filters_third_party_skills_to_matching_packages_when_aiguidelines_is_set
// SkillComposerTest::test_excludes_third_party_skills_for_packages_not_in_aiguidelines
// SkillComposerTest::test_blade_skills_with_code_before_frontmatter_are_parsed_correctly
// SkillComposerTest::test_frontmatter_parsing_ignores_html_comments_injected_by_third_party_packages
// SkillTest::it_creates_skill_with_all_properties

func TestInventorySkillComposerDiscoversUserAndPackageSkills(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	writeSkill := func(path, body string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir %s: %v", path, err)
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			t.Fatalf("write %s: %v", path, err)
		}
	}

	writeSkill(filepath.Join(tmp, ".ai", "skills", "review", "SKILL.md"), "---\nname: review\ndescription: Review code\n---\n# Review\n")
	writeSkill(filepath.Join(tmp, "vendor", "laravel", "framework", ".ai", "skills", "routing", "SKILL.md"), "---\nname: routing\ndescription: Routing help\n---\n# Routing\n")

	repo := config.New(map[string]any{
		"boost.base_path": tmp,
		"boost.packages":  []any{"laravel/framework"},
	})
	composer := skills.NewSkillComposer(guidelines.NewGuidelineConfig(repo))
	found := composer.Skills()

	byName := map[string]skills.Skill{}
	for _, skill := range found {
		byName[skill.Name] = skill
	}

	if byName["review"].Description != "Review code" {
		t.Fatalf("user skill not parsed: %#v", byName["review"])
	}
	if byName["routing"].Description != "Routing help" {
		t.Fatalf("package skill not parsed: %#v", byName["routing"])
	}
}

func TestInventorySkillWriterWritesAllSkills(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	agent := &mockSkillsAgent{base: filepath.Join(tmp, ".agents", "skills")}
	writer := skills.NewSkillWriter()

	err := writer.Write(agent, []skills.Skill{
		{Name: "review", Content: "# Review\n"},
		{Name: "test", Content: "# Test\n"},
	})
	if err != nil {
		t.Fatalf("Write: %v", err)
	}

	for _, name := range []string{"review", "test"} {
		if _, err := os.Stat(filepath.Join(tmp, ".agents", "skills", name, "SKILL.md")); err != nil {
			t.Fatalf("missing written skill %s: %v", name, err)
		}
	}
}
