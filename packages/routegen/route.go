package routegen

import "regexp"

// Route describes a single named route in the registry.

type Route struct {
	Name    string `json:"name"`
	Method  string `json:"method"`
	Pattern string `json:"pattern"`
}

var paramRegex = regexp.MustCompile(`\{(\w+)\}`)

// Params extracts the parameter names from the route pattern.
func (r Route) Params() []string {
	matches := paramRegex.FindAllStringSubmatch(r.Pattern, -1)
	params := make([]string, 0, len(matches))

	for _, match := range matches {
		params = append(params, match[1])
	}

	return params
}
