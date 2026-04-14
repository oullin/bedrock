package routing

import (
	"errors"
	"testing"
)

// Translation of upstream/framework tests/Routing/RouteCollectionTest.php.

func TestRouteCollection_Add(t *testing.T) {
	t.Run("test_add_returns_route", func(t *testing.T) {
		c := NewRouteCollection()
		r := NewRoute("GET", "/foo", func() {})
		got := c.Add(r)

		if got != r {
			t.Error("Add should return the same route")
		}

		if c.Count() != 1 {
			t.Errorf("count = %d, want 1", c.Count())
		}
	})

	t.Run("test_add_indexes_by_method", func(t *testing.T) {
		c := NewRouteCollection()
		c.Add(NewRoute("GET", "/foo", func() {}))
		c.Add(NewRoute("POST", "/foo", func() {}))

		if len(c.Get("GET")) != 1 {
			t.Errorf("GET count = %d", len(c.Get("GET")))
		}

		if len(c.Get("POST")) != 1 {
			t.Errorf("POST count = %d", len(c.Get("POST")))
		}
	})

	t.Run("test_get_with_empty_string_returns_all", func(t *testing.T) {
		c := NewRouteCollection()
		c.Add(NewRoute("GET", "/foo", func() {}))
		c.Add(NewRoute("POST", "/bar", func() {}))

		if len(c.Get("")) != 2 {
			t.Errorf("get all = %d", len(c.Get("")))
		}
	})
}

func TestRouteCollection_NameLookups(t *testing.T) {
	t.Run("test_get_by_name", func(t *testing.T) {
		c := NewRouteCollection()
		r := NewRoute("GET", "/users", func() {}).Name("users.index")
		c.Add(r)

		if c.GetByName("users.index") != r {
			t.Error("GetByName failed")
		}

		if !c.HasNamedRoute("users.index") {
			t.Error("HasNamedRoute false")
		}

		if c.GetByName("missing") != nil {
			t.Error("GetByName(missing) should be nil")
		}
	})

	t.Run("test_refresh_name_lookups", func(t *testing.T) {
		c := NewRouteCollection()
		r := NewRoute("GET", "/users", func() {})
		c.Add(r)
		r.Name("users.index")
		c.RefreshNameLookups()

		if c.GetByName("users.index") != r {
			t.Error("after refresh, name lookup should resolve")
		}
	})

	t.Run("test_first_name_wins", func(t *testing.T) {
		c := NewRouteCollection()
		r1 := NewRoute("GET", "/a", func() {}).Name("dup")
		r2 := NewRoute("GET", "/b", func() {}).Name("dup")
		c.Add(r1)
		c.Add(r2)

		if c.GetByName("dup") != r1 {
			t.Error("first registration should win")
		}
	})
}

func TestRouteCollection_ActionLookups(t *testing.T) {
	t.Run("test_get_by_action", func(t *testing.T) {
		c := NewRouteCollection()
		r := NewRoute("GET", "/users", "App\\Http\\Controllers\\UserController@index")
		c.Add(r)

		if c.GetByAction("App\\Http\\Controllers\\UserController@index") != r {
			t.Error("GetByAction failed")
		}
	})
}

func TestRouteCollection_Match(t *testing.T) {
	t.Run("test_match_returns_bound_route", func(t *testing.T) {
		c := NewRouteCollection()
		c.Add(NewRoute("GET", "/users/{user}", func() {}))
		got, err := c.Match(fakeRequest{method: "GET", path: "/users/42"})

		if err != nil {
			t.Fatal(err)
		}

		if got.Parameter("user", "") != "42" {
			t.Errorf("user = %q", got.Parameter("user", ""))
		}
	})

	t.Run("test_match_returns_not_found", func(t *testing.T) {
		c := NewRouteCollection()
		c.Add(NewRoute("GET", "/users", func() {}))
		_, err := c.Match(fakeRequest{method: "GET", path: "/missing"})

		if !errors.Is(err, ErrRouteNotFound) {
			t.Errorf("err = %v, want ErrRouteNotFound", err)
		}
	})

	t.Run("test_match_returns_method_not_allowed", func(t *testing.T) {
		c := NewRouteCollection()
		c.Add(NewRoute("GET", "/users", func() {}))
		_, err := c.Match(fakeRequest{method: "POST", path: "/users"})

		var mna *MethodNotAllowedError

		if !errors.As(err, &mna) {
			t.Fatalf("err = %v, want MethodNotAllowedError", err)
		}
		// HEAD is auto-added to GET routes, so both GET and HEAD should appear.
		hasGet := false

		for _, m := range mna.Allowed {
			if m == "GET" {
				hasGet = true
			}
		}

		if !hasGet {
			t.Errorf("allowed = %v, want GET present", mna.Allowed)
		}
	})

	t.Run("test_fallback_route_matches_last", func(t *testing.T) {
		c := NewRouteCollection()
		concrete := NewRoute("GET", "/users", func() {})
		fallback := NewRoute("GET", "/{any}", func() {}).Where("any", ".*").Fallback()
		c.Add(concrete)
		c.Add(fallback)
		got, err := c.Match(fakeRequest{method: "GET", path: "/users"})

		if err != nil {
			t.Fatal(err)
		}

		if got != concrete {
			t.Error("concrete route should win over fallback")
		}

		got, err = c.Match(fakeRequest{method: "GET", path: "/anything"})

		if err != nil {
			t.Fatal(err)
		}

		if got != fallback {
			t.Error("fallback should match unhandled paths")
		}
	})
}

func TestCompiledRouteCollection(t *testing.T) {
	t.Run("test_compiled_collection_matches", func(t *testing.T) {
		routes := []*Route{
			NewRoute("GET", "/users", func() {}).Name("users.index"),
			NewRoute("POST", "/users", func() {}),
		}
		c := NewCompiledRouteCollection(routes, nil)

		if c.Count() != 2 {
			t.Errorf("count = %d", c.Count())
		}

		if c.GetByName("users.index") != routes[0] {
			t.Error("name lookup failed")
		}

		got, err := c.Match(fakeRequest{method: "GET", path: "/users"})

		if err != nil || got != routes[0] {
			t.Errorf("match failed: %v %v", got, err)
		}
	})

	t.Run("test_compiled_collection_method_not_allowed", func(t *testing.T) {
		c := NewCompiledRouteCollection([]*Route{
			NewRoute("GET", "/users", func() {}),
		}, nil)
		_, err := c.Match(fakeRequest{method: "POST", path: "/users"})

		var mna *MethodNotAllowedError

		if !errors.As(err, &mna) {
			t.Errorf("err = %v", err)
		}
	})
}

// Compile-time check that both collection types satisfy the interface.
var _ RouteCollectionInterface = (*RouteCollection)(nil)
var _ RouteCollectionInterface = (*CompiledRouteCollection)(nil)
