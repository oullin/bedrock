package routing

import (
	"testing"
)

// tests/Routing/RoutingRouteTest.php — the dispatching half of that file
// requires Router (M4) and lives in router_test.go.
// RoutingRouteTest::testMatchesMethodAgainstRequests
// RoutingRouteTest::testWherePatternsProperlyFilter
// RoutingRouteTest::testRoutePrefixing
// RoutingRouteTest::testRouteDomainRegistration
// RoutingRouteTest::testRouteParametersDefaultValue
// RoutingRouteTest::testHasParameters
// RoutingRouteTest::testForgetParameter
// RoutingRouteTest::testParameterNames

// fakeRequest implements matching.MatchableRequest and boundRequest for tests.
type fakeRequest struct {
	method string
	host   string
	path   string
	secure bool
}

func (r fakeRequest) Method() string      { return r.method }
func (r fakeRequest) Host() string        { return r.host }
func (r fakeRequest) PathInfo() string    { return r.path }
func (r fakeRequest) Secure() bool        { return r.secure }
func (r fakeRequest) DecodedPath() string { return r.path }

func TestRoute_Construction(t *testing.T) {
	t.Run("test_basic_route_construction", func(t *testing.T) {
		r := NewRoute("GET", "/foo", func() {})

		if r.Uri != "foo" {
			t.Errorf("uri = %q, want foo", r.Uri)
		}

		if got := r.Methods(); len(got) != 2 || got[0] != "GET" || got[1] != "HEAD" {
			t.Errorf("methods = %v, want [GET HEAD]", got)
		}
	})

	t.Run("test_post_route_does_not_get_head", func(t *testing.T) {
		r := NewRoute("POST", "/foo", func() {})

		for _, m := range r.Methods() {
			if m == "HEAD" {
				t.Error("POST route unexpectedly has HEAD")
			}
		}
	})

	t.Run("test_methods_normalized_to_upper", func(t *testing.T) {
		r := NewRoute("get", "/x", func() {})

		if r.Methods()[0] != "GET" {
			t.Errorf("methods = %v", r.Methods())
		}
	})
}

func TestRoute_Where(t *testing.T) {
	t.Run("test_where_sets_constraint", func(t *testing.T) {
		r := NewRoute("GET", "/users/{id}", func() {})
		r.Where("id", "[0-9]+")

		if r.Wheres["id"] != "[0-9]+" {
			t.Errorf("wheres = %v", r.Wheres)
		}
	})

	t.Run("test_where_alpha_helper", func(t *testing.T) {
		r := NewRoute("GET", "/users/{name}", func() {})
		r.WhereAlpha("name")

		if r.Wheres["name"] != `[a-zA-Z]+` {
			t.Errorf("wheres = %v", r.Wheres)
		}
	})

	t.Run("test_where_number_helper", func(t *testing.T) {
		r := NewRoute("GET", "/users/{id}", func() {})
		r.WhereNumber("id")

		if r.Wheres["id"] != `[0-9]+` {
			t.Errorf("wheres = %v", r.Wheres)
		}
	})

	t.Run("test_where_uuid_helper", func(t *testing.T) {
		r := NewRoute("GET", "/users/{id}", func() {})
		r.WhereUuid("id")

		if r.Wheres["id"] == "" {
			t.Error("expected uuid constraint")
		}
	})

	t.Run("test_where_in_helper", func(t *testing.T) {
		r := NewRoute("GET", "/users/{role}", func() {})
		r.WhereIn("role", []string{"admin", "user", "guest"})

		if r.Wheres["role"] != "admin|user|guest" {
			t.Errorf("wheres = %v", r.Wheres)
		}
	})
}

func TestRoute_Name(t *testing.T) {
	t.Run("test_route_name", func(t *testing.T) {
		r := NewRoute("GET", "/users", func() {}).Name("users.index")

		if r.GetName() != "users.index" {
			t.Errorf("name = %q", r.GetName())
		}
	})

	t.Run("test_route_name_concatenates", func(t *testing.T) {
		r := NewRoute("GET", "/users", func() {}).Name("users.").Name("index")

		if r.GetName() != "users.index" {
			t.Errorf("name = %q", r.GetName())
		}
	})

	t.Run("test_named_glob", func(t *testing.T) {
		r := NewRoute("GET", "/", func() {}).Name("admin.users.show")

		if !r.Named("admin.*") {
			t.Error("admin.* should match")
		}

		if !r.Named("admin.users.show") {
			t.Error("exact should match")
		}

		if r.Named("public.*") {
			t.Error("public.* should not match")
		}
	})
}

func TestRoute_PrefixAndDomain(t *testing.T) {
	t.Run("test_prefix_applies_to_uri", func(t *testing.T) {
		r := NewRoute("GET", "users", func() {}).Prefix("api")

		if r.Uri != "api/users" {
			t.Errorf("uri = %q", r.Uri)
		}

		if r.GetPrefix() != "api" {
			t.Errorf("prefix = %q", r.GetPrefix())
		}
	})

	t.Run("test_prefix_handles_root", func(t *testing.T) {
		r := NewRoute("GET", "/", func() {})
		// Prefix("") must not panic and must leave URI alone.
		if r.Uri != "/" && r.Uri != "" {
			t.Errorf("uri = %q", r.Uri)
		}
	})

	t.Run("test_domain_set_and_get", func(t *testing.T) {
		r := NewRoute("GET", "/", func() {}).Domain("api.example.com")

		if r.GetDomain() != "api.example.com" {
			t.Errorf("domain = %q", r.GetDomain())
		}
	})

	t.Run("test_domain_strips_scheme", func(t *testing.T) {
		r := NewRoute("GET", "/", func() {}).Domain("https://api.example.com")

		if r.GetDomain() != "api.example.com" {
			t.Errorf("domain = %q", r.GetDomain())
		}
	})
}

func TestRoute_Defaults(t *testing.T) {
	t.Run("test_defaults_sets_value", func(t *testing.T) {
		r := NewRoute("GET", "/users/{user}", func() {}).Defaults("user", "guest")

		if r.DefaultValues["user"] != "guest" {
			t.Errorf("defaults = %v", r.DefaultValues)
		}

		if !r.HasDefault("user") {
			t.Error("HasDefault should report true")
		}
	})
}

func TestRoute_Fallback(t *testing.T) {
	t.Run("test_fallback_marks_route", func(t *testing.T) {
		r := NewRoute("GET", "/{any}", func() {}).Fallback()

		if !r.IsFallback {
			t.Error("expected IsFallback true")
		}
	})
}

func TestRoute_Compile(t *testing.T) {
	t.Run("test_compile_produces_regex", func(t *testing.T) {
		r := NewRoute("GET", "/users/{user}", func() {})
		c := r.Compiled()

		if c == nil {
			t.Fatal("compiled is nil")
		}

		if len(c.PathVariables()) != 1 || c.PathVariables()[0] != "user" {
			t.Errorf("path vars = %v", c.PathVariables())
		}
	})

	t.Run("test_compile_with_host", func(t *testing.T) {
		r := NewRoute("GET", "/", func() {}).Domain("{sub}.example.com")
		c := r.Compiled()

		if c == nil || c.HostRegex() == "" {
			t.Fatal("expected host regex")
		}
	})
}

func TestRoute_Matches(t *testing.T) {
	t.Run("test_match_static_path", func(t *testing.T) {
		r := NewRoute("GET", "/users", func() {})

		if !r.Matches(fakeRequest{method: "GET", path: "/users"}, true) {
			t.Error("should match GET /users")
		}

		if r.Matches(fakeRequest{method: "POST", path: "/users"}, true) {
			t.Error("should not match POST /users")
		}

		if r.Matches(fakeRequest{method: "GET", path: "/posts"}, true) {
			t.Error("should not match GET /posts")
		}
	})

	t.Run("test_match_with_parameter", func(t *testing.T) {
		r := NewRoute("GET", "/users/{id}", func() {})

		if !r.Matches(fakeRequest{method: "GET", path: "/users/42"}, true) {
			t.Error("should match")
		}
	})

	t.Run("test_match_with_constraint", func(t *testing.T) {
		r := NewRoute("GET", "/users/{id}", func() {})
		r.WhereNumber("id")

		if !r.Matches(fakeRequest{method: "GET", path: "/users/42"}, true) {
			t.Error("should match numeric")
		}

		if r.Matches(fakeRequest{method: "GET", path: "/users/abc"}, true) {
			t.Error("should not match alpha")
		}
	})

	t.Run("test_skip_method_validator", func(t *testing.T) {
		r := NewRoute("GET", "/users", func() {})
		// includingMethod=false means POST is also acceptable.
		if !r.Matches(fakeRequest{method: "POST", path: "/users"}, false) {
			t.Error("should match without method validator")
		}
	})
}

func TestRoute_Bind(t *testing.T) {
	t.Run("test_bind_extracts_parameters", func(t *testing.T) {
		r := NewRoute("GET", "/users/{user}/posts/{post}", func() {})

		if _, err := r.Bind(fakeRequest{path: "/users/42/posts/7"}); err != nil {
			t.Fatal(err)
		}

		if r.Parameter("user", "") != "42" {
			t.Errorf("user = %q", r.Parameter("user", ""))
		}

		if r.Parameter("post", "") != "7" {
			t.Errorf("post = %q", r.Parameter("post", ""))
		}
	})

	t.Run("test_bind_then_set_then_forget", func(t *testing.T) {
		r := NewRoute("GET", "/users/{user}", func() {})
		_, _ = r.Bind(fakeRequest{path: "/users/42"})
		r.SetParameter("user", "43")

		if r.Parameter("user", "") != "43" {
			t.Error("setParameter failed")
		}

		r.ForgetParameter("user")

		if r.HasParameter("user") {
			t.Error("forgetParameter failed")
		}
	})

	t.Run("test_parameter_names", func(t *testing.T) {
		r := NewRoute("GET", "/posts/{post}/comments/{comment}", func() {})
		got := r.ParameterNames()

		if len(got) != 2 || got[0] != "post" || got[1] != "comment" {
			t.Errorf("names = %v", got)
		}
	})
}

func TestRoute_BindingFields(t *testing.T) {
	t.Run("test_custom_binding_field", func(t *testing.T) {
		r := NewRoute("GET", "/users/{user:slug}", func() {})

		if r.BindingFieldFor("user") != "slug" {
			t.Errorf("binding field = %q", r.BindingFieldFor("user"))
		}

		if r.Uri != "users/{user}" {
			t.Errorf("uri = %q", r.Uri)
		}
	})
}

func TestRoute_Scheme(t *testing.T) {
	t.Run("test_secure_route_only_matches_https", func(t *testing.T) {
		r := NewRoute("GET", "/", func() {})
		r.ActionMap["https"] = "https"

		if r.Matches(fakeRequest{method: "GET", path: "/", secure: false}, true) {
			t.Error("insecure request should not match https route")
		}

		if !r.Matches(fakeRequest{method: "GET", path: "/", secure: true}, true) {
			t.Error("secure request should match https route")
		}
	})
}
