package routegen

import (
	"testing"
)

// TestSafeMethod verifies that SafeMethod produces the same output as
// TypeScript::safeMethod() from the PHP implementation, covering all edge
// cases exercised by the upstream RouteGen test suite.
func TestSafeMethod(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input  string
		suffix string
		want   string
	}{
		// Reserved keywords get the suffix appended.
		{"delete", "Method", "deleteMethod"},
		{"export", "Method", "exportMethod"},
		{"default", "Method", "defaultMethod"},
		{"class", "Method", "classMethod"},
		{"return", "Method", "returnMethod"},
		{"import", "Method", "importMethod"},
		{"let", "Method", "letMethod"},
		{"const", "Method", "constMethod"},
		{"var", "Method", "varMethod"},
		{"function", "Method", "functionMethod"},
		{"new", "Method", "newMethod"},
		{"null", "Method", "nullMethod"},
		{"true", "Method", "trueMethod"},
		{"false", "Method", "falseMethod"},

		// Names starting with digits get the suffix prepended (lowercased).
		{"404", "Method", "method404"},
		{"2fa", "Method", "method2fa"},
		{"123abc", "Method", "method123abc"},

		// Normal names pass through unchanged.
		{"dashboard", "Method", "dashboard"},
		{"index", "Method", "index"},
		{"show", "Method", "show"},
		{"PostController", "Method", "PostController"},

		// Kebab-case converted to camelCase.
		{"invalid-js-name", "Method", "invalidJsName"},
		{"my-route", "Method", "myRoute"},

		// Non-word characters replaced by underscore then camel.
		{"some method", "Method", "some_method"},
	}

	for _, c := range cases {
		c := c
		t.Run(c.input, func(t *testing.T) {
			t.Parallel()
			got := SafeMethod(c.input, c.suffix)

			if got != c.want {
				t.Errorf("SafeMethod(%q, %q) = %q, want %q", c.input, c.suffix, got, c.want)
			}
		})
	}
}

func TestQuoteIfNeeded(t *testing.T) {
	t.Parallel()

	cases := []struct{ input, want string }{
		{"404", "404"},       // pure integer — no quotes
		{"2fa", `"2fa"`},     // starts with digit → quoted
		{"delete", "delete"}, // doesn't start with digit → no quotes
		{"index", "index"},
	}

	for _, c := range cases {
		c := c
		t.Run(c.input, func(t *testing.T) {
			t.Parallel()
			got := QuoteIfNeeded(c.input)

			if got != c.want {
				t.Errorf("QuoteIfNeeded(%q) = %q, want %q", c.input, got, c.want)
			}
		})
	}
}

// TestNewVerb verifies HTTP verb normalisation.
func TestNewVerb(t *testing.T) {
	t.Parallel()

	cases := []struct {
		method   string
		actual   string
		formSafe string
	}{
		{"GET", "get", "get"},
		{"get", "get", "get"},
		{"HEAD", "head", "get"},
		{"OPTIONS", "options", "get"},
		{"POST", "post", "post"},
		{"PUT", "put", "post"},
		{"PATCH", "patch", "post"},
		{"DELETE", "delete", "post"},
	}

	for _, c := range cases {
		c := c
		t.Run(c.method, func(t *testing.T) {
			t.Parallel()
			v := NewVerb(c.method)

			if v.Actual != c.actual {
				t.Errorf("Verb.Actual = %q, want %q", v.Actual, c.actual)
			}

			if v.FormSafe != c.formSafe {
				t.Errorf("Verb.FormSafe = %q, want %q", v.FormSafe, c.formSafe)
			}
		})
	}
}

// TestRouteInfoHelpers exercises the derived fields on RouteInfo.
func TestRouteInfoHelpers(t *testing.T) {
	t.Parallel()

	r := &RouteInfo{
		URI:        "/posts/{post}",
		Methods:    []string{"get", "head"},
		Name:       "posts.show",
		Controller: "App\\Http\\Controllers\\PostController@show",
	}

	if got := r.ControllerClass(); got != "App\\Http\\Controllers\\PostController" {
		t.Errorf("ControllerClass() = %q", got)
	}

	if got := r.ActionMethod(); got != "show" {
		t.Errorf("ActionMethod() = %q", got)
	}

	if got := r.DotNamespace(); got != "App.Http.Controllers.PostController" {
		t.Errorf("DotNamespace() = %q", got)
	}

	if got := r.JsMethod(); got != "show" {
		t.Errorf("JsMethod() = %q", got)
	}

	if got := r.NamedMethod(); got != "show" {
		t.Errorf("NamedMethod() = %q", got)
	}
}

// TestRouteInfoInvokable exercises invokable controller detection.
func TestRouteInfoInvokable(t *testing.T) {
	t.Parallel()

	r := &RouteInfo{
		URI:         "/invokable-controller",
		Methods:     []string{"get", "head"},
		Controller:  "App\\Http\\Controllers\\InvokableController@Invoke",
		IsInvokable: true,
	}

	if got := r.OriginalJsMethod(); got != "InvokableController" {
		t.Errorf("OriginalJsMethod() = %q, want InvokableController", got)
	}
}

// TestFullURI verifies URI assembly with domain, scheme, and base-path.
func TestFullURI(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		route    RouteInfo
		wantJSON string
	}{
		{
			name:     "simple",
			route:    RouteInfo{URI: "/posts"},
			wantJSON: `"/posts"`,
		},
		{
			name:     "with_base_path",
			route:    RouteInfo{URI: "/posts", BasePath: "/v2"},
			wantJSON: `"/v2/posts"`,
		},
		{
			name:     "with_domain",
			route:    RouteInfo{URI: "/fixed-domain/{param}", Domain: "example.test", Scheme: "//"},
			wantJSON: `"//example.test/fixed-domain/{param}"`,
		},
		{
			name: "with_optional_default",
			route: RouteInfo{
				URI:      "/with-defaults/{locale}",
				Defaults: map[string]string{"locale": "en"},
			},
			wantJSON: `"/with-defaults/{locale?}"`,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			got := c.route.FullURI()

			if got != c.wantJSON {
				t.Errorf("FullURI() = %s, want %s", got, c.wantJSON)
			}
		})
	}
}
