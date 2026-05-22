package rules

import "strings"

func init_required() {
	RegisterImplicit("Required", validateRequired)
	RegisterImplicit("RequiredIf", validateRequiredIf)
	RegisterImplicit("RequiredIfAccepted", validateRequiredIfAccepted)
	RegisterImplicit("RequiredIfDeclined", validateRequiredIfDeclined)
	RegisterImplicit("RequiredUnless", validateRequiredUnless)
	RegisterImplicit("RequiredWith", validateRequiredWith)
	RegisterImplicit("RequiredWithAll", validateRequiredWithAll)
	RegisterImplicit("RequiredWithout", validateRequiredWithout)
	RegisterImplicit("RequiredWithoutAll", validateRequiredWithoutAll)
}

func validateRequired(attribute string, value any, _ []string, ctx RuleContext) bool {
	if !ctx.IsPresent(attribute) {
		return false
	}

	return isFilled(value)
}

// validateRequiredIf: field is required when otherField == one of values.
// Params: [otherField, value1, value2, ...]
func validateRequiredIf(attribute string, value any, params []string, ctx RuleContext) bool {
	if len(params) < 2 {
		return true
	}

	other := ctx.GetValue(params[0])
	otherStr := stringify(other)

	for _, v := range params[1:] {
		if otherStr == v {
			return validateRequired(attribute, value, nil, ctx)
		}
	}

	return true
}

// validateRequiredIfAccepted: field is required when otherField is accepted.
func validateRequiredIfAccepted(attribute string, value any, params []string, ctx RuleContext) bool {
	if len(params) < 1 {
		return true
	}

	other := ctx.GetValue(params[0])

	if !isAcceptable(other) {
		return true
	}

	return validateRequired(attribute, value, nil, ctx)
}

// validateRequiredIfDeclined: field is required when otherField is declined.
func validateRequiredIfDeclined(attribute string, value any, params []string, ctx RuleContext) bool {
	if len(params) < 1 {
		return true
	}

	other := ctx.GetValue(params[0])

	if !isDeclinable(other) {
		return true
	}

	return validateRequired(attribute, value, nil, ctx)
}

// validateRequiredUnless: field is required unless otherField is in values.
// Params: [otherField, value1, value2, ...]
func validateRequiredUnless(attribute string, value any, params []string, ctx RuleContext) bool {
	if len(params) < 2 {
		return true
	}

	other := ctx.GetValue(params[0])
	otherStr := stringify(other)

	for _, v := range params[1:] {
		if otherStr == v {
			return true // condition excludes this field
		}
	}

	return validateRequired(attribute, value, nil, ctx)
}

// validateRequiredWith: field is required when ANY of the other fields are present.
func validateRequiredWith(attribute string, value any, params []string, ctx RuleContext) bool {
	for _, p := range params {
		if ctx.IsPresent(p) && isFilled(ctx.GetValue(p)) {
			return validateRequired(attribute, value, nil, ctx)
		}
	}

	return true
}

// validateRequiredWithAll: field is required when ALL of the other fields are present.
func validateRequiredWithAll(attribute string, value any, params []string, ctx RuleContext) bool {
	allPresent := true

	for _, p := range params {
		if !ctx.IsPresent(p) || !isFilled(ctx.GetValue(p)) {
			allPresent = false

			break
		}
	}

	if allPresent {
		return validateRequired(attribute, value, nil, ctx)
	}

	return true
}

// validateRequiredWithout: field is required when ANY of the other fields are absent.
func validateRequiredWithout(attribute string, value any, params []string, ctx RuleContext) bool {
	for _, p := range params {
		if !ctx.IsPresent(p) || !isFilled(ctx.GetValue(p)) {
			return validateRequired(attribute, value, nil, ctx)
		}
	}

	return true
}

// validateRequiredWithoutAll: field is required when ALL of the other fields are absent.
func validateRequiredWithoutAll(attribute string, value any, params []string, ctx RuleContext) bool {
	allAbsent := true

	for _, p := range params {
		if ctx.IsPresent(p) && isFilled(ctx.GetValue(p)) {
			allAbsent = false

			break
		}
	}

	if allAbsent {
		return validateRequired(attribute, value, nil, ctx)
	}

	return true
}

// init_filled registers Filled (value must be present AND non-blank).
func init_filled() {
	RegisterImplicit("Filled", validateFilled)
}

func validateFilled(attribute string, value any, _ []string, ctx RuleContext) bool {
	if !ctx.IsPresent(attribute) {
		return true // not present is OK for Filled (unlike Required)
	}

	return isFilled(value)
}

// init_present registers Present family rules.
func init_present() {
	RegisterImplicit("Present", validatePresent)
	RegisterImplicit("PresentIf", validatePresentIf)
	RegisterImplicit("PresentUnless", validatePresentUnless)
	RegisterImplicit("PresentWith", validatePresentWith)
	RegisterImplicit("PresentWithAll", validatePresentWithAll)
}

func validatePresent(attribute string, _ any, _ []string, ctx RuleContext) bool {
	return ctx.IsPresent(attribute)
}

func validatePresentIf(attribute string, value any, params []string, ctx RuleContext) bool {
	if len(params) < 2 {
		return true
	}

	other := ctx.GetValue(params[0])
	otherStr := stringify(other)

	for _, v := range params[1:] {
		if otherStr == v {
			return validatePresent(attribute, value, nil, ctx)
		}
	}

	return true
}

func validatePresentUnless(attribute string, value any, params []string, ctx RuleContext) bool {
	if len(params) < 2 {
		return true
	}

	other := ctx.GetValue(params[0])
	otherStr := stringify(other)

	for _, v := range params[1:] {
		if otherStr == v {
			return true
		}
	}

	return validatePresent(attribute, value, nil, ctx)
}

func validatePresentWith(attribute string, value any, params []string, ctx RuleContext) bool {
	for _, p := range params {
		if ctx.IsPresent(p) {
			return validatePresent(attribute, value, nil, ctx)
		}
	}

	return true
}

func validatePresentWithAll(attribute string, value any, params []string, ctx RuleContext) bool {
	for _, p := range params {
		if !ctx.IsPresent(p) {
			return true
		}
	}

	return validatePresent(attribute, value, nil, ctx)
}

// init_missing registers Missing family rules.
func init_missing() {
	RegisterImplicit("Missing", validateMissing)
	RegisterImplicit("MissingIf", validateMissingIf)
	RegisterImplicit("MissingUnless", validateMissingUnless)
	RegisterImplicit("MissingWith", validateMissingWith)
	RegisterImplicit("MissingWithAll", validateMissingWithAll)
}

func validateMissing(attribute string, _ any, _ []string, ctx RuleContext) bool {
	return !ctx.IsPresent(attribute)
}

func validateMissingIf(attribute string, value any, params []string, ctx RuleContext) bool {
	if len(params) < 2 {
		return true
	}

	other := ctx.GetValue(params[0])
	otherStr := stringify(other)

	for _, v := range params[1:] {
		if otherStr == v {
			return validateMissing(attribute, value, nil, ctx)
		}
	}

	return true
}

func validateMissingUnless(attribute string, value any, params []string, ctx RuleContext) bool {
	if len(params) < 2 {
		return true
	}

	other := ctx.GetValue(params[0])
	otherStr := stringify(other)

	for _, v := range params[1:] {
		if otherStr == v {
			return true
		}
	}

	return validateMissing(attribute, value, nil, ctx)
}

func validateMissingWith(attribute string, value any, params []string, ctx RuleContext) bool {
	for _, p := range params {
		if ctx.IsPresent(p) {
			return validateMissing(attribute, value, nil, ctx)
		}
	}

	return true
}

func validateMissingWithAll(attribute string, value any, params []string, ctx RuleContext) bool {
	for _, p := range params {
		if !ctx.IsPresent(p) {
			return true
		}
	}

	return validateMissing(attribute, value, nil, ctx)
}

// sometimes is a no-op marker rule handled by the Validator directly.
func validateSometimes(_ string, _ any, _ []string, _ RuleContext) bool {
	return true
}

// nullable is a no-op marker rule — the validator skips other rules when the
// value is nil/null.
func validateNullable(_ string, _ any, _ []string, _ RuleContext) bool {
	return true
}

// bail is a no-op marker rule — the validator stops on first failure for that field.
func validateBail(_ string, _ any, _ []string, _ RuleContext) bool {
	return true
}

// init registers all rules in this file.
func init_required_markers() {
	RegisterImplicit("Sometimes", validateSometimes)
	Register("Nullable", validateNullable)
	Register("Bail", validateBail)
}

// Helper: allParamsInOther checks whether all params values exist in a
// comma-separated "other" field value.
func allParamsIn(params []string, values []string) bool {
	for _, p := range params {
		if !containsString(values, strings.TrimSpace(p)) {
			return false
		}
	}

	return true
}
