package rules

import (
	"math"
	"strings"
)

func init_numeric() {
	Register("Digits", validateDigits)
	Register("DigitsBetween", validateDigitsBetween)
	Register("Decimal", validateDecimal)
	Register("MinDigits", validateMinDigits)
	Register("MaxDigits", validateMaxDigits)
	Register("MultipleOf", validateMultipleOf)
}

// validateDigits: value must be exactly n digits.
func validateDigits(_ string, value any, params []string, _ RuleContext) bool {
	s := stringify(value)

	// must be pure digits
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}

	n, ok := parseFloat64Param(params, 0)

	if !ok {
		return false
	}

	return float64(len(s)) == n
}

// validateDigitsBetween: digit count must be between min and max.
func validateDigitsBetween(_ string, value any, params []string, _ RuleContext) bool {
	if len(params) < 2 {
		return true
	}

	s := stringify(value)

	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}

	min, ok1 := parseFloat64Param(params, 0)
	max, ok2 := parseFloat64Param(params, 1)

	if !ok1 || !ok2 {
		return true
	}

	l := float64(len(s))

	return l >= min && l <= max
}

// validateDecimal: value must have exactly n decimal places (or between min,max).
// Params: [places] or [min,max]
func validateDecimal(_ string, value any, params []string, _ RuleContext) bool {
	s := stringify(value)
	s = strings.TrimSpace(s)

	// strip sign
	if len(s) > 0 && (s[0] == '-' || s[0] == '+') {
		s = s[1:]
	}

	parts := strings.Split(s, ".")

	if len(parts) != 2 {
		// No decimal point — 0 decimal places
		if len(params) == 1 {
			n, ok := parseFloat64Param(params, 0)

			return ok && n == 0
		}

		if len(params) >= 2 {
			min, ok1 := parseFloat64Param(params, 0)
			_, ok2 := parseFloat64Param(params, 1)

			return ok1 && ok2 && min == 0
		}

		return false
	}

	decimals := float64(len(parts[1]))

	if len(params) == 1 {
		n, ok := parseFloat64Param(params, 0)

		return ok && decimals == n
	}

	if len(params) >= 2 {
		min, ok1 := parseFloat64Param(params, 0)
		max, ok2 := parseFloat64Param(params, 1)

		return ok1 && ok2 && decimals >= min && decimals <= max
	}

	return false
}

// validateMinDigits: number of digits must be >= min.
func validateMinDigits(_ string, value any, params []string, _ RuleContext) bool {
	s := stringify(value)

	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}

	min, ok := parseFloat64Param(params, 0)

	if !ok {
		return true
	}

	return float64(len(s)) >= min
}

// validateMaxDigits: number of digits must be <= max.
func validateMaxDigits(_ string, value any, params []string, _ RuleContext) bool {
	s := stringify(value)

	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}

	max, ok := parseFloat64Param(params, 0)

	if !ok {
		return true
	}

	return float64(len(s)) <= max
}

// validateMultipleOf: value must be a multiple of divisor.
func validateMultipleOf(_ string, value any, params []string, _ RuleContext) bool {
	f, ok := toFloat64(value)

	if !ok {
		return false
	}

	divisor, ok := parseFloat64Param(params, 0)

	if !ok || divisor == 0 {
		return false
	}

	// Use remainder with tolerance for floating-point
	remainder := math.Mod(f, divisor)

	return math.Abs(remainder) < 1e-10 || math.Abs(remainder-divisor) < 1e-10
}
