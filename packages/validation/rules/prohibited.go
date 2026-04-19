package rules

func init_prohibited() {
	RegisterImplicit("Prohibited", validateProhibited)
	RegisterImplicit("ProhibitedIf", validateProhibitedIf)
	RegisterImplicit("ProhibitedIfAccepted", validateProhibitedIfAccepted)
	RegisterImplicit("ProhibitedIfDeclined", validateProhibitedIfDeclined)
	RegisterImplicit("ProhibitedUnless", validateProhibitedUnless)
	RegisterImplicit("Prohibits", validateProhibits)
}

// validateProhibited: field must not be present (or must be blank/null).
func validateProhibited(attribute string, value any, _ []string, ctx RuleContext) bool {
	if !ctx.IsPresent(attribute) {
		return true
	}

	return isBlank(value)
}

// validateProhibitedIf: field is prohibited when otherField equals one of values.
func validateProhibitedIf(attribute string, value any, params []string, ctx RuleContext) bool {
	if len(params) < 2 {
		return true
	}

	other := ctx.GetValue(params[0])
	otherStr := stringify(other)

	for _, v := range params[1:] {
		if otherStr == v {
			return validateProhibited(attribute, value, nil, ctx)
		}
	}

	return true
}

// validateProhibitedIfAccepted: field is prohibited when otherField is accepted.
func validateProhibitedIfAccepted(attribute string, value any, params []string, ctx RuleContext) bool {
	if len(params) < 1 {
		return true
	}

	other := ctx.GetValue(params[0])

	if isAcceptable(other) {
		return validateProhibited(attribute, value, nil, ctx)
	}

	return true
}

// validateProhibitedIfDeclined: field is prohibited when otherField is declined.
func validateProhibitedIfDeclined(attribute string, value any, params []string, ctx RuleContext) bool {
	if len(params) < 1 {
		return true
	}

	other := ctx.GetValue(params[0])

	if isDeclinable(other) {
		return validateProhibited(attribute, value, nil, ctx)
	}

	return true
}

// validateProhibitedUnless: field is prohibited unless otherField is in values.
func validateProhibitedUnless(attribute string, value any, params []string, ctx RuleContext) bool {
	if len(params) < 2 {
		return true
	}

	other := ctx.GetValue(params[0])
	otherStr := stringify(other)

	for _, v := range params[1:] {
		if otherStr == v {
			return true // condition excludes prohibition
		}
	}

	return validateProhibited(attribute, value, nil, ctx)
}

// validateProhibits: if this field is present, the listed fields must be absent.
// Params: [otherField1, otherField2, ...]
func validateProhibits(attribute string, value any, params []string, ctx RuleContext) bool {
	if !ctx.IsPresent(attribute) || isBlank(value) {
		return true
	}

	for _, p := range params {
		if ctx.IsPresent(p) && isFilled(ctx.GetValue(p)) {
			return false
		}
	}

	return true
}
