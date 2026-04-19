// Package guidelines provides guideline composition, project introspection, and
// file writing for the boost coding-assistant integration layer.
package guidelines

import (
	"github.com/bedrock/packages/config"
)

// GuidelineConfig holds settings for guideline composition. It wraps a
// config.Repository so callers get structured access via typed methods.
// Mirrors Laravel\Boost\Install\GuidelineConfig.
type GuidelineConfig struct {
	repo *config.Repository
}

// NewGuidelineConfig constructs a GuidelineConfig backed by the given repository.
func NewGuidelineConfig(repo *config.Repository) *GuidelineConfig {
	return &GuidelineConfig{repo: repo}
}

// BasePath returns the project base path for guideline discovery.
// Defaults to "." when not configured.
func (c *GuidelineConfig) BasePath() string {
	v, _ := c.repo.String("boost.base_path", ".")

	return v
}

// Packages returns the list of package names whose guidelines should be included.
func (c *GuidelineConfig) Packages() []string {
	v, _ := c.repo.Array("boost.packages", []any{})
	result := make([]string, 0, len(v))

	for _, item := range v {
		if s, ok := item.(string); ok {
			result = append(result, s)
		}
	}

	return result
}

// CustomPath returns the custom guidelines directory path.
// Defaults to ".ai/guidelines".
func (c *GuidelineConfig) CustomPath() string {
	v, _ := c.repo.String("boost.custom_guideline_path", ".ai/guidelines")

	return v
}

// HasSkillsEnabled reports whether skills support is enabled.
func (c *GuidelineConfig) HasSkillsEnabled() bool {
	v, _ := c.repo.Boolean("boost.skills.enabled", false)

	return v
}

// HasMcpEnabled reports whether MCP server support is enabled.
func (c *GuidelineConfig) HasMcpEnabled() bool {
	v, _ := c.repo.Boolean("boost.mcp.enabled", false)

	return v
}
