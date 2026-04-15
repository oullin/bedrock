package guidelines

import (
	"os"
	"path/filepath"
	"strings"
)

// Guideline holds a named piece of markdown to be included in the composed output.
type Guideline struct {
	Key     string
	Content string
}

// GuidelineComposer assembles guidelines from all configured sources into a
// single markdown document.
// Mirrors Upstream\Boost\Install\GuidelineComposer.
type GuidelineComposer struct {
	config *GuidelineConfig
	used   []string
}

// NewGuidelineComposer constructs a GuidelineComposer.
func NewGuidelineComposer(cfg *GuidelineConfig) *GuidelineComposer {
	return &GuidelineComposer{config: cfg}
}

// Compose discovers all guidelines and concatenates them into a single markdown
// string. Populates the internal used-keys list.
func (c *GuidelineComposer) Compose() string {
	guidelines := c.Guidelines()
	return ComposeGuidelines(guidelines)
}

// Guidelines discovers and returns all active guidelines from:
//  1. Custom guideline files in config.CustomPath()
//  2. Package-specific guideline files listed in config.Packages()
func (c *GuidelineComposer) Guidelines() []Guideline {
	c.used = nil
	var all []Guideline

	// 1. Package-specific guidelines
	for _, pkg := range c.config.Packages() {
		content := c.loadPackageGuideline(pkg)
		if content != "" {
			all = append(all, Guideline{Key: pkg, Content: content})
			c.used = append(c.used, pkg)
		}
	}

	// 2. Custom user guidelines from .ai/guidelines/
	custom := c.discoverCustomGuidelines()
	all = append(all, custom...)

	return all
}

// Used returns the list of guideline keys that were included in the last Compose call.
func (c *GuidelineComposer) Used() []string {
	if c.used == nil {
		return []string{}
	}

	return c.used
}

// CustomGuidelinePath returns the absolute path to the custom guideline directory.
func (c *GuidelineComposer) CustomGuidelinePath(path string) string {
	return filepath.Join(c.config.BasePath(), c.config.CustomPath(), path)
}

// ComposeGuidelines concatenates a slice of Guideline values into a single
// markdown string. Static equivalent for use without a full GuidelineComposer.
func ComposeGuidelines(guidelines []Guideline) string {
	var sb strings.Builder

	for i, g := range guidelines {
		if i > 0 {
			sb.WriteString("\n\n")
		}
		sb.WriteString(strings.TrimSpace(g.Content))
	}

	return sb.String()
}

// loadPackageGuideline attempts to load a guideline file for the named package
// from the custom path directory. Returns "" when not found.
func (c *GuidelineComposer) loadPackageGuideline(pkg string) string {
	candidates := []string{
		filepath.Join(c.config.BasePath(), c.config.CustomPath(), pkg+".md"),
		filepath.Join(c.config.BasePath(), c.config.CustomPath(), pkg, "guidelines.md"),
	}

	for _, path := range candidates {
		data, err := os.ReadFile(path)
		if err == nil {
			return string(data)
		}
	}

	return ""
}

// discoverCustomGuidelines scans the custom guidelines directory for .md files.
func (c *GuidelineComposer) discoverCustomGuidelines() []Guideline {
	dir := filepath.Join(c.config.BasePath(), c.config.CustomPath())
	entries, err := os.ReadDir(dir)

	if err != nil {
		return nil
	}

	var result []Guideline

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".md" {
			continue
		}

		key := strings.TrimSuffix(entry.Name(), ".md")
		data, readErr := os.ReadFile(filepath.Join(dir, entry.Name()))

		if readErr != nil {
			continue
		}

		result = append(result, Guideline{Key: "custom:" + key, Content: string(data)})
		c.used = append(c.used, "custom:"+key)
	}

	return result
}
