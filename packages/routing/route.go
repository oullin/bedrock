package routing

import (
	nethttp "net/http"
	"regexp"
	"strings"
)

// Route represents a single registered route with its metadata, constraints,
// middleware, and handler. It supports fluent chaining for configuration.
type Route struct {
	name              string
	methods           []string
	uri               string
	handler           HandlerFunc
	middleware        []MiddlewareFunc
	withoutMiddleware []MiddlewareFunc
	constraints       map[string]string
	compiled          map[string]*regexp.Regexp
	domain            string
	prefix            string
	defaults          map[string]string
	fallback          bool
	namePrefix        string
	collection        *RouteCollection
}

// newRoute creates a route with the given methods, URI, and handler.
func newRoute(methods []string, uri string, handler HandlerFunc) *Route {
	return &Route{
		methods:     methods,
		uri:         uri,
		handler:     handler,
		constraints: make(map[string]string),
		compiled:    make(map[string]*regexp.Regexp),
		defaults:    make(map[string]string),
	}
}

// Name sets the route name and returns the route for chaining.
func (r *Route) Name(name string) *Route {
	if r.namePrefix != "" {
		r.name = r.namePrefix + "." + name
	} else {
		r.name = name
	}

	if r.collection != nil {
		r.collection.updateName(r)
	}

	return r
}

// GetName returns the route name.
func (r *Route) GetName() string {
	return r.name
}

// Methods returns the HTTP methods this route responds to.
func (r *Route) Methods() []string {
	return r.methods
}

// URI returns the route URI pattern.
func (r *Route) URI() string {
	return r.uri
}

// GetDomain returns the route domain constraint.
func (r *Route) GetDomain() string {
	return r.domain
}

// GetPrefix returns the route prefix.
func (r *Route) GetPrefix() string {
	return r.prefix
}

// Where adds a regex constraint for a route parameter.
func (r *Route) Where(param, pattern string) *Route {
	r.constraints[param] = pattern
	r.compiled[param] = regexp.MustCompile("^" + pattern + "$")

	return r
}

// WhereNumber constrains the given parameters to numeric values.
func (r *Route) WhereNumber(params ...string) *Route {
	for _, p := range params {
		r.Where(p, `[0-9]+`)
	}

	return r
}

// WhereAlpha constrains the given parameters to alphabetic values.
func (r *Route) WhereAlpha(params ...string) *Route {
	for _, p := range params {
		r.Where(p, `[a-zA-Z]+`)
	}

	return r
}

// WhereAlphaNumeric constrains the given parameters to alphanumeric values.
func (r *Route) WhereAlphaNumeric(params ...string) *Route {
	for _, p := range params {
		r.Where(p, `[a-zA-Z0-9]+`)
	}

	return r
}

// WhereIn constrains the given parameter to a set of allowed values.
func (r *Route) WhereIn(param string, values []string) *Route {
	pattern := "(" + strings.Join(values, "|") + ")"
	r.Where(param, pattern)

	return r
}

// WhereUUID constrains the given parameters to UUID format.
func (r *Route) WhereUUID(params ...string) *Route {
	for _, p := range params {
		r.Where(p, `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)
	}

	return r
}

// Middleware appends middleware to the route.
func (r *Route) Middleware(mw ...MiddlewareFunc) *Route {
	r.middleware = append(r.middleware, mw...)

	return r
}

// WithoutMiddleware marks middleware to be excluded from this route.
func (r *Route) WithoutMiddleware(mw ...MiddlewareFunc) *Route {
	r.withoutMiddleware = append(r.withoutMiddleware, mw...)

	return r
}

// Domain sets a domain constraint for the route.
func (r *Route) Domain(domain string) *Route {
	r.domain = domain

	return r
}

// Defaults sets default parameter values.
func (r *Route) Defaults(defaults map[string]string) *Route {
	for k, v := range defaults {
		r.defaults[k] = v
	}

	return r
}

// SetFallback marks the route as a fallback route.
func (r *Route) SetFallback(fallback bool) *Route {
	r.fallback = fallback

	return r
}

// IsFallback returns whether this is a fallback route.
func (r *Route) IsFallback() bool {
	return r.fallback
}

// HasParameters returns true if the route URI contains parameters.
func (r *Route) HasParameters() bool {
	return strings.Contains(r.uri, "{")
}

// ParameterNames returns the parameter names extracted from the URI pattern.
func (r *Route) ParameterNames() []string {
	matches := paramRegex.FindAllStringSubmatch(r.uri, -1)
	params := make([]string, 0, len(matches))

	for _, m := range matches {
		params = append(params, m[1])
	}

	return params
}

// GetConstraints returns the parameter constraints map.
func (r *Route) GetConstraints() map[string]string {
	result := make(map[string]string, len(r.constraints))

	for k, v := range r.constraints {
		result[k] = v
	}

	return result
}

// GetDefaults returns the default parameter values.
func (r *Route) GetDefaults() map[string]string {
	result := make(map[string]string, len(r.defaults))

	for k, v := range r.defaults {
		result[k] = v
	}

	return result
}

// GetMiddleware returns the route-level middleware.
func (r *Route) GetMiddleware() []MiddlewareFunc {
	return r.middleware
}

// MatchesConstraints checks whether the given parameter values satisfy all
// compiled constraints. Returns true if there are no constraints.
func (r *Route) MatchesConstraints(params map[string]string) bool {
	for param, re := range r.compiled {
		val, ok := params[param]

		if !ok {
			if _, hasDef := r.defaults[param]; hasDef {
				continue
			}

			return false
		}

		if !re.MatchString(val) {
			return false
		}
	}

	return true
}

// Matches checks if the given request satisfies this route's constraints.
func (r *Route) Matches(req *nethttp.Request) bool {
	if r.domain != "" && !matchesDomain(r.domain, req.Host) {
		return false
	}

	params := extractParams(r.uri, req)

	return r.MatchesConstraints(params)
}

// matchesDomain checks if the request host matches the route domain pattern.
func matchesDomain(pattern, host string) bool {
	host = strings.Split(host, ":")[0]
	pattern = strings.Split(pattern, ":")[0]

	if !strings.Contains(pattern, "{") {
		return strings.EqualFold(pattern, host)
	}

	regexStr := "^" + paramRegex.ReplaceAllString(pattern, `([^.]+)`) + "$"

	re, err := regexp.Compile(regexStr)

	if err != nil {
		return false
	}

	return re.MatchString(host)
}

// extractParams extracts path parameter values from the request using the
// Go 1.22+ PathValue API.
func extractParams(uri string, req *nethttp.Request) map[string]string {
	matches := paramRegex.FindAllStringSubmatch(uri, -1)
	params := make(map[string]string, len(matches))

	for _, m := range matches {
		params[m[1]] = req.PathValue(m[1])
	}

	return params
}
