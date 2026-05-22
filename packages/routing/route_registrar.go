package routing

// RouteRegistrar is the fluent builder returned when the user starts a route
// registration chain with attribute-only calls (e.g. r.Middleware("auth").Get(...)).
//
// Each builder method clones the registrar with the new attribute applied and
// returns it; the terminal verb call (Get/Post/etc.) drains the accumulated
// attributes into the created route via the router's group machinery.
//
// Ref: @bedrock/code-0339
type RouteRegistrar struct {
	router     *Router
	attributes map[string]any
	CreatesRegularExpressionRouteConstraints
}

// NewRouteRegistrar constructs an empty registrar wired to router.

// routeRegistrarTarget adapts RouteRegistrar to [whereTarget]. WhereAlpha
// helpers funnel through here so they can populate the registrar's where map.
type routeRegistrarTarget struct{ r *RouteRegistrar }

func NewRouteRegistrar(router *Router) *RouteRegistrar {
	r := &RouteRegistrar{router: router, attributes: map[string]any{}}
	r.CreatesRegularExpressionRouteConstraints.Bind(routeRegistrarTarget{r})

	return r
}

func (t routeRegistrarTarget) Where(name, expression string) *Route {
	w, _ := t.r.attributes["where"].(map[string]string)

	if w == nil {
		w = map[string]string{}
	}

	w[name] = expression
	t.r.attributes["where"] = w

	return nil
}

// =====================================================================
// Attribute setters
// =====================================================================

// Middleware adds middleware to the registrar's accumulated attributes.
func (r *RouteRegistrar) Middleware(middleware ...any) *RouteRegistrar {
	existing, _ := r.attributes["middleware"].([]any)

	for _, m := range middleware {
		if m == nil {
			continue
		}

		existing = append(existing, m)
	}

	if len(existing) == 0 {
		delete(r.attributes, "middleware")
	} else {
		r.attributes["middleware"] = existing
	}

	return r
}

// WithoutMiddleware accumulates middleware names to be excluded from the
// resulting routes.
func (r *RouteRegistrar) WithoutMiddleware(middleware ...any) *RouteRegistrar {
	existing, _ := r.attributes["excluded_middleware"].([]any)

	for _, m := range middleware {
		if m == nil {
			continue
		}

		existing = append(existing, m)
	}

	if len(existing) == 0 {
		delete(r.attributes, "excluded_middleware")
	} else {
		r.attributes["excluded_middleware"] = existing
	}

	return r
}

// Prefix sets the URI prefix.
func (r *RouteRegistrar) Prefix(prefix string) *RouteRegistrar {
	r.attributes["prefix"] = prefix

	return r
}

// Domain sets the host pattern.
func (r *RouteRegistrar) Domain(domain string) *RouteRegistrar {
	r.attributes["domain"] = domain

	return r
}

// Name (alias: As) prepends to the route name. Use NamePrefix for clarity in
// chained calls.
func (r *RouteRegistrar) Name(name string) *RouteRegistrar { return r.As(name) }

// As sets the name prefix.
func (r *RouteRegistrar) As(value string) *RouteRegistrar {
	r.attributes["as"] = value

	return r
}

// Namespace sets the controller namespace prefix.
func (r *RouteRegistrar) Namespace(ns string) *RouteRegistrar {
	r.attributes["namespace"] = ns

	return r
}

// Controller sets the group-level controller, allowing terminal verb calls to
// pass just a method name as their action.
func (r *RouteRegistrar) Controller(controller string) *RouteRegistrar {
	r.attributes["controller"] = controller

	return r
}

// ScopeBindings enables scoped implicit bindings on the resulting routes.
func (r *RouteRegistrar) ScopeBindings() *RouteRegistrar {
	r.attributes["scope_bindings"] = true

	return r
}

// WithoutScopedBindings disables scoped implicit bindings.
func (r *RouteRegistrar) WithoutScopedBindings() *RouteRegistrar {
	r.attributes["scope_bindings"] = false

	return r
}

// =====================================================================
// Terminal verb calls
// =====================================================================

// Get registers a GET route within the registrar's attribute scope.
func (r *RouteRegistrar) Get(uri string, action any) *Route {
	return r.registerRoute("Get", uri, action)
}

// Post registers a POST route.
func (r *RouteRegistrar) Post(uri string, action any) *Route {
	return r.registerRoute("Post", uri, action)
}

// Put registers a PUT route.
func (r *RouteRegistrar) Put(uri string, action any) *Route {
	return r.registerRoute("Put", uri, action)
}

// Patch registers a PATCH route.
func (r *RouteRegistrar) Patch(uri string, action any) *Route {
	return r.registerRoute("Patch", uri, action)
}

// Delete registers a DELETE route.
func (r *RouteRegistrar) Delete(uri string, action any) *Route {
	return r.registerRoute("Delete", uri, action)
}

// Options registers an OPTIONS route.
func (r *RouteRegistrar) Options(uri string, action any) *Route {
	return r.registerRoute("Options", uri, action)
}

// Any registers a route handling all standard verbs.
func (r *RouteRegistrar) Any(uri string, action any) *Route {
	return r.registerRoute("Any", uri, action)
}

// Group runs the supplied callback inside the registrar's attribute scope.
func (r *RouteRegistrar) Group(routes func(*Router)) {
	r.router.Group(r.attributes, routes)
}

// Resource registers a RESTful resource controller within scope.
func (r *RouteRegistrar) Resource(name, controller string, options map[string]any) *PendingResourceRegistration {
	merged := mergeOptions(r.attributes, options)
	rr := NewResourceRegistrar(r.router)

	return NewPendingResourceRegistration(rr, name, controller, merged)
}

// =====================================================================
// Internals
// =====================================================================

func (r *RouteRegistrar) registerRoute(verb, uri string, action any) *Route {
	var route *Route

	r.router.Group(r.attributes, func(rt *Router) {
		switch verb {
		case "Get":
			route = rt.Get(uri, action)
		case "Post":
			route = rt.Post(uri, action)
		case "Put":
			route = rt.Put(uri, action)
		case "Patch":
			route = rt.Patch(uri, action)
		case "Delete":
			route = rt.Delete(uri, action)
		case "Options":
			route = rt.Options(uri, action)
		case "Any":
			route = rt.Any(uri, action)
		}
	})

	return route
}

func mergeOptions(a, b map[string]any) map[string]any {
	out := map[string]any{}

	for k, v := range a {
		out[k] = v
	}

	for k, v := range b {
		out[k] = v
	}

	return out
}
