package routing

import (
	"errors"
	"testing"
)

// Ref: @bedrock/code-0393
// RouteCollectionTest::testRouteCollectionCanAddRoute
// RouteCollectionTest::testRouteCollectionAddReturnsTheRoute
// RouteCollectionTest::testRouteCollectionCanRetrieveByName
// RouteCollectionTest::testRouteCollectionCanRetrieveByAction
// RouteCollectionTest::testRouteCollectionCanRetrieveByMethod
// RouteCollectionTest::testRouteCollectionCanGetAllRoutes
// RouteCollectionTest::testRouteCollectionCanGetRoutesByName
// RouteCollectionTest::testRouteCollectionCanGetRoutesByMethod
// RouteCollectionTest::testRouteCollectionCanRefreshNameLookups
// RouteCollectionTest::testCannotCacheDuplicateRouteNames
// RouteCollectionTest::testRouteCollectionCanHandleSameRoute
// RouteCollectionTest::testRouteCollectionCleansUpOverwrittenRoutes
// RouteCollectionTest::testRouteCollectionCanGetIterator
// RouteCollectionTest::testRouteCollectionCanGetIteratorWhenEmpty
// RouteCollectionTest::testRouteCollectionCanGetIteratorWhenRouteAreAdded
// RouteCollectionTest::testRouteCollectionRequestMethodNotAllowed
// RouteCollectionTest::testRouteCollectionDontMatchNonMatchingDoubleSlashes
// RouteCollectionTest::testOverlappingRoutesMatchesFirstRoute
// RouteCollectionTest::testHasNameRouteMethod
// RouteCollectionTest::testPrependsRoutesWithDomain
// RouteCollectionTest::testToSymfonyRouteCollection

func TestRouteCollection_Add(t *testing.T) {
	// RouteCollectionTest::testRouteCollectionAddReturnsTheRoute
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

	// RouteCollectionTest::testRouteCollectionCanRetrieveByMethod
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

	// RouteCollectionTest::testRouteCollectionCanGetRoutesByMethod
	t.Run("test_get_routes_by_method_returns_registered_method_indexes", func(t *testing.T) {
		c := NewRouteCollection()
		index := NewRoute("GET", "/foo/index", func() {}).Name("foo.index")
		show := NewRoute("GET", "/foo/show", func() {}).Name("foo.show")
		create := NewRoute("POST", "/bar", func() {}).Name("bar.create")
		c.Add(index)
		c.Add(show)
		c.Add(create)

		got := c.GetRoutesByMethod()

		if len(got["GET"]) != 2 || got["GET"][0] != index || got["GET"][1] != show {
			t.Errorf("GET routes = %v", got["GET"])
		}

		if len(got["HEAD"]) != 2 || got["HEAD"][0] != index || got["HEAD"][1] != show {
			t.Errorf("HEAD routes = %v", got["HEAD"])
		}

		if len(got["POST"]) != 1 || got["POST"][0] != create {
			t.Errorf("POST routes = %v", got["POST"])
		}
	})

	// RouteCollectionTest::testRouteCollectionCanGetAllRoutes
	t.Run("test_get_with_empty_string_returns_all", func(t *testing.T) {
		c := NewRouteCollection()
		c.Add(NewRoute("GET", "/foo", func() {}))
		c.Add(NewRoute("POST", "/bar", func() {}))

		if len(c.Get("")) != 2 {
			t.Errorf("get all = %d", len(c.Get("")))
		}
	})

	// RouteCollectionTest::testRouteCollectionCanHandleSameRoute
	t.Run("test_same_route_overwrites_in_place", func(t *testing.T) {
		c := NewRouteCollection()
		first := NewRoute("GET", "/foo", func() {})
		second := NewRoute("GET", "/foo", func() {})
		c.Add(first)
		c.Add(second)

		if c.Count() != 1 {
			t.Errorf("count = %d, want 1", c.Count())
		}

		got := c.Get("GET")

		if len(got) != 1 || got[0] != second {
			t.Errorf("GET routes = %v, want overwritten route", got)
		}

		all := c.GetRoutes()

		if len(all) != 1 || all[0] != second {
			t.Errorf("all routes = %v, want overwritten route", all)
		}
	})

	// RouteCollectionTest::testRouteCollectionCleansUpOverwrittenRoutes
	t.Run("test_refresh_lookups_cleans_up_overwritten_routes", func(t *testing.T) {
		c := NewRouteCollection()
		first := NewRoute("GET", "/product", map[string]any{
			"controller": "View@view",
			"as":         "routeA",
		})
		second := NewRoute("GET", "/product", map[string]any{
			"controller": "OverwrittenView@view",
			"as":         "overwrittenRouteA",
		})
		c.Add(first)
		c.Add(second)

		if c.GetByName("routeA") != first || c.GetByAction("View@view") != first {
			t.Fatal("stale lookups should exist before refresh, matching the upstream lookup-cache behavior")
		}

		c.RefreshNameLookups()
		c.RefreshActionLookups()

		if c.GetByName("routeA") != nil {
			t.Error("overwritten route name should be removed after refresh")
		}

		if c.GetByAction("View@view") != nil {
			t.Error("overwritten route action should be removed after refresh")
		}

		if c.GetByName("overwrittenRouteA") != second || c.GetByAction("OverwrittenView@view") != second {
			t.Error("replacement route lookups should remain after refresh")
		}
	})

	// RouteCollectionTest::testRouteCollectionCanGetIterator
	t.Run("test_get_routes_returns_copy_in_registration_order", func(t *testing.T) {
		c := NewRouteCollection()
		first := NewRoute("GET", "/a", func() {})
		second := NewRoute("POST", "/b", func() {})
		c.Add(first)
		c.Add(second)

		got := c.GetRoutes()

		if len(got) != 2 || got[0] != first || got[1] != second {
			t.Errorf("routes = %v", got)
		}

		got[0] = nil

		if c.GetRoutes()[0] != first {
			t.Error("GetRoutes should return a defensive copy")
		}
	})

	// RouteCollectionTest::testRouteCollectionCanGetIteratorWhenEmpty
	t.Run("test_get_routes_empty", func(t *testing.T) {
		c := NewRouteCollection()

		if len(c.GetRoutes()) != 0 {
			t.Errorf("routes = %v, want empty", c.GetRoutes())
		}
	})
}

func TestRouteCollection_NameLookups(t *testing.T) {
	// RouteCollectionTest::testRouteCollectionCanRetrieveByName
	// RouteCollectionTest::testHasNameRouteMethod
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

	// RouteCollectionTest::testRouteCollectionCanGetRoutesByName
	t.Run("test_get_routes_by_name_returns_registered_name_lookup", func(t *testing.T) {
		c := NewRouteCollection()
		index := NewRoute("GET", "/foo/index", func() {}).Name("foo_index")
		show := NewRoute("GET", "/foo/show", func() {}).Name("foo_show")
		create := NewRoute("POST", "/bar", func() {}).Name("bar_create")
		c.Add(index)
		c.Add(show)
		c.Add(create)

		got := c.GetRoutesByName()

		if len(got) != 3 {
			t.Fatalf("routes by name = %v", got)
		}

		if got["foo_index"] != index || got["foo_show"] != show || got["bar_create"] != create {
			t.Errorf("routes by name = %v", got)
		}
	})

	// RouteCollectionTest::testRouteCollectionCanRefreshNameLookups
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

	// RouteCollectionTest::testCannotCacheDuplicateRouteNames
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
	// RouteCollectionTest::testRouteCollectionCanRetrieveByAction
	t.Run("test_get_by_action", func(t *testing.T) {
		c := NewRouteCollection()
		r := NewRoute("GET", "/users", "App\\Http\\Controllers\\UserController@index")
		c.Add(r)

		if c.GetByAction("App\\Http\\Controllers\\UserController@index") != r {
			t.Error("GetByAction failed")
		}
	})

	t.Run("test_refresh_action_lookups_normalizes_leading_backslash", func(t *testing.T) {
		c := NewRouteCollection()
		r := NewRoute("GET", "/users", map[string]any{
			"controller": "\\App\\Http\\Controllers\\UserController@index",
		})
		c.Add(r)

		if c.GetByAction("App\\Http\\Controllers\\UserController@index") != r {
			t.Fatal("GetByAction should normalize leading backslashes on add")
		}

		c.RefreshActionLookups()

		if c.GetByAction("App\\Http\\Controllers\\UserController@index") != r {
			t.Fatal("GetByAction should normalize leading backslashes on refresh")
		}
	})
}

func TestRouteCollection_DomainHandling(t *testing.T) {
	// RouteCollectionTest::testPrependsRoutesWithDomain
	t.Run("test_prepends_routes_with_domain", func(t *testing.T) {
		c := NewRouteCollection()
		admin := NewRoute("GET", "/users", func() {}).Domain("admin.example.com")
		api := NewRoute("GET", "/users", func() {}).Domain("api.example.com")
		c.Add(admin)
		c.Add(api)

		if c.Count() != 2 {
			t.Fatalf("count = %d, want 2", c.Count())
		}

		if len(c.Get("GET")) != 2 {
			t.Fatalf("GET routes = %d, want 2", len(c.Get("GET")))
		}
	})
}

func TestRouteCollection_ToCompiledCollection(t *testing.T) {
	// RouteCollectionTest::testToSymfonyRouteCollection
	t.Run("test_to_symfony_route_collection", func(t *testing.T) {
		route := NewRoute("GET", "/users", map[string]any{
			"controller": "\\App\\Http\\Controllers\\UserController@index",
		}).Name("users.index")
		c := NewCompiledRouteCollection([]*Route{route}, nil)

		if c.Count() != 1 {
			t.Fatalf("count = %d, want 1", c.Count())
		}

		if c.GetByName("users.index") != route {
			t.Fatal("compiled collection should preserve named route lookups")
		}

		if c.GetByAction("App\\Http\\Controllers\\UserController@index") != route {
			t.Fatal("compiled collection should normalize controller lookups")
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

	// RouteCollectionTest::testRouteCollectionDontMatchNonMatchingDoubleSlashes
	t.Run("test_double_slash_path_does_not_match_single_slash_route", func(t *testing.T) {
		c := NewRouteCollection()
		c.Add(NewRoute("GET", "/foo", func() {}))
		_, err := c.Match(fakeRequest{method: "GET", path: "//foo"})

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

		if got == concrete {
			t.Error("match should return a request-scoped route instance")
		}

		if got.Uri != concrete.Uri {
			t.Error("concrete route should win over fallback")
		}

		got, err = c.Match(fakeRequest{method: "GET", path: "/anything"})

		if err != nil {
			t.Fatal(err)
		}

		if got == fallback {
			t.Error("match should return a request-scoped fallback instance")
		}

		if got.Uri != fallback.Uri {
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

		if err != nil {
			t.Errorf("match failed: %v %v", got, err)
		}

		if got == routes[0] {
			t.Error("match should return a request-scoped route instance")
		}

		if got.Uri != routes[0].Uri {
			t.Errorf("matched URI = %q, want %q", got.Uri, routes[0].Uri)
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
