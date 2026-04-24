package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/handler"
)

func TestRouteSetUsesCanonicalRouteGenRoutes(t *testing.T) {
	t.Parallel()

	routes := billing.NewRouteRegistry()

	cases := map[string]string{
		billing.RoutePortal:                    "/billing",
		billing.RouteState:                     "/billing/state",
		billing.RouteRouteGen:                 "/billing/routegen",
		billing.RouteSubscriptionStore:         "/billing/subscription",
		billing.RouteSubscriptionUpdate:        "/billing/subscription",
		billing.RouteSubscriptionCancel:        "/billing/subscription/cancel",
		billing.RouteSubscriptionResume:        "/billing/subscription/resume",
		billing.RouteSubscriptionPaymentMethod: "/billing/subscription/payment-method",
		billing.RoutePendingCheckout:           "/billing/pending-checkout",
		billing.RouteInvoiceDownload:           "/billing/{type}/{id}/invoices/{transaction}/download",
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

func TestRouteSetServesRouteGenManifest(t *testing.T) {
	t.Parallel()

	routeSet := handler.NewRouteSet(&handler.Handlers{
		Portal: &handler.PortalHandler{},
	})

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/billing/routegen", nil)
	routeSet.Handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	if !strings.Contains(rec.Body.String(), billing.RouteSubscriptionPaymentMethod) {
		t.Fatalf("manifest does not contain payment method route: %s", rec.Body.String())
	}
}
