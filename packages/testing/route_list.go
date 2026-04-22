package testing

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// RouteListEntry captures the small slice of route metadata needed by the
// route-list compliance tests.
type RouteListEntry struct {
	Name          string
	Action        string
	URI           string
	Closure       bool
	Path          string
	Vendor        bool
	BindingFields map[string]string
}

// RouteListOptions controls the route-list renderer.
type RouteListOptions struct {
	Verbose      bool
	JSON         bool
	NameFilter   string
	ActionFilter string
	ExceptVendor bool
}

// RenderRouteList renders a filtered list of routes in either plain text or
// JSON form.
func RenderRouteList(routes []RouteListEntry, opts RouteListOptions) (string, error) {
	filtered := make([]RouteListEntry, 0, len(routes))
	for _, route := range routes {
		if opts.ExceptVendor && route.Vendor {
			continue
		}

		if opts.NameFilter != "" && !strings.Contains(route.Name, opts.NameFilter) {
			continue
		}

		if opts.ActionFilter != "" && !strings.Contains(route.Action, opts.ActionFilter) {
			continue
		}

		filtered = append(filtered, route)
	}

	if opts.JSON {
		payload := make([]map[string]any, 0, len(filtered))
		for _, route := range filtered {
			item := map[string]any{
				"action": route.Action,
				"name":   route.Name,
				"uri":    route.URI,
			}

			if route.Closure {
				item["path"] = route.Path
			} else {
				item["path"] = nil
			}

			if len(route.BindingFields) > 0 {
				fields := make(map[string]string, len(route.BindingFields))
				for key, value := range route.BindingFields {
					fields[key] = value
				}
				item["bindingFields"] = fields
			}

			payload = append(payload, item)
		}

		data, err := json.Marshal(payload)
		if err != nil {
			return "", err
		}

		return string(data), nil
	}

	lines := make([]string, 0, len(filtered))
	for _, route := range filtered {
		parts := []string{route.Name, route.Action, route.URI}

		if opts.Verbose && route.Closure && route.Path != "" {
			parts = append(parts, route.Path)
		}

		if len(route.BindingFields) > 0 {
			keys := make([]string, 0, len(route.BindingFields))
			for key := range route.BindingFields {
				keys = append(keys, key)
			}
			sort.Strings(keys)

			fields := make([]string, 0, len(keys))
			for _, key := range keys {
				fields = append(fields, fmt.Sprintf("%s:%s", key, route.BindingFields[key]))
			}

			parts = append(parts, strings.Join(fields, ","))
		}

		lines = append(lines, strings.Join(parts, " | "))
	}

	return strings.Join(lines, "\n"), nil
}
