package routing

import (
	"reflect"
	"testing"

	"github.com/bedrock/packages/routing/contracts"
)

// tests/Routing/ImplicitRouteBindingTest.php that don't depend on Orm.
// ImplicitRouteBindingTest::test_it_can_resolve_the_implicit_model_route_bindings_for_the_given_route

// fakeUser is a UrlRoutable used by the implicit-binding tests.
type fakeUser struct {
	ID   string
	Slug string
}

// fakeContainer satisfies DependencyContainer for tests, minting fresh
// instances for the requested type via reflect.New.
type fakeContainer struct{}

func (u *fakeUser) GetRouteKey() any        { return u.ID }
func (u *fakeUser) GetRouteKeyName() string { return "id" }
func (u *fakeUser) ResolveRouteBinding(value, field string) (any, error) {
	if field == "slug" {
		if value == "alice" {
			return &fakeUser{ID: "1", Slug: "alice"}, nil
		}

		return nil, nil
	}

	if value == "1" {
		return &fakeUser{ID: "1", Slug: "alice"}, nil
	}

	return nil, nil
}
func (u *fakeUser) ResolveChildRouteBinding(childType, value, field string) (any, error) {
	return nil, nil
}

func (fakeContainer) MakeFor(t reflect.Type) (any, error) {
	if t.Kind() == reflect.Ptr {
		return reflect.New(t.Elem()).Interface(), nil
	}

	return reflect.New(t).Elem().Interface(), nil
}

// Compile-time check that *fakeUser satisfies UrlRoutable.
var _ contracts.UrlRoutable = (*fakeUser)(nil)

func TestImplicitRouteBinding(t *testing.T) {
	t.Run("test_resolves_url_routable_parameter", func(t *testing.T) {
		// A handler typed as *fakeUser triggers implicit binding on the {user} param.
		r := NewRoute("GET", "/users/{user}", func(user *fakeUser) {})
		_, _ = r.Bind(fakeRequest{path: "/users/1"})

		if err := (ImplicitRouteBinding{}).ResolveForRoute(fakeContainer{}, r); err != nil {
			t.Fatal(err)
		}

		bound := r.BoundModel("user")

		if bound == nil {
			t.Fatal("bound model nil")
		}

		u, ok := bound.(*fakeUser)

		if !ok {
			t.Fatalf("bound type = %T", bound)
		}

		if u.ID != "1" {
			t.Errorf("id = %s", u.ID)
		}
	})

	t.Run("test_unresolved_returns_not_found", func(t *testing.T) {
		r := NewRoute("GET", "/users/{user}", func(user *fakeUser) {})
		_, _ = r.Bind(fakeRequest{path: "/users/missing"})
		err := (ImplicitRouteBinding{}).ResolveForRoute(fakeContainer{}, r)

		if err == nil {
			t.Fatal("expected error")
		}

		if _, ok := err.(*ModelNotFoundError); !ok {
			t.Errorf("err = %v", err)
		}
	})

	t.Run("test_uses_binding_field", func(t *testing.T) {
		r := NewRoute("GET", "/users/{user:slug}", func(user *fakeUser) {})
		_, _ = r.Bind(fakeRequest{path: "/users/alice"})

		if err := (ImplicitRouteBinding{}).ResolveForRoute(fakeContainer{}, r); err != nil {
			t.Fatal(err)
		}

		u := r.BoundModel("user").(*fakeUser)

		if u.Slug != "alice" {
			t.Errorf("slug = %s", u.Slug)
		}
	})
}
