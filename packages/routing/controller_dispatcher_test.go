package routing

import (
	"testing"

	"github.com/bedrock/packages/routing/controllers"
)

// Translation of the dispatcher portions of laravel/framework
// tests/Routing/RoutingControllerAttributeTest.php and the controller
// dispatch parts of RoutingRouteTest.
// RoutingControllerAttributeTest::testControllerMiddlewareAttributesAreInherited
// RoutingControllerAttributeTest::testControllerMiddlewareAttributesAreInheritedInDeclarationOrder
// RoutingRouteTest::testControllerCallActionMethodParameters

// userController is a fake controller used in dispatch tests.
type userController struct {
	Controller
	lastID  int
	lastTag string
}

// authController declares middleware via the HasMiddleware interface.
type authController struct{ Controller }

func (c *userController) Show(id int) string {
	c.lastID = id

	return "show"
}

func (c *userController) ShowTagged(tag string) string {
	c.lastTag = tag

	return "tagged:" + tag
}

func (c *authController) Middleware() []controllers.Middleware {
	return []controllers.Middleware{
		controllers.NewMiddleware("auth").WithExcept("Public"),
		controllers.NewMiddleware("verified").WithOnly("Settings"),
	}
}
func (c *authController) Settings() string { return "ok" }
func (c *authController) Public() string   { return "ok" }

func TestCallableDispatcher(t *testing.T) {
	t.Run("test_dispatch_func_with_string_param", func(t *testing.T) {
		d := NewCallableDispatcher(nil)
		r := NewRoute("GET", "/users/{user}", func(user string) string { return "u=" + user })
		_, _ = r.Bind(fakeRequest{path: "/users/alice"})
		got, err := d.Dispatch(r, r.ActionMap["uses"])

		if err != nil {
			t.Fatal(err)
		}

		if got != "u=alice" {
			t.Errorf("got %v", got)
		}
	})

	t.Run("test_dispatch_func_with_int_param", func(t *testing.T) {
		d := NewCallableDispatcher(nil)
		r := NewRoute("GET", "/users/{id}", func(id int) int { return id * 2 })
		_, _ = r.Bind(fakeRequest{path: "/users/21"})
		got, err := d.Dispatch(r, r.ActionMap["uses"])

		if err != nil {
			t.Fatal(err)
		}

		if got != 42 {
			t.Errorf("got %v, want 42", got)
		}
	})
}

func TestControllerDispatcher(t *testing.T) {
	// RoutingRouteTest::testControllerCallActionMethodParameters
	t.Run("test_dispatch_calls_method", func(t *testing.T) {
		d := NewControllerDispatcher(nil)
		ctrl := &userController{}
		r := NewRoute("GET", "/users/{id}", "userController@Show")
		_, _ = r.Bind(fakeRequest{path: "/users/7"})
		got, err := d.Dispatch(r, ctrl, "Show")

		if err != nil {
			t.Fatal(err)
		}

		if got != "show" || ctrl.lastID != 7 {
			t.Errorf("got %v, lastID=%d", got, ctrl.lastID)
		}
	})

	t.Run("test_dispatch_string_param", func(t *testing.T) {
		d := NewControllerDispatcher(nil)
		ctrl := &userController{}
		r := NewRoute("GET", "/tags/{tag}", "userController@ShowTagged")
		_, _ = r.Bind(fakeRequest{path: "/tags/golang"})
		got, err := d.Dispatch(r, ctrl, "ShowTagged")

		if err != nil {
			t.Fatal(err)
		}

		if got != "tagged:golang" {
			t.Errorf("got %v", got)
		}
	})

	t.Run("test_dispatch_missing_method", func(t *testing.T) {
		d := NewControllerDispatcher(nil)
		ctrl := &userController{}
		r := NewRoute("GET", "/x", "userController@Missing")
		_, err := d.Dispatch(r, ctrl, "Missing")

		var mc *MissingControllerMethodError

		if err == nil || !asMissing(err, &mc) {
			t.Errorf("expected MissingControllerMethodError, got %v", err)
		}
	})

	// RoutingControllerAttributeTest::testControllerMiddlewareAttributesAreInherited
	t.Run("test_get_middleware_filters_by_only_except", func(t *testing.T) {
		d := NewControllerDispatcher(nil)
		ctrl := &authController{}

		// Settings → "auth" applies (Public excluded), "verified" applies (only=Settings)
		got := d.GetMiddleware(ctrl, "Settings")

		if len(got) != 2 || got[0] != "auth" || got[1] != "verified" {
			t.Errorf("settings middleware = %v, want [auth verified]", got)
		}

		// public → "auth" excluded, "verified" not in only-list
		got = d.GetMiddleware(ctrl, "Public")

		if len(got) != 0 {
			t.Errorf("public middleware should be empty, got %v", got)
		}
	})

	// RoutingControllerAttributeTest::testControllerMiddlewareAttributesAreInherited
	t.Run("test_controller_middleware_attributes_are_inherited", func(t *testing.T) {
		type inheritedController struct{ Controller }

		d := NewControllerDispatcher(nil)
		ctrl := &inheritedController{}
		ctrl.Use("auth").Only("Show")
		ctrl.Use("throttle").Except("Public")

		got := d.GetMiddleware(ctrl, "Show")

		if len(got) != 2 || got[0] != "auth" || got[1] != "throttle" {
			t.Errorf("show middleware = %v", got)
		}

		got = d.GetMiddleware(ctrl, "Public")

		if len(got) != 0 {
			t.Errorf("public middleware = %v, want []", got)
		}
	})

	// RoutingControllerAttributeTest::testControllerMiddlewareAttributesAreInheritedInDeclarationOrder
	t.Run("test_controller_middleware_attributes_are_in_declaration_order", func(t *testing.T) {
		type orderedController struct{ Controller }

		d := NewControllerDispatcher(nil)
		ctrl := &orderedController{}
		ctrl.Use("first").Only("Show")
		ctrl.Use("second").Only("Show")

		got := d.GetMiddleware(ctrl, "Show")

		if len(got) != 2 || got[0] != "first" || got[1] != "second" {
			t.Errorf("show middleware order = %v", got)
		}
	})
}

// Helper for type-asserting the missing-method error without importing errors.
func asMissing(err error, dst **MissingControllerMethodError) bool {
	if e, ok := err.(*MissingControllerMethodError); ok {
		*dst = e

		return true
	}

	return false
}
