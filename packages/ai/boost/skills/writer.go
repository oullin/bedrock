package skills

import (
	"fmt"
	"os"
	"path/filepath"
)

// SupportsSkills is a local re-declaration of the interface subset needed
// by SkillWriter, avoiding an import of the parent boost package.
type SupportsSkills interface {
	SkillsPath() string
}

// SkillWriter writes SKILL.md files to an agent's skills directory.
// Mirrors Laravel\Boost\Install\SkillWriter.
type SkillWriter struct{}

// NewSkillWriter constructs a SkillWriter.
func NewSkillWriter() *SkillWriter { return &SkillWriter{} }

// Write writes each skill's content to <agent.SkillsPath()>/<skill.Name>/SKILL.md.
// Creates parent directories as needed.
func (w *SkillWriter) Write(agent SupportsSkills, skills []Skill) error {
	base := agent.SkillsPath()

	for _, skill := range skills {
		dir := filepath.Join(base, skill.Name)

		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("boost: skills mkdir %s: %w", dir, err)
		}

		path := filepath.Join(dir, "SKILL.md")

		if err := os.WriteFile(path, []byte(skill.Content), 0o644); err != nil {
			return fmt.Errorf("boost: skills write %s: %w", path, err)
		}
	}

	return nil
}
