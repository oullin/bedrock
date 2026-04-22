package boost

import (
	"errors"

	"github.com/bedrock/packages/ai/boost/internal/boosterr"
)

var (
	// ErrAgentAlreadyRegistered is returned by RegisterAgent when the given key
	// is already present in the manager.
	ErrAgentAlreadyRegistered = errors.New("boost: agent key already registered")

	// ErrNoMcpConfigPath is returned when an agent's McpConfigPath is empty.
	ErrNoMcpConfigPath = boosterr.ErrNoMcpConfigPath

	// ErrMcpInstallFailed is returned when an MCP installation attempt fails.
	ErrMcpInstallFailed = errors.New("boost: MCP installation failed")
)
