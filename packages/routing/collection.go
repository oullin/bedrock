package routing

import (
	nethttp "net/http"
	"strings"
	"sync"
)

// RouteCollection manages registered routes and supports lookup by name,
// method, and request matching. It is safe for concurrent use.
type RouteCollection struct {
	mu        sync.RWMutex
	routes    []*Route
	nameMap   map[string]*Route
	methodMap map[string][]*Route
}

// NewRouteCollection creates an empty route collection.
func NewRouteCollection() *RouteCollection {
	return &RouteCollection{
		nameMap:   make(map[string]*Route),
		methodMap: make(map[string][]*Route),
	}
}

// Add registers a route in the collection, indexing it by name and method.
func (c *RouteCollection) Add(route *Route) *Route {
	c.mu.Lock()

	defer c.mu.Unlock()

	route.collection = c
	c.routes = append(c.routes, route)

	if route.name != "" {
		c.nameMap[route.name] = route
	}

	for _, method := range route.methods {
		m := strings.ToUpper(method)
		c.methodMap[m] = append(c.methodMap[m], route)
	}

	if len(route.methods) == 0 {
		for _, m := range allMethods {
			c.methodMap[m] = append(c.methodMap[m], route)
		}
	}

	return route
}

// Match finds a route matching the given request. Returns ErrMethodNotAllowed
// if the URI matches but the method does not, or ErrRouteNotFound if no route
// matches at all.
func (c *RouteCollection) Match(req *nethttp.Request) (*Route, error) {
	c.mu.RLock()

	defer c.mu.RUnlock()

	method := strings.ToUpper(req.Method)
	uriMatched := false

	candidates := c.methodMap[method]

	for _, route := range candidates {
		if matchesURI(route.uri, req.URL.Path) && route.Matches(req) {
			return route, nil
		}
	}

	for _, routes := range c.methodMap {
		for _, route := range routes {
			if matchesURI(route.uri, req.URL.Path) {
				uriMatched = true

				break
			}
		}

		if uriMatched {
			break
		}
	}

	if uriMatched {
		return nil, ErrMethodNotAllowed
	}

	return nil, ErrRouteNotFound
}

// ByName returns the route with the given name.
func (c *RouteCollection) ByName(name string) (*Route, bool) {
	c.mu.RLock()

	defer c.mu.RUnlock()

	route, ok := c.nameMap[name]

	return route, ok
}

// ByMethod returns all routes registered for the given HTTP method.
func (c *RouteCollection) ByMethod(method string) []*Route {
	c.mu.RLock()

	defer c.mu.RUnlock()

	routes := c.methodMap[strings.ToUpper(method)]
	result := make([]*Route, len(routes))
	copy(result, routes)

	return result
}

// HasNamedRoute returns true if a route with the given name exists.
func (c *RouteCollection) HasNamedRoute(name string) bool {
	c.mu.RLock()

	defer c.mu.RUnlock()

	_, ok := c.nameMap[name]

	return ok
}

// All returns all registered routes in insertion order.
func (c *RouteCollection) All() []*Route {
	c.mu.RLock()

	defer c.mu.RUnlock()

	result := make([]*Route, len(c.routes))
	copy(result, c.routes)

	return result
}

// Len returns the number of registered routes.
func (c *RouteCollection) Len() int {
	c.mu.RLock()

	defer c.mu.RUnlock()

	return len(c.routes)
}

// updateName updates the name index when a route's name changes.
func (c *RouteCollection) updateName(route *Route) {
	c.mu.Lock()

	defer c.mu.Unlock()

	if route.name != "" {
		c.nameMap[route.name] = route
	}
}

// RefreshNameLookups rebuilds the name index from the current routes.
func (c *RouteCollection) RefreshNameLookups() {
	c.mu.Lock()

	defer c.mu.Unlock()

	c.nameMap = make(map[string]*Route, len(c.routes))

	for _, route := range c.routes {
		if route.name != "" {
			c.nameMap[route.name] = route
		}
	}
}

// matchesURI performs a basic check whether a URI pattern could match a path.
// This is a simplified matcher—full matching is handled by ServeMux.
func matchesURI(pattern, path string) bool {
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	if len(patternParts) != len(pathParts) {
		return false
	}

	for i, part := range patternParts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			continue
		}

		if part != pathParts[i] {
			return false
		}
	}

	return true
}

var allMethods = []string{
	nethttp.MethodGet,
	nethttp.MethodHead,
	nethttp.MethodPost,
	nethttp.MethodPut,
	nethttp.MethodPatch,
	nethttp.MethodDelete,
	nethttp.MethodOptions,
}
