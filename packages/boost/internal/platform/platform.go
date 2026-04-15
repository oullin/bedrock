// Package platform defines shared enums used by coding agents.
package platform

import "runtime"

// Platform identifies the host operating system.
// Mirrors Laravel\Boost\Install\Enums\Platform.
type Platform int

const (
	// Darwin represents macOS.
	Darwin Platform = iota
	// Linux represents any Linux distribution.
	Linux
	// Windows represents Microsoft Windows.
	Windows
)

// Current returns the Platform matching the current OS.
func Current() Platform {
	switch runtime.GOOS {
	case "darwin":
		return Darwin
	case "windows":
		return Windows
	default:
		return Linux
	}
}

// String returns the canonical lowercase name.
func (p Platform) String() string {
	switch p {
	case Darwin:
		return "darwin"
	case Windows:
		return "windows"
	default:
		return "linux"
	}
}

// McpInstallationStrategy determines how an agent writes its MCP server config.
// Mirrors Laravel\Boost\Install\Enums\McpInstallationStrategy.
type McpInstallationStrategy int

const (
	// McpStrategyFile writes a JSON (or TOML) config file.
	McpStrategyFile McpInstallationStrategy = iota
	// McpStrategyShell runs a shell command to install the MCP server.
	McpStrategyShell
	// McpStrategyNone indicates no MCP support for this agent.
	McpStrategyNone
)
