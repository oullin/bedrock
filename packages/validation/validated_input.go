package validation

import "encoding/json"

// ValidatedInput is the result of Validator.Safe().  It wraps the validated
// data and provides helpers for extracting subsets.
type ValidatedInput struct {
	data map[string]any
}

// All returns the full validated data map.
func (v *ValidatedInput) All() map[string]any {
	return v.data
}

// Get returns the value for key, along with whether it was present.
func (v *ValidatedInput) Get(key string) (any, bool) {
	val, ok := v.data[key]

	return val, ok
}

// Only returns a new map containing only the specified keys.
func (v *ValidatedInput) Only(keys ...string) map[string]any {
	out := make(map[string]any, len(keys))

	for _, k := range keys {
		if val, ok := v.data[k]; ok {
			out[k] = val
		}
	}

	return out
}

// Except returns a new map excluding the specified keys.
func (v *ValidatedInput) Except(keys ...string) map[string]any {
	exclude := make(map[string]bool, len(keys))

	for _, k := range keys {
		exclude[k] = true
	}

	out := make(map[string]any, len(v.data))

	for k, val := range v.data {
		if !exclude[k] {
			out[k] = val
		}
	}

	return out
}

// ToJSON serialises the validated data as JSON.
func (v *ValidatedInput) ToJSON() ([]byte, error) {
	return json.Marshal(v.data)
}
