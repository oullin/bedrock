package rules

import "unicode"

// validatePassword validates a password against minimum requirements.
// This is the basic form. The full Password rule object (PasswordRule) is in
// rule.go and allows chaining (Min, Letters, Numbers, Symbols, Uncompromised).
// When used as a bare string rule "password" it just checks minimum length 8.

// PasswordOptions holds the requirements for the PasswordRule object.
type PasswordOptions struct {
	Min           int
	Max           int // 0 = no limit
	Letters       bool
	Numbers       bool
	Symbols       bool
	Mixed         bool // mixed case required
	Uncompromised bool
}

func init_password() {
	Register("Password", validatePassword)
}

func validatePassword(_ string, value any, params []string, _ RuleContext) bool {
	s, ok := value.(string)

	if !ok {
		return false
	}

	minLen := 8

	if len(params) > 0 {
		if n, ok2 := parseFloat64Param(params, 0); ok2 {
			minLen = int(n)
		}
	}

	return len([]rune(s)) >= minLen
}

// CheckPassword validates a password against a set of options.
func CheckPassword(s string, opts PasswordOptions) bool {
	runes := []rune(s)

	if opts.Min > 0 && len(runes) < opts.Min {
		return false
	}

	if opts.Max > 0 && len(runes) > opts.Max {
		return false
	}

	if opts.Letters {
		hasLetter := false

		for _, r := range runes {
			if unicode.IsLetter(r) {
				hasLetter = true

				break
			}
		}

		if !hasLetter {
			return false
		}
	}

	if opts.Numbers {
		hasDigit := false

		for _, r := range runes {
			if unicode.IsDigit(r) {
				hasDigit = true

				break
			}
		}

		if !hasDigit {
			return false
		}
	}

	if opts.Symbols {
		hasSymbol := false

		for _, r := range runes {
			if !unicode.IsLetter(r) && !unicode.IsDigit(r) && !unicode.IsSpace(r) {
				hasSymbol = true

				break
			}
		}

		if !hasSymbol {
			return false
		}
	}

	if opts.Mixed {
		hasUpper, hasLower := false, false

		for _, r := range runes {
			if unicode.IsUpper(r) {
				hasUpper = true
			}

			if unicode.IsLower(r) {
				hasLower = true
			}
		}

		if !hasUpper || !hasLower {
			return false
		}
	}

	return true
}
