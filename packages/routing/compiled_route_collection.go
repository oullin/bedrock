package routing

import "github.com/bedrock/packages/routing/matching"

// CompiledRouteCollection is the production-mode counterpart to
// [RouteCollection]. In Laravel/Symfony it is backed by a dumped matcher; in
// the Go port it stores pre-bound route metadata in a slice and delegates
// matching to the same compiled regexes.
//
// The constructor mirrors Laravel's: it accepts a "compiled" payload (the
// dumped matcher data, opaque to consumers) and an "attributes" payload that
// the dumper produces alongside it. In Go we store a slice of Routes that the
// router built from the cached form; M11 will provide a real cache loader.
//
// Mirrors Illuminate\Routing\CompiledRouteCollection.
type CompiledRouteCollection struct {
	AbstractRouteCollection
	routes     []*Route
	nameList   map[string]*Route
	actionList map[string]*Route
}

// NewCompiledRouteCollection builds a compiled collection from the supplied
// routes. The "attributes" map is accepted for parity with the PHP signature
// but is currently unused; M11's cache loader will populate it.
func NewCompiledRouteCollection(routes []*Route, _ map[string]any) *CompiledRouteCollection {
	c := &CompiledRouteCollection{
		routes:     append([]*Route(nil), routes...),
		nameList:   map[string]*Route{},
		actionList: map[string]*Route{},
	}
	c.RefreshNameLookups()
	c.RefreshActionLookups()

	return c
}

// Add appends a route and updates the lookup tables.
func (c *CompiledRouteCollection) Add(route *Route) *Route {
	c.routes = append(c.routes, route)

	if name := route.GetName(); name != "" {
		if _, ok := c.nameList[name]; !ok {
			c.nameList[name] = route
		}
	}

	if v, ok := route.ActionMap["controller"]; ok {
		if s, ok := v.(string); ok && s != "" {
			if _, ok := c.actionList[s]; !ok {
				c.actionList[s] = route
			}
		}
	}

	return route
}

// RefreshNameLookups rebuilds the name index.
func (c *CompiledRouteCollection) RefreshNameLookups() {
	c.nameList = map[string]*Route{}

	for _, r := range c.routes {
		if name := r.GetName(); name != "" {
			if _, ok := c.nameList[name]; !ok {
				c.nameList[name] = r
			}
		}
	}
}

// RefreshActionLookups rebuilds the action index.
func (c *CompiledRouteCollection) RefreshActionLookups() {
	c.actionList = map[string]*Route{}

	for _, r := range c.routes {
		if v, ok := r.ActionMap["controller"]; ok {
			if s, ok := v.(string); ok && s != "" {
				if _, ok := c.actionList[s]; !ok {
					c.actionList[s] = r
				}
			}
		}
	}
}

// Match scans the compiled set in registration order, just like
// RouteCollection.Match.
func (c *CompiledRouteCollection) Match(request matching.MatchableRequest) (*Route, error) {
	matched := c.MatchAgainstRoutes(c.Get(request.Method()), request, true)

	return c.HandleMatchedRoute(c.Get, request, matched)
}

// Get returns routes filtered by method, or all routes when method is "".
func (c *CompiledRouteCollection) Get(method string) []*Route {
	if method == "" {
		return c.GetRoutes()
	}

	out := make([]*Route, 0, len(c.routes))

	for _, r := range c.routes {
		for _, m := range r.Methods() {
			if m == method {
				out = append(out, r)

				break
			}
		}
	}

	return out
}

// HasNamedRoute reports whether name is registered.
func (c *CompiledRouteCollection) HasNamedRoute(name string) bool {
	return c.nameList[name] != nil
}

// GetByName returns the named route or nil.
func (c *CompiledRouteCollection) GetByName(name string) *Route { return c.nameList[name] }

// GetByAction returns the route for the given action or nil.
func (c *CompiledRouteCollection) GetByAction(action string) *Route { return c.actionList[action] }

// GetRoutes returns a defensive copy of the internal slice.
func (c *CompiledRouteCollection) GetRoutes() []*Route {
	out := make([]*Route, len(c.routes))
	copy(out, c.routes)

	return out
}

// GetRoutesByMethod groups routes by HTTP verb.
func (c *CompiledRouteCollection) GetRoutesByMethod() map[string][]*Route {
	out := map[string][]*Route{}

	for _, r := range c.routes {
		for _, m := range r.Methods() {
			out[m] = append(out[m], r)
		}
	}

	return out
}

// GetRoutesByName returns the name index.
func (c *CompiledRouteCollection) GetRoutesByName() map[string]*Route {
	out := make(map[string]*Route, len(c.nameList))

	for k, v := range c.nameList {
		out[k] = v
	}

	return out
}

// Count reports the total number of routes.
func (c *CompiledRouteCollection) Count() int { return len(c.routes) }
