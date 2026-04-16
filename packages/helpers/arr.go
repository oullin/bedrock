package helpers

import (
	"fmt"
	"sort"
)

// ArrAdd adds a key/value pair to the map if the key does not already exist.
// Supports dot-notation keys for nested access.
// Mirrors Arr::add().
func ArrAdd(m map[string]any, key string, value any) map[string]any {
	if !dotHas(m, key) {
		dotSet(m, key, value)
	}

	return m
}

// ArrGet retrieves a value from a nested map using dot-notation.
// Returns the default value if the key is not found.
// Mirrors Arr::get().
func ArrGet(m map[string]any, key string, def ...any) any {
	val, ok := dotGet(m, key)
	if ok {
		return val
	}

	if len(def) > 0 {
		return def[0]
	}

	return nil
}

// ArrSet sets a value in a nested map using dot-notation, creating intermediate
// maps as needed. Mutates m and returns it for chaining.
// Mirrors Arr::set().
func ArrSet(m map[string]any, key string, value any) map[string]any {
	dotSet(m, key, value)
	return m
}

// ArrHas checks whether the given key(s) exist in the map using dot-notation.
// Returns true only if ALL keys exist.
// Mirrors Arr::has().
func ArrHas(m map[string]any, keys ...string) bool {
	if len(keys) == 0 {
		return false
	}

	for _, key := range keys {
		if !dotHas(m, key) {
			return false
		}
	}

	return true
}

// ArrForget removes one or more keys from the map using dot-notation.
// Mirrors Arr::forget().
func ArrForget(m map[string]any, keys ...string) {
	for _, key := range keys {
		dotForget(m, key)
	}
}

// ArrPull gets a value from the map and removes it.
// Returns the default if the key is not found.
// Mirrors Arr::pull().
func ArrPull(m map[string]any, key string, def ...any) any {
	val := ArrGet(m, key, def...)
	ArrForget(m, key)

	return val
}

// ArrDot flattens a nested map into a single-level map with dot-notation keys.
// Mirrors Arr::dot().
func ArrDot(m map[string]any, prepend ...string) map[string]any {
	p := ""
	if len(prepend) > 0 {
		p = prepend[0]
	}

	result := make(map[string]any)
	dotFlatten(m, p, result)

	return result
}

// ArrExcept returns all key/value pairs except those with the specified keys.
// Mirrors Arr::except().
func ArrExcept(m map[string]any, keys ...string) map[string]any {
	exclude := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		exclude[k] = struct{}{}
	}

	result := make(map[string]any, len(m))
	for k, v := range m {
		if _, skip := exclude[k]; !skip {
			result[k] = v
		}
	}

	return result
}

// ArrOnly returns a map containing only the specified keys.
// Mirrors Arr::only().
func ArrOnly(m map[string]any, keys ...string) map[string]any {
	result := make(map[string]any, len(keys))
	for _, k := range keys {
		if v, ok := m[k]; ok {
			result[k] = v
		}
	}

	return result
}

// ArrDivide splits a map into two slices: one of keys and one of values.
// Keys are returned in sorted order for deterministic output.
// Mirrors Arr::divide().
func ArrDivide(m map[string]any) ([]string, []any) {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	values := make([]any, len(keys))
	for i, k := range keys {
		values[i] = m[k]
	}

	return keys, values
}

// ArrPluck extracts a list of values for a given key from a slice of maps.
// If indexKey is provided, the result is a map[string]any keyed by that field.
// Otherwise, the result is []any.
// Mirrors Arr::pluck().
func ArrPluck(items []map[string]any, valueKey string, indexKey ...string) any {
	if len(indexKey) > 0 && indexKey[0] != "" {
		result := make(map[string]any, len(items))

		for _, item := range items {
			val := item[valueKey]
			idx := fmt.Sprint(item[indexKey[0]])
			result[idx] = val
		}

		return result
	}

	result := make([]any, 0, len(items))
	for _, item := range items {
		result = append(result, item[valueKey])
	}

	return result
}

// ArrSortRecursive recursively sorts a map by keys and any nested slices/maps.
// Mirrors Arr::sortRecursive().
func ArrSortRecursive(m map[string]any, descending ...bool) map[string]any {
	desc := len(descending) > 0 && descending[0]

	result := make(map[string]any, len(m))
	for k, v := range m {
		switch child := v.(type) {
		case map[string]any:
			result[k] = ArrSortRecursive(child, desc)
		case []any:
			result[k] = sortSliceRecursive(child, desc)
		default:
			result[k] = v
		}
	}

	return result
}

func sortSliceRecursive(items []any, desc bool) []any {
	sorted := make([]any, len(items))
	copy(sorted, items)

	for i, v := range sorted {
		switch child := v.(type) {
		case map[string]any:
			sorted[i] = ArrSortRecursive(child, desc)
		case []any:
			sorted[i] = sortSliceRecursive(child, desc)
		}
	}

	sort.SliceStable(sorted, func(i, j int) bool {
		si := fmt.Sprint(sorted[i])
		sj := fmt.Sprint(sorted[j])

		if desc {
			return si > sj
		}

		return si < sj
	})

	return sorted
}
