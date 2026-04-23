package config

import (
	"fmt"
	"strconv"
	"strings"

	collection "github.com/bedrock/packages/collection/collection"
	"github.com/spf13/viper"
)

// Repository is a Laravel-inspired configuration store backed by Viper. It
// provides dot-notation key access, type-safe getters, and array manipulation
// helpers. The underlying Viper instance can be configured directly for YAML
// file loading and environment variable binding.
//
// Repository is not safe for concurrent use. If concurrent access is needed,
// callers must synchronise externally.
type Repository struct {
	v     *viper.Viper
	items map[string]any
}

// New creates a Repository pre-loaded with the given key-value pairs. The map
// may contain nested map[string]any values accessible via dot notation.
//
// Items are stored verbatim in the repository and are not pushed into the
// underlying Viper instance: doing so would split dot-notation keys (e.g.
// "a.b") into nested paths, which — combined with Go's randomised map
// iteration order — can leave Viper with a non-deterministic view when the
// input map mixes literal dotted keys with nested maps under the same prefix.
// Lookups consult items first via lookupExplicit, so callers see the
// unaltered structure they provided.
func New(items map[string]any) *Repository {
	return &Repository{v: viper.New(), items: cloneMap(items)}
}

// NewFromViper wraps an already-configured Viper instance. Use this when you
// have set up config file paths, environment prefixes, or other Viper options
// before creating the repository.
func NewFromViper(v *viper.Viper) *Repository {
	return &Repository{v: v, items: map[string]any{}}
}

// NewWithDefaults creates a Repository with the given key-value pairs registered
// as defaults. File values and environment variables override these defaults.
func NewWithDefaults(defaults map[string]any) *Repository {
	v := viper.New()

	for key, value := range defaults {
		v.SetDefault(key, value)
	}

	return &Repository{v: v, items: map[string]any{}}
}

// Viper returns the underlying Viper instance so consumers can configure file
// paths, environment prefixes, key replacers, and call ReadInConfig directly.
func (r *Repository) Viper() *viper.Viper {
	return r.v
}

// Has reports whether the given key is set in any configuration source.
func (r *Repository) Has(key string) bool {
	if _, ok := r.lookupExplicit(key); ok {
		return true
	}

	return r.v.IsSet(key)
}

// Get returns the value for key. If the key is not set, the first fallback
// value is returned (or nil when no fallback is provided).
func (r *Repository) Get(key string, fallback ...any) any {
	if value, ok := r.lookupExplicit(key); ok {
		return value
	}

	if r.v.IsSet(key) {
		return r.v.Get(key)
	}

	if len(fallback) > 0 {
		return fallback[0]
	}

	return nil
}

// GetMany returns a map of values for the given keys. Keys that are not set map
// to nil unless an optional per-key fallback map is provided.
func (r *Repository) GetMany(keys []string, defaults ...map[string]any) map[string]any {
	result := make(map[string]any, len(keys))

	for _, key := range keys {
		if !r.Has(key) && len(defaults) > 0 {
			if fallback, ok := defaults[0][key]; ok {
				result[key] = fallback

				continue
			}
		}

		result[key] = r.Get(key)
	}

	return result
}

// Set stores a value at the given dot-notation key.
func (r *Repository) Set(key string, value any) {
	if r.items == nil {
		r.items = map[string]any{}
	}

	if _, ok := r.items[key]; ok || !strings.Contains(key, ".") {
		r.items[key] = value
		r.v.Set(key, value)

		return
	}

	setDot(r.items, key, value)
	r.v.Set(key, value)
}

// SetMany stores multiple key-value pairs at once.
func (r *Repository) SetMany(values map[string]any) {
	for key, value := range values {
		r.Set(key, value)
	}
}

// Unset marks a key as explicitly present with a nil value. This mirrors
// Laravel's repository offset unset behavior while preserving Go's explicit
// method surface.
func (r *Repository) Unset(key string) {
	r.Set(key, nil)
}

// Prepend inserts value at the beginning of the slice stored at key. If the
// key does not exist, a new single-element slice is created.
func (r *Repository) Prepend(key string, value any) {
	existing := r.Get(key)

	switch v := existing.(type) {
	case []any:
		r.Set(key, append([]any{value}, v...))
	default:
		r.Set(key, []any{value})
	}
}

// Push appends value to the end of the slice stored at key. If the key does
// not exist, a new single-element slice is created.
func (r *Repository) Push(key string, value any) {
	existing := r.Get(key)

	switch v := existing.(type) {
	case []any:
		r.Set(key, append(v, value))
	default:
		r.Set(key, []any{value})
	}
}

// All returns every configuration item as a nested map.
func (r *Repository) All() map[string]any {
	all := cloneMap(r.v.AllSettings())
	mergeMap(all, r.items)

	return all
}

// String returns the string value for key. If the key is not set, the first
// fallback is returned. An error wrapping ErrInvalidType is returned when the
// stored value is not a string.
func (r *Repository) String(key string, fallback ...string) (string, error) {
	value, ok := r.lookupExplicit(key)

	if !ok && r.v.IsSet(key) {
		value = r.v.Get(key)
		ok = true
	}

	if !ok {
		if len(fallback) > 0 {
			return fallback[0], nil
		}

		return "", nil
	}

	s, ok := value.(string)

	if !ok {
		return "", fmt.Errorf("%w: value for key [%s] must be a string, %T given", ErrInvalidType, key, value)
	}

	return s, nil
}

// Integer returns the int value for key. If the key is not set, the first
// fallback is returned. An error wrapping ErrInvalidType is returned when the
// stored value is not an int.
func (r *Repository) Integer(key string, fallback ...int) (int, error) {
	value, ok := r.lookupExplicit(key)

	if !ok && r.v.IsSet(key) {
		value = r.v.Get(key)
		ok = true
	}

	if !ok {
		if len(fallback) > 0 {
			return fallback[0], nil
		}

		return 0, nil
	}

	i, ok := value.(int)

	if !ok {
		return 0, fmt.Errorf("%w: value for key [%s] must be an int, %T given", ErrInvalidType, key, value)
	}

	return i, nil
}

// Float returns the float64 value for key. If the key is not set, the first
// fallback is returned. An error wrapping ErrInvalidType is returned when the
// stored value is not a float64.
func (r *Repository) Float(key string, fallback ...float64) (float64, error) {
	value, ok := r.lookupExplicit(key)

	if !ok && r.v.IsSet(key) {
		value = r.v.Get(key)
		ok = true
	}

	if !ok {
		if len(fallback) > 0 {
			return fallback[0], nil
		}

		return 0, nil
	}

	f, ok := value.(float64)

	if !ok {
		return 0, fmt.Errorf("%w: value for key [%s] must be a float64, %T given", ErrInvalidType, key, value)
	}

	return f, nil
}

// Boolean returns the bool value for key. If the key is not set, the first
// fallback is returned. An error wrapping ErrInvalidType is returned when the
// stored value is not a bool.
func (r *Repository) Boolean(key string, fallback ...bool) (bool, error) {
	value, ok := r.lookupExplicit(key)

	if !ok && r.v.IsSet(key) {
		value = r.v.Get(key)
		ok = true
	}

	if !ok {
		if len(fallback) > 0 {
			return fallback[0], nil
		}

		return false, nil
	}

	b, ok := value.(bool)

	if !ok {
		return false, fmt.Errorf("%w: value for key [%s] must be a bool, %T given", ErrInvalidType, key, value)
	}

	return b, nil
}

// Array returns the []any value for key. If the key is not set, the first
// fallback is returned. An error wrapping ErrInvalidType is returned when the
// stored value is not a []any.
func (r *Repository) Array(key string, fallback ...[]any) ([]any, error) {
	value, ok := r.lookupExplicit(key)

	if !ok && r.v.IsSet(key) {
		value = r.v.Get(key)
		ok = true
	}

	if !ok {
		if len(fallback) > 0 {
			return fallback[0], nil
		}

		return nil, nil
	}

	a, ok := value.([]any)

	if !ok {
		return nil, fmt.Errorf("%w: value for key [%s] must be a []any, %T given", ErrInvalidType, key, value)
	}

	return a, nil
}

// Collection returns the value for key wrapped in Bedrock's slice collection.
// An error wrapping ErrInvalidType is returned when the stored value is not a
// []any.
func (r *Repository) Collection(key string, fallback ...[]any) (*collection.Collection[any], error) {
	items, err := r.Array(key, fallback...)

	if err != nil {
		return nil, err
	}

	return collection.Collect(items), nil
}

func (r *Repository) lookupExplicit(key string) (any, bool) {
	if r.items == nil {
		return nil, false
	}

	if value, ok := r.items[key]; ok {
		return value, true
	}

	return lookupDot(r.items, key)
}

func lookupDot(items map[string]any, key string) (any, bool) {
	if key == "." || key == "" || strings.Contains(key, "..") {
		return nil, false
	}

	parts := strings.Split(key, ".")

	var current any = items

	for _, part := range parts {
		if part == "" {
			return nil, false
		}

		switch typed := current.(type) {
		case map[string]any:
			value, ok := typed[part]

			if !ok {
				return nil, false
			}

			current = value
		case []any:
			index, err := strconv.Atoi(part)

			if err != nil || index < 0 || index >= len(typed) {
				return nil, false
			}

			current = typed[index]
		default:
			return nil, false
		}
	}

	return current, true
}

func setDot(items map[string]any, key string, value any) {
	if key == "." || key == "" || strings.Contains(key, "..") {
		items[key] = value

		return
	}

	parts := strings.Split(key, ".")
	current := items

	for _, part := range parts[:len(parts)-1] {
		next, ok := current[part].(map[string]any)

		if !ok {
			next = map[string]any{}
			current[part] = next
		}

		current = next
	}

	current[parts[len(parts)-1]] = value
}

func mergeMap(target map[string]any, source map[string]any) {
	for key, value := range source {
		sourceMap, sourceIsMap := value.(map[string]any)
		targetMap, targetIsMap := target[key].(map[string]any)

		if sourceIsMap && targetIsMap {
			mergeMap(targetMap, sourceMap)

			continue
		}

		target[key] = cloneValue(value)
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

		for i, item := range typed {
			cloned[i] = cloneValue(item)
		}

		return cloned
	default:
		return typed
	}
}
