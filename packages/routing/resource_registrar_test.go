package routing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/routing"
)

func TestResourceRegistrarFullCRUD(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	reg := routing.NewResourceRegistrar(r)

	reg.Register("posts", routing.ResourceHandlers{
		Index:   func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "index") },
		Create:  func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "create") },
		Store:   func(ctx *routing.Context) error { return ctx.Text(http.StatusCreated, "store") },
		Show:    func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "show:"+ctx.Param("post")) },
		Edit:    func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "edit:"+ctx.Param("post")) },
		Update:  func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "update:"+ctx.Param("post")) },
		Destroy: func(ctx *routing.Context) error { return ctx.NoContent() },
	})

	cases := []struct {
		method, path, body string
		code               int
	}{
		{http.MethodGet, "/posts", "index", http.StatusOK},
		{http.MethodGet, "/posts/create", "create", http.StatusOK},
		{http.MethodPost, "/posts", "store", http.StatusCreated},
		{http.MethodGet, "/posts/1", "show:1", http.StatusOK},
		{http.MethodGet, "/posts/1/edit", "edit:1", http.StatusOK},
		{http.MethodPut, "/posts/1", "update:1", http.StatusOK},
		{http.MethodDelete, "/posts/1", "", http.StatusNoContent},
	}

	for _, tc := range cases {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)

		if rec.Code != tc.code {
			t.Errorf("%s %s: expected %d, got %d", tc.method, tc.path, tc.code, rec.Code)
		}

		if tc.body != "" && rec.Body.String() != tc.body {
			t.Errorf("%s %s: expected %q, got %q", tc.method, tc.path, tc.body, rec.Body.String())
		}
	}
}

func TestResourceRegistrarOnly(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	reg := routing.NewResourceRegistrar(r)

	reg.Register("items", routing.ResourceHandlers{
		Index: func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "index") },
		Show:  func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "show") },
		Store: func(ctx *routing.Context) error { return ctx.Text(http.StatusCreated, "store") },
	}, routing.ResourceOptions{
		Only: []string{"index", "show"},
	})

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for index, got %d", rec.Code)
	}

	req = httptest.NewRequest(http.MethodPost, "/items", nil)
	rec = httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405 for excluded store, got %d", rec.Code)
	}
}

func TestResourceRegistrarExcept(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	reg := routing.NewResourceRegistrar(r)

	reg.Register("things", routing.ResourceHandlers{
		Index:   func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "index") },
		Store:   func(ctx *routing.Context) error { return ctx.Text(http.StatusCreated, "store") },
		Destroy: func(ctx *routing.Context) error { return ctx.NoContent() },
	}, routing.ResourceOptions{
		Except: []string{"destroy"},
	})

	req := httptest.NewRequest(http.MethodGet, "/things", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestResourceRegistrarCustomNames(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	reg := routing.NewResourceRegistrar(r)

	reg.Register("posts", routing.ResourceHandlers{
		Index: func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "index") },
	}, routing.ResourceOptions{
		Names: map[string]string{"index": "blog.posts"},
	})

	if !r.Collection().HasNamedRoute("blog.posts") {
		t.Fatal("expected custom route name 'blog.posts'")
	}
}

func TestResourceRegistrarCustomParameters(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	reg := routing.NewResourceRegistrar(r)

	reg.Register("posts", routing.ResourceHandlers{
		Show: func(ctx *routing.Context) error {
			return ctx.Text(http.StatusOK, "slug:"+ctx.Param("slug"))
		},
	}, routing.ResourceOptions{
		Parameters: map[string]string{"posts": "slug"},
	})

	req := httptest.NewRequest(http.MethodGet, "/posts/hello-world", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Body.String() != "slug:hello-world" {
		t.Fatalf("expected 'slug:hello-world', got %q", rec.Body.String())
	}
}

func TestResourceRegistrarWithMiddleware(t *testing.T) {
	t.Parallel()

	var called bool

	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			next.ServeHTTP(w, r)
		})
	}

	r := routing.New(nil)
	reg := routing.NewResourceRegistrar(r)

	reg.Register("posts", routing.ResourceHandlers{
		Index: func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "index") },
	}, routing.ResourceOptions{
		Middleware: []routing.MiddlewareFunc{mw},
	})

	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if !called {
		t.Fatal("expected resource middleware to be called")
	}
}

func TestAPIResource(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	reg := routing.NewResourceRegistrar(r)

	reg.APIResource("articles", routing.ResourceHandlers{
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
		{http.MethodGet, "/articles", http.StatusOK},
		{http.MethodPost, "/articles", http.StatusCreated},
		{http.MethodGet, "/articles/1", http.StatusOK},
		{http.MethodPut, "/articles/1", http.StatusOK},
		{http.MethodDelete, "/articles/1", http.StatusNoContent},
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

func TestSingletonResource(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	reg := routing.NewResourceRegistrar(r)

	reg.Singleton("profile", routing.SingletonHandlers{
		Show:   func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "show") },
		Edit:   func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "edit") },
		Update: func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "update") },
	})

	cases := []struct {
		method, path, body string
		code               int
	}{
		{http.MethodGet, "/profile", "show", http.StatusOK},
		{http.MethodGet, "/profile/edit", "edit", http.StatusOK},
		{http.MethodPut, "/profile", "update", http.StatusOK},
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

func TestSingletonWithDestroy(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	reg := routing.NewResourceRegistrar(r)

	reg.Singleton("profile", routing.SingletonHandlers{
		Show:    func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "show") },
		Update:  func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "update") },
		Destroy: func(ctx *routing.Context) error { return ctx.NoContent() },
	})

	req := httptest.NewRequest(http.MethodDelete, "/profile", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
}

func TestAPISingleton(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	reg := routing.NewResourceRegistrar(r)

	reg.APISingleton("settings", routing.SingletonHandlers{
		Show:   func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "show") },
		Update: func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "update") },
	})

	cases := []struct {
		method, path, body string
		code               int
	}{
		{http.MethodGet, "/settings", "show", http.StatusOK},
		{http.MethodPut, "/settings", "update", http.StatusOK},
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

func TestResourceRegistrarAutoNames(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	reg := routing.NewResourceRegistrar(r)

	reg.Register("posts", routing.ResourceHandlers{
		Index: func(ctx *routing.Context) error { return nil },
		Show:  func(ctx *routing.Context) error { return nil },
		Store: func(ctx *routing.Context) error { return nil },
	})

	if !r.Collection().HasNamedRoute("posts.index") {
		t.Fatal("expected 'posts.index' named route")
	}

	if !r.Collection().HasNamedRoute("posts.show") {
		t.Fatal("expected 'posts.show' named route")
	}

	if !r.Collection().HasNamedRoute("posts.store") {
		t.Fatal("expected 'posts.store' named route")
	}
}

func TestResourceRegistrarSingularizeParam(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	reg := routing.NewResourceRegistrar(r)

	reg.Register("categories", routing.ResourceHandlers{
		Show: func(ctx *routing.Context) error {
			return ctx.Text(http.StatusOK, "cat:"+ctx.Param("category"))
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/categories/5", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Body.String() != "cat:5" {
		t.Fatalf("expected 'cat:5', got %q", rec.Body.String())
	}
}

func TestAPIResourceExclude(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	reg := routing.NewResourceRegistrar(r)

	reg.APIResource("posts", routing.ResourceHandlers{
		Index:   func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "index") },
		Store:   func(ctx *routing.Context) error { return ctx.Text(http.StatusCreated, "store") },
		Show:    func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "show") },
		Update:  func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "update") },
		Destroy: func(ctx *routing.Context) error { return ctx.NoContent() },
	}, routing.ResourceOptions{
		Except: []string{"destroy"},
	})

	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}

func TestAPIResourceOnly(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	reg := routing.NewResourceRegistrar(r)

	reg.APIResource("posts", routing.ResourceHandlers{
		Index: func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "index") },
		Show:  func(ctx *routing.Context) error { return ctx.Text(http.StatusOK, "show") },
	}, routing.ResourceOptions{
		Only: []string{"index", "show"},
	})

	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
