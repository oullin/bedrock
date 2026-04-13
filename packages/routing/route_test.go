package routing_test

import (
	"net/http"
	"testing"

	"github.com/bedrock/packages/routing"
)

func TestRouteName(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	route := r.Get("/users", func(ctx *routing.Context) error { return nil })
	route.Name("users.index")

	if route.GetName() != "users.index" {
		t.Fatalf("expected 'users.index', got %q", route.GetName())
	}
}

func TestRouteMethods(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	route := r.Get("/users", func(ctx *routing.Context) error { return nil })

	methods := route.Methods()

	if len(methods) != 1 || methods[0] != http.MethodGet {
		t.Fatalf("expected [GET], got %v", methods)
	}
}

func TestRouteURI(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	route := r.Get("/users/{id}", func(ctx *routing.Context) error { return nil })

	if route.URI() != "/users/{id}" {
		t.Fatalf("expected '/users/{id}', got %q", route.URI())
	}
}

func TestRouteHasParameters(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)

	withParams := r.Get("/users/{id}", func(ctx *routing.Context) error { return nil })

	if !withParams.HasParameters() {
		t.Fatal("expected route to have parameters")
	}

	withoutParams := r.Get("/about", func(ctx *routing.Context) error { return nil })

	if withoutParams.HasParameters() {
		t.Fatal("expected route to have no parameters")
	}
}

func TestRouteParameterNames(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	route := r.Get("/posts/{id}/comments/{cid}", func(ctx *routing.Context) error { return nil })

	names := route.ParameterNames()

	if len(names) != 2 || names[0] != "id" || names[1] != "cid" {
		t.Fatalf("expected [id cid], got %v", names)
	}
}

func TestRouteWhereConstraint(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	route := r.Get("/users/{id}", func(ctx *routing.Context) error { return nil })
	route.Where("id", `[0-9]+`)

	constraints := route.GetConstraints()

	if constraints["id"] != `[0-9]+` {
		t.Fatalf("expected numeric constraint, got %q", constraints["id"])
	}

	if !route.MatchesConstraints(map[string]string{"id": "42"}) {
		t.Fatal("expected numeric value to match")
	}

	if route.MatchesConstraints(map[string]string{"id": "abc"}) {
		t.Fatal("expected alphabetic value to not match")
	}
}

func TestRouteWhereNumber(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	route := r.Get("/users/{id}", func(ctx *routing.Context) error { return nil })
	route.WhereNumber("id")

	if !route.MatchesConstraints(map[string]string{"id": "123"}) {
		t.Fatal("expected numeric value to match")
	}

	if route.MatchesConstraints(map[string]string{"id": "abc"}) {
		t.Fatal("expected alphabetic value to not match")
	}
}

func TestRouteWhereAlpha(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	route := r.Get("/categories/{slug}", func(ctx *routing.Context) error { return nil })
	route.WhereAlpha("slug")

	if !route.MatchesConstraints(map[string]string{"slug": "books"}) {
		t.Fatal("expected alpha value to match")
	}

	if route.MatchesConstraints(map[string]string{"slug": "books123"}) {
		t.Fatal("expected alphanumeric value to not match alpha constraint")
	}
}

func TestRouteWhereAlphaNumeric(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	route := r.Get("/items/{code}", func(ctx *routing.Context) error { return nil })
	route.WhereAlphaNumeric("code")

	if !route.MatchesConstraints(map[string]string{"code": "abc123"}) {
		t.Fatal("expected alphanumeric value to match")
	}

	if route.MatchesConstraints(map[string]string{"code": "abc-123"}) {
		t.Fatal("expected hyphenated value to not match")
	}
}

func TestRouteWhereIn(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	route := r.Get("/lang/{locale}", func(ctx *routing.Context) error { return nil })
	route.WhereIn("locale", []string{"en", "es", "fr"})

	if !route.MatchesConstraints(map[string]string{"locale": "en"}) {
		t.Fatal("expected 'en' to match")
	}

	if route.MatchesConstraints(map[string]string{"locale": "de"}) {
		t.Fatal("expected 'de' to not match")
	}
}

func TestRouteWhereUUID(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	route := r.Get("/items/{uuid}", func(ctx *routing.Context) error { return nil })
	route.WhereUUID("uuid")

	if !route.MatchesConstraints(map[string]string{"uuid": "550e8400-e29b-41d4-a716-446655440000"}) {
		t.Fatal("expected UUID to match")
	}

	if route.MatchesConstraints(map[string]string{"uuid": "not-a-uuid"}) {
		t.Fatal("expected non-UUID to not match")
	}
}

func TestRouteChaining(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	route := r.Get("/users/{id}", func(ctx *routing.Context) error { return nil }).
		Name("users.show").
		WhereNumber("id").
		Domain("api.example.com")

	if route.GetName() != "users.show" {
		t.Fatalf("expected 'users.show', got %q", route.GetName())
	}

	if route.GetDomain() != "api.example.com" {
		t.Fatalf("expected 'api.example.com', got %q", route.GetDomain())
	}

	if !route.MatchesConstraints(map[string]string{"id": "42"}) {
		t.Fatal("expected constraint to match")
	}
}

func TestRouteFallback(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	route := r.Get("/fallback", func(ctx *routing.Context) error { return nil })

	if route.IsFallback() {
		t.Fatal("expected route to not be fallback by default")
	}

	route.SetFallback(true)

	if !route.IsFallback() {
		t.Fatal("expected route to be fallback")
	}
}

func TestRouteDefaults(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	route := r.Get("/page/{page}", func(ctx *routing.Context) error { return nil })
	route.Defaults(map[string]string{"page": "1"})

	defaults := route.GetDefaults()

	if defaults["page"] != "1" {
		t.Fatalf("expected default page=1, got %q", defaults["page"])
	}
}

func TestRouteMiddlewareAppend(t *testing.T) {
	t.Parallel()

	mw1 := func(next http.Handler) http.Handler { return next }
	mw2 := func(next http.Handler) http.Handler { return next }

	r := routing.New(nil)
	route := r.Get("/test", func(ctx *routing.Context) error { return nil })
	route.Middleware(mw1, mw2)

	if len(route.GetMiddleware()) != 2 {
		t.Fatalf("expected 2 middleware, got %d", len(route.GetMiddleware()))
	}
}

func TestRouteConstraintsWithMultipleParams(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	route := r.Get("/posts/{id}/comments/{cid}", func(ctx *routing.Context) error { return nil })
	route.WhereNumber("id", "cid")

	if !route.MatchesConstraints(map[string]string{"id": "1", "cid": "2"}) {
		t.Fatal("expected both numeric values to match")
	}

	if route.MatchesConstraints(map[string]string{"id": "1", "cid": "abc"}) {
		t.Fatal("expected mixed values to not match")
	}
}

func TestRouteConstraintsWithMissingParam(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	route := r.Get("/users/{id}", func(ctx *routing.Context) error { return nil })
	route.WhereNumber("id")

	if route.MatchesConstraints(map[string]string{}) {
		t.Fatal("expected missing param to not match")
	}
}

func TestRouteConstraintsWithDefaultForMissingParam(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	route := r.Get("/page/{page}", func(ctx *routing.Context) error { return nil })
	route.WhereNumber("page")
	route.Defaults(map[string]string{"page": "1"})

	if !route.MatchesConstraints(map[string]string{}) {
		t.Fatal("expected missing param with default to pass")
	}
}

func TestRouteNoConstraints(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	route := r.Get("/any/{value}", func(ctx *routing.Context) error { return nil })

	if !route.MatchesConstraints(map[string]string{"value": "anything-goes"}) {
		t.Fatal("expected unconstrained param to match anything")
	}
}

func TestRouteGetPrefix(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	r.Group("/api", nil, func(sub *routing.Router) {
		route := sub.Get("/users", func(ctx *routing.Context) error { return nil })

		if route.GetPrefix() != "/api" {
			t.Fatalf("expected prefix '/api', got %q", route.GetPrefix())
		}
	})
}

func TestRouteWhereChaining(t *testing.T) {
	t.Parallel()

	r := routing.New(nil)
	route := r.Get("/posts/{id}/tags/{slug}", func(ctx *routing.Context) error { return nil })
	route.Where("id", `[0-9]+`).Where("slug", `[a-z-]+`)

	if !route.MatchesConstraints(map[string]string{"id": "42", "slug": "hello-world"}) {
		t.Fatal("expected chained constraints to match")
	}

	if route.MatchesConstraints(map[string]string{"id": "abc", "slug": "hello-world"}) {
		t.Fatal("expected id constraint to fail")
	}
}
