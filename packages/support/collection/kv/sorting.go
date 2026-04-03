package kv

import "sort"

// Sort returns a new map with the same entries, iterated by sorted keys.
// Note: Go maps do not guarantee iteration order; the returned map
// is sorted at construction time.
func Sort(items map[string]any) map[string]any {
	keys := make([]string, 0, len(items))

	for k := range items {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	result := make(map[string]any, len(items))

	for _, k := range keys {
		result[k] = items[k]
	}

	return result
}

// SortRecursive returns a new map with entries sorted by key,
// recursing into any nested maps.
func SortRecursive(items map[string]any) map[string]any {
	result := make(map[string]any, len(items))
	keys := make([]string, 0, len(items))

	for k := range items {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		v := items[k]

		if nested, ok := v.(map[string]any); ok {
			result[k] = SortRecursive(nested)
		} else {
			result[k] = v
		}
	}

	return result
}
