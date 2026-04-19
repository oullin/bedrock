package skills

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/bedrock/packages/boost/guidelines"
)

// SkillComposer discovers skills from boost built-ins, third-party packages,
// and user-defined .ai/skills/ directories.
// Mirrors Laravel\Boost\Install\SkillComposer (16 public/protected methods).
type SkillComposer struct {
	config *guidelines.GuidelineConfig
}

// NewSkillComposer constructs a SkillComposer.
func NewSkillComposer(cfg *guidelines.GuidelineConfig) *SkillComposer {
	return &SkillComposer{config: cfg}
}

// Skills returns all discovered skills from all sources.
func (c *SkillComposer) Skills() []Skill {
	var all []Skill

	all = append(all, c.BoostSkills()...)
	all = append(all, c.ThirdPartySkills()...)
	all = append(all, c.UserSkills()...)

	return all
}

// BoostSkills returns the built-in boost skills bundled with the package.
// In Go this is an empty list because bundled PHP guideline templates are not
// included; the caller registers skills via the user skills path.
func (c *SkillComposer) BoostSkills() []Skill { return []Skill{} }

// ThirdPartySkills returns skills discovered from configured third-party packages.
// Scans <basePath>/vendor/<pkg>/.ai/skills/ for each configured package.
func (c *SkillComposer) ThirdPartySkills() []Skill {
	var result []Skill

	for _, pkg := range c.config.Packages() {
		dir := filepath.Join(c.config.BasePath(), "vendor", pkg, ".ai", "skills")
		result = append(result, c.DiscoverSkillsFromPath(dir)...)
	}

	return result
}

// UserSkills returns all user-defined skills from .ai/skills/.
func (c *SkillComposer) UserSkills() []Skill {
	return c.DiscoverExplicitUserSkills()
}

// DiscoverExplicitUserSkills returns skills from <basePath>/.ai/skills/.
func (c *SkillComposer) DiscoverExplicitUserSkills() []Skill {
	dir := filepath.Join(c.config.BasePath(), ".ai", "skills")

	return c.DiscoverSkillsFromPath(dir)
}

// DiscoverSkillsFromPath walks the given directory and returns all SKILL.md files.
func (c *SkillComposer) DiscoverSkillsFromPath(dir string) []Skill {
	entries, err := os.ReadDir(dir)

	if err != nil {
		return []Skill{}
	}

	var result []Skill

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		skillFile := filepath.Join(dir, entry.Name(), "SKILL.md")
		skill, parseErr := c.ParseSkill(skillFile)

		if parseErr == nil && skill != nil {
			result = append(result, *skill)
		}
	}

	return result
}

// ParseSkill reads and parses a SKILL.md file at the given path.
// Returns nil if the file does not exist or cannot be parsed.
func (c *SkillComposer) ParseSkill(path string) (*Skill, error) {
	data, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	content := string(data)
	frontmatter := c.ParseSkillFrontmatter(content)

	name := frontmatter["name"]

	if name == "" {
		name = filepath.Base(filepath.Dir(path))
	}

	return &Skill{
		Name:        name,
		Description: frontmatter["description"],
		Path:        path,
		Content:     content,
	}, nil
}

// ParseSkillFrontmatter extracts YAML-style frontmatter from a SKILL.md file.
// Only handles simple key: value lines; no nested structures.
func (c *SkillComposer) ParseSkillFrontmatter(content string) map[string]string {
	result := make(map[string]string)
	scanner := bufio.NewScanner(strings.NewReader(content))

	inFrontmatter := false

	for scanner.Scan() {
		line := scanner.Text()

		if line == "---" {
			if !inFrontmatter {
				inFrontmatter = true

				continue
			}

			break // closing ---
		}

		if inFrontmatter {
			if idx := strings.IndexByte(line, ':'); idx > 0 {
				key := strings.TrimSpace(line[:idx])
				val := strings.TrimSpace(line[idx+1:])
				result[key] = val
			}
		}
	}

	return result
}
