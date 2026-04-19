package rules

func init_array() {
	Register("Distinct", validateDistinct)
	Register("RequiredArrayKeys", validateRequiredArrayKeys)
	Register("InArrayKeys", validateInArrayKeys)
	Register("Contains", validateContains)
	Register("DoesntContain", validateDoesntContain)
}

// validateDistinct: no duplicate values in the array.
// Params: optional "strict" or "ignore_case"
func validateDistinct(_ string, value any, params []string, _ RuleContext) bool {
	arr, ok := value.([]any)

	if !ok {
		return true // non-array values are considered passing
	}

	ignoreCase := containsString(params, "ignore_case")
	seen := make(map[string]bool, len(arr))

	for _, item := range arr {
		key := stringify(item)

		if ignoreCase {
			key = toLower(key)
		}

		if seen[key] {
			return false
		}

		seen[key] = true
	}

	return true
}

// validateRequiredArrayKeys: array must contain all specified keys.
// Params: [key1, key2, ...]
func validateRequiredArrayKeys(_ string, value any, params []string, _ RuleContext) bool {
	m, ok := value.(map[string]any)

	if !ok {
		return false
	}

	for _, k := range params {
		if _, exists := m[k]; !exists {
			return false
		}
	}

	return true
}

// validateInArrayKeys: array may only contain the specified keys.
// Params: [key1, key2, ...]
func validateInArrayKeys(_ string, value any, params []string, _ RuleContext) bool {
	m, ok := value.(map[string]any)

	if !ok {
		return true
	}

	for k := range m {
		if !containsString(params, k) {
			return false
		}
	}

	return true
}

// validateContains: array must contain all specified values.
// Params: [value1, value2, ...]
func validateContains(_ string, value any, params []string, _ RuleContext) bool {
	arr, ok := value.([]any)

	if !ok {
		return false
	}

	for _, p := range params {
		found := false

		for _, item := range arr {
			if stringify(item) == p {
				found = true

				break
			}
		}

		if !found {
			return false
		}
	}

	return true
}

// validateDoesntContain: array must not contain any of the specified values.
func validateDoesntContain(_ string, value any, params []string, _ RuleContext) bool {
	arr, ok := value.([]any)

	if !ok {
		return true
	}

	for _, item := range arr {
		if containsString(params, stringify(item)) {
			return false
		}
	}

	return true
}

func toLower(s string) string {
	result := make([]byte, len(s))

	for i := range s {
		c := s[i]

		if c >= 'A' && c <= 'Z' {
			c += 32
		}

		result[i] = c
	}

	return string(result)
}
