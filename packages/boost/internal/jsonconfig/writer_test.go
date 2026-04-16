package jsonconfig_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/bedrock/packages/boost/internal/jsonconfig"
)

func TestWriteEntryCreatesFile(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.json")
	written, err := jsonconfig.WriteEntry(path, "mcpServers", "boost", map[string]any{
		"command": "go",
		"args":    []any{"run", "."},
	}, nil)

	if err != nil {
		t.Fatalf("WriteEntry: %v", err)
	}
	if !written {
		t.Error("expected written=true for new entry")
	}

	root := readJSON(t, path)
	servers, ok := root["mcpServers"].(map[string]any)
	if !ok {
		t.Fatal("mcpServers key missing or wrong type")
	}
	if _, ok := servers["boost"]; !ok {
		t.Error("boost entry not found")
	}
}

func TestWriteEntryUsesSkeleton(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.json")
	skeleton := map[string]any{"version": "1.0"}

	_, err := jsonconfig.WriteEntry(path, "mcpServers", "boost", map[string]any{"command": "go"}, skeleton)
	if err != nil {
		t.Fatalf("WriteEntry: %v", err)
	}

	root := readJSON(t, path)
	if root["version"] != "1.0" {
		t.Errorf("skeleton key version=%v, want 1.0", root["version"])
	}
}

func TestWriteEntryIdempotent(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.json")
	cfg := map[string]any{"command": "go"}

	if _, err := jsonconfig.WriteEntry(path, "mcpServers", "boost", cfg, nil); err != nil {
		t.Fatalf("first WriteEntry: %v", err)
	}

	written, err := jsonconfig.WriteEntry(path, "mcpServers", "boost", cfg, nil)
	if err != nil {
		t.Fatalf("second WriteEntry: %v", err)
	}
	if written {
		t.Error("expected written=false for existing entry")
	}
}

func TestWriteEntryMergesExisting(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.json")

	existing := map[string]any{
		"mcpServers": map[string]any{
			"other": map[string]any{"command": "other"},
		},
	}
	data, _ := json.MarshalIndent(existing, "", "    ")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	_, err := jsonconfig.WriteEntry(path, "mcpServers", "new-server", map[string]any{"command": "go"}, nil)
	if err != nil {
		t.Fatalf("WriteEntry: %v", err)
	}

	root := readJSON(t, path)
	servers, _ := root["mcpServers"].(map[string]any)

	if _, ok := servers["other"]; !ok {
		t.Error("existing 'other' entry was removed")
	}
	if _, ok := servers["new-server"]; !ok {
		t.Error("new 'new-server' entry not found")
	}
}

func TestWriteEntryNoTempFileLeftBehind(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	if _, err := jsonconfig.WriteEntry(path, "mcpServers", "s", map[string]any{}, nil); err != nil {
		t.Fatalf("WriteEntry: %v", err)
	}

	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		name := e.Name()
		if name != "config.json" && name != "config.json.lock" {
			t.Errorf("unexpected file left behind: %s", name)
		}
	}
}

func TestWriteEntryConcurrent(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.json")
	const n = 20

	var wg sync.WaitGroup
	errs := make([]error, n)

	for i := range n {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			key := "server-" + string(rune('a'+idx))
			_, err := jsonconfig.WriteEntry(path, "mcpServers", key, map[string]any{"id": idx}, nil)
			errs[idx] = err
		}(i)
	}

	wg.Wait()

	for i, err := range errs {
		if err != nil {
			t.Errorf("goroutine %d: %v", i, err)
		}
	}

	root := readJSON(t, path)
	servers, ok := root["mcpServers"].(map[string]any)
	if !ok {
		t.Fatal("mcpServers missing or wrong type")
	}

	if len(servers) != n {
		t.Errorf("got %d entries, want %d", len(servers), n)
	}
}

func readJSON(t *testing.T, path string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile %s: %v", path, err)
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("Unmarshal %s: %v", path, err)
	}
	return m
}
