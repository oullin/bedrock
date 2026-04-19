// Package log is the facade for the log manager. It forwards common
// operations to the LogManager bound under "log" in the global bedrock
// Application.
package log

import (
	"sync"

	"github.com/bedrock/packages/bedrock"
	logpkg "github.com/bedrock/packages/log"
)

var (
	mu     sync.Mutex
	cached *logpkg.LogManager
)

// Manager returns the log manager from the global Application. Resolved
// once per process and cached.
func Manager() *logpkg.LogManager {
	mu.Lock()

	defer mu.Unlock()

	if cached == nil {
		cached = bedrock.Resolve[*logpkg.LogManager]("log")
	}

	return cached
}

// Reset clears the cached manager. Tests must call this after reinstalling
// a different Application via bedrock.SetApp.
func Reset() {
	mu.Lock()

	defer mu.Unlock()

	cached = nil
}

// Channel returns a named channel logger. Pass nothing for the default channel.
func Channel(name ...string) (*logpkg.Logger, error) {
	return Manager().Channel(name...)
}

// Info logs to the default channel at info level. Errors resolving the
// channel are silently dropped — call Channel() directly if you need
// error handling.
func Info(message string, context ...map[string]any) {
	if l, err := Manager().Channel(); err == nil {
		l.Info(message, context...)
	}
}

// Error logs to the default channel at error level.
func Error(message string, context ...map[string]any) {
	if l, err := Manager().Channel(); err == nil {
		l.Error(message, context...)
	}
}

// Warning logs to the default channel at warning level.
func Warning(message string, context ...map[string]any) {
	if l, err := Manager().Channel(); err == nil {
		l.Warning(message, context...)
	}
}

// Debug logs to the default channel at debug level.
func Debug(message string, context ...map[string]any) {
	if l, err := Manager().Channel(); err == nil {
		l.Debug(message, context...)
	}
}
