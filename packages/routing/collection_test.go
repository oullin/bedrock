package routing_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/routing"
)

func TestCollectionAdd(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Get("/users", func(ctx *routing.Context) error { return nil })

	if r.Collection().Len() != 1 {
		t.Fatalf("expected 1 route, got %d", r.Collection().Len())
	}
}

func TestCollectionAddReturnsRoute(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	route := r.Get("/users", func(ctx *routing.Context) error { return nil })

	if route == nil {
		t.Fatal("expected route to be returned")
	}

	if route.URI() != "/users" {
		t.Fatalf("expected URI '/users', got %q", route.URI())
	}
}

func TestCollectionByName(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Get("/users", func(ctx *routing.Context) error { return nil }).Name("users.index")

	found, ok := r.Collection().ByName("users.index")

	if !ok || found.GetName() != "users.index" {
		t.Fatalf("expected to find route by name, got ok=%v name=%q", ok, found.GetName())
	}
}

func TestCollectionByNameMissing(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	_, ok := r.Collection().ByName("missing")

	if ok {
		t.Fatal("expected not found")
	}
}

func TestCollectionHasNamedRoute(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Get("/users", func(ctx *routing.Context) error { return nil }).Name("users.index")

	if !r.Collection().HasNamedRoute("users.index") {
		t.Fatal("expected named route to exist")
	}

	if r.Collection().HasNamedRoute("missing") {
		t.Fatal("expected missing route to not exist")
	}
}

func TestCollectionByMethod(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Get("/users", func(ctx *routing.Context) error { return nil })
	r.Post("/users", func(ctx *routing.Context) error { return nil })
	r.Get("/posts", func(ctx *routing.Context) error { return nil })

	gets := r.Collection().ByMethod(http.MethodGet)

	if len(gets) != 2 {
		t.Fatalf("expected 2 GET routes, got %d", len(gets))
	}

	posts := r.Collection().ByMethod(http.MethodPost)

	if len(posts) != 1 {
		t.Fatalf("expected 1 POST route, got %d", len(posts))
	}
}

func TestCollectionAll(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Get("/a", func(ctx *routing.Context) error { return nil })
	r.Post("/b", func(ctx *routing.Context) error { return nil })
	r.Put("/c", func(ctx *routing.Context) error { return nil })

	all := r.Collection().All()

	if len(all) != 3 {
		t.Fatalf("expected 3 routes, got %d", len(all))
	}
}

func TestCollectionLen(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)

	if r.Collection().Len() != 0 {
		t.Fatal("expected 0 routes initially")
	}

	r.Get("/a", func(ctx *routing.Context) error { return nil })

	if r.Collection().Len() != 1 {
		t.Fatalf("expected 1 route, got %d", r.Collection().Len())
	}
}

func TestCollectionMatch(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Get("/users", func(ctx *routing.Context) error { return nil }).Name("users.index")
	r.Get("/users/{id}", func(ctx *routing.Context) error { return nil }).Name("users.show")

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	route, err := r.Collection().Match(req)

	if err != nil {
		t.Fatalf("expected match, got error: %v", err)
	}

	if route.GetName() != "users.index" {
		t.Fatalf("expected 'users.index', got %q", route.GetName())
	}
}

func TestCollectionMatchWithParam(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Get("/users/{id}", func(ctx *routing.Context) error { return nil }).Name("users.show")

	req := httptest.NewRequest(http.MethodGet, "/users/42", nil)
	route, err := r.Collection().Match(req)

	if err != nil {
		t.Fatalf("expected match, got error: %v", err)
	}

	if route.GetName() != "users.show" {
		t.Fatalf("expected 'users.show', got %q", route.GetName())
	}
}

func TestCollectionMatchNotFound(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Get("/users", func(ctx *routing.Context) error { return nil })

	req := httptest.NewRequest(http.MethodGet, "/posts", nil)
	_, err := r.Collection().Match(req)

	if err != routing.ErrRouteNotFound {
		t.Fatalf("expected ErrRouteNotFound, got %v", err)
	}
}

func TestCollectionMatchMethodNotAllowed(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Post("/users", func(ctx *routing.Context) error { return nil })

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	_, err := r.Collection().Match(req)

	if err != routing.ErrMethodNotAllowed {
		t.Fatalf("expected ErrMethodNotAllowed, got %v", err)
	}
}

func TestCollectionRefreshNameLookups(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	route := r.Get("/users", func(ctx *routing.Context) error { return nil })

	if r.Collection().HasNamedRoute("users.index") {
		t.Fatal("expected no named route before naming")
	}

	route.Name("users.index")
	r.Collection().RefreshNameLookups()

	if !r.Collection().HasNamedRoute("users.index") {
		t.Fatal("expected named route after refresh")
	}
}

func TestCollectionInsertionOrderPreserved(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Get("/first", func(ctx *routing.Context) error { return nil }).Name("first")
	r.Get("/second", func(ctx *routing.Context) error { return nil }).Name("second")
	r.Get("/third", func(ctx *routing.Context) error { return nil }).Name("third")

	all := r.Collection().All()

	if len(all) != 3 {
		t.Fatalf("expected 3 routes, got %d", len(all))
	}

	if all[0].GetName() != "first" || all[2].GetName() != "third" {
		t.Fatalf("expected insertion order, got %v", []string{all[0].GetName(), all[1].GetName(), all[2].GetName()})
	}
}

func TestCollectionAnyRouteMatchesAllMethods(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Any("/catch-all", func(ctx *routing.Context) error { return nil }).Name("catch")

	for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete} {
		req := httptest.NewRequest(method, "/catch-all", nil)
		route, err := r.Collection().Match(req)

		if err != nil {
			t.Fatalf("%s: expected match, got error: %v", method, err)
		}

		if route.GetName() != "catch" {
			t.Fatalf("%s: expected 'catch', got %q", method, route.GetName())
		}
	}
}

func TestCollectionMultipleRoutesSameURI(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Get("/items", func(ctx *routing.Context) error { return nil }).Name("items.index")
	r.Post("/items", func(ctx *routing.Context) error { return nil }).Name("items.store")

	getReq := httptest.NewRequest(http.MethodGet, "/items", nil)
	route, err := r.Collection().Match(getReq)

	if err != nil || route.GetName() != "items.index" {
		t.Fatalf("GET: expected 'items.index', got %q err=%v", route.GetName(), err)
	}

	postReq := httptest.NewRequest(http.MethodPost, "/items", nil)
	route, err = r.Collection().Match(postReq)

	if err != nil || route.GetName() != "items.store" {
		t.Fatalf("POST: expected 'items.store', got %q err=%v", route.GetName(), err)
	}
}
