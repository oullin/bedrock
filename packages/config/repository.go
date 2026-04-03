package config

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// Repository stores configuration items using dotted-path access semantics.
type Repository struct {
	items map[string]any
}

// NewRepository creates a repository from the supplied item map.
func NewRepository(items map[string]any) *Repository {
	if items == nil {
		items = map[string]any{}
	}

	return &Repository{items: cloneMap(items)}
}

// Has reports whether a key exists.
func (r *Repository) Has(key string) bool {
	_, ok := lookup(r.items, key)

	return ok
}

// Get returns the value stored at key, falling back to defaultValue when absent.
func (r *Repository) Get(key string, defaultValue any) any {
	value, ok := lookup(r.items, key)

	if !ok {
		return defaultValue
	}

	return value
}

// GetMany returns a map of key lookups using the provided defaults.
func (r *Repository) GetMany(keys map[string]any) map[string]any {
	if len(keys) == 0 {
		return map[string]any{}
	}

	values := make(map[string]any, len(keys))

	for key, defaultValue := range keys {
		values[key] = r.Get(key, defaultValue)
	}

	return values
}

// Set writes one or more values into the repository.
func (r *Repository) Set(key any, value ...any) {
	switch typed := key.(type) {
	case string:
		var item any

		if len(value) > 0 {
			item = value[0]
		}

		r.setPath(typed, item)
	case map[string]any:
		for nestedKey, nestedValue := range typed {
			r.setPath(nestedKey, nestedValue)
		}
	}
}

// Prepend prepends a value to the slice stored at key.
func (r *Repository) Prepend(key string, value any) {
	items := append([]any{cloneValue(value)}, r.sliceValue(key)...)
	r.setPath(key, items)
}

// Push appends a value to the slice stored at key.
func (r *Repository) Push(key string, value any) {
	items := append(r.sliceValue(key), cloneValue(value))
	r.setPath(key, items)
}

// All returns a deep clone of all items.
func (r *Repository) All() map[string]any {
	return cloneMap(r.items)
}

// String returns a string value.
func (r *Repository) String(key string) (string, error) {
	value, ok := lookup(r.items, key)

	if !ok {
		return "", fmt.Errorf("config: missing key %q", key)
	}

	text, ok := value.(string)

	if !ok {
		return "", typeError(key, "string", value)
	}

	return text, nil
}

// Int returns an integer value.
func (r *Repository) Int(key string) (int, error) {
	value, ok := lookup(r.items, key)

	if !ok {
		return 0, fmt.Errorf("config: missing key %q", key)
	}

	switch typed := value.(type) {
	case int:
		return typed, nil
	case int64:
		return int(typed), nil
	case float64:
		return int(typed), nil
	case string:
		parsed, err := strconv.Atoi(strings.TrimSpace(typed))

		if err != nil {
			return 0, typeError(key, "int", value)
		}

		return parsed, nil
	default:
		return 0, typeError(key, "int", value)
	}
}

// Bool returns a boolean value.
func (r *Repository) Bool(key string) (bool, error) {
	value, ok := lookup(r.items, key)

	if !ok {
		return false, fmt.Errorf("config: missing key %q", key)
	}

	switch typed := value.(type) {
	case bool:
		return typed, nil
	case string:
		parsed, err := strconv.ParseBool(strings.TrimSpace(typed))

		if err != nil {
			return false, typeError(key, "bool", value)
		}

		return parsed, nil
	default:
		return false, typeError(key, "bool", value)
	}
}

// Duration returns a parsed Go duration.
func (r *Repository) Duration(key string) (time.Duration, error) {
	value, ok := lookup(r.items, key)

	if !ok {
		return 0, fmt.Errorf("config: missing key %q", key)
	}

	switch typed := value.(type) {
	case time.Duration:
		return typed, nil
	case string:
		duration, err := time.ParseDuration(strings.TrimSpace(typed))

		if err != nil {
			return 0, typeError(key, "duration", value)
		}

		return duration, nil
	default:
		return 0, typeError(key, "duration", value)
	}
}

// StringSlice returns a string slice.
func (r *Repository) StringSlice(key string) ([]string, error) {
	value, ok := lookup(r.items, key)

	if !ok {
		return nil, fmt.Errorf("config: missing key %q", key)
	}

	switch typed := value.(type) {
	case []string:
		return append([]string(nil), typed...), nil
	case []any:
		out := make([]string, 0, len(typed))

		for _, item := range typed {
			text, ok := item.(string)

			if !ok {
				return nil, typeError(key, "[]string", value)
			}

			out = append(out, text)
		}

		return out, nil
	case string:
		if strings.TrimSpace(typed) == "" {
			return []string{}, nil
		}

		parts := strings.Split(typed, ",")
		out := make([]string, 0, len(parts))

		for _, part := range parts {
			out = append(out, strings.TrimSpace(part))
		}

		return out, nil
	default:
		return nil, typeError(key, "[]string", value)
	}
}

// Map returns a cloned map value.
func (r *Repository) Map(key string) (map[string]any, error) {
	value, ok := lookup(r.items, key)

	if !ok {
		return nil, fmt.Errorf("config: missing key %q", key)
	}

	typed, ok := value.(map[string]any)

	if !ok {
		return nil, typeError(key, "map[string]any", value)
	}

	return cloneMap(typed), nil
}

func lookup(items map[string]any, key string) (any, bool) {
	if value, ok := items[key]; ok {
		return cloneValue(value), true
	}

	segments := splitKey(key)

	if len(segments) == 0 {
		return items, true
	}

	current := any(items)

	for index, segment := range segments {
		mapped, ok := current.(map[string]any)

		if !ok {
			return nil, false
		}

		remaining := strings.Join(segments[index:], ".")

		if value, ok := mapped[remaining]; ok {
			return cloneValue(value), true
		}

		next, ok := mapped[segment]

		if !ok {
			return nil, false
		}

		current = next
	}

	return cloneValue(current), true
}

func splitKey(key string) []string {
	key = strings.TrimSpace(key)

	if key == "" {
		return nil
	}

	return strings.Split(key, ".")
}

func typeError(key string, want string, value any) error {
	return fmt.Errorf("config: key %q must be %s, got %T", key, want, value)
}

func (r *Repository) setPath(key string, value any) {
	segments := splitKey(key)

	if len(segments) == 0 {
		return
	}

	current := r.items

	for _, segment := range segments[:len(segments)-1] {
		next, ok := current[segment]

		if !ok {
			child := map[string]any{}
			current[segment] = child
			current = child

			continue
		}

		child, ok := next.(map[string]any)

		if !ok {
			child = map[string]any{}
			current[segment] = child
		}

		current = child
	}

	current[segments[len(segments)-1]] = cloneValue(value)
}

func (r *Repository) sliceValue(key string) []any {
	value, ok := lookup(r.items, key)

	if !ok {
		return []any{}
	}

	switch typed := value.(type) {
	case []any:
		return append([]any(nil), typed...)
	case []string:
		out := make([]any, 0, len(typed))

		for _, item := range typed {
			out = append(out, item)
		}

		return out
	default:
		return []any{}
	}
}

func cloneMap(items map[string]any) map[string]any {
	cloned := make(map[string]any, len(items))

	for key, value := range items {
		cloned[key] = cloneValue(value)
	}

	return cloned
}

func cloneValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		return cloneMap(typed)
	case []any:
		cloned := make([]any, len(typed))

		for index, item := range typed {
			cloned[index] = cloneValue(item)
		}

		return cloned
	case []string:
		return append([]string(nil), typed...)
	default:
		return typed
	}
}
