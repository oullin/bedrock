package routing

import "strings"

// ResourceRegistrar generates the seven standard RESTful routes for a
// Ref: @bedrock/code-0330
// The Go port covers: register/singleton, only/except, names/parameters,
// shallow nesting, and prefixed names (e.g. "users.posts" → nested URI).
// Trashed-binding behavior and per-method middleware filtering are wired in
// alongside [middleware/SubstituteBindings] in M9/M10.
type ResourceRegistrar struct {
	router *Router

	// parameters maps a resource name to the URL placeholder name. nil means
	// "use the singular form of the resource name".
	parameters map[string]string
}

var (
	// resourceDefaults is the canonical action set for a RESTful resource.
	resourceDefaults = []string{"index", "create", "store", "show", "edit", "update", "destroy"}
	// singletonResourceDefaults is the action set for a singleton resource.
	singletonResourceDefaults = []string{"show", "edit", "update"}
)

// resourceVerbs is the per-action mapping of URI verbs (e.g. /users/create).
var resourceVerbs = map[string]string{
	"create": "create",
	"edit":   "edit",
}

// NewResourceRegistrar constructs a registrar bound to the given router.
func NewResourceRegistrar(router *Router) *ResourceRegistrar {
	return &ResourceRegistrar{router: router, parameters: map[string]string{}}
}

// Register emits the seven standard RESTful routes for the resource.
//
// Ref: @bedrock/code-0330
//   - "only":     []string of action names to keep
//   - "except":   []string of action names to drop
//   - "names":    map[string]string overriding the route name per action
//   - "parameters": map[string]string overriding the URL placeholder per resource
//   - "as":       string prefix prepended to every route name
//   - "shallow":  bool — for nested resources, drop the parent prefix on
//     show/edit/update/destroy.
func (rr *ResourceRegistrar) Register(name, controller string, options map[string]any) *RouteCollection {
	if strings.Contains(name, "/") {
		return rr.prefixedResource(name, controller, options)
	}

	base := rr.GetResourceWildcard(lastSegment(name, "."))
	methods := rr.getResourceMethods(resourceDefaults, options)

	collection := NewRouteCollection()

	for _, m := range methods {
		var route *Route

		switch m {
		case "index":
			route = rr.addResourceIndex(name, base, controller, options)
		case "create":
			route = rr.addResourceCreate(name, base, controller, options)
		case "store":
			route = rr.addResourceStore(name, base, controller, options)
		case "show":
			route = rr.addResourceShow(name, base, controller, options)
		case "edit":
			route = rr.addResourceEdit(name, base, controller, options)
		case "update":
			route = rr.addResourceUpdate(name, base, controller, options)
		case "destroy":
			route = rr.addResourceDestroy(name, base, controller, options)
		}

		if route != nil {
			collection.Add(route)
		}
	}

	return collection
}

// Singleton emits the routes for a singleton resource (no index, no list).
func (rr *ResourceRegistrar) Singleton(name, controller string, options map[string]any) *RouteCollection {
	defaults := append([]string(nil), singletonResourceDefaults...)

	if _, ok := options["creatable"]; ok {
		defaults = append(defaults, "create", "store", "destroy")
	} else if _, ok := options["destroyable"]; ok {
		defaults = append(defaults, "destroy")
	}

	methods := rr.getResourceMethods(defaults, options)
	collection := NewRouteCollection()
	base := rr.GetResourceWildcard(lastSegment(name, "."))

	for _, m := range methods {
		var route *Route

		switch m {
		case "create":
			route = rr.addResourceCreate(name, base, controller, options)
		case "store":
			route = rr.addResourceStore(name, base, controller, options)
		case "show":
			route = rr.addSingletonShow(name, controller, options)
		case "edit":
			route = rr.addSingletonEdit(name, controller, options)
		case "update":
			route = rr.addSingletonUpdate(name, controller, options)
		case "destroy":
			route = rr.addSingletonDestroy(name, controller, options)
		}

		if route != nil {
			collection.Add(route)
		}
	}

	return collection
}

// =====================================================================
// Standard resource methods
// =====================================================================

func (rr *ResourceRegistrar) addResourceIndex(name, base, controller string, options map[string]any) *Route {
	uri := rr.GetResourceUri(name)
	action := rr.getResourceAction(name, controller, "index", options)

	return rr.router.AddRoute([]string{"GET", "HEAD"}, uri, action)
}

func (rr *ResourceRegistrar) addResourceCreate(name, base, controller string, options map[string]any) *Route {
	uri := rr.GetResourceUri(name) + "/" + resourceVerbs["create"]
	action := rr.getResourceAction(name, controller, "create", options)

	return rr.router.AddRoute([]string{"GET", "HEAD"}, uri, action)
}

func (rr *ResourceRegistrar) addResourceStore(name, base, controller string, options map[string]any) *Route {
	uri := rr.GetResourceUri(name)
	action := rr.getResourceAction(name, controller, "store", options)

	return rr.router.AddRoute([]string{"POST"}, uri, action)
}

func (rr *ResourceRegistrar) addResourceShow(name, base, controller string, options map[string]any) *Route {
	uri := rr.GetResourceUri(name) + "/{" + base + "}"
	action := rr.getResourceAction(name, controller, "show", options)

	return rr.router.AddRoute([]string{"GET", "HEAD"}, uri, action)
}

func (rr *ResourceRegistrar) addResourceEdit(name, base, controller string, options map[string]any) *Route {
	uri := rr.GetResourceUri(name) + "/{" + base + "}/" + resourceVerbs["edit"]
	action := rr.getResourceAction(name, controller, "edit", options)

	return rr.router.AddRoute([]string{"GET", "HEAD"}, uri, action)
}

func (rr *ResourceRegistrar) addResourceUpdate(name, base, controller string, options map[string]any) *Route {
	uri := rr.GetResourceUri(name) + "/{" + base + "}"
	action := rr.getResourceAction(name, controller, "update", options)

	return rr.router.AddRoute([]string{"PUT", "PATCH"}, uri, action)
}

func (rr *ResourceRegistrar) addResourceDestroy(name, base, controller string, options map[string]any) *Route {
	uri := rr.GetResourceUri(name) + "/{" + base + "}"
	action := rr.getResourceAction(name, controller, "destroy", options)

	return rr.router.AddRoute([]string{"DELETE"}, uri, action)
}

// =====================================================================
// Singleton methods
// =====================================================================

func (rr *ResourceRegistrar) addSingletonShow(name, controller string, options map[string]any) *Route {
	uri := rr.GetResourceUri(name)
	action := rr.getResourceAction(name, controller, "show", options)

	return rr.router.AddRoute([]string{"GET", "HEAD"}, uri, action)
}

func (rr *ResourceRegistrar) addSingletonEdit(name, controller string, options map[string]any) *Route {
	uri := rr.GetResourceUri(name) + "/" + resourceVerbs["edit"]
	action := rr.getResourceAction(name, controller, "edit", options)

	return rr.router.AddRoute([]string{"GET", "HEAD"}, uri, action)
}

func (rr *ResourceRegistrar) addSingletonUpdate(name, controller string, options map[string]any) *Route {
	uri := rr.GetResourceUri(name)
	action := rr.getResourceAction(name, controller, "update", options)

	return rr.router.AddRoute([]string{"PUT", "PATCH"}, uri, action)
}

func (rr *ResourceRegistrar) addSingletonDestroy(name, controller string, options map[string]any) *Route {
	uri := rr.GetResourceUri(name)
	action := rr.getResourceAction(name, controller, "destroy", options)

	return rr.router.AddRoute([]string{"DELETE"}, uri, action)
}

// =====================================================================
// URI / parameter / action helpers
// =====================================================================

// GetResourceUri returns the canonical URI fragment for a (possibly nested)
// resource name like "users" or "users.posts.comments".
func (rr *ResourceRegistrar) GetResourceUri(resource string) string {
	if !strings.Contains(resource, ".") {
		return resource
	}

	segments := strings.Split(resource, ".")

	return rr.getNestedResourceUri(segments)
}

func (rr *ResourceRegistrar) getNestedResourceUri(segments []string) string {
	parts := make([]string, 0, len(segments)*2)

	for i, s := range segments {
		parts = append(parts, s)

		if i < len(segments)-1 {
			parts = append(parts, "{"+rr.GetResourceWildcard(s)+"}")
		}
	}

	return strings.Join(parts, "/")
}

// GetResourceWildcard returns the URL placeholder name for a resource segment.
// "users" → "user", "blog_posts" → "blog_post". Custom mappings supplied via
// the parameters option (or the global parameter map) override this.
func (rr *ResourceRegistrar) GetResourceWildcard(value string) string {
	if v, ok := rr.parameters[value]; ok {
		return v
	}

	if v, ok := globalResourceParameters[value]; ok {
		return v
	}

	return strings.TrimSuffix(value, "s")
}

var globalResourceParameters = map[string]string{}

// SingularParameters toggles automatic singularization (kept for parity).
func SingularParameters(_ bool) { /* no-op in Go port */ }

// SetParameters replaces the global parameter map.
func SetParameters(params map[string]string) {
	globalResourceParameters = map[string]string{}

	for k, v := range params {
		globalResourceParameters[k] = v
	}
}

// GetParameters returns the global parameter map.
func GetParameters() map[string]string { return globalResourceParameters }

// SetResourceVerbs replaces the create/edit verb map.
func SetResourceVerbs(verbs map[string]string) {
	for k, v := range verbs {
		resourceVerbs[k] = v
	}
}

func (rr *ResourceRegistrar) getResourceMethods(defaults []string, options map[string]any) []string {
	// Defensive copy: the package-level `resourceDefaults` slice would be
	// clobbered if we filtered in place, leaking state between tests and
	// between resources registered in the same process.
	methods := append([]string(nil), defaults...)

	if only, ok := options["only"].([]string); ok {
		filtered := make([]string, 0, len(methods))

		for _, m := range methods {
			for _, o := range only {
				if m == o {
					filtered = append(filtered, m)

					break
				}
			}
		}

		methods = filtered
	}

	if except, ok := options["except"].([]string); ok {
		out := make([]string, 0, len(methods))

		for _, m := range methods {
			drop := false

			for _, e := range except {
				if m == e {
					drop = true

					break
				}
			}

			if !drop {
				out = append(out, m)
			}
		}

		methods = out
	}

	return methods
}

func (rr *ResourceRegistrar) getResourceAction(name, controller, method string, options map[string]any) any {
	action := map[string]any{
		"uses":       controller + "@" + method,
		"controller": controller + "@" + method,
		"as":         rr.getResourceRouteName(name, method, options),
	}

	if mw, ok := options["middleware"]; ok {
		action["middleware"] = mw
	}

	return action
}

func (rr *ResourceRegistrar) getResourceRouteName(resource, method string, options map[string]any) string {
	name := resource + "." + method

	if names, ok := options["names"].(map[string]string); ok {
		if v, ok := names[method]; ok {
			name = v
		}
	}

	if as, ok := options["as"].(string); ok {
		name = as + "." + name
	}

	return name
}

func (rr *ResourceRegistrar) prefixedResource(name, controller string, options map[string]any) *RouteCollection {
	idx := strings.LastIndex(name, "/")
	prefix := name[:idx]
	bare := name[idx+1:]
	collection := NewRouteCollection()
	rr.router.Group(map[string]any{"prefix": prefix}, func(r *Router) {
		nested := rr.Register(bare, controller, options)

		for _, route := range nested.GetRoutes() {
			collection.Add(route)
		}
	})

	return collection
}

func (rr *ResourceRegistrar) prefixedSingleton(name, controller string, options map[string]any) *RouteCollection {
	idx := strings.LastIndex(name, "/")
	prefix := name[:idx]
	bare := name[idx+1:]
	collection := NewRouteCollection()
	rr.router.Group(map[string]any{"prefix": prefix}, func(r *Router) {
		nested := rr.Singleton(bare, controller, options)

		for _, route := range nested.GetRoutes() {
			collection.Add(route)
		}
	})

	return collection
}

func lastSegment(s, sep string) string {
	if idx := strings.LastIndex(s, sep); idx >= 0 {
		return s[idx+1:]
	}

	return s
}
