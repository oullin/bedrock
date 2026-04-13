package log

import (
	"encoding/json"
	"fmt"
	"time"
)

// Formatter converts a log record into a byte representation.
type Formatter interface {
	Format(record Record) ([]byte, error)
}

// LineFormatter produces human-readable single-line log output.
type LineFormatter struct {
	DateFormat     string
	IncludeContext bool
	LineSeparator  string
}

var _ Formatter = (*LineFormatter)(nil)

// NewLineFormatter creates a LineFormatter with sensible defaults.
func NewLineFormatter() *LineFormatter {
	return &LineFormatter{
		DateFormat:     time.RFC3339,
		IncludeContext: true,
		LineSeparator:  "\n",
	}
}

// Format produces output in the form:
// [2024-01-15T10:30:00Z] channel.LEVEL: Message {"key":"value"}
func (f *LineFormatter) Format(record Record) ([]byte, error) {
	ts := record.Time.Format(f.DateFormat)
	level := LevelName(record.Level)

	line := fmt.Sprintf("[%s] %s.%s: %s", ts, record.Channel, level, record.Message)

	if f.IncludeContext && len(record.Context) > 0 {
		ctx, err := json.Marshal(record.Context)

		if err != nil {
			return nil, fmt.Errorf("log: failed to marshal context: %w", err)
		}

		line += " " + string(ctx)
	}

	if len(record.Extra) > 0 {
		extra, err := json.Marshal(record.Extra)

		if err != nil {
			return nil, fmt.Errorf("log: failed to marshal extra: %w", err)
		}

		line += " " + string(extra)
	}

	line += f.LineSeparator

	return []byte(line), nil
}
