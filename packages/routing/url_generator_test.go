package routing_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/bedrock/packages/routing"
)

func newGenerator() (*routing.UrlGenerator, *routing.Registry) {
	reg := routing.NewRegistry()
	gen := routing.NewUrlGenerator(reg, "https://example.com", []byte("secret-key"))

	return gen, reg
}

func TestUrlGeneratorTo(t *testing.T) {
	t.Parallel()

	gen, _ := newGenerator()

	if u := gen.To("/users"); u != "https://example.com/users" {
		t.Fatalf("expected 'https://example.com/users', got %q", u)
	}
}

func TestUrlGeneratorToWithQuery(t *testing.T) {
	t.Parallel()

	gen, _ := newGenerator()
	q := url.Values{"page": {"2"}, "sort": {"name"}}
	u := gen.To("/users", q)

	if u != "https://example.com/users?page=2&sort=name" {
		t.Fatalf("unexpected URL: %q", u)
	}
}

func TestUrlGeneratorToFullURL(t *testing.T) {
	t.Parallel()

	gen, _ := newGenerator()

	if u := gen.To("https://other.com/path"); u != "https://other.com/path" {
		t.Fatalf("expected full URL passthrough, got %q", u)
	}
}

func TestUrlGeneratorSecure(t *testing.T) {
	t.Parallel()

	gen := routing.NewUrlGenerator(routing.NewRegistry(), "http://example.com", nil)

	if u := gen.Secure("/users"); u != "https://example.com/users" {
		t.Fatalf("expected HTTPS URL, got %q", u)
	}
}

func TestUrlGeneratorRoute(t *testing.T) {
	t.Parallel()

	gen, reg := newGenerator()
	reg.Add("users.show", "GET", "/users/{id}")

	u := gen.Route("users.show", map[string]string{"id": "42"})

	if u != "https://example.com/users/42" {
		t.Fatalf("expected 'https://example.com/users/42', got %q", u)
	}
}

func TestUrlGeneratorRouteWithQuery(t *testing.T) {
	t.Parallel()

	gen, reg := newGenerator()
	reg.Add("users.index", "GET", "/users")

	q := url.Values{"page": {"3"}}
	u := gen.Route("users.index", nil, q)

	if u != "https://example.com/users?page=3" {
		t.Fatalf("unexpected URL: %q", u)
	}
}

func TestUrlGeneratorRouteWithDefaults(t *testing.T) {
	t.Parallel()

	gen, reg := newGenerator()
	reg.Add("users.show", "GET", "/users/{id}")
	gen.SetDefaults(map[string]string{"id": "1"})

	u := gen.Route("users.show", nil)

	if u != "https://example.com/users/1" {
		t.Fatalf("expected default param, got %q", u)
	}
}

func TestUrlGeneratorRouteParamsOverrideDefaults(t *testing.T) {
	t.Parallel()

	gen, reg := newGenerator()
	reg.Add("users.show", "GET", "/users/{id}")
	gen.SetDefaults(map[string]string{"id": "1"})

	u := gen.Route("users.show", map[string]string{"id": "42"})

	if u != "https://example.com/users/42" {
		t.Fatalf("expected param to override default, got %q", u)
	}
}

func TestUrlGeneratorRouteUnknown(t *testing.T) {
	t.Parallel()

	gen, _ := newGenerator()
	u := gen.Route("nonexistent", nil)

	if u == "" || !contains_(u, "unknown") {
		t.Fatalf("expected fallback for unknown route, got %q", u)
	}
}

func TestUrlGeneratorCurrent(t *testing.T) {
	t.Parallel()

	gen, _ := newGenerator()
	req := httptest.NewRequest(http.MethodGet, "/users/42", nil)
	gen.SetRequest(req)

	if u := gen.Current(); u != "https://example.com/users/42" {
		t.Fatalf("expected current URL, got %q", u)
	}
}

func TestUrlGeneratorCurrentNoRequest(t *testing.T) {
	t.Parallel()

	gen, _ := newGenerator()

	if u := gen.Current(); u != "https://example.com/" {
		t.Fatalf("expected root URL, got %q", u)
	}
}

func TestUrlGeneratorPrevious(t *testing.T) {
	t.Parallel()

	gen, _ := newGenerator()
	req := httptest.NewRequest(http.MethodGet, "/current", nil)
	req.Header.Set("Referer", "https://example.com/previous")
	gen.SetRequest(req)

	if u := gen.Previous(); u != "https://example.com/previous" {
		t.Fatalf("expected referer URL, got %q", u)
	}
}

func TestUrlGeneratorPreviousFallback(t *testing.T) {
	t.Parallel()

	gen, _ := newGenerator()
	req := httptest.NewRequest(http.MethodGet, "/current", nil)
	gen.SetRequest(req)

	if u := gen.Previous("/fallback"); u != "/fallback" {
		t.Fatalf("expected fallback URL, got %q", u)
	}
}

func TestUrlGeneratorPreviousNoRequest(t *testing.T) {
	t.Parallel()

	gen, _ := newGenerator()

	if u := gen.Previous(); u != "https://example.com/" {
		t.Fatalf("expected root URL as fallback, got %q", u)
	}
}

func TestUrlGeneratorPreviousPath(t *testing.T) {
	t.Parallel()

	gen, _ := newGenerator()
	req := httptest.NewRequest(http.MethodGet, "/current", nil)
	req.Header.Set("Referer", "https://example.com/previous/page?q=1")
	gen.SetRequest(req)

	if p := gen.PreviousPath(); p != "/previous/page" {
		t.Fatalf("expected '/previous/page', got %q", p)
	}
}

func TestUrlGeneratorAsset(t *testing.T) {
	t.Parallel()

	gen, _ := newGenerator()

	if u := gen.Asset("css/app.css"); u != "https://example.com/css/app.css" {
		t.Fatalf("expected asset URL, got %q", u)
	}
}

func TestUrlGeneratorSetRootURL(t *testing.T) {
	t.Parallel()

	gen, _ := newGenerator()
	gen.SetRootURL("https://new.example.com")

	if u := gen.To("/test"); u != "https://new.example.com/test" {
		t.Fatalf("expected new root URL, got %q", u)
	}
}

func TestSignedRoute(t *testing.T) {
	t.Parallel()

	gen, reg := newGenerator()
	reg.Add("download", "GET", "/files/{id}")

	signed, err := gen.SignedRoute("download", map[string]string{"id": "42"})

	if err != nil {
		t.Fatal(err)
	}

	parsed, _ := url.Parse(signed)

	if parsed.Query().Get("signature") == "" {
		t.Fatal("expected signature in URL")
	}

	req := httptest.NewRequest(http.MethodGet, signed, nil)
	req.URL = parsed

	if !gen.HasValidSignature(req) {
		t.Fatal("expected valid signature")
	}
}

func TestSignedRouteWithExpiration(t *testing.T) {
	t.Parallel()

	gen, reg := newGenerator()
	reg.Add("download", "GET", "/files/{id}")

	signed, err := gen.SignedRoute("download", map[string]string{"id": "42"}, 1*time.Hour)

	if err != nil {
		t.Fatal(err)
	}

	parsed, _ := url.Parse(signed)

	if parsed.Query().Get("expires") == "" {
		t.Fatal("expected expires in URL")
	}

	req := httptest.NewRequest(http.MethodGet, signed, nil)
	req.URL = parsed

	if !gen.HasValidSignature(req) {
		t.Fatal("expected valid signature with expiration")
	}
}

func TestTemporarySignedRoute(t *testing.T) {
	t.Parallel()

	gen, reg := newGenerator()
	reg.Add("invite", "GET", "/invite/{token}")

	signed, err := gen.TemporarySignedRoute("invite", 30*time.Minute, map[string]string{"token": "abc"})

	if err != nil {
		t.Fatal(err)
	}

	parsed, _ := url.Parse(signed)
	req := httptest.NewRequest(http.MethodGet, signed, nil)
	req.URL = parsed

	if !gen.HasValidSignature(req) {
		t.Fatal("expected valid temporary signature")
	}
}

func TestSignedRouteExpired(t *testing.T) {
	t.Parallel()

	gen, reg := newGenerator()
	reg.Add("link", "GET", "/link/{id}")

	signed, _ := gen.SignedRoute("link", map[string]string{"id": "1"}, -1*time.Hour)
	parsed, _ := url.Parse(signed)
	req := httptest.NewRequest(http.MethodGet, signed, nil)
	req.URL = parsed

	if gen.HasValidSignature(req) {
		t.Fatal("expected expired signature to be invalid")
	}

	if !gen.HasCorrectSignature(req) {
		t.Fatal("expected correct signature despite expiration")
	}

	if gen.SignatureHasNotExpired(req) {
		t.Fatal("expected signature to be expired")
	}
}

func TestInvalidSignature(t *testing.T) {
	t.Parallel()

	gen, reg := newGenerator()
	reg.Add("link", "GET", "/link/{id}")

	signed, _ := gen.SignedRoute("link", map[string]string{"id": "1"})
	parsed, _ := url.Parse(signed)

	q := parsed.Query()
	q.Set("signature", "tampered")
	parsed.RawQuery = q.Encode()

	req := httptest.NewRequest(http.MethodGet, parsed.String(), nil)
	req.URL = parsed

	if gen.HasCorrectSignature(req) {
		t.Fatal("expected tampered signature to be invalid")
	}
}

func TestNoSignature(t *testing.T) {
	t.Parallel()

	gen, _ := newGenerator()
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	if gen.HasCorrectSignature(req) {
		t.Fatal("expected no signature to be invalid")
	}
}

func TestSignatureWithNoExpires(t *testing.T) {
	t.Parallel()

	gen, reg := newGenerator()
	reg.Add("link", "GET", "/link/{id}")

	signed, _ := gen.SignedRoute("link", map[string]string{"id": "1"})
	parsed, _ := url.Parse(signed)
	req := httptest.NewRequest(http.MethodGet, signed, nil)
	req.URL = parsed

	if !gen.SignatureHasNotExpired(req) {
		t.Fatal("expected no-expiry signature to be valid")
	}
}

func TestUrlGeneratorMultipleParams(t *testing.T) {
	t.Parallel()

	gen, reg := newGenerator()
	reg.Add("posts.comments.show", "GET", "/posts/{postId}/comments/{commentId}")

	u := gen.Route("posts.comments.show", map[string]string{"postId": "5", "commentId": "3"})

	if u != "https://example.com/posts/5/comments/3" {
		t.Fatalf("unexpected URL: %q", u)
	}
}

func TestUrlGeneratorToStripsLeadingSlash(t *testing.T) {
	t.Parallel()

	gen, _ := newGenerator()

	u1 := gen.To("/path")
	u2 := gen.To("path")

	if u1 != u2 {
		t.Fatalf("expected same URL, got %q vs %q", u1, u2)
	}
}

func contains_(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsStr(s, substr))
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}
