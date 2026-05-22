package rules

import "strings"

func init_accepted() {
	RegisterImplicit("Accepted", validateAccepted)
	RegisterImplicit("AcceptedIf", validateAcceptedIf)
	Register("Declined", validateDeclined)
	Register("DeclinedIf", validateDeclinedIf)
}

// Passes if the value is "yes", "on", "1", "true", 1, or true.
func validateAccepted(_ string, value any, _ []string, _ RuleContext) bool {
	return isAcceptable(value)
}

// validateAcceptedIf passes when the value is accepted AND the other field
// equals one of the given values.
// Params: [otherField, value1, value2, ...]
func validateAcceptedIf(attribute string, value any, params []string, ctx RuleContext) bool {
	if len(params) < 2 {
		return true
	}

	other := ctx.GetValue(params[0])
	otherStr := stringify(other)

	for _, v := range params[1:] {
		if otherStr == v {
			return isAcceptable(value)
		}
	}

	return true // condition not met — rule doesn't apply
}

// validateDeclined passes if the value is "no", "off", "0", "false", 0, or false.
func validateDeclined(_ string, value any, _ []string, _ RuleContext) bool {
	return isDeclinable(value)
}

// validateDeclinedIf passes when the value is declined AND the other field
// equals one of the given values.
func validateDeclinedIf(attribute string, value any, params []string, ctx RuleContext) bool {
	if len(params) < 2 {
		return true
	}

	other := ctx.GetValue(params[0])
	otherStr := stringify(other)

	for _, v := range params[1:] {
		if otherStr == v {
			return isDeclinable(value)
		}
	}

	return true
}

func isAcceptable(value any) bool {
	switch v := value.(type) {
	case bool:
		return v
	case int:
		return v == 1
	case int64:
		return v == 1
	case float64:
		return v == 1
	case string:
		lower := strings.ToLower(strings.TrimSpace(v))

		return lower == "yes" || lower == "on" || lower == "1" || lower == "true"
	}

	return false
}

func isDeclinable(value any) bool {
	switch v := value.(type) {
	case bool:
		return !v
	case int:
		return v == 0
	case int64:
		return v == 0
	case float64:
		return v == 0
	case string:
		lower := strings.ToLower(strings.TrimSpace(v))

		return lower == "no" || lower == "off" || lower == "0" || lower == "false"
	}

	return false
}
