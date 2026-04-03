package kv

// Dot flattens a nested map into a single-level map with dot-notated keys.
// An optional prefix is prepended to each key.
func Dot(items map[string]any, prepend ...string) map[string]any {
	prefix := ""

	if len(prepend) > 0 {
		prefix = prepend[0]
	}

	result := make(map[string]any)
	dotRecursive(items, prefix, result)

	return result
}

func dotRecursive(items map[string]any, prefix string, result map[string]any) {
	for key, value := range items {
		fullKey := key

		if prefix != "" {
			fullKey = prefix + "." + key
		}

		if nested, ok := value.(map[string]any); ok {
			dotRecursive(nested, fullKey, result)
		} else {
			result[fullKey] = value
		}
	}
}

// Undot expands a flat map with dot-notated keys into a nested map.
func Undot(items map[string]any) map[string]any {
	result := make(map[string]any)

	for key, value := range items {
		Set(result, key, value)
	}

	return result
}

// Map applies a callback to each value in the map and returns a new map
// with the transformed values.
func Map[V any, R any](items map[string]V, callback func(V, string) R) map[string]R {
	result := make(map[string]R, len(items))

	for k, v := range items {
		result[k] = callback(v, k)
	}

	return result
}

// Where returns a new map containing only the entries for which
// the callback returns true.
func Where[V any](items map[string]V, callback func(V, string) bool) map[string]V {
	result := make(map[string]V)

	for k, v := range items {
		if callback(v, k) {
			result[k] = v
		}
	}

	return result
}

// PrependKeysWith returns a new map with the given prefix prepended to every key.
func PrependKeysWith[V any](items map[string]V, prefix string) map[string]V {
	result := make(map[string]V, len(items))

	for k, v := range items {
		result[prefix+k] = v
	}

	return result
}

// Replace returns a new map with entries from items, overwritten by
// any matching keys from the replacement maps.
func Replace(items map[string]any, replacements ...map[string]any) map[string]any {
	result := make(map[string]any, len(items))

	for k, v := range items {
		result[k] = v
	}

	for _, replacement := range replacements {
		for k, v := range replacement {
			result[k] = v
		}
	}

	return result
}
