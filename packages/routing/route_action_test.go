package routing

import "testing"

// Translation of laravel/framework tests/Routing/RouteActionTest.php.
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
}
