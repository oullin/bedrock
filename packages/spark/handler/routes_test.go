package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/routing"
	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/spark/handler"
	"github.com/bedrock/packages/wayfinder"
)

func TestRouteSetUsesCanonicalWayfinderRoutes(t *testing.T) {
	t.Parallel()

	routes := spark.NewRouteRegistry()

	cases := map[string]string{
		spark.RoutePortal:                    "/billing",
		spark.RouteState:                     "/spark/state",
		spark.RouteWayfinder:                 "/spark/wayfinder",
		spark.RouteSubscriptionStore:         "/spark/subscription",
		spark.RouteSubscriptionUpdate:        "/spark/subscription",
		spark.RouteSubscriptionCancel:        "/spark/subscription/cancel",
		spark.RouteSubscriptionResume:        "/spark/subscription/resume",
		spark.RouteSubscriptionPaymentMethod: "/spark/subscription/payment-method",
		spark.RoutePendingCheckout:           "/spark/pending-checkout",
		spark.RouteInvoiceDownload:           "/spark/{type}/{id}/invoices/{transaction}/download",
	}

	for name, want := range cases {
		route, ok := routes.Lookup(name)

		if !ok {
			t.Fatalf("route %q not registered", name)
		}

		if route.Pattern != want {
			t.Fatalf("%s pattern = %q, want %q", name, route.Pattern, want)
		}
	}
}

func TestRouteSetServesWayfinderManifest(t *testing.T) {
	t.Parallel()

	routeSet := handler.NewRouteSet(&handler.Handlers{
		Portal: &handler.PortalHandler{},
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/spark/wayfinder", nil)
	routeSet.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if !strings.Contains(rec.Body.String(), spark.RouteSubscriptionPaymentMethod) {
		t.Fatalf("manifest does not contain payment method route: %s", rec.Body.String())
	}
}

func TestNewRouteSet_ReturnsRouterRegistryAndHandler(t *testing.T) {
	t.Parallel()

	set := handler.NewRouteSet(&handler.Handlers{
		Portal: &handler.PortalHandler{},
	})

	if set.Router == nil {
		t.Fatal("Router is nil")
	}

	if set.Registry == nil {
		t.Fatal("Registry is nil")
	}

	if set.Handler == nil {
		t.Fatal("Handler is nil")
	}

	if _, ok := set.Registry.Lookup(spark.RoutePortal); !ok {
		t.Fatal("registry missing RoutePortal")
	}
}

func TestRegisterRoutes_NilRegistryFallsBackToDefault(t *testing.T) {
	t.Parallel()

	router := routing.NewRouter(nil, nil)

	handler.RegisterRoutes(router, nil, &handler.Handlers{
		Portal: &handler.PortalHandler{},
	})
}

func TestRegisterRoutes_EmptyRegistryUsesBlankPatterns(t *testing.T) {
	t.Parallel()

	router := routing.NewRouter(nil, nil)
	empty := wayfinder.New()

	handler.RegisterRoutes(router, empty, &handler.Handlers{
		Portal: &handler.PortalHandler{},
	})
}
