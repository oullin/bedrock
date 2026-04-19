package rules

import "strings"

func init_comparison() {
	Register("Same", validateSame)
	Register("Different", validateDifferent)
	Register("Confirmed", validateConfirmed)
	Register("In", validateIn)
	Register("NotIn", validateNotIn)
	Register("InArray", validateInArray)
}

// validateSame: value must match the value of another field.
// Params: [otherField]
func validateSame(_ string, value any, params []string, ctx RuleContext) bool {
	if len(params) == 0 {
		return true
	}

	other := ctx.GetValue(params[0])

	return stringify(value) == stringify(other)
}

// validateDifferent: value must differ from the value of another field.
func validateDifferent(_ string, value any, params []string, ctx RuleContext) bool {
	if len(params) == 0 {
		return true
	}

	other := ctx.GetValue(params[0])

	return stringify(value) != stringify(other)
}

// validateConfirmed: value must match the value of "{attribute}_confirmation".
func validateConfirmed(attribute string, value any, _ []string, ctx RuleContext) bool {
	confirmation := ctx.GetValue(attribute + "_confirmation")

	return stringify(value) == stringify(confirmation)
}

// validateIn: value must be in the given list.
// Params: [value1, value2, ...]
func validateIn(_ string, value any, params []string, _ RuleContext) bool {
	s := stringify(value)

	for _, p := range params {
		if s == strings.TrimSpace(p) {
			return true
		}
	}

	return false
}

// validateNotIn: value must NOT be in the given list.
func validateNotIn(_ string, value any, params []string, ctx RuleContext) bool {
	return !validateIn("", value, params, ctx)
}

// validateInArray: value must exist in the values of another array field.
// Params: [otherField]
func validateInArray(_ string, value any, params []string, ctx RuleContext) bool {
	if len(params) == 0 {
		return true
	}

	other := ctx.GetValue(params[0])

	arr, ok := other.([]any)

	if !ok {
		return false
	}

	s := stringify(value)

	for _, item := range arr {
		if stringify(item) == s {
			return true
		}
	}

	return false
}
