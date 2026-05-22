package handler

import (
	"net/http"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/httpx/routingx"
	"github.com/bedrock/packages/routing"
	"github.com/bedrock/packages/routegen"
)

// Handlers bundles all canonical Billing HTTP handler instances.
type Handlers struct {
	Portal        *PortalHandler
	NewSub        *NewSubscriptionHandler
	UpdateSub     *UpdateSubscriptionHandler
	CancelSub     *CancelSubscriptionHandler
	ResumeSub     *ResumeSubscriptionHandler
	Pending       *PendingCheckoutHandler
	Invoice       *DownloadInvoiceHandler
	PaymentMethod *PaymentMethodsHandler
}

// RouteSet exposes Billing's router, route registry, and dispatch handler.
type RouteSet struct {
	Router   *routing.Router
	Registry *routegen.Registry
	Handler  http.Handler
}

// NewRouteSet builds the canonical Billing route set.
func NewRouteSet(h *Handlers) *RouteSet {
	registry := billing.NewRouteRegistry()
	router := routing.NewRouter(nil, nil)

	RegisterRoutes(router, registry, h)

	return &RouteSet{
		Router:   router,
		Registry: registry,
		Handler:  routingx.NewHandler(router),
	}
}

// RegisterRoutes registers all canonical Billing routes.
func RegisterRoutes(router *routing.Router, registry *routegen.Registry, h *Handlers) {
	if registry == nil {
		registry = billing.NewRouteRegistry()
	}

	if h.Portal != nil {
		h.Portal.WithRoutes(registry)
	}

	router.Post(routePattern(registry, billing.RouteSubscriptionStore), h.NewSub.Create).Name(billing.RouteSubscriptionStore)
	router.Put(routePattern(registry, billing.RouteSubscriptionUpdate), h.UpdateSub.Update).Name(billing.RouteSubscriptionUpdate)
	router.Put(routePattern(registry, billing.RouteSubscriptionCancel), h.CancelSub.Cancel).Name(billing.RouteSubscriptionCancel)
	router.Put(routePattern(registry, billing.RouteSubscriptionResume), h.ResumeSub.Resume).Name(billing.RouteSubscriptionResume)
	router.Put(routePattern(registry, billing.RouteSubscriptionPaymentMethod), h.PaymentMethod.Setup).Name(billing.RouteSubscriptionPaymentMethod)
	router.Post(routePattern(registry, billing.RoutePendingCheckout), h.Pending.Create).Name(billing.RoutePendingCheckout)
	router.Get(routePattern(registry, billing.RouteInvoiceDownload), h.Invoice.Download).Name(billing.RouteInvoiceDownload)
	router.Get(routePattern(registry, billing.RoutePortal), h.Portal.Show).Name(billing.RoutePortal)
	router.Get(routePattern(registry, billing.RouteState), h.Portal.State).Name(billing.RouteState)
	router.Get(routePattern(registry, billing.RouteRouteGen), routegen.Handler(registry).ServeHTTP).Name(billing.RouteRouteGen)
	router.Get(routePattern(registry, billing.RoutePortalForType), h.Portal.Show).Name(billing.RoutePortalForType)
	router.Get(routePattern(registry, billing.RoutePortalForBillable), h.Portal.Show).Name(billing.RoutePortalForBillable)
}

func routePattern(registry *routegen.Registry, name string) string {
	route, ok := registry.Lookup(name)

	if !ok {
		return ""
	}

	return route.Pattern
}
