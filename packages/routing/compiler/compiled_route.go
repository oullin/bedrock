package compiler

import "regexp"

// Token is one element of a compiled route pattern.
//
// Kind is either "text" (a literal preceding a variable or trailing the
// pattern) or "variable" (a placeholder).
//
// For text tokens, only Prefix is meaningful and holds the literal text.
//
// For variable tokens:
//   - Prefix:    the separator character that immediately precedes the variable
//     in the URI (typically "/", "." or empty).
//   - Regexp:    the inner regular expression that matches the variable's value.
//   - Name:      the variable name as it appears between the braces.
//   - Utf8:      true when the variable was authored with UTF-8 semantics.
//   - Important: true when the variable was prefixed with "!" in the source
//     pattern (Symfony "important variable" — never optional).
//
// ([type, prefix, regexp, name, utf8, important]) so a port of any test that
// inspects compiled tokens remains source-compatible.
type Token struct {
	Kind      string
	Prefix    string
	Regexp    string
	Name      string
	Utf8      bool
	Important bool
}

// CompiledRoute is the result of compiling a route URI (and optionally a host
// pattern).
type CompiledRoute struct {
	staticPrefix  string
	regex         string
	tokens        []Token
	pathVariables []string
	hostRegex     string
	hostTokens    []Token
	hostVariables []string
	variables     []string

	compiledRegex     *regexp.Regexp
	compiledHostRegex *regexp.Regexp
}

// NewCompiledRoute constructs a CompiledRoute. Used by [Compile]; rarely needed
// by user code.
func NewCompiledRoute(
	staticPrefix, regex string,
	tokens []Token,
	pathVariables []string,
	hostRegex string,
	hostTokens []Token,
	hostVariables []string,
	variables []string,
) *CompiledRoute {
	cr := &CompiledRoute{
		staticPrefix:  staticPrefix,
		regex:         regex,
		tokens:        tokens,
		pathVariables: pathVariables,
		hostRegex:     hostRegex,
		hostTokens:    hostTokens,
		hostVariables: hostVariables,
		variables:     variables,
	}

	if regex != "" {
		cr.compiledRegex = regexp.MustCompile(regex)
	}

	if hostRegex != "" {
		cr.compiledHostRegex = regexp.MustCompile(hostRegex)
	}

	return cr
}

// StaticPrefix returns the longest static prefix the path regex can begin with.
// Empty for routes that begin with a variable.
func (c *CompiledRoute) StaticPrefix() string { return c.staticPrefix }

// Regex returns the path-matching regular expression as a Go RE2 source string.
func (c *CompiledRoute) Regex() string { return c.regex }

// Tokens returns the path tokens in source order.
func (c *CompiledRoute) Tokens() []Token { return c.tokens }

// PathVariables returns the variable names declared in the path pattern.
func (c *CompiledRoute) PathVariables() []string { return c.pathVariables }

// HostRegex returns the host-matching regular expression, or the empty string
// if the route has no host constraint.
func (c *CompiledRoute) HostRegex() string { return c.hostRegex }

// HostTokens returns the host tokens in source order.
func (c *CompiledRoute) HostTokens() []Token { return c.hostTokens }

// HostVariables returns the variable names declared in the host pattern.
func (c *CompiledRoute) HostVariables() []string { return c.hostVariables }

// Variables returns the union of host and path variable names in source order
// (deduplicated).
func (c *CompiledRoute) Variables() []string { return c.variables }

// CompiledRegex returns the lazily-cached compiled path regex.
func (c *CompiledRoute) CompiledRegex() *regexp.Regexp { return c.compiledRegex }

// CompiledHostRegex returns the lazily-cached compiled host regex, or nil.
func (c *CompiledRoute) CompiledHostRegex() *regexp.Regexp { return c.compiledHostRegex }
