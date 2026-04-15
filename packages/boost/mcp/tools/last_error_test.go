package tools_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bedrock/packages/boost/mcp/tools"
)

func TestLastErrorNoFile(t *testing.T) {
	t.Parallel()

	tool := &tools.LastError{LogFilePath: filepath.Join(t.TempDir(), "app.log")}
	resp, err := tool.Handle(tools.McpRequest{})

	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if resp.IsError {
		t.Errorf("Handle should not return error when log file missing")
	}

	data, ok := resp.Content[0].Data.(map[string]any)
	if !ok {
		t.Fatalf("data type = %T", resp.Content[0].Data)
	}

	if data["error"] != nil {
		t.Errorf("error = %v, want nil", data["error"])
	}
}

func TestLastErrorFindsError(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	path := filepath.Join(tmp, "app.log")
	logContent := "[2024-01-01T00:00:00] app.INFO: Starting up\n[2024-01-01T00:00:01] app.ERROR: Something failed here\n"

	if err := os.WriteFile(path, []byte(logContent), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	tool := &tools.LastError{LogFilePath: path}
	resp, err := tool.Handle(tools.McpRequest{})

	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if resp.IsError {
		t.Errorf("Handle returned error response: %v", resp.Content)
	}

	data, ok := resp.Content[0].Data.(map[string]any)
	if !ok {
		t.Fatalf("data type = %T", resp.Content[0].Data)
	}

	if data["error"] == nil {
		t.Error("expected error entry to be non-nil")
	}
}

func TestLastErrorNoErrors(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	path := filepath.Join(tmp, "app.log")

	if err := os.WriteFile(path, []byte("[2024-01-01] app.INFO: all good\n"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	tool := &tools.LastError{LogFilePath: path}
	resp, _ := tool.Handle(tools.McpRequest{})

	data, ok := resp.Content[0].Data.(map[string]any)
	if !ok {
		t.Fatalf("data type = %T", resp.Content[0].Data)
	}

	if data["error"] != nil {
		t.Errorf("error = %v, want nil when no error lines", data["error"])
	}
}

func TestLastErrorIsReadOnly(t *testing.T) {
	t.Parallel()

	if !(&tools.LastError{}).IsReadOnly() {
		t.Error("LastError should be read-only")
	}
}
