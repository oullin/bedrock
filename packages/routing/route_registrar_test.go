package routing

import "testing"

// Ref: @bedrock/code-0394
// the resource-registration portions of RoutingRouteTest.
// RouteRegistrarTest::testMiddlewareFluentRegistration
// RouteRegistrarTest::testMiddlewareAsNull
// RouteRegistrarTest::testCanRegisterGroupWithPrefix
// RouteRegistrarTest::testCanRegisterGroupWithDomain
// RouteRegistrarTest::testFallbackRoute
// RouteRegistrarTest::testCanRegisterMatchRouteWithClosureAction
// RouteRegistrarTest::testCanRegisterAnyRouteWithClosureAction
// RouteRegistrarTest::testWithoutMiddlewareRegistration
// RouteRegistrarTest::testCanSetScopeBindings
// RouteRegistrarTest::testCanSetWithoutScopedBindings
// RouteRegistrarTest::testWhereNumberRegistration
// RouteRegistrarTest::testWhereAlphaRegistration
// RouteRegistrarTest::testWhereAlphaNumericRegistration
// RouteRegistrarTest::testWhereInRegistration
// RouteRegistrarTest::testCanRegisterRouteWithControllerActionArray
// RouteRegistrarTest::testCanRegisterNamespacedGroupRouteWithControllerActionArray
// RouteRegistrarTest::testCanRegisterRouteWithArrayAndControllerAction
// RouteRegistrarTest::testCanRegisterResource
// RouteRegistrarTest::testCanRegisterResourcesWithOnlyOption
// RouteRegistrarTest::testCanRegisterResourcesWithExceptOption
// RouteRegistrarTest::testCanLimitMethodsOnRegisteredResource
// RouteRegistrarTest::testCanExcludeMethodsOnRegisteredResource
// RouteRegistrarTest::testCanLimitAndExcludeMethodsOnRegisteredResource
// RouteRegistrarTest::testCanExcludeMethodsOnRegisteredApiResource
// RouteRegistrarTest::testCanRegisterApiResourcesWithExceptOption
// RouteRegistrarTest::testCanRegisterApiResourcesWithOnlyOption
// RouteRegistrarTest::testCanRegisterApiResourcesWithoutOption
// RouteRegistrarTest::testCanNameRoutesOnRegisteredResource
// RouteRegistrarTest::testCanSetMiddlewareOnRegisteredResource
// RouteRegistrarTest::testCanRegisterCreatableSingleton
// RouteRegistrarTest::testSingletonCreatableNotDestroyable
// RouteRegistrarTest::testCanRegisterSingleton
// RouteRegistrarTest::testCanRegisterApiSingleton
// RouteRegistrarTest::testCanRegisterCreatableApiSingleton
// RouteRegistrarTest::testApiSingletonCreatableNotDestroyable
// RouteRegistrarTest::testApiSingletonCanBeOnlyCreatable
// RouteRegistrarTest::testApiSingletonCanIncludeAnySingletonMethods
// RouteRegistrarTest::testCanSetRouteName
// RouteRegistrarTest::testCanSetRouteNameUsingNameAlias
// RouteRegistrarTest::testCanOverrideParametersOnRegisteredResource
// RouteRegistrarTest::testCanAccessRegisteredResourceRoutesAsRouteCollection
// RouteRegistrarTest::testCanRegisterResourcesWithoutOption
// RouteRegistrarTest::testSingletonCanBeDestroyable
// RouteRegistrarTest::testSingletonCanBeOnlyCreatable
// RouteRegistrarTest::testSingletonDoesntAllowIncludingUnsupportedMethods
// RouteRegistrarTest::testApiSingletonCanBeDestroyable
// RouteRegistrarTest::testCanRegisterGroupWithNamespace
// RouteRegistrarTest::testCanRegisterGroupWithDomainAndNamePrefix
// RouteRegistrarTest::testCanRegisterGroupWithController
// RouteRegistrarTest::testCanOverrideGroupControllerWithStringSyntax

func TestRouteRegistrar_Fluent(t *testing.T) {
	// RouteRegistrarTest::testMiddlewareFluentRegistration
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

	// RouteRegistrarTest::testCanSetRouteNameUsingNameAlias
	t.Run("test_as_prefixes_name", func(t *testing.T) {
		router := NewRouter(nil, nil)
		registrar := NewRouteRegistrar(router)
		route := registrar.As("users.").Get("/users", func() {}).Name("index")

		if route.GetName() != "users.index" {
			t.Errorf("name = %q", route.GetName())
		}
	})

	// RouteRegistrarTest::testCanRegisterRouteWithControllerActionArray
	t.Run("test_route_action_array_sets_controller", func(t *testing.T) {
		router := NewRouter(nil, nil)
		route := router.Get("/users", map[string]any{"uses": "UserController@show"})

		if route.ActionMap["controller"] != "UserController@show" {
			t.Errorf("controller = %v", route.ActionMap["controller"])
		}
	})

	// RouteRegistrarTest::testCanRegisterGroupWithDomain
	t.Run("test_domain_attribute", func(t *testing.T) {
		router := NewRouter(nil, nil)
		registrar := NewRouteRegistrar(router)
		route := registrar.Domain("api.example.com").Get("/", func() {})

		if route.GetDomain() != "api.example.com" {
			t.Errorf("domain = %q", route.GetDomain())
		}
	})

	// RouteRegistrarTest::testCanRegisterNamespacedGroupRouteWithControllerActionArray
	t.Run("test_group_namespace_prefixes_controller_action", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Group(map[string]any{"namespace": "App\\Http\\Controllers"}, func(r *Router) {
			route := r.Get("/users", map[string]any{"uses": "UserController@show"})

			if route.ActionMap["controller"] != "App\\Http\\Controllers\\UserController@show" {
				t.Errorf("controller = %v", route.ActionMap["controller"])
			}
		})
	})

	// RouteRegistrarTest::testCanRegisterGroupWithDomainAndNamePrefix
	t.Run("test_group_domain_and_name_prefix_merge", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Group(map[string]any{
			"domain": "api.example.com",
			"as":     "api.",
		}, func(r *Router) {
			route := r.Get("/", func() {}).Name("index")

			if route.GetDomain() != "api.example.com" {
				t.Errorf("domain = %q", route.GetDomain())
			}

			if route.GetName() != "api.index" {
				t.Errorf("name = %q", route.GetName())
			}
		})
	})

	// RouteRegistrarTest::testCanRegisterGroupWithController
	t.Run("test_group_controller_prefixes_string_action", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Group(map[string]any{"controller": "App\\Http\\Controllers\\UserController"}, func(r *Router) {
			route := r.Get("/users", "show")

			if route.ActionMap["controller"] != "App\\Http\\Controllers\\UserController@show" {
				t.Errorf("controller = %v", route.ActionMap["controller"])
			}
		})
	})

	// RouteRegistrarTest::testCanOverrideGroupControllerWithStringSyntax
	t.Run("test_group_controller_does_not_override_explicit_action", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Group(map[string]any{"controller": "App\\Http\\Controllers\\UserController"}, func(r *Router) {
			route := r.Get("/users", "OtherController@show")

			if route.ActionMap["controller"] != "OtherController@show" {
				t.Errorf("controller = %v", route.ActionMap["controller"])
			}
		})
	})

	// RouteRegistrarTest::testMiddlewareAsNull
	t.Run("test_middleware_nil_is_ignored", func(t *testing.T) {
		router := NewRouter(nil, nil)
		route := NewRouteRegistrar(router).Middleware(nil).Get("/users", func() {})

		if _, ok := route.ActionMap["middleware"]; ok {
			t.Errorf("middleware = %v, want no middleware entry", route.ActionMap["middleware"])
		}
	})

	// RouteRegistrarTest::testWithoutMiddlewareRegistration
	t.Run("test_without_middleware_registration", func(t *testing.T) {
		router := NewRouter(nil, nil)
		route := NewRouteRegistrar(router).WithoutMiddleware("csrf").Get("/users", func() {})
		excluded, _ := route.ActionMap["excluded_middleware"].([]any)

		if len(excluded) != 1 || excluded[0] != "csrf" {
			t.Errorf("excluded middleware = %v", excluded)
		}
	})

	// RouteRegistrarTest::testCanSetScopeBindings
	t.Run("test_scope_bindings_are_propagated", func(t *testing.T) {
		router := NewRouter(nil, nil)
		route := NewRouteRegistrar(router).ScopeBindings().Get("/users/{user}", func() {})

		if v, ok := route.ActionMap["scope_bindings"].(bool); !ok || !v {
			t.Errorf("scope_bindings = %v", route.ActionMap["scope_bindings"])
		}
	})

	// RouteRegistrarTest::testCanSetWithoutScopedBindings
	t.Run("test_without_scoped_bindings_are_propagated", func(t *testing.T) {
		router := NewRouter(nil, nil)
		route := NewRouteRegistrar(router).WithoutScopedBindings().Get("/users/{user}", func() {})

		if v, ok := route.ActionMap["scope_bindings"].(bool); !ok || v {
			t.Errorf("scope_bindings = %v, want false", route.ActionMap["scope_bindings"])
		}
	})

	t.Run("test_where_constraints_are_propagated", func(t *testing.T) {
		// RouteRegistrarTest::testWhereNumberRegistration
		t.Run("number", func(t *testing.T) {
			router := NewRouter(nil, nil)
			registrar := NewRouteRegistrar(router)
			registrar.WhereNumber("id")
			route := registrar.Get("/users/{id}", func() {})

			if route.Wheres["id"] != "[0-9]+" {
				t.Errorf("wheres = %v", route.Wheres)
			}
		})

		// RouteRegistrarTest::testWhereAlphaRegistration
		t.Run("alpha", func(t *testing.T) {
			router := NewRouter(nil, nil)
			registrar := NewRouteRegistrar(router)
			registrar.WhereAlpha("slug")
			route := registrar.Get("/posts/{slug}", func() {})

			if route.Wheres["slug"] != "[a-zA-Z]+" {
				t.Errorf("wheres = %v", route.Wheres)
			}
		})

		// RouteRegistrarTest::testWhereAlphaNumericRegistration
		t.Run("alpha_numeric", func(t *testing.T) {
			router := NewRouter(nil, nil)
			registrar := NewRouteRegistrar(router)
			registrar.WhereAlphaNumeric("code")
			route := registrar.Get("/codes/{code}", func() {})

			if route.Wheres["code"] != "[a-zA-Z0-9]+" {
				t.Errorf("wheres = %v", route.Wheres)
			}
		})

		// RouteRegistrarTest::testWhereInRegistration
		t.Run("in", func(t *testing.T) {
			router := NewRouter(nil, nil)
			registrar := NewRouteRegistrar(router)
			registrar.WhereIn("role", []string{"admin", "user"})
			route := registrar.Get("/roles/{role}", func() {})

			if route.Wheres["role"] != "admin|user" {
				t.Errorf("wheres = %v", route.Wheres)
			}
		})
	})
}

func TestResourceRegistrar(t *testing.T) {
	// RouteRegistrarTest::testCanRegisterResource
	t.Run("test_register_emits_seven_routes", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Resource("users", "UserController", nil).Register()
		count := router.GetRoutes().(*RouteCollection).Count()

		if count != 7 {
			t.Errorf("count = %d, want 7", count)
		}
	})

	// RouteRegistrarTest::testCanRegisterResourcesWithOnlyOption
	// RouteRegistrarTest::testCanLimitMethodsOnRegisteredResource
	t.Run("test_only_filters_actions", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Resource("users", "UserController", nil).Only("index", "show", "destroy").Register()
		count := router.GetRoutes().(*RouteCollection).Count()

		if count != 3 {
			t.Errorf("count = %d, want 3", count)
		}

		assertNamedRoutes(t, router, []string{"users.index", "users.show", "users.destroy"}, nil)
	})

	// RouteRegistrarTest::testCanRegisterResourcesWithExceptOption
	// RouteRegistrarTest::testCanExcludeMethodsOnRegisteredResource
	t.Run("test_except_filters_actions", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Resource("users", "UserController", nil).Except("index", "create", "store", "show", "edit").Register()
		count := router.GetRoutes().(*RouteCollection).Count()

		if count != 2 {
			t.Errorf("count = %d, want 2", count)
		}

		assertNamedRoutes(t, router, []string{"users.update", "users.destroy"}, []string{"users.index", "users.show"})
	})

	// RouteRegistrarTest::testCanLimitAndExcludeMethodsOnRegisteredResource
	t.Run("test_only_and_except_filters_actions", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Resource("users", "UserController", nil).
			Only("index", "show", "destroy").
			Except("destroy").
			Register()

		if count := router.GetRoutes().(*RouteCollection).Count(); count != 2 {
			t.Errorf("count = %d, want 2", count)
		}

		assertNamedRoutes(t, router, []string{"users.index", "users.show"}, []string{"users.destroy"})
	})

	// RouteRegistrarTest::testCanRegisterApiResourcesWithoutOption
	t.Run("test_api_resource_excludes_create_edit", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.ApiResource("users", "UserController", nil).Register()
		count := router.GetRoutes().(*RouteCollection).Count()

		if count != 5 {
			t.Errorf("count = %d, want 5", count)
		}
	})

	// RouteRegistrarTest::testCanExcludeMethodsOnRegisteredApiResource
	t.Run("test_api_resource_except_keeps_api_exclusions", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.ApiResource("users", "UserController", nil).Except("index", "show", "store").Register()

		if count := router.GetRoutes().(*RouteCollection).Count(); count != 2 {
			t.Errorf("count = %d, want 2", count)
		}

		assertNamedRoutes(t, router,
			[]string{"users.update", "users.destroy"},
			[]string{"users.index", "users.store", "users.show", "users.create", "users.edit"},
		)
	})

	// RouteRegistrarTest::testCanRegisterApiResourcesWithExceptOption
	t.Run("test_api_resources_with_except_option", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.ApiResources(map[string]string{
			"resource-one": "OneController",
			"resource-two": "TwoController",
		}, map[string]any{"except": []string{"create", "show"}})

		if count := router.GetRoutes().(*RouteCollection).Count(); count != 8 {
			t.Errorf("count = %d, want 8", count)
		}

		for _, resource := range []string{"resource-one", "resource-two"} {
			assertNamedRoutes(t, router,
				[]string{resource + ".index", resource + ".store", resource + ".update", resource + ".destroy"},
				[]string{resource + ".create", resource + ".show", resource + ".edit"},
			)
		}
	})

	// RouteRegistrarTest::testCanRegisterApiResourcesWithOnlyOption
	t.Run("test_api_resources_with_only_option", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.ApiResources(map[string]string{
			"resource-one": "OneController",
			"resource-two": "TwoController",
		}, map[string]any{"only": []string{"index", "show"}})

		if count := router.GetRoutes().(*RouteCollection).Count(); count != 4 {
			t.Errorf("count = %d, want 4", count)
		}

		for _, resource := range []string{"resource-one", "resource-two"} {
			assertNamedRoutes(t, router,
				[]string{resource + ".index", resource + ".show"},
				[]string{resource + ".store", resource + ".update", resource + ".destroy", resource + ".create", resource + ".edit"},
			)
		}
	})

	// RouteRegistrarTest::testCanRegisterResourcesWithoutOption
	t.Run("test_resources_without_option_registers_each_resource", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Resources(map[string]string{
			"users": "UserController",
			"posts": "PostController",
		}, nil)

		if count := router.GetRoutes().(*RouteCollection).Count(); count != 14 {
			t.Errorf("count = %d, want 14", count)
		}

		assertNamedRoutes(t, router, []string{"users.index", "users.destroy", "posts.index", "posts.destroy"}, nil)
	})

	t.Run("test_show_uri_uses_singular_param", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Resource("users", "UserController", nil).Only("show").Register()
		routes := router.GetRoutes().GetRoutes()

		if routes[0].Uri != "users/{user}" {
			t.Errorf("uri = %q", routes[0].Uri)
		}
	})

	// RouteRegistrarTest::testCanNameRoutesOnRegisteredResource
	t.Run("test_route_names_default", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Resource("users", "UserController", nil).Only("index").Register()
		route := router.GetRoutes().GetRoutes()[0]

		if route.GetName() != "users.index" {
			t.Errorf("name = %q", route.GetName())
		}
	})

	// RouteRegistrarTest::testCanAccessRegisteredResourceRoutesAsRouteCollection
	t.Run("test_registered_resource_routes_are_returned_as_collection", func(t *testing.T) {
		router := NewRouter(nil, nil)
		resource := router.Resource("users", "UserController", nil).Register()

		if resource.Count() != 7 {
			t.Errorf("resource count = %d, want 7", resource.Count())
		}

		assertNamedRoutes(t, router,
			[]string{"users.index", "users.create", "users.store", "users.show", "users.edit", "users.update", "users.destroy"},
			nil,
		)
	})

	// RouteRegistrarTest::testCanRegisterSingleton
	t.Run("test_singleton_emits_three_routes", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Singleton("profile", "ProfileController", nil).Register()
		count := router.GetRoutes().(*RouteCollection).Count()

		if count != 3 {
			t.Errorf("count = %d, want 3", count)
		}
	})

	// RouteRegistrarTest::testCanRegisterCreatableSingleton
	t.Run("test_singleton_creatable_adds_create_store_destroy", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Singleton("profile", "ProfileController", nil).Creatable().Register()
		count := router.GetRoutes().(*RouteCollection).Count()

		if count != 6 {
			t.Errorf("count = %d, want 6", count)
		}
	})

	// RouteRegistrarTest::testSingletonCreatableNotDestroyable
	t.Run("test_singleton_creatable_can_exclude_destroy", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Singleton("profile", "ProfileController", nil).Creatable().Except("destroy").Register()

		if count := router.GetRoutes().(*RouteCollection).Count(); count != 5 {
			t.Errorf("count = %d, want 5", count)
		}

		assertNamedRoutes(t, router,
			[]string{"profile.create", "profile.store", "profile.show", "profile.edit", "profile.update"},
			[]string{"profile.destroy"},
		)
	})

	// RouteRegistrarTest::testCanRegisterApiSingleton
	t.Run("test_api_singleton_excludes_edit", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.ApiSingleton("profile", "ProfileController", nil).Register()

		if count := router.GetRoutes().(*RouteCollection).Count(); count != 2 {
			t.Errorf("count = %d, want 2", count)
		}

		assertNamedRoutes(t, router, []string{"profile.show", "profile.update"}, []string{"profile.edit"})
	})

	// RouteRegistrarTest::testCanRegisterCreatableApiSingleton
	t.Run("test_creatable_api_singleton_excludes_create_edit", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.ApiSingleton("profile", "ProfileController", nil).Creatable().Register()

		if count := router.GetRoutes().(*RouteCollection).Count(); count != 4 {
			t.Errorf("count = %d, want 4", count)
		}

		assertNamedRoutes(t, router,
			[]string{"profile.store", "profile.show", "profile.update", "profile.destroy"},
			[]string{"profile.create", "profile.edit"},
		)
	})

	// RouteRegistrarTest::testApiSingletonCreatableNotDestroyable
	t.Run("test_creatable_api_singleton_can_exclude_destroy", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.ApiSingleton("profile", "ProfileController", nil).Creatable().Except("destroy").Register()

		if count := router.GetRoutes().(*RouteCollection).Count(); count != 3 {
			t.Errorf("count = %d, want 3", count)
		}

		assertNamedRoutes(t, router,
			[]string{"profile.store", "profile.show", "profile.update"},
			[]string{"profile.create", "profile.edit", "profile.destroy"},
		)
	})

	t.Run("test_nested_resource_uri", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Resource("users.posts", "PostController", nil).Only("show").Register()
		route := router.GetRoutes().GetRoutes()[0]

		if route.Uri != "users/{user}/posts/{post}" {
			t.Errorf("uri = %q", route.Uri)
		}
	})

	// RouteRegistrarTest::testCanOverrideParametersOnRegisteredResource
	t.Run("test_resource_parameters_override_wildcards", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Resource("users.posts", "PostController", nil).
			Parameters(map[string]string{"users": "account", "posts": "article"}).
			Only("show").
			Register()
		route := router.GetRoutes().GetRoutes()[0]

		if route.Uri != "users/{account}/posts/{article}" {
			t.Errorf("uri = %q", route.Uri)
		}
	})

	// RouteRegistrarTest::testCanSetMiddlewareOnRegisteredResource
	t.Run("test_resource_middleware_is_applied", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Resource("users", "UserController", nil).Middleware("auth").Only("index").Register()
		route := router.GetRoutes().GetRoutes()[0]

		mw, _ := route.ActionMap["middleware"].([]any)

		if len(mw) != 1 || mw[0] != "auth" {
			t.Errorf("middleware = %v", mw)
		}
	})

	// RouteRegistrarTest::testSingletonCanBeDestroyable
	t.Run("test_singleton_destroyable_adds_destroy_route", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Singleton("profile", "ProfileController", nil).Destroyable().Register()

		if count := router.GetRoutes().(*RouteCollection).Count(); count != 4 {
			t.Errorf("count = %d, want 4", count)
		}
	})

	// RouteRegistrarTest::testApiSingletonCanBeDestroyable
	t.Run("test_api_singleton_destroyable_adds_destroy_route", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.ApiSingleton("profile", "ProfileController", nil).Destroyable().Register()

		if count := router.GetRoutes().(*RouteCollection).Count(); count != 3 {
			t.Errorf("count = %d, want 3", count)
		}
	})

	// RouteRegistrarTest::testSingletonCanBeOnlyCreatable
	t.Run("test_singleton_can_be_only_creatable", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Singleton("profile", "ProfileController", nil).Creatable().Only("create", "store").Register()

		if count := router.GetRoutes().(*RouteCollection).Count(); count != 2 {
			t.Errorf("count = %d, want 2", count)
		}

		assertNamedRoutes(t, router, []string{"profile.create", "profile.store"}, []string{"profile.show"})
	})

	// RouteRegistrarTest::testApiSingletonCanBeOnlyCreatable
	t.Run("test_api_singleton_can_be_only_creatable", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.ApiSingleton("profile", "ProfileController", nil).Creatable().Only("store").Register()

		if count := router.GetRoutes().(*RouteCollection).Count(); count != 1 {
			t.Errorf("count = %d, want 1", count)
		}

		assertNamedRoutes(t, router, []string{"profile.store"}, []string{"profile.show", "profile.create"})
	})

	// RouteRegistrarTest::testSingletonDoesntAllowIncludingUnsupportedMethods
	t.Run("test_singleton_rejects_unsupported_methods", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.Singleton("profile", "ProfileController", nil).Only("index", "store", "create", "destroy").Register()

		if count := router.GetRoutes().(*RouteCollection).Count(); count != 0 {
			t.Errorf("count = %d, want 0", count)
		}

		router.ApiSingleton("account", "AccountController", nil).Only("index", "store", "create", "destroy").Register()

		if count := router.GetRoutes().(*RouteCollection).Count(); count != 0 {
			t.Errorf("count after api singleton = %d, want 0", count)
		}
	})

	// RouteRegistrarTest::testApiSingletonCanIncludeAnySingletonMethods
	t.Run("test_api_singleton_can_explicitly_include_edit", func(t *testing.T) {
		router := NewRouter(nil, nil)
		router.ApiSingleton("profile", "ProfileController", nil).Only("edit").Register()

		if count := router.GetRoutes().(*RouteCollection).Count(); count != 1 {
			t.Errorf("count = %d, want 1", count)
		}

		assertNamedRoutes(t, router, []string{"profile.edit"}, []string{"profile.show", "profile.update"})
	})
}

func assertNamedRoutes(t *testing.T, router *Router, present []string, absent []string) {
	t.Helper()

	for _, name := range present {
		if !router.GetRoutes().HasNamedRoute(name) {
			t.Errorf("expected route %q to be registered", name)
		}
	}

	for _, name := range absent {
		if router.GetRoutes().HasNamedRoute(name) {
			t.Errorf("expected route %q to be absent", name)
		}
	}
}
