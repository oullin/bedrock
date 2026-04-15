package guidelines

import (
	"fmt"
	"os"
	"path/filepath"
)

// SupportsGuidelines is a local re-declaration of the interface subset needed
// by GuidelineWriter, avoiding an import of the parent boost package.
type SupportsGuidelines interface {
	GuidelinesPath() string
	Frontmatter() bool
	TransformGuidelines(markdown string) string
}

// GuidelineWriter writes composed guidelines markdown to an agent's designated
// guidelines file.
// Mirrors Laravel\Boost\Install\GuidelineWriter.
type GuidelineWriter struct{}

// NewGuidelineWriter constructs a GuidelineWriter.
func NewGuidelineWriter() *GuidelineWriter { return &GuidelineWriter{} }

// Write transforms and writes content to agent.GuidelinesPath().
// Creates parent directories as needed.
func (w *GuidelineWriter) Write(agent SupportsGuidelines, content string) error {
	transformed := agent.TransformGuidelines(content)
	path := agent.GuidelinesPath()

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("boost: guidelines mkdir %s: %w", filepath.Dir(path), err)
	}

	if err := os.WriteFile(path, []byte(transformed), 0o644); err != nil {
		return fmt.Errorf("boost: guidelines write %s: %w", path, err)
	}

	return nil
}
