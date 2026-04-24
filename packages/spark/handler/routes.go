package handler

import (
	"net/http"

	"github.com/bedrock/packages/httpx/routingx"
	"github.com/bedrock/packages/routing"
	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/wayfinder"
)

// Handlers bundles all canonical Spark HTTP handler instances.
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

// RouteSet exposes Spark's router, route registry, and dispatch handler.
type RouteSet struct {
	Router   *routing.Router
	Registry *wayfinder.Registry
	Handler  http.Handler
}

// NewRouteSet builds the canonical Spark route set.
func NewRouteSet(h *Handlers) *RouteSet {
	registry := spark.NewRouteRegistry()
	router := routing.NewRouter(nil, nil)

	RegisterRoutes(router, registry, h)

	return &RouteSet{
		Router:   router,
		Registry: registry,
		Handler:  routingx.NewHandler(router),
	}
}

// RegisterRoutes registers all canonical Spark routes.
func RegisterRoutes(router *routing.Router, registry *wayfinder.Registry, h *Handlers) {
	if registry == nil {
		registry = spark.NewRouteRegistry()
	}

	if h.Portal != nil {
		h.Portal.WithRoutes(registry)
	}

	router.Post(routePattern(registry, spark.RouteSubscriptionStore), h.NewSub.Create).Name(spark.RouteSubscriptionStore)
	router.Put(routePattern(registry, spark.RouteSubscriptionUpdate), h.UpdateSub.Update).Name(spark.RouteSubscriptionUpdate)
	router.Put(routePattern(registry, spark.RouteSubscriptionCancel), h.CancelSub.Cancel).Name(spark.RouteSubscriptionCancel)
	router.Put(routePattern(registry, spark.RouteSubscriptionResume), h.ResumeSub.Resume).Name(spark.RouteSubscriptionResume)
	router.Put(routePattern(registry, spark.RouteSubscriptionPaymentMethod), h.PaymentMethod.Setup).Name(spark.RouteSubscriptionPaymentMethod)
	router.Post(routePattern(registry, spark.RoutePendingCheckout), h.Pending.Create).Name(spark.RoutePendingCheckout)
	router.Get(routePattern(registry, spark.RouteInvoiceDownload), h.Invoice.Download).Name(spark.RouteInvoiceDownload)
	router.Get(routePattern(registry, spark.RoutePortal), h.Portal.Show).Name(spark.RoutePortal)
	router.Get(routePattern(registry, spark.RoutePortalForType), h.Portal.Show).Name(spark.RoutePortalForType)
	router.Get(routePattern(registry, spark.RoutePortalForBillable), h.Portal.Show).Name(spark.RoutePortalForBillable)
	router.Get(routePattern(registry, spark.RouteState), h.Portal.State).Name(spark.RouteState)
	router.Get(routePattern(registry, spark.RouteWayfinder), wayfinder.Handler(registry).ServeHTTP).Name(spark.RouteWayfinder)
}

func routePattern(registry *wayfinder.Registry, name string) string {
	route, ok := registry.Lookup(name)

	if !ok {
		return ""
	}

	return route.Pattern
}
