package install

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SkillWriter writes Skill definitions as SKILL.md files into the agent's skills
// directory. Mirrors Upstream\Boost\Install\SkillWriter.
type SkillWriter struct{}

// SupportsSkillsPath is the minimal interface required by SkillWriter.
type SupportsSkillsPath interface {
	SkillsPath() string
}

// Skill is a local interface that carries the information a SkillWriter needs.
// It is satisfied by skills.Skill without a direct import dependency.
type Skill interface {
	SkillName() string
	SkillContent() string
}

// SkillEntry is a simple concrete that satisfies Skill. Callers may use this
// or implement their own.
type SkillEntry struct {
	Name    string
	Content string
}

func (s *SkillEntry) SkillName() string    { return s.Name }
func (s *SkillEntry) SkillContent() string { return s.Content }

// Write creates <agentSkillsPath>/<skill.Name>/SKILL.md for every provided
// skill. Parent directories are created as needed.
func (w *SkillWriter) Write(agent SupportsSkillsPath, skills []Skill) error {
	base := agent.SkillsPath()

	if base == "" {
		return fmt.Errorf("install: agent returned an empty skills path")
	}

	for _, s := range skills {
		name := s.SkillName()

		if name == "" {
			continue
		}

		dir := filepath.Join(base, name)

		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("install: create skill dir %s: %w", dir, err)
		}

		content := s.SkillContent()

		if !strings.HasSuffix(content, "\n") {
			content += "\n"
		}

		dest := filepath.Join(dir, "SKILL.md")

		if err := os.WriteFile(dest, []byte(content), 0644); err != nil {
			return fmt.Errorf("install: write skill %s: %w", dest, err)
		}
	}

	return nil
}
