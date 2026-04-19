package rules

import "regexp"

func init_regex() {
	Register("Regex", validateRegex)
	Register("NotRegex", validateNotRegex)
}

func validateRegex(_ string, value any, params []string, _ RuleContext) bool {
	if len(params) == 0 {
		return true
	}

	s, ok := value.(string)

	if !ok {
		s = stringify(value)
	}

	pattern := params[0]

	// Strip surrounding slashes if present (e.g. /^[a-z]+$/)
	if len(pattern) >= 2 && pattern[0] == '/' && pattern[len(pattern)-1] == '/' {
		pattern = pattern[1 : len(pattern)-1]
	}

	re, err := regexp.Compile(pattern)

	if err != nil {
		return false
	}

	return re.MatchString(s)
}

func validateNotRegex(_ string, value any, params []string, ctx RuleContext) bool {
	return !validateRegex("", value, params, ctx)
}
