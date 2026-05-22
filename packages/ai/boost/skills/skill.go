// Package skills provides skill discovery, composition, and writing for the
// boost coding-assistant integration layer.
package skills

// Skill represents a single SKILL.md skill definition.
type Skill struct {
	// Name is the canonical skill identifier (e.g. "pest-testing").
	Name string
	// Description is the one-line skill description from the SKILL.md frontmatter.
	Description string
	// Package is the originating package name (empty for user-defined skills).
	Package string
	// Path is the absolute filesystem path to the SKILL.md file.
	Path string
	// Content is the full text of the SKILL.md file.
	Content string
}
