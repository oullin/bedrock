package rules

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// StringifyValue converts any value to its string representation.
// Exported so the parent validation package can reuse it.
func StringifyValue(v any) string {
	return stringify(v)
}

// stringify converts any value to its string representation.
func stringify(v any) string {
	if v == nil {
		return ""
	}

	switch val := v.(type) {
	case string:
		return val
	case bool:
		if val {
			return "true"
		}

		return "false"
	case int:
		return strconv.Itoa(val)
	case int8:
		return strconv.FormatInt(int64(val), 10)
	case int16:
		return strconv.FormatInt(int64(val), 10)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case int64:
		return strconv.FormatInt(val, 10)
	case uint:
		return strconv.FormatUint(uint64(val), 10)
	case uint64:
		return strconv.FormatUint(val, 10)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", val)
	}
}

// isBlank reports whether value is nil, an empty string, or an empty slice/map.
func isBlank(value any) bool {
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case string:
		return strings.TrimSpace(v) == ""
	case []any:
		return len(v) == 0
	case map[string]any:
		return len(v) == 0
	}

	rv := reflect.ValueOf(value)

	switch rv.Kind() {
	case reflect.Slice, reflect.Map:
		return rv.Len() == 0
	case reflect.Ptr:
		return rv.IsNil()
	}

	return false
}

// isFilled is the inverse of isBlank.
func isFilled(value any) bool {
	return !isBlank(value)
}

// toFloat64 coerces a value to float64.  Returns (value, true) on success.
func toFloat64(v any) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case int8:
		return float64(val), true
	case int16:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint8:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	case float32:
		return float64(val), true
	case float64:
		return val, true
	case string:
		f, err := strconv.ParseFloat(strings.TrimSpace(val), 64)

		if err == nil {
			return f, true
		}
	}

	return 0, false
}

// isNumericValue reports whether v can be interpreted as a number.
func isNumericValue(v any) bool {
	_, ok := toFloat64(v)

	return ok
}

// isIntegerValue reports whether v is (or represents) an integer.
func isIntegerValue(v any) bool {
	f, ok := toFloat64(v)

	if !ok {
		return false
	}

	return f == math.Trunc(f)
}

// getSize returns the "size" of a value for comparison rules (min, max, between).
// For strings → len(runes); for numeric → float value; for arrays → len.
func getSize(v any) (float64, string) {
	if v == nil {
		return 0, "string"
	}

	switch val := v.(type) {
	case string:
		return float64(len([]rune(val))), "string"
	case []any:
		return float64(len(val)), "array"
	case map[string]any:
		return float64(len(val)), "array"
	}

	rv := reflect.ValueOf(v)

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		return float64(rv.Len()), "array"
	case reflect.Map:
		return float64(rv.Len()), "array"
	}

	if f, ok := toFloat64(v); ok {
		return f, "numeric"
	}

	s := stringify(v)

	return float64(len([]rune(s))), "string"
}

// parseFloat64Param parses the first parameter as float64, returning 0 on
// failure.
func parseFloat64Param(params []string, idx int) (float64, bool) {
	if idx >= len(params) {
		return 0, false
	}

	f, err := strconv.ParseFloat(strings.TrimSpace(params[idx]), 64)

	if err != nil {
		return 0, false
	}

	return f, true
}

// containsString reports whether ss contains s.
func containsString(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}

	return false
}
