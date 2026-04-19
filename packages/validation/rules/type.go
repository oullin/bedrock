package rules

import (
	"encoding/json"
	"reflect"
	"strings"
)

func init_type() {
	Register("String", validateString)
	Register("Numeric", validateNumericType)
	Register("Integer", validateInteger)
	Register("Boolean", validateBoolean)
	Register("Array", validateArray)
	Register("List", validateList)
	Register("Json", validateJson)
}

func validateString(_ string, value any, _ []string, _ RuleContext) bool {
	if value == nil {
		return false
	}

	_, ok := value.(string)

	return ok
}

func validateNumericType(_ string, value any, _ []string, _ RuleContext) bool {
	return isNumericValue(value)
}

func validateInteger(_ string, value any, _ []string, _ RuleContext) bool {
	return isIntegerValue(value)
}

func validateBoolean(_ string, value any, _ []string, _ RuleContext) bool {
	if value == nil {
		return false
	}

	switch v := value.(type) {
	case bool:
		return true
	case int:
		return v == 0 || v == 1
	case int8:
		return v == 0 || v == 1
	case int16:
		return v == 0 || v == 1
	case int32:
		return v == 0 || v == 1
	case int64:
		return v == 0 || v == 1
	case uint:
		return v == 0 || v == 1
	case uint8:
		return v == 0 || v == 1
	case uint16:
		return v == 0 || v == 1
	case uint32:
		return v == 0 || v == 1
	case uint64:
		return v == 0 || v == 1
	case float64:
		return v == 0 || v == 1
	case string:
		lower := strings.ToLower(strings.TrimSpace(v))

		return lower == "true" || lower == "false" || lower == "1" || lower == "0"
	}

	return false
}

func validateArray(_ string, value any, params []string, _ RuleContext) bool {
	if value == nil {
		return false
	}

	rv := reflect.ValueOf(value)

	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Map {
		_, ok := value.([]any)

		if !ok {
			_, ok = value.(map[string]any)

			if !ok {
				return false
			}
		}
	}

	// If params are given they specify the only allowed keys
	if len(params) > 0 {
		m, ok := value.(map[string]any)

		if !ok {
			return true // non-map array — key constraint doesn't apply
		}

		for k := range m {
			if !containsString(params, k) {
				return false
			}
		}
	}

	return true
}

// validateList requires the value to be a sequential (0-indexed) array, not a map.
func validateList(_ string, value any, _ []string, _ RuleContext) bool {
	_, ok := value.([]any)

	return ok
}

func validateJson(_ string, value any, _ []string, _ RuleContext) bool {
	s, ok := value.(string)

	if !ok {
		return false
	}

	var out any

	return json.Unmarshal([]byte(strings.TrimSpace(s)), &out) == nil
}
