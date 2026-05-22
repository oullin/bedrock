package str

import (
	"strings"

	"github.com/jinzhu/inflection"
)

// pluralOverrides corrects jinzhu/inflection words that are mishandled
// by overly broad suffix rules (e.g. "man" → "men" matches "human").
var pluralOverrides = map[string]string{
	"human": "humans",
}

// StrPlural returns the plural form of the given word.
// If count is provided and equals 1, the singular form is returned.
// Note: English-only. Non-English pluralization is not supported.
func StrPlural(value string, count ...int) string {
	if len(count) > 0 && count[0] == 1 {
		return value
	}

	lower := strings.ToLower(value)

	if override, ok := pluralOverrides[lower]; ok {
		// Preserve original capitalisation
		if len(value) > 0 && strings.ToUpper(value[:1]) == value[:1] {
			return strings.ToUpper(override[:1]) + override[1:]
		}

		return override
	}

	return inflection.Plural(value)
}

// StrSingular returns the singular form of the given word.
// Note: English-only via go-flect.
func StrSingular(value string) string {
	return inflection.Singular(value)
}

// StrPluralStudly pluralizes the last "studly" word in the string.
func StrPluralStudly(value string, count ...int) string {
	if len(count) > 0 && count[0] == 1 {
		return value
	}

	// Split studly words
	parts := StrUcsplit(value)

	if len(parts) == 0 {
		return value
	}

	// Pluralize last part
	last := parts[len(parts)-1]
	parts[len(parts)-1] = StrStudly(StrPlural(strings.ToLower(last)))

	return strings.Join(parts, "")
}

// StrPluralPascal pluralizes the last "pascal" word in the string.
func StrPluralPascal(value string, count ...int) string {
	return StrPluralStudly(value, count...)
}
