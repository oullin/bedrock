package httpx_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/httpx"
)

func TestNewRequest(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/hello", nil)
	req := httpx.NewRequest(raw)

	if req.Raw() != raw {
		t.Fatal("Raw() should return the underlying request")
	}
}

func TestRequestMethod(t *testing.T) {
	t.Parallel()

	req := httpx.NewRequest(httptest.NewRequest(http.MethodPost, "/", nil))

	if req.Method() != http.MethodPost {
		t.Fatalf("expected POST, got %s", req.Method())
	}

	if !req.IsMethod("post") {
		t.Fatal("IsMethod should be case-insensitive")
	}

	if req.IsMethod("GET") {
		t.Fatal("IsMethod should not match GET for a POST request")
	}
}

func TestRequestURL(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/foo?bar=1", nil)
	req := httpx.NewRequest(raw)

	if req.URL() != "/foo?bar=1" {
		t.Fatalf("expected /foo?bar=1, got %s", req.URL())
	}

	if req.Path() != "/foo" {
		t.Fatalf("expected /foo, got %s", req.Path())
	}

	if req.QueryString() != "bar=1" {
		t.Fatalf("expected bar=1, got %s", req.QueryString())
	}
}

func TestRequestFullURL(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/path?q=1", nil)
	raw.Host = "example.com"
	req := httpx.NewRequest(raw)

	expected := "http://example.com/path?q=1"

	if req.FullURL() != expected {
		t.Fatalf("expected %s, got %s", expected, req.FullURL())
	}
}

func TestRequestFullURLWithForwardedProto(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/path", nil)
	raw.Host = "example.com"
	raw.Header.Set("X-Forwarded-Proto", "https")
	req := httpx.NewRequest(raw)

	if req.FullURL() != "https://example.com/path" {
		t.Fatalf("expected https scheme, got %s", req.FullURL())
	}
}

func TestRequestSegments(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/foo/bar/baz", nil)
	req := httpx.NewRequest(raw)

	segments := req.Segments()

	if len(segments) != 3 {
		t.Fatalf("expected 3 segments, got %d", len(segments))
	}

	if segments[0] != "foo" || segments[1] != "bar" || segments[2] != "baz" {
		t.Fatalf("unexpected segments: %v", segments)
	}

	if req.Segment(2) != "bar" {
		t.Fatalf("expected bar, got %s", req.Segment(2))
	}

	if req.Segment(10, "default") != "default" {
		t.Fatal("out-of-range segment should return fallback")
	}
}

func TestRequestIs(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/admin/users/123", nil)
	req := httpx.NewRequest(raw)

	if !req.Is("/admin/*") {
		t.Fatal("expected /admin/* to match")
	}

	if req.Is("/blog/*") {
		t.Fatal("expected /blog/* to not match")
	}
}

func TestRequestIP(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.RemoteAddr = "192.168.1.1:1234"
	req := httpx.NewRequest(raw)

	if req.IP() != "192.168.1.1" {
		t.Fatalf("expected 192.168.1.1, got %s", req.IP())
	}
}

func TestRequestIPFromForwardedFor(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.Header.Set("X-Forwarded-For", "10.0.0.1, 10.0.0.2")
	req := httpx.NewRequest(raw)

	if req.IP() != "10.0.0.1" {
		t.Fatalf("expected 10.0.0.1, got %s", req.IP())
	}

	ips := req.IPs()

	if len(ips) != 2 {
		t.Fatalf("expected 2 IPs, got %d", len(ips))
	}
}

func TestRequestUserAgent(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.Header.Set("User-Agent", "TestBot/1.0")
	req := httpx.NewRequest(raw)

	if req.UserAgent() != "TestBot/1.0" {
		t.Fatalf("expected TestBot/1.0, got %s", req.UserAgent())
	}
}

func TestRequestSecure(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	req := httpx.NewRequest(raw)

	if req.Secure() {
		t.Fatal("expected non-secure for plain HTTP")
	}

	raw.Header.Set("X-Forwarded-Proto", "https")

	if !req.Secure() {
		t.Fatal("expected secure with X-Forwarded-Proto: https")
	}
}

func TestRequestSchemeAndHost(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.Host = "example.com"
	req := httpx.NewRequest(raw)

	if req.SchemeAndHost() != "http://example.com" {
		t.Fatalf("expected http://example.com, got %s", req.SchemeAndHost())
	}
}

func TestRequestHost(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/", nil)
	raw.Host = "example.com:8080"
	req := httpx.NewRequest(raw)

	if req.Host() != "example.com:8080" {
		t.Fatalf("expected example.com:8080, got %s", req.Host())
	}
}

func TestRequestFingerprint(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/test", nil)
	raw.Host = "example.com"
	raw.RemoteAddr = "1.2.3.4:5678"
	raw.Header.Set("User-Agent", "bot")
	req := httpx.NewRequest(raw)

	fp := req.Fingerprint()

	if fp == "" {
		t.Fatal("fingerprint should not be empty")
	}
}

func TestRequestSessionAttachment(t *testing.T) {
	t.Parallel()

	req := httpx.NewRequest(httptest.NewRequest(http.MethodGet, "/", nil))

	if req.Session() != nil {
		t.Fatal("session should be nil by default")
	}

	store := &stubSession{}
	req.SetSession(store)

	if req.Session() != store {
		t.Fatal("session should be the attached store")
	}
}

func TestRequestSegmentsEmpty(t *testing.T) {
	t.Parallel()

	req := httpx.NewRequest(httptest.NewRequest(http.MethodGet, "/", nil))

	if segments := req.Segments(); segments != nil {
		t.Fatalf("expected nil segments for root path, got %v", segments)
	}
}

func TestRequestQueryValues(t *testing.T) {
	t.Parallel()

	raw := httptest.NewRequest(http.MethodGet, "/?a=1&b=2", nil)
	req := httpx.NewRequest(raw)
	vals := req.QueryValues()

	if vals.Get("a") != "1" || vals.Get("b") != "2" {
		t.Fatalf("unexpected query values: %v", vals)
	}
}
