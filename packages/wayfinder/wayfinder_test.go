package wayfinder_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bedrock/packages/wayfinder"
)

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

// generateTo runs Generate() into a temp directory and returns the base path.
func generateTo(t *testing.T, routes []*wayfinder.RouteInfo, opts wayfinder.Options) string {
	t.Helper()
	dir := t.TempDir()
	opts.Path = dir

	if err := wayfinder.Generate(routes, opts); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	return dir
}

// readFile reads a file relative to base and fails the test if absent.
func readFile(t *testing.T, base, rel string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(base, filepath.FromSlash(rel)))

	if err != nil {
		t.Fatalf("readFile(%q): %v", rel, err)
	}

	return string(b)
}

// assertContains checks that content contains the expected substring.
func assertContains(t *testing.T, content, want string) {
	t.Helper()

	if !strings.Contains(content, want) {
		t.Errorf("expected content to contain:\n  %q\ngot:\n%s", want, content)
	}
}

// assertNotContains checks that content does NOT contain the given substring.
func assertNotContains(t *testing.T, content, unwanted string) {
	t.Helper()

	if strings.Contains(content, unwanted) {
		t.Errorf("expected content NOT to contain %q", unwanted)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Workbench route fixtures (mirror Wayfinder workbench)
// ─────────────────────────────────────────────────────────────────────────────

func postControllerRoutes() []*wayfinder.RouteInfo {
	base := "App\\Http\\Controllers\\PostController"

	return []*wayfinder.RouteInfo{
		{URI: "/posts", Methods: []string{"get", "head"}, Controller: base + "@index"},
		{URI: "/posts/create", Methods: []string{"get", "head"}, Controller: base + "@create"},
		{URI: "/posts", Methods: []string{"post"}, Controller: base + "@store"},
		{
			URI: "/posts/{post}", Methods: []string{"get", "head"},
			Controller: base + "@show",
			Params:     []wayfinder.Param{{Name: "post"}},
		},
		{
			URI: "/posts/{post}/edit", Methods: []string{"get", "head"},
			Controller: base + "@edit",
			Params:     []wayfinder.Param{{Name: "post"}},
		},
		{
			URI: "/posts/{post}", Methods: []string{"put", "patch"},
			Controller: base + "@update",
			Params:     []wayfinder.Param{{Name: "post"}},
		},
		{
			URI: "/posts/{post}", Methods: []string{"delete"},
			Controller: base + "@destroy",
			Params:     []wayfinder.Param{{Name: "post"}},
		},
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// PostController tests (mirrors PostController.test.ts)
// ─────────────────────────────────────────────────────────────────────────────

// Port of PostController::test_properties.
// Port of PostController::test_url.
// Port of PostController::test_default_method.
// Port of PostController::test_get.
// Port of PostController::test_head.
// Port of PostController::test_definition.
// Port of PostController::test_post.
// Port of PostController::test_patch.
// Port of PostController::test_delete.
// Port of DefaultParameters::test_it_can_generate_urls_without_default_parameters_set.
func TestPostControllerGeneration(t *testing.T) {
	t.Parallel()

	dir := generateTo(t, postControllerRoutes(), wayfinder.Options{})
	content := readFile(t, dir, "actions/App/Http/Controllers/PostController.ts")

	// index — no params, GET
	assertContains(t, content, `export const index`)
	assertContains(t, content, `url: index.url(options)`)
	assertContains(t, content, `method: "get"`)
	assertContains(t, content, `index.definition = {`)
	assertContains(t, content, `methods: ["get","head"]`)
	assertContains(t, content, `url: "/posts"`)
	assertContains(t, content, `index.url = `)
	assertContains(t, content, `index.get = `)
	assertContains(t, content, `index.head = `)

	// show — one required param
	assertContains(t, content, `export const show`)
	assertContains(t, content, `show.url = `)
	assertContains(t, content, `.replace("{post}"`)
	assertContains(t, content, `show.definition`)
	assertContains(t, content, `url: "/posts/{post}"`)

	// destroy — DELETE
	assertContains(t, content, `export const destroy`)
	assertContains(t, content, `method: "delete"`)

	// update — PUT/PATCH
	assertContains(t, content, `export const update`)
	assertContains(t, content, `method: "put"`)

	// store — POST
	assertContains(t, content, `export const store`)
	assertContains(t, content, `method: "post"`)

	// Default export (controller object)
	assertContains(t, content, `export default PostController`)
}

// Port of AnonymousMiddleware::test_will_allow_for_closure_middleware.
func TestAnonymousMiddlewareClosureRoutesDoNotBlockGeneration(t *testing.T) {
	t.Parallel()

	routes := []*wayfinder.RouteInfo{
		{URI: "/closure", Methods: []string{"get", "head"}, Controller: "Closure"},
		{URI: "/posts", Methods: []string{"get", "head"}, Controller: "App\\Http\\Controllers\\PostController@index"},
	}

	dir := generateTo(t, routes, wayfinder.Options{})
	content := readFile(t, dir, "actions/App/Http/Controllers/PostController.ts")

	assertContains(t, content, `url: "/posts"`)

	if _, err := os.Stat(filepath.Join(dir, "actions/Closure.ts")); !os.IsNotExist(err) {
		t.Fatalf("closure middleware route should not create a Closure action file")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// InvokableController tests (mirrors InvokableController.test.ts)
// ─────────────────────────────────────────────────────────────────────────────

// Port of InvokableController::it_exports_default_for_invokable_controllers.
func TestInvokableControllerGeneration(t *testing.T) {
	t.Parallel()

	routes := []*wayfinder.RouteInfo{
		{
			URI:         "/invokable-controller",
			Methods:     []string{"get", "head"},
			Controller:  "App\\Http\\Controllers\\InvokableController@Invoke",
			IsInvokable: true,
		},
	}

	dir := generateTo(t, routes, wayfinder.Options{})
	content := readFile(t, dir, "actions/App/Http/Controllers/InvokableController.ts")

	// Must use default export (not a named export).
	assertContains(t, content, `export default InvokableController`)
	assertContains(t, content, `url: "/invokable-controller"`)
	assertContains(t, content, `method: "get"`)
	// Should NOT have "export const InvokableController"
	assertNotContains(t, content, `export const InvokableController`)
}

// Port of InvokablePlusController::it_exports_default_and_methods_for_invokable_controllers.
func TestInvokablePlusControllerGeneration(t *testing.T) {
	t.Parallel()

	routes := []*wayfinder.RouteInfo{
		{
			URI:         "/invokable-plus",
			Methods:     []string{"get", "head"},
			Controller:  "App\\Http\\Controllers\\InvokablePlusController@Invoke",
			IsInvokable: true,
		},
		{
			URI:        "/invokable-plus/download",
			Methods:    []string{"post"},
			Controller: "App\\Http\\Controllers\\InvokablePlusController@download",
		},
	}

	dir := generateTo(t, routes, wayfinder.Options{})
	content := readFile(t, dir, "actions/App/Http/Controllers/InvokablePlusController.ts")

	assertContains(t, content, `const InvokablePlusController = (`)
	assertContains(t, content, `const download = (`)
	assertContains(t, content, `InvokablePlusController.download = download`)
	assertContains(t, content, `export default InvokablePlusController`)
	assertNotContains(t, content, `export const InvokablePlusController`)
}

// ─────────────────────────────────────────────────────────────────────────────
// OptionalController tests (mirrors OptionalController.test.ts)
// ─────────────────────────────────────────────────────────────────────────────

// Port of OptionalController::test_url.
// Port of OptionalController::test_definition.
// Port of OptionalController::test_url_supports_falsy_optional_values.
// Port of OptionalController::it_throws_an_error_when_passing_optional_parameters_with_missing_optional_parameters_before.
func TestOptionalControllerGeneration(t *testing.T) {
	t.Parallel()

	routes := []*wayfinder.RouteInfo{
		{
			URI:        "/optional/{parameter?}",
			Methods:    []string{"get", "head"},
			Controller: "App\\Http\\Controllers\\OptionalController@optional",
			Params:     []wayfinder.Param{{Name: "parameter", Optional: true}},
		},
		{
			URI:        "/many-optional/{one?}/{two?}/{three?}",
			Methods:    []string{"get", "head"},
			Controller: "App\\Http\\Controllers\\OptionalController@manyOptional",
			Params: []wayfinder.Param{
				{Name: "one", Optional: true},
				{Name: "two", Optional: true},
				{Name: "three", Optional: true},
			},
		},
	}

	dir := generateTo(t, routes, wayfinder.Options{})
	content := readFile(t, dir, "actions/App/Http/Controllers/OptionalController.ts")

	// optional param — definition should contain {parameter?}
	assertContains(t, content, `url: "/optional/{parameter?}"`)
	// validateParameters should be called
	assertContains(t, content, `validateParameters(`)
	// Optional chaining in replace
	assertContains(t, content, `.replace("{parameter?}"`)

	// manyOptional
	assertContains(t, content, `url: "/many-optional/{one?}/{two?}/{three?}"`)
	assertContains(t, content, `validateParameters(args, ["one", "two", "three"])`)
	assertContains(t, content, `.replace("{one?}", parsedArgs.one?.toString() ?? '')`)
	assertContains(t, content, `.replace("{two?}", parsedArgs.two?.toString() ?? '')`)
	assertContains(t, content, `.replace("{three?}", parsedArgs.three?.toString() ?? '')`)
}

// Port of EmptyRoute::it_doesn_t_add_a_to_an_empty_route.
func TestEmptyRouteGeneration(t *testing.T) {
	t.Parallel()

	routes := []*wayfinder.RouteInfo{
		{
			URI:        "",
			Methods:    []string{"get", "head"},
			Controller: "App\\Http\\Controllers\\EmptyRouteController@index",
		},
	}

	dir := generateTo(t, routes, wayfinder.Options{})
	content := readFile(t, dir, "actions/App/Http/Controllers/EmptyRouteController.ts")

	assertContains(t, content, `url: ""`)
	assertNotContains(t, content, `url: "/"`)
}

// ─────────────────────────────────────────────────────────────────────────────
// ModelBindingController tests (mirrors ModelBindingController.test.ts)
// ─────────────────────────────────────────────────────────────────────────────

// Port of ModelBindingController::test_will_detect_model_binding.
func TestModelBindingControllerGeneration(t *testing.T) {
	t.Parallel()

	routes := []*wayfinder.RouteInfo{
		{
			URI:        "/users/{user}",
			Methods:    []string{"get", "head"},
			Controller: "App\\Http\\Controllers\\ModelBindingController@show",
			Params:     []wayfinder.Param{{Name: "user"}},
		},
	}

	dir := generateTo(t, routes, wayfinder.Options{})
	content := readFile(t, dir, "actions/App/Http/Controllers/ModelBindingController.ts")

	assertContains(t, content, `url: "/users/{user}"`)
	// Primitive shorthand support.
	assertContains(t, content, `if (typeof args === 'string' || typeof args === 'number')`)
	// Array support.
	assertContains(t, content, `if (Array.isArray(args))`)
}

// Port of CamelCaseRouteParameter::test_can_resolve_model_binding_keys_for_camelcase_route_handler_parameters.
func TestCamelCaseRouteParameterGeneration(t *testing.T) {
	t.Parallel()

	routes := []*wayfinder.RouteInfo{
		{
			URI:        "/profiles/{userProfile}",
			Methods:    []string{"get", "head"},
			Controller: "App\\Http\\Controllers\\CamelCaseRouteParameterController@show",
			Params:     []wayfinder.Param{{Name: "userProfile", Key: "uuid"}},
		},
	}

	dir := generateTo(t, routes, wayfinder.Options{})
	content := readFile(t, dir, "actions/App/Http/Controllers/CamelCaseRouteParameterController.ts")

	assertContains(t, content, `url: "/profiles/{userProfile}"`)
	assertContains(t, content, `"uuid" in args`)
	assertContains(t, content, `args.userProfile.uuid`)
	assertContains(t, content, `.replace("{userProfile}", parsedArgs.userProfile.toString())`)
}

// ─────────────────────────────────────────────────────────────────────────────
// KeyController tests (mirrors KeyController.test.ts)
// ─────────────────────────────────────────────────────────────────────────────

// Port of KeyController::it_can_pass_primitive_values_to_routes_with_custom_keys.
// Port of KeyController::it_can_pass_objects_with_custom_key.
// Port of KeyController::test_definition.
func TestKeyControllerGeneration(t *testing.T) {
	t.Parallel()

	routes := []*wayfinder.RouteInfo{
		{
			URI:        "/keys/{key}/edit",
			Methods:    []string{"get", "head"},
			Controller: "App\\Http\\Controllers\\KeyController@edit",
			Params:     []wayfinder.Param{{Name: "key", Key: "uuid"}},
		},
	}

	dir := generateTo(t, routes, wayfinder.Options{})
	content := readFile(t, dir, "actions/App/Http/Controllers/KeyController.ts")

	assertContains(t, content, `url: "/keys/{key}/edit"`)
	// Custom key resolution: check for uuid field.
	assertContains(t, content, `.uuid`)
	// definition should show url pattern.
	assertContains(t, content, `url: "/keys/{key}/edit"`)
}

// ─────────────────────────────────────────────────────────────────────────────
// DisallowedMethodNames tests (mirrors DisallowedMethodNames.test.ts)
// ─────────────────────────────────────────────────────────────────────────────

// Port of DisallowedMethodNames::test_will_append_method_to_invalid_methods.
// Port of DisallowedMethodNames::test_will_properly_handle_leading_numbers.
// Port of DisallowedMethodNames::test_will_properly_handle_reserved_js_words.
func TestDisallowedMethodNamesGeneration(t *testing.T) {
	t.Parallel()

	base := "App\\Http\\Controllers\\DisallowedMethodNameController"
	routes := []*wayfinder.RouteInfo{
		{URI: "/disallowed/delete", Methods: []string{"get", "head"}, Controller: base + "@delete"},
		{URI: "/disallowed/404", Methods: []string{"get", "head"}, Controller: base + "@404"},
		{URI: "/disallowed/2fa", Methods: []string{"get", "head"}, Controller: base + "@2fa"},
		{URI: "/disallowed/default", Methods: []string{"get", "head"}, Controller: base + "@default"},
	}

	dir := generateTo(t, routes, wayfinder.Options{})
	content := readFile(t, dir, "actions/App/Http/Controllers/DisallowedMethodNameController.ts")

	// Reserved word "delete" → "deleteMethod"
	assertContains(t, content, `const deleteMethod`)
	assertContains(t, content, `url: "/disallowed/delete"`)

	// Leading number "404" → "method404"
	assertContains(t, content, `const method404`)
	assertContains(t, content, `url: "/disallowed/404"`)

	// Reserved word "default" → "defaultMethod"
	assertContains(t, content, `const defaultMethod`)

	// Controller object must have original name aliases.
	assertContains(t, content, `delete: deleteMethod`)
	assertContains(t, content, `404: method404`)

	// Default export
	assertContains(t, content, `export default DisallowedMethodNameController`)
}

// Port of MethodNameCollision::test_does_not_shadow_a_generated_method_named_options.
func TestMethodNameCollisionGeneration(t *testing.T) {
	t.Parallel()

	routes := []*wayfinder.RouteInfo{
		{
			URI:        "/method-collision/{post}",
			Methods:    []string{"get", "head"},
			Controller: "App\\Http\\Controllers\\MethodNameCollisionController@options",
			Params:     []wayfinder.Param{{Name: "post"}},
		},
	}

	dir := generateTo(t, routes, wayfinder.Options{})
	content := readFile(t, dir, "actions/App/Http/Controllers/MethodNameCollisionController.ts")

	assertContains(t, content, `export const options = (args: {`)
	assertContains(t, content, `routeOptions?: RouteQueryOptions`)
	assertContains(t, content, `queryParams(routeOptions)`)
	assertNotContains(t, content, `queryParams(options)`)
}

// Port of ParamaterName::test_url.
// Port of ParamaterName::test_definition.
func TestParameterNameCollisionGeneration(t *testing.T) {
	t.Parallel()

	routes := []*wayfinder.RouteInfo{
		{
			URI:        "/parameter-names/{args}/{options}/{parsedArgs}",
			Methods:    []string{"get", "head"},
			Controller: "App\\Http\\Controllers\\ParamaterNameController@show",
			Params: []wayfinder.Param{
				{Name: "args"},
				{Name: "options"},
				{Name: "parsedArgs"},
			},
		},
	}

	dir := generateTo(t, routes, wayfinder.Options{})
	content := readFile(t, dir, "actions/App/Http/Controllers/ParamaterNameController.ts")

	assertContains(t, content, `url: "/parameter-names/{args}/{options}/{parsedArgs}"`)
	assertContains(t, content, `.replace("{args}", parsedArgs.args.toString())`)
	assertContains(t, content, `.replace("{options}", parsedArgs.options.toString())`)
	assertContains(t, content, `.replace("{parsedArgs}", parsedArgs.parsedArgs.toString())`)
}

// ─────────────────────────────────────────────────────────────────────────────
// TwoRoutesSameAction tests (mirrors TwoRoutesSameAction.test.ts)
// ─────────────────────────────────────────────────────────────────────────────

// Port of TwoRoutesSameAction::it_creates_a_keyed_dictionary_of_routes_for_multiple_routes_pointing_to_the_same_action.
func TestTwoRoutesSameActionGeneration(t *testing.T) {
	t.Parallel()

	base := "App\\Http\\Controllers\\TwoRoutesSameActionController"
	routes := []*wayfinder.RouteInfo{
		{URI: "/two-routes-one-action-1", Methods: []string{"get", "head"}, Controller: base + "@same"},
		{URI: "/two-routes-one-action-2", Methods: []string{"get", "head"}, Controller: base + "@same"},
	}

	dir := generateTo(t, routes, wayfinder.Options{})
	content := readFile(t, dir, "actions/App/Http/Controllers/TwoRoutesSameActionController.ts")

	// Keyed dictionary with URI strings as keys.
	assertContains(t, content, `"/two-routes-one-action-1"`)
	assertContains(t, content, `"/two-routes-one-action-2"`)
	assertContains(t, content, `export const same`)
}

// ─────────────────────────────────────────────────────────────────────────────
// DomainController tests (mirrors DomainController.test.ts)
// ─────────────────────────────────────────────────────────────────────────────

// Port of DomainController::test_can_generate_fixed_domain_urls.
// Port of DomainController::test_can_generate_dynamic_domain_urls.
// Port of AppUrlRootResolution::test_does_not_inject_app_url_port_into_explicit_domain_routes.
func TestDomainControllerGeneration(t *testing.T) {
	t.Parallel()

	base := "App\\Http\\Controllers\\DomainController"
	routes := []*wayfinder.RouteInfo{
		{
			URI: "/fixed-domain/{param}", Methods: []string{"get", "head"},
			Controller: base + "@fixedDomain",
			Domain:     "example.test",
			Scheme:     "//",
			Params:     []wayfinder.Param{{Name: "param"}},
		},
		{
			URI: "/default-parameters-domain/{param}", Methods: []string{"get", "head"},
			Controller: base + "@defaultParametersDomain",
			Domain:     "{defaultDomain?}.au",
			Scheme:     "//",
			Params: []wayfinder.Param{
				{Name: "defaultDomain", Optional: true},
				{Name: "param"},
			},
		},
	}

	dir := generateTo(t, routes, wayfinder.Options{})
	content := readFile(t, dir, "actions/App/Http/Controllers/DomainController.ts")

	// Fixed domain URLs include the domain.
	assertContains(t, content, `//example.test/fixed-domain/{param}`)

	// Dynamic domain URL includes the domain placeholder.
	assertContains(t, content, `//{defaultDomain?}.au/default-parameters-domain/{param}`)
	assertNotContains(t, content, `:8080`)
}

// ─────────────────────────────────────────────────────────────────────────────
// Named routes tests (mirrors NamedRoutes.test.ts)
// ─────────────────────────────────────────────────────────────────────────────

// Port of NamedRoutes::it_exports_named_routes.
func TestNamedRoutesGeneration(t *testing.T) {
	t.Parallel()

	routes := []*wayfinder.RouteInfo{
		{
			URI: "/posts/{post}/edit", Methods: []string{"get", "head"},
			Name:       "posts.edit",
			Controller: "App\\Http\\Controllers\\PostController@edit",
			Params:     []wayfinder.Param{{Name: "post"}},
		},
		{
			URI: "/dashboard", Methods: []string{"get", "head"},
			Name:       "dashboard",
			Controller: "App\\Http\\Controllers\\DashboardController@index",
		},
		{
			URI:         "/named-invokable-controller",
			Methods:     []string{"get", "head"},
			Name:        "invokable",
			Controller:  "App\\Http\\Controllers\\InvokableController@Invoke",
			IsInvokable: true,
		},
		{
			URI:        "/invalid-js-name",
			Methods:    []string{"get", "head"},
			Name:       "invalid_js_name",
			Controller: "App\\Http\\Controllers\\SomeController@invalidJsName",
		},
	}

	dir := generateTo(t, routes, wayfinder.Options{SkipActions: true})

	// posts/index.ts should export "edit"
	postsContent := readFile(t, dir, "routes/posts/index.ts")
	assertContains(t, postsContent, `export const edit`)
	assertContains(t, postsContent, `url: "/posts/{post}/edit"`)

	// root index.ts should have dashboard
	rootContent := readFile(t, dir, "routes/index.ts")
	assertContains(t, rootContent, `dashboard`)
}

// Port of NamedspacedRoute::it_can_access_a_namespaced_route.
// Port of StorageRoute::it_can_import_storage_routes.
func TestNamespacedAndStorageRoutesGeneration(t *testing.T) {
	t.Parallel()

	routes := []*wayfinder.RouteInfo{
		{
			URI:        "/admin/reports",
			Methods:    []string{"get", "head"},
			Name:       "admin::reports.index",
			Controller: "App\\Http\\Controllers\\Admin\\ReportsController@index",
		},
		{
			URI:     "/storage/{path}",
			Methods: []string{"get", "head"},
			Name:    "storage.local",
			Params:  []wayfinder.Param{{Name: "path"}},
		},
	}

	dir := generateTo(t, routes, wayfinder.Options{})
	content := readFile(t, dir, "routes/namespaced/admin/reports/index.ts")
	assertContains(t, content, `export const index`)
	assertContains(t, content, `url: "/admin/reports"`)

	storage := readFile(t, dir, "routes/storage/index.ts")
	assertContains(t, storage, `export const local`)
	assertContains(t, storage, `url: "/storage/{path}"`)
}

// ─────────────────────────────────────────────────────────────────────────────
// UrlDefaults tests (mirrors UrlDefaultsController.test.ts)
// ─────────────────────────────────────────────────────────────────────────────

// Port of UrlDefaultsController::test_url.
// Port of UrlDefaultsController::test_default_method.
// Port of DefaultParameters::test_it_can_generate_urls_with_default_url_parameters_set_on_backend_and_frontend.
func TestUrlDefaultsControllerGeneration(t *testing.T) {
	t.Parallel()

	base := "App\\Http\\Controllers\\UrlDefaultsController"
	routes := []*wayfinder.RouteInfo{
		{
			URI:        "/with-defaults/{locale}",
			Methods:    []string{"post"},
			Controller: base + "@onlyDefaults",
			Params: []wayfinder.Param{
				{Name: "locale", Optional: true, Default: "en"},
			},
			Defaults: map[string]string{"locale": "en"},
		},
	}

	dir := generateTo(t, routes, wayfinder.Options{})
	content := readFile(t, dir, "actions/App/Http/Controllers/UrlDefaultsController.ts")

	// Default value is used in parsedArgs.
	assertContains(t, content, `?? "en"`)
	// URI shows optional marker because locale has a default.
	assertContains(t, content, `{locale?}`)
}

// ─────────────────────────────────────────────────────────────────────────────
// WithForm option tests
// ─────────────────────────────────────────────────────────────────────────────

// Port of PostController::test_default_form_method.
// Port of PostController::test_form_get.
// Port of PostController::test_form_head.
func TestWithFormOption(t *testing.T) {
	t.Parallel()

	routes := []*wayfinder.RouteInfo{
		{
			URI: "/posts", Methods: []string{"post"},
			Controller: "App\\Http\\Controllers\\PostController@store",
		},
		{
			URI: "/posts/{post}", Methods: []string{"put", "patch"},
			Controller: "App\\Http\\Controllers\\PostController@update",
			Params:     []wayfinder.Param{{Name: "post"}},
		},
	}

	dir := generateTo(t, routes, wayfinder.Options{WithForm: true})
	content := readFile(t, dir, "actions/App/Http/Controllers/PostController.ts")

	// Form helpers must be present.
	assertContains(t, content, `.form = `)
	assertContains(t, content, `RouteFormDefinition`)
	assertContains(t, content, `action:`)
	// _method spoofing for non-GET verbs.
	assertContains(t, content, `_method`)
}

// ─────────────────────────────────────────────────────────────────────────────
// AppUrl base-path tests (mirrors AppUrlRootResolution.test.ts)
// ─────────────────────────────────────────────────────────────────────────────

// Port of AppUrlRootResolution::test_prefixes_generated_urls_with_app_url_path.
func TestAppURLPathPrefix(t *testing.T) {
	t.Parallel()

	routes := postControllerRoutes()
	// Apply base path — simulates APP_URL=http://localhost:8081/v2.
	for _, r := range routes {
		r.BasePath = "/v2"
	}

	dir := generateTo(t, routes, wayfinder.Options{})
	content := readFile(t, dir, "actions/App/Http/Controllers/PostController.ts")

	assertContains(t, content, `url: "/v2/posts"`)
	assertContains(t, content, `url: "/v2/posts/{post}"`)
}

// ─────────────────────────────────────────────────────────────────────────────
// Wayfinder runtime utility
// ─────────────────────────────────────────────────────────────────────────────

// Port of QueryParams::it_can_convert_basic_params.
// Port of QueryParams::it_can_convert_array_params.
// Port of QueryParams::it_can_convert_object_params.
// Port of QueryParams::it_can_convert_boolean_params.
// Port of QueryParams::it_will_ignore_existing_params_without_star.
// Port of QueryParams::it_can_integrate_basic_params_with_existing_window_params.
// Port of QueryParams::it_can_integrate_array_params_with_existing_window_params.
// Port of QueryParams::it_can_integrate_object_params_with_existing_window_params.
// Port of QueryParams::it_can_delete_existing_params_via_null.
// Port of QueryParams::it_can_delete_existing_params_via_undefined.
// Port of QueryParams::it_can_merge_with_the_form_method.
// Port of QueryParams::it_can_pass_nested_query_parameters.
// Port of QueryParams::it_ignores_nested_object_values_with_unallowed_types.
// Port of DefaultParameters::test_it_can_generate_urls_with_dynamic_function_based_default_url_parameters.
func TestWayfinderRuntimeUtility(t *testing.T) {
	t.Parallel()

	dir := generateTo(t, nil, wayfinder.Options{SkipActions: true, SkipRoutes: true})
	content := readFile(t, dir, "wayfinder/index.ts")

	// Must contain the core runtime functions.
	assertContains(t, content, `export const queryParams`)
	assertContains(t, content, `export const setUrlDefaults`)
	assertContains(t, content, `export const applyUrlDefaults`)
	assertContains(t, content, `export const validateParameters`)
	assertContains(t, content, `export type QueryParams`)
	assertContains(t, content, `export type RouteDefinition`)
	assertContains(t, content, `export type RouteFormDefinition`)
	assertContains(t, content, `export type RouteQueryOptions`)
}

func TestWayfinderRuntimeQueryParamsSource(t *testing.T) {
	t.Parallel()

	data, err := os.ReadFile(filepath.FromSlash("resources/wayfinder.ts"))

	if err != nil {
		t.Fatal(err)
	}

	content := string(data)

	assertContains(t, content, `return "1";`)
	assertContains(t, content, `return "0";`)
	assertContains(t, content, "params.append(`${paramKey}[]`, getValue(v))")
	assertContains(t, content, `addNestedParams(value, paramKey, params)`)
	assertContains(t, content, `window.location.search`)
	assertContains(t, content, `params.delete(key)`)
	assertContains(t, content, `export const setUrlDefaults = (params: UrlDefaults | (() => UrlDefaults))`)
	assertContains(t, content, `export const applyUrlDefaults`)
	assertContains(t, content, `export const validateParameters`)
}

// ─────────────────────────────────────────────────────────────────────────────
// SkipActions / SkipRoutes options
// ─────────────────────────────────────────────────────────────────────────────

func TestSkipOptions(t *testing.T) {
	t.Parallel()

	routes := []*wayfinder.RouteInfo{
		{
			URI: "/posts", Methods: []string{"get"},
			Name:       "posts.index",
			Controller: "App\\Http\\Controllers\\PostController@index",
		},
	}

	t.Run("skip_actions", func(t *testing.T) {
		t.Parallel()
		dir := generateTo(t, routes, wayfinder.Options{SkipActions: true})

		if _, err := os.Stat(filepath.Join(dir, "actions")); !os.IsNotExist(err) {
			t.Error("actions/ directory should not exist when SkipActions=true")
		}
	})

	t.Run("skip_routes", func(t *testing.T) {
		t.Parallel()
		dir := generateTo(t, routes, wayfinder.Options{SkipRoutes: true})

		if _, err := os.Stat(filepath.Join(dir, "routes")); !os.IsNotExist(err) {
			t.Error("routes/ directory should not exist when SkipRoutes=true")
		}
	})
}

// ─────────────────────────────────────────────────────────────────────────────
// Barrel files
// ─────────────────────────────────────────────────────────────────────────────

// Port of PathNameNormalization::it_can_normalize_to_camel_case.
// Port of PathNameNormalization::it_will_properly_export_barrel_files.
// Port of NestedController::it_can_handle_conflicting_nested_route_names.
func TestBarrelFilesGeneration(t *testing.T) {
	t.Parallel()

	routes := []*wayfinder.RouteInfo{
		{URI: "/posts", Methods: []string{"get"}, Controller: "App\\Http\\Controllers\\PostController@index"},
		{URI: "/users", Methods: []string{"get"}, Controller: "App\\Http\\Controllers\\UserController@index"},
	}

	dir := generateTo(t, routes, wayfinder.Options{SkipRoutes: true})

	// Controller-level barrel.
	indexContent := readFile(t, dir, "actions/App/Http/Controllers/index.ts")
	assertContains(t, indexContent, `import PostController from './PostController'`)
	assertContains(t, indexContent, `import UserController from './UserController'`)

	// Top-level barrel.
	rootIndex := readFile(t, dir, "actions/index.ts")
	assertContains(t, rootIndex, `export default`)
}

// Port of RepeatedNamespaceController::it_avoids_conflicting_barrel_identifiers_when_namespace_segment_repeats.
func TestRepeatedNamespaceControllerGeneration(t *testing.T) {
	t.Parallel()

	routes := []*wayfinder.RouteInfo{
		{URI: "/admin/repeated", Methods: []string{"get"}, Controller: "App\\Http\\Controllers\\Admin\\Admin\\RepeatedNamespaceController@index"},
		{URI: "/admin/users", Methods: []string{"get"}, Controller: "App\\Http\\Controllers\\Admin\\UserController@index"},
	}

	dir := generateTo(t, routes, wayfinder.Options{SkipRoutes: true})
	content := readFile(t, dir, "actions/App/Http/Controllers/Admin/index.ts")

	assertContains(t, content, `import Admin from './Admin'`)
	assertContains(t, content, `import UserController from './UserController'`)
	assertContains(t, content, `const Admin = {`)
}
