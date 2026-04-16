package helpers

import "github.com/bedrock/packages/support"

// CamelCase converts a string to camelCase.
// Delegates to support.StrCamel. Mirrors Upstream's camel_case().
func CamelCase(value string) string { return support.StrCamel(value) }

// EndsWith determines if a string ends with any of the given suffixes.
// Delegates to support.StrEndsWith. Mirrors Upstream's ends_with().
func EndsWith(subject string, suffixes ...string) bool {
	return support.StrEndsWith(subject, suffixes...)
}

// KebabCase converts a string to kebab-case.
// Delegates to support.StrKebab. Mirrors Upstream's kebab_case().
func KebabCase(value string) string { return support.StrKebab(value) }

// SnakeCase converts a string to snake_case.
// Delegates to support.StrSnake. Mirrors Upstream's snake_case().
func SnakeCase(value string, delimiter ...string) string {
	return support.StrSnake(value, delimiter...)
}

// StartsWith determines if a string starts with any of the given prefixes.
// Delegates to support.StrStartsWith. Mirrors Upstream's starts_with().
func StartsWith(subject string, prefixes ...string) bool {
	return support.StrStartsWith(subject, prefixes...)
}

// StrAfter returns the portion after the first occurrence of search.
// Delegates to support.StrAfter. Mirrors Upstream's str_after().
func StrAfter(subject, search string) string { return support.StrAfter(subject, search) }

// StrBefore returns the portion before the first occurrence of search.
// Delegates to support.StrBefore. Mirrors Upstream's str_before().
func StrBefore(subject, search string) string { return support.StrBefore(subject, search) }

// StrContains determines if a string contains any of the given substrings.
// Delegates to support.StrContains. Mirrors Upstream's str_contains().
func StrContains(haystack string, needles ...string) bool {
	return support.StrContains(haystack, needles...)
}

// StrFinish appends a single instance of the given cap to the string.
// Delegates to support.StrFinish. Mirrors Upstream's str_finish().
func StrFinish(value, cap string) string { return support.StrFinish(value, cap) }

// StrIs determines if a string matches a given pattern.
// Delegates to support.StrIs. Mirrors Upstream's str_is().
func StrIs(pattern, value string) bool { return support.StrIs(pattern, value) }

// StrLimit truncates a string to the given length.
// Delegates to support.StrLimit. Mirrors Upstream's str_limit().
func StrLimit(value string, limit int, end ...string) string {
	return support.StrLimit(value, limit, end...)
}

// StrPlural returns the plural form of the given word.
// Delegates to support.StrPlural. Mirrors Upstream's str_plural().
func StrPlural(value string, count ...int) string { return support.StrPlural(value, count...) }

// StrRandom generates a random alpha-numeric string.
// Delegates to support.StrRandom. Mirrors Upstream's str_random().
func StrRandom(length ...int) string { return support.StrRandom(length...) }

// StrReplaceArray sequentially replaces occurrences of search.
// Delegates to support.StrReplaceArray. Mirrors Upstream's str_replace_array().
func StrReplaceArray(search string, replacements []string, subject string) string {
	return support.StrReplaceArray(search, replacements, subject)
}

// StrReplaceFirst replaces the first occurrence of search.
// Delegates to support.StrReplaceFirst. Mirrors Upstream's str_replace_first().
func StrReplaceFirst(search, replace, subject string) string {
	return support.StrReplaceFirst(search, replace, subject)
}

// StrReplaceLast replaces the last occurrence of search.
// Delegates to support.StrReplaceLast. Mirrors Upstream's str_replace_last().
func StrReplaceLast(search, replace, subject string) string {
	return support.StrReplaceLast(search, replace, subject)
}

// StrSingular returns the singular form of the given word.
// Delegates to support.StrSingular. Mirrors Upstream's str_singular().
func StrSingular(value string) string { return support.StrSingular(value) }

// StrSlug generates a URL-friendly slug.
// Delegates to support.StrSlug. Mirrors Upstream's str_slug().
func StrSlug(title string, separator ...string) string {
	return support.StrSlug(title, separator...)
}

// StrStart prepends a single instance of the given prefix.
// Delegates to support.StrStart. Mirrors Upstream's str_start().
func StrStart(value, prefix string) string { return support.StrStart(value, prefix) }

// StudlyCase converts a string to StudlyCase (PascalCase).
// Delegates to support.StrStudly. Mirrors Upstream's studly_case().
func StudlyCase(value string) string { return support.StrStudly(value) }

// TitleCase converts a string to Title Case.
// Delegates to support.StrTitle. Mirrors Upstream's title_case().
func TitleCase(value string) string { return support.StrTitle(value) }
