package watchers

import (
	"github.com/bedrock/packages/telescope"
)

// ViewWatcher monitors template/view rendering and records entries as Telescope
// entries. It mirrors Laravel's ViewWatcher class.
type ViewWatcher struct {
	telescope.BaseWatcher
}

// NewViewWatcher creates a ViewWatcher with the given options.
func NewViewWatcher(t *telescope.Telescope, options map[string]any) *ViewWatcher {
	w := &ViewWatcher{}
	w.SetTelescope(t)
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

	entry := telescope.NewEntry(telescope.EntryTypeView, content)

	w.Scope().RecordView(entry)
}
