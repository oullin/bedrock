package watchers

import "github.com/bedrock/packages/telescope"

// BatchWatcher records queued batch dispatches as Telescope batch entries.
type BatchWatcher struct {
	telescope.BaseWatcher
}

// BatchDispatch captures a queued batch dispatch event.
type BatchDispatch struct {
	ID    string
	Name  string
	Total int
}

// NewBatchWatcher creates a BatchWatcher with the given options.
func NewBatchWatcher(t *telescope.Telescope, options map[string]any) *BatchWatcher {
	w := &BatchWatcher{}
	w.SetTelescope(t)
	w.Options = options

	return w
}

// Register is a no-op for BatchWatcher; callers drive it via Record.
func (w *BatchWatcher) Register(_ any) error { return nil }

// Record stores a batch dispatch entry.
func (w *BatchWatcher) Record(batch BatchDispatch) {
	entry := telescope.NewEntry(telescope.EntryTypeBatch, map[string]any{
		"id":    batch.ID,
		"name":  batch.Name,
		"total": batch.Total,
	})

	w.Scope().Record(entry)
}
