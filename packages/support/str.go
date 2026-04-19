package support

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"unicode"
	"unicode/utf8"
)

// randomMu guards package-level random string factory state.

// Casing caches (camel, studly, snake)

// StrAfter returns the portion of the string after the first occurrence of search.
// If search is not found, the full subject is returned.
// Mirrors Str::after().

// StrAfterLast returns the portion after the last occurrence of search.
// Mirrors Str::afterLast().

// StrBefore returns the portion before the first occurrence of search.
// Mirrors Str::before().

// StrBeforeLast returns the portion before the last occurrence of search.
// Mirrors Str::beforeLast().

// StrBetween returns the portion between from and to.
// Mirrors Str::between().

// StrBetweenFirst returns the smallest portion between from and to.
// Mirrors Str::betweenFirst().

// StrCamel converts a string to camelCase.
// Mirrors Str::camel().

// StrStudly converts a string to StudlyCase (PascalCase).
// Mirrors Str::studly().

// Replace common separators with spaces, then title-case

// studlySplit splits a string on separators (-, _, space) and camel-case boundaries.

// Replace hyphens, underscores with spaces

// StrPascal converts a string to PascalCase (alias for Studly).
// Mirrors Str::pascal().

// StrSnake converts a string to snake_case.
// Mirrors Str::snake().

// Insert separator before uppercase letters that follow lowercase or digits

// StrKebab converts a string to kebab-case.
// Mirrors Str::kebab().

// StrTitle converts a string to Title Case.
// Mirrors Str::title().

//nolint:staticcheck

// StrHeadline converts a string to Headline Case.
// Each word is capitalized, separating camelCase and underscored strings.
// Mirrors Str::headline().

// Handle camelCase by inserting spaces

// Replace separators

// Title case each word

// StrApa converts a string to APA title case.
// Mirrors Str::apa().

// APA: lowercase articles, conjunctions, short prepositions

// Capitalize

// StrUcfirst uppercases the first character of the string.
// Mirrors Str::ucfirst().

// StrLcfirst lowercases the first character of the string.
// Mirrors Str::lcfirst().

// StrUcwords uppercases the first character of each word.
// Mirrors Str::ucwords().

//nolint:staticcheck

// Word-boundary title case with custom separators

// StrUcsplit splits a string by uppercase characters.
// Mirrors Str::ucsplit().

// StrContains determines if the haystack contains the given needle(s).
// Case-sensitive by default.
// Mirrors Str::contains().

// StrContainsAll determines if the haystack contains all given needles.
// Mirrors Str::containsAll().

// StrDoesntContain determines if the haystack does not contain the given needle(s).
// Mirrors Str::doesntContain().

// StrStartsWith determines if the subject starts with the given prefix(es).
// Mirrors Str::startsWith().

// StrDoesntStartWith determines if the subject does not start with any of the given values.
// Mirrors Str::doesntStartWith().

// StrEndsWith determines if the subject ends with the given suffix(es).
// Mirrors Str::endsWith().

// StrDoesntEndWith determines if the subject does not end with any of the given values.
// Mirrors Str::doesntEndWith().

// StrIs determines if the value matches the given pattern.
// An asterisk (*) may be used as a wildcard value.
// Mirrors Str::is().

// Convert glob pattern to regex

// StrIsMatch determines if the value matches any of the given patterns (regex).
// Mirrors Str::isMatch().

// StrIsAscii determines if a string is 7-bit ASCII.
// Mirrors Str::isAscii().

// StrIsJson determines if the given string is valid JSON.
// Mirrors Str::isJson().

// StrIsUrl determines if the given string is a valid URL.
// Mirrors Str::isUrl().

// StrIsUuid determines if the given string is a valid UUID.
// Mirrors Str::isUuid().

// StrIsUlid determines if the given string is a valid ULID.
// Mirrors Str::isUlid().

// StrLength returns the length of the string in runes (characters).
// Mirrors Str::length().

// StrLimit truncates the string to the given number of characters.
// Mirrors Str::limit().

// StrWords limits the number of words in a string.
// Mirrors Str::words().

// Split preserving spaces around words

// StrLower converts the string to lowercase.
// Mirrors Str::lower().

// StrUpper converts the string to uppercase.
// Mirrors Str::upper().

// StrReverse reverses the string (UTF-8 aware).
// Mirrors Str::reverse().

// StrReplace replaces occurrences of search in subject with replace.
// Mirrors Str::replace().

// StrReplaceArray replaces each search occurrence with the next value from replacements.
// Mirrors Str::replaceArray().

// StrReplaceFirst replaces the first occurrence of search in subject.
// Mirrors Str::replaceFirst().

// StrReplaceLast replaces the last occurrence of search in subject.
// Mirrors Str::replaceLast().

// StrReplaceStart replaces the search value if it appears at the start of subject.
// Mirrors Str::replaceStart().

// StrReplaceEnd replaces the search value if it appears at the end of subject.
// Mirrors Str::replaceEnd().

// StrReplaceMatches replaces all occurrences of a regex pattern in subject.
// Mirrors Str::replaceMatches().

// StrRemove removes all occurrences of search from subject.
// Mirrors Str::remove().

// StrRepeat repeats the string the given number of times.
// Mirrors Str::repeat().

// StrSquish removes all extra blank space (including between words).
// Mirrors Str::squish().

// StrDeduplicate replaces consecutive occurrences of characters with a single instance.
// Mirrors Str::deduplicate().

// StrTrim removes leading and trailing whitespace (or given chars).
// Mirrors Str::trim().

// StrLtrim removes leading whitespace (or given chars).
// Mirrors Str::ltrim().

// StrRtrim removes trailing whitespace (or given chars).
// Mirrors Str::rtrim().

// StrStart prepends a single instance of value to the subject
// if it does not already start with it.
// Mirrors Str::start().

// StrFinish appends a single instance of cap to value
// if it does not already end with it.
// Mirrors Str::finish().

// StrWrap wraps the string with the given strings.
// Mirrors Str::wrap().

// StrUnwrap removes the given prefix and suffix from the string.
// Mirrors Str::unwrap().

// StrChopStart removes a single leading occurrence of needle (or the first matching needle).
// Mirrors Str::chopStart().

// StrChopEnd removes a single trailing occurrence of needle (or the first matching needle).
// Mirrors Str::chopEnd().

// StrMask replaces a portion of the string with a repeated character.
// Mirrors Str::mask().

// StrExcerpt returns a string excerpt around the given phrase.
// Mirrors Str::excerpt().

// StrMatch returns the first regex match in subject.
// Mirrors Str::match().

// StrMatchAll returns all regex matches in subject.
// Mirrors Str::matchAll().

// StrNumbers extracts all numeric characters from the string.
// Mirrors Str::numbers().

// StrPadBoth pads the string on both sides to the given length.
// Mirrors Str::padBoth().

// StrPadLeft pads the left side of the string to the given length.
// Mirrors Str::padLeft().

// StrPadRight pads the right side of the string to the given length.
// Mirrors Str::padRight().

// StrPosition finds the position of needle in haystack.
// Returns the byte position and whether it was found.
// Mirrors Str::position().

// StrSubstr returns a substring from start for the given length.
// Mirrors Str::substr().

// StrSubstrCount counts the occurrences of needle in haystack.
// Mirrors Str::substrCount().

// StrSubstrReplace replaces a portion of the string.
// Mirrors Str::substrReplace().

// StrSwap swaps multiple keyword pairs in the subject.
// Mirrors Str::swap().

// Build patterns ordered by length (longer first to avoid partial replacements)

// StrTake returns the first n characters, or from the end if n is negative.
// Mirrors Str::take().

// negative: from end

// StrWordCount counts the number of words in the string.
// Mirrors Str::wordCount().

// StrWordWrap wraps the string to the given number of characters.
// Mirrors Str::wordWrap().

// StrToBase64 encodes the string to base64.
// Mirrors Str::toBase64().

// StrFromBase64 decodes a base64-encoded string.
// Mirrors Str::fromBase64().

// Try URL-safe base64

// StrSlug generates a URL-friendly slug from the given string.
// Mirrors Str::slug().

// Convert to ASCII

// Replace @ with separator (common convention)

// Replace non-alphanumeric with separator

// Remove leading/trailing separators and deduplicate

// StrAscii transliterates the string to its closest ASCII representation.
// Mirrors Str::ascii().

// Unknown non-ASCII chars are dropped (like Upstream's behavior)

// StrTransliterate converts the string to its closest ASCII form.
// Mirrors Str::transliterate().

// StrInitials returns the initials of the given name.
// Mirrors Str::initials().

// StrCharAt returns the character at the given index (UTF-8 aware).
// Mirrors Str::charAt().

// StrParseCallback parses a "Class@method" or "Class::method" callback string.
// Returns [class, method].
// Mirrors Str::parseCallback().

// StrRandom generates a random alphanumeric string of the given length.
// Mirrors Str::random().

// fallback used when sequence exhausted

// CreateRandomStringsUsing sets a custom random string factory.
// Mirrors Str::createRandomStringsUsing().

// CreateRandomStringsUsingSequence sets a sequence of strings to use for random generation.
// Mirrors Str::createRandomStringsUsingSequence().

// CreateRandomStringsNormally resets random string generation to default.
// Mirrors Str::createRandomStringsNormally().

// StrPassword generates a secure random password.
// Mirrors Str::password().

// FlushCache clears the casing caches.
// Mirrors Str::flushCache().

// StrConvertCase converts the string case using the given mode.
// Mirrors Str::convertCase().

// MB_CASE_UPPER

// MB_CASE_LOWER

// MB_CASE_TITLE
//nolint:staticcheck

// Of creates a new StringBuilder for fluent string manipulation.
// Mirrors Str::of().

// StringBuilder provides a fluent interface for string manipulation.
// Mirrors Framework\Support\Stringable.
type StringBuilder struct {
	value string
}

var (
	randomMu       sync.Mutex
	randomFactory  func(int) string
	randomSequence []string
	randomFallback func(int) string

	camelCache  sync.Map
	studlyCache sync.Map
	snakeCache  sync.Map
)

func StrAfter(subject, search string) string {
	if search == "" {
		return subject
	}

	idx := strings.Index(subject, search)

	if idx == -1 {
		return subject
	}

	return subject[idx+len(search):]
}

func StrAfterLast(subject, search string) string {
	if search == "" {
		return subject
	}

	idx := strings.LastIndex(subject, search)

	if idx == -1 {
		return subject
	}

	return subject[idx+len(search):]
}

func StrBefore(subject, search string) string {
	if search == "" {
		return subject
	}

	idx := strings.Index(subject, search)

	if idx == -1 {
		return subject
	}

	return subject[:idx]
}

func StrBeforeLast(subject, search string) string {
	if search == "" {
		return subject
	}

	idx := strings.LastIndex(subject, search)

	if idx == -1 {
		return subject
	}

	return subject[:idx]
}

func StrBetween(subject, from, to string) string {
	if from == "" || to == "" {
		return subject
	}

	return StrBeforeLast(StrAfter(subject, from), to)
}

func StrBetweenFirst(subject, from, to string) string {
	if from == "" || to == "" {
		return subject
	}

	return StrBefore(StrAfter(subject, from), to)
}

func StrCamel(value string) string {
	if cached, ok := camelCache.Load(value); ok {
		return cached.(string)
	}

	result := strLcfirst(StrStudly(value))
	camelCache.Store(value, result)

	return result
}

func StrStudly(value string) string {
	if cached, ok := studlyCache.Load(value); ok {
		return cached.(string)
	}

	words := studlySplit(value)

	var builder strings.Builder

	for _, w := range words {
		if len(w) > 0 {
			r, size := utf8.DecodeRuneInString(w)
			builder.WriteRune(unicode.ToUpper(r))
			builder.WriteString(w[size:])
		}
	}

	result := builder.String()
	studlyCache.Store(value, result)

	return result
}

func studlySplit(value string) []string {

	value = regexp.MustCompile(`[-_\s]+`).ReplaceAllString(value, " ")

	return strings.Fields(value)
}

func StrPascal(value string) string {
	return StrStudly(value)
}

func StrSnake(value string, delimiter ...string) string {
	sep := "_"

	if len(delimiter) > 0 {
		sep = delimiter[0]
	}

	key := value + sep

	if cached, ok := snakeCache.Load(key); ok {
		return cached.(string)
	}

	var result strings.Builder
	runes := []rune(value)

	for i, r := range runes {
		if unicode.IsUpper(r) {
			if i > 0 && (unicode.IsLower(runes[i-1]) || unicode.IsDigit(runes[i-1])) {
				result.WriteString(sep)
			} else if i > 0 && unicode.IsUpper(runes[i-1]) && i+1 < len(runes) && unicode.IsLower(runes[i+1]) {
				result.WriteString(sep)
			}

			result.WriteRune(unicode.ToLower(r))
		} else if r == '-' || r == ' ' {
			result.WriteString(sep)
		} else {
			result.WriteRune(r)
		}
	}

	res := strings.Trim(result.String(), sep)
	snakeCache.Store(key, res)

	return res
}

func StrKebab(value string) string {
	return StrSnake(value, "-")
}

func StrTitle(value string) string {
	return strings.Title(strings.ToLower(value))
}

func StrHeadline(value string) string {

	re := regexp.MustCompile(`([a-z])([A-Z])`)
	value = re.ReplaceAllString(value, "$1 $2")

	value = regexp.MustCompile(`[-_]+`).ReplaceAllString(value, " ")

	words := strings.Fields(value)

	for i, w := range words {
		if len(w) == 0 {
			continue
		}

		r, size := utf8.DecodeRuneInString(w)
		words[i] = string(unicode.ToUpper(r)) + strings.ToLower(w[size:])
	}

	return strings.Join(words, " ")
}

func StrApa(value string) string {

	minors := map[string]bool{
		"a": true, "an": true, "the": true,
		"and": true, "but": true, "for": true, "nor": true, "or": true, "so": true, "yet": true,
		"as": true, "at": true, "by": true, "in": true, "of": true, "on": true, "to": true, "up": true,
		"via": true,
	}

	words := strings.Fields(value)

	for i, w := range words {
		lower := strings.ToLower(w)

		if i == 0 || i == len(words)-1 || !minors[lower] {

			if len(w) > 0 {
				r, size := utf8.DecodeRuneInString(w)
				words[i] = string(unicode.ToUpper(r)) + w[size:]
			}
		} else {
			words[i] = lower
		}
	}

	return strings.Join(words, " ")
}

func StrUcfirst(value string) string {
	if value == "" {
		return value
	}

	r, size := utf8.DecodeRuneInString(value)

	return string(unicode.ToUpper(r)) + value[size:]
}

func StrLcfirst(value string) string {
	return strLcfirst(value)
}

func strLcfirst(value string) string {
	if value == "" {
		return value
	}

	r, size := utf8.DecodeRuneInString(value)

	return string(unicode.ToLower(r)) + value[size:]
}

func StrUcwords(value string, separators ...string) string {
	seps := " \t\r\n\f\v"

	if len(separators) > 0 {
		seps = separators[0]
	}

	if seps == " " || seps == " \t\r\n\f\v" {
		return strings.Title(value)
	}

	runes := []rune(value)
	newWord := true

	for i, r := range runes {
		if strings.ContainsRune(seps, r) {
			newWord = true
		} else if newWord {
			runes[i] = unicode.ToUpper(r)
			newWord = false
		}
	}

	return string(runes)
}

func StrUcsplit(value string) []string {
	re := regexp.MustCompile(`(?m)[A-Z][^A-Z]*`)
	matches := re.FindAllString(value, -1)

	if len(matches) == 0 {
		return []string{value}
	}

	return matches
}

func StrContains(haystack string, needles ...string) bool {
	for _, needle := range needles {
		if needle == "" || strings.Contains(haystack, needle) {
			return true
		}
	}

	return false
}

func StrContainsAll(haystack string, needles []string) bool {
	for _, needle := range needles {
		if !strings.Contains(haystack, needle) {
			return false
		}
	}

	return true
}

func StrDoesntContain(haystack string, needles ...string) bool {
	return !StrContains(haystack, needles...)
}

func StrStartsWith(subject string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(subject, prefix) {
			return true
		}
	}

	return false
}

func StrDoesntStartWith(subject string, prefixes ...string) bool {
	return !StrStartsWith(subject, prefixes...)
}

func StrEndsWith(subject string, suffixes ...string) bool {
	for _, suffix := range suffixes {
		if strings.HasSuffix(subject, suffix) {
			return true
		}
	}

	return false
}

func StrDoesntEndWith(subject string, suffixes ...string) bool {
	return !StrEndsWith(subject, suffixes...)
}

func StrIs(pattern, value string, ignoreCase ...bool) bool {
	ci := len(ignoreCase) > 0 && ignoreCase[0]

	if ci {
		pattern = strings.ToLower(pattern)
		value = strings.ToLower(value)
	}

	if pattern == value {
		return true
	}

	if pattern == "*" {
		return true
	}

	re := globToRegex(pattern)
	matched, err := regexp.MatchString(re, value)

	return err == nil && matched
}

func globToRegex(pattern string) string {
	var sb strings.Builder

	sb.WriteString("^")

	for _, r := range pattern {
		switch r {
		case '*':
			sb.WriteString(".*")
		case '?':
			sb.WriteString(".")
		case '.', '+', '(', ')', '[', ']', '{', '}', '^', '$', '|', '\\':
			sb.WriteRune('\\')
			sb.WriteRune(r)
		default:
			sb.WriteRune(r)
		}
	}

	sb.WriteString("$")

	return sb.String()
}

func StrIsMatch(patterns []string, value string) bool {
	for _, pattern := range patterns {
		matched, err := regexp.MatchString(pattern, value)

		if err == nil && matched {
			return true
		}
	}

	return false
}

func StrIsAscii(value string) bool {
	for _, r := range value {
		if r > 127 {
			return false
		}
	}

	return true
}

func StrIsJson(value string) bool {
	if value == "" {
		return false
	}

	var v any

	return json.Unmarshal([]byte(value), &v) == nil
}

func StrIsUrl(value string, protocols ...string) bool {
	u, err := url.ParseRequestURI(value)

	if err != nil {
		return false
	}

	if u.Scheme == "" || u.Host == "" {
		return false
	}

	if len(protocols) > 0 {
		for _, p := range protocols {
			if u.Scheme == p {
				return true
			}
		}

		return false
	}

	return true
}

func StrIsUuid(value string) bool {
	re := regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

	return re.MatchString(value)
}

func StrIsUlid(value string) bool {
	if len(value) != 26 {
		return false
	}

	re := regexp.MustCompile(`^[0-7][0-9A-HJKMNP-TV-Z]{25}$`)

	return re.MatchString(strings.ToUpper(value))
}

func StrLength(value string) int {
	return utf8.RuneCountInString(value)
}

func StrLimit(value string, limit int, end ...string) string {
	suffix := "..."

	if len(end) > 0 {
		suffix = end[0]
	}

	runes := []rune(value)

	if len(runes) <= limit {
		return value
	}

	return string(runes[:limit]) + suffix
}

func StrWords(value string, words int, end ...string) string {
	suffix := "..."

	if len(end) > 0 {
		suffix = end[0]
	}

	wordList := regexp.MustCompile(`\s+`).Split(strings.TrimSpace(value), -1)

	if len(wordList) <= words {
		return value
	}

	return strings.Join(wordList[:words], " ") + suffix
}

func StrLower(value string) string {
	return strings.ToLower(value)
}

func StrUpper(value string) string {
	return strings.ToUpper(value)
}

func StrReverse(value string) string {
	runes := []rune(value)

	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}

	return string(runes)
}

func StrReplace(search, replace any, subject string, caseSensitive ...bool) string {
	cs := len(caseSensitive) == 0 || caseSensitive[0]

	doReplace := func(s, r string) string {
		if cs {
			return strings.ReplaceAll(subject, s, r)
		}

		re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(s))

		return re.ReplaceAllString(subject, r)
	}

	switch s := search.(type) {
	case string:
		r, _ := replace.(string)
		subject = doReplace(s, r)
	case []string:
		switch r := replace.(type) {
		case []string:
			for i, sv := range s {
				rep := ""

				if i < len(r) {
					rep = r[i]
				}

				subject = doReplace(sv, rep)
			}
		case string:
			for _, sv := range s {
				subject = doReplace(sv, r)
			}
		}
	}

	return subject
}

func StrReplaceArray(search string, replacements []string, subject string) string {
	idx := 0

	for {
		pos := strings.Index(subject, search)

		if pos == -1 {
			break
		}

		rep := ""

		if idx < len(replacements) {
			rep = replacements[idx]
			idx++
		}

		subject = subject[:pos] + rep + subject[pos+len(search):]
	}

	return subject
}

func StrReplaceFirst(search, replace, subject string) string {
	if search == "" {
		return subject
	}

	idx := strings.Index(subject, search)

	if idx == -1 {
		return subject
	}

	return subject[:idx] + replace + subject[idx+len(search):]
}

func StrReplaceLast(search, replace, subject string) string {
	if search == "" {
		return subject
	}

	idx := strings.LastIndex(subject, search)

	if idx == -1 {
		return subject
	}

	return subject[:idx] + replace + subject[idx+len(search):]
}

func StrReplaceStart(search, replace, subject string) string {
	if search == "" || strings.HasPrefix(subject, search) {
		return replace + subject[len(search):]
	}

	return subject
}

func StrReplaceEnd(search, replace, subject string) string {
	if search != "" && strings.HasSuffix(subject, search) {
		return subject[:len(subject)-len(search)] + replace
	}

	return subject
}

func StrReplaceMatches(pattern, replace, subject string, limit ...int) string {
	re, err := regexp.Compile(pattern)

	if err != nil {
		return subject
	}

	n := -1

	if len(limit) > 0 {
		n = limit[0]
	}

	if n < 0 {
		return re.ReplaceAllString(subject, replace)
	}

	count := 0

	return re.ReplaceAllStringFunc(subject, func(match string) string {
		if count >= n {
			return match
		}

		count++

		return re.ReplaceAllString(match, replace)
	})
}

func StrRemove(search any, subject string, caseSensitive ...bool) string {
	return StrReplace(search, "", subject, caseSensitive...)
}

func StrRepeat(str string, times int) string {
	return strings.Repeat(str, times)
}

func StrSquish(value string) string {
	re := regexp.MustCompile(`\s+`)

	return strings.TrimSpace(re.ReplaceAllString(value, " "))
}

func StrDeduplicate(str string, characters ...string) string {
	chars := " "

	if len(characters) > 0 {
		chars = characters[0]
	}

	for _, ch := range chars {
		re := regexp.MustCompile(regexp.QuoteMeta(string(ch)) + `+`)
		str = re.ReplaceAllString(str, string(ch))
	}

	return str
}

func StrTrim(value string, chars ...string) string {
	if len(chars) > 0 {
		return strings.Trim(value, chars[0])
	}

	return strings.TrimSpace(value)
}

func StrLtrim(value string, chars ...string) string {
	if len(chars) > 0 {
		return strings.TrimLeft(value, chars[0])
	}

	return strings.TrimLeftFunc(value, unicode.IsSpace)
}

func StrRtrim(value string, chars ...string) string {
	if len(chars) > 0 {
		return strings.TrimRight(value, chars[0])
	}

	return strings.TrimRightFunc(value, unicode.IsSpace)
}

func StrStart(value, prefix string) string {
	if strings.HasPrefix(value, prefix) {
		return value
	}

	return prefix + value
}

func StrFinish(value, cap string) string {
	if strings.HasSuffix(value, cap) {
		return value
	}

	return value + cap
}

func StrWrap(value, before string, after ...string) string {
	a := before

	if len(after) > 0 {
		a = after[0]
	}

	return before + value + a
}

func StrUnwrap(value, before, after string) string {
	if after == "" {
		after = before
	}

	if strings.HasPrefix(value, before) && strings.HasSuffix(value, after) {
		value = value[len(before):]

		if strings.HasSuffix(value, after) {
			value = value[:len(value)-len(after)]
		}
	}

	return value
}

func StrChopStart(subject string, needles ...string) string {
	for _, needle := range needles {
		if strings.HasPrefix(subject, needle) {
			return subject[len(needle):]
		}
	}

	return subject
}

func StrChopEnd(subject string, needles ...string) string {
	for _, needle := range needles {
		if strings.HasSuffix(subject, needle) {
			return subject[:len(subject)-len(needle)]
		}
	}

	return subject
}

func StrMask(str, character string, index int, length ...int) string {
	runes := []rune(str)
	runeCount := len(runes)

	if index < 0 {
		index = max(runeCount+index, 0)
	}

	end := runeCount

	if len(length) > 0 && length[0] >= 0 {
		end = index + length[0]

		if end > runeCount {
			end = runeCount
		}
	}

	maskRune, _ := utf8.DecodeRuneInString(character)

	for i := index; i < end; i++ {
		runes[i] = maskRune
	}

	return string(runes)
}

func StrExcerpt(text, phrase string, radius int, omission ...string) string {
	om := "..."

	if len(omission) > 0 {
		om = omission[0]
	}

	if phrase == "" {
		runes := []rune(text)

		if len(runes) <= radius {
			return text
		}

		return string(runes[:radius]) + om
	}

	idx := strings.Index(strings.ToLower(text), strings.ToLower(phrase))

	if idx == -1 {
		runes := []rune(text)

		if len(runes) <= radius {
			return text
		}

		return string(runes[:radius]) + om
	}

	runes := []rune(text)
	phraseRunes := []rune(phrase)
	phraseStart := utf8.RuneCountInString(text[:idx])
	phraseEnd := phraseStart + len(phraseRunes)

	start := max(0, phraseStart-radius)
	end := min(len(runes), phraseEnd+radius)

	var sb strings.Builder

	if start > 0 {
		sb.WriteString(om)
	}

	sb.WriteString(string(runes[start:end]))

	if end < len(runes) {
		sb.WriteString(om)
	}

	return sb.String()
}

func StrMatch(pattern, subject string) string {
	re, err := regexp.Compile(pattern)

	if err != nil {
		return ""
	}

	match := re.FindStringSubmatch(subject)

	if len(match) == 0 {
		return ""
	}

	if len(match) > 1 {
		return match[1]
	}

	return match[0]
}

func StrMatchAll(pattern, subject string) []string {
	re, err := regexp.Compile(pattern)

	if err != nil {
		return nil
	}

	matches := re.FindAllStringSubmatch(subject, -1)

	var result []string

	for _, m := range matches {
		if len(m) > 1 {
			result = append(result, m[1])
		} else if len(m) > 0 {
			result = append(result, m[0])
		}
	}

	return result
}

func StrNumbers(value string) string {
	var sb strings.Builder

	for _, r := range value {
		if unicode.IsDigit(r) {
			sb.WriteRune(r)
		}
	}

	return sb.String()
}

func StrPadBoth(value string, length int, pad ...string) string {
	padStr := " "

	if len(pad) > 0 {
		padStr = pad[0]
	}

	runes := []rune(value)
	total := length - len(runes)

	if total <= 0 {
		return value
	}

	leftPad := total / 2
	rightPad := total - leftPad

	return strPad(padStr, leftPad) + value + strPad(padStr, rightPad)
}

func StrPadLeft(value string, length int, pad ...string) string {
	padStr := " "

	if len(pad) > 0 {
		padStr = pad[0]
	}

	runes := []rune(value)
	total := length - len(runes)

	if total <= 0 {
		return value
	}

	return strPad(padStr, total) + value
}

func StrPadRight(value string, length int, pad ...string) string {
	padStr := " "

	if len(pad) > 0 {
		padStr = pad[0]
	}

	runes := []rune(value)
	total := length - len(runes)

	if total <= 0 {
		return value
	}

	return value + strPad(padStr, total)
}

func strPad(pad string, n int) string {
	if pad == "" || n <= 0 {
		return ""
	}

	repeated := strings.Repeat(pad, (n/len([]rune(pad)))+1)

	return string([]rune(repeated)[:n])
}

func StrPosition(haystack, needle string, offset ...int) (int, bool) {
	off := 0

	if len(offset) > 0 {
		off = offset[0]
	}

	runes := []rune(haystack)

	if off < 0 {
		off = len(runes) + off
	}

	if off < 0 || off >= len(runes) {
		return -1, false
	}

	sub := string(runes[off:])
	idx := strings.Index(sub, needle)

	if idx == -1 {
		return -1, false
	}

	return off + utf8.RuneCountInString(sub[:idx]), true
}

func StrSubstr(str string, start int, length ...int) string {
	runes := []rune(str)
	runeLen := len(runes)

	if start < 0 {
		start = max(runeLen+start, 0)
	}

	if start >= runeLen {
		return ""
	}

	end := runeLen

	if len(length) > 0 {
		l := length[0]

		if l < 0 {
			end = runeLen + l
		} else {
			end = start + l
		}
	}

	if end > runeLen {
		end = runeLen
	}

	if start >= end {
		return ""
	}

	return string(runes[start:end])
}

func StrSubstrCount(haystack, needle string, offset ...int) int {
	if needle == "" {
		return 0
	}

	off := 0

	if len(offset) > 0 {
		off = offset[0]
	}

	if off > 0 {
		runes := []rune(haystack)

		if off >= len(runes) {
			return 0
		}

		haystack = string(runes[off:])
	}

	return strings.Count(haystack, needle)
}

func StrSubstrReplace(str, replace string, offset int, length ...int) string {
	runes := []rune(str)
	runeLen := len(runes)

	if offset < 0 {
		offset = max(runeLen+offset, 0)
	}

	if offset > runeLen {
		offset = runeLen
	}

	end := runeLen

	if len(length) > 0 {
		l := length[0]

		if l < 0 {
			end = max(runeLen+l, offset)
		} else {
			end = min(offset+l, runeLen)
		}
	}

	result := make([]rune, 0, runeLen)
	result = append(result, runes[:offset]...)
	result = append(result, []rune(replace)...)
	result = append(result, runes[end:]...)

	return string(result)
}

func StrSwap(m map[string]string, subject string) string {
	if len(m) == 0 {
		return subject
	}

	for k, v := range m {
		subject = strings.ReplaceAll(subject, k, v)
	}

	return subject
}

func StrTake(str string, n int) string {
	runes := []rune(str)

	if n >= 0 {
		if n >= len(runes) {
			return str
		}

		return string(runes[:n])
	}

	start := len(runes) + n

	if start < 0 {
		start = 0
	}

	return string(runes[start:])
}

func StrWordCount(str string) int {
	return len(strings.Fields(str))
}

func StrWordWrap(str string, width int, breakStr ...string) string {
	brk := "\n"

	if len(breakStr) > 0 {
		brk = breakStr[0]
	}

	words := strings.Fields(str)

	if len(words) == 0 {
		return str
	}

	var lines []string
	currentLine := words[0]

	for _, word := range words[1:] {
		if len([]rune(currentLine))+1+len([]rune(word)) <= width {
			currentLine += " " + word
		} else {
			lines = append(lines, currentLine)
			currentLine = word
		}
	}

	lines = append(lines, currentLine)

	return strings.Join(lines, brk)
}

func StrToBase64(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

func StrFromBase64(value string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(value)

	if err != nil {

		decoded, err = base64.URLEncoding.DecodeString(value)

		if err != nil {
			return "", err
		}
	}

	return string(decoded), nil
}

func StrSlug(title string, separator ...string) string {
	sep := "-"

	if len(separator) > 0 {
		sep = separator[0]
	}

	title = StrAscii(title)
	title = strings.ToLower(title)

	title = strings.ReplaceAll(title, "@", sep+"at"+sep)

	re := regexp.MustCompile(`[^a-z0-9]+`)
	title = re.ReplaceAllString(title, sep)

	title = strings.Trim(title, sep)

	if sep != "" {
		re2 := regexp.MustCompile(regexp.QuoteMeta(sep) + `+`)
		title = re2.ReplaceAllString(title, sep)
	}

	return title
}

func StrAscii(value string, language ...string) string {
	var sb strings.Builder

	for _, r := range value {
		if r <= 127 {
			sb.WriteRune(r)

			continue
		}

		if s, ok := asciiMap[r]; ok {
			sb.WriteString(s)
		}

	}

	return sb.String()
}

func StrTransliterate(str string, unknown ...string) string {
	unk := "?"

	if len(unknown) > 0 {
		unk = unknown[0]
	}

	var sb strings.Builder

	for _, r := range str {
		if r <= 127 {
			sb.WriteRune(r)

			continue
		}

		if s, ok := asciiMap[r]; ok {
			sb.WriteString(s)
		} else {
			sb.WriteString(unk)
		}
	}

	return sb.String()
}

func StrInitials(name string, delimiter ...string) string {
	sep := ""

	if len(delimiter) > 0 {
		sep = delimiter[0]
	}

	words := strings.Fields(name)

	var initials []string

	for _, w := range words {
		r, _ := utf8.DecodeRuneInString(w)
		initials = append(initials, string(unicode.ToUpper(r)))
	}

	return strings.Join(initials, sep)
}

func StrCharAt(subject string, index int) string {
	runes := []rune(subject)

	if index < 0 {
		index = len(runes) + index
	}

	if index < 0 || index >= len(runes) {
		return ""
	}

	return string(runes[index])
}

func StrParseCallback(callback string, def ...string) [2]string {
	d := ""

	if len(def) > 0 {
		d = def[0]
	}

	for _, sep := range []string{"@", "::"} {
		if idx := strings.LastIndex(callback, sep); idx != -1 {
			return [2]string{callback[:idx], callback[idx+len(sep):]}
		}
	}

	return [2]string{callback, d}
}

func StrRandom(length ...int) string {
	l := 16

	if len(length) > 0 {
		l = length[0]
	}

	randomMu.Lock()
	factory := randomFactory
	sequence := randomSequence
	fallback := randomFallback
	randomMu.Unlock()

	if factory != nil {
		return factory(l)
	}

	if len(sequence) > 0 {
		randomMu.Lock()
		val := randomSequence[0]
		randomSequence = randomSequence[1:]
		randomMu.Unlock()

		return val
	}

	_ = fallback

	return generateRandom(l)
}

func generateRandom(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := range result {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		result[i] = charset[n.Int64()]
	}

	return string(result)
}

func CreateRandomStringsUsing(factory func(int) string) {
	randomMu.Lock()
	randomFactory = factory
	randomMu.Unlock()
}

func CreateRandomStringsUsingSequence(sequence []string, whenMissing ...func(int) string) {
	randomMu.Lock()
	randomSequence = sequence

	if len(whenMissing) > 0 {
		randomFallback = whenMissing[0]
	}

	randomMu.Unlock()
}

func CreateRandomStringsNormally() {
	randomMu.Lock()
	randomFactory = nil
	randomSequence = nil
	randomFallback = nil
	randomMu.Unlock()
}

func StrPassword(length ...int) (string, error) {
	l := 32

	if len(length) > 0 {
		l = length[0]
	}

	const (
		letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		numbers = "0123456789"
		symbols = "!@#$%^&*()_+-=[]{}|;':\",./<>?"
	)
	charset := letters + numbers + symbols

	result := make([]byte, l)

	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))

		if err != nil {
			return "", err
		}

		result[i] = charset[n.Int64()]
	}

	return string(result), nil
}

func FlushCache() {
	camelCache = sync.Map{}
	studlyCache = sync.Map{}
	snakeCache = sync.Map{}
}

func StrConvertCase(str string, mode int) string {
	switch mode {
	case 0:
		return strings.ToUpper(str)
	case 1:
		return strings.ToLower(str)
	case 2:
		return strings.Title(str)
	default:
		return strings.ToLower(str)
	}
}

func Of(value string) *StringBuilder {
	return &StringBuilder{value: value}
}

func (s *StringBuilder) String() string { return s.value }
func (s *StringBuilder) Value() string  { return s.value }

func (s *StringBuilder) After(search string) *StringBuilder {
	return &StringBuilder{value: StrAfter(s.value, search)}
}
func (s *StringBuilder) AfterLast(search string) *StringBuilder {
	return &StringBuilder{value: StrAfterLast(s.value, search)}
}
func (s *StringBuilder) Before(search string) *StringBuilder {
	return &StringBuilder{value: StrBefore(s.value, search)}
}
func (s *StringBuilder) BeforeLast(search string) *StringBuilder {
	return &StringBuilder{value: StrBeforeLast(s.value, search)}
}
func (s *StringBuilder) Between(from, to string) *StringBuilder {
	return &StringBuilder{value: StrBetween(s.value, from, to)}
}
func (s *StringBuilder) BetweenFirst(from, to string) *StringBuilder {
	return &StringBuilder{value: StrBetweenFirst(s.value, from, to)}
}
func (s *StringBuilder) Camel() *StringBuilder {
	return &StringBuilder{value: StrCamel(s.value)}
}
func (s *StringBuilder) Studly() *StringBuilder {
	return &StringBuilder{value: StrStudly(s.value)}
}
func (s *StringBuilder) Pascal() *StringBuilder {
	return &StringBuilder{value: StrPascal(s.value)}
}
func (s *StringBuilder) Snake(delimiter ...string) *StringBuilder {
	return &StringBuilder{value: StrSnake(s.value, delimiter...)}
}
func (s *StringBuilder) Kebab() *StringBuilder {
	return &StringBuilder{value: StrKebab(s.value)}
}
func (s *StringBuilder) Lower() *StringBuilder {
	return &StringBuilder{value: StrLower(s.value)}
}
func (s *StringBuilder) Upper() *StringBuilder {
	return &StringBuilder{value: StrUpper(s.value)}
}
func (s *StringBuilder) Title() *StringBuilder {
	return &StringBuilder{value: StrTitle(s.value)}
}
func (s *StringBuilder) Headline() *StringBuilder {
	return &StringBuilder{value: StrHeadline(s.value)}
}
func (s *StringBuilder) Apa() *StringBuilder {
	return &StringBuilder{value: StrApa(s.value)}
}
func (s *StringBuilder) Ucfirst() *StringBuilder {
	return &StringBuilder{value: StrUcfirst(s.value)}
}
func (s *StringBuilder) Lcfirst() *StringBuilder {
	return &StringBuilder{value: StrLcfirst(s.value)}
}
func (s *StringBuilder) Slug(separator ...string) *StringBuilder {
	return &StringBuilder{value: StrSlug(s.value, separator...)}
}
func (s *StringBuilder) Contains(needles ...string) bool {
	return StrContains(s.value, needles...)
}
func (s *StringBuilder) ContainsAll(needles []string) bool {
	return StrContainsAll(s.value, needles)
}
func (s *StringBuilder) StartsWith(prefixes ...string) bool {
	return StrStartsWith(s.value, prefixes...)
}
func (s *StringBuilder) EndsWith(suffixes ...string) bool {
	return StrEndsWith(s.value, suffixes...)
}
func (s *StringBuilder) Is(pattern string, ignoreCase ...bool) bool {
	return StrIs(pattern, s.value, ignoreCase...)
}
func (s *StringBuilder) IsAscii() bool    { return StrIsAscii(s.value) }
func (s *StringBuilder) IsJson() bool     { return StrIsJson(s.value) }
func (s *StringBuilder) IsEmpty() bool    { return s.value == "" }
func (s *StringBuilder) IsNotEmpty() bool { return s.value != "" }
func (s *StringBuilder) Length() int      { return StrLength(s.value) }
func (s *StringBuilder) Limit(limit int, end ...string) *StringBuilder {
	return &StringBuilder{value: StrLimit(s.value, limit, end...)}
}
func (s *StringBuilder) Words(words int, end ...string) *StringBuilder {
	return &StringBuilder{value: StrWords(s.value, words, end...)}
}
func (s *StringBuilder) Mask(character string, index int, length ...int) *StringBuilder {
	return &StringBuilder{value: StrMask(s.value, character, index, length...)}
}
func (s *StringBuilder) Match(pattern string) *StringBuilder {
	return &StringBuilder{value: StrMatch(pattern, s.value)}
}
func (s *StringBuilder) MatchAll(pattern string) []string {
	return StrMatchAll(pattern, s.value)
}
func (s *StringBuilder) Replace(search, replace any, caseSensitive ...bool) *StringBuilder {
	return &StringBuilder{value: StrReplace(search, replace, s.value, caseSensitive...)}
}
func (s *StringBuilder) ReplaceFirst(search, replace string) *StringBuilder {
	return &StringBuilder{value: StrReplaceFirst(search, replace, s.value)}
}
func (s *StringBuilder) ReplaceLast(search, replace string) *StringBuilder {
	return &StringBuilder{value: StrReplaceLast(search, replace, s.value)}
}
func (s *StringBuilder) ReplaceMatches(pattern, replace string) *StringBuilder {
	return &StringBuilder{value: StrReplaceMatches(pattern, replace, s.value)}
}
func (s *StringBuilder) Remove(search any, caseSensitive ...bool) *StringBuilder {
	return &StringBuilder{value: StrRemove(search, s.value, caseSensitive...)}
}
func (s *StringBuilder) Reverse() *StringBuilder {
	return &StringBuilder{value: StrReverse(s.value)}
}
func (s *StringBuilder) Squish() *StringBuilder {
	return &StringBuilder{value: StrSquish(s.value)}
}
func (s *StringBuilder) Trim(chars ...string) *StringBuilder {
	return &StringBuilder{value: StrTrim(s.value, chars...)}
}
func (s *StringBuilder) Ltrim(chars ...string) *StringBuilder {
	return &StringBuilder{value: StrLtrim(s.value, chars...)}
}
func (s *StringBuilder) Rtrim(chars ...string) *StringBuilder {
	return &StringBuilder{value: StrRtrim(s.value, chars...)}
}
func (s *StringBuilder) Start(prefix string) *StringBuilder {
	return &StringBuilder{value: StrStart(s.value, prefix)}
}
func (s *StringBuilder) Finish(cap string) *StringBuilder {
	return &StringBuilder{value: StrFinish(s.value, cap)}
}
func (s *StringBuilder) Wrap(before string, after ...string) *StringBuilder {
	return &StringBuilder{value: StrWrap(s.value, before, after...)}
}
func (s *StringBuilder) Unwrap(before, after string) *StringBuilder {
	return &StringBuilder{value: StrUnwrap(s.value, before, after)}
}
func (s *StringBuilder) ChopStart(needles ...string) *StringBuilder {
	return &StringBuilder{value: StrChopStart(s.value, needles...)}
}
func (s *StringBuilder) ChopEnd(needles ...string) *StringBuilder {
	return &StringBuilder{value: StrChopEnd(s.value, needles...)}
}
func (s *StringBuilder) Substr(start int, length ...int) *StringBuilder {
	return &StringBuilder{value: StrSubstr(s.value, start, length...)}
}
func (s *StringBuilder) Take(n int) *StringBuilder {
	return &StringBuilder{value: StrTake(s.value, n)}
}
func (s *StringBuilder) PadLeft(length int, pad ...string) *StringBuilder {
	return &StringBuilder{value: StrPadLeft(s.value, length, pad...)}
}
func (s *StringBuilder) PadRight(length int, pad ...string) *StringBuilder {
	return &StringBuilder{value: StrPadRight(s.value, length, pad...)}
}
func (s *StringBuilder) PadBoth(length int, pad ...string) *StringBuilder {
	return &StringBuilder{value: StrPadBoth(s.value, length, pad...)}
}
func (s *StringBuilder) Swap(m map[string]string) *StringBuilder {
	return &StringBuilder{value: StrSwap(m, s.value)}
}
func (s *StringBuilder) Append(values ...string) *StringBuilder {
	return &StringBuilder{value: s.value + strings.Join(values, "")}
}
func (s *StringBuilder) Prepend(values ...string) *StringBuilder {
	return &StringBuilder{value: strings.Join(values, "") + s.value}
}
func (s *StringBuilder) Numbers() *StringBuilder {
	return &StringBuilder{value: StrNumbers(s.value)}
}
func (s *StringBuilder) Ascii(language ...string) *StringBuilder {
	return &StringBuilder{value: StrAscii(s.value, language...)}
}
func (s *StringBuilder) ToBase64() *StringBuilder {
	return &StringBuilder{value: StrToBase64(s.value)}
}
func (s *StringBuilder) FromBase64() (*StringBuilder, error) {
	v, err := StrFromBase64(s.value)

	return &StringBuilder{value: v}, err
}
func (s *StringBuilder) Plural(count ...int) *StringBuilder {
	return &StringBuilder{value: StrPlural(s.value, count...)}
}
func (s *StringBuilder) Singular() *StringBuilder {
	return &StringBuilder{value: StrSingular(s.value)}
}
func (s *StringBuilder) Initials(delimiter ...string) *StringBuilder {
	return &StringBuilder{value: StrInitials(s.value, delimiter...)}
}
func (s *StringBuilder) WordCount() int { return StrWordCount(s.value) }
func (s *StringBuilder) WordWrap(width int, breakStr ...string) *StringBuilder {
	return &StringBuilder{value: StrWordWrap(s.value, width, breakStr...)}
}
func (s *StringBuilder) Excerpt(phrase string, radius int, omission ...string) *StringBuilder {
	return &StringBuilder{value: StrExcerpt(s.value, phrase, radius, omission...)}
}
func (s *StringBuilder) Ucsplit() []string { return StrUcsplit(s.value) }
func (s *StringBuilder) Deduplicate(chars ...string) *StringBuilder {
	return &StringBuilder{value: StrDeduplicate(s.value, chars...)}
}
func (s *StringBuilder) Markdown(options ...map[string]any) *StringBuilder {
	return &StringBuilder{value: StrMarkdown(s.value, options...)}
}
func (s *StringBuilder) InlineMarkdown(options ...map[string]any) *StringBuilder {
	return &StringBuilder{value: StrInlineMarkdown(s.value, options...)}
}

// min/max helpers (go 1.21+ has these as builtins, but keeping for clarity)
func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
