package install_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/bedrock/packages/ai/boost/install"
)

type mockMcpAgent struct {
	configPath string
	configKey  string
}

func (m *mockMcpAgent) McpConfigPath() string { return m.configPath }
func (m *mockMcpAgent) McpConfigKey() string  { return m.configKey }

func TestMcpWriterCreatesFile(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	agent := &mockMcpAgent{
		configPath: filepath.Join(tmp, ".mcp.json"),
		configKey:  "mcpServers",
	}

	serverCfg := map[string]any{
		"command": "go",
		"args":    []string{"run", ".", "mcp"},
	}

	w := &install.McpWriter{}
	written, err := w.Write(agent, "boost", serverCfg)

	if err != nil {
		t.Fatalf("Write: %v", err)
	}

	if !written {
		t.Error("Write should return true when entry is new")
	}

	data, err := os.ReadFile(agent.configPath)

	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	var root map[string]any

	if err := json.Unmarshal(data, &root); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}

	servers, ok := root["mcpServers"].(map[string]any)

	if !ok {
		t.Fatalf("mcpServers not found or wrong type")
	}

	if _, exists := servers["boost"]; !exists {
		t.Error("boost entry not found in mcpServers")
	}
}

func TestMcpWriterIdempotent(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	agent := &mockMcpAgent{
		configPath: filepath.Join(tmp, ".mcp.json"),
		configKey:  "mcpServers",
	}

	serverCfg := map[string]any{"command": "go", "args": []string{"run", "."}}
	w := &install.McpWriter{}

	// Write once.
	if _, err := w.Write(agent, "boost", serverCfg); err != nil {
		t.Fatalf("first Write: %v", err)
	}

	// Write again — must return false (already exists).
	written, err := w.Write(agent, "boost", serverCfg)

	if err != nil {
		t.Fatalf("second Write: %v", err)
	}

	if written {
		t.Error("second Write should return false (idempotent)")
	}
}

func TestMcpWriterMergesExistingKeys(t *testing.T) {
	t.Parallel()

	tmp := t.TempDir()
	path := filepath.Join(tmp, ".mcp.json")

	// Pre-create the file with an existing entry.
	existing := map[string]any{
		"mcpServers": map[string]any{
			"other": map[string]any{"command": "other"},
		},
	}
	data, _ := json.MarshalIndent(existing, "", "    ")
	_ = os.WriteFile(path, data, 0644)

	agent := &mockMcpAgent{configPath: path, configKey: "mcpServers"}
	w := &install.McpWriter{}

	if _, err := w.Write(agent, "boost", map[string]any{"command": "go"}); err != nil {
		t.Fatalf("Write: %v", err)
	}

	raw, _ := os.ReadFile(path)

	var root map[string]any

	_ = json.Unmarshal(raw, &root)

	servers, _ := root["mcpServers"].(map[string]any)

	if _, ok := servers["other"]; !ok {
		t.Error("existing 'other' entry was removed")
	}

	if _, ok := servers["boost"]; !ok {
		t.Error("new 'boost' entry not found")
	}
}
