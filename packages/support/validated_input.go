package support

import "strconv"

// ValidatedInput wraps validated request data with typed access helpers.
type ValidatedInput struct {
	values map[string]any
}

// NewValidatedInput creates a validated input wrapper.
func NewValidatedInput(values map[string]any) ValidatedInput {
	copied := make(map[string]any, len(values))

	for key, value := range values {
		copied[key] = value
	}

	return ValidatedInput{values: copied}
}

// All returns all validated values.
func (v ValidatedInput) All() map[string]any {
	copied := make(map[string]any, len(v.values))

	for key, value := range v.values {
		copied[key] = value
	}

	return copied
}

// Input returns a value by key or default.
func (v ValidatedInput) Input(key string, defaults ...any) any {
	if value, ok := v.values[key]; ok {
		return value
	}

	if len(defaults) > 0 {
		return defaults[0]
	}

	return nil
}

// Exists reports whether all keys are present.
func (v ValidatedInput) Exists(keys ...string) bool {
	for _, key := range keys {
		if _, ok := v.values[key]; !ok {
			return false
		}
	}

	return len(keys) > 0
}

// Has reports whether all keys are present and non-empty.
func (v ValidatedInput) Has(keys ...string) bool {
	for _, key := range keys {
		value, ok := v.values[key]

		if !ok || Blank(value) {
			return false
		}
	}

	return len(keys) > 0
}

// HasAny reports whether any key is present and non-empty.
func (v ValidatedInput) HasAny(keys ...string) bool {
	for _, key := range keys {
		if value, ok := v.values[key]; ok && !Blank(value) {
			return true
		}
	}

	return false
}

// Missing reports whether key is absent.
func (v ValidatedInput) Missing(key string) bool {
	_, ok := v.values[key]

	return !ok
}

// Filled reports whether key is present and non-blank.
func (v ValidatedInput) Filled(key string) bool {
	value, ok := v.values[key]

	return ok && !Blank(value)
}

// AnyFilled reports whether any key is present and non-blank.
func (v ValidatedInput) AnyFilled(keys ...string) bool {
	for _, key := range keys {
		if v.Filled(key) {
			return true
		}
	}

	return false
}

// Only returns selected keys.
func (v ValidatedInput) Only(keys ...string) map[string]any {
	return ArrOnly(v.values, keys...)
}

// Except returns all values except selected keys.
func (v ValidatedInput) Except(keys ...string) map[string]any {
	return ArrExcept(v.values, keys...)
}

// Keys returns all keys.
func (v ValidatedInput) Keys() []string {
	keys, _ := ArrDivide(v.values)

	return keys
}

// Merge returns a new input wrapper with values overwritten by extra.
func (v ValidatedInput) Merge(extra map[string]any) ValidatedInput {
	merged := v.All()

	for key, value := range extra {
		merged[key] = value
	}

	return NewValidatedInput(merged)
}

// Bool returns a bool value.
func (v ValidatedInput) Bool(key string, defaults ...bool) bool {
	return NewFluent(v.values).Bool(key, defaults...)
}

// Int returns an int value.
func (v ValidatedInput) Int(key string, defaults ...int) int {
	return NewFluent(v.values).Int(key, defaults...)
}

// Float returns a float value.
func (v ValidatedInput) Float(key string, defaults ...float64) float64 {
	return NewFluent(v.values).Float(key, defaults...)
}

// String returns a string value.
func (v ValidatedInput) String(key string, defaults ...string) string {
	value := v.Input(key)

	if value == nil {
		if len(defaults) > 0 {
			return defaults[0]
		}

		return ""
	}

	if stringValue, ok := value.(string); ok {
		return stringValue
	}

	return strconv.FormatInt(int64(v.Int(key)), 10)
}
