package routing

import "strings"

// MiddlewareNameResolver expands middleware aliases and group names into
// concrete class names (or closures).
//
// Ref: @bedrock/code-0320
type MiddlewareNameResolver struct{}

// Resolve translates name into one of:
//   - a closure ([]any{name}) — returned unchanged for inline middleware
//   - a single resolved class string (alias lookup)
//   - a flat slice of class strings (group expansion)
//
// Ref: @bedrock/code-0320
// callers can flatten it into the middleware pipeline uniformly.
func (MiddlewareNameResolver) Resolve(name any, aliases map[string]any, groups map[string][]any) []any {
	switch v := name.(type) {
	case func():
		return []any{v}
	case string:
		// Alias to closure → return closure directly.
		if alias, ok := aliases[v]; ok {
			if _, isFunc := alias.(func()); isFunc {
				return []any{alias}
			}
		}
		// Group expansion.
		if _, ok := groups[v]; ok {
			return parseMiddlewareGroup(v, aliases, groups, map[string]bool{})
		}
		// String → "alias[:params]" form: split params, look up alias, re-attach.
		base, params := splitColon(v)

		if alias, ok := aliases[base]; ok {
			if s, ok := alias.(string); ok {
				return []any{s + params}
			}

			return []any{alias}
		}

		return []any{v}
	default:
		return []any{name}
	}
}

func parseMiddlewareGroup(name string, aliases map[string]any, groups map[string][]any, seen map[string]bool) []any {
	if seen[name] {
		return nil
	}

	seen[name] = true
	out := make([]any, 0, len(groups[name]))

	for _, m := range groups[name] {
		s, isStr := m.(string)

		if isStr {
			if _, isGroup := groups[s]; isGroup {
				out = append(out, parseMiddlewareGroup(s, aliases, groups, seen)...)

				continue
			}

			base, params := splitColon(s)

			if alias, ok := aliases[base]; ok {
				if cls, ok := alias.(string); ok {
					out = append(out, cls+params)

					continue
				}

				out = append(out, alias)

				continue
			}

			out = append(out, s)

			continue
		}

		out = append(out, m)
	}

	return out
}

func splitColon(s string) (base, params string) {
	if idx := strings.Index(s, ":"); idx >= 0 {
		return s[:idx], s[idx:]
	}

	return s, ""
}
