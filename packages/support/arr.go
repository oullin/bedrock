package support

import (
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// ArrAccessible reports whether value is a map, slice, array, or string.

// ArrMap maps each item in a slice.

// ArrMapWithKeys maps each item into key/value pairs.

// ArrString returns a string value for the key.

// ArrInteger returns an integer value for the key.

// ArrFloat returns a float64 value for the key.

// ArrBoolean returns a boolean value for the key.

// ArrArray returns a []any value for the key.

// ArrReject filters a map by rejecting values accepted by predicate.

// ArrWhereKey filters a map by key predicate.

// ArrFrom converts supported values into a []any.

// ArrKeyBy returns a map keyed by each item's selected value.

// ArrPrependKeysWith prefixes every key in the map.

// ArrSelect returns a slice of maps containing only the selected keys.

// ArrIsAssoc reports whether a map has non-list integer keys or string keys.

// ArrIsList reports whether value is a zero-based list.

// ArrAdd adds a key/value pair to the map if the key does not already exist.
// Supports dot-notation keys for nested access.

// ArrGet retrieves a value from a nested map using dot-notation.
// Returns the default value if the key is not found.

// ArrSet sets a value in a nested map using dot-notation, creating intermediate
// maps as needed. Mutates m and returns it for chaining.

// ArrHas checks whether the given key(s) exist in the map using dot-notation.
// Returns true only if ALL keys exist.

// ArrHasAny checks whether any of the given keys exist in the map.

// ArrForget removes one or more keys from the map using dot-notation.

// ArrExists checks whether the given key exists in the map.
// Supports dot-notation keys for nested access.

// ArrWhereNotNull returns a copy of the map without nil values.

// ArrExceptValues returns all key/value pairs whose values are not in the excluded list.

// ArrUndot expands dot-notation keys into nested maps.

// ArrJoin joins the slice into a string, using the final glue for the last separator when provided.

// ArrTake returns the first n items, or the last -n items when count is negative.

// ArrPush appends values to the end of the slice.

// ArrQuery encodes a map into a query string.

// ArrPull gets a value from the map and removes it.
// Returns the default if the key is not found.

// ArrDot flattens a nested map into a single-level map with dot-notation keys.

// ArrExcept returns all key/value pairs except those with the specified keys.

// ArrOnly returns a map containing only the specified keys.

// ArrOnlyValues returns values for the selected keys, preserving key order.

// ArrDivide splits a map into two slices: one of keys and one of values.
// Keys are returned in sorted order for deterministic output.

// ArrPluck extracts a list of values for a given key from a slice of maps.
// If indexKey is provided, the result is a map[string]any keyed by that field.
// Otherwise, the result is []any.

// ArrSortDesc returns a reverse-sorted copy of items.

// SortDirection controls an ArrSortByMany sort clause.
type SortDirection string

// SortAsc sorts a clause in ascending order.

// SortDesc sorts a clause in descending order.

// SortClause describes one key and direction for ArrSortByMany.
type SortClause struct {
	Key       string
	Direction SortDirection
}

func ArrAccessible(value any) bool {
	if value == nil {
		return false
	}

	switch value.(type) {
	case map[string]any, []any, string:
		return true
	}

	rv := reflect.ValueOf(value)

	switch rv.Kind() {
	case reflect.Map, reflect.Slice, reflect.Array:
		return true
	default:
		return false
	}
}

func ArrMap[T any, R any](items []T, mapper func(T) R) []R {
	result := make([]R, 0, len(items))

	for _, item := range items {
		result = append(result, mapper(item))
	}

	return result
}

func ArrMapWithKeys[T any](items []T, mapper func(T) map[string]any) map[string]any {
	result := make(map[string]any, len(items))

	for _, item := range items {
		for key, value := range mapper(item) {
			result[key] = value
		}
	}

	return result
}

func ArrString(m map[string]any, key string, def ...string) string {
	value := ArrGet(m, key)

	if value == nil {
		if len(def) > 0 {
			return def[0]
		}

		return ""
	}

	return fmt.Sprint(value)
}

func ArrInteger(m map[string]any, key string, def ...int) int {
	value := ArrGet(m, key)

	switch v := value.(type) {
	case int:
		return v
	case int64:
		return int(v)
	case float64:
		return int(v)
	case string:
		if parsed, err := strconv.Atoi(v); err == nil {
			return parsed
		}
	}

	if len(def) > 0 {
		return def[0]
	}

	return 0
}

func ArrFloat(m map[string]any, key string, def ...float64) float64 {
	value := ArrGet(m, key)

	switch v := value.(type) {
	case float64:
		return v
	case float32:
		return float64(v)
	case int:
		return float64(v)
	case string:
		if parsed, err := strconv.ParseFloat(v, 64); err == nil {
			return parsed
		}
	}

	if len(def) > 0 {
		return def[0]
	}

	return 0
}

func ArrBoolean(m map[string]any, key string, def ...bool) bool {
	value := ArrGet(m, key)

	switch v := value.(type) {
	case bool:
		return v
	case string:
		if parsed, err := strconv.ParseBool(v); err == nil {
			return parsed
		}
	}

	if len(def) > 0 {
		return def[0]
	}

	return false
}

func ArrArray(m map[string]any, key string, def ...[]any) []any {
	value := ArrGet(m, key)

	if items, ok := anySlice(value); ok {
		return items
	}

	if len(def) > 0 {
		return def[0]
	}

	return []any{}
}

func ArrReject[T any](items []T, predicate func(T, int) bool) []T {
	return ArrWhere(items, func(item T, index int) bool {
		return !predicate(item, index)
	})
}

func ArrWhereKey(m map[string]any, predicate func(string) bool) map[string]any {
	result := make(map[string]any)

	for key, value := range m {
		if predicate(key) {
			result[key] = value
		}
	}

	return result
}

func ArrFrom(value any) []any {
	if value == nil {
		return []any{}
	}

	if items, ok := anySlice(value); ok {
		return items
	}

	return []any{value}
}

func ArrKeyBy(items []map[string]any, key string) map[string]map[string]any {
	result := make(map[string]map[string]any, len(items))

	for _, item := range items {
		result[fmt.Sprint(ArrGet(item, key))] = item
	}

	return result
}

func ArrPrependKeysWith(m map[string]any, prefix string) map[string]any {
	result := make(map[string]any, len(m))

	for key, value := range m {
		result[prefix+key] = value
	}

	return result
}

func ArrSelect(items []map[string]any, keys ...string) []map[string]any {
	result := make([]map[string]any, 0, len(items))

	for _, item := range items {
		result = append(result, ArrOnly(item, keys...))
	}

	return result
}

func ArrIsAssoc(value any) bool {
	switch v := value.(type) {
	case map[string]any:
		return len(v) > 0
	case map[int]any:
		return !ArrIsList(v)
	default:
		return false
	}
}

func ArrIsList(value any) bool {
	switch v := value.(type) {
	case []any:
		return true
	case []string:
		return true
	case []int:
		return true
	case map[int]any:
		for i := 0; i < len(v); i++ {
			if _, ok := v[i]; !ok {
				return false
			}
		}

		return true
	default:
		rv := reflect.ValueOf(value)

		switch rv.Kind() {
		case reflect.Slice, reflect.Array:
			return true
		default:
			return false
		}
	}
}

func ArrAdd(m map[string]any, key string, value any) map[string]any {
	if !dotHas(m, key) {
		dotSet(m, key, value)
	}

	return m
}

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

func ArrSet(m map[string]any, key string, value any) map[string]any {
	dotSet(m, key, value)

	return m
}

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

func ArrHasAny(m map[string]any, keys ...string) bool {
	for _, key := range keys {
		if dotHas(m, key) {
			return true
		}
	}

	return false
}

func ArrForget(m map[string]any, keys ...string) {
	for _, key := range keys {
		dotForget(m, key)
	}
}

func ArrExists(m map[string]any, key string) bool {
	return dotHas(m, key)
}

func ArrWhereNotNull(m map[string]any) map[string]any {
	result := make(map[string]any, len(m))

	for k, v := range m {
		if v != nil {
			result[k] = v
		}
	}

	return result
}

func ArrExceptValues(m map[string]any, values ...any) map[string]any {
	result := make(map[string]any, len(m))

	for k, v := range m {
		if !containsValue(values, v) {
			result[k] = v
		}
	}

	return result
}

func ArrUndot(m map[string]any) map[string]any {
	result := make(map[string]any, len(m))

	for key, value := range m {
		dotSet(result, key, value)
	}

	return result
}

func ArrJoin(items []string, glue string, finalGlue ...string) string {
	switch len(items) {
	case 0:
		return ""
	case 1:
		return items[0]
	case 2:
		if len(finalGlue) > 0 {
			return items[0] + finalGlue[0] + items[1]
		}

		return items[0] + glue + items[1]
	}

	if len(finalGlue) == 0 {
		return strings.Join(items, glue)
	}

	return strings.Join(items[:len(items)-1], glue) + finalGlue[0] + items[len(items)-1]
}

func ArrTake[T any](items []T, count int) []T {
	if count == 0 {
		return []T{}
	}

	if count > 0 {
		if count >= len(items) {
			return append([]T(nil), items...)
		}

		return append([]T(nil), items[:count]...)
	}

	n := -count

	if n >= len(items) {
		return append([]T(nil), items...)
	}

	return append([]T(nil), items[len(items)-n:]...)
}

func ArrPush[T any](items []T, values ...T) []T {
	return append(items, values...)
}

func ArrQuery(m map[string]any) string {
	if len(m) == 0 {
		return ""
	}

	values := url.Values{}

	for key, value := range m {
		switch v := value.(type) {
		case []string:
			for _, item := range v {
				values.Add(key, item)
			}
		case []any:
			for _, item := range v {
				values.Add(key, fmt.Sprint(item))
			}
		default:
			rv := reflect.ValueOf(value)

			if rv.IsValid() && rv.Kind() == reflect.Slice {
				for i := 0; i < rv.Len(); i++ {
					values.Add(key, fmt.Sprint(rv.Index(i).Interface()))
				}

				continue
			}

			values.Set(key, fmt.Sprint(value))
		}
	}

	return values.Encode()
}

func ArrPull(m map[string]any, key string, def ...any) any {
	val := ArrGet(m, key, def...)
	ArrForget(m, key)

	return val
}

func ArrDot(m map[string]any, prepend ...string) map[string]any {
	p := ""

	if len(prepend) > 0 {
		p = prepend[0]
	}

	result := make(map[string]any)
	dotFlatten(m, p, result)

	return result
}

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

func ArrOnly(m map[string]any, keys ...string) map[string]any {
	result := make(map[string]any, len(keys))

	for _, k := range keys {
		if v, ok := m[k]; ok {
			result[k] = v
		}
	}

	return result
}

func ArrOnlyValues(m map[string]any, keys ...string) []any {
	result := make([]any, 0, len(keys))

	for _, key := range keys {
		if value, ok := m[key]; ok {
			result = append(result, value)
		}
	}

	return result
}

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

func ArrPluck(items []map[string]any, valueKey string, indexKey ...string) any {
	if len(indexKey) > 0 && indexKey[0] != "" {
		result := make(map[string]any, len(items))

		for _, item := range items {
			val := ArrGet(item, valueKey)
			idx := fmt.Sprint(ArrGet(item, indexKey[0]))
			result[idx] = val
		}

		return result
	}

	result := make([]any, 0, len(items))

	for _, item := range items {
		result = append(result, ArrGet(item, valueKey))
	}

	return result
}

func ArrSortDesc[T ~string | ~int | ~int64 | ~float64](items []T) []T {
	result := append([]T(nil), items...)

	sort.Slice(result, func(i, j int) bool {
		return result[i] > result[j]
	})

	return result
}

const (
	SortAsc SortDirection = "asc"

	SortDesc SortDirection = "desc"
)

// ArrSortByMany sorts maps by multiple dot-notation keys.
func ArrSortByMany(items []map[string]any, clauses ...SortClause) []map[string]any {
	result := append([]map[string]any(nil), items...)

	sort.SliceStable(result, func(i, j int) bool {
		for _, clause := range clauses {
			left := fmt.Sprint(ArrGet(result[i], clause.Key))
			right := fmt.Sprint(ArrGet(result[j], clause.Key))

			if left == right {
				continue
			}

			if clause.Direction == SortDesc {
				return left > right
			}

			return left < right
		}

		return false
	})

	return result
}

func containsValue(values []any, target any) bool {
	for _, value := range values {
		if reflect.DeepEqual(value, target) {
			return true
		}
	}

	return false
}

// ArrSortRecursive recursively sorts a map by keys and any nested slices/maps.
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

func anySlice(value any) ([]any, bool) {
	if items, ok := value.([]any); ok {
		return items, true
	}

	rv := reflect.ValueOf(value)

	if !rv.IsValid() {
		return nil, false
	}

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		result := make([]any, rv.Len())

		for i := 0; i < rv.Len(); i++ {
			result[i] = rv.Index(i).Interface()
		}

		return result, true
	default:
		return nil, false
	}
}
