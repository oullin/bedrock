package routing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/anvil/routing"
)

// Laravel: testBasicRouting
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

// Laravel: testRouteParameterBinding
func TestRouteParameterBinding(t *testing.T) {
	t.Parallel()

	router := routing.New(nil)
	router.Get("/users/{id}", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, "user:"+ctx.Param("id"))
	})
	router.Get("/posts/{postId}/comments/{commentId}", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, ctx.Param("postId")+":"+ctx.Param("commentId"))
	})

	t.Run("single parameter", func(t *testing.T) {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/users/42", nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
		if rec.Body.String() != "user:42" {
			t.Fatalf("expected 'user:42', got %q", rec.Body.String())
		}
	})

	t.Run("multiple parameters", func(t *testing.T) {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/posts/10/comments/5", nil))
		if rec.Code != http.StatusOK {
			t.Fatalf("expected 200, got %d", rec.Code)
		}
		if rec.Body.String() != "10:5" {
			t.Fatalf("expected '10:5', got %q", rec.Body.String())
		}
	})
}

// Laravel: test POST/PUT/DELETE route registration
func TestPostRouteDispatch(t *testing.T) {
	t.Parallel()

	router := routing.New(nil)
	router.Post("/items", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusCreated, "created")
	})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/items", nil))
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	// GET should be rejected.
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/items", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 for GET on POST route, got %d", rec.Code)
	}
}

func TestPutRouteDispatch(t *testing.T) {
	t.Parallel()

	router := routing.New(nil)
	router.Put("/items/{id}", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, "updated:"+ctx.Param("id"))
	})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodPut, "/items/7", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "updated:7" {
		t.Fatalf("expected 'updated:7', got %q", rec.Body.String())
	}
}

func TestDeleteRouteDispatch(t *testing.T) {
	t.Parallel()

	router := routing.New(nil)
	router.Delete("/items/{id}", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusNoContent, "")
	})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodDelete, "/items/3", nil))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
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

// ======================== PATCH / OPTIONS TESTS ========================

func TestPatchRouteDispatch(t *testing.T) {
	t.Parallel()

	router := routing.New(nil)
	router.Patch("/items/{id}", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, "patched-"+ctx.Param("id"))
	})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodPatch, "/items/5", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "patched-5" {
		t.Fatalf("expected 'patched-5', got %q", rec.Body.String())
	}
}

func TestOptionsRouteDispatch(t *testing.T) {
	t.Parallel()

	router := routing.New(nil)
	router.Options("/cors", func(ctx *routing.Context) error {
		ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		return ctx.Text(http.StatusNoContent, "")
	})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest("OPTIONS", "/cors", nil))
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("expected CORS header")
	}
}

// ======================== GROUP TESTS ========================

func TestGroupPrefixAndMiddleware(t *testing.T) {
	t.Parallel()

	router := routing.New(nil)

	auth := routing.Middleware(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Auth", "applied")
			next.ServeHTTP(w, r)
		})
	})

	router.Group("/api", []routing.Middleware{auth}, func(r *routing.Router) {
		r.Get("/users", func(ctx *routing.Context) error {
			return ctx.Text(http.StatusOK, "users list")
		})
	})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/users", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "users list" {
		t.Fatalf("expected 'users list', got %q", rec.Body.String())
	}
	if rec.Header().Get("X-Auth") != "applied" {
		t.Fatal("expected group middleware to be applied")
	}
}

func TestNestedGroups(t *testing.T) {
	t.Parallel()

	router := routing.New(nil)
	router.Group("/api", nil, func(api *routing.Router) {
		api.Group("/v1", nil, func(v1 *routing.Router) {
			v1.Get("/status", func(ctx *routing.Context) error {
				return ctx.Text(http.StatusOK, "v1-ok")
			})
		})
	})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/api/v1/status", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() != "v1-ok" {
		t.Fatalf("expected 'v1-ok', got %q", rec.Body.String())
	}
}

func TestGroupMiddlewareDoesNotAffectOuterRoutes(t *testing.T) {
	t.Parallel()

	router := routing.New(nil)
	router.Get("/public", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, "public")
	})

	router.Group("/admin", []routing.Middleware{func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("X-Admin", "yes")
			next.ServeHTTP(w, r)
		})
	}}, func(r *routing.Router) {
		r.Get("/dashboard", func(ctx *routing.Context) error {
			return ctx.Text(http.StatusOK, "dashboard")
		})
	})

	// Public route should NOT have admin middleware.
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/public", nil))
	if rec.Header().Get("X-Admin") != "" {
		t.Fatal("expected public route to NOT have admin middleware")
	}

	// Admin route should have admin middleware.
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/admin/dashboard", nil))
	if rec.Header().Get("X-Admin") != "yes" {
		t.Fatal("expected admin route to have admin middleware")
	}
}

// ======================== RESOURCE TESTS ========================

func TestResourceRouting(t *testing.T) {
	t.Parallel()

	router := routing.New(nil)
	router.Resource("posts", routing.ResourceHandlers{
		Index: func(ctx *routing.Context) error {
			return ctx.Text(http.StatusOK, "index")
		},
		Show: func(ctx *routing.Context) error {
			return ctx.Text(http.StatusOK, "show-"+ctx.Param("id"))
		},
		Store: func(ctx *routing.Context) error {
			return ctx.Text(http.StatusCreated, "stored")
		},
		Update: func(ctx *routing.Context) error {
			return ctx.Text(http.StatusOK, "updated-"+ctx.Param("id"))
		},
		Destroy: func(ctx *routing.Context) error {
			return ctx.Text(http.StatusNoContent, "")
		},
	})

	tests := []struct {
		method string
		path   string
		code   int
		body   string
	}{
		{http.MethodGet, "/posts", http.StatusOK, "index"},
		{http.MethodGet, "/posts/42", http.StatusOK, "show-42"},
		{http.MethodPost, "/posts", http.StatusCreated, "stored"},
		{http.MethodPut, "/posts/42", http.StatusOK, "updated-42"},
		{http.MethodDelete, "/posts/42", http.StatusNoContent, ""},
	}

	for _, tc := range tests {
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, httptest.NewRequest(tc.method, tc.path, nil))
		if rec.Code != tc.code {
			t.Fatalf("%s %s: expected %d, got %d", tc.method, tc.path, tc.code, rec.Code)
		}
		if tc.body != "" && rec.Body.String() != tc.body {
			t.Fatalf("%s %s: expected %q, got %q", tc.method, tc.path, tc.body, rec.Body.String())
		}
	}
}

func TestPartialResourceRouting(t *testing.T) {
	t.Parallel()

	router := routing.New(nil)
	router.Resource("comments", routing.ResourceHandlers{
		Index: func(ctx *routing.Context) error {
			return ctx.Text(http.StatusOK, "comments")
		},
	})

	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/comments", nil))
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

// ======================== NAMED ROUTE TESTS ========================

func TestNamedRouteUrlGeneration(t *testing.T) {
	t.Parallel()

	router := routing.New(nil)
	router.Get("/users/{id}", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, "user")
	})
	router.Name("/users/{id}", "user.show")

	url := router.Route("user.show", map[string]string{"id": "42"})
	if url != "/users/42" {
		t.Fatalf("expected '/users/42', got %q", url)
	}
}

func TestNamedRouteUnknownReturnsEmpty(t *testing.T) {
	t.Parallel()

	router := routing.New(nil)
	if url := router.Route("nonexistent", nil); url != "" {
		t.Fatalf("expected empty string for unknown route, got %q", url)
	}
}

func TestNamedRouteWithGroupPrefix(t *testing.T) {
	t.Parallel()

	router := routing.New(nil)
	router.Group("/api", nil, func(api *routing.Router) {
		api.Get("/users/{id}", func(ctx *routing.Context) error {
			return ctx.Text(http.StatusOK, "user")
		})
		api.Name("/users/{id}", "api.user.show")
	})

	url := router.Route("api.user.show", map[string]string{"id": "7"})
	if url != "/api/users/7" {
		t.Fatalf("expected '/api/users/7', got %q", url)
	}
}
