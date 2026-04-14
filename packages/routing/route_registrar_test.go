package routing

import "testing"

// Translation of upstream/framework tests/Routing/RouteRegistrarTest.php and
// the resource-registration portions of RoutingRouteTest.

func TestRouteRegistrar_Fluent(t *testing.T) {
	t.Run("test_middleware_then_get", func(t *testing.T) {
		router := NewRouter(nil, nil)
		registrar := NewRouteRegistrar(router)
		route := registrar.Middleware("auth").Prefix("api").Get("/users", func() {})

		if route.Uri != "api/users" {
			t.Errorf("uri = %q", route.Uri)
		}

		mw, _ := route.ActionMap["middleware"].([]any)

		if len(mw) != 1 || mw[0] != "auth" {
			t.Errorf("middleware = %v", mw)
		}
	})

	t.Run("test_as_prefixes_name", func(t *testing.T) {
		router := NewRouter(nil, nil)
		registrar := NewRouteRegistrar(router)
		route := registrar.As("users.").Get("/users", func() {}).Name("index")

		if route.GetName() != "users.index" {
			t.Errorf("name = %q", route.GetName())
		}
	})

	t.Run("test_domain_attribute", func(t *testing.T) {
		router := NewRouter(nil, nil)
		registrar := NewRouteRegistrar(router)
		route := registrar.Domain("api.example.com").Get("/", func() {})

		if route.GetDomain() != "api.example.com" {
			t.Errorf("domain = %q", route.GetDomain())
		}
	})
}

func TestResourceRegistrar(t *testing.T) {
	t.Run("test_register_emits_seven_routes", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Resource("users", "UserController", nil).Register()
		count := router.GetRoutes().(*RouteCollection).Count()

		if count != 7 {
			t.Errorf("count = %d, want 7", count)
		}
	})

	t.Run("test_only_filters_actions", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Resource("users", "UserController", nil).Only("index", "show").Register()
		count := router.GetRoutes().(*RouteCollection).Count()

		if count != 2 {
			t.Errorf("count = %d, want 2", count)
		}
	})

	t.Run("test_except_filters_actions", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Resource("users", "UserController", nil).Except("destroy").Register()
		count := router.GetRoutes().(*RouteCollection).Count()

		if count != 6 {
			t.Errorf("count = %d, want 6", count)
		}
	})

	t.Run("test_api_resource_excludes_create_edit", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.ApiResource("users", "UserController", nil).Register()
		count := router.GetRoutes().(*RouteCollection).Count()

		if count != 5 {
			t.Errorf("count = %d, want 5", count)
		}
	})

	t.Run("test_show_uri_uses_singular_param", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Resource("users", "UserController", nil).Only("show").Register()
		routes := router.GetRoutes().GetRoutes()

		if routes[0].Uri != "users/{user}" {
			t.Errorf("uri = %q", routes[0].Uri)
		}
	})

	t.Run("test_route_names_default", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Resource("users", "UserController", nil).Only("index").Register()
		route := router.GetRoutes().GetRoutes()[0]

		if route.GetName() != "users.index" {
			t.Errorf("name = %q", route.GetName())
		}
	})

	t.Run("test_singleton_emits_three_routes", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Singleton("profile", "ProfileController", nil).Register()
		count := router.GetRoutes().(*RouteCollection).Count()

		if count != 3 {
			t.Errorf("count = %d, want 3", count)
		}
	})

	t.Run("test_singleton_creatable_adds_create_store_destroy", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Singleton("profile", "ProfileController", nil).Creatable().Register()
		count := router.GetRoutes().(*RouteCollection).Count()

		if count != 6 {
			t.Errorf("count = %d, want 6", count)
		}
	})

	t.Run("test_nested_resource_uri", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Resource("users.posts", "PostController", nil).Only("show").Register()
		route := router.GetRoutes().GetRoutes()[0]

		if route.Uri != "users/{user}/posts/{post}" {
			t.Errorf("uri = %q", route.Uri)
		}
	})
}
