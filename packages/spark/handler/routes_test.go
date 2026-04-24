package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/spark/handler"
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
