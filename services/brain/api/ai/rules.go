package ai

import (
	"errors"
	"path/filepath"

	"github.com/bedrock/packages/filesystem"
)

// RuleTarget describes one editor-rules file written by GenerateRules.
// The list mirrors laravel-brain's RulesExporter::TARGETS so any service
// scanned by brain ends up with the same AI assistant onboarding regardless
// of language. packages/ai/boost will own this table in a future commit.
type RuleTarget struct {
	Path  string
	Label string
}

// Targets is the canonical list of editor-rules files brain writes.
var Targets = []RuleTarget{
	{Path: "CLAUDE.md", Label: "Claude Code / Claude.ai"},
	{Path: "AGENTS.md", Label: "OpenAI Codex / generic AGENTS.md"},
	{Path: ".cursorrules", Label: "Cursor"},
	{Path: ".windsurfrules", Label: "Windsurf"},
	{Path: ".github/copilot-instructions.md", Label: "GitHub Copilot"},
	{Path: ".junie/guidelines.md", Label: "JetBrains Junie"},
	{Path: ".aider.conf.yml", Label: "Aider"},
}

// ErrConflict is returned when a target file already exists and force=false.
var ErrConflict = errors.New("rule target file already exists; pass force=true to overwrite")

// GenerateRules writes every Target under `root`. If a file exists and
// force is false, no files are written and ErrConflict is returned.
// On success, the returned slice lists the paths written.
//
// Uses packages/filesystem so directory creation and atomic-style writes
// stay consistent with the rest of bedrock.
func GenerateRules(root, body string, force bool) ([]string, error) {
	fs := filesystem.New()
	if !force {
		for _, t := range Targets {
			if fs.Exists(filepath.Join(root, t.Path)) {
				return nil, ErrConflict
			}
		}
	}
	written := make([]string, 0, len(Targets))
	for _, t := range Targets {
		full := filepath.Join(root, t.Path)
		if err := fs.Put(full, []byte(body)); err != nil {
			return written, err
		}
		written = append(written, t.Path)
	}
	return written, nil
}

