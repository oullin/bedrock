package support

import (
	"fmt"
	"html"
	"os"
	"reflect"
	"strings"
)

// Blank determines whether the given value is "blank".
// A value is blank if it is nil, an empty string, an empty slice,
// an empty map, or whitespace-only string.
func Blank(value any) bool {
	if value == nil {
		return true
	}

	if s, ok := value.(string); ok {
		return strings.TrimSpace(s) == ""
	}

	rv := reflect.ValueOf(value)

	switch rv.Kind() {
	case reflect.Slice, reflect.Map, reflect.Array:
		return rv.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return rv.IsNil()
	}

	return false
}

// Filled determines whether the given value is not "blank".
func Filled(value any) bool {
	return !Blank(value)
}

// Tap calls the given closure with the value, then returns the value.
func Tap[T any](value T, callbacks ...func(T)) T {
	for _, cb := range callbacks {
		cb(value)
	}

	return value
}

// Value returns the value of the given value. If the value is a
// function it will be called and its return value used.
func Value[T any](value any, args ...any) T {
	switch fn := value.(type) {
	case func() T:
		return fn()
	case func(any) T:
		if len(args) > 0 {
			return fn(args[0])
		}

		return fn(nil)
	case func(...any) T:
		return fn(args...)
	case T:
		return fn
	}

	// Try a plain any→T cast via reflection
	if v, ok := value.(T); ok {
		return v
	}

	var zero T

	return zero
}

// With returns the value, or calls the callback with the value if provided.
func With[T any](value T, fn ...func(T) T) T {
	if len(fn) > 0 && fn[0] != nil {
		return fn[0](value)
	}

	return value
}

// Transform transforms the given value if it is not blank.
// If the value is blank and a default is provided, the default is returned.
func Transform[T, U any](value T, fn func(T) U, def ...U) (U, bool) {
	if Blank(value) {
		if len(def) > 0 {
			return def[0], false
		}

		var zero U

		return zero, false
	}

	return fn(value), true
}

// E HTML-encodes the given string, converting special characters to HTML entities.
func E(value any) string {
	if value == nil {
		return ""
	}

	var raw string

	switch v := value.(type) {
	case string:
		raw = v
	case []byte:
		raw = string(v)
	default:
		raw = fmt.Sprint(v)
	}

	return html.EscapeString(raw)
}

// Env returns the value of the environment variable named by key.
// If the variable is not set or is empty, the optional default value is returned.
// The strings "true", "false", "null", and "empty" are converted to their
// Go equivalents: "true"→"true", "false"→"false", "null"→"", "empty"→"".
func Env(key string, def ...string) string {
	val, ok := os.LookupEnv(key)

	if !ok || val == "" {
		if len(def) > 0 {
			return def[0]
		}

		return ""
	}

	// Unquote escaped quoted strings (e.g. "\"hello\"" → "hello")
	if len(val) >= 2 && val[0] == '"' && val[len(val)-1] == '"' {
		val = val[1 : len(val)-1]
	}

	switch strings.ToLower(val) {
	case "true", "(true)":
		return "true"
	case "false", "(false)":
		return "false"
	case "null", "(null)":
		return ""
	case "empty", "(empty)":
		return ""
	}

	return val
}

// EnvBool returns the environment variable as a boolean.
// "true", "1", "yes", "on" → true; everything else → false.
func EnvBool(key string, def ...bool) bool {
	val, ok := os.LookupEnv(key)

	if !ok {
		if len(def) > 0 {
			return def[0]
		}

		return false
	}

	switch strings.ToLower(val) {
	case "true", "1", "yes", "on", "(true)":
		return true
	}

	return false
}

// Head returns the first item in a slice.
func Head[T any](items []T) (T, bool) {
	if len(items) == 0 {
		var zero T

		return zero, false
	}

	return items[0], true
}

// Last returns the last item in a slice.
func Last[T any](items []T) (T, bool) {
	if len(items) == 0 {
		var zero T

		return zero, false
	}

	return items[len(items)-1], true
}

// ClassBasename returns the unqualified name for a type or class-like string.
func ClassBasename(value any) string {
	if value == nil {
		return ""
	}

	if name, ok := value.(string); ok {
		name = strings.Trim(name, "\\/")

		if idx := strings.LastIndexAny(name, "\\/"); idx >= 0 {
			return name[idx+1:]
		}

		return name
	}

	typ := reflect.TypeOf(value)

	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	return typ.Name()
}
