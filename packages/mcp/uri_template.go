package mcp

import (
	"fmt"
	"regexp"
	"strings"
)

// UriTemplate implements RFC 6570 Level 1 URI templates: patterns that
// contain {varname} placeholders. It supports matching an incoming URI
// against a template and extracting the variable bindings, as well as
// expanding a template given a variable map.
type UriTemplate struct {
	raw      string
	pattern  *regexp.Regexp
	varNames []string
}

// NewUriTemplate compiles a URI template string. Returns an error if the
// template is syntactically invalid (e.g. empty variable name).
func NewUriTemplate(template string) (*UriTemplate, error) {
	varNames := []string{}
	// Build a regex by replacing each {varname} placeholder with a named
	// capture group. The rest of the template is regex-escaped to prevent
	// ReDoS or accidental wildcards.
	var sb strings.Builder
	sb.WriteString("^")

	remaining := template
	for {
		open := strings.Index(remaining, "{")
		if open == -1 {
			sb.WriteString(regexp.QuoteMeta(remaining))
			break
		}
		close := strings.Index(remaining[open:], "}")
		if close == -1 {
			return nil, fmt.Errorf("mcp: uri template missing closing brace in %q", template)
		}
		close += open

		varName := remaining[open+1 : close]
		if varName == "" {
			return nil, fmt.Errorf("mcp: uri template has empty variable name in %q", template)
		}
		if strings.ContainsAny(varName, "{}") {
			return nil, fmt.Errorf("mcp: uri template has nested braces in %q", template)
		}

		sb.WriteString(regexp.QuoteMeta(remaining[:open]))
		sb.WriteString("(?P<")
		sb.WriteString(regexp.QuoteMeta(varName))
		sb.WriteString(">[^/]+)")

		varNames = append(varNames, varName)
		remaining = remaining[close+1:]
	}

	sb.WriteString("$")

	pattern, err := regexp.Compile(sb.String())
	if err != nil {
		return nil, fmt.Errorf("mcp: uri template compile error: %w", err)
	}

	return &UriTemplate{raw: template, pattern: pattern, varNames: varNames}, nil
}

// Match reports whether uri matches the template. When it does, vars contains
// the extracted variable bindings.
func (t *UriTemplate) Match(uri string) (vars map[string]string, ok bool) {
	m := t.pattern.FindStringSubmatch(uri)
	if m == nil {
		return nil, false
	}
	vars = make(map[string]string, len(t.varNames))
	for _, name := range t.varNames {
		idx := t.pattern.SubexpIndex(regexp.QuoteMeta(name))
		if idx >= 0 && idx < len(m) {
			vars[name] = m[idx]
		}
	}
	return vars, true
}

// Expand fills the template with the given variable bindings and returns the
// resulting URI string. Variables missing from vars are left as-is.
func (t *UriTemplate) Expand(vars map[string]string) string {
	result := t.raw
	for k, v := range vars {
		result = strings.ReplaceAll(result, "{"+k+"}", v)
	}
	return result
}

// Template returns the original template string.
func (t *UriTemplate) Template() string { return t.raw }

// IsTemplate reports whether a string contains at least one {varname}
// placeholder. Use this to distinguish static URIs from URI templates.
func IsTemplate(uri string) bool {
	open := strings.Index(uri, "{")
	if open == -1 {
		return false
	}
	close := strings.Index(uri[open:], "}")
	return close > 1 // at least one character between braces
}
