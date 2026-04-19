package rules

import "strings"

func init_size() {
	Register("Min", validateMin)
	Register("Max", validateMax)
	Register("Between", validateBetween)
	Register("Size", validateSize)
	Register("Gt", validateGt)
	Register("Gte", validateGte)
	Register("Lt", validateLt)
	Register("Lte", validateLte)
}

func validateMin(_ string, value any, params []string, ctx RuleContext) bool {
	min, ok := parseFloat64Param(params, 0)

	if !ok {
		return true
	}

	size, _ := getSize(value)

	return size >= min
}

func validateMax(_ string, value any, params []string, ctx RuleContext) bool {
	max, ok := parseFloat64Param(params, 0)

	if !ok {
		return true
	}

	size, _ := getSize(value)

	return size <= max
}

func validateBetween(_ string, value any, params []string, ctx RuleContext) bool {
	if len(params) < 2 {
		return true
	}

	min, ok1 := parseFloat64Param(params, 0)
	max, ok2 := parseFloat64Param(params, 1)

	if !ok1 || !ok2 {
		return true
	}

	size, _ := getSize(value)

	return size >= min && size <= max
}

func validateSize(_ string, value any, params []string, _ RuleContext) bool {
	target, ok := parseFloat64Param(params, 0)

	if !ok {
		return true
	}

	size, _ := getSize(value)

	return size == target
}

// validateGt: value must be greater than another field's value or a literal.
// Params: [otherFieldOrLiteral]
func validateGt(_ string, value any, params []string, ctx RuleContext) bool {
	if len(params) == 0 {
		return true
	}

	other := ctx.GetValue(params[0])

	if other == nil {
		// treat as literal numeric
		target, ok := parseFloat64Param(params, 0)

		if !ok {
			return true
		}

		size, _ := getSize(value)

		return size > target
	}

	otherSize, _ := getSize(other)
	mySize, _ := getSize(value)

	return mySize > otherSize
}

func validateGte(_ string, value any, params []string, ctx RuleContext) bool {
	if len(params) == 0 {
		return true
	}

	other := ctx.GetValue(params[0])

	if other == nil {
		target, ok := parseFloat64Param(params, 0)

		if !ok {
			return true
		}

		size, _ := getSize(value)

		return size >= target
	}

	otherSize, _ := getSize(other)
	mySize, _ := getSize(value)

	return mySize >= otherSize
}

func validateLt(_ string, value any, params []string, ctx RuleContext) bool {
	if len(params) == 0 {
		return true
	}

	other := ctx.GetValue(params[0])

	if other == nil {
		target, ok := parseFloat64Param(params, 0)

		if !ok {
			return true
		}

		size, _ := getSize(value)

		return size < target
	}

	otherSize, _ := getSize(other)
	mySize, _ := getSize(value)

	return mySize < otherSize
}

func validateLte(_ string, value any, params []string, ctx RuleContext) bool {
	if len(params) == 0 {
		return true
	}

	other := ctx.GetValue(params[0])

	if other == nil {
		target, ok := parseFloat64Param(params, 0)

		if !ok {
			return true
		}

		size, _ := getSize(value)

		return size <= target
	}

	otherSize, _ := getSize(other)
	mySize, _ := getSize(value)

	return mySize <= otherSize
}

// MessageTypeForSize returns the dot-qualified rule key for size-related messages.
// e.g. ("Max", "hello") → "Max.string"
func MessageTypeForSize(rule string, value any) string {
	if value == nil {
		return rule
	}

	_, kind := getSize(value)

	if kind == "string" || kind == "numeric" || kind == "array" {
		return rule + "." + kind
	}

	s, ok := value.(string)

	if ok {
		_ = strings.TrimSpace(s) // just to use the import
	}

	return rule
}
