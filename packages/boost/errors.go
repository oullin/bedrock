package boost

import "errors"

var (
	// ErrAgentAlreadyRegistered is returned by RegisterAgent when the given key
	// is already present in the manager.
	ErrAgentAlreadyRegistered = errors.New("boost: agent key already registered")

	// ErrNoMcpConfigPath is returned when an agent's McpConfigPath is empty.
	ErrNoMcpConfigPath = errors.New("boost: agent has no MCP config path")

	// ErrMcpInstallFailed is returned when an MCP installation attempt fails.
	ErrMcpInstallFailed = errors.New("boost: MCP installation failed")
)
