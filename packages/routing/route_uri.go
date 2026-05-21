package routing

import (
	"regexp"
	"strings"
)

// RouteUri is a parsed route URI pattern. The Uri field holds the URI with
// any "{param:field}" segments rewritten to the bare "{param}" form, and
// BindingFields maps a parameter name to its custom binding field (the part
// after the colon).
//
// Ref: @bedrock/code-0341
type RouteUri struct {
	Uri           string
	BindingFields map[string]string
}

// NewRouteUri constructs a RouteUri from an already-parsed URI and binding
// field map. Most callers want [ParseRouteUri].
func NewRouteUri(uri string, bindingFields map[string]string) *RouteUri {
	if bindingFields == nil {
		bindingFields = map[string]string{}
	}

	return &RouteUri{Uri: uri, BindingFields: bindingFields}
}

// routeUriPlaceholderRe matches a single "{name}" or "{name:field}" segment,
// optionally followed by "?" to mark the parameter as optional.
//
// Ref: @bedrock/code-0341
var routeUriPlaceholderRe = regexp.MustCompile(`\{([\w:]+?)\??\}`)

// ParseRouteUri parses a upstream-style URI pattern, extracting any
// custom binding fields ("{user:slug}") into BindingFields and rewriting the
// URI so downstream consumers see the canonical "{user}" form.
//
// Ref: @bedrock/code-0341
func ParseRouteUri(uri string) *RouteUri {
	bindingFields := map[string]string{}

	for _, match := range routeUriPlaceholderRe.FindAllString(uri, -1) {
		if !strings.Contains(match, ":") {
			continue
		}
		// Strip {, }, and trailing ? from the matched segment, then split on ":".
		trimmed := strings.TrimSuffix(strings.TrimSuffix(strings.TrimPrefix(match, "{"), "}"), "?")
		segments := strings.SplitN(trimmed, ":", 2)

		if len(segments) != 2 {
			continue
		}

		bindingFields[segments[0]] = segments[1]

		var replacement string

		if strings.Contains(match, "?") {
			replacement = "{" + segments[0] + "?}"
		} else {
			replacement = "{" + segments[0] + "}"
		}

		uri = strings.Replace(uri, match, replacement, 1)
	}

	return NewRouteUri(uri, bindingFields)
}
