package tools_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bedrock/packages/ai/boost/mcp/tools"
)

func TestReadLogEntriesNoFile(t *testing.T) {
	t.Parallel()

	tool := &tools.ReadLogEntries{LogFilePath: filepath.Join(t.TempDir(), "app.log")}
	resp, err := tool.Handle(tools.McpRequest{})

	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	if resp.IsError {
		t.Errorf("Handle should not error when file missing")
	}

	data, ok := resp.Content[0].Data.(map[string]any)

	if !ok {
		t.Fatalf("data type = %T", resp.Content[0].Data)
	}

	entries, _ := data["entries"].([]any)

	if len(entries) != 0 {
		t.Errorf("entries = %v, want empty", entries)
	}
}

func TestReadLogEntriesPSR3(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	path := filepath.Join(tmp, "app.log")
	content := "[2024-01-01T00:00:00] app.INFO: Entry one\n[2024-01-01T00:00:01] app.INFO: Entry two\n[2024-01-01T00:00:02] app.ERROR: Entry three\n"

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	tool := &tools.ReadLogEntries{LogFilePath: path}
	resp, err := tool.Handle(tools.McpRequest{Args: map[string]any{"entries": float64(2)}})

	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	data, ok := resp.Content[0].Data.(map[string]any)

	if !ok {
		t.Fatalf("data type = %T", resp.Content[0].Data)
	}

	entries, _ := data["entries"].([]any)

	if len(entries) != 2 {
		t.Errorf("entries count = %d, want 2", len(entries))
	}
}

func TestReadLogEntriesJSON(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	path := filepath.Join(tmp, "app.log")
	content := `{"level":"info","message":"one"}` + "\n" + `{"level":"error","message":"two"}` + "\n"

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}

	tool := &tools.ReadLogEntries{LogFilePath: path}
	resp, err := tool.Handle(tools.McpRequest{})

	if err != nil {
		t.Fatalf("Handle: %v", err)
	}

	data, ok := resp.Content[0].Data.(map[string]any)

	if !ok {
		t.Fatalf("data type = %T", resp.Content[0].Data)
	}

	entries, _ := data["entries"].([]any)

	if len(entries) != 2 {
		t.Errorf("JSON entries count = %d, want 2", len(entries))
	}
}

func TestReadLogEntriesIsReadOnly(t *testing.T) {
	t.Parallel()

	if !(&tools.ReadLogEntries{}).IsReadOnly() {
		t.Error("ReadLogEntries should be read-only")
	}
}
