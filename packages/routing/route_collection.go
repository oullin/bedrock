package routing

import (
	"strings"

	"github.com/bedrock/packages/routing/matching"
)

// RouteCollection is the in-memory route store used by [Router] during
// registration and request dispatching.
//
// Internally it maintains four indexes that mirror the PHP class:
//
//   - routes:     methodVerb → (domain+uri → Route), preserving insertion order
//     within each verb so dispatch is deterministic.
//   - allRoutes:  every route in registration order, keyed by methods+uri
//     so re-registration of the same route overwrites in place.
//   - nameList:   route name → Route, populated when [Route.Name] is set.
//   - actionList: "Controller@method" → Route, populated when a controller
//     action is supplied.
//
// Mirrors Illuminate\Routing\RouteCollection.
type RouteCollection struct {
	AbstractRouteCollection

	// Insertion order for each verb. We keep an explicit slice alongside the
	// map because Go maps are unordered, which would break dispatch.
	routesByMethod  map[string][]*Route
	routesByMethodK map[string]map[string]int // method → key → index in slice

	allRoutes  []*Route
	allKeys    map[string]int // method-set+domain+uri → index in allRoutes
	nameList   map[string]*Route
	actionList map[string]*Route
}

// NewRouteCollection constructs an empty collection.
func NewRouteCollection() *RouteCollection {
	return &RouteCollection{
		routesByMethod:  map[string][]*Route{},
		routesByMethodK: map[string]map[string]int{},
		allKeys:         map[string]int{},
		nameList:        map[string]*Route{},
		actionList:      map[string]*Route{},
	}
}

// Add registers a route and returns it. Mirrors RouteCollection::add.
func (c *RouteCollection) Add(route *Route) *Route {
	c.addToCollections(route)
	c.addLookups(route)

	return route
}

func (c *RouteCollection) addToCollections(route *Route) {
	domainAndUri := route.GetDomain() + route.Uri

	for _, method := range route.Methods() {
		key := domainAndUri

		if existing, ok := c.routesByMethodK[method]; ok {
			if idx, found := existing[key]; found {
				c.routesByMethod[method][idx] = route

				continue
			}
		} else {
			c.routesByMethodK[method] = map[string]int{}
		}

		idx := len(c.routesByMethod[method])
		c.routesByMethod[method] = append(c.routesByMethod[method], route)
		c.routesByMethodK[method][key] = idx
	}

	allKey := strings.Join(route.Methods(), "|") + domainAndUri

	if idx, ok := c.allKeys[allKey]; ok {
		c.allRoutes[idx] = route

		return
	}

	c.allKeys[allKey] = len(c.allRoutes)
	c.allRoutes = append(c.allRoutes, route)
}

func (c *RouteCollection) addLookups(route *Route) {
	if name := route.GetName(); name != "" {
		if _, ok := c.nameList[name]; !ok {
			c.nameList[name] = route
		}
	}

	if v, ok := route.ActionMap["controller"]; ok {
		if controller, ok := v.(string); ok && controller != "" {
			controller = strings.TrimLeft(controller, `\`)

			if _, present := c.actionList[controller]; !present {
				c.actionList[controller] = route
			}
		}
	}
}

// RefreshNameLookups rebuilds the name → Route table from the current routes.
func (c *RouteCollection) RefreshNameLookups() {
	c.nameList = map[string]*Route{}

	for _, r := range c.allRoutes {
		if name := r.GetName(); name != "" {
			if _, ok := c.nameList[name]; !ok {
				c.nameList[name] = r
			}
		}
	}
}

// RefreshActionLookups rebuilds the action → Route table.
func (c *RouteCollection) RefreshActionLookups() {
	c.actionList = map[string]*Route{}

	for _, r := range c.allRoutes {
		if v, ok := r.ActionMap["controller"]; ok {
			if controller, ok := v.(string); ok && controller != "" {
				controller = strings.TrimLeft(controller, `\`)

				if _, present := c.actionList[controller]; !present {
					c.actionList[controller] = r
				}
			}
		}
	}
}

// Match returns the route that matches the request or an error.
func (c *RouteCollection) Match(request matching.MatchableRequest) (*Route, error) {
	routes := c.Get(request.Method())
	matched := c.MatchAgainstRoutes(routes, request, true)

	return c.HandleMatchedRoute(c.Get, request, matched)
}

// Get returns routes for the given method, or all routes when method is "".
func (c *RouteCollection) Get(method string) []*Route {
	if method == "" {
		return c.GetRoutes()
	}

	return c.routesByMethod[method]
}

// HasNamedRoute reports whether a name is registered.
//
// If the lookup misses, the name table is refreshed once before giving up —
// this handles the common pattern of calling [Route.Name] AFTER the route was
// added to the collection.
func (c *RouteCollection) HasNamedRoute(name string) bool {
	return c.GetByName(name) != nil
}

// GetByName returns the named route or nil. Lazy-refreshes the lookup table
// once if the first probe misses.
func (c *RouteCollection) GetByName(name string) *Route {
	if r, ok := c.nameList[name]; ok {
		return r
	}

	c.RefreshNameLookups()

	return c.nameList[name]
}

// GetByAction returns the action's route or nil.
func (c *RouteCollection) GetByAction(action string) *Route { return c.actionList[action] }

// GetRoutes returns every route in registration order.
func (c *RouteCollection) GetRoutes() []*Route {
	out := make([]*Route, len(c.allRoutes))
	copy(out, c.allRoutes)

	return out
}

// GetRoutesByMethod returns the method → routes map.
func (c *RouteCollection) GetRoutesByMethod() map[string][]*Route {
	out := make(map[string][]*Route, len(c.routesByMethod))

	for k, v := range c.routesByMethod {
		clone := make([]*Route, len(v))
		copy(clone, v)
		out[k] = clone
	}

	return out
}

// GetRoutesByName returns the name → route map.
func (c *RouteCollection) GetRoutesByName() map[string]*Route {
	out := make(map[string]*Route, len(c.nameList))

	for k, v := range c.nameList {
		out[k] = v
	}

	return out
}

// Count reports the total number of routes (mirrors PHP Countable).
func (c *RouteCollection) Count() int { return len(c.allRoutes) }
