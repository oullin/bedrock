package watchers

import (
	"reflect"
	"strings"

	"github.com/bedrock/packages/telescope"
)

// ignoredEventPrefixes contains event name prefixes from the framework that
// should not be recorded, mirroring the upstream EventWatcher ignore list.

// EventWatcher monitors application events and records them as Telescope
// entries. It mirrors the upstream EventWatcher class.
//
// Options:
//   - "ignore" ([]string): additional event names or prefixes to ignore.
type EventWatcher struct {
	telescope.BaseWatcher
}

var ignoredEventPrefixes = []string{
	"github.com/bedrock/packages/telescope",
	"@bedrock\\",
	"Octane\\",
}

// NewEventWatcher creates an EventWatcher with the given options.
func NewEventWatcher(t *telescope.Telescope, options map[string]any) *EventWatcher {
	w := &EventWatcher{}
	w.SetTelescope(t)
	w.Options = options

	return w
}

// Register is a no-op for EventWatcher; callers drive it by calling Record
// directly. When integrated with the bedrock events package, pass all events
// to Record as they fire.
func (w *EventWatcher) Register(_ any) error { return nil }

// ShouldIgnore reports whether the event should be skipped, mirroring
// EventWatcher::shouldIgnore().
func (w *EventWatcher) ShouldIgnore(eventName string) bool {
	for _, prefix := range ignoredEventPrefixes {
		if strings.HasPrefix(eventName, prefix) {
			return true
		}
	}

	if ignored := w.StringsOption("ignore"); len(ignored) > 0 {
		for _, name := range ignored {
			if eventName == name || strings.HasPrefix(eventName, name) {
				return true
			}
		}
	}

	return false
}

// Record records an application event entry. payload is the event data.
// listenerNames lists the registered listeners for this event.
func (w *EventWatcher) Record(eventName string, payload any, listenerNames []string) {
	if w.ShouldIgnore(eventName) {
		return
	}

	content := map[string]any{
		"name":      eventName,
		"payload":   extractPayload(payload),
		"listeners": listenerNames,
	}

	entry := telescope.NewEntry(telescope.EntryTypeEvent, content)

	w.Scope().RecordEvent(entry)
}

// extractPayload converts the event payload to a JSON-safe representation,
// mirroring EventWatcher::extractPayload().
func extractPayload(payload any) any {
	if payload == nil {
		return nil
	}

	switch v := payload.(type) {
	case map[string]any:
		return v
	case string, bool, int, int64, float64:
		return payload
	default:
		return structToMap(v)
	}
}

// structToMap converts an arbitrary value to a map[string]any using
// reflection. Exported fields only; zero-value fields are included.
func structToMap(v any) map[string]any {
	if v == nil {
		return nil
	}

	rv := reflect.ValueOf(v)

	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}

		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return map[string]any{"value": v}
	}

	rt := rv.Type()
	out := make(map[string]any, rt.NumField())

	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)

		if !field.IsExported() {
			continue
		}

		out[field.Name] = rv.Field(i).Interface()
	}

	return out
}
