package handler_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/routing"
	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/handler"
	"github.com/bedrock/packages/routegen"
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

	if _, ok := set.Registry.Lookup(billing.RoutePortal); !ok {
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
	empty := routegen.New()

	handler.RegisterRoutes(router, empty, &handler.Handlers{
		Portal: &handler.PortalHandler{},
	})
}
