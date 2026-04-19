package echo

import "strings"

// EventFormatter normalizes event names for the Echo broadcaster.
// It optionally prepends a namespace and converts dot separators to backslashes.
type EventFormatter struct {
	namespace string
}

// NewEventFormatter creates an EventFormatter with the given namespace.
// Pass an empty string to disable namespace prepending, which is equivalent
// to namespace: false in the TypeScript library.
func NewEventFormatter(namespace string) *EventFormatter {
	return &EventFormatter{namespace: namespace}
}

// Format normalizes an event name according to Echo's rules:
//
//  1. If the event starts with '.' or '\\': strip the first character and
//     return the remainder unchanged — no namespace is prepended and dots
//     are not replaced.
//
//  2. Otherwise: if namespace is non-empty, prepend "namespace."; then
//     replace every '.' with '\\' in the entire result.
//
// Examples with namespace "App.Events":
//
//	"Users.UserCreated"         → "App\\Events\\Users\\UserCreated"
//	".App\\Users\\UserCreated"  → "App\\Users\\UserCreated"
//	"\\App\\Users\\UserCreated" → "App\\Users\\UserCreated"
//
// Example with namespace "" (disabled):
//
//	"Users.UserCreated"         → "Users\\UserCreated"
func (f *EventFormatter) Format(event string) string {
	if len(event) == 0 {
		return event
	}

	first := event[0]

	if first == '.' || first == '\\' {
		return event[1:]
	}

	full := event

	if f.namespace != "" {
		full = f.namespace + "." + event
	}

	return strings.ReplaceAll(full, ".", "\\")
}

// SetNamespace updates the namespace used for event formatting.
// Pass an empty string to disable namespace prepending.
func (f *EventFormatter) SetNamespace(namespace string) {
	f.namespace = namespace
}
