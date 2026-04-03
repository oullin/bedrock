package kv

import "strings"

// Get retrieves a value from a nested map using dot notation.
// If the key is not found, the first default value is returned, or nil.
func Get(target map[string]any, key string, defaults ...any) any {
	if key == "" {
		return target
	}

	segments := strings.Split(key, ".")

	var current any = target

	for _, segment := range segments {
		switch v := current.(type) {
		case map[string]any:
			val, ok := v[segment]

			if !ok {
				if len(defaults) > 0 {
					return defaults[0]
				}

				return nil
			}

			current = val
		default:
			if len(defaults) > 0 {
				return defaults[0]
			}

			return nil
		}
	}

	return current
}

// Set sets a value in a nested map using dot notation.
// By default, existing values are overwritten; pass false to preserve them.
func Set(target map[string]any, key string, value any, overwrite ...bool) map[string]any {
	shouldOverwrite := true

	if len(overwrite) > 0 {
		shouldOverwrite = overwrite[0]
	}

	segments := strings.Split(key, ".")
	current := target

	for i, segment := range segments {
		if i == len(segments)-1 {
			if shouldOverwrite {
				current[segment] = value
			} else {
				if _, exists := current[segment]; !exists {
					current[segment] = value
				}
			}
		} else {
			if _, exists := current[segment]; !exists {
				current[segment] = make(map[string]any)
			}

			if next, ok := current[segment].(map[string]any); ok {
				current = next
			} else {
				newMap := make(map[string]any)
				current[segment] = newMap
				current = newMap
			}
		}
	}

	return target
}

// Has reports whether the given dot-notated key exists in the nested map.
func Has(target map[string]any, key string) bool {
	if key == "" {
		return false
	}

	segments := strings.Split(key, ".")

	var current any = target

	for _, segment := range segments {
		switch v := current.(type) {
		case map[string]any:
			val, ok := v[segment]

			if !ok {
				return false
			}

			current = val
		default:
			return false
		}
	}

	return true
}

// Fill sets the value at the given dot-notated key only if it does not already exist.
func Fill(target map[string]any, key string, value any) map[string]any {
	return Set(target, key, value, false)
}

// Add is an alias for [Fill]. It sets a value only if the key does not already exist.
func Add(target map[string]any, key string, value any) map[string]any {
	return Fill(target, key, value)
}

// Forget removes the value at the given dot-notated key from the nested map.
func Forget(target map[string]any, key string) map[string]any {
	segments := strings.Split(key, ".")

	if len(segments) == 1 {
		delete(target, key)

		return target
	}

	current := target

	for i := 0; i < len(segments)-1; i++ {
		val, ok := current[segments[i]]

		if !ok {
			return target
		}

		next, ok := val.(map[string]any)

		if !ok {
			return target
		}

		current = next
	}

	delete(current, segments[len(segments)-1])

	return target
}

// ForgetMany removes one or more dot-notated keys from the map.
func ForgetMany(items map[string]any, keys ...string) map[string]any {
	for _, key := range keys {
		Forget(items, key)
	}

	return items
}

// Pull removes a key from the map and returns its value.
// If the key does not exist, the first default value is returned.
func Pull(items map[string]any, key string, defaults ...any) any {
	value := Get(items, key, defaults...)
	Forget(items, key)

	return value
}

// HasAll reports whether all the given dot-notated keys exist in the map.
func HasAll(items map[string]any, keys ...string) bool {
	for _, key := range keys {
		if !Has(items, key) {
			return false
		}
	}

	return true
}

// HasAny reports whether at least one of the given dot-notated keys exists in the map.
func HasAny(items map[string]any, keys ...string) bool {
	for _, key := range keys {
		if Has(items, key) {
			return true
		}
	}

	return false
}
