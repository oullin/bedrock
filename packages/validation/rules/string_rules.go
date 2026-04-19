package rules

import (
	"strings"
	"unicode"
)

func init_string() {
	Register("Alpha", validateAlpha)
	Register("AlphaDash", validateAlphaDash)
	Register("AlphaNum", validateAlphaNum)
	Register("Ascii", validateAscii)
	Register("Lowercase", validateLowercase)
	Register("Uppercase", validateUppercase)
	Register("StartsWith", validateStartsWith)
	Register("DoesntStartWith", validateDoesntStartWith)
	Register("EndsWith", validateEndsWith)
	Register("DoesntEndWith", validateDoesntEndWith)
	Register("HexColor", validateHexColor)
}

func validateAlpha(_ string, value any, _ []string, _ RuleContext) bool {
	s, ok := value.(string)

	if !ok {
		return false
	}

	if s == "" {
		return false
	}

	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}

	return true
}

func validateAlphaDash(_ string, value any, _ []string, _ RuleContext) bool {
	s, ok := value.(string)

	if !ok {
		return false
	}

	if s == "" {
		return false
	}

	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '_' {
			return false
		}
	}

	return true
}

func validateAlphaNum(_ string, value any, _ []string, _ RuleContext) bool {
	s, ok := value.(string)

	if !ok {
		return false
	}

	if s == "" {
		return false
	}

	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}

func validateAscii(_ string, value any, _ []string, _ RuleContext) bool {
	s, ok := value.(string)

	if !ok {
		return false
	}

	for _, r := range s {
		if r > 127 {
			return false
		}
	}

	return true
}

func validateLowercase(_ string, value any, _ []string, _ RuleContext) bool {
	s, ok := value.(string)

	if !ok {
		return false
	}

	return s == strings.ToLower(s)
}

func validateUppercase(_ string, value any, _ []string, _ RuleContext) bool {
	s, ok := value.(string)

	if !ok {
		return false
	}

	return s == strings.ToUpper(s)
}

func validateStartsWith(_ string, value any, params []string, _ RuleContext) bool {
	s, ok := value.(string)

	if !ok {
		return false
	}

	for _, p := range params {
		if strings.HasPrefix(s, p) {
			return true
		}
	}

	return false
}

func validateDoesntStartWith(_ string, value any, params []string, _ RuleContext) bool {
	s, ok := value.(string)

	if !ok {
		return true
	}

	for _, p := range params {
		if strings.HasPrefix(s, p) {
			return false
		}
	}

	return true
}

func validateEndsWith(_ string, value any, params []string, _ RuleContext) bool {
	s, ok := value.(string)

	if !ok {
		return false
	}

	for _, p := range params {
		if strings.HasSuffix(s, p) {
			return true
		}
	}

	return false
}

func validateDoesntEndWith(_ string, value any, params []string, _ RuleContext) bool {
	s, ok := value.(string)

	if !ok {
		return true
	}

	for _, p := range params {
		if strings.HasSuffix(s, p) {
			return false
		}
	}

	return true
}

func validateHexColor(_ string, value any, _ []string, _ RuleContext) bool {
	s, ok := value.(string)

	if !ok {
		return false
	}

	if len(s) == 0 || s[0] != '#' {
		return false
	}

	hex := s[1:]

	if len(hex) != 3 && len(hex) != 4 && len(hex) != 6 && len(hex) != 8 {
		return false
	}

	for _, r := range hex {
		if !((r >= '0' && r <= '9') || (r >= 'a' && r <= 'f') || (r >= 'A' && r <= 'F')) {
			return false
		}
	}

	return true
}
