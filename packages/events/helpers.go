package events

import (
	"reflect"
	"strings"
)

// eventName resolves an event to its string name.
// String values are returned as-is. Struct values use reflect.TypeOf().String().
func eventName(event any) string {
	if s, ok := event.(string); ok {
		return s
	}

	t := reflect.TypeOf(event)

	if t == nil {
		return ""
	}

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.String()
}

// isWildcardPattern reports whether the pattern contains a wildcard character.
func isWildcardPattern(pattern string) bool {
	return strings.Contains(pattern, "*")
}
