package routing

import (
	"strings"
	"testing"

	"github.com/bedrock/packages/routing/contracts"
)

// Translation of upstream/framework tests/Routing/RoutingUrlGeneratorTest.php.
//
// Byte-level signed URL parity with Upstream cannot be asserted here without a
// PHP runtime to dump fixtures. The tests below verify the round-trip
// invariants (Sign → HasValidSignature) and the canonical encoding rules.
// RoutingUrlGeneratorTest::testBasicGeneration
// RoutingUrlGeneratorTest::testQueryGeneration
// RoutingUrlGeneratorTest::testAssetGeneration
// RoutingUrlGeneratorTest::testBasicGenerationWithPathFormatting
// RoutingUrlGeneratorTest::testUseRootUrl
// RoutingUrlGeneratorTest::testForceHttps
// RoutingUrlGeneratorTest::testBasicRouteGeneration
// RoutingUrlGeneratorTest::testFluentRouteNameDefinitions
// RoutingUrlGeneratorTest::testRoutableInterfaceRouting
// RoutingUrlGeneratorTest::testRoutableInterfaceRoutingWithCustomBindingField
// RoutingUrlGeneratorTest::testRoutableInterfaceRoutingAsQueryString
// RoutingUrlGeneratorTest::testRoutableInterfaceRoutingWithSingleParameter
// RoutingUrlGeneratorTest::testRoutesMaintainRequestScheme
// RoutingUrlGeneratorTest::testRoutesWithDomains
// RoutingUrlGeneratorTest::testRoutesWithDomainsAndPorts
// RoutingUrlGeneratorTest::testRoutesWithDomainsStripsProtocols
// RoutingUrlGeneratorTest::testUrlGenerationForControllersRequiresPassingOfRequiredParameters
// RoutingUrlGeneratorTest::testUrlGenerationThrowsExceptionForMissingParametersWithMeaningfulMessage
// RoutingUrlGeneratorTest::testSetAssetUrl
// RoutingUrlGeneratorTest::testSignedUrl
// RoutingUrlGeneratorTest::testSignedRelativeUrl
// RoutingUrlGeneratorTest::testTemporarySignedRoute
// RoutingUrlGeneratorTest::testSignedUrlParameterCannotBeNamedSignature
// RoutingUrlGeneratorTest::testSignedUrlParameterCannotBeNamedExpires
// RoutingUrlGeneratorTest::testRouteGenerationWithBackedEnums
// RoutingUrlGeneratorTest::testRouteGenerationWithNestedBackedEnums
// RoutingUrlGeneratorTest::testSignedUrlWithKeyResolver
// RoutingUrlGeneratorTest::testRouteNotDefinedException
// RoutingUrlGeneratorTest::testMissingNamedRouteResolution
// RoutingUrlGeneratorTest::testPassedParametersHavePrecedenceOverDefaults
// RoutingUrlGeneratorTest::testComplexRouteGenerationWithDefaultsAndBindingFields
// RoutingUrlGeneratorTest::testDefaultsCanBeCombinedWithExtraQueryParameters
// RoutingUrlGeneratorTest::testUrlGenerationWithOptionalParameters

// fakeURLRequest implements [URLRequest] for tests.
type fakeURLRequest struct {
	scheme string
	host   string
	url    string
	path   string
	query  map[string]string
	qs     string
}

func (r fakeURLRequest) Scheme() string           { return r.scheme }
func (r fakeURLRequest) Host() string             { return r.host }
func (r fakeURLRequest) URL() string              { return r.url }
func (r fakeURLRequest) Path() string             { return r.path }
func (r fakeURLRequest) Query(name string) string { return r.query[name] }
func (r fakeURLRequest) QueryString() string      { return r.qs }

type urlGeneratorRoutable struct {
	key any
}

func (u urlGeneratorRoutable) GetRouteKey() any                                     { return u.key }
func (u urlGeneratorRoutable) GetRouteKeyName() string                              { return "id" }
func (u urlGeneratorRoutable) ResolveRouteBinding(value, field string) (any, error) { return nil, nil }
func (u urlGeneratorRoutable) ResolveChildRouteBinding(_, _, _ string) (any, error) { return nil, nil }

var _ contracts.UrlRoutable = urlGeneratorRoutable{}

type urlGeneratorStatus string

func (s urlGeneratorStatus) BackingValue() string { return string(s) }

func newGen(t *testing.T) (*UrlGenerator, *Router) {
	t.Helper()
	router := NewRouter(nil, nil)
	req := fakeURLRequest{scheme: "http", host: "example.com"}
	gen := NewUrlGenerator(router.GetRoutes(), req, "")

	return gen, router
}

func TestUrlGenerator_To(t *testing.T) {
	// RoutingUrlGeneratorTest::testBasicGeneration
	t.Run("test_to_returns_absolute", func(t *testing.T) {
		gen, _ := newGen(t)
		got := gen.To("/foo", nil, nil)

		if got != "http://example.com/foo" {
			t.Errorf("got %q", got)
		}
	})

	// RoutingUrlGeneratorTest::testBasicGenerationWithPathFormatting
	t.Run("test_to_escapes_extra_path_segments", func(t *testing.T) {
		gen, _ := newGen(t)
		got := gen.To("/users", []string{"Taylor Otwell"}, nil)

		if got != "http://example.com/users/Taylor%20Otwell" {
			t.Errorf("got %q", got)
		}
	})

	// RoutingUrlGeneratorTest::testForceHttps
	t.Run("test_to_secure_forces_https", func(t *testing.T) {
		gen, _ := newGen(t)
		got := gen.Secure("/foo", nil)

		if !strings.HasPrefix(got, "https://") {
			t.Errorf("got %q, want https prefix", got)
		}
	})

	t.Run("test_to_passthrough_absolute", func(t *testing.T) {
		gen, _ := newGen(t)
		got := gen.To("https://other.example/foo", nil, nil)

		if got != "https://other.example/foo" {
			t.Errorf("got %q", got)
		}
	})
}

func TestUrlGenerator_Route(t *testing.T) {
	// RoutingUrlGeneratorTest::testBasicRouteGeneration
	t.Run("test_named_route_url", func(t *testing.T) {
		gen, router := newGen(t)
		router.Get("/users/{user}", func() {}).Name("users.show")
		got, err := gen.Route("users.show", map[string]any{"user": "alice"}, true)

		if err != nil {
			t.Fatal(err)
		}

		if got != "http://example.com/users/alice" {
			t.Errorf("got %q", got)
		}
	})

	t.Run("test_named_route_relative", func(t *testing.T) {
		gen, router := newGen(t)
		router.Get("/users/{user}", func() {}).Name("users.show")
		got, err := gen.Route("users.show", map[string]any{"user": "alice"}, false)

		if err != nil {
			t.Fatal(err)
		}

		if got != "/users/alice" {
			t.Errorf("got %q", got)
		}
	})

	// RoutingUrlGeneratorTest::testRoutesMaintainRequestScheme
	t.Run("test_named_route_uses_request_scheme", func(t *testing.T) {
		router := NewRouter(nil, nil)
		req := fakeURLRequest{scheme: "https", host: "example.com"}
		gen := NewUrlGenerator(router.GetRoutes(), req, "")
		router.Get("/users", func() {}).Name("users.index")
		got, err := gen.Route("users.index", nil, true)

		if err != nil {
			t.Fatal(err)
		}

		if got != "https://example.com/users" {
			t.Errorf("got %q", got)
		}
	})

	// RoutingUrlGeneratorTest::testQueryGeneration
	t.Run("test_named_route_extra_params_become_query", func(t *testing.T) {
		gen, router := newGen(t)
		router.Get("/search", func() {}).Name("search")
		got, err := gen.Route("search", map[string]any{"q": "go", "page": 2}, true)

		if err != nil {
			t.Fatal(err)
		}
		// Sorted: page=2&q=go
		if !strings.HasSuffix(got, "/search?page=2&q=go") {
			t.Errorf("got %q", got)
		}
	})

	// RoutingUrlGeneratorTest::testAssetGeneration
	t.Run("test_asset_generation_uses_asset_root", func(t *testing.T) {
		gen, _ := newGen(t)
		gen.assetRoot = "https://cdn.example/assets"
		got := gen.Asset("images/logo.svg", nil)

		if got != "https://cdn.example/assets/images/logo.svg" {
			t.Errorf("got %q", got)
		}
	})

	// RoutingUrlGeneratorTest::testSetAssetUrl
	t.Run("test_constructor_asset_root_sets_asset_url", func(t *testing.T) {
		router := NewRouter(nil, nil)
		gen := NewUrlGenerator(router.GetRoutes(), fakeURLRequest{scheme: "http", host: "example.com"}, "https://assets.example")
		got := gen.Asset("/css/app.css", nil)

		if got != "https://assets.example/css/app.css" {
			t.Errorf("got %q", got)
		}
	})

	// RoutingUrlGeneratorTest::testUseRootUrl
	t.Run("test_force_root_url_overrides_host", func(t *testing.T) {
		gen, _ := newGen(t)
		gen.ForceRootUrl("cdn.example/base")
		gen.ForceHttps(true)
		got := gen.To("/foo", nil, nil)

		if got != "https://cdn.example/base/foo" {
			t.Errorf("got %q", got)
		}
	})

	// RoutingUrlGeneratorTest::testUrlGenerationForControllersRequiresPassingOfRequiredParameters
	// RoutingUrlGeneratorTest::testUrlGenerationThrowsExceptionForMissingParametersWithMeaningfulMessage
	t.Run("test_missing_parameter_errors", func(t *testing.T) {
		gen, router := newGen(t)
		router.Get("/users/{user}", func() {}).Name("users.show")
		_, err := gen.Route("users.show", nil, true)

		if err == nil {
			t.Fatal("expected error")
		}

		if !strings.Contains(err.Error(), "users.show") || !strings.Contains(err.Error(), "user") {
			t.Fatalf("err = %v", err)
		}
	})

	// RoutingUrlGeneratorTest::testRouteNotDefinedException
	t.Run("test_unknown_route_errors", func(t *testing.T) {
		gen, _ := newGen(t)
		_, err := gen.Route("missing", nil, true)

		if err == nil {
			t.Fatal("expected error")
		}
	})

	// RoutingUrlGeneratorTest::testRoutesWithDomains
	t.Run("test_named_route_uses_route_domain", func(t *testing.T) {
		gen, router := newGen(t)
		router.Get("/users", func() {}).Domain("admin.example.com").Name("users.index")
		got, err := gen.Route("users.index", nil, true)

		if err != nil {
			t.Fatal(err)
		}

		if got != "http://admin.example.com/users" {
			t.Errorf("got %q", got)
		}
	})

	// RoutingUrlGeneratorTest::testRoutesWithDomainsAndPorts
	t.Run("test_named_route_preserves_route_domain_port", func(t *testing.T) {
		gen, router := newGen(t)
		router.Get("/users", func() {}).Domain("admin.example.com:8443").Name("users.index")
		got, err := gen.Route("users.index", nil, true)

		if err != nil {
			t.Fatal(err)
		}

		if got != "http://admin.example.com:8443/users" {
			t.Errorf("got %q", got)
		}
	})

	// RoutingUrlGeneratorTest::testRoutesWithDomainsStripsProtocols
	t.Run("test_named_route_strips_protocol_from_route_domain", func(t *testing.T) {
		gen, router := newGen(t)
		router.Get("/users", func() {}).Domain("https://admin.example.com").Name("users.index")
		got, err := gen.Route("users.index", nil, true)

		if err != nil {
			t.Fatal(err)
		}

		if got != "http://admin.example.com/users" {
			t.Errorf("got %q", got)
		}
	})

	// RoutingUrlGeneratorTest::testUrlGenerationWithOptionalParameters
	t.Run("test_optional_route_parameters_are_removed_when_missing", func(t *testing.T) {
		gen, router := newGen(t)
		router.Get("/users/{user?}", func() {}).Name("users.optional")
		got, err := gen.Route("users.optional", nil, true)

		if err != nil {
			t.Fatal(err)
		}

		if got != "http://example.com/users" {
			t.Errorf("got %q", got)
		}
	})

	// RoutingUrlGeneratorTest::testPassedParametersHavePrecedenceOverDefaults
	t.Run("test_passed_parameters_override_route_defaults", func(t *testing.T) {
		gen, router := newGen(t)
		router.Get("/{locale}/users", func() {}).Name("localized.users").Defaults("locale", "en")
		got, err := gen.Route("localized.users", map[string]any{"locale": "fr"}, true)

		if err != nil {
			t.Fatal(err)
		}

		if got != "http://example.com/fr/users" {
			t.Errorf("got %q", got)
		}
	})

	// RoutingUrlGeneratorTest::testDefaultsCanBeCombinedWithExtraQueryParameters
	t.Run("test_route_defaults_combine_with_extra_query_parameters", func(t *testing.T) {
		gen, router := newGen(t)
		router.Get("/{locale}/users", func() {}).Name("localized.users").Defaults("locale", "en")
		got, err := gen.Route("localized.users", map[string]any{"page": 2}, true)

		if err != nil {
			t.Fatal(err)
		}

		if got != "http://example.com/en/users?page=2" {
			t.Errorf("got %q", got)
		}
	})

	// RoutingUrlGeneratorTest::testComplexRouteGenerationWithDefaultsAndBindingFields
	t.Run("test_route_defaults_and_binding_fields_generate_url", func(t *testing.T) {
		gen, router := newGen(t)
		router.Get("/{locale}/users/{user:slug}", func() {}).Name("localized.users.show").Defaults("locale", "en")
		got, err := gen.Route("localized.users.show", map[string]any{
			"user": urlGeneratorRoutable{key: "taylor"},
			"tab":  "profile",
		}, true)

		if err != nil {
			t.Fatal(err)
		}

		if got != "http://example.com/en/users/taylor?tab=profile" {
			t.Errorf("got %q", got)
		}
	})

	// RoutingUrlGeneratorTest::testRoutableInterfaceRouting
	// RoutingUrlGeneratorTest::testRoutableInterfaceRoutingWithSingleParameter
	t.Run("test_routable_parameter_uses_route_key", func(t *testing.T) {
		gen, router := newGen(t)
		router.Get("/users/{user}", func() {}).Name("users.show")
		got, err := gen.Route("users.show", map[string]any{"user": urlGeneratorRoutable{key: 42}}, true)

		if err != nil {
			t.Fatal(err)
		}

		if got != "http://example.com/users/42" {
			t.Errorf("got %q", got)
		}
	})

	// RoutingUrlGeneratorTest::testRoutableInterfaceRoutingWithCustomBindingField
	t.Run("test_routable_parameter_with_custom_binding_field_uses_route_key", func(t *testing.T) {
		gen, router := newGen(t)
		router.Get("/users/{user:slug}", func() {}).Name("users.slug")
		got, err := gen.Route("users.slug", map[string]any{"user": urlGeneratorRoutable{key: "taylor"}}, true)

		if err != nil {
			t.Fatal(err)
		}

		if got != "http://example.com/users/taylor" {
			t.Errorf("got %q", got)
		}
	})

	// RoutingUrlGeneratorTest::testRoutableInterfaceRoutingAsQueryString
	t.Run("test_routable_extra_parameter_is_query_string", func(t *testing.T) {
		gen, router := newGen(t)
		router.Get("/search", func() {}).Name("search")
		got, err := gen.Route("search", map[string]any{"user": urlGeneratorRoutable{key: 42}}, true)

		if err != nil {
			t.Fatal(err)
		}

		if got != "http://example.com/search?user=42" {
			t.Errorf("got %q", got)
		}
	})

	// RoutingUrlGeneratorTest::testRouteGenerationWithBackedEnums
	t.Run("test_backed_enum_parameter_uses_backing_value", func(t *testing.T) {
		gen, router := newGen(t)
		router.Get("/statuses/{status}", func() {}).Name("statuses.show")
		got, err := gen.Route("statuses.show", map[string]any{"status": urlGeneratorStatus("published")}, true)

		if err != nil {
			t.Fatal(err)
		}

		if got != "http://example.com/statuses/published" {
			t.Errorf("got %q", got)
		}
	})

	// RoutingUrlGeneratorTest::testRouteGenerationWithNestedBackedEnums
	t.Run("test_nested_backed_enum_parameters_use_backing_values", func(t *testing.T) {
		gen, router := newGen(t)
		router.Get("/teams/{team}/statuses/{status}", func() {}).Name("teams.statuses.show")
		got, err := gen.Route("teams.statuses.show", map[string]any{
			"team":   urlGeneratorStatus("core"),
			"status": urlGeneratorStatus("active"),
		}, true)

		if err != nil {
			t.Fatal(err)
		}

		if got != "http://example.com/teams/core/statuses/active" {
			t.Errorf("got %q", got)
		}
	})
}

func TestUrlGenerator_Signed(t *testing.T) {
	// RoutingUrlGeneratorTest::testSignedUrl
	t.Run("test_signed_route_round_trip", func(t *testing.T) {
		gen, router := newGen(t)
		gen.SetKeyResolver("test-key-12345")
		router.Get("/download/{file}", func() {}).Name("download")
		got, err := gen.SignedRoute("download", map[string]any{"file": "a.zip"}, 0, true)

		if err != nil {
			t.Fatal(err)
		}
		// Should contain a signature query string.
		if !strings.Contains(got, "signature=") {
			t.Errorf("missing signature: %q", got)
		}
		// Round-trip: a request carrying that URL should validate.
		idx := strings.Index(got, "?")
		req := fakeURLRequest{
			scheme: "http",
			host:   "example.com",
			url:    got[:idx],
			path:   "/download/a.zip",
			query:  parseQuery(got[idx+1:]),
			qs:     got[idx+1:],
		}

		if !gen.HasValidSignature(req, true) {
			t.Errorf("signature should validate, got url %q", got)
		}
	})

	// RoutingUrlGeneratorTest::testSignedRelativeUrl
	t.Run("test_signed_relative_route_round_trip", func(t *testing.T) {
		gen, router := newGen(t)
		gen.SetKeyResolver("test-key-12345")
		router.Get("/download/{file}", func() {}).Name("download")
		got, err := gen.SignedRoute("download", map[string]any{"file": "a.zip"}, 0, false)

		if err != nil {
			t.Fatal(err)
		}

		idx := strings.Index(got, "?")
		req := fakeURLRequest{
			scheme: "http",
			host:   "example.com",
			url:    "http://example.com" + got,
			path:   "/download/a.zip",
			query:  parseQuery(got[idx+1:]),
			qs:     got[idx+1:],
		}

		if !gen.HasValidSignature(req, false) {
			t.Errorf("relative signature should validate, got url %q", got)
		}
	})

	// RoutingUrlGeneratorTest::testTemporarySignedRoute
	t.Run("test_temporary_signed_route_has_expires", func(t *testing.T) {
		gen, router := newGen(t)
		gen.SetKeyResolver("k")
		router.Get("/dl/{f}", func() {}).Name("dl")
		got, err := gen.TemporarySignedRoute("dl", 60, map[string]any{"f": "x"}, true)

		if err != nil {
			t.Fatal(err)
		}

		if !strings.Contains(got, "expires=") {
			t.Errorf("missing expires: %q", got)
		}
	})

	// RoutingUrlGeneratorTest::testSignedUrlParameterCannotBeNamedSignature
	// RoutingUrlGeneratorTest::testSignedUrlParameterCannotBeNamedExpires
	t.Run("test_signed_route_rejects_reserved_params", func(t *testing.T) {
		gen, router := newGen(t)
		gen.SetKeyResolver("k")
		router.Get("/x", func() {}).Name("x")

		if _, err := gen.SignedRoute("x", map[string]any{"signature": "x"}, 0, true); err == nil {
			t.Error("expected reserved-param error")
		}

		if _, err := gen.SignedRoute("x", map[string]any{"expires": "1"}, 0, true); err == nil {
			t.Error("expected reserved-param error")
		}
	})

	t.Run("test_invalid_signature_fails", func(t *testing.T) {
		gen, router := newGen(t)
		gen.SetKeyResolver("k")
		router.Get("/dl/{f}", func() {}).Name("dl")
		req := fakeURLRequest{
			scheme: "http",
			host:   "example.com",
			url:    "http://example.com/dl/a",
			path:   "/dl/a",
			query:  map[string]string{"signature": "deadbeef"},
			qs:     "signature=deadbeef",
		}

		if gen.HasValidSignature(req, true) {
			t.Error("bogus signature should not validate")
		}
	})
}

func parseQuery(s string) map[string]string {
	out := map[string]string{}

	for _, pair := range strings.Split(s, "&") {
		if i := strings.Index(pair, "="); i >= 0 {
			out[pair[:i]] = pair[i+1:]
		}
	}

	return out
}
