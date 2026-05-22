package routing

import "testing"

// Ref: @bedrock/code-0391
// RouteActionTest::test_it_can_detect_a_serialized_closure
// RoutingRouteTest::testRouteGetControllerClass
// RoutingRouteTest::testRouteFlushController
func TestRouteAction(t *testing.T) {
	t.Run("test_parse_action_with_callable", func(t *testing.T) {
		called := false
		fn := func() { called = true }
		a, err := ParseAction("/x", fn)

		if err != nil {
			t.Fatal(err)
		}

		if a.Uses == nil {
			t.Fatal("uses nil")
		}
		// Sanity: invoke and ensure the original closure runs.
		a.Uses.(func())()

		if !called {
			t.Error("closure not invoked")
		}
	})

	t.Run("test_parse_action_with_controller_string", func(t *testing.T) {
		a, err := ParseAction("/x", "App\\Http\\Controllers\\UserController@show")

		if err != nil {
			t.Fatal(err)
		}

		if a.Controller != "App\\Http\\Controllers\\UserController@show" {
			t.Errorf("controller = %q", a.Controller)
		}
	})

	t.Run("test_parse_action_invokable_string", func(t *testing.T) {
		a, err := ParseAction("/x", "App\\Http\\Controllers\\Invokable")

		if err != nil {
			t.Fatal(err)
		}

		want := "App\\Http\\Controllers\\Invokable@Invoke"

		if a.Controller != want {
			t.Errorf("controller = %q, want %q", a.Controller, want)
		}
	})

	t.Run("test_parse_action_with_nil_returns_missing_action", func(t *testing.T) {
		a, err := ParseAction("/x", nil)

		if err != nil {
			t.Fatal(err)
		}

		if a.Uses == nil {
			t.Fatal("expected placeholder uses")
		}

		if got := a.Uses.(func() error)(); got == nil {
			t.Error("expected missing-action error")
		}
	})

	// RoutingRouteTest::testRouteGetControllerClass
	t.Run("test_route_get_controller_class", func(t *testing.T) {
		route := NewRoute("GET", "/users/{user}", "App\\Http\\Controllers\\UserController@show")

		if got := route.GetControllerClass(); got != "App\\Http\\Controllers\\UserController" {
			t.Errorf("controller class = %q", got)
		}

		if got := route.GetActionMethod(); got != "show" {
			t.Errorf("action method = %q", got)
		}
	})

	// RoutingRouteTest::testRouteFlushController
	t.Run("test_route_flush_controller", func(t *testing.T) {
		route := NewRoute("GET", "/users/{user}", "UserController@show")
		route.Controller = &userController{}
		route.FlushController()

		if route.Controller != nil {
			t.Errorf("controller = %v, want nil", route.Controller)
		}
	})

	t.Run("test_parse_action_with_map", func(t *testing.T) {
		a, err := ParseAction("/x", map[string]any{
			"uses":       func() {},
			"middleware": []any{"auth"},
			"as":         "users.show",
			"prefix":     "api",
		})

		if err != nil {
			t.Fatal(err)
		}

		if a.As != "users.show" || a.Prefix != "api" {
			t.Errorf("a = %+v", a)
		}

		if len(a.Middleware) != 1 || a.Middleware[0] != "auth" {
			t.Errorf("middleware = %v", a.Middleware)
		}
	})

	// RouteActionTest::test_it_can_detect_a_serialized_closure
	t.Run("test_contains_serialized_closure_is_false", func(t *testing.T) {
		if ContainsSerializedClosure(&Action{}) {
			t.Error("ContainsSerializedClosure should always return false")
		}
	})
}
