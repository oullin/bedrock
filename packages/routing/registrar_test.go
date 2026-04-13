package routing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/routing"
)

func TestRegistrarGetWithMiddleware(t *testing.T) {
	t.Parallel()

	var called bool

	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			next.ServeHTTP(w, r)
		})
	}

	r := routing.New(nil)
	r.WithMiddleware(mw).Get("/test", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if !called {
		t.Fatal("expected middleware to be called")
	}

	if rec.Body.String() != "ok" {
		t.Fatalf("expected 'ok', got %q", rec.Body.String())
	}
}

func TestRegistrarPostWithMiddleware(t *testing.T) {
	t.Parallel()

	var called bool

	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			next.ServeHTTP(w, r)
		})
	}

	r := routing.New(nil)
	r.WithMiddleware(mw).Post("/submit", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusCreated, "created")
	})

	req := httptest.NewRequest(http.MethodPost, "/submit", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if !called {
		t.Fatal("expected middleware to be called")
	}

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}

func TestRegistrarPrefixGroup(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.WithPrefix("/api").Group(func(sub *routing.Router) {
		sub.Get("/users", func(ctx *routing.Context) error {
			return ctx.Text(http.StatusOK, "api-users")
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Body.String() != "api-users" {
		t.Fatalf("expected 'api-users', got %q", rec.Body.String())
	}
}

func TestRegistrarDomainGroup(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.WithDomain("api.example.com").Group(func(sub *routing.Router) {
		sub.Get("/test", func(ctx *routing.Context) error {
			return ctx.Text(http.StatusOK, "domain-test")
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestRegistrarNameGroup(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.WithName("api").Group(func(sub *routing.Router) {
		route := sub.Get("/users", func(ctx *routing.Context) error { return nil })
		route.Name("users.index")

		if route.GetName() != "api.users.index" {
			t.Fatalf("expected 'api.users.index', got %q", route.GetName())
		}
	})
}

func TestRegistrarChaining(t *testing.T) {
	t.Parallel()

	var mwCalled bool

	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mwCalled = true
			next.ServeHTTP(w, r)
		})
	}

	r := routing.New(nil)
	r.WithMiddleware(mw).Prefix("/api").Name("api").Group(func(sub *routing.Router) {
		sub.Get("/test", func(ctx *routing.Context) error {
			return ctx.Text(http.StatusOK, "chained")
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if !mwCalled {
		t.Fatal("expected middleware to be called")
	}

	if rec.Body.String() != "chained" {
		t.Fatalf("expected 'chained', got %q", rec.Body.String())
	}
}

func TestRegistrarWhereConstraint(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.WithPrefix("/api").Where("id", `[0-9]+`).Group(func(sub *routing.Router) {
		sub.Get("/users/{id}", func(ctx *routing.Context) error {
			return ctx.Text(http.StatusOK, ctx.Param("id"))
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/users/42", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodGet, "/api/users/abc", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestRegistrarAny(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.WithPrefix("/api").Any("/catch", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, "caught")
	})

	for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodPut} {
		req := httptest.NewRequest(method, "/api/catch", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Body.String() != "caught" {
			t.Fatalf("%s: expected 'caught', got %q", method, rec.Body.String())
		}
	}
}

func TestRegistrarMatch(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.WithPrefix("/api").Match(
		[]string{http.MethodGet, http.MethodPost},
		"/both",
		func(ctx *routing.Context) error {
			return ctx.Text(http.StatusOK, "matched")
		},
	)

	for _, method := range []string{http.MethodGet, http.MethodPost} {
		req := httptest.NewRequest(method, "/api/both", nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Body.String() != "matched" {
			t.Fatalf("%s: expected 'matched', got %q", method, rec.Body.String())
		}
	}
}

func TestRegistrarPut(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.WithPrefix("/api").Put("/items/{id}", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, "updated:"+ctx.Param("id"))
	})

	req := httptest.NewRequest(http.MethodPut, "/api/items/5", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Body.String() != "updated:5" {
		t.Fatalf("expected 'updated:5', got %q", rec.Body.String())
	}
}

func TestRegistrarDelete(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.WithPrefix("/api").Delete("/items/{id}", func(ctx *routing.Context) error {
		return ctx.NoContent()
	})

	req := httptest.NewRequest(http.MethodDelete, "/api/items/5", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}

func TestRegistrarPatch(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.WithPrefix("/api").Patch("/items/{id}", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, "patched")
	})

	req := httptest.NewRequest(http.MethodPatch, "/api/items/5", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Body.String() != "patched" {
		t.Fatalf("expected 'patched', got %q", rec.Body.String())
	}
}

func TestRegistrarOptions(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.WithPrefix("/api").Options("/cors", func(ctx *routing.Context) error {
		ctx.SetHeader("Access-Control-Allow-Origin", "*")

		return ctx.NoContent()
	})

	req := httptest.NewRequest(http.MethodOptions, "/api/cors", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}

	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("expected CORS header")
	}
}

func TestRegistrarResource(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.WithPrefix("/api").Resource("posts", routing.ResourceHandlers{
		Index: func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "index") },
		Store: func(ctx *routing.Context) error { return ctx.Text(http.StatusCreated, "store") },
		Show:  func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "show") },
	})

	cases := []struct {
		method, path, body string
		code               int
	}{
		{http.MethodGet, "/api/posts", "index", http.StatusOK},
		{http.MethodPost, "/api/posts", "store", http.StatusCreated},
		{http.MethodGet, "/api/posts/1", "show", http.StatusOK},
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

func TestRegistrarAPIResource(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.WithPrefix("/api").APIResource("posts", routing.ResourceHandlers{
		Index:   func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "index") },
		Store:   func(ctx *routing.Context) error { return ctx.Text(http.StatusCreated, "store") },
		Show:    func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "show") },
		Update:  func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "update") },
		Destroy: func(ctx *routing.Context) error { return ctx.NoContent() },
	})

	cases := []struct {
		method, path string
		code         int
	}{
		{http.MethodGet, "/api/posts", http.StatusOK},
		{http.MethodPost, "/api/posts", http.StatusCreated},
		{http.MethodGet, "/api/posts/1", http.StatusOK},
		{http.MethodPut, "/api/posts/1", http.StatusOK},
		{http.MethodDelete, "/api/posts/1", http.StatusNoContent},
	}

	for _, tc := range cases {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != tc.code {
			t.Errorf("%s %s: expected %d, got %d", tc.method, tc.path, tc.code, rec.Code)
		}
	}
}

func TestRegistrarMultipleMiddleware(t *testing.T) {
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

	r := routing.New(nil)
	r.WithMiddleware(mw1).Middleware(mw2).Get("/test", func(ctx *routing.Context) error {
		return ctx.Text(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if len(order) != 2 || order[0] != "mw1" || order[1] != "mw2" {
		t.Fatalf("expected [mw1, mw2], got %v", order)
	}
}

func TestRegistrarGroupWithMiddleware(t *testing.T) {
	t.Parallel()

	var called bool

	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			next.ServeHTTP(w, r)
		})
	}

	r := routing.New(nil)
	r.WithMiddleware(mw).Prefix("/api").Group(func(sub *routing.Router) {
		sub.Get("/test1", func(ctx *routing.Context) error {
			return ctx.Text(http.StatusOK, "test1")
		})
		sub.Get("/test2", func(ctx *routing.Context) error {
			return ctx.Text(http.StatusOK, "test2")
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/test1", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if !called {
		t.Fatal("expected middleware to be called for grouped routes")
	}
}

func TestRegistrarGroupWithoutPrefix(t *testing.T) {
	t.Parallel()

	var called bool

	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			next.ServeHTTP(w, r)
		})
	}

	r := routing.New(nil)
	r.WithMiddleware(mw).Group(func(sub *routing.Router) {
		sub.Get("/no-prefix", func(ctx *routing.Context) error {
			return ctx.Text(http.StatusOK, "ok")
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/no-prefix", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if !called {
		t.Fatal("expected middleware without prefix")
	}
}
