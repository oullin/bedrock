package watchers

import (
	"github.com/bedrock/packages/debugbar"
)

// ignoredCommands lists CLI commands that should not be recorded by default,
// mirroring Upstream's CommandWatcher ignore list.
var ignoredCommands = []string{
	"schedule:run",
	"schedule:finish",
	"package:discover",
}

// CommandWatcher monitors CLI command execution and records entries as
// DebugBar entries. It mirrors Upstream's CommandWatcher class.
//
// Options:
//   - "ignore" ([]string): additional command names to skip.
type CommandWatcher struct {
	debugbar.BaseWatcher
}

// NewCommandWatcher creates a CommandWatcher with the given options.
func NewCommandWatcher(t *debugbar.DebugBar, options map[string]any) *CommandWatcher {
	w := &CommandWatcher{}
	w.SetDebugBar(t)
	w.Options = options

	return w
}

// Register is a no-op for CommandWatcher; callers drive it via Record.
func (w *CommandWatcher) Register(_ any) error { return nil }

// ShouldIgnore reports whether the command should be skipped.
func (w *CommandWatcher) ShouldIgnore(name string) bool {
	for _, cmd := range ignoredCommands {
		if cmd == name {
			return true
		}
	}

	for _, cmd := range w.StringsOption("ignore") {
		if cmd == name {
			return true
		}
	}

	return false
}

// Record records a CLI command execution entry. name is the command name,
// exitCode is the process exit code (0 = success), arguments and flags are
// the parsed command-line inputs.
func (w *CommandWatcher) Record(name string, exitCode int, arguments, flags map[string]any) {
	if w.ShouldIgnore(name) {
		return
	}

	content := map[string]any{
		"command":   name,
		"exit_code": exitCode,
		"arguments": arguments,
		"options":   flags,
	}

	entry := debugbar.NewEntry(debugbar.EntryTypeCommand, content)

	w.Scope().RecordCommand(entry)
}
