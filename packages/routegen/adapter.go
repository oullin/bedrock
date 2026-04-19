package routegen

import (
	"net/url"
	"strings"

	"github.com/bedrock/packages/routing"
)

// AdapterOptions configures how bedrock routes are converted to RouteInfo.
type AdapterOptions struct {
	// AppURL is the base application URL (e.g. "http://localhost:8000/v2").
	// When set, its path component is prepended to all generated URIs,
	// and its host/scheme override routes that do not have their own domain.
	AppURL string

	// ForcedScheme overrides the scheme for all routes ("http://", "https://").
	ForcedScheme string
}

// FromRouteCollection converts a bedrock RouteCollectionInterface into a slice
// of RouteInfo values suitable for passing to Generate.
//
// It mirrors the route mapping logic in GenerateCommand::handle() from the PHP
// implementation, extracting controller metadata, parameters, defaults, and
// domain information from each registered route.
func FromRouteCollection(collection routing.RouteCollectionInterface, opts AdapterOptions) []*RouteInfo {
	routes := collection.GetRoutes()
	out := make([]*RouteInfo, 0, len(routes))

	basePath, forcedScheme := parseAppURL(opts.AppURL)

	if opts.ForcedScheme != "" {
		forcedScheme = opts.ForcedScheme
	}

	for _, r := range routes {
		info := adaptRoute(r, basePath, forcedScheme)
		out = append(out, info)
	}

	return out
}

// adaptRoute converts a single routing.Route to a RouteInfo.
func adaptRoute(r *routing.Route, basePath, forcedScheme string) *RouteInfo {
	// Lower-case HTTP methods.
	methods := make([]string, len(r.HTTPMethods))

	for i, m := range r.HTTPMethods {
		methods[i] = strings.ToLower(m)
	}

	// Controller string ("Class@Method" or "Closure").
	controller := r.GetActionName()
	actionMethod := r.GetActionMethod()

	// Invokable detection: the Go routing package registers invokable
	// controllers with method "Invoke" (capital I) appended by ParseAction.
	isInvokable := actionMethod == "Invoke" || actionMethod == "__invoke"

	// Route name — strip generated:: prefix and trailing dots.
	name := cleanRouteName(r.GetName())

	// Parameters: walk ParameterNames and check URI for optional marker.
	paramNames := r.ParameterNames()
	bindingFields := r.BindingFields()
	defaults := r.DefaultValues
	params := make([]Param, 0, len(paramNames))

	for _, pName := range paramNames {
		optional := isOptionalParam(r.Uri, pName) || hasDefault(defaults, pName)
		key := bindingFields[pName]
		defaultVal := defaultString(defaults, pName)

		params = append(params, Param{
			Name:     pName,
			Optional: optional,
			Key:      key,
			Default:  defaultVal,
		})
	}

	// Domain handling.
	domain := r.GetDomain()
	scheme := ""

	if domain != "" {
		scheme = "//"
	}

	if forcedScheme != "" && domain == "" {
		// Only apply forced scheme when route has no own domain.
		scheme = forcedScheme
	}

	return &RouteInfo{
		URI:         r.Uri,
		Methods:     methods,
		Name:        name,
		Controller:  controller,
		IsInvokable: isInvokable,
		Params:      params,
		Domain:      domain,
		Scheme:      scheme,
		BasePath:    basePath,
	}
}

// isOptionalParam reports whether the named parameter appears as {name?}
// (with trailing question mark) in the URI string.
func isOptionalParam(uri, name string) bool {
	return strings.Contains(uri, "{"+name+"?}")
}

// hasDefault reports whether name has an entry in the defaults map.
func hasDefault(defaults map[string]any, name string) bool {
	if defaults == nil {
		return false
	}

	_, ok := defaults[name]

	return ok
}

// defaultString returns the string representation of a default value,
// or empty string when none exists.
func defaultString(defaults map[string]any, name string) string {
	if defaults == nil {
		return ""
	}

	v, ok := defaults[name]

	if !ok {
		return ""
	}

	switch s := v.(type) {
	case string:
		return s
	case int:
		return strings.TrimRight(strings.TrimRight(
			strings.TrimSuffix(string(rune('0'+s)), "0"), "."), "")
	default:
		return ""
	}
}

// parseAppURL extracts the base path and scheme from an APP_URL string.
// "http://localhost:8000/v2" → basePath="/v2", scheme=""
// "https://example.com"     → basePath="", scheme=""
func parseAppURL(appURL string) (basePath, scheme string) {
	if appURL == "" {
		return "", ""
	}

	if strings.HasPrefix(appURL, "//") {
		appURL = "http:" + appURL
	}

	u, err := url.Parse(appURL)

	if err != nil {
		return "", ""
	}

	path := strings.TrimRight(u.Path, "/")

	if path == "/" || path == "" {
		return "", ""
	}

	return path, ""
}
