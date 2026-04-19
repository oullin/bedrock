package watchers

import (
	"github.com/bedrock/packages/debugbar"
)

// ViewWatcher monitors template/view rendering and records entries as DebugBar
// entries. It mirrors Upstream's ViewWatcher class.
type ViewWatcher struct {
	debugbar.BaseWatcher
}

// NewViewWatcher creates a ViewWatcher with the given options.
func NewViewWatcher(t *debugbar.DebugBar, options map[string]any) *ViewWatcher {
	w := &ViewWatcher{}
	w.SetDebugBar(t)
	w.Options = options

	return w
}

// Register is a no-op for ViewWatcher; callers drive it via Record.
func (w *ViewWatcher) Register(_ any) error { return nil }

// Record records a view-rendering entry. name is the template identifier
// (e.g. "users.index"), dataKeys lists the data variable names available to
// the template, and composers lists any composers/creators attached.
func (w *ViewWatcher) Record(name string, dataKeys []string, composers []string) {
	content := map[string]any{
		"name":      name,
		"data":      dataKeys,
		"composers": composers,
	}

	entry := debugbar.NewEntry(debugbar.EntryTypeView, content)

	w.Scope().RecordView(entry)
}
