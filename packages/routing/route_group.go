package routing

import "strings"

// It holds no state — its
// public surface is the static [RouteGroup.Merge] helper that the router uses
// when nesting [Router.Group] calls.
type RouteGroup struct{}

// MergeRouteGroup merges a new group attribute set into an old one.
//
// like "namespace", "prefix", "where", "as", "domain", "controller",
// "middleware". Unknown keys are passed through unchanged.
//
// prependExistingPrefix controls whether the old prefix is placed before the
// new prefix (the default for nested groups) or after it (used by
// [Router.mergeGroupAttributesIntoRoute] when merging the route's own
// attributes back into the group stack).
func MergeRouteGroup(new, old map[string]any, prependExistingPrefix bool) map[string]any {
	if old == nil {
		old = map[string]any{}
	}

	if new == nil {
		new = map[string]any{}
	}

	if _, ok := new["domain"]; ok {
		delete(old, "domain")
	}

	if _, ok := new["controller"]; ok {
		delete(old, "controller")
	}

	new = formatAs(new, old)
	new["namespace"] = formatNamespace(new, old)
	new["prefix"] = formatPrefix(new, old, prependExistingPrefix)
	new["where"] = formatWhere(new, old)

	out := map[string]any{}

	for k, v := range old {
		switch k {
		case "namespace", "prefix", "where", "as":
			continue
		}

		out[k] = v
	}

	for k, v := range new {
		// Mirror PHP array_merge_recursive: for "middleware" keys, append.
		if k == "middleware" {
			if existing, ok := out["middleware"]; ok {
				out["middleware"] = mergeMiddleware(existing, v)

				continue
			}
		}

		out[k] = v
	}

	return out
}

func formatNamespace(new, old map[string]any) any {
	if v, ok := new["namespace"]; ok {
		s, _ := v.(string)

		if oldNs, ok := old["namespace"]; ok && !strings.HasPrefix(s, `\`) {
			return strings.Trim(oldNs.(string), `\`) + `\` + strings.Trim(s, `\`)
		}

		return strings.Trim(s, `\`)
	}

	return old["namespace"]
}

func formatPrefix(new, old map[string]any, prependExisting bool) any {
	oldPrefix, _ := old["prefix"].(string)

	if prependExisting {
		if v, ok := new["prefix"]; ok {
			return strings.Trim(oldPrefix, "/") + "/" + strings.Trim(v.(string), "/")
		}

		return oldPrefix
	}

	if v, ok := new["prefix"]; ok {
		return strings.Trim(v.(string), "/") + "/" + strings.Trim(oldPrefix, "/")
	}

	return oldPrefix
}

func formatWhere(new, old map[string]any) map[string]string {
	out := map[string]string{}

	if v, ok := old["where"].(map[string]string); ok {
		for k, val := range v {
			out[k] = val
		}
	}

	if v, ok := new["where"].(map[string]string); ok {
		for k, val := range v {
			out[k] = val
		}
	}

	return out
}

func formatAs(new, old map[string]any) map[string]any {
	out := map[string]any{}

	for k, v := range new {
		out[k] = v
	}

	if v, ok := old["as"]; ok {
		oldAs, _ := v.(string)
		newAs, _ := out["as"].(string)
		out["as"] = oldAs + newAs
	}

	return out
}

func mergeMiddleware(a, b any) []any {
	out := []any{}

	if as, ok := a.([]any); ok {
		out = append(out, as...)
	}

	if bs, ok := b.([]any); ok {
		out = append(out, bs...)
	}

	return out
}
