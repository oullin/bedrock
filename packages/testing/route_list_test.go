package testing_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/bedrock/packages/routing"
	packagetesting "github.com/bedrock/packages/testing"
)

func TestRouteListRendering(t *testing.T) {
	t.Parallel()

	baseRoutes := []packagetesting.RouteListEntry{
		{
			Name:   "dashboard",
			Action: "DashboardController@index",
			URI:    "/dashboard",
		},
		{
			Name:   "users.show",
			Action: "UserController@show",
			URI:    routing.ParseRouteUri("/users/{user:slug}").Uri,
			BindingFields: map[string]string{
				"user": "slug",
			},
		},
		{
			Name:    "closure.route",
			Action:  "Closure",
			URI:     "/closures",
			Closure: true,
			Path:    "routes/web.go:42",
		},
		{
			Name:   "vendor.metrics",
			Action: "VendorController@index",
			URI:    "/vendor/metrics",
			Vendor: true,
		},
	}

	t.Run("Console/RouteListCommandTest::testDisplayRoutesForCli", func(t *testing.T) {
		got, err := packagetesting.RenderRouteList(baseRoutes[:2], packagetesting.RouteListOptions{})

		if err != nil {
			t.Fatalf("expected route list render to succeed, got %v", err)
		}

		if !strings.Contains(got, "dashboard") || !strings.Contains(got, "/users/{user}") {
			t.Fatalf("unexpected route list output: %q", got)
		}
	})

	t.Run("Console/RouteListCommandTest::testDisplayRoutesForCliInVerboseMode", func(t *testing.T) {
		got, err := packagetesting.RenderRouteList(baseRoutes[:2], packagetesting.RouteListOptions{Verbose: true})

		if err != nil {
			t.Fatalf("expected verbose route list render to succeed, got %v", err)
		}

		if !strings.Contains(got, "DashboardController@index") {
			t.Fatalf("expected verbose route list output to include actions, got %q", got)
		}
	})

	t.Run("Console/RouteListCommandTest::testRouteCanBeFilteredByName", func(t *testing.T) {
		got, err := packagetesting.RenderRouteList(baseRoutes, packagetesting.RouteListOptions{NameFilter: "dashboard"})

		if err != nil {
			t.Fatalf("expected route list render to succeed, got %v", err)
		}

		if strings.Contains(got, "users.show") || !strings.Contains(got, "dashboard") {
			t.Fatalf("expected name filter to narrow route list, got %q", got)
		}
	})

	t.Run("Console/RouteListCommandTest::testRouteCanBeFilteredByAction", func(t *testing.T) {
		got, err := packagetesting.RenderRouteList(baseRoutes, packagetesting.RouteListOptions{ActionFilter: "UserController"})

		if err != nil {
			t.Fatalf("expected route list render to succeed, got %v", err)
		}

		if strings.Contains(got, "dashboard") || !strings.Contains(got, "users.show") {
			t.Fatalf("expected action filter to narrow route list, got %q", got)
		}
	})

	t.Run("Console/RouteListCommandTest::testClosurePathIsDisplayedInVerboseMode", func(t *testing.T) {
		got, err := packagetesting.RenderRouteList([]packagetesting.RouteListEntry{baseRoutes[2]}, packagetesting.RouteListOptions{Verbose: true})

		if err != nil {
			t.Fatalf("expected verbose route list render to succeed, got %v", err)
		}

		if !strings.Contains(got, "routes/web.go:42") {
			t.Fatalf("expected closure path to be visible in verbose mode, got %q", got)
		}
	})

	t.Run("Console/RouteListCommandTest::testClosurePathIsDisplayedInNonVerboseMode", func(t *testing.T) {
		got, err := packagetesting.RenderRouteList([]packagetesting.RouteListEntry{baseRoutes[2]}, packagetesting.RouteListOptions{})

		if err != nil {
			t.Fatalf("expected route list render to succeed, got %v", err)
		}

		if strings.Contains(got, "routes/web.go:42") {
			t.Fatalf("expected closure path to be hidden outside verbose mode, got %q", got)
		}
	})

	t.Run("Console/RouteListCommandTest::testClosurePathIsIncludedInJsonOutput", func(t *testing.T) {
		got, err := packagetesting.RenderRouteList([]packagetesting.RouteListEntry{baseRoutes[2]}, packagetesting.RouteListOptions{JSON: true})

		if err != nil {
			t.Fatalf("expected JSON route list render to succeed, got %v", err)
		}

		if !strings.Contains(got, `"path":"routes/web.go:42"`) {
			t.Fatalf("expected closure path in JSON output, got %q", got)
		}
	})

	t.Run("Console/RouteListCommandTest::testControllerRouteHasNullPathInJsonOutput", func(t *testing.T) {
		got, err := packagetesting.RenderRouteList([]packagetesting.RouteListEntry{baseRoutes[0]}, packagetesting.RouteListOptions{JSON: true})

		if err != nil {
			t.Fatalf("expected JSON route list render to succeed, got %v", err)
		}

		if !strings.Contains(got, `"path":null`) {
			t.Fatalf("expected controller route to emit a null path, got %q", got)
		}
	})

	t.Run("Console/RouteListCommandTest::testDisplayRoutesExceptVendor", func(t *testing.T) {
		got, err := packagetesting.RenderRouteList(baseRoutes, packagetesting.RouteListOptions{ExceptVendor: true})

		if err != nil {
			t.Fatalf("expected route list render to succeed, got %v", err)
		}

		if strings.Contains(got, "vendor.metrics") {
			t.Fatalf("expected vendor route to be filtered out, got %q", got)
		}
	})

	t.Run("Console/RouteListCommandTest::testDisplayRoutesWithBindingFields", func(t *testing.T) {
		got, err := packagetesting.RenderRouteList([]packagetesting.RouteListEntry{baseRoutes[1]}, packagetesting.RouteListOptions{})

		if err != nil {
			t.Fatalf("expected route list render to succeed, got %v", err)
		}

		if !strings.Contains(got, "user:slug") {
			t.Fatalf("expected binding fields to be displayed, got %q", got)
		}
	})

	t.Run("Console/RouteListCommandTest::testDisplayRoutesWithBindingFieldsAsJson", func(t *testing.T) {
		got, err := packagetesting.RenderRouteList([]packagetesting.RouteListEntry{baseRoutes[1]}, packagetesting.RouteListOptions{JSON: true})

		if err != nil {
			t.Fatalf("expected JSON route list render to succeed, got %v", err)
		}

		var payload []map[string]any

		if err := json.Unmarshal([]byte(got), &payload); err != nil {
			t.Fatalf("expected valid JSON, got %v", err)
		}

		bindingFields, ok := payload[0]["bindingFields"].(map[string]any)

		if !ok || bindingFields["user"] != "slug" {
			t.Fatalf("expected binding fields in JSON output, got %v", payload[0])
		}
	})
}
