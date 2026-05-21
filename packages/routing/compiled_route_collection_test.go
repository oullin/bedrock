package routing

import "testing"

// CompiledRouteCollectionTest::testCompiledRouteCollectionCanRetrieveByActionWithLeadingBackslash

func TestCompiledRouteCollection_ActionLookups(t *testing.T) {
	// CompiledRouteCollectionTest::testCompiledRouteCollectionCanRetrieveByActionWithLeadingBackslash
	t.Run("test_get_by_action_normalizes_leading_backslash", func(t *testing.T) {
		route := NewRoute("GET", "/users", map[string]any{
			"controller": "\\App\\Http\\Controllers\\UserController@index",
		})
		c := NewCompiledRouteCollection([]*Route{route}, nil)

		if c.GetByAction("App\\Http\\Controllers\\UserController@index") != route {
			t.Fatal("GetByAction should normalize leading backslashes on init")
		}

		c.RefreshActionLookups()

		if c.GetByAction("App\\Http\\Controllers\\UserController@index") != route {
			t.Fatal("GetByAction should normalize leading backslashes after refresh")
		}
	})
}
