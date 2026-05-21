package routing

import (
	"strings"

	"github.com/bedrock/packages/routing/compiler"
)

// boundRoute is the minimum surface RouteParameterBinder needs from a route.
//
// The concrete [Route] type (M2) satisfies this interface; tests use a
// fake. Decoupling here keeps M1 buildable before [Route] exists.
type boundRoute interface {
	Compiled() *compiler.CompiledRoute
	ParameterNames() []string
	DefaultsMap() map[string]any
}

// boundRequest is the minimum surface RouteParameterBinder needs from a
// request. httpx.Request satisfies this; tests provide a fake.
type boundRequest interface {
	DecodedPath() string
	Host() string
}

// RouteParameterBinder extracts route parameter values from a request by
// running the route's compiled regex against the request's path (and host,
// when the route has a host pattern), then merging in declared defaults.
//
// Mirrors @bedrock\Routing\RouteParameterBinder.
type RouteParameterBinder struct {
	route boundRoute
}

// NewRouteParameterBinder wraps the given route.
func NewRouteParameterBinder(route boundRoute) *RouteParameterBinder {
	return &RouteParameterBinder{route: route}
}

// Parameters returns the parameter map for a successful match. Missing
// optional parameters are filled in from the route's declared defaults.
func (b *RouteParameterBinder) Parameters(req boundRequest) map[string]string {
	parameters := b.bindPathParameters(req)

	if c := b.route.Compiled(); c != nil && c.HostRegex() != "" {
		parameters = b.bindHostParameters(req, parameters)
	}

	return b.replaceDefaults(parameters)
}

func (b *RouteParameterBinder) bindPathParameters(req boundRequest) map[string]string {
	c := b.route.Compiled()

	if c == nil || c.CompiledRegex() == nil {
		return map[string]string{}
	}

	path := "/" + strings.TrimLeft(req.DecodedPath(), "/")
	m := c.CompiledRegex().FindStringSubmatch(path)

	if m == nil {
		return map[string]string{}
	}

	return b.matchToKeys(c.CompiledRegex().SubexpNames(), m)
}

func (b *RouteParameterBinder) bindHostParameters(req boundRequest, parameters map[string]string) map[string]string {
	c := b.route.Compiled()
	hostRe := c.CompiledHostRegex()

	if hostRe == nil {
		return parameters
	}

	m := hostRe.FindStringSubmatch(req.Host())

	if m == nil {
		return parameters
	}

	hostParams := b.matchToKeys(hostRe.SubexpNames(), m)

	for k, v := range parameters {
		hostParams[k] = v
	}

	return hostParams
}

func (b *RouteParameterBinder) matchToKeys(names, matches []string) map[string]string {
	wanted := map[string]struct{}{}

	for _, n := range b.route.ParameterNames() {
		wanted[n] = struct{}{}
	}

	out := map[string]string{}

	for i, name := range names {
		if name == "" || i >= len(matches) {
			continue
		}

		if _, ok := wanted[name]; !ok {
			continue
		}

		if matches[i] == "" {
			continue
		}

		out[name] = matches[i]
	}

	return out
}

func (b *RouteParameterBinder) replaceDefaults(parameters map[string]string) map[string]string {
	defaults := b.route.DefaultsMap()

	for k, v := range parameters {
		if v == "" {
			if dv, ok := defaults[k]; ok {
				parameters[k] = anyToString(dv)
			}
		}
	}

	for k, v := range defaults {
		if _, ok := parameters[k]; !ok {
			parameters[k] = anyToString(v)
		}
	}

	return parameters
}

func anyToString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}

	return ""
}
