package validation

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"

	contract "github.com/bedrock/packages/contracts/validation"
)

// ParsedRule represents a single rule with its name and parameters.
type ParsedRule struct {
	// Name is the StudlyCase rule name, e.g. "Required", "Max", "RequiredIf".
	Name string

	// Parameters holds the colon-separated arguments, e.g. ["255"] for "max:255".
	Parameters []string

	// Object is set when the rule was supplied as a ValidationRule object rather
	// than a string.
	Object contract.ValidationRule
}

// IsObject reports whether this rule was provided as a ValidationRule instance.
func (r ParsedRule) IsObject() bool {
	return r.Object != nil
}

// String returns a human-readable representation used in error messages.
func (r ParsedRule) String() string {
	if r.IsObject() {
		return fmt.Sprintf("%T", r.Object)
	}

	if len(r.Parameters) == 0 {
		return r.Name
	}

	return r.Name + ":" + strings.Join(r.Parameters, ",")
}

// Parse parses a single rule string like "max:255" or "required_if:field,value"
// into a ParsedRule.  The name is normalised to StudlyCase.
func Parse(rule string) ParsedRule {
	rule = strings.TrimSpace(rule)

	// regex rules keep the full string after "regex:" intact
	lower := strings.ToLower(rule)

	if strings.HasPrefix(lower, "regex:") || strings.HasPrefix(lower, "not_regex:") {
		parts := strings.SplitN(rule, ":", 2)

		return ParsedRule{
			Name:       StudlyCase(parts[0]),
			Parameters: []string{parts[1]},
		}
	}

	idx := strings.IndexByte(rule, ':')

	if idx == -1 {
		return ParsedRule{Name: StudlyCase(rule)}
	}

	name := rule[:idx]
	param := rule[idx+1:]

	return ParsedRule{
		Name:       StudlyCase(name),
		Parameters: parseParameters(param),
	}
}

// parseParameters splits a comma-separated parameter string, honouring
// double-quoted values (upstream uses CSV with " as quote and \ as escape).
func parseParameters(s string) []string {
	if s == "" {
		return nil
	}

	// Simple split — no quoting complexity needed for most rules.
	// For values that may contain commas (e.g. regex), callers use SplitN
	// at the caller site.
	return strings.Split(s, ",")
}

// Explode converts any rule representation into a slice of ParsedRule.
// Accepted input types:
//   - string:                  "required|email|max:255"
//   - []string:                []string{"required", "email"}
//   - []any:                   mix of strings and ValidationRule objects
//   - ValidationRule:          single rule object
//   - []ParsedRule:            passed through as-is
func Explode(rules any) []ParsedRule {
	switch v := rules.(type) {
	case string:
		return explodeString(v)
	case []string:
		return explodeStringSlice(v)
	case []any:
		return explodeAnySlice(v)
	case contract.ValidationRule:
		return []ParsedRule{{Object: v}}
	case []ParsedRule:
		return v
	case ParsedRule:
		return []ParsedRule{v}
	default:
		return nil
	}
}

func explodeString(s string) []ParsedRule {
	parts := strings.Split(s, "|")
	out := make([]ParsedRule, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)

		if p == "" {
			continue
		}

		out = append(out, Parse(p))
	}

	return out
}

func explodeStringSlice(ss []string) []ParsedRule {
	out := make([]ParsedRule, 0, len(ss))

	for _, s := range ss {
		out = append(out, explodeString(s)...)
	}

	return out
}

func explodeAnySlice(items []any) []ParsedRule {
	out := make([]ParsedRule, 0, len(items))

	for _, item := range items {
		switch v := item.(type) {
		case string:
			out = append(out, explodeString(v)...)
		case contract.ValidationRule:
			out = append(out, ParsedRule{Object: v})
		case ParsedRule:
			out = append(out, v)
		}
	}

	return out
}

// StudlyCase converts a snake_case or kebab-case rule name to StudlyCase
// (PascalCase).  For example: "required_if" → "RequiredIf", "max" → "Max".
func StudlyCase(s string) string {
	if s == "" {
		return ""
	}

	s = strings.ReplaceAll(s, "-", "_")
	parts := strings.Split(s, "_")

	var b strings.Builder

	for _, p := range parts {
		if p == "" {
			continue
		}

		runes := []rune(p)
		runes[0] = unicode.ToUpper(runes[0])
		b.WriteString(string(runes))
	}

	return b.String()
}

// ExpandWildcards resolves wildcard patterns like "items.*.name" against a
// flattened dot-notation data map, returning the concrete attribute paths that
// match.  If the pattern contains no wildcard it is returned unchanged.
func ExpandWildcards(pattern string, flat map[string]any) []string {
	if !strings.Contains(pattern, "*") {
		return []string{pattern}
	}

	// Build a regex from the pattern:  items.*.name → ^items\.\d+\.name$
	regexStr := wildcardToRegex(pattern)
	re, err := regexp.Compile(regexStr)

	if err != nil {
		return []string{pattern}
	}

	var matches []string

	for key := range flat {
		if re.MatchString(key) {
			matches = append(matches, key)
		}
	}

	return matches
}

func wildcardToRegex(pattern string) string {
	// Escape regex meta-characters except our own *
	var b strings.Builder

	b.WriteString("^")

	for _, ch := range pattern {
		switch ch {
		case '*':
			b.WriteString(`[^.]+`)
		case '.':
			b.WriteString(`\.`)
		case '?', '+', '(', ')', '[', ']', '{', '}', '^', '$', '|', '\\':
			b.WriteRune('\\')
			b.WriteRune(ch)
		default:
			b.WriteRune(ch)
		}
	}

	b.WriteString("$")

	return b.String()
}

// FlattenData converts a nested map[string]any into a flat dot-notation map.
// e.g. {"user": {"name": "Alice"}} → {"user.name": "Alice"}
func FlattenData(data map[string]any) map[string]any {
	out := make(map[string]any)
	flattenInto(out, data, "")

	return out
}

func flattenInto(out map[string]any, data map[string]any, prefix string) {
	for k, v := range data {
		key := k

		if prefix != "" {
			key = prefix + "." + k
		}

		switch child := v.(type) {
		case map[string]any:
			// Keep the map itself so top-level array/list/shape rules can run.
			out[key] = v
			flattenInto(out, child, key)
		case []any:
			for i, item := range child {
				indexKey := fmt.Sprintf("%s.%d", key, i)

				if nested, ok := item.(map[string]any); ok {
					flattenInto(out, nested, indexKey)
				} else {
					out[indexKey] = item
				}
			}

			// Also store the array itself under its key for array-level rules
			out[key] = v
		default:
			out[key] = v
		}
	}
}
