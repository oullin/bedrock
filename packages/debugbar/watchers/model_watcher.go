package watchers

import (
	"github.com/bedrock/packages/debugbar"
)

// ModelAction constants mirror the Orm lifecycle events recorded by
// Upstream's ModelWatcher.

// ModelWatcher monitors model lifecycle events and records them as DebugBar
// entries. It mirrors Upstream's ModelWatcher class.
//
// Options:
//   - "events" ([]string): specific model actions to record (default: all).
//   - "ignore" ([]string): fully-qualified model type names to skip.
//   - "hydrations" (bool): whether to record retrieved (hydration) events.
type ModelWatcher struct {
	debugbar.BaseWatcher
}

const (
	ModelActionCreated   = "created"
	ModelActionUpdated   = "updated"
	ModelActionDeleted   = "deleted"
	ModelActionRestored  = "restored"
	ModelActionRetrieved = "retrieved"
)

// NewModelWatcher creates a ModelWatcher with the given options.
func NewModelWatcher(t *debugbar.DebugBar, options map[string]any) *ModelWatcher {
	w := &ModelWatcher{}
	w.SetDebugBar(t)
	w.Options = options

	return w
}

// Register is a no-op for ModelWatcher; callers drive it via Record.
func (w *ModelWatcher) Register(_ any) error { return nil }

// ShouldIgnore reports whether the model type should be skipped.
func (w *ModelWatcher) ShouldIgnore(modelType string) bool {
	for _, name := range w.StringsOption("ignore") {
		if name == modelType {
			return true
		}
	}

	return false
}

// Record records a model lifecycle event. modelType is the Go type name,
// action is one of the ModelAction* constants, and key is the primary key
// value for the model instance.
func (w *ModelWatcher) Record(modelType, action string, key any) {
	if w.ShouldIgnore(modelType) {
		return
	}

	// Skip hydration events unless explicitly requested.
	if action == ModelActionRetrieved && !w.BoolOption("hydrations") {
		return
	}

	content := map[string]any{
		"model":  modelType,
		"action": action,
		"key":    key,
	}

	entry := debugbar.NewEntry(debugbar.EntryTypeModel, content)
	entry.AddTags(modelType)

	w.Scope().RecordModel(entry)
}
