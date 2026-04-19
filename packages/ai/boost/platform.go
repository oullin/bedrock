package boost

import "github.com/bedrock/packages/ai/boost/internal/platform"

const (
	// PlatformDarwin represents macOS.
	PlatformDarwin = platform.Darwin
	// PlatformLinux represents any Linux distribution.
	PlatformLinux = platform.Linux
	// PlatformWindows represents Microsoft Windows.
	PlatformWindows = platform.Windows
)

// CurrentPlatform returns the Platform for the current OS.
func CurrentPlatform() Platform { return platform.Current() }
