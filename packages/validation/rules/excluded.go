package rules

func init_excluded() {
	// Exclude rules are marker rules — the validator removes the field from
	// the validated output.  They always "pass" (return true) since they do not
	// produce error messages; the field is simply excluded.
	RegisterImplicit("Exclude", validateExclude)
	RegisterImplicit("ExcludeIf", validateExcludeIf)
	RegisterImplicit("ExcludeUnless", validateExcludeUnless)
	RegisterImplicit("ExcludeWith", validateExcludeWith)
	RegisterImplicit("ExcludeWithout", validateExcludeWithout)
}

func validateExclude(_ string, _ any, _ []string, _ RuleContext) bool {
	return true // handled as a special case in the validator; no error message
}

func validateExcludeIf(_ string, _ any, _ []string, _ RuleContext) bool {
	return true
}

func validateExcludeUnless(_ string, _ any, _ []string, _ RuleContext) bool {
	return true
}

func validateExcludeWith(_ string, _ any, _ []string, _ RuleContext) bool {
	return true
}

func validateExcludeWithout(_ string, _ any, _ []string, _ RuleContext) bool {
	return true
}

// ShouldExclude evaluates whether an attribute should be excluded from the
// validated output given its rules, data, and context.
func ShouldExclude(attribute string, parsedRuleName string, params []string, ctx RuleContext) bool {
	switch parsedRuleName {
	case "Exclude":
		return true
	case "ExcludeIf":
		if len(params) < 2 {
			return false
		}

		other := ctx.GetValue(params[0])
		otherStr := stringify(other)

		for _, v := range params[1:] {
			if otherStr == v {
				return true
			}
		}

		return false
	case "ExcludeUnless":
		if len(params) < 2 {
			return false
		}

		other := ctx.GetValue(params[0])
		otherStr := stringify(other)

		for _, v := range params[1:] {
			if otherStr == v {
				return false // condition says keep it
			}
		}

		return true // not in the unless list → exclude
	case "ExcludeWith":
		for _, p := range params {
			if ctx.IsPresent(p) && isFilled(ctx.GetValue(p)) {
				return true
			}
		}

		return false
	case "ExcludeWithout":
		for _, p := range params {
			if !ctx.IsPresent(p) || !isFilled(ctx.GetValue(p)) {
				return true
			}
		}

		return false
	}

	return false
}
