package watchers

import (
	"fmt"

	"github.com/bedrock/packages/debugbar"
)

// DumpWatcher monitors debug dump() calls and records them as DebugBar
// entries.
type DumpWatcher struct {
	debugbar.BaseWatcher
}

// NewDumpWatcher creates a DumpWatcher with the given options.
func NewDumpWatcher(t *debugbar.DebugBar, options map[string]any) *DumpWatcher {
	w := &DumpWatcher{}
	w.SetDebugBar(t)
	w.Options = options

	return w
}

// Register is a no-op for DumpWatcher; callers drive it via Record.
func (w *DumpWatcher) Register(_ any) error { return nil }

// Record records a debug dump entry. value is the value being dumped; it is
// formatted using fmt.Sprintf("%#v", ...) to produce a Go-syntax
func (w *DumpWatcher) Record(value any) {
	content := map[string]any{
		"dump": fmt.Sprintf("%#v", value),
	}

	entry := debugbar.NewEntry(debugbar.EntryTypeDump, content)

	w.Scope().RecordDump(entry)
}
