package guidelines

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// GuidelineAssist provides project-introspection helpers that guideline composers
// use to produce context-aware output.
// Mirrors upstream Boost\Install\GuidelineAssist (21 public methods).
type GuidelineAssist struct {
	config *GuidelineConfig
}

// NewGuidelineAssist constructs a GuidelineAssist.
func NewGuidelineAssist(cfg *GuidelineConfig) *GuidelineAssist {
	return &GuidelineAssist{config: cfg}
}

// ---- Project structure helpers ------------------------------------------

// Models scans the project base path for exported struct types defined in
// .go files and returns a map of typeName → absolute file path.
// Adapts PHP's models() which returns Eloquent model class names.
func (g *GuidelineAssist) Models() map[string]string {
	result := make(map[string]string)
	base := g.config.BasePath()

	_ = filepath.WalkDir(base, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			switch d.Name() {
			case "vendor", "node_modules", ".git":
				return filepath.SkipDir
			}

			return nil
		}

		if filepath.Ext(path) != ".go" {
			return nil
		}

		data, readErr := os.ReadFile(path)

		if readErr != nil {
			return nil
		}

		for _, line := range strings.Split(string(data), "\n") {
			trimmed := strings.TrimSpace(line)

			if strings.HasPrefix(trimmed, "type ") && strings.Contains(trimmed, " struct") {
				parts := strings.Fields(trimmed)

				if len(parts) >= 3 && parts[0] == "type" {
					result[parts[1]] = path
				}
			}
		}

		return nil
	})

	return result
}

// ShouldEnforceStrictTypes always returns false in Go (Go has no PHP strict_types).
func (g *GuidelineAssist) ShouldEnforceStrictTypes() bool { return false }

// NodePackageManager returns the Node package manager in use by inspecting the
// project base path for lock files (pnpm-lock.yaml, yarn.lock, package-lock.json).
func (g *GuidelineAssist) NodePackageManager() string {
	base := g.config.BasePath()

	if _, err := os.Stat(filepath.Join(base, "pnpm-lock.yaml")); err == nil {
		return "pnpm"
	}

	if _, err := os.Stat(filepath.Join(base, "yarn.lock")); err == nil {
		return "yarn"
	}

	return "npm"
}

// NodePackageManagerCommand returns the full command for the given npm script.
func (g *GuidelineAssist) NodePackageManagerCommand(command string) string {
	return g.NodePackageManager() + " run " + command
}

// HasPackage reports whether a named Go module dependency is listed in go.mod.
func (g *GuidelineAssist) HasPackage(pkg string) bool {
	base := g.config.BasePath()
	data, err := os.ReadFile(filepath.Join(base, "go.mod"))

	if err != nil {
		return false
	}

	return strings.Contains(string(data), pkg)
}

// AppPath joins base path with the given relative path.
// Adapts PHP's app_path() helper.
func (g *GuidelineAssist) AppPath(path string) string {
	return filepath.Join(g.config.BasePath(), path)
}

// ArtisanCommand returns the Go-equivalent entry-point command.
// Adapts PHP's artisan() which returns "php artisan {command}".
func (g *GuidelineAssist) ArtisanCommand(command string) string {
	return "go run . " + command
}

// ComposerCommand returns the Go module management equivalent.
// Adapts PHP's composerCommand() which returns "composer {command}".
func (g *GuidelineAssist) ComposerCommand(command string) string {
	return "go " + command
}

// BinCommand returns the ./vendor/bin equivalent for Go.
func (g *GuidelineAssist) BinCommand(command string) string {
	return "go run " + command
}

// SailBinaryPath returns the path to a Docker/compose entry-point if present.
// Adapts PHP's sailBinaryPath().
func (g *GuidelineAssist) SailBinaryPath() string {
	base := g.config.BasePath()

	if _, err := os.Stat(filepath.Join(base, "docker-compose.yml")); err == nil {
		return "docker compose"
	}

	if _, err := os.Stat(filepath.Join(base, "docker-compose.yaml")); err == nil {
		return "docker compose"
	}

	return ""
}

// SupportsPintAgentFormatter always returns false; this is PHP-specific.
func (g *GuidelineAssist) SupportsPintAgentFormatter() bool { return false }

// Skills returns the list of skills returned by SkillComposer. Callers inject
// the skill list to avoid a circular dependency.
func (g *GuidelineAssist) Skills() []string { return g.config.Packages() }

// HasSkillsEnabled delegates to config.
func (g *GuidelineAssist) HasSkillsEnabled() bool { return g.config.HasSkillsEnabled() }

// HasMcpEnabled delegates to config.
func (g *GuidelineAssist) HasMcpEnabled() bool { return g.config.HasMcpEnabled() }

// ---- Package detection --------------------------------------------------

// packageJSONField reads a string field from the project package.json.
func (g *GuidelineAssist) packageJSONField(field string) string {
	data, err := os.ReadFile(filepath.Join(g.config.BasePath(), "package.json"))

	if err != nil {
		return ""
	}

	var m map[string]any

	if err := json.Unmarshal(data, &m); err != nil {
		return ""
	}

	if v, ok := m[field].(string); ok {
		return v
	}

	return ""
}

// PackageName returns the project name from package.json or go.mod.
func (g *GuidelineAssist) PackageName() string {
	if name := g.packageJSONField("name"); name != "" {
		return name
	}

	base := g.config.BasePath()
	data, err := os.ReadFile(filepath.Join(base, "go.mod"))

	if err != nil {
		return ""
	}

	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimPrefix(line, "module ")
		}
	}

	return ""
}
