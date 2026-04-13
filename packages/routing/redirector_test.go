package routing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/routing"
)

func newRedirector() (*routing.Redirector, *routing.Registry) {
	reg := routing.NewRegistry()
	gen := routing.NewUrlGenerator(reg, "https://example.com", nil)

	return routing.NewRedirector(gen), reg
}

func TestRedirectorTo(t *testing.T) {
	t.Parallel()

	rd, _ := newRedirector()

	r := routing.New(nil)
	r.Get("/old", func(ctx *routing.Context) error {
		return rd.To(ctx, "/new")
	})

	req := httptest.NewRequest(http.MethodGet, "/old", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}

	if loc := rec.Header().Get("Location"); loc != "https://example.com/new" {
		t.Fatalf("expected redirect to https://example.com/new, got %q", loc)
	}
}

func TestRedirectorToCustomStatus(t *testing.T) {
	t.Parallel()

	rd, _ := newRedirector()

	r := routing.New(nil)
	r.Get("/old", func(ctx *routing.Context) error {
		return rd.To(ctx, "/new", http.StatusMovedPermanently)
	})

	req := httptest.NewRequest(http.MethodGet, "/old", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusMovedPermanently {
		t.Fatalf("expected 301, got %d", rec.Code)
	}
}

func TestRedirectorAway(t *testing.T) {
	t.Parallel()

	rd, _ := newRedirector()

	r := routing.New(nil)
	r.Get("/external", func(ctx *routing.Context) error {
		return rd.Away(ctx, "https://google.com")
	})

	req := httptest.NewRequest(http.MethodGet, "/external", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if loc := rec.Header().Get("Location"); loc != "https://google.com" {
		t.Fatalf("expected redirect to https://google.com, got %q", loc)
	}
}

func TestRedirectorSecure(t *testing.T) {
	t.Parallel()

	reg := routing.NewRegistry()
	gen := routing.NewUrlGenerator(reg, "http://example.com", nil)
	rd := routing.NewRedirector(gen)

	r := routing.New(nil)
	r.Get("/insecure", func(ctx *routing.Context) error {
		return rd.Secure(ctx, "/secure-page")
	})

	req := httptest.NewRequest(http.MethodGet, "/insecure", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if loc := rec.Header().Get("Location"); loc != "https://example.com/secure-page" {
		t.Fatalf("expected HTTPS redirect, got %q", loc)
	}
}

func TestRedirectorBack(t *testing.T) {
	t.Parallel()

	rd, _ := newRedirector()

	r := routing.New(nil)
	r.Get("/back", func(ctx *routing.Context) error {
		return rd.Back(ctx)
	})

	req := httptest.NewRequest(http.MethodGet, "/back", nil)
	req.Header.Set("Referer", "https://example.com/previous")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if loc := rec.Header().Get("Location"); loc != "https://example.com/previous" {
		t.Fatalf("expected redirect to referer, got %q", loc)
	}
}

func TestRedirectorBackFallback(t *testing.T) {
	t.Parallel()

	rd, _ := newRedirector()

	r := routing.New(nil)
	r.Get("/back", func(ctx *routing.Context) error {
		return rd.Back(ctx, "/fallback")
	})

	req := httptest.NewRequest(http.MethodGet, "/back", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if loc := rec.Header().Get("Location"); loc != "/fallback" {
		t.Fatalf("expected redirect to fallback, got %q", loc)
	}
}

func TestRedirectorRefresh(t *testing.T) {
	t.Parallel()

	rd, _ := newRedirector()

	r := routing.New(nil)
	r.Get("/refresh", func(ctx *routing.Context) error {
		return rd.Refresh(ctx)
	})

	req := httptest.NewRequest(http.MethodGet, "/refresh", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if loc := rec.Header().Get("Location"); loc != "https://example.com/refresh" {
		t.Fatalf("expected redirect to current URL, got %q", loc)
	}
}

func TestRedirectorRoute(t *testing.T) {
	t.Parallel()

	rd, reg := newRedirector()
	reg.Add("users.show", "GET", "/users/{id}")

	r := routing.New(nil)
	r.Get("/go-to-user", func(ctx *routing.Context) error {
		return rd.Route(ctx, "users.show", map[string]string{"id": "42"})
	})

	req := httptest.NewRequest(http.MethodGet, "/go-to-user", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if loc := rec.Header().Get("Location"); loc != "https://example.com/users/42" {
		t.Fatalf("expected redirect to named route, got %q", loc)
	}
}

func TestRedirectorRouteCustomStatus(t *testing.T) {
	t.Parallel()

	rd, reg := newRedirector()
	reg.Add("home", "GET", "/")

	r := routing.New(nil)
	r.Get("/redirect-home", func(ctx *routing.Context) error {
		return rd.Route(ctx, "home", nil, http.StatusMovedPermanently)
	})

	req := httptest.NewRequest(http.MethodGet, "/redirect-home", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusMovedPermanently {
		t.Fatalf("expected 301, got %d", rec.Code)
	}
}

func TestRedirectorGetUrlGenerator(t *testing.T) {
	t.Parallel()

	rd, _ := newRedirector()

	if rd.GetUrlGenerator() == nil {
		t.Fatal("expected URL generator")
	}
}
