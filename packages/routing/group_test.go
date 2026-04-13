package routing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/routing"
)

func TestMergeGroupPrefix(t *testing.T) {
	t.Parallel()

	parent := routing.GroupAttributes{Prefix: "/api"}
	child := routing.GroupAttributes{Prefix: "/v1"}

	merged := routing.MergeGroup(parent, child)

	if merged.Prefix != "/api/v1" {
		t.Fatalf("expected '/api/v1', got %q", merged.Prefix)
	}
}

func TestMergeGroupDomainOverride(t *testing.T) {
	t.Parallel()

	parent := routing.GroupAttributes{Domain: "example.com"}
	child := routing.GroupAttributes{Domain: "api.example.com"}

	merged := routing.MergeGroup(parent, child)

	if merged.Domain != "api.example.com" {
		t.Fatalf("expected 'api.example.com', got %q", merged.Domain)
	}
}

func TestMergeGroupDomainInherit(t *testing.T) {
	t.Parallel()

	parent := routing.GroupAttributes{Domain: "example.com"}
	child := routing.GroupAttributes{}

	merged := routing.MergeGroup(parent, child)

	if merged.Domain != "example.com" {
		t.Fatalf("expected 'example.com', got %q", merged.Domain)
	}
}

func TestMergeGroupNamePrefix(t *testing.T) {
	t.Parallel()

	parent := routing.GroupAttributes{NamePrefix: "api"}
	child := routing.GroupAttributes{NamePrefix: "users"}

	merged := routing.MergeGroup(parent, child)

	if merged.NamePrefix != "api.users" {
		t.Fatalf("expected 'api.users', got %q", merged.NamePrefix)
	}
}

func TestMergeGroupNamePrefixParentOnly(t *testing.T) {
	t.Parallel()

	parent := routing.GroupAttributes{NamePrefix: "api"}
	child := routing.GroupAttributes{}

	merged := routing.MergeGroup(parent, child)

	if merged.NamePrefix != "api" {
		t.Fatalf("expected 'api', got %q", merged.NamePrefix)
	}
}

func TestMergeGroupNamePrefixChildOnly(t *testing.T) {
	t.Parallel()

	parent := routing.GroupAttributes{}
	child := routing.GroupAttributes{NamePrefix: "users"}

	merged := routing.MergeGroup(parent, child)

	if merged.NamePrefix != "users" {
		t.Fatalf("expected 'users', got %q", merged.NamePrefix)
	}
}

func TestMergeGroupMiddleware(t *testing.T) {
	t.Parallel()

	mw1 := func(next http.Handler) http.Handler { return next }
	mw2 := func(next http.Handler) http.Handler { return next }

	parent := routing.GroupAttributes{Middleware: []routing.MiddlewareFunc{mw1}}
	child := routing.GroupAttributes{Middleware: []routing.MiddlewareFunc{mw2}}

	merged := routing.MergeGroup(parent, child)

	if len(merged.Middleware) != 2 {
		t.Fatalf("expected 2 middleware, got %d", len(merged.Middleware))
	}
}

func TestMergeGroupWhere(t *testing.T) {
	t.Parallel()

	parent := routing.GroupAttributes{Where: map[string]string{"id": `[0-9]+`}}
	child := routing.GroupAttributes{Where: map[string]string{"slug": `[a-z-]+`}}

	merged := routing.MergeGroup(parent, child)

	if merged.Where["id"] != `[0-9]+` {
		t.Fatalf("expected parent constraint to be inherited")
	}

	if merged.Where["slug"] != `[a-z-]+` {
		t.Fatalf("expected child constraint to be present")
	}
}

func TestMergeGroupWhereOverride(t *testing.T) {
	t.Parallel()

	parent := routing.GroupAttributes{Where: map[string]string{"id": `[0-9]+`}}
	child := routing.GroupAttributes{Where: map[string]string{"id": `[a-z]+`}}

	merged := routing.MergeGroup(parent, child)

	if merged.Where["id"] != `[a-z]+` {
		t.Fatalf("expected child constraint to override parent, got %q", merged.Where["id"])
	}
}

func TestGroupWithAttributes(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.GroupWith(routing.GroupAttributes{
		Prefix:     "/api",
		NamePrefix: "api",
	}, func(sub *routing.Router) {
		sub.Get("/users", func(ctx *routing.Context) error {
			return ctx.Text(http.StatusOK, "users")
		}).Name("users.index")
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

func TestGroupWithMiddleware(t *testing.T) {
	t.Parallel()

	var called bool

	mw := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			called = true
			next.ServeHTTP(w, r)
		})
	}

	r := routing.New(nil)
	r.GroupWith(routing.GroupAttributes{
		Prefix:     "/api",
		Middleware: []routing.MiddlewareFunc{mw},
	}, func(sub *routing.Router) {
		sub.Get("/test", func(ctx *routing.Context) error {
			return ctx.Text(http.StatusOK, "ok")
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/test", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if !called {
		t.Fatal("expected group middleware to be called")
	}
}

func TestNestedGroups(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.GroupWith(routing.GroupAttributes{Prefix: "/api"}, func(api *routing.Router) {
		api.GroupWith(routing.GroupAttributes{Prefix: "/v1"}, func(v1 *routing.Router) {
			v1.Get("/users", func(ctx *routing.Context) error {
				return ctx.Text(http.StatusOK, "v1-users")
			})
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	if rec.Body.String() != "v1-users" {
		t.Fatalf("expected 'v1-users', got %q", rec.Body.String())
	}
}

func TestGroupWithNamePrefix(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.GroupWith(routing.GroupAttributes{
		Prefix:     "/api",
		NamePrefix: "api",
	}, func(sub *routing.Router) {
		route := sub.Get("/users", func(ctx *routing.Context) error { return nil })
		route.Name("users.index")

		if route.GetName() != "api.users.index" {
			t.Fatalf("expected 'api.users.index', got %q", route.GetName())
		}
	})
}

func TestNestedGroupNamePrefix(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.GroupWith(routing.GroupAttributes{
		Prefix:     "/api",
		NamePrefix: "api",
	}, func(api *routing.Router) {
		api.GroupWith(routing.GroupAttributes{
			Prefix:     "/v1",
			NamePrefix: "v1",
		}, func(v1 *routing.Router) {
			route := v1.Get("/users", func(ctx *routing.Context) error { return nil })
			route.Name("users.index")

			if route.GetName() != "api.v1.users.index" {
				t.Fatalf("expected 'api.v1.users.index', got %q", route.GetName())
			}
		})
	})
}

func TestGroupWithWhere(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.GroupWith(routing.GroupAttributes{
		Prefix: "/api",
		Where:  map[string]string{"id": `[0-9]+`},
	}, func(sub *routing.Router) {
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
		t.Fatalf("expected 404 for non-numeric id, got %d", rec.Code)
	}
}

func TestMergeGroupEmpty(t *testing.T) {
	t.Parallel()

	parent := routing.GroupAttributes{}
	child := routing.GroupAttributes{}

	merged := routing.MergeGroup(parent, child)

	if merged.Prefix != "" || merged.Domain != "" || merged.NamePrefix != "" {
		t.Fatal("expected all empty on merge of empty groups")
	}
}
