package config

import (
	"fmt"

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
	v *viper.Viper
}

// New creates a Repository pre-loaded with the given key-value pairs. The map
// may contain nested map[string]any values accessible via dot notation.
func New(items map[string]any) *Repository {
	v := viper.New()

	for key, value := range items {
		v.Set(key, value)
	}

	return &Repository{v: v}
}

// NewFromViper wraps an already-configured Viper instance. Use this when you
// have set up config file paths, environment prefixes, or other Viper options
// before creating the repository.
func NewFromViper(v *viper.Viper) *Repository {
	return &Repository{v: v}
}

// Viper returns the underlying Viper instance so consumers can configure file
// paths, environment prefixes, key replacers, and call ReadInConfig directly.
func (r *Repository) Viper() *viper.Viper {
	return r.v
}

// Has reports whether the given key is set in any configuration source.
func (r *Repository) Has(key string) bool {
	return r.v.IsSet(key)
}

// Get returns the value for key. If the key is not set, the first fallback
// value is returned (or nil when no fallback is provided).
func (r *Repository) Get(key string, fallback ...any) any {
	if !r.v.IsSet(key) {
		if len(fallback) > 0 {
			return fallback[0]
		}

		return nil
	}

	return r.v.Get(key)
}

// GetMany returns a map of values for the given keys. Keys that are not set
// map to nil.
func (r *Repository) GetMany(keys []string) map[string]any {
	result := make(map[string]any, len(keys))

	for _, key := range keys {
		result[key] = r.Get(key)
	}

	return result
}

// Set stores a value at the given dot-notation key.
func (r *Repository) Set(key string, value any) {
	r.v.Set(key, value)
}

// SetMany stores multiple key-value pairs at once.
func (r *Repository) SetMany(values map[string]any) {
	for key, value := range values {
		r.v.Set(key, value)
	}
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

// All returns every configuration item as a flat map.
func (r *Repository) All() map[string]any {
	return r.v.AllSettings()
}

// String returns the string value for key. If the key is not set, the first
// fallback is returned. An error wrapping ErrInvalidType is returned when the
// stored value is not a string.
func (r *Repository) String(key string, fallback ...string) (string, error) {
	if !r.v.IsSet(key) {
		if len(fallback) > 0 {
			return fallback[0], nil
		}

		return "", nil
	}

	value := r.v.Get(key)

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
	if !r.v.IsSet(key) {
		if len(fallback) > 0 {
			return fallback[0], nil
		}

		return 0, nil
	}

	value := r.v.Get(key)

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
	if !r.v.IsSet(key) {
		if len(fallback) > 0 {
			return fallback[0], nil
		}

		return 0, nil
	}

	value := r.v.Get(key)

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
	if !r.v.IsSet(key) {
		if len(fallback) > 0 {
			return fallback[0], nil
		}

		return false, nil
	}

	value := r.v.Get(key)

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
	if !r.v.IsSet(key) {
		if len(fallback) > 0 {
			return fallback[0], nil
		}

		return nil, nil
	}

	value := r.v.Get(key)

	a, ok := value.([]any)
	if !ok {
		return nil, fmt.Errorf("%w: value for key [%s] must be a []any, %T given", ErrInvalidType, key, value)
	}

	return a, nil
}
