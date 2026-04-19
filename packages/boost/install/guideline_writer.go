package install

import (
	"fmt"
	"os"
	"path/filepath"
)

// GuidelineWriter writes composed guideline Markdown content to the path
// dictated by an agent that implements SupportsGuidelinesPath.
// Mirrors Laravel\Boost\Install\GuidelineWriter.
type GuidelineWriter struct {
	formatter *MarkdownFormatter
}

// NewGuidelineWriter returns a GuidelineWriter with the default MarkdownFormatter.

// SupportsGuidelinesPath is the minimal interface required by GuidelineWriter;
// satisfied by any coding agent that has a guidelines file path.
type SupportsGuidelinesPath interface {
	GuidelinesPath() string
}

func NewGuidelineWriter() *GuidelineWriter {
	return &GuidelineWriter{formatter: &MarkdownFormatter{}}
}

// Write writes content to the agent's guidelines path, creating any missing
// parent directories. Returns an error if the write fails.
func (w *GuidelineWriter) Write(agent SupportsGuidelinesPath, content string) error {
	path := agent.GuidelinesPath()

	if path == "" {
		return fmt.Errorf("install: agent returned an empty guidelines path")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("install: create guidelines dir: %w", err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("install: write guidelines: %w", err)
	}

	return nil
}

// WriteFormatted applies the MarkdownFormatter before writing.
func (w *GuidelineWriter) WriteFormatted(agent SupportsGuidelinesPath, content string) error {
	return w.Write(agent, w.formatter.Trim(content))
}
