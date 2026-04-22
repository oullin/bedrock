package support

import (
	"encoding/json"
	"iter"
	"reflect"
	"strconv"
	"strings"
)

// Fluent provides a dynamic key-value bag backed by map[string]any.
// It supports dot-notation access for nested values.
// Mirrors Illuminate\Support\Fluent.
type Fluent struct {
	attributes map[string]any
}

// NewFluent creates a new Fluent instance.
// Accepts an optional initial value that can be a map or a struct with exported fields.
// Mirrors new Fluent($attributes).
func NewFluent(attrs ...any) *Fluent {
	f := &Fluent{attributes: make(map[string]any)}

	if len(attrs) > 0 {
		for k, v := range fluentAttributes(attrs[0]) {
			f.attributes[k] = v
		}
	}

	return f
}

// Get returns the value for the given key, or the default if not found.
// Supports simple key access (no dot notation, for compatibility with
// Laravel's Fluent which uses data_get internally).
// Mirrors Fluent::get().
func (f *Fluent) Get(key string, def ...any) any {
	val, ok := f.attributes[key]

	if !ok {
		if len(def) > 0 {
			if fn, isFn := def[0].(func() any); isFn {
				return fn()
			}

			return def[0]
		}

		return nil
	}

	return val
}

// Set sets the value for the given key.
// Mirrors Fluent::set().
func (f *Fluent) Set(key string, value any) *Fluent {
	f.attributes[key] = value

	return f
}

// Has reports whether the key exists in the attributes.
func (f *Fluent) Has(key string) bool {
	_, ok := f.attributes[key]

	return ok
}

// Missing reports whether the key is absent from the attributes.
func (f *Fluent) Missing(key string) bool {
	return !f.Has(key)
}

// All returns all attributes as a map.
// Mirrors Fluent::all().
func (f *Fluent) All() map[string]any {
	result := make(map[string]any, len(f.attributes))

	for k, v := range f.attributes {
		result[k] = v
	}

	return result
}

// Array returns all attributes as a map.
// Mirrors Fluent::toArray() / array access helpers.
func (f *Fluent) Array() map[string]any {
	return f.All()
}

// Only returns a map containing only the specified keys.
// Mirrors Fluent::only().
func (f *Fluent) Only(keys ...string) map[string]any {
	result := make(map[string]any, len(keys))

	for _, k := range keys {
		if v, ok := f.attributes[k]; ok {
			result[k] = v
		}
	}

	return result
}

// Except returns a map excluding the specified keys.
// Mirrors Fluent::except() (via Arr::except equivalent).
func (f *Fluent) Except(keys ...string) map[string]any {
	excluded := make(map[string]bool, len(keys))

	for _, k := range keys {
		excluded[k] = true
	}

	result := make(map[string]any)

	for k, v := range f.attributes {
		if !excluded[k] {
			result[k] = v
		}
	}

	return result
}

// Fill merges attributes into the Fluent, overwriting existing keys.
// Mirrors Fluent::fill().
func (f *Fluent) Fill(attrs map[string]any) *Fluent {
	for k, v := range attrs {
		f.attributes[k] = v
	}

	return f
}

// Merge merges attributes without overwriting existing keys.
// Mirrors the behavior of Fluent when used with array_merge where existing keys survive.
func (f *Fluent) Merge(attrs map[string]any) *Fluent {
	for k, v := range attrs {
		if _, exists := f.attributes[k]; !exists {
			f.attributes[k] = v
		}
	}

	return f
}

// Scope returns a new Fluent containing only keys that start with the given prefix.
// The prefix is stripped from the returned keys.
// Mirrors Fluent::scope().
func (f *Fluent) Scope(key string) *Fluent {
	result := make(map[string]any)
	prefix := key + "."

	for k, v := range f.attributes {
		if strings.HasPrefix(k, prefix) {
			result[k[len(prefix):]] = v
		} else if k == key {
			if nested, ok := v.(map[string]any); ok {
				return NewFluent(nested)
			}
		}
	}

	return NewFluent(result)
}

// IsEmpty reports whether the Fluent has no attributes.
func (f *Fluent) IsEmpty() bool {
	return len(f.attributes) == 0
}

// IsNotEmpty reports whether the Fluent has at least one attribute.
func (f *Fluent) IsNotEmpty() bool {
	return !f.IsEmpty()
}

// Count returns the number of attributes.
func (f *Fluent) Count() int {
	return len(f.attributes)
}

// String returns the string value of the given key.
// Mirrors Fluent::string().
func (f *Fluent) String(key string, def ...string) string {
	v := f.Get(key)

	if v == nil {
		if len(def) > 0 {
			return def[0]
		}

		return ""
	}

	switch s := v.(type) {
	case string:
		return s
	case []byte:
		return string(s)
	default:
		return strconv.Itoa(0) // fallback
	}
}

// Bool returns the bool value of the given key.
// Truthy string values ("yes", "on", "true", "1") return true.
// Mirrors Fluent::boolean().
func (f *Fluent) Bool(key string, def ...bool) bool {
	v := f.Get(key)

	if v == nil {
		if len(def) > 0 {
			return def[0]
		}

		return false
	}

	switch b := v.(type) {
	case bool:
		return b
	case int:
		return b != 0
	case float64:
		return b != 0
	case string:
		switch strings.ToLower(b) {
		case "true", "yes", "on", "1":
			return true
		}

		return false
	}

	return false
}

// Int returns the integer value of the given key.
// Mirrors Fluent::integer().
func (f *Fluent) Int(key string, def ...int) int {
	v := f.Get(key)

	if v == nil {
		if len(def) > 0 {
			return def[0]
		}

		return 0
	}

	switch n := v.(type) {
	case int:
		return n
	case int64:
		return int(n)
	case float64:
		return int(n)
	case string:
		i, _ := strconv.Atoi(n)

		return i
	}

	return 0
}

// Float returns the float value of the given key.
// Mirrors Fluent::float().
func (f *Fluent) Float(key string, def ...float64) float64 {
	v := f.Get(key)

	if v == nil {
		if len(def) > 0 {
			return def[0]
		}

		return 0
	}

	switch n := v.(type) {
	case float64:
		return n
	case float32:
		return float64(n)
	case int:
		return float64(n)
	case string:
		f, _ := strconv.ParseFloat(n, 64)

		return f
	}

	return 0
}

// ToMap returns all attributes as a map.
// Mirrors Fluent::toArray().
func (f *Fluent) ToMap() map[string]any {
	return f.All()
}

// MarshalJSON implements json.Marshaler.
func (f *Fluent) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.attributes)
}

// UnmarshalJSON implements json.Unmarshaler.
func (f *Fluent) UnmarshalJSON(data []byte) error {
	f.attributes = make(map[string]any)

	return json.Unmarshal(data, &f.attributes)
}

// ToJSON returns the JSON-encoded attributes.
func (f *Fluent) ToJSON() ([]byte, error) {
	return json.Marshal(f.attributes)
}

// ToPrettyJSON returns the JSON-encoded attributes with indentation.
func (f *Fluent) ToPrettyJSON() (string, error) {
	data, err := json.MarshalIndent(f.attributes, "", "  ")
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func fluentAttributes(value any) map[string]any {
	result := make(map[string]any)

	if value == nil {
		return result
	}

	switch v := value.(type) {
	case map[string]any:
		return v
	case map[string]string:
		for key, item := range v {
			result[key] = item
		}

		return result
	case *Fluent:
		return v.All()
	case Fluent:
		return v.All()
	case iter.Seq2[string, any]:
		v(func(key string, item any) bool {
			result[key] = item

			return true
		})

		return result
	}

	rv := reflect.ValueOf(value)
	if !rv.IsValid() {
		return result
	}

	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return result
		}

		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return result
	}

	rt := rv.Type()
	for i := 0; i < rv.NumField(); i++ {
		field := rt.Field(i)
		if field.PkgPath != "" {
			continue
		}

		name := field.Name
		if tag := field.Tag.Get("json"); tag != "" && tag != "-" {
			if comma := strings.Index(tag, ","); comma >= 0 {
				name = tag[:comma]
			} else {
				name = tag
			}
		}

		result[name] = rv.Field(i).Interface()
	}

	return result
}
