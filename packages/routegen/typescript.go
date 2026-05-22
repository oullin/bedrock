package routegen

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// reservedKeywords is the complete list of TypeScript/JavaScript reserved
// words that cannot be used as identifiers without quoting.
var reservedKeywords = map[string]struct{}{
	"await": {}, "break": {}, "case": {}, "catch": {}, "class": {},
	"const": {}, "continue": {}, "debugger": {}, "default": {}, "delete": {},
	"do": {}, "else": {}, "enum": {}, "export": {}, "extends": {},
	"false": {}, "finally": {}, "for": {}, "function": {}, "if": {},
	"implements": {}, "import": {}, "in": {}, "instanceof": {}, "interface": {},
	"let": {}, "new": {}, "null": {}, "package": {}, "private": {},
	"protected": {}, "public": {}, "return": {}, "static": {}, "super": {},
	"switch": {}, "this": {}, "throw": {}, "true": {}, "try": {},
	"typeof": {}, "var": {}, "void": {}, "while": {}, "with": {}, "yield": {},
}

// nonIdentRe matches characters that are not valid in a TypeScript identifier
// (keeping unicode letters, digits, underscore, and $).
var nonIdentRe = regexp.MustCompile(`[^\p{L}\p{Nd}_$-]`)

// SafeMethod sanitizes a raw method name so it is safe to use as a TypeScript
// identifier.
//
// suffix is "Method" or "Param". It is appended (lowercased + ucfirst) when
// the name is a reserved keyword, or prepended (lowercased) when the name
// starts with an invalid first character.
func SafeMethod(method, suffix string) string {
	// 1. Replace characters that are not unicode letters, digits, _ or $ with _.
	s := nonIdentRe.ReplaceAllString(method, "_")

	// 2. If the result still contains a hyphen convert to camelCase.
	if strings.Contains(s, "-") {
		s = toCamel(s)
	}

	sfxLower := strings.ToLower(suffix)
	sfxUpper := strings.ToUpper(sfxLower[:1]) + sfxLower[1:]

	// 3. Reserved keyword → append lowered suffix with ucfirst.
	if _, ok := reservedKeywords[s]; ok {
		return s + sfxUpper
	}

	// 4. Starts with a non-[a-zA-Z_$] character → prepend lowered suffix.
	first, _ := utf8.DecodeRuneInString(s)

	if first != utf8.RuneError && !unicode.IsLetter(first) && first != '_' && first != '$' {
		return sfxLower + s
	}

	return s
}

// QuoteIfNeeded returns the key quoted when it starts with a digit but is not
// a pure integer (e.g., "2fa" → `"2fa"`). Pure numbers are returned as-is.
func QuoteIfNeeded(name string) string {
	if isNumeric(name) {
		return name
	}

	first := name[0]

	if first >= '0' && first <= '9' {
		return `"` + name + `"`
	}

	return name
}

// isNumeric reports whether s represents a non-negative integer.
func isNumeric(s string) bool {
	if s == "" {
		return false
	}

	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}

	return true
}

// toCamel converts a kebab-case string to camelCase.
// "invalid-js-name" → "invalidJsName".
func toCamel(s string) string {
	parts := strings.Split(s, "-")

	for i := 1; i < len(parts); i++ {
		if len(parts[i]) > 0 {
			parts[i] = strings.ToUpper(parts[i][:1]) + parts[i][1:]
		}
	}

	return strings.Join(parts, "")
}

// cleanUpRe holds pre-compiled regexps used by CleanUp.
var (
	cleanUpArrowRe   = regexp.MustCompile(`=>\s*\{\n{2,}`)
	cleanUpReplaceRe = regexp.MustCompile(`\\\s+\.replace`)
	cleanUpQueryRe   = regexp.MustCompile(`\s+\+ queryParams\(options\)`)
	cleanUpMultiNlRe = regexp.MustCompile(`\n{3,}`)
)

// CleanUp normalises the whitespace and indentation of a generated TypeScript
// file.
func CleanUp(src string) string {
	// --- simple string replacements (order matters) ---
	replacements := [][2]string{
		{" ,", ","},
		{"[ ", "["},
		{", }", " }"},
		{"} )", "})"},
		{" )", ")"},
		{"( ", "("},
	}

	for _, r := range replacements {
		src = strings.ReplaceAll(src, r[0], r[1])
	}

	// --- re-indent based on brace depth ---
	src = reindent(src)

	// --- regex replacements ---
	src = cleanUpArrowRe.ReplaceAllString(src, "=> {\n")
	src = cleanUpReplaceRe.ReplaceAllStringFunc(src, func(m string) string {
		return "\n            .replace"
	})
	src = cleanUpQueryRe.ReplaceAllString(src, " + queryParams(options)")
	src = cleanUpMultiNlRe.ReplaceAllString(src, "\n\n")

	return strings.TrimSpace(src) + "\n"
}

// reindent strips leading whitespace from every line then re-applies
// indentation based on open/close brace/bracket depth (4 spaces per level).
func reindent(src string) string {
	lines := strings.Split(src, "\n")
	depth := 0

	var b strings.Builder

	b.Grow(len(src))

	for i, raw := range lines {
		line := strings.TrimSpace(raw)

		if line == "" {
			b.WriteByte('\n')

			continue
		}

		// Decrease depth before writing closing braces/brackets.
		if strings.HasPrefix(line, "}") || strings.HasPrefix(line, "]") {
			if depth > 0 {
				depth--
			}
		}

		b.WriteString(strings.Repeat("    ", depth))
		b.WriteString(line)

		if i < len(lines)-1 {
			b.WriteByte('\n')
		}

		// Increase depth after writing lines that open a block.
		if strings.HasSuffix(line, "{") || strings.HasSuffix(line, "[") {
			depth++
		}
	}

	return b.String()
}
