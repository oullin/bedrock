package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ReadLogEntries reads the last N log entries from the application log.
// Handles both PSR-3 text format and JSON (structured) log formats.
// Mirrors upstream Boost\Mcp\Tools\ReadLogEntries.
// Tagged IsReadOnly.
type ReadLogEntries struct {
	// LogFilePath overrides the default "storage/logs/app.log".
	LogFilePath string
}

func (t *ReadLogEntries) Name() string     { return "read_log_entries" }
func (t *ReadLogEntries) IsReadOnly() bool { return true }

func (t *ReadLogEntries) Description() string {
	return "Read the last N log entries from the application log. " +
		"Handles both PSR-3 text and JSON structured log formats."
}

func (t *ReadLogEntries) Schema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"entries": map[string]any{
				"type":        "integer",
				"description": "Number of log entries to return. Defaults to 50.",
				"default":     50,
			},
		},
		"required": []string{},
	}
}

// psr3Re matches the start of a PSR-3 log line: [timestamp] channel.LEVEL: ...
var psr3Re = regexp.MustCompile(`^\[\d{4}-\d{2}-\d{2}[T ]\d{2}:\d{2}:\d{2}`)

// Handle reads the last N log entries.
func (t *ReadLogEntries) Handle(req McpRequest) (McpResponse, error) {
	n := 50

	if v, ok := req.Args["entries"]; ok {
		switch val := v.(type) {
		case int:
			n = val
		case float64:
			n = int(val)
		case string:
			if parsed, err := strconv.Atoi(val); err == nil {
				n = parsed
			}
		}
	}

	path := t.LogFilePath

	if path == "" {
		path = "storage/logs/app.log"
	}

	data, err := os.ReadFile(path)

	if err != nil {
		if os.IsNotExist(err) {
			return OkResponse(map[string]any{"entries": []any{}}), nil
		}

		return ErrorResponse(fmt.Sprintf("read_log_entries: %v", err)), nil
	}

	content := string(data)
	entries := parseLogEntries(content, n)

	return OkResponse(map[string]any{"entries": entries, "count": len(entries)}), nil
}

// parseLogEntries splits the log content into individual entries and returns the
// last n of them. Handles both PSR-3 and newline-delimited JSON.
func parseLogEntries(content string, n int) []any {
	lines := strings.Split(strings.TrimSpace(content), "\n")

	if len(lines) == 0 {
		return []any{}
	}

	// Detect JSON format by testing the first non-empty line.
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		if strings.HasPrefix(strings.TrimSpace(line), "{") {
			return parseJSONLog(lines, n)
		}

		break
	}

	return parsePSR3Log(lines, n)
}

func parsePSR3Log(lines []string, n int) []any {
	var entries []string

	var current strings.Builder

	for _, line := range lines {
		if psr3Re.MatchString(line) {
			if current.Len() > 0 {
				entries = append(entries, strings.TrimSpace(current.String()))
				current.Reset()
			}
		}

		current.WriteString(line)
		current.WriteByte('\n')
	}

	if current.Len() > 0 {
		entries = append(entries, strings.TrimSpace(current.String()))
	}

	if len(entries) > n {
		entries = entries[len(entries)-n:]
	}

	result := make([]any, len(entries))

	for i, e := range entries {
		result[i] = e
	}

	return result
}

func parseJSONLog(lines []string, n int) []any {
	var entries []any

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			continue
		}

		var m map[string]any

		if err := json.Unmarshal([]byte(trimmed), &m); err == nil {
			entries = append(entries, m)
		}
	}

	if len(entries) > n {
		entries = entries[len(entries)-n:]
	}

	return entries
}
