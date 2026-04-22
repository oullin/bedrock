package routing

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/bedrock/packages/routing/compiler"
	"github.com/bedrock/packages/routing/matching"
)

// Route is the Go translation of Framework\Routing\Route.
//
// In PHP, Route composes several traits (Conditionable, Macroable,
// CreatesRegularExpressionRouteConstraints, FiltersControllerMiddleware,
// ResolvesRouteDependencies). In Go, the constraint helper is composed via
// embedding [CreatesRegularExpressionRouteConstraints]; the dispatch traits
// are realized as methods on Route directly in M5.
//
// Field visibility deliberately mirrors the PHP class — public fields are
// exposed (capitalized), protected ones are unexported. This keeps tests
// translatable without reaching for getters on every line.
type Route struct {
	// Embedded helper supplying WhereAlpha / WhereNumber / etc.
	CreatesRegularExpressionRouteConstraints

	Uri                string
	HTTPMethods        []string
	ActionMap          map[string]any // mirrors PHP "$action" associative array
	IsFallback         bool
	Controller         any
	DefaultValues      map[string]any
	Wheres             map[string]string
	Parameters         map[string]string
	parameterNames     []string
	originalParameters map[string]string
	withTrashedBinds   bool
	lockSeconds        *int
	waitSeconds        *int
	computedMiddleware []any
	compiled           *compiler.CompiledRoute
	router             any // *Router — typed as any to defer the import cycle to M4
	container          BindingContainer
	bindingFields      map[string]string
	missingHandler     any
	boundModels        map[string]any // populated by ImplicitRouteBinding
}

// storeBoundModel records a resolved model instance so consumers can fetch
// the typed value (rather than the string parameter) after binding.
//
// Used by [ImplicitRouteBinding] in M9.
func (r *Route) storeBoundModel(name string, model any) {
	if r.boundModels == nil {
		r.boundModels = map[string]any{}
	}

	r.boundModels[name] = model
}

// BoundModel returns the typed model resolved for parameter name, or nil if
// no implicit binding has populated one.
func (r *Route) BoundModel(name string) any { return r.boundModels[name] }

// NewRoute constructs a Route.
//
// methods may be a single HTTP verb or a slice of verbs. action follows the
// same rules as [ParseAction]: a func, "Controller@method" string, map, or
// nil for fluent registration.
//
// Mirrors Route::__construct.
func NewRoute(methods any, uri string, action any) *Route {
	r := &Route{
		Uri:           uri,
		HTTPMethods:   normalizeMethods(methods),
		DefaultValues: map[string]any{},
		Wheres:        map[string]string{},
		bindingFields: map[string]string{},
	}
	r.CreatesRegularExpressionRouteConstraints.Bind(r)

	parsed, err := ParseAction(uri, action)

	if err != nil {
		// Mirror Upstream's behavior: store a missing-action sentinel rather
		// than panicking at construction.
		parsed = missingAction(uri)
	}

	r.ActionMap = actionToMap(parsed)
	delete(r.ActionMap, "prefix")

	// Append HEAD whenever the route handles GET — Symfony/HTTP convention.
	hasGet, hasHead := false, false

	for _, m := range r.HTTPMethods {
		if m == "GET" {
			hasGet = true
		}

		if m == "HEAD" {
			hasHead = true
		}
	}

	if hasGet && !hasHead {
		r.HTTPMethods = append(r.HTTPMethods, "HEAD")
	}

	if m, ok := action.(map[string]any); ok {
		if p, ok := m["prefix"].(string); ok {
			r.Prefix(p)
		} else {
			r.Prefix("")
		}
	} else {
		r.Prefix("")
	}

	// Re-parse the URI through ParseRouteUri so {user:slug} bindings populate
	// bindingFields. Route.uri at this point may already contain "{user}" form.
	parsedUri := ParseRouteUri(r.Uri)
	r.Uri = parsedUri.Uri

	for k, v := range parsedUri.BindingFields {
		r.bindingFields[k] = v
	}

	return r
}

func normalizeMethods(m any) []string {
	switch v := m.(type) {
	case string:
		return []string{strings.ToUpper(v)}
	case []string:
		out := make([]string, len(v))

		for i, s := range v {
			out[i] = strings.ToUpper(s)
		}

		return out
	default:
		return nil
	}
}

func actionToMap(a *Action) map[string]any {
	m := map[string]any{}

	if a == nil {
		return m
	}

	if a.Uses != nil {
		m["uses"] = a.Uses
	}

	if a.Controller != "" {
		m["controller"] = a.Controller
	}

	if len(a.Middleware) > 0 {
		m["middleware"] = a.Middleware
	}

	if len(a.Where) > 0 {
		m["where"] = a.Where
	}

	if len(a.Defaults) > 0 {
		m["defaults"] = a.Defaults
	}

	if a.Domain != "" {
		m["domain"] = a.Domain
	}

	if a.As != "" {
		m["as"] = a.As
	}

	if a.Namespace != "" {
		m["namespace"] = a.Namespace
	}

	for k, v := range a.Extras {
		m[k] = v
	}

	return m
}

// =====================================================================
// SourceRoute interface (compiler) — Path, Host, Defaults, Requirements
// =====================================================================

// Path returns the route URI prefixed with "/" so [compiler.Compile] sees a
// canonical leading slash. Mirrors the implicit slashing PHP performs on the
// path inside the Symfony layer.
func (r *Route) Path() string {
	if strings.HasPrefix(r.Uri, "/") {
		return r.Uri
	}

	return "/" + r.Uri
}

// Host returns the host pattern declared via [Route.Domain].
func (r *Route) Host() string { return r.GetDomain() }

// Requirements returns the where-constraint map.
func (r *Route) Requirements() map[string]string { return r.Wheres }

// HasDefault reports whether name has a declared default value.
func (r *Route) HasDefault(name string) bool {
	_, ok := r.DefaultValues[name]

	return ok
}

// =====================================================================
// Compile / Bind / Matches
// =====================================================================

// CompileRoute compiles the route into a [*compiler.CompiledRoute] if it has
// not been compiled yet, then returns the cached value.
//
// Mirrors Route::compileRoute.
func (r *Route) CompileRoute() (*compiler.CompiledRoute, error) {
	if r.compiled != nil {
		return r.compiled, nil
	}

	c, err := compiler.Compile(r)

	if err != nil {
		return nil, err
	}

	r.compiled = c

	return c, nil
}

// Compiled returns the cached compiled route (compiling on first call).
//
// The error from a first-time compile is swallowed to keep parity with
// PHP's getCompiled which throws — Go callers should call [Route.CompileRoute]
// directly when they need the error.
func (r *Route) Compiled() *compiler.CompiledRoute {
	c, _ := r.CompileRoute()

	return c
}

// GetCompiled is a parity alias for [Route.Compiled].
func (r *Route) GetCompiled() *compiler.CompiledRoute { return r.Compiled() }

// Bind binds the route to a request, populating Parameters via the
// [RouteParameterBinder]. Mirrors Route::bind.
func (r *Route) Bind(req boundRequest) (*Route, error) {
	if _, err := r.CompileRoute(); err != nil {
		return nil, err
	}

	r.Parameters = NewRouteParameterBinder(r).Parameters(req)
	r.originalParameters = make(map[string]string, len(r.Parameters))

	for k, v := range r.Parameters {
		r.originalParameters[k] = v
	}

	return r, nil
}

// Matches reports whether the route matches the given request.
//
// includingMethod=false skips the method validator — used by the router when
// gathering "method not allowed" candidates.
//
// Mirrors Route::matches.
func (r *Route) Matches(req matching.MatchableRequest, includingMethod bool) bool {
	if _, err := r.CompileRoute(); err != nil {
		return false
	}

	for _, v := range matching.All() {
		if !includingMethod {
			if _, isMethod := v.(matching.MethodValidator); isMethod {
				continue
			}
		}

		if !v.Matches(r, req) {
			return false
		}
	}

	return true
}

// =====================================================================
// Parameters
// =====================================================================

// ErrRouteNotBound is returned by [Route.ParametersOrErr] /
// [Route.OriginalParametersOrErr] when the route has not yet been bound.
var ErrRouteNotBound = errors.New("route is not bound")

// HasParameters reports whether the route has been bound.
func (r *Route) HasParameters() bool { return r.Parameters != nil }

// HasParameter reports whether the bound parameter map contains name.
func (r *Route) HasParameter(name string) bool {
	if !r.HasParameters() {
		return false
	}

	_, ok := r.Parameters[name]

	return ok
}

// Parameter returns the bound value of name or fallback when absent.
func (r *Route) Parameter(name, fallback string) string {
	if v, ok := r.Parameters[name]; ok {
		return v
	}

	return fallback
}

// OriginalParameter returns the pre-binding-resolution value of name.
func (r *Route) OriginalParameter(name, fallback string) string {
	if v, ok := r.originalParameters[name]; ok {
		return v
	}

	return fallback
}

// SetParameter sets a single parameter, overriding any prior value.
func (r *Route) SetParameter(name, value string) {
	if r.Parameters == nil {
		r.Parameters = map[string]string{}
	}

	r.Parameters[name] = value
}

// ForgetParameter removes name from the bound parameter map.
func (r *Route) ForgetParameter(name string) {
	delete(r.Parameters, name)
}

// ParametersOrErr returns the bound parameter map or [ErrRouteNotBound].
func (r *Route) ParametersOrErr() (map[string]string, error) {
	if r.Parameters == nil {
		return nil, ErrRouteNotBound
	}

	return r.Parameters, nil
}

// OriginalParameters returns the original (pre-binding) parameter map.
func (r *Route) OriginalParameters() map[string]string { return r.originalParameters }

// ParametersWithoutNulls returns the bound parameter map with empty values
// stripped.
func (r *Route) ParametersWithoutNulls() map[string]string {
	out := map[string]string{}

	for k, v := range r.Parameters {
		if v != "" {
			out[k] = v
		}
	}

	return out
}

// ParameterNames returns the parameter names appearing in the URI and host
// pattern. The result is computed once and cached.
func (r *Route) ParameterNames() []string {
	if r.parameterNames != nil {
		return r.parameterNames
	}

	r.parameterNames = compileParameterNames(r.GetDomain() + r.Uri)

	return r.parameterNames
}

var paramNameRe = regexp.MustCompile(`\{(.*?)\}`)

func compileParameterNames(pattern string) []string {
	matches := paramNameRe.FindAllStringSubmatch(pattern, -1)
	out := make([]string, 0, len(matches))

	for _, m := range matches {
		out = append(out, strings.TrimSuffix(m[1], "?"))
	}

	return out
}

// SignatureParameters delegates to [RouteSignatureParameters.FromAction].
func (r *Route) SignatureParameters(conditions map[string]any) []SignatureParameter {
	a := &Action{}

	if uses, ok := r.ActionMap["uses"]; ok {
		a.Uses = uses
	}

	return RouteSignatureParameters{}.FromAction(a, conditions)
}

// DefaultsMap returns the route default-values map. Implements [boundRoute].
func (r *Route) DefaultsMap() map[string]any { return r.DefaultValues }

// =====================================================================
// Binding fields / scoped bindings
// =====================================================================

// BindingFieldFor returns the custom binding field declared for parameter
// (e.g. "slug" for "{user:slug}").
func (r *Route) BindingFieldFor(parameter string) string {
	return r.bindingFields[parameter]
}

// BindingFields returns the full custom binding field map.
func (r *Route) BindingFields() map[string]string { return r.bindingFields }

// SetBindingFields replaces the custom binding field map.
func (r *Route) SetBindingFields(fields map[string]string) *Route {
	r.bindingFields = fields

	return r
}

// ParentOfParameter returns the parameter that immediately precedes the named
// parameter in the bound parameter list, or "" for the first parameter.
func (r *Route) ParentOfParameter(parameter string) string {
	keys := make([]string, 0, len(r.Parameters))

	for _, n := range r.ParameterNames() {
		if _, ok := r.Parameters[n]; ok {
			keys = append(keys, n)
		}
	}

	for i, k := range keys {
		if k == parameter {
			if i == 0 {
				return ""
			}

			return keys[i-1]
		}
	}

	return ""
}

// WithTrashed enables retrieval of soft-deleted models for implicit bindings.
func (r *Route) WithTrashed(withTrashed ...bool) *Route {
	v := true

	if len(withTrashed) > 0 {
		v = withTrashed[0]
	}

	r.withTrashedBinds = v

	return r
}

// AllowsTrashedBindings reports whether soft-deleted models may be bound.
func (r *Route) AllowsTrashedBindings() bool { return r.withTrashedBinds }

// =====================================================================
// Defaults / Where
// =====================================================================

// Defaults sets a single default value for a parameter and returns the route.
func (r *Route) Defaults(key string, value any) *Route {
	r.DefaultValues[key] = value

	return r
}

// SetDefaults replaces the route's default-value map.
func (r *Route) SetDefaults(defaults map[string]any) *Route {
	r.DefaultValues = defaults

	return r
}

// Where adds a regular-expression requirement for one parameter and returns
// the route for chaining.
func (r *Route) Where(name, expression string) *Route {
	r.Wheres[name] = expression

	return r
}

// SetWheres replaces all where requirements at once.
func (r *Route) SetWheres(wheres map[string]string) *Route {
	for k, v := range wheres {
		r.Wheres[k] = v
	}

	return r
}

// =====================================================================
// Fallback / scheme / methods
// =====================================================================

// Fallback marks the route as a fallback (last-resort) route.
func (r *Route) Fallback() *Route {
	r.IsFallback = true

	return r
}

// SetFallback explicitly sets the fallback flag.
func (r *Route) SetFallback(v bool) *Route {
	r.IsFallback = v

	return r
}

// Methods returns the HTTP verbs the route responds to.
func (r *Route) Methods() []string { return r.HTTPMethods }

// HttpOnly reports whether the route only accepts plain-HTTP requests.
func (r *Route) HttpOnly() bool { return actionContains(r.ActionMap, "http") }

// HttpsOnly reports whether the route only accepts HTTPS requests.
func (r *Route) HttpsOnly() bool { return r.Secure() }

// Secure reports whether the route is HTTPS-only.
func (r *Route) Secure() bool { return actionContains(r.ActionMap, "https") }

func actionContains(action map[string]any, target string) bool {
	for _, v := range action {
		if s, ok := v.(string); ok && s == target {
			return true
		}
	}

	return false
}

// =====================================================================
// Domain
// =====================================================================

// Domain sets (or, when called with the empty string, returns) the host
// pattern for the route. The two-mode signature mirrors PHP's domain($d=null).
func (r *Route) Domain(domain string) *Route {
	parsed := ParseRouteUri(domain)
	r.ActionMap["domain"] = parsed.Uri

	for k, v := range parsed.BindingFields {
		r.bindingFields[k] = v
	}

	return r
}

// GetDomain returns the host pattern with any scheme prefix stripped.
func (r *Route) GetDomain() string {
	v, ok := r.ActionMap["domain"]

	if !ok {
		return ""
	}

	s, _ := v.(string)
	s = strings.TrimPrefix(s, "http://")
	s = strings.TrimPrefix(s, "https://")

	return s
}

// =====================================================================
// Prefix / URI
// =====================================================================

// GetPrefix returns the URI prefix attached to the route (or "").
func (r *Route) GetPrefix() string {
	if v, ok := r.ActionMap["prefix"]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}

	return ""
}

// Prefix prepends the supplied prefix to the route URI and stores it in the
// action map so group attributes can be merged later.
func (r *Route) Prefix(prefix string) *Route {
	r.updatePrefixOnAction(prefix)
	uri := strings.TrimRight(prefix, "/") + "/" + strings.TrimLeft(r.Uri, "/")

	if uri != "/" {
		uri = strings.Trim(uri, "/")
	}

	r.SetUri(uri)

	return r
}

func (r *Route) updatePrefixOnAction(prefix string) {
	existing := ""

	if v, ok := r.ActionMap["prefix"]; ok {
		existing, _ = v.(string)
	}

	combined := strings.Trim(strings.TrimRight(prefix, "/")+"/"+strings.TrimLeft(existing, "/"), "/")

	if combined != "" {
		r.ActionMap["prefix"] = combined
	}
}

// GetUri returns the canonical URI string.
func (r *Route) GetUri() string { return r.Uri }

// SetUri rewrites the URI, re-parsing any "{name:field}" segments.
func (r *Route) SetUri(uri string) *Route {
	r.bindingFields = map[string]string{}
	parsed := ParseRouteUri(uri)
	r.Uri = parsed.Uri

	for k, v := range parsed.BindingFields {
		r.bindingFields[k] = v
	}

	return r
}

// =====================================================================
// Name / action / uses
// =====================================================================

// GetName returns the route's name (the "as" key), or "" if unnamed.
func (r *Route) GetName() string {
	if v, ok := r.ActionMap["as"]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}

	return ""
}

// Name appends to the route's name. Mirrors Route::name which concatenates
// when a name is already present (used by RouteRegistrar to compose group
// name prefixes).
func (r *Route) Name(name string) *Route {
	existing := r.GetName()
	r.ActionMap["as"] = existing + name

	return r
}

// Named reports whether the route's name matches any of the supplied glob
// patterns. "users.*" matches any name beginning with "users.".
func (r *Route) Named(patterns ...string) bool {
	name := r.GetName()

	if name == "" {
		return false
	}

	for _, p := range patterns {
		if matchPattern(p, name) {
			return true
		}
	}

	return false
}

// matchPattern is a minimal Str::is implementation supporting a single "*"
// wildcard per segment. Faithful enough for the routing test surface.
func matchPattern(pattern, value string) bool {
	if pattern == value {
		return true
	}

	if !strings.Contains(pattern, "*") {
		return false
	}

	parts := strings.Split(pattern, "*")

	if !strings.HasPrefix(value, parts[0]) {
		return false
	}

	value = value[len(parts[0]):]

	for i := 1; i < len(parts); i++ {
		idx := strings.Index(value, parts[i])

		if idx < 0 {
			return false
		}

		value = value[idx+len(parts[i]):]
	}

	return strings.HasSuffix(pattern, "*") || value == ""
}

// Uses sets (or replaces) the handler for the route.
func (r *Route) Uses(action any) *Route {
	parsed, _ := ParseAction(r.Uri, action)

	for k, v := range actionToMap(parsed) {
		r.ActionMap[k] = v
	}

	return r
}

// GetActionName returns the canonical "Class@method" name, or "Closure".
func (r *Route) GetActionName() string {
	if v, ok := r.ActionMap["controller"]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}

	if v, ok := r.ActionMap["uses"]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}

	return "Closure"
}

// GetControllerClass returns the class portion of a controller action.
func (r *Route) GetControllerClass() string {
	name := r.GetActionName()

	if name == "Closure" {
		return ""
	}

	if idx := strings.Index(name, "@"); idx >= 0 {
		return name[:idx]
	}

	return name
}

// FlushController clears any cached controller instance.
func (r *Route) FlushController() *Route {
	r.Controller = nil

	return r
}

// GetActionMethod returns the method portion of "Class@method", or "" for
// closures.
func (r *Route) GetActionMethod() string {
	name := r.GetActionName()

	if idx := strings.Index(name, "@"); idx >= 0 {
		return name[idx+1:]
	}

	return ""
}

// GetAction returns the raw action map. When key is non-empty the value at
// that key is returned; otherwise the entire map is returned.
func (r *Route) GetAction(key ...string) any {
	if len(key) == 0 || key[0] == "" {
		return r.ActionMap
	}

	return r.ActionMap[key[0]]
}

// SetAction replaces the action map and returns the route.
func (r *Route) SetAction(action map[string]any) *Route {
	r.ActionMap = action

	return r
}

// =====================================================================
// Missing / scoped bindings / blocking
// =====================================================================

// Missing assigns a callback that runs when implicit binding cannot resolve
// a parameter for this route.
func (r *Route) Missing(handler func(req any) any) *Route {
	r.missingHandler = handler

	return r
}

// GetMissing returns the missing-binding handler set via [Route.Missing].
func (r *Route) GetMissing() any { return r.missingHandler }

// ScopeBindings enables scoped implicit binding for child route parameters.
func (r *Route) ScopeBindings() *Route {
	r.ActionMap["scope_bindings"] = true

	return r
}

// WithoutScopedBindings explicitly disables scoped implicit binding.
func (r *Route) WithoutScopedBindings() *Route {
	r.ActionMap["scope_bindings"] = false

	return r
}

// EnforcesScopedBindings reports whether scoped binding is on.
func (r *Route) EnforcesScopedBindings() bool {
	v, ok := r.ActionMap["scope_bindings"]

	if !ok {
		return false
	}

	b, _ := v.(bool)

	return b
}

// PreventsScopedBindings reports whether scoped binding has been disabled.
func (r *Route) PreventsScopedBindings() bool {
	v, ok := r.ActionMap["scope_bindings"]

	if !ok {
		return false
	}

	b, _ := v.(bool)

	return !b
}

// Block enables session blocking with the supplied lock and wait windows
// (in seconds). Mirrors Route::block.
func (r *Route) Block(lockSeconds, waitSeconds int) *Route {
	r.lockSeconds = &lockSeconds
	r.waitSeconds = &waitSeconds

	return r
}

// WithoutBlocking disables previously-configured session blocking.
func (r *Route) WithoutBlocking() *Route {
	r.lockSeconds = nil
	r.waitSeconds = nil

	return r
}

// LocksFor returns the configured lock duration in seconds.
func (r *Route) LocksFor() int {
	if r.lockSeconds == nil {
		return 0
	}

	return *r.lockSeconds
}

// WaitsFor returns the configured wait duration in seconds.
func (r *Route) WaitsFor() int {
	if r.waitSeconds == nil {
		return 0
	}

	return *r.waitSeconds
}

// =====================================================================
// Wiring
// =====================================================================

// SetRouter records the parent router. The parameter is typed as any to keep
// the package free of an import cycle on M4's [Router] type.
func (r *Route) SetRouter(router any) *Route {
	r.router = router

	return r
}

// SetContainer records the container used to resolve dependencies.
func (r *Route) SetContainer(c BindingContainer) *Route {
	r.container = c

	return r
}

// =====================================================================
// String form
// =====================================================================

func (r *Route) String() string {
	return fmt.Sprintf("Route{%s %s -> %s}", strings.Join(r.HTTPMethods, "|"), r.Uri, r.GetActionName())
}
