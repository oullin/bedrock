package support

import "strings"

var irregularPlurals = map[string]string{
	"child":  "children",
	"person": "people",
	"mouse":  "mice",
}

var irregularSingulars = map[string]string{
	"children": "child",
	"people":   "person",
	"mice":     "mouse",
}

// Plural returns a simple English plural form.
func Plural(word string, count ...int) string {
	if len(count) > 0 && count[0] == 1 {
		return word
	}

	if word == "" || !isAlphaNumeric(word[len(word)-1]) {
		return word
	}

	lower := strings.ToLower(word)

	if plural, ok := irregularPlurals[lower]; ok {
		return preserveCase(word, plural)
	}

	switch {
	case strings.HasSuffix(lower, "y") && len(word) > 1 && !isVowel(lower[len(lower)-2]):
		return word[:len(word)-1] + preserveCase(word[len(word)-1:], "ies")
	case strings.HasSuffix(lower, "s"), strings.HasSuffix(lower, "x"), strings.HasSuffix(lower, "ch"), strings.HasSuffix(lower, "sh"):
		return word + "es"
	default:
		return word + "s"
	}
}

// Singular returns a simple English singular form.
func Singular(word string) string {
	lower := strings.ToLower(word)

	if singular, ok := irregularSingulars[lower]; ok {
		return preserveCase(word, singular)
	}

	switch {
	case strings.HasSuffix(lower, "ies") && len(word) > 3:
		return word[:len(word)-3] + preserveCase(word[len(word)-3:], "y")
	case strings.HasSuffix(lower, "es") && (strings.HasSuffix(lower, "ses") || strings.HasSuffix(lower, "xes") || strings.HasSuffix(lower, "ches") || strings.HasSuffix(lower, "shes")):
		return word[:len(word)-2]
	case strings.HasSuffix(lower, "s") && len(word) > 1:
		return word[:len(word)-1]
	default:
		return word
	}
}

// PluralStudly pluralizes the last StudlyCase word segment.
func PluralStudly(word string, count ...int) string {
	if word == "" {
		return word
	}

	start := len(word) - 1

	for start > 0 && (word[start] < 'A' || word[start] > 'Z') {
		start--
	}

	if start == len(word)-1 {
		return Plural(word, count...)
	}

	return word[:start] + Plural(word[start:], count...)
}

func isAlphaNumeric(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9')
}

func isVowel(b byte) bool {
	return strings.ContainsRune("aeiou", rune(b))
}

func preserveCase(original, replacement string) string {
	if original == strings.ToUpper(original) {
		return strings.ToUpper(replacement)
	}

	if len(original) > 0 && original[0] >= 'A' && original[0] <= 'Z' {
		return strings.ToUpper(replacement[:1]) + replacement[1:]
	}

	return replacement
}
