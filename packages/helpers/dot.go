package helpers

import "strings"

// dotGet retrieves a value from a nested map using dot-notation.
// It checks the literal key first, then splits on ".".
// Returns (value, true) if found, (nil, false) otherwise.
func dotGet(m map[string]any, key string) (any, bool) {
	if val, ok := m[key]; ok {
		return val, true
	}

	parts := strings.SplitN(key, ".", 2)
	val, ok := m[parts[0]]

	if !ok {
		return nil, false
	}

	if len(parts) == 1 {
		return val, true
	}

	sub, ok := val.(map[string]any)

	if !ok {
		return nil, false
	}

	return dotGet(sub, parts[1])
}

// dotSet sets a value in a nested map using dot-notation,
// creating intermediate maps as needed.
func dotSet(m map[string]any, key string, value any) {
	parts := strings.SplitN(key, ".", 2)

	if len(parts) == 1 {
		m[key] = value

		return
	}

	sub, ok := m[parts[0]].(map[string]any)

	if !ok {
		sub = make(map[string]any)
		m[parts[0]] = sub
	}

	dotSet(sub, parts[1], value)
}

// dotHas checks if a key exists in a nested map using dot-notation.
func dotHas(m map[string]any, key string) bool {
	if _, ok := m[key]; ok {
		return true
	}

	parts := strings.SplitN(key, ".", 2)
	val, ok := m[parts[0]]

	if !ok {
		return false
	}

	if len(parts) == 1 {
		return true
	}

	sub, ok := val.(map[string]any)

	if !ok {
		return false
	}

	return dotHas(sub, parts[1])
}

// dotForget removes a key from a nested map using dot-notation.
func dotForget(m map[string]any, key string) {
	parts := strings.SplitN(key, ".", 2)

	if len(parts) == 1 {
		delete(m, key)

		return
	}

	sub, ok := m[parts[0]].(map[string]any)

	if !ok {
		return
	}

	dotForget(sub, parts[1])
}

// dotFlatten converts a nested map into a flat map with dot-notation keys.
func dotFlatten(m map[string]any, prepend string, result map[string]any) {
	for k, v := range m {
		key := k

		if prepend != "" {
			key = prepend + "." + k
		}

		sub, ok := v.(map[string]any)

		if ok {
			dotFlatten(sub, key, result)
		} else {
			result[key] = v
		}
	}
}
