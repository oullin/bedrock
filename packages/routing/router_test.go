package routing_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/routing"
)

type testError struct{ msg string }

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

func (e *testError) Error() string { return e.msg }

func TestAnyRoute(t *testing.T) {
	t.Parallel()

	r := newRouter()
	r.Any("/catch", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, "caught:"+ctx.Method())
	})

	for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch} {
		req := httptest.NewRequest(method, "/catch", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("%s: expected 200, got %d", method, rec.Code)
		}
	}
}

func TestMatchRoute(t *testing.T) {
	t.Parallel()

	r := newRouter()
	r.Match([]string{http.MethodGet, http.MethodPost}, "/both", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, "matched")
	})

	for _, method := range []string{http.MethodGet, http.MethodPost} {
		req := httptest.NewRequest(method, "/both", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("%s: expected 200, got %d", method, rec.Code)
		}
	}

	req := httptest.NewRequest(http.MethodDelete, "/both", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 for DELETE, got %d", rec.Code)
	}
}

func TestRedirectRoute(t *testing.T) {
	t.Parallel()

	r := newRouter()
	r.Redirect("/old", "/new")

	req := httptest.NewRequest(http.MethodGet, "/old", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusFound {
		t.Fatalf("expected 302, got %d", rec.Code)
	}

	if loc := rec.Header().Get("Location"); loc != "/new" {
		t.Fatalf("expected Location=/new, got %q", loc)
	}
}

func TestRedirectRouteCustomStatus(t *testing.T) {
	t.Parallel()

	r := newRouter()
	r.Redirect("/old", "/new", http.StatusTemporaryRedirect)

	req := httptest.NewRequest(http.MethodGet, "/old", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusTemporaryRedirect {
		t.Fatalf("expected 307, got %d", rec.Code)
	}
}

func TestPermanentRedirectRoute(t *testing.T) {
	t.Parallel()

	r := newRouter()
	r.PermanentRedirect("/old", "/new")

	req := httptest.NewRequest(http.MethodGet, "/old", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusMovedPermanently {
		t.Fatalf("expected 301, got %d", rec.Code)
	}
}

func TestFallbackRoute(t *testing.T) {
	t.Parallel()

	r := newRouter()
	r.Get("/known", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, "known")
	})
	r.Fallback(func(ctx *routing.Context) error {
		return ctx.Text(http.StatusNotFound, "fallback")
	})

	req := httptest.NewRequest(http.MethodGet, "/known", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Body.String() != "known" {
		t.Fatalf("expected 'known', got %q", rec.Body.String())
	}

	req = httptest.NewRequest(http.MethodGet, "/unknown", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Body.String() != "fallback" {
		t.Fatalf("expected 'fallback', got %q", rec.Body.String())
	}
}

func TestGlobalPattern(t *testing.T) {
	t.Parallel()

	r := newRouter()
	r.Pattern("id", `[0-9]+`)
	r.Get("/users/{id}", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, ctx.Param("id"))
	})

	req := httptest.NewRequest(http.MethodGet, "/users/42", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/users/abc", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestRouteReturnsRouteObject(t *testing.T) {
	t.Parallel()

	r := newRouter()
	route := r.Get("/test", func(ctx *routing.Context) error { return nil })

	if route == nil {
		t.Fatal("expected route to be returned")
	}

	if route.URI() != "/test" {
		t.Fatalf("expected URI '/test', got %q", route.URI())
	}
}

func TestRouteFluentChaining(t *testing.T) {
	t.Parallel()

	r := newRouter()
	route := r.Get("/users/{id}", func(ctx *routing.Context) error { return nil }).
		Name("users.show").
		WhereNumber("id")

	if route.GetName() != "users.show" {
		t.Fatalf("expected 'users.show', got %q", route.GetName())
	}
}

func TestCollectionAccess(t *testing.T) {
	t.Parallel()

	r := newRouter()
	r.Get("/a", func(ctx *routing.Context) error { return nil })
	r.Post("/b", func(ctx *routing.Context) error { return nil })

	if r.Collection().Len() != 2 {
		t.Fatalf("expected 2 routes in collection, got %d", r.Collection().Len())
	}
}

func TestResourceWithCreateAndEdit(t *testing.T) {
	t.Parallel()

	r := newRouter()
	r.Resource("articles", routing.ResourceHandlers{
		Index:  func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "index") },
		Create: func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "create") },
		Store:  func(ctx *routing.Context) error { return ctx.Text(http.StatusCreated, "store") },
		Show:   func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "show") },
		Edit:   func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "edit") },
		Update: func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "update") },
	})

	cases := []struct {
		method, path, body string
		code               int
	}{
		{http.MethodGet, "/articles", "index", http.StatusOK},
		{http.MethodGet, "/articles/create", "create", http.StatusOK},
		{http.MethodPost, "/articles", "store", http.StatusCreated},
		{http.MethodGet, "/articles/1", "show", http.StatusOK},
		{http.MethodGet, "/articles/1/edit", "edit", http.StatusOK},
		{http.MethodPut, "/articles/1", "update", http.StatusOK},
	}

	for _, tc := range cases {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != tc.code {
			t.Errorf("%s %s: expected %d, got %d", tc.method, tc.path, tc.code, rec.Code)
		}

		if rec.Body.String() != tc.body {
			t.Errorf("%s %s: expected %q, got %q", tc.method, tc.path, tc.body, rec.Body.String())
		}
	}
}

func TestMiddlewareOrder(t *testing.T) {
	t.Parallel()

	var order []string

	mw1 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "mw1")
			next.ServeHTTP(w, r)
		})
	}

	mw2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "mw2")
			next.ServeHTTP(w, r)
		})
	}

	r := newRouter()
	r.Use(mw1, mw2)
	r.Get("/test", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if len(order) != 2 || order[0] != "mw1" || order[1] != "mw2" {
		t.Fatalf("expected [mw1, mw2], got %v", order)
	}
}

func TestGroupMiddlewareInherited(t *testing.T) {
	t.Parallel()

	var parentCalled, childCalled bool

	parentMw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			parentCalled = true
			next.ServeHTTP(w, r)
		})
	}

	childMw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			childCalled = true
			next.ServeHTTP(w, r)
		})
	}

	r := newRouter()
	r.Use(parentMw)
	r.Group("/api", []routing.MiddlewareFunc{childMw}, func(sub *routing.Router) {
		sub.Get("/test", func(ctx *routing.Context) error {
			return ctx.Text(http.StatusOK, "ok")
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if !parentCalled {
		t.Fatal("expected parent middleware to be called")
	}

	if !childCalled {
		t.Fatal("expected child middleware to be called")
	}
}

func TestPatchRoute(t *testing.T) {
	t.Parallel()

	r := newRouter()
	r.Patch("/items/{id}", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, "patched:"+ctx.Param("id"))
	})

	req := httptest.NewRequest(http.MethodPatch, "/items/5", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Body.String() != "patched:5" {
		t.Fatalf("expected 'patched:5', got %q", rec.Body.String())
	}
}

func TestOptionsRoute(t *testing.T) {
	t.Parallel()

	r := newRouter()
	r.Options("/cors", func(ctx *routing.Context) error {
		ctx.SetHeader("Access-Control-Allow-Origin", "*")

		return ctx.NoContent()
	})

	req := httptest.NewRequest(http.MethodOptions, "/cors", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}

	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("expected CORS header")
	}
}
