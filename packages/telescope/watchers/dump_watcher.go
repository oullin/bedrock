package watchers

import (
	"fmt"

	"github.com/bedrock/packages/telescope"
)

// DumpWatcher monitors debug dump() calls and records them as Telescope
// entries. It mirrors Laravel's DumpWatcher class.
type DumpWatcher struct {
	telescope.BaseWatcher
}

// NewDumpWatcher creates a DumpWatcher with the given options.
func NewDumpWatcher(t *telescope.Telescope, options map[string]any) *DumpWatcher {
	w := &DumpWatcher{}
	w.SetTelescope(t)
	w.Options = options

	return w
}

// Register is a no-op for DumpWatcher; callers drive it via Record.
func (w *DumpWatcher) Register(_ any) error { return nil }

// Record records a debug dump entry. value is the value being dumped; it is
// formatted using fmt.Sprintf("%#v", ...) to produce a Go-syntax
// representation mirroring PHP's var_dump output.
func (w *DumpWatcher) Record(value any) {
	content := map[string]any{
		"dump": fmt.Sprintf("%#v", value),
	}

	entry := telescope.NewEntry(telescope.EntryTypeDump, content)

	w.Scope().RecordDump(entry)
}
