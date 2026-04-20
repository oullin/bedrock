package arr

import "strings"

// Wrap wraps the given value in a single-element slice.
func Wrap[T any](value T) []T {
	return []T{value}
}

// WrapSlice returns the given slice unchanged.
// Use this to wrap a value that is already a slice without nesting it.
func WrapSlice[T any](value []T) []T {
	return value
}

// Join concatenates string slice elements with a glue string.
// Optional final glue is used between the last two elements.
func Join(items []string, glue string, finalGlues ...string) string {
	if len(items) == 0 {
		return ""
	}

	if len(items) == 1 {
		return items[0]
	}

	if len(finalGlues) > 0 && finalGlues[0] != "" {
		last := items[len(items)-1]
		rest := items[:len(items)-1]

		return strings.Join(rest, glue) + finalGlues[0] + last
	}

	return strings.Join(items, glue)
}
