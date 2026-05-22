package routegen

import (
	"strings"
)

// RouteInfo contains all route metadata required for TypeScript generation.
// It is intentionally decoupled from the bedrock routing package so the
// generator can be used without a live Router (e.g. from JSON manifests).
type RouteInfo struct {
	// URI is the raw route URI pattern as registered (e.g. "/posts/{post}").
	URI string

	// Methods holds the lowercase HTTP methods accepted by the route
	// (e.g. ["get", "head"]).
	Methods []string

	// Name is the route name (e.g. "posts.show").  Empty for unnamed routes.
	Name string

	// Controller is the fully-qualified "Class@Method" action string
	// (e.g. "App\\Http\\Controllers\\PostController@show").
	// Empty for closure-based routes.
	Controller string

	// IsInvokable is true when the action class and action method refer to
	// the same thing — i.e. the controller has a single __invoke / Invoke
	// method rather than named action methods.
	IsInvokable bool

	// Params is the ordered list of URI parameters.
	Params []Param

	// Domain is the host pattern for domain-scoped routes
	// (e.g. "{tenant}.example.com" or "api.example.com").
	Domain string

	// Scheme is a forced URL scheme ("http://", "https://", or "//").
	// Empty means protocol-relative (//domain/path) when Domain is set.
	Scheme string

	// BasePath is the path prefix derived from APP_URL (e.g. "/v2").
	BasePath string

	// Defaults holds URL default values populated by middleware
	// (equivalent to URL::defaults() calls in PHP middleware).
	Defaults map[string]string
}

// HasController reports whether this route is backed by a named controller
// rather than a closure.
func (r *RouteInfo) HasController() bool {
	return r.Controller != "" && r.Controller != "\\Closure" && r.Controller != "Closure"
}

// ControllerClass returns just the class portion of the Controller string,
// stripping the "@Method" suffix.
//
// "App\\Http\\Controllers\\PostController@show" → "App\\Http\\Controllers\\PostController"
func (r *RouteInfo) ControllerClass() string {
	if idx := strings.LastIndex(r.Controller, "@"); idx >= 0 {
		return r.Controller[:idx]
	}

	return r.Controller
}

// ActionMethod returns the method portion of the Controller string.
//
// "App\\Http\\Controllers\\PostController@show" → "show"
func (r *RouteInfo) ActionMethod() string {
	if idx := strings.LastIndex(r.Controller, "@"); idx >= 0 {
		return r.Controller[idx+1:]
	}

	return r.Controller
}

// DotNamespace converts the controller class path into a dot-separated
// namespace suitable for grouping and file-path construction.
//
// "App\\Http\\Controllers\\PostController@show" → "App.Http.Controllers.PostController"
func (r *RouteInfo) DotNamespace() string {
	cls := r.ControllerClass()
	// Strip leading backslash if present.
	cls = strings.TrimPrefix(cls, "\\")
	// Convert both Go backslash style and forward-slash style.
	cls = strings.ReplaceAll(cls, "\\", ".")
	cls = strings.ReplaceAll(cls, "/", ".")

	return cls
}

// OriginalJsMethod returns the raw method name before TypeScript sanitisation.
// For invokable controllers this is the last segment of the class name.
func (r *RouteInfo) OriginalJsMethod() string {
	if r.IsInvokable {
		cls := r.ControllerClass()
		cls = strings.TrimPrefix(cls, "\\")
		parts := strings.FieldsFunc(cls, func(c rune) bool {
			return c == '\\' || c == '/'
		})

		if len(parts) > 0 {
			return parts[len(parts)-1]
		}

		return cls
	}

	return r.ActionMethod()
}

// JsMethod returns the TypeScript-safe method name for use as an export
// identifier or object property.
func (r *RouteInfo) JsMethod() string {
	return SafeMethod(r.OriginalJsMethod(), "Method")
}

// NamedMethod returns the TypeScript-safe name derived from the last
// dot-segment of the route name.
func (r *RouteInfo) NamedMethod() string {
	name := r.Name

	if idx := strings.LastIndex(name, "."); idx >= 0 {
		name = name[idx+1:]
	}

	return SafeMethod(name, "Method")
}

// Verbs returns the HTTP verbs for this route as Verb values.
func (r *RouteInfo) Verbs() []Verb {
	return verbsFromMethods(r.Methods)
}

// FullURI assembles the full URI including optional scheme, domain, and base
// path prefix. It is the equivalent of Route::uri() in the PHP implementation.
func (r *RouteInfo) FullURI() string {
	// Apply defaults: mark defaulted parameters as optional in the URI.
	uri := r.URI

	if uri == "" {
		return jsonString("")
	}

	if !strings.HasPrefix(uri, "/") {
		uri = "/" + uri
	}

	// Apply base path prefix.
	if r.BasePath != "" {
		base := strings.TrimRight(r.BasePath, "/")
		uri = base + "/" + strings.TrimLeft(uri, "/")
	}

	// Convert defaulted params to optional in the URI string.
	for k := range r.Defaults {
		placeholder := "{" + k + "}"
		optional := "{" + k + "?}"
		uri = strings.ReplaceAll(uri, placeholder, optional)
	}

	// Prepend domain if present.
	if r.Domain != "" {
		scheme := r.Scheme

		if scheme == "" {
			scheme = "//"
		}

		uri = scheme + r.Domain + uri
	}

	// Trim trailing slash (except for bare "/").
	if uri != "/" {
		uri = strings.TrimRight(uri, "/")
	}

	return jsonString(uri)
}

// methodActuals returns the lowercase HTTP method strings for use in
// definition.methods JSON.
func (r *RouteInfo) methodActuals() []string {
	out := make([]string, len(r.Methods))

	for i, m := range r.Methods {
		out[i] = strings.ToLower(m)
	}

	return out
}

// routeNameToFileParts splits a dot-separated route name into file path parts
// for the routes/ directory, appending "index" as the final segment.
//
// "posts.edit" → ["posts", "index"]
// "storage.export" → ["storage", "index"]
func routeNameToFileParts(name string) []string {
	parts := strings.Split(name, ".")

	if len(parts) == 0 {
		return []string{"index"}
	}
	// Pop the last segment (the leaf method name) and replace with "index".
	parts[len(parts)-1] = "index"

	return parts
}

// routeNamePrefix returns the portion of the route name up to (but not
// including) the final dot-segment.
//
// "posts.edit" → "posts"
// "dashboard"  → ""   (no prefix — goes into the root index.ts)
func routeNamePrefix(name string) string {
	if idx := strings.LastIndex(name, "."); idx >= 0 {
		return name[:idx]
	}

	return ""
}
