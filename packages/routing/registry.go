package routing

import (
	"encoding/json"
	"log"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

// Route describes a single named route entry in the registry.
type Route struct {
	Name    string `json:"name"`
	Method  string `json:"method"`
	Pattern string `json:"pattern"`
}

// Params extracts the path parameter names from the route pattern.

// Registry holds named routes and provides URL resolution and export
// capabilities. It is safe for concurrent use.
type Registry struct {
	mu     sync.RWMutex
	routes map[string]Route
	order  []string
}

// NewRegistry creates an empty route registry.

// Add registers a named route. method may be empty for path-only entries.
// The pattern uses {param} placeholders (e.g., "/users/{id}").

// Group calls fn with a scoped builder that prefixes both the route name
// (dot-separated) and the URL pattern.

// URL resolves a named route substituting {param} placeholders with params.
// Unknown names return a fallback string and log a warning.

// Lookup returns the Route for the given name, or ok=false.

// Manifest returns a name-to-pattern map suitable for sharing as props.

// ManifestProps returns the manifest as map[string]any.

// Export returns all registered routes in insertion order.

// ToJSON serializes all routes as a JSON array.

// RegistryGroup is a scoped builder that prefixes route names and patterns.
type RegistryGroup struct {
	registry   *Registry
	namePrefix string
	pathPrefix string
}

var paramRegex = regexp.MustCompile(`\{(\w+)\}`)

func (r Route) Params() []string {
	matches := paramRegex.FindAllStringSubmatch(r.Pattern, -1)
	params := make([]string, 0, len(matches))

	for _, m := range matches {
		params = append(params, m[1])
	}

	return params
}

func NewRegistry() *Registry {
	return &Registry{routes: make(map[string]Route)}
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

func (r *Registry) Group(namePrefix, pathPrefix string, fn func(*RegistryGroup)) *Registry {
	fn(&RegistryGroup{registry: r, namePrefix: namePrefix, pathPrefix: pathPrefix})

	return r
}

func (r *Registry) URL(name string, params map[string]string) string {
	r.mu.RLock()
	route, ok := r.routes[name]
	r.mu.RUnlock()

	if !ok {
		log.Printf("routing: unknown route %q, returning fallback", name)

		return "#!routing:unknown-route"
	}

	result := route.Pattern

	for k, v := range params {
		result = strings.ReplaceAll(result, "{"+k+"}", url.PathEscape(v))
	}

	return result
}

func (r *Registry) Lookup(name string) (Route, bool) {
	r.mu.RLock()

	defer r.mu.RUnlock()

	route, ok := r.routes[name]

	return route, ok
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

func (r *Registry) ToJSON() ([]byte, error) {
	return json.Marshal(r.Export())
}

// Add registers a route within the group's scope.
func (g *RegistryGroup) Add(name, method, pattern string) *RegistryGroup {
	g.registry.Add(g.namePrefix+"."+name, method, g.pathPrefix+pattern)

	return g
}

// Group creates a nested sub-group.
func (g *RegistryGroup) Group(namePrefix, pathPrefix string, fn func(*RegistryGroup)) *RegistryGroup {
	sub := &RegistryGroup{
		registry:   g.registry,
		namePrefix: g.namePrefix + "." + namePrefix,
		pathPrefix: g.pathPrefix + pathPrefix,
	}

	fn(sub)

	return g
}
