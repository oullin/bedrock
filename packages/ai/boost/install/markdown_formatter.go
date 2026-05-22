// Package install provides utilities for writing agent config files (MCP JSON,
// guideline Markdown, and skills) to disk.
// namespace.
package install

import (
	"regexp"
	"strings"
)

// frontmatterRe matches a YAML front-matter block at the very start of a file.
// Uses [\s\S]*? so it handles both non-empty and empty (---\n---\n) blocks.

// MarkdownFormatter applies lightweight formatting transformations to Markdown
// content.
type MarkdownFormatter struct{}

var frontmatterRe = regexp.MustCompile(`(?s)^---\n[\s\S]*?---\n?`)

// StripFrontmatter removes the YAML front-matter block (if any) from content.
func (f *MarkdownFormatter) StripFrontmatter(content string) string {
	return frontmatterRe.ReplaceAllString(content, "")
}

// AddFrontmatter prepends a YAML front-matter block to content.
// If content already has front-matter it is replaced.
func (f *MarkdownFormatter) AddFrontmatter(content string, data map[string]string) string {
	stripped := f.StripFrontmatter(content)

	if len(data) == 0 {
		return stripped
	}

	var sb strings.Builder

	sb.WriteString("---\n")

	for k, v := range data {
		sb.WriteString(k)
		sb.WriteString(": ")
		sb.WriteString(v)
		sb.WriteByte('\n')
	}

	sb.WriteString("---\n")
	sb.WriteString(stripped)

	return sb.String()
}

// NormalizeHeadings ensures the first heading in content is an H1 (`#`).
// If the document starts at H2 (`##`) or deeper all headings are promoted by
// one level.
func (f *MarkdownFormatter) NormalizeHeadings(content string) string {
	lines := strings.Split(content, "\n")

	// Find the minimum heading level used in the document.
	minLevel := 0

	for _, line := range lines {
		trimmed := strings.TrimLeft(line, "#")
		level := len(line) - len(trimmed)

		if level > 0 && (minLevel == 0 || level < minLevel) {
			minLevel = level
		}
	}

	if minLevel <= 1 {
		return content
	}

	promote := minLevel - 1
	result := make([]string, len(lines))

	for i, line := range lines {
		trimmed := strings.TrimLeft(line, "#")
		level := len(line) - len(trimmed)

		if level > 0 {
			result[i] = strings.Repeat("#", level-promote) + trimmed
		} else {
			result[i] = line
		}
	}

	return strings.Join(result, "\n")
}

// Trim removes leading and trailing whitespace from content.
func (f *MarkdownFormatter) Trim(content string) string {
	return strings.TrimSpace(content)
}
