package tools

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// BrowserLogs reads the last N entries from storage/logs/browser.log.
// Mirrors Upstream\Boost\Mcp\Tools\BrowserLogs.
// Tagged IsReadOnly.
type BrowserLogs struct {
	// LogFilePath overrides the default "storage/logs/browser.log".
	LogFilePath string
}

func (t *BrowserLogs) Name() string     { return "browser_logs" }
func (t *BrowserLogs) IsReadOnly() bool { return true }

func (t *BrowserLogs) Description() string {
	return "Retrieve the most recent browser/JavaScript log entries from storage/logs/browser.log."
}

func (t *BrowserLogs) Schema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"entries": map[string]any{
				"type":        "integer",
				"description": "Number of log entries to return. Defaults to 20.",
				"default":     20,
			},
		},
		"required": []string{},
	}
}

// Handle reads the browser log and returns the last N entries.
func (t *BrowserLogs) Handle(req McpRequest) (McpResponse, error) {
	n := 20
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
		path = "storage/logs/browser.log"
	}

	entries, err := readLastLines(path, n)
	if err != nil {
		return ErrorResponse(fmt.Sprintf("browser_logs: %v", err)), nil
	}

	return OkResponse(map[string]any{"entries": entries}), nil
}

// readLastLines returns the last n lines of the file at path.
func readLastLines(path string, n int) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}

		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if len(lines) > n {
		lines = lines[len(lines)-n:]
	}

	return lines, scanner.Err()
}
