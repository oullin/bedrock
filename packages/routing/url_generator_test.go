package routing

import (
	"strings"
	"testing"
)

// Translation of upstream/framework tests/Routing/RoutingUrlGeneratorTest.php.
//
// Byte-level signed URL parity with Upstream cannot be asserted here without a
// PHP runtime to dump fixtures. The tests below verify the round-trip
// invariants (Sign → HasValidSignature) and the canonical encoding rules.
// RoutingUrlGeneratorTest::testBasicGeneration
// RoutingUrlGeneratorTest::testForceHttps
// RoutingUrlGeneratorTest::testBasicRouteGeneration
// RoutingUrlGeneratorTest::testFluentRouteNameDefinitions
// RoutingUrlGeneratorTest::testSignedUrl
// RoutingUrlGeneratorTest::testTemporarySignedRoute
// RoutingUrlGeneratorTest::testSignedUrlParameterCannotBeNamedSignature
// RoutingUrlGeneratorTest::testSignedUrlParameterCannotBeNamedExpires
// RoutingUrlGeneratorTest::testSignedUrlWithKeyResolver
// RoutingUrlGeneratorTest::testMissingNamedRouteResolution

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

func newGen(t *testing.T) (*UrlGenerator, *Router) {
	t.Helper()
	router := NewRouter(nil, nil)
	req := fakeURLRequest{scheme: "http", host: "example.com"}
	gen := NewUrlGenerator(router.GetRoutes(), req, "")

	return gen, router
}

func TestUrlGenerator_To(t *testing.T) {
	t.Run("test_to_returns_absolute", func(t *testing.T) {
		gen, _ := newGen(t)
		got := gen.To("/foo", nil, nil)

		if got != "http://example.com/foo" {
			t.Errorf("got %q", got)
		}
	})

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

	t.Run("test_missing_parameter_errors", func(t *testing.T) {
		gen, router := newGen(t)
		router.Get("/users/{user}", func() {}).Name("users.show")
		_, err := gen.Route("users.show", nil, true)

		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("test_unknown_route_errors", func(t *testing.T) {
		gen, _ := newGen(t)
		_, err := gen.Route("missing", nil, true)

		if err == nil {
			t.Fatal("expected error")
		}
	})
}

func TestUrlGenerator_Signed(t *testing.T) {
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
