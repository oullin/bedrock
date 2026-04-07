package routing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/routing"
)

// Upstream: testBasicRouting
func TestGetRouteRegistrationAndDispatch(t *testing.T) {
	t.Parallel()

	router := routing.New(nil)
	router.Get("/hello", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, "hello world")
	})

	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "hello world" {
		t.Fatalf("expected 'hello world', got %q", rec.Body.String())
	}
}

// Test method-specific routing
func TestHandleWithMethodRestriction(t *testing.T) {
	t.Parallel()

	router := routing.New(nil)
	router.Handle(http.MethodPost, "/submit", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusCreated, "created")
	})

	// POST should work
	req := httptest.NewRequest(http.MethodPost, "/submit", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	// GET should be method not allowed
	req = httptest.NewRequest(http.MethodGet, "/submit", nil)
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

// Test context text response
func TestContextTextResponse(t *testing.T) {
	t.Parallel()

	router := routing.New(nil)
	router.Get("/text", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusAccepted, "accepted")
	})

	req := httptest.NewRequest(http.MethodGet, "/text", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "text/plain; charset=utf-8" {
		t.Fatalf("expected text/plain content-type, got %q", ct)
	}
	if rec.Body.String() != "accepted" {
		t.Fatalf("expected 'accepted', got %q", rec.Body.String())
	}
}

// Test handler error returns 500
func TestHandlerErrorReturns500(t *testing.T) {
	t.Parallel()

	router := routing.New(nil)
	router.Get("/error", func(_ *routing.Context) error {
		return http.ErrNoCookie // any error
	})

	req := httptest.NewRequest(http.MethodGet, "/error", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

// Test context has request and writer
func TestContextHasRequestAndWriter(t *testing.T) {
	t.Parallel()

	router := routing.New(nil)
	var gotPath string
	router.Get("/check", func(ctx *routing.Context) error {
		gotPath = ctx.Request.URL.Path
		return ctx.Text(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/check", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	if gotPath != "/check" {
		t.Fatalf("expected path '/check', got %q", gotPath)
	}
}

// Test context render without renderer returns error
func TestContextRenderWithoutRendererReturnsError(t *testing.T) {
	t.Parallel()

	router := routing.New(nil)
	router.Get("/render", func(ctx *routing.Context) error {
		return ctx.Render("welcome", nil)
	})

	req := httptest.NewRequest(http.MethodGet, "/render", nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Should return 500 because renderer is nil
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 when renderer is nil, got %d", rec.Code)
	}
}

// Test multiple routes
func TestMultipleRoutes(t *testing.T) {
	t.Parallel()

	router := routing.New(nil)
	router.Get("/a", func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "route-a") })
	router.Get("/b", func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "route-b") })

	for _, tc := range []struct {
		path string
		want string
	}{
		{"/a", "route-a"},
		{"/b", "route-b"},
	} {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, tc.path, nil))
		if rec.Body.String() != tc.want {
			t.Fatalf("path %s: expected %q, got %q", tc.path, tc.want, rec.Body.String())
		}
	}
}

// Test Handle with empty method allows all methods
func TestHandleWithEmptyMethodAllowsAll(t *testing.T) {
	t.Parallel()

	router := routing.New(nil)
	router.Handle("", "/any", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, "any method")
	})

	for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete} {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(method, "/any", nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("method %s: expected 200, got %d", method, rec.Code)
		}
	}
}
