package routing_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/routing"
)

func newRouter() *routing.Router {
	return routing.New(nil)
}

func TestGetRoute(t *testing.T) {
	t.Parallel()

	r := newRouter()
	r.Get("/hello", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, "hello world")
	})

	req := httptest.NewRequest(http.MethodGet, "/hello", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "hello world" {
		t.Fatalf("expected 'hello world', got %q", rec.Body.String())
	}
}

func TestMethodRestriction(t *testing.T) {
	t.Parallel()

	r := newRouter()
	r.Post("/submit", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusCreated, "created")
	})

	req := httptest.NewRequest(http.MethodPost, "/submit", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/submit", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rec.Code)
	}
}

func TestPathParam(t *testing.T) {
	t.Parallel()

	r := newRouter()
	r.Get("/users/{id}", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, ctx.Param("id"))
	})

	req := httptest.NewRequest(http.MethodGet, "/users/42", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Body.String() != "42" {
		t.Fatalf("expected '42', got %q", rec.Body.String())
	}
}

func TestMiddlewareApplied(t *testing.T) {
	t.Parallel()

	var called bool

	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			next.ServeHTTP(w, r)
		})
	}

	r := newRouter()
	r.Use(mw)
	r.Get("/mw", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/mw", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if !called {
		t.Fatal("middleware was not called")
	}
}

func TestGroup(t *testing.T) {
	t.Parallel()

	r := newRouter()
	r.Group("/api", nil, func(sub *routing.Router) {
		sub.Get("/users", func(ctx *routing.Context) error {
			return ctx.Text(http.StatusOK, "users")
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "users" {
		t.Fatalf("expected 'users', got %q", rec.Body.String())
	}
}

func TestResource(t *testing.T) {
	t.Parallel()

	r := newRouter()
	r.Resource("posts", routing.ResourceHandlers{
		Index: func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "index") },
		Show:  func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "show:"+ctx.Param("id")) },
		Store: func(ctx *routing.Context) error { return ctx.Text(http.StatusCreated, "store") },
	})

	cases := []struct {
		method, path, body string
		code               int
	}{
		{http.MethodGet, "/posts", "index", http.StatusOK},
		{http.MethodGet, "/posts/99", "show:99", http.StatusOK},
		{http.MethodPost, "/posts", "store", http.StatusCreated},
	}

	for _, tc := range cases {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != tc.code {
			t.Errorf("%s %s: expected %d, got %d", tc.method, tc.path, tc.code, rec.Code)
		}
		if !strings.Contains(rec.Body.String(), tc.body) {
			t.Errorf("%s %s: expected body %q, got %q", tc.method, tc.path, tc.body, rec.Body.String())
		}
	}
}

func TestJSONResponse(t *testing.T) {
	t.Parallel()

	r := newRouter()
	r.Get("/json", func(ctx *routing.Context) error {
		return ctx.JSON(http.StatusOK, map[string]string{"key": "val"})
	})

	req := httptest.NewRequest(http.MethodGet, "/json", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected application/json, got %q", ct)
	}
	if !strings.Contains(rec.Body.String(), `"key"`) {
		t.Fatalf("unexpected body: %q", rec.Body.String())
	}
}

func TestNamedRoute(t *testing.T) {
	t.Parallel()

	reg := routing.NewRegistry()
	r := routing.New(reg)
	r.Get("/users/{id}", func(ctx *routing.Context) error { return nil })
	reg.Add("users.show", "GET", "/users/{id}")

	url := r.Route("users.show", map[string]string{"id": "7"})
	if url != "/users/7" {
		t.Fatalf("expected /users/7, got %q", url)
	}
}

func TestHandlerError(t *testing.T) {
	t.Parallel()

	r := newRouter()
	r.Get("/err", func(ctx *routing.Context) error {
		return &testError{msg: "boom"}
	})

	req := httptest.NewRequest(http.MethodGet, "/err", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

type testError struct{ msg string }

func (e *testError) Error() string { return e.msg }
