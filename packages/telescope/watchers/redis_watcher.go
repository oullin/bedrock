package watchers

import (
	"strings"
	"time"

	"github.com/bedrock/packages/telescope"
)

// ignoredRedisCommands lists commands that should not be recorded (pipeline /
// transaction control), mirroring Laravel's RedisWatcher.

// RedisWatcher monitors Redis command execution and records entries as
// Telescope entries. It mirrors Laravel's RedisWatcher class.
type RedisWatcher struct {
	telescope.BaseWatcher
}

var ignoredRedisCommands = []string{
	"MULTI",
	"EXEC",
	"DISCARD",
	"PIPELINE",
}

// NewRedisWatcher creates a RedisWatcher with the given options.
func NewRedisWatcher(t *telescope.Telescope, options map[string]any) *RedisWatcher {
	w := &RedisWatcher{}
	w.SetTelescope(t)
	w.Options = options

	return w
}

// Register is a no-op for RedisWatcher; callers drive it via Record.
func (w *RedisWatcher) Register(_ any) error { return nil }

// ShouldIgnore reports whether the Redis command should be skipped.
func (w *RedisWatcher) ShouldIgnore(command string) bool {
	upper := strings.ToUpper(command)

	for _, cmd := range ignoredRedisCommands {
		if cmd == upper {
			return true
		}
	}

	return false
}

// Record records a Redis command execution entry. command is the Redis
// command (e.g. "SET"), connection is the connection name, and duration is
// the command execution time.
func (w *RedisWatcher) Record(command, connection string, duration time.Duration) {
	if w.ShouldIgnore(command) {
		return
	}

	content := map[string]any{
		"command":    command,
		"connection": connection,
		"time":       float64(duration.Microseconds()) / 1000.0,
	}

	entry := telescope.NewEntry(telescope.EntryTypeRedis, content)

	w.Scope().RecordRedis(entry)
}
