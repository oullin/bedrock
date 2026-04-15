package install

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// McpWriter upserts an MCP server entry into an agent's MCP config file.
// The file is expected to be a JSON object whose structure is:
//
//	{ "<configKey>": { "<serverKey>": { ... } } }
//
// If the file does not exist it is created. If the top-level key or the server
// entry already exists, the call is idempotent (returns false, nil).
// Mirrors Upstream\Boost\Install\McpWriter.
type McpWriter struct{}

// SupportsMcpConfigPath is the minimal interface required by McpWriter.
type SupportsMcpConfigPath interface {
	McpConfigPath() string
	McpConfigKey() string
}

// Write merges serverConfig under configPath[configKey][serverKey].
// Returns true if the file was modified, false if the entry already existed.
func (w *McpWriter) Write(agent SupportsMcpConfigPath, serverKey string, serverConfig map[string]any) (bool, error) {
	configPath := agent.McpConfigPath()
	configKey := agent.McpConfigKey()

	if configPath == "" {
		return false, fmt.Errorf("install: agent returned an empty MCP config path")
	}

	return writeJSONConfigEntry(configPath, configKey, serverKey, serverConfig, nil)
}

// writeJSONConfigEntry is the shared idempotent JSON-merge helper used by both
// McpWriter and agents.BaseAgent.InstallMcp.
func writeJSONConfigEntry(
	configPath, configKey, serverKey string,
	serverConfig map[string]any,
	skeleton map[string]any,
) (bool, error) {
	// Read existing file (or start from skeleton / empty object).
	var root map[string]any

	data, err := os.ReadFile(configPath)
	if err != nil && !os.IsNotExist(err) {
		return false, fmt.Errorf("mcp_writer: read %s: %w", configPath, err)
	}

	if len(data) > 0 {
		if err := json.Unmarshal(data, &root); err != nil {
			return false, fmt.Errorf("mcp_writer: unmarshal %s: %w", configPath, err)
		}
	}

	if root == nil {
		if skeleton != nil {
			root = skeleton
		} else {
			root = map[string]any{}
		}
	}

	// Ensure the top-level config key exists.
	if _, ok := root[configKey]; !ok {
		root[configKey] = map[string]any{}
	}

	servers, ok := root[configKey].(map[string]any)
	if !ok {
		// Unexpected type; reset.
		servers = map[string]any{}
		root[configKey] = servers
	}

	// Idempotent: if the server key already exists don't overwrite.
	if _, exists := servers[serverKey]; exists {
		return false, nil
	}

	servers[serverKey] = serverConfig

	// Write back.
	if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
		return false, fmt.Errorf("mcp_writer: mkdir %s: %w", filepath.Dir(configPath), err)
	}

	out, err := json.MarshalIndent(root, "", "    ")
	if err != nil {
		return false, fmt.Errorf("mcp_writer: marshal: %w", err)
	}

	if err := os.WriteFile(configPath, append(out, '\n'), 0644); err != nil {
		return false, fmt.Errorf("mcp_writer: write %s: %w", configPath, err)
	}

	return true, nil
}
