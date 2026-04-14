package routing

import (
	"fmt"
	"strings"

	"github.com/bedrock/packages/routing/matching"
)

// EventDispatcher is the minimal event surface the router needs.
//
// Mirrors Illuminate\Contracts\Events\Dispatcher narrowly. The bedrock
// packages/events Dispatcher satisfies this; tests can supply a noop.
type EventDispatcher interface {
	Dispatch(event any)
}

// noopEvents is the zero-value dispatcher used when none is supplied.
type noopEvents struct{}

// Router is the central router, mirroring Illuminate\Routing\Router.
//
// It owns the [RouteCollection], the middleware aliases and groups, the
// global parameter patterns, the binder map, and the group attribute stack.
// Most public methods correspond 1:1 to a PHP method of the same name; where
// PHP overloads (e.g. domain($d=null) acts as both setter and getter) the Go
// surface splits them into two methods.
type Router struct {
	events                  EventDispatcher
	container               BindingContainer
	routes                  RouteCollectionInterface
	current                 *Route
	currentRequest          matching.MatchableRequest
	middleware              map[string]any   // alias name → class string or closure
	middlewareGroups        map[string][]any // group name → middleware list
	MiddlewarePriority      []string
	binders                 map[string]BindingResolver
	patterns                map[string]string
	groupStack              []map[string]any
	implicitBindingCallback any
}

func (noopEvents) Dispatch(any) {}

// NewRouter constructs a router with the given event dispatcher and
// container. Either or both may be nil — defaults are substituted.
func NewRouter(events EventDispatcher, container BindingContainer) *Router {
	if events == nil {
		events = noopEvents{}
	}

	return &Router{
		events:           events,
		container:        container,
		routes:           NewRouteCollection(),
		middleware:       map[string]any{},
		middlewareGroups: map[string][]any{},
		binders:          map[string]BindingResolver{},
		patterns:         map[string]string{},
	}
}

// =====================================================================
// Verb shortcuts
// =====================================================================

// Get registers a GET (and HEAD) route.
func (r *Router) Get(uri string, action any) *Route {
	return r.AddRoute([]string{"GET", "HEAD"}, uri, action)
}

// Post registers a POST route.
func (r *Router) Post(uri string, action any) *Route {
	return r.AddRoute([]string{"POST"}, uri, action)
}

// Put registers a PUT route.
func (r *Router) Put(uri string, action any) *Route {
	return r.AddRoute([]string{"PUT"}, uri, action)
}

// Patch registers a PATCH route.
func (r *Router) Patch(uri string, action any) *Route {
	return r.AddRoute([]string{"PATCH"}, uri, action)
}

// Delete registers a DELETE route.
func (r *Router) Delete(uri string, action any) *Route {
	return r.AddRoute([]string{"DELETE"}, uri, action)
}

// Options registers an OPTIONS route.
func (r *Router) Options(uri string, action any) *Route {
	return r.AddRoute([]string{"OPTIONS"}, uri, action)
}

// Any registers a route that matches all of the verbs in [HTTPVerbs].
func (r *Router) Any(uri string, action any) *Route {
	verbs := append([]string(nil), HTTPVerbs...)

	return r.AddRoute(verbs, uri, action)
}

// Match registers a route bound to the supplied verbs.
func (r *Router) Match(methods []string, uri string, action any) *Route {
	upper := make([]string, len(methods))

	for i, m := range methods {
		upper[i] = strings.ToUpper(m)
	}

	return r.AddRoute(upper, uri, action)
}

// Fallback registers a fallback route that runs when no other route matches.
func (r *Router) Fallback(action any) *Route {
	route := r.AddRoute([]string{"GET", "HEAD"}, "{fallbackPlaceholder}", action)
	route.Where("fallbackPlaceholder", `.*`)
	route.Fallback()

	return route
}

// Redirect registers a route that issues an HTTP redirect.
func (r *Router) Redirect(uri, destination string, status int) *Route {
	if status == 0 {
		status = 302
	}

	return r.Any(uri, map[string]any{
		"uses": func() (string, int) { return destination, status },
	})
}

// PermanentRedirect registers a 301 redirect.
func (r *Router) PermanentRedirect(uri, destination string) *Route {
	return r.Redirect(uri, destination, 301)
}

// =====================================================================
// Resource entry points (mirror Router::resource / ::apiResource etc.)
// =====================================================================

// Resource starts a fluent registration for a RESTful resource controller.
func (r *Router) Resource(name, controller string, options map[string]any) *PendingResourceRegistration {
	return NewPendingResourceRegistration(NewResourceRegistrar(r), name, controller, options)
}

// ApiResource starts a fluent registration that excludes the create and edit
// form actions (the JSON-API subset).
func (r *Router) ApiResource(name, controller string, options map[string]any) *PendingResourceRegistration {
	if options == nil {
		options = map[string]any{}
	}

	if _, ok := options["only"]; !ok {
		options["except"] = []string{"create", "edit"}
	}

	return r.Resource(name, controller, options)
}

// Singleton starts a fluent registration for a singleton resource controller.
func (r *Router) Singleton(name, controller string, options map[string]any) *PendingSingletonResourceRegistration {
	return NewPendingSingletonResourceRegistration(NewResourceRegistrar(r), name, controller, options)
}

// ApiSingleton is the JSON-API counterpart of [Router.Singleton].
func (r *Router) ApiSingleton(name, controller string, options map[string]any) *PendingSingletonResourceRegistration {
	if options == nil {
		options = map[string]any{}
	}

	if _, ok := options["only"]; !ok {
		options["except"] = []string{"edit"}
	}

	return r.Singleton(name, controller, options)
}

// Resources iterates name → controller pairs, registering each as a resource.
func (r *Router) Resources(resources map[string]string, options map[string]any) {
	for name, controller := range resources {
		r.Resource(name, controller, options).Register()
	}
}

// ApiResources is the bulk counterpart of [Router.ApiResource].
func (r *Router) ApiResources(resources map[string]string, options map[string]any) {
	for name, controller := range resources {
		r.ApiResource(name, controller, options).Register()
	}
}

// Singletons is the bulk counterpart of [Router.Singleton].
func (r *Router) Singletons(resources map[string]string, options map[string]any) {
	for name, controller := range resources {
		r.Singleton(name, controller, options).Register()
	}
}

// =====================================================================
// Group machinery
// =====================================================================

// Group registers routes inside a shared attribute scope.
//
// Mirrors Router::group. The routes argument is a Go function so registration
// happens within a clearly-scoped block; PHP's "string filename" form is
// supported by passing the result of [NewRouteFileRegistrar].Register to a
// closure.
func (r *Router) Group(attributes map[string]any, routes func(*Router)) *Router {
	r.updateGroupStack(attributes)
	routes(r)
	r.groupStack = r.groupStack[:len(r.groupStack)-1]

	return r
}

func (r *Router) updateGroupStack(attributes map[string]any) {
	if r.HasGroupStack() {
		attributes = r.MergeWithLastGroup(attributes, true)
	}

	r.groupStack = append(r.groupStack, attributes)
}

// MergeWithLastGroup merges the given attributes with the top of the group
// stack. Mirrors Router::mergeWithLastGroup.
func (r *Router) MergeWithLastGroup(new map[string]any, prependExistingPrefix bool) map[string]any {
	if !r.HasGroupStack() {
		return new
	}

	return MergeRouteGroup(new, r.groupStack[len(r.groupStack)-1], prependExistingPrefix)
}

// GetLastGroupPrefix returns the prefix of the top group in the stack.
func (r *Router) GetLastGroupPrefix() string {
	if !r.HasGroupStack() {
		return ""
	}

	last := r.groupStack[len(r.groupStack)-1]

	if v, ok := last["prefix"]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}

	return ""
}

// HasGroupStack reports whether any groups are currently active.
func (r *Router) HasGroupStack() bool { return len(r.groupStack) > 0 }

// GetGroupStack returns the current group stack.
func (r *Router) GetGroupStack() []map[string]any { return r.groupStack }

// =====================================================================
// Route creation
// =====================================================================

// AddRoute creates and registers a new route.
func (r *Router) AddRoute(methods []string, uri string, action any) *Route {
	return r.routes.Add(r.createRoute(methods, uri, action))
}

func (r *Router) createRoute(methods []string, uri string, action any) *Route {
	if r.actionReferencesController(action) {
		action = r.convertToControllerAction(action)
	}

	route := r.NewRoute(methods, r.prefixUri(uri), action)

	if r.HasGroupStack() {
		r.mergeGroupAttributesIntoRoute(route)
	}

	r.addWhereClausesToRoute(route)

	return route
}

func (r *Router) actionReferencesController(action any) bool {
	switch v := action.(type) {
	case nil:
		return false
	case string:
		return true
	case map[string]any:
		if uses, ok := v["uses"]; ok {
			_, isString := uses.(string)

			return isString
		}
	}

	return false
}

func (r *Router) convertToControllerAction(action any) any {
	var m map[string]any

	switch v := action.(type) {
	case string:
		m = map[string]any{"uses": v}
	case map[string]any:
		m = map[string]any{}

		for k, val := range v {
			m[k] = val
		}
	}

	if r.HasGroupStack() {
		if uses, ok := m["uses"].(string); ok {
			uses = r.prependGroupController(uses)
			uses = r.prependGroupNamespace(uses)
			m["uses"] = uses
		}
	}

	if uses, ok := m["uses"].(string); ok {
		m["controller"] = uses
	}

	return m
}

func (r *Router) prependGroupNamespace(class string) string {
	if !r.HasGroupStack() {
		return class
	}

	last := r.groupStack[len(r.groupStack)-1]
	ns, ok := last["namespace"].(string)

	if !ok || ns == "" {
		return class
	}

	if strings.HasPrefix(class, `\`) || strings.HasPrefix(class, ns) {
		return class
	}

	return ns + `\` + class
}

func (r *Router) prependGroupController(class string) string {
	if !r.HasGroupStack() {
		return class
	}

	last := r.groupStack[len(r.groupStack)-1]
	groupCtrl, ok := last["controller"].(string)

	if !ok || groupCtrl == "" {
		return class
	}

	if strings.Contains(class, "@") {
		return class
	}

	return groupCtrl + "@" + class
}

// NewRoute is the parity-named factory used by [Router.AddRoute] and by
// callers that want a wired-but-unregistered route.
func (r *Router) NewRoute(methods []string, uri string, action any) *Route {
	return NewRoute(methods, uri, action).
		SetRouter(r).
		SetContainer(r.container)
}

func (r *Router) prefixUri(uri string) string {
	prefix := strings.Trim(r.GetLastGroupPrefix(), "/")
	combined := strings.Trim(prefix+"/"+strings.Trim(uri, "/"), "/")

	if combined == "" {
		return "/"
	}

	return combined
}

func (r *Router) addWhereClausesToRoute(route *Route) {
	for k, v := range r.patterns {
		if _, ok := route.Wheres[k]; !ok {
			route.Wheres[k] = v
		}
	}

	if where, ok := route.ActionMap["where"]; ok {
		if m, ok := where.(map[string]string); ok {
			for k, v := range m {
				route.Wheres[k] = v
			}
		}
	}
}

func (r *Router) mergeGroupAttributesIntoRoute(route *Route) {
	merged := r.MergeWithLastGroup(route.ActionMap, false)
	route.SetAction(merged)
}

// =====================================================================
// Dispatch
// =====================================================================

// Dispatch routes an incoming request to a Route and returns whatever the
// handler returned. Response normalization is M8's responsibility.
func (r *Router) Dispatch(request matching.MatchableRequest) (any, error) {
	r.currentRequest = request

	return r.DispatchToRoute(request)
}

// DispatchToRoute matches the request against the routes and runs the result.
func (r *Router) DispatchToRoute(request matching.MatchableRequest) (any, error) {
	route, err := r.findRoute(request)

	if err != nil {
		return nil, err
	}

	return r.runRoute(request, route)
}

func (r *Router) findRoute(request matching.MatchableRequest) (*Route, error) {
	route, err := r.routes.Match(request)

	if err != nil {
		return nil, err
	}

	r.current = route
	route.SetContainer(r.container)

	return route, nil
}

// runRoute is the simplified Go form: middleware-less direct invocation. Full
// pipeline support arrives in M5 once ControllerDispatcher/CallableDispatcher
// land. For now we call the route's "uses" handler if it's a Go function.
func (r *Router) runRoute(request matching.MatchableRequest, route *Route) (any, error) {
	uses, ok := route.ActionMap["uses"]

	if !ok || uses == nil {
		return nil, fmt.Errorf("route %s has no handler", route.Uri)
	}

	switch fn := uses.(type) {
	case func():
		fn()

		return nil, nil
	case func() error:
		return nil, fn()
	case func() any:
		return fn(), nil
	case func() (any, error):
		return fn()
	}

	return uses, nil
}

// =====================================================================
// Bindings / patterns / middleware aliases
// =====================================================================

// Bind registers an explicit binder for the named parameter.
func (r *Router) Bind(key string, binder any) {
	key = strings.ReplaceAll(key, "-", "_")
	r.binders[key] = ForCallback(r.container, binder)
}

// Model registers a model binder for the named parameter.
func (r *Router) Model(key, class string, callback BindingResolver) {
	r.Bind(key, ForModel(r.container, class, callback))
}

// GetBindingCallback returns the binder registered for key, or nil.
func (r *Router) GetBindingCallback(key string) BindingResolver {
	key = strings.ReplaceAll(key, "-", "_")

	return r.binders[key]
}

// GetPatterns returns the global where-pattern map.
func (r *Router) GetPatterns() map[string]string { return r.patterns }

// Pattern sets a global where-pattern that will be applied to every newly
// registered route.
func (r *Router) Pattern(key, pattern string) { r.patterns[key] = pattern }

// Patterns sets multiple global where-patterns at once.
func (r *Router) Patterns(patterns map[string]string) {
	for k, v := range patterns {
		r.Pattern(k, v)
	}
}

// AliasMiddleware registers a short-hand name for a middleware class.
func (r *Router) AliasMiddleware(name string, class any) *Router {
	r.middleware[name] = class

	return r
}

// GetMiddleware returns the alias map.
func (r *Router) GetMiddleware() map[string]any { return r.middleware }

// HasMiddlewareGroup reports whether a middleware group is registered.
func (r *Router) HasMiddlewareGroup(name string) bool {
	_, ok := r.middlewareGroups[name]

	return ok
}

// GetMiddlewareGroups returns the group map.
func (r *Router) GetMiddlewareGroups() map[string][]any { return r.middlewareGroups }

// MiddlewareGroup registers a named group of middleware.
func (r *Router) MiddlewareGroup(name string, middleware []any) *Router {
	r.middlewareGroups[name] = middleware

	return r
}

// PrependMiddlewareToGroup adds a middleware to the front of an existing
// group, deduplicating against the existing entries.
func (r *Router) PrependMiddlewareToGroup(group string, middleware any) *Router {
	if existing, ok := r.middlewareGroups[group]; ok {
		for _, m := range existing {
			if m == middleware {
				return r
			}
		}

		r.middlewareGroups[group] = append([]any{middleware}, existing...)
	}

	return r
}

// PushMiddlewareToGroup adds a middleware to the end of a group, creating the
// group if it does not yet exist.
func (r *Router) PushMiddlewareToGroup(group string, middleware any) *Router {
	if _, ok := r.middlewareGroups[group]; !ok {
		r.middlewareGroups[group] = []any{}
	}

	for _, m := range r.middlewareGroups[group] {
		if m == middleware {
			return r
		}
	}

	r.middlewareGroups[group] = append(r.middlewareGroups[group], middleware)

	return r
}

// RemoveMiddlewareFromGroup removes a middleware from the named group.
func (r *Router) RemoveMiddlewareFromGroup(group string, middleware any) *Router {
	if !r.HasMiddlewareGroup(group) {
		return r
	}

	out := r.middlewareGroups[group][:0]

	for _, m := range r.middlewareGroups[group] {
		if m != middleware {
			out = append(out, m)
		}
	}

	r.middlewareGroups[group] = out

	return r
}

// FlushMiddlewareGroups clears all middleware groups.
func (r *Router) FlushMiddlewareGroups() *Router {
	r.middlewareGroups = map[string][]any{}

	return r
}

// =====================================================================
// Current route accessors
// =====================================================================

// Current returns the most recently dispatched route, or nil.
func (r *Router) Current() *Route { return r.current }

// GetCurrentRoute is a parity alias for [Router.Current].
func (r *Router) GetCurrentRoute() *Route { return r.current }

// GetCurrentRequest returns the most recently dispatched request.
func (r *Router) GetCurrentRequest() matching.MatchableRequest { return r.currentRequest }

// CurrentRouteName returns the name of the current route, or "".
func (r *Router) CurrentRouteName() string {
	if r.current == nil {
		return ""
	}

	return r.current.GetName()
}

// CurrentRouteAction returns the controller action of the current route.
func (r *Router) CurrentRouteAction() string {
	if r.current == nil {
		return ""
	}

	return r.current.GetActionName()
}

// Has reports whether every supplied name is registered.
func (r *Router) Has(names ...string) bool {
	for _, n := range names {
		if !r.routes.HasNamedRoute(n) {
			return false
		}
	}

	return true
}

// Is reports whether the current route name matches any of the given glob
// patterns.
func (r *Router) Is(patterns ...string) bool {
	if r.current == nil {
		return false
	}

	return r.current.Named(patterns...)
}

// CurrentRouteNamed is an alias for [Router.Is] kept for parity.
func (r *Router) CurrentRouteNamed(patterns ...string) bool { return r.Is(patterns...) }

// Uses reports whether the current route's action matches any of the given
// glob patterns. Useful in middleware that wants to scope behavior to a
// controller namespace.
func (r *Router) Uses(patterns ...string) bool {
	if r.current == nil {
		return false
	}

	action := r.current.GetActionName()

	for _, p := range patterns {
		if matchPattern(p, action) {
			return true
		}
	}

	return false
}

// CurrentRouteUses reports whether the current route's action equals action.
func (r *Router) CurrentRouteUses(action string) bool {
	return r.current != nil && r.current.GetActionName() == action
}

// =====================================================================
// Routes accessor / replace
// =====================================================================

// GetRoutes returns the underlying route collection.
func (r *Router) GetRoutes() RouteCollectionInterface { return r.routes }

// SetRoutes swaps the underlying route collection. After a swap each route's
// container reference is refreshed so future dispatches resolve through the
// router's container.
func (r *Router) SetRoutes(routes RouteCollectionInterface) {
	for _, route := range routes.GetRoutes() {
		route.SetContainer(r.container)
	}

	r.routes = routes
}

// SetCompiledRoutes installs a [CompiledRouteCollection] built from the
// supplied routes — the production-mode entry point used by the route cache.
func (r *Router) SetCompiledRoutes(routes []*Route) {
	r.SetRoutes(NewCompiledRouteCollection(routes, nil))
}

// SetContainer rebinds the container instance.
func (r *Router) SetContainer(c BindingContainer) { r.container = c }

// =====================================================================
// Bindings
// =====================================================================

// SubstituteBindings runs the explicit binders on the route's parameters,
// replacing each value with whatever the binder returned.
//
// Mirrors Router::substituteBindings.
func (r *Router) SubstituteBindings(route *Route) error {
	if route.Parameters == nil {
		return nil
	}

	for k, v := range route.Parameters {
		if binder, ok := r.binders[k]; ok {
			result, err := binder(v, route)

			if err != nil {
				return err
			}

			if s, ok := result.(string); ok {
				route.SetParameter(k, s)
			}
		}
	}

	return nil
}

// SubstituteImplicitBindings delegates to [ImplicitRouteBinding.ResolveForRoute].
//
// dependencyResolver is the type-keyed lookup the resolver uses to
// instantiate model parameters; M11 wires in the bedrock container adapter.
func (r *Router) SubstituteImplicitBindings(route *Route, dependencyResolver DependencyContainer) error {
	if dependencyResolver == nil {
		return nil
	}

	return ImplicitRouteBinding{}.ResolveForRoute(dependencyResolver, route)
}
