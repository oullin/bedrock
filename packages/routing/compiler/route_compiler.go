package compiler

import (
	"fmt"
	"regexp"
	"strings"
)

// SourceRoute is the minimum surface a compilable route must expose.
//
// The concrete [routing.Route] type implements this interface. Tests can
// substitute their own implementation without pulling in the full Route type.
type SourceRoute interface {
	// Path returns the route URI pattern (e.g. "/users/{user}").
	Path() string
	// Host returns the host pattern, or the empty string if none.
	Host() string
	// Requirements returns parameter regex requirements keyed by name.
	Requirements() map[string]string
	// HasDefault reports whether the named parameter has a default value.
	HasDefault(name string) bool
}

// Separators are the characters that may serve as a separator immediately
// before an optional variable.

// VariableMaximumLength is the longest a single PCRE/RE2 named group may be.

// variableRe matches a single placeholder of the form {name} or {!name}.
//
// The leading "!" denotes a Symfony "important variable" that must never be
// optional. The Go RE2 character class \w is equivalent to PCRE \w for ASCII
// names; upstream does not currently emit non-ASCII parameter names.

// Compile produces a [CompiledRoute] from a [SourceRoute].
//
// It compiles the host pattern (when present), then the path pattern, and
// merges the two variable lists.

type compileResult struct {
	staticPrefix string
	regex        string
	tokens       []Token
	variables    []string
}

const Separators = "/,;.:-_~+*=@|"

const VariableMaximumLength = 32

var variableRe = regexp.MustCompile(`\{(!)?([\w\x80-\xff]+)\}`)

// validGroupName matches variable names that are safe for RE2 named capture
// groups (?P<name>...). RE2 requires [A-Za-z_][A-Za-z0-9_]*.
var validGroupName = regexp.MustCompile(`^[A-Za-z_]\w*$`)

func Compile(route SourceRoute) (*CompiledRoute, error) {
	var (
		hostVariables []string
		variables     []string
		hostRegex     string
		hostTokens    []Token
	)

	if host := route.Host(); host != "" {
		res, err := compilePattern(route, host, true)

		if err != nil {
			return nil, err
		}

		hostVariables = res.variables
		variables = append(variables, hostVariables...)
		hostTokens = res.tokens
		hostRegex = res.regex
	}

	path := route.Path()
	res, err := compilePattern(route, path, false)

	if err != nil {
		return nil, err
	}

	for _, v := range res.variables {
		if v == "_fragment" {
			return nil, fmt.Errorf(`route pattern %q cannot contain "_fragment" as a path parameter`, path)
		}
	}

	pathVariables := res.variables
	variables = append(variables, pathVariables...)

	return NewCompiledRoute(
		res.staticPrefix,
		res.regex,
		res.tokens,
		pathVariables,
		hostRegex,
		hostTokens,
		hostVariables,
		uniqueStrings(variables),
	), nil
}

func compilePattern(route SourceRoute, pattern string, isHost bool) (*compileResult, error) {
	var (
		tokens    []Token
		variables []string
		pos       int
	)
	defaultSeparator := "/"

	if isHost {
		defaultSeparator = "."
	}

	matches := variableRe.FindAllStringSubmatchIndex(pattern, -1)

	for _, m := range matches {
		// m = [matchStart, matchEnd, importantStart, importantEnd, nameStart, nameEnd]
		matchStart, matchEnd := m[0], m[1]
		important := m[2] >= 0
		varName := pattern[m[4]:m[5]]

		precedingText := pattern[pos:matchStart]
		pos = matchEnd

		var precedingChar string

		if precedingText != "" {
			// last byte; upstream route names are ASCII so byte slicing is safe.
			precedingChar = precedingText[len(precedingText)-1:]
		}

		isSeparator := precedingChar != "" && strings.Contains(Separators, precedingChar)

		if len(varName) > 0 && varName[0] >= '0' && varName[0] <= '9' {
			return nil, fmt.Errorf("variable name %q cannot start with a digit in route pattern %q", varName, pattern)
		}

		if !validGroupName.MatchString(varName) {
			return nil, fmt.Errorf("variable name %q contains characters not supported by RE2 named groups in route pattern %q", varName, pattern)
		}

		for _, existing := range variables {
			if existing == varName {
				return nil, fmt.Errorf("route pattern %q cannot reference variable name %q more than once", pattern, varName)
			}
		}

		if len(varName) > VariableMaximumLength {
			return nil, fmt.Errorf("variable name %q cannot be longer than %d characters in route pattern %q", varName, VariableMaximumLength, pattern)
		}

		if isSeparator && precedingText != precedingChar {
			tokens = append(tokens, Token{Kind: "text", Prefix: precedingText[:len(precedingText)-len(precedingChar)]})
		} else if !isSeparator && precedingText != "" {
			tokens = append(tokens, Token{Kind: "text", Prefix: precedingText})
		}

		regexp := route.Requirements()[varName]

		if regexp == "" {
			followingPattern := pattern[pos:]
			nextSeparator := findNextSeparator(followingPattern)
			extra := ""

			if defaultSeparator != nextSeparator && nextSeparator != "" {
				extra = quoteMeta(nextSeparator)
			}

			regexp = fmt.Sprintf("[^%s%s]+", quoteMeta(defaultSeparator), extra)
		} else {
			regexp = transformCapturingGroupsToNonCapturing(regexp)
		}

		tokenPrefix := ""

		if isSeparator {
			tokenPrefix = precedingChar
		}

		tokens = append(tokens, Token{
			Kind:      "variable",
			Prefix:    tokenPrefix,
			Regexp:    regexp,
			Name:      varName,
			Important: important,
		})
		variables = append(variables, varName)
	}

	if pos < len(pattern) {
		tokens = append(tokens, Token{Kind: "text", Prefix: pattern[pos:]})
	}

	// Find first optional token (a variable with a default and no important
	// flag, with only optional variables following it).
	firstOptional := -1

	if !isHost {
		for i := len(tokens) - 1; i >= 0; i-- {
			t := tokens[i]

			if t.Kind == "variable" && !t.Important && route.HasDefault(t.Name) {
				firstOptional = i
			} else {
				break
			}
		}
	}

	// Compute the matching regexp.
	var regexpBuilder strings.Builder

	for i := range tokens {
		regexpBuilder.WriteString(computeRegexp(tokens, i, firstOptional))
	}

	prefix := "(?s)"

	if isHost {
		prefix = "(?si)"
	}

	regexpStr := prefix + "^" + regexpBuilder.String() + "$"

	// Reverse tokens (Symfony stores them reversed for the dumper).
	reversed := make([]Token, len(tokens))

	for i, t := range tokens {
		reversed[len(tokens)-1-i] = t
	}

	return &compileResult{
		staticPrefix: determineStaticPrefix(route, tokens),
		regex:        regexpStr,
		tokens:       reversed,
		variables:    variables,
	}, nil
}

// determineStaticPrefix returns the longest constant-text prefix that the
// regex can begin with.
func determineStaticPrefix(route SourceRoute, tokens []Token) string {
	if len(tokens) == 0 {
		return ""
	}

	if tokens[0].Kind != "text" {
		if route.HasDefault(tokens[0].Name) || tokens[0].Prefix == "/" {
			return ""
		}

		return tokens[0].Prefix
	}

	prefix := tokens[0].Prefix

	if len(tokens) > 1 && tokens[1].Prefix != "/" && !route.HasDefault(tokens[1].Name) {
		prefix += tokens[1].Prefix
	}

	return prefix
}

// findNextSeparator returns the next static separator character (after
// stripping inner placeholders) or the empty string.
func findNextSeparator(pattern string) string {
	if pattern == "" {
		return ""
	}

	stripped := variableRe.ReplaceAllString(pattern, "")

	if stripped == "" {
		return ""
	}

	first := stripped[0:1]

	if strings.Contains(Separators, first) {
		return first
	}

	return ""
}

// computeRegexp builds the regex fragment for a single token.
//
// For the leading optional variable, the surrounding separator is included
// inside the optional group so the entire segment can be elided.
func computeRegexp(tokens []Token, index, firstOptional int) string {
	t := tokens[index]

	if t.Kind == "text" {
		return quoteMeta(t.Prefix)
	}
	// Variable token.
	if index == 0 && firstOptional == 0 {
		return fmt.Sprintf("%s(?P<%s>%s)?", quoteMeta(t.Prefix), t.Name, t.Regexp)
	}

	regex := fmt.Sprintf("%s(?P<%s>%s)", quoteMeta(t.Prefix), t.Name, t.Regexp)

	if firstOptional >= 0 && index >= firstOptional {
		// Wrap this and following optional tokens in a single optional group
		// at the firstOptional position. The actual nesting is built inside
		// out by the caller via successive computeRegexp calls; here we just
		// emit "(?:..." once at firstOptional and ")?" at the end.
		nbOptional := index - firstOptional

		if index == firstOptional {
			regex = "(?:" + regex
		}

		if index == len(tokens)-1 {
			regex += strings.Repeat(")?", nbOptional+1)
		}
	}

	return regex
}

// transformCapturingGroupsToNonCapturing converts user-supplied PCRE capturing
// groups to non-capturing so the named groups produced by the compiler remain
// the only addressable captures. RE2 does not support backreferences so this
// is a strict syntactic rewrite.
func transformCapturingGroupsToNonCapturing(pattern string) string {
	var b strings.Builder

	for i := 0; i < len(pattern); i++ {
		if pattern[i] == '\\' && i+1 < len(pattern) {
			b.WriteByte(pattern[i])
			b.WriteByte(pattern[i+1])
			i++

			continue
		}

		if pattern[i] == '(' && i+1 < len(pattern) && pattern[i+1] != '?' {
			b.WriteString("(?:")

			continue
		}

		b.WriteByte(pattern[i])
	}

	return b.String()
}

// It escapes regex metacharacters
// in literal text portions of the pattern.
func quoteMeta(s string) string { return regexp.QuoteMeta(s) }

func uniqueStrings(in []string) []string {
	seen := make(map[string]struct{}, len(in))
	out := make([]string, 0, len(in))

	for _, s := range in {
		if _, ok := seen[s]; ok {
			continue
		}

		seen[s] = struct{}{}
		out = append(out, s)
	}

	return out
}
