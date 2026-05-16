package ai

import (
	"errors"
	"os"
	"path/filepath"
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
func GenerateRules(root, body string, force bool) ([]string, error) {
	written := make([]string, 0, len(Targets))
	if !force {
		for _, t := range Targets {
			if _, err := os.Stat(filepath.Join(root, t.Path)); err == nil {
				return nil, ErrConflict
			}
		}
	}
	for _, t := range Targets {
		full := filepath.Join(root, t.Path)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			return written, err
		}
		if err := os.WriteFile(full, []byte(body), 0o644); err != nil {
			return written, err
		}
		written = append(written, t.Path)
	}
	return written, nil
}
