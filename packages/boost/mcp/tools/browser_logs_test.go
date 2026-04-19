package tools_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bedrock/packages/boost/mcp/tools"
)

func TestBrowserLogsNoFile(t *testing.T) {
	t.Parallel()

	tool := &tools.BrowserLogs{LogFilePath: filepath.Join(t.TempDir(), "browser.log")}
	resp, err := tool.Handle(tools.McpRequest{})

	if err != nil {
		t.Fatalf("Handle: unexpected error: %v", err)
	}

	if resp.IsError {
		t.Errorf("Handle returned error response when file missing: %v", resp.Content)
	}

	data, ok := resp.Content[0].Data.(map[string]any)

	if !ok {
		t.Fatalf("content data type = %T, want map[string]any", resp.Content[0].Data)
	}

	entriesLen := 0

	switch v := data["entries"].(type) {
	case []string:
		entriesLen = len(v)
	case []any:
		entriesLen = len(v)
	}

	if entriesLen != 0 {
		t.Errorf("entries count = %d, want empty (type: %T)", entriesLen, data["entries"])
	}
}

func TestBrowserLogsReadsFile(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	path := filepath.Join(tmp, "browser.log")
	content := "line1\nline2\nline3\n"

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	tool := &tools.BrowserLogs{LogFilePath: path}
	resp, err := tool.Handle(tools.McpRequest{Args: map[string]any{"entries": float64(2)}})

	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	data, ok := resp.Content[0].Data.(map[string]any)

	if !ok {
		t.Fatalf("content data type = %T", resp.Content[0].Data)
	}

	// The entries value is []string from the implementation.
	entriesLen := 0

	switch v := data["entries"].(type) {
	case []string:
		entriesLen = len(v)
	case []any:
		entriesLen = len(v)
	}

	if entriesLen != 2 {
		t.Errorf("entries count = %d, want 2 (raw type: %T)", entriesLen, data["entries"])
	}
}

func TestBrowserLogsIsReadOnly(t *testing.T) {
	t.Parallel()

	if !(&tools.BrowserLogs{}).IsReadOnly() {
		t.Error("BrowserLogs should be read-only")
	}
}
