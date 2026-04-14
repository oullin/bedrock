package routing

import "testing"

// Translation of laravel/framework tests/Routing/RouteUriTest.php.
// Each subtest preserves the upstream PHP method name for grep-ability.
func TestRouteUri(t *testing.T) {
	t.Run("test_parsing_uri_with_no_binding_fields", func(t *testing.T) {
		ru := ParseRouteUri("/users/{user}")

		if ru.Uri != "/users/{user}" {
			t.Errorf("uri = %q", ru.Uri)
		}

		if len(ru.BindingFields) != 0 {
			t.Errorf("expected no binding fields, got %v", ru.BindingFields)
		}
	})

	t.Run("test_parsing_uri_with_binding_fields", func(t *testing.T) {
		ru := ParseRouteUri("/users/{user:slug}")

		if ru.Uri != "/users/{user}" {
			t.Errorf("uri = %q, want /users/{user}", ru.Uri)
		}

		if got := ru.BindingFields["user"]; got != "slug" {
			t.Errorf("binding = %q, want slug", got)
		}
	})

	t.Run("test_parsing_uri_with_optional_binding_fields", func(t *testing.T) {
		ru := ParseRouteUri("/users/{user:slug?}")

		if ru.Uri != "/users/{user?}" {
			t.Errorf("uri = %q, want /users/{user?}", ru.Uri)
		}

		if got := ru.BindingFields["user"]; got != "slug" {
			t.Errorf("binding = %q, want slug", got)
		}
	})

	t.Run("test_parsing_uri_with_multiple_binding_fields", func(t *testing.T) {
		ru := ParseRouteUri("/posts/{post:uuid}/comments/{comment:id}")

		if ru.Uri != "/posts/{post}/comments/{comment}" {
			t.Errorf("uri = %q", ru.Uri)
		}

		if ru.BindingFields["post"] != "uuid" || ru.BindingFields["comment"] != "id" {
			t.Errorf("bindings = %v", ru.BindingFields)
		}
	})
}
