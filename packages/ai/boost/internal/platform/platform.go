// Package platform defines shared enums used by coding agents.
package platform

import "runtime"

// Platform identifies the host operating system.
type Platform int

// Darwin represents macOS.

// Linux represents any Linux distribution.

// Windows represents Microsoft Windows.

// Current returns the Platform matching the current OS.

// String returns the canonical lowercase name.

// McpInstallationStrategy determines how an agent writes its MCP server config.
type McpInstallationStrategy int

const (
	Darwin Platform = iota

	Linux

	Windows
)

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

const (
	// McpStrategyFile writes a JSON (or TOML) config file.
	McpStrategyFile McpInstallationStrategy = iota
	// McpStrategyShell runs a shell command to install the MCP server.
	McpStrategyShell
	// McpStrategyNone indicates no MCP support for this agent.
	McpStrategyNone
)
