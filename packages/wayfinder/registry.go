package wayfinder

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// Registry holds named routes and provides URL resolution and export
// capabilities. It is safe for concurrent use.
type Registry struct {
	mu     sync.RWMutex
	routes map[string]Route
	order  []string
}

// New creates an empty route registry.

// Add registers a named route. The pattern uses {param} placeholders
// (e.g., "/contacts/{contact}"). Method is the HTTP method ("GET",
// "POST", etc.).

// Group calls fn with a scoped builder that prefixes both the route
// name (dot-separated) and the URL pattern.

// URL resolves a named route with the given parameters. Unknown route
// names return a non-navigating fallback and log a warning. Parameters replace {param} placeholders.

// Manifest returns a name-to-pattern map suitable for sharing as
// Inertia props.

// ManifestProps returns the manifest as map[string]any, ready for
// use with inertia.ShareProps or inertia.SetProps.

// Export returns all registered routes in insertion order.

// Lookup returns the Route for the given name, or ok=false.

// ToJSON serializes all routes as a JSON array.

// Group is a scoped builder that prefixes route names and patterns.
type Group struct {
	registry   *Registry
	namePrefix string
	pathPrefix string
}

func New() *Registry {
	return &Registry{
		routes: make(map[string]Route),
	}
}

func (r *Registry) Add(name, method, pattern string) *Registry {
	r.mu.Lock()

	defer r.mu.Unlock()

	route := Route{
		Name:    name,
		Method:  strings.ToUpper(method),
		Pattern: pattern,
	}

	if _, exists := r.routes[name]; !exists {
		r.order = append(r.order, name)
	}

	r.routes[name] = route

	return r
}

func (r *Registry) Group(namePrefix, pathPrefix string, fn func(g *Group)) *Registry {
	g := &Group{
		registry:   r,
		namePrefix: namePrefix,
		pathPrefix: pathPrefix,
	}

	fn(g)

	return r
}

func (r *Registry) URL(name string, params map[string]string) string {
	r.mu.RLock()

	route, ok := r.routes[name]

	r.mu.RUnlock()

	if !ok {
		log.Printf("wayfinder: unknown route %q, returning fallback", name)

		return "#!wayfinder:unknown-route"
	}

	result := route.Pattern

	for key, value := range params {
		result = strings.ReplaceAll(result, "{"+key+"}", url.PathEscape(value))
	}

	return result
}

// Handle registers an HTTP handler on the given mux using the pattern
// stored for the named route. Unknown names log a warning and skip.
func (r *Registry) Handle(name string, handler http.Handler, mux *http.ServeMux) {
	r.mu.RLock()

	route, ok := r.routes[name]

	r.mu.RUnlock()

	if !ok {
		log.Printf("wayfinder: Handle: unknown route %q, skipping", name)

		return
	}

	mux.Handle(route.Pattern, handler)
}

func (r *Registry) Manifest() map[string]string {
	r.mu.RLock()

	defer r.mu.RUnlock()

	m := make(map[string]string, len(r.routes))

	for name, route := range r.routes {
		m[name] = route.Pattern
	}

	return m
}

func (r *Registry) ManifestProps() map[string]any {
	manifest := r.Manifest()
	props := make(map[string]any, len(manifest))

	for k, v := range manifest {
		props[k] = v
	}

	return props
}

func (r *Registry) Export() []Route {
	r.mu.RLock()

	defer r.mu.RUnlock()

	routes := make([]Route, 0, len(r.order))

	for _, name := range r.order {
		routes = append(routes, r.routes[name])
	}

	return routes
}

func (r *Registry) Lookup(name string) (Route, bool) {
	r.mu.RLock()

	defer r.mu.RUnlock()

	route, ok := r.routes[name]

	return route, ok
}

func (r *Registry) ToJSON() ([]byte, error) {
	return json.Marshal(r.Export())
}

// Add registers a route within the group's scope.
func (g *Group) Add(name, method, pattern string) *Group {
	fullName := g.namePrefix + "." + name
	fullPattern := g.pathPrefix + pattern

	g.registry.Add(fullName, method, fullPattern)

	return g
}

// Handle registers an HTTP handler via the parent registry using the
// group-scoped route name.
func (g *Group) Handle(name string, handler http.Handler, mux *http.ServeMux) {
	g.registry.Handle(g.namePrefix+"."+name, handler, mux)
}

// Group creates a nested sub-group.
func (g *Group) Group(namePrefix, pathPrefix string, fn func(g *Group)) *Group {
	sub := &Group{
		registry:   g.registry,
		namePrefix: g.namePrefix + "." + namePrefix,
		pathPrefix: g.pathPrefix + pathPrefix,
	}

	fn(sub)

	return g
}
