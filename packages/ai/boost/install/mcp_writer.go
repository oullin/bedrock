package install

import (
	"github.com/bedrock/packages/ai/boost/internal/boosterr"
	"github.com/bedrock/packages/ai/boost/internal/jsonconfig"
)

// McpWriter upserts an MCP server entry into an agent's MCP config file.
// The file is expected to be a JSON object whose structure is:
//
//	{ "<configKey>": { "<serverKey>": { ... } } }
//
// If the file does not exist it is created. If the top-level key or the server
// entry already exists, the call is idempotent (returns false, nil).
// Mirrors upstream Boost\Install\McpWriter.
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
		return false, boosterr.ErrNoMcpConfigPath
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
	return jsonconfig.WriteEntry(configPath, configKey, serverKey, serverConfig, skeleton)
}
