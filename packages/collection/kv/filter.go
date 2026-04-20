package kv

// Only returns a new map containing only the specified keys.
func Only(items map[string]any, keys ...string) map[string]any {
	result := make(map[string]any)

	for _, key := range keys {
		if v, ok := items[key]; ok {
			result[key] = v
		}
	}

	return result
}

// Except returns a new map with the specified keys removed.
func Except(items map[string]any, keys ...string) map[string]any {
	result := make(map[string]any, len(items))

	for k, v := range items {
		result[k] = v
	}

	for _, key := range keys {
		delete(result, key)
	}

	return result
}

// IsAssoc reports whether the map is non-empty.
func IsAssoc(items map[string]any) bool {
	return len(items) > 0
}
