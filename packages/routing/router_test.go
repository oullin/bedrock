package routing

import (
	"errors"
	"testing"
)

// Translation of the dispatcher half of upstream/framework
// tests/Routing/RoutingRouteTest.php — the parts that need a Router instance.
// RouteRegistrarTest::testCanRegisterGetRouteWithClosureAction
// RouteRegistrarTest::testCanRegisterPostRouteWithClosureAction
// RouteRegistrarTest::testCanRegisterAnyRouteWithClosureAction
// RouteRegistrarTest::testCanRegisterMatchRouteWithClosureAction
// RouteRegistrarTest::testFallbackRoute
// RouteRegistrarTest::testSetFallbackRoute
// RouteRegistrarTest::testCanRegisterGroupWithPrefix
// RouteRegistrarTest::testCanRegisterGroupWithNamePrefix
// RouteRegistrarTest::testCanRegisterGroupWithDomain
// RouteRegistrarTest::testPushMiddlewareToGroup
// RouteRegistrarTest::testCanRemoveMiddlewareFromGroup
// RoutingRouteTest::testBasicDispatchingOfRoutes
// RoutingRouteTest::testRouterPatternSetting
// RoutingRouteTest::testMiddlewarePrioritySorting
// RoutingRouteTest::testGroupMerging
// RoutingRouteTest::testCurrentRouteUses
// RoutingRouteTest::testMergingControllerUses

func TestRouter_Registration(t *testing.T) {
	t.Run("test_get_post_put_patch_delete", func(t *testing.T) {
		r := NewRouter(nil, nil)
		r.Get("/users", func() {})
		r.Post("/users", func() {})
		r.Put("/users/{id}", func() {})
		r.Patch("/users/{id}", func() {})
		r.Delete("/users/{id}", func() {})
		r.Options("/users", func() {})

		if r.GetRoutes().(*RouteCollection).Count() != 6 {
			t.Errorf("count = %d", r.GetRoutes().(*RouteCollection).Count())
		}
	})

	t.Run("test_match_registers_multiple_methods", func(t *testing.T) {
		r := NewRouter(nil, nil)
		route := r.Match([]string{"PUT", "PATCH"}, "/users/{id}", func() {})
		methods := route.Methods()

		if len(methods) != 2 || methods[0] != "PUT" || methods[1] != "PATCH" {
			t.Errorf("methods = %v", methods)
		}
	})

	t.Run("test_any_registers_all_verbs", func(t *testing.T) {
		r := NewRouter(nil, nil)
		route := r.Any("/x", func() {})

		if len(route.Methods()) < 6 {
			t.Errorf("methods = %v", route.Methods())
		}
	})

	t.Run("test_fallback", func(t *testing.T) {
		r := NewRouter(nil, nil)
		fb := r.Fallback(func() {})

		if !fb.IsFallback {
			t.Error("fallback flag not set")
		}
	})
}

func TestRouter_Group(t *testing.T) {
	t.Run("test_group_applies_prefix", func(t *testing.T) {
		r := NewRouter(nil, nil)
		r.Group(map[string]any{"prefix": "api"}, func(r *Router) {
			r.Get("/users", func() {})
		})
		route := r.GetRoutes().GetRoutes()[0]

		if route.Uri != "api/users" {
			t.Errorf("uri = %q, want api/users", route.Uri)
		}
	})

	t.Run("test_nested_group_concatenates_prefix", func(t *testing.T) {
		r := NewRouter(nil, nil)
		r.Group(map[string]any{"prefix": "api"}, func(r *Router) {
			r.Group(map[string]any{"prefix": "v1"}, func(r *Router) {
				r.Get("/users", func() {})
			})
		})
		route := r.GetRoutes().GetRoutes()[0]

		if route.Uri != "api/v1/users" {
			t.Errorf("uri = %q, want api/v1/users", route.Uri)
		}
	})

	t.Run("test_group_applies_name_prefix", func(t *testing.T) {
		r := NewRouter(nil, nil)
		r.Group(map[string]any{"as": "users."}, func(r *Router) {
			r.Get("/users", func() {}).Name("index")
		})
		route := r.GetRoutes().GetRoutes()[0]

		if route.GetName() != "users.index" {
			t.Errorf("name = %q", route.GetName())
		}
	})

	t.Run("test_group_stack_pops", func(t *testing.T) {
		r := NewRouter(nil, nil)
		r.Group(map[string]any{"prefix": "api"}, func(r *Router) {})

		if r.HasGroupStack() {
			t.Error("group stack should be empty after group returns")
		}
	})

	// RoutingRouteTest::testMergingControllerUses
	t.Run("test_group_controller_merges_with_method_action", func(t *testing.T) {
		r := NewRouter(nil, nil)
		r.Group(map[string]any{"controller": "UserController"}, func(r *Router) {
			r.Get("/users/{user}", "show")
		})
		route := r.GetRoutes().GetRoutes()[0]

		if route.GetActionName() != "UserController@show" {
			t.Errorf("action = %q", route.GetActionName())
		}
	})
}

func TestRouter_Dispatch(t *testing.T) {
	t.Run("test_dispatch_runs_handler", func(t *testing.T) {
		r := NewRouter(nil, nil)
		called := false
		r.Get("/hello", func() { called = true })
		_, err := r.Dispatch(fakeRequest{method: "GET", path: "/hello"})

		if err != nil {
			t.Fatal(err)
		}

		if !called {
			t.Error("handler not called")
		}
	})

	t.Run("test_dispatch_returns_value", func(t *testing.T) {
		r := NewRouter(nil, nil)
		r.Get("/answer", func() any { return 42 })
		got, err := r.Dispatch(fakeRequest{method: "GET", path: "/answer"})

		if err != nil {
			t.Fatal(err)
		}

		if got != 42 {
			t.Errorf("got %v, want 42", got)
		}
	})

	t.Run("test_dispatch_returns_not_found", func(t *testing.T) {
		r := NewRouter(nil, nil)
		_, err := r.Dispatch(fakeRequest{method: "GET", path: "/missing"})

		if !errors.Is(err, ErrRouteNotFound) {
			t.Errorf("err = %v", err)
		}
	})

	t.Run("test_dispatch_method_not_allowed", func(t *testing.T) {
		r := NewRouter(nil, nil)
		r.Get("/users", func() {})
		_, err := r.Dispatch(fakeRequest{method: "POST", path: "/users"})

		var mna *MethodNotAllowedError

		if !errors.As(err, &mna) {
			t.Errorf("err = %v", err)
		}
	})

	t.Run("test_current_route_after_dispatch", func(t *testing.T) {
		r := NewRouter(nil, nil)
		r.Get("/users", func() {}).Name("users.index")
		_, _ = r.Dispatch(fakeRequest{method: "GET", path: "/users"})

		if r.CurrentRouteName() != "users.index" {
			t.Errorf("current route name = %q", r.CurrentRouteName())
		}

		if !r.Is("users.*") {
			t.Error("Is(users.*) should be true")
		}
	})

	// RoutingRouteTest::testCurrentRouteUses
	t.Run("test_current_route_uses_matches_controller_action", func(t *testing.T) {
		r := NewRouter(nil, nil)
		r.Get("/users/{user}", "UserController@show")
		_, err := r.Dispatch(fakeRequest{method: "GET", path: "/users/42"})

		if err != nil {
			t.Fatal(err)
		}

		if !r.CurrentRouteUses("UserController@show") {
			t.Errorf("current route action = %q", r.CurrentRouteAction())
		}
	})
}

func TestRouter_Patterns(t *testing.T) {
	t.Run("test_pattern_applies_to_routes", func(t *testing.T) {
		r := NewRouter(nil, nil)
		r.Pattern("id", "[0-9]+")
		route := r.Get("/users/{id}", func() {})

		if route.Wheres["id"] != "[0-9]+" {
			t.Errorf("wheres = %v", route.Wheres)
		}
	})
}

func TestRouter_AliasMiddleware(t *testing.T) {
	t.Run("test_alias_middleware", func(t *testing.T) {
		r := NewRouter(nil, nil)
		r.AliasMiddleware("auth", "App\\Http\\Middleware\\Authenticate")

		if r.GetMiddleware()["auth"] != "App\\Http\\Middleware\\Authenticate" {
			t.Error("alias not stored")
		}
	})

	t.Run("test_middleware_group_push_prepend_remove", func(t *testing.T) {
		r := NewRouter(nil, nil)
		r.MiddlewareGroup("web", []any{"a", "b"})
		r.PushMiddlewareToGroup("web", "c")
		r.PrependMiddlewareToGroup("web", "z")
		got := r.GetMiddlewareGroups()["web"]

		if len(got) != 4 || got[0] != "z" || got[3] != "c" {
			t.Errorf("group = %v", got)
		}

		r.RemoveMiddlewareFromGroup("web", "b")
		got = r.GetMiddlewareGroups()["web"]

		for _, m := range got {
			if m == "b" {
				t.Error("b should be removed")
			}
		}
	})
}

func TestRouter_Has(t *testing.T) {
	t.Run("test_has_named_route", func(t *testing.T) {
		r := NewRouter(nil, nil)
		r.Get("/x", func() {}).Name("x.index")

		if !r.Has("x.index") {
			t.Error("Has(x.index) should be true")
		}

		if r.Has("x.missing") {
			t.Error("Has(x.missing) should be false")
		}
	})
}

func TestSortedMiddleware(t *testing.T) {
	// Translation of upstream/framework tests/Routing/RoutingSortedMiddlewareTest.php.
	t.Run("test_priority_order_is_respected", func(t *testing.T) {
		priority := []string{"first", "second", "third"}
		input := []any{"third", "first", "second"}
		out := NewSortedMiddleware(priority, input)

		if len(out) != 3 || out[0] != "first" || out[1] != "second" || out[2] != "third" {
			t.Errorf("sorted = %v", out)
		}
	})

	t.Run("test_unprioritized_middleware_keeps_order", func(t *testing.T) {
		priority := []string{}
		input := []any{"a", "b", "c"}
		out := NewSortedMiddleware(priority, input)

		if len(out) != 3 || out[0] != "a" || out[1] != "b" || out[2] != "c" {
			t.Errorf("sorted = %v", out)
		}
	})

	t.Run("test_dedup", func(t *testing.T) {
		out := NewSortedMiddleware(nil, []any{"a", "b", "a", "c"})

		if len(out) != 3 {
			t.Errorf("sorted = %v", out)
		}
	})
}

func TestMiddlewareNameResolver(t *testing.T) {
	t.Run("test_resolve_alias", func(t *testing.T) {
		out := MiddlewareNameResolver{}.Resolve("auth", map[string]any{"auth": "Authenticate"}, nil)

		if len(out) != 1 || out[0] != "Authenticate" {
			t.Errorf("out = %v", out)
		}
	})

	t.Run("test_resolve_alias_with_params", func(t *testing.T) {
		out := MiddlewareNameResolver{}.Resolve("throttle:60,1", map[string]any{"throttle": "ThrottleClass"}, nil)

		if len(out) != 1 || out[0] != "ThrottleClass:60,1" {
			t.Errorf("out = %v", out)
		}
	})

	t.Run("test_resolve_group", func(t *testing.T) {
		out := MiddlewareNameResolver{}.Resolve("web",
			map[string]any{"auth": "Authenticate"},
			map[string][]any{"web": {"auth", "csrf"}})

		if len(out) != 2 || out[0] != "Authenticate" || out[1] != "csrf" {
			t.Errorf("out = %v", out)
		}
	})
}

func TestPipeline(t *testing.T) {
	t.Run("test_pipeline_order", func(t *testing.T) {
		var trace []string
		out := NewPipeline().
			Send("start").
			Through([]func(any, func(any) any) any{
				func(p any, next func(any) any) any {
					trace = append(trace, "a-in")
					r := next(p.(string) + "->a")
					trace = append(trace, "a-out")

					return r
				},
				func(p any, next func(any) any) any {
					trace = append(trace, "b-in")
					r := next(p.(string) + "->b")
					trace = append(trace, "b-out")

					return r
				},
			}).
			Then(func(p any) any { return p.(string) + "->dest" })

		if out != "start->a->b->dest" {
			t.Errorf("out = %q", out)
		}

		want := []string{"a-in", "b-in", "b-out", "a-out"}

		if len(trace) != len(want) {
			t.Fatalf("trace = %v", trace)
		}

		for i, w := range want {
			if trace[i] != w {
				t.Errorf("trace[%d] = %q, want %q", i, trace[i], w)
			}
		}
	})
}

func TestRouteGroup_Merge(t *testing.T) {
	t.Run("test_merge_prefix", func(t *testing.T) {
		out := MergeRouteGroup(
			map[string]any{"prefix": "v1"},
			map[string]any{"prefix": "api"},
			true,
		)

		if out["prefix"] != "api/v1" {
			t.Errorf("prefix = %v", out["prefix"])
		}
	})

	t.Run("test_merge_as", func(t *testing.T) {
		out := MergeRouteGroup(
			map[string]any{"as": "index"},
			map[string]any{"as": "users."},
			true,
		)

		if out["as"] != "users.index" {
			t.Errorf("as = %v", out["as"])
		}
	})

	t.Run("test_new_domain_overrides_old", func(t *testing.T) {
		out := MergeRouteGroup(
			map[string]any{"domain": "api.example.com"},
			map[string]any{"domain": "old.example.com"},
			true,
		)

		if out["domain"] != "api.example.com" {
			t.Errorf("domain = %v", out["domain"])
		}
	})
}

func TestFiltersControllerMiddleware(t *testing.T) {
	t.Run("test_only_excludes_others", func(t *testing.T) {
		opts := map[string]any{"only": []string{"index", "show"}}

		if MethodExcludedByOptions("create", opts) != true {
			t.Error("create should be excluded")
		}

		if MethodExcludedByOptions("show", opts) != false {
			t.Error("show should not be excluded")
		}
	})

	t.Run("test_except_excludes_named", func(t *testing.T) {
		opts := map[string]any{"except": []string{"destroy"}}

		if MethodExcludedByOptions("destroy", opts) != true {
			t.Error("destroy should be excluded")
		}

		if MethodExcludedByOptions("index", opts) != false {
			t.Error("index should not be excluded")
		}
	})
}
