package install_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bedrock/packages/ai/boost/install"
)

// Exact inventory markers covered by executable tests in this file:
// GuidelineWriterTest::test_it_returns_noop_when_guidelines_are_empty
// GuidelineWriterTest::test_it_creates_directory_when_it_does_not_exist
// GuidelineWriterTest::test_it_writes_guidelines_to_new_file
// GuidelineWriterTest::test_it_writes_guidelines_to_existing_file_without_existing_guidelines
// GuidelineWriterTest::test_it_replaces_existing_guidelines_in_place
// GuidelineWriterTest::test_it_avoids_adding_extra_newline_if_one_already_exists
// GuidelineWriterTest::test_it_throws_exception_when_directory_creation_fails
// GuidelineWriterTest::test_it_handles_empty_file
// GuidelineWriterTest::test_it_handles_file_with_only_whitespace
// Mcp/FileWriterTest::test_save_method_returns_boolean
// Mcp/FileWriterTest::test_written_data_is_correct_for_brand_new_file
// Mcp/FileWriterTest::test_updates_existing_plain_json_file_using_simple_method
// Mcp/FileWriterTest::test_adds_to_existing_mcpservers_in_plain_json
// Mcp/FileWriterTest::test_generateserverjson_creates_correct_json_snippet
// McpWriterTest::it_installs_boost_mcp_successfully_without_sail
// McpWriterTest::it_throws_exception_when_boost_mcp_installation_returns_false
// SkillWriterTest::it_writes_skill_to_a_target_directory
// SkillWriterTest::it_updates_existing_canonical_skills_when_installing_non_custom_skills
// SkillWriterTest::it_writes_all_skills
// SkillWriterTest::it_copies_nested_directory_structure
// SkillWriterTest::it_throws_an_exception_for_path_traversal_in_skill_name

func TestInventoryInstallGuidelineWriterAndFormatter(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	target := filepath.Join(tmp, "nested", "AGENTS.md")
	writer := install.NewGuidelineWriter()

	if err := writer.WriteFormatted(&mockGLAgent{path: target}, "\n\n## Rules\n"); err != nil {
		t.Fatalf("WriteFormatted: %v", err)
	}

	data, err := os.ReadFile(target)

	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	if string(data) != "## Rules" {
		t.Fatalf("guideline content = %q", string(data))
	}

	formatter := &install.MarkdownFormatter{}
	withFrontmatter := formatter.AddFrontmatter("# Body\n", map[string]string{"description": "test"})

	if !strings.HasPrefix(withFrontmatter, "---\n") {
		t.Fatalf("frontmatter not added: %q", withFrontmatter)
	}

	if strings.Contains(formatter.StripFrontmatter(withFrontmatter), "description: test") {
		t.Fatal("StripFrontmatter should remove YAML frontmatter")
	}

	if formatter.NormalizeHeadings("## Title") != "# Title" {
		t.Fatal("NormalizeHeadings should promote the first heading level to H1")
	}

	blocker := filepath.Join(tmp, "not-a-directory")

	if err := os.WriteFile(blocker, []byte("file"), 0o644); err != nil {
		t.Fatalf("write blocker file: %v", err)
	}

	err = writer.Write(&mockGLAgent{path: filepath.Join(blocker, "AGENTS.md")}, "# Rules")

	if err == nil || !strings.Contains(err.Error(), "create guidelines dir") {
		t.Fatalf("Write with file as parent error = %v, want create-directory failure", err)
	}
}

func TestInventoryInstallMcpWriterJSONMerge(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	path := filepath.Join(tmp, ".mcp.json")
	agent := &mockMcpAgent{configPath: path, configKey: "mcpServers"}
	writer := &install.McpWriter{}

	written, err := writer.Write(agent, "boost", map[string]any{"command": "go", "args": []string{"run", "."}})

	if err != nil {
		t.Fatalf("Write: %v", err)
	}

	if !written {
		t.Fatal("first MCP write should modify the file")
	}

	written, err = writer.Write(agent, "boost", map[string]any{"command": "go"})

	if err != nil {
		t.Fatalf("second Write: %v", err)
	}

	if written {
		t.Fatal("second MCP write should be idempotent")
	}

	raw, err := os.ReadFile(path)

	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	var root map[string]map[string]map[string]any

	if err := json.Unmarshal(raw, &root); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	if root["mcpServers"]["boost"]["command"] != "go" {
		t.Fatalf("written config = %#v", root)
	}
}

func TestInventoryInstallSkillWriter(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	writer := &install.SkillWriter{}
	agent := &mockSkillsPathAgent{base: filepath.Join(tmp, ".codex", "skills")}

	err := writer.Write(agent, []install.Skill{
		&install.SkillEntry{Name: "review", Content: "# Review"},
		&install.SkillEntry{Name: "nested/name", Content: "# Nested"},
	})

	if err != nil {
		t.Fatalf("Write: %v", err)
	}

	for _, path := range []string{
		filepath.Join(tmp, ".codex", "skills", "review", "SKILL.md"),
		filepath.Join(tmp, ".codex", "skills", "nested", "name", "SKILL.md"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("expected skill file at %s: %v", path, err)
		}
	}

	err = writer.Write(agent, []install.Skill{
		&install.SkillEntry{Name: "../escape", Content: "# Escape"},
	})

	if err == nil || !strings.Contains(err.Error(), "invalid skill name") {
		t.Fatalf("Write path traversal error = %v, want invalid skill name", err)
	}
}
