package tools

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// LastError reads the application log and returns the most recent error entry.
// Mirrors Laravel\Boost\Mcp\Tools\LastError.
// Tagged IsReadOnly.
type LastError struct {
	// LogFilePath overrides the default "storage/logs/app.log".
	LogFilePath string
}

func (t *LastError) Name() string     { return "last_error" }
func (t *LastError) IsReadOnly() bool { return true }

func (t *LastError) Description() string {
	return "Retrieve the last error or exception from the application log file."
}

func (t *LastError) Schema() map[string]any {
	return map[string]any{
		"type":       "object",
		"properties": map[string]any{},
		"required":   []string{},
	}
}

// errorLineRe matches PSR-3-style log lines at ERROR level or above.
var errorLineRe = regexp.MustCompile(`(?i)(ERROR|CRITICAL|EMERGENCY|ALERT|level=error|"level":"error")`)

// Handle reads the log file and returns the last error entry.
func (t *LastError) Handle(_ McpRequest) (McpResponse, error) {
	path := t.LogFilePath

	if path == "" {
		path = "storage/logs/app.log"
	}

	data, err := os.ReadFile(path)

	if err != nil {
		if os.IsNotExist(err) {
			return OkResponse(map[string]any{"error": nil, "message": "no log file found"}), nil
		}

		return ErrorResponse(fmt.Sprintf("last_error: %v", err)), nil
	}

	lines := strings.Split(string(data), "\n")
	// Walk backwards to find the last error line.
	for i := len(lines) - 1; i >= 0; i-- {
		if errorLineRe.MatchString(lines[i]) {
			// Collect the entry (current line + any following non-timestamp lines).
			entry := collectLogEntry(lines, i)

			return OkResponse(map[string]any{"error": entry}), nil
		}
	}

	return OkResponse(map[string]any{"error": nil, "message": "no errors found in log"}), nil
}

// collectLogEntry returns the log entry starting at index i, including
// continuation lines that don't start a new log entry.
func collectLogEntry(lines []string, start int) string {
	newEntryRe := regexp.MustCompile(`^\[?\d{4}-\d{2}-\d{2}`)

	var sb strings.Builder

	sb.WriteString(lines[start])

	for i := start + 1; i < len(lines); i++ {
		if newEntryRe.MatchString(lines[i]) {
			break
		}

		if lines[i] != "" {
			sb.WriteByte('\n')
			sb.WriteString(lines[i])
		}
	}

	return sb.String()
}
