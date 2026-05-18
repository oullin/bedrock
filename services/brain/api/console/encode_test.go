package console

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteJSONCreatesParentDirs(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "deeply", "nested", "out.json")

	payload := map[string]any{"hello": "world", "count": 3}

	if err := writeJSON(path, payload); err != nil {
		t.Fatalf("writeJSON: %v", err)
	}

	body, err := os.ReadFile(path)

	if err != nil {
		t.Fatalf("read written file: %v", err)
	}

	var got map[string]any

	if err := json.Unmarshal(body, &got); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if got["hello"] != "world" {
		t.Errorf("hello = %v, want world", got["hello"])
	}
	// Trailing newline per writeJSON's append('\n').
	if body[len(body)-1] != '\n' {
		t.Errorf("expected trailing newline; last byte = %q", body[len(body)-1])
	}
}

func TestWriteJSONReturnsEncodeError(t *testing.T) {
	// channels cannot be JSON-marshaled.
	err := writeJSON(filepath.Join(t.TempDir(), "x.json"), make(chan int))

	if err == nil {
		t.Fatal("expected encode error, got nil")
	}
}
