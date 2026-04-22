package boosterr

import "errors"

// ErrNoMcpConfigPath is returned when an agent's MCP config path is empty.
var ErrNoMcpConfigPath = errors.New("boost: agent has no MCP config path")
