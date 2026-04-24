package spark

import (
	"strconv"

	"github.com/bedrock/packages/wayfinder"
)

const (
	RoutePortal                    = "spark.portal"
	RoutePortalForType             = "spark.portal.type"
	RoutePortalForBillable         = "spark.portal.billable"
	RouteState                     = "spark.state"
	RouteWayfinder                 = "spark.wayfinder"
	RouteSubscriptionStore         = "spark.subscription.store"
	RouteSubscriptionUpdate        = "spark.subscription.update"
	RouteSubscriptionCancel        = "spark.subscription.cancel"
	RouteSubscriptionResume        = "spark.subscription.resume"
	RouteSubscriptionPaymentMethod = "spark.subscription.payment-method"
	RoutePendingCheckout           = "spark.pending-checkout"
	RouteInvoiceDownload           = "spark.invoices.download"
)

// NewRouteRegistry returns the canonical Spark route registry shared by the
// backend router, generated frontend helpers, and portal shell.
func NewRouteRegistry() *wayfinder.Registry {
	routes := wayfinder.New()
	RegisterWayfinderRoutes(routes)

	return routes
}

// RegisterWayfinderRoutes registers canonical Spark routes on registry.
func RegisterWayfinderRoutes(routes *wayfinder.Registry) {
	routes.Add(RoutePortal, "GET", "/billing")

	routes.Add(RoutePortalForType, "GET", "/billing/{type}")

	routes.Add(RoutePortalForBillable, "GET", "/billing/{type}/{id}")

	routes.Add(RouteState, "GET", "/spark/state")

	routes.Add(RouteWayfinder, "GET", "/spark/wayfinder")

	routes.Add(RouteSubscriptionStore, "POST", "/spark/subscription")

	routes.Add(RouteSubscriptionUpdate, "PUT", "/spark/subscription")

	routes.Add(RouteSubscriptionCancel, "PUT", "/spark/subscription/cancel")

	routes.Add(RouteSubscriptionResume, "PUT", "/spark/subscription/resume")

	routes.Add(RouteSubscriptionPaymentMethod, "PUT", "/spark/subscription/payment-method")

	routes.Add(RoutePendingCheckout, "POST", "/spark/pending-checkout")

	routes.Add(RouteInvoiceDownload, "GET", "/spark/{type}/{id}/invoices/{transaction}/download")
}

// InvoiceDownloadURL resolves the canonical invoice download URL.
func InvoiceDownloadURL(routes *wayfinder.Registry, billable Billable, transaction Transaction) string {
	id := transaction.PaddleID

	if id == "" {
		id = transaction.InvoiceNumber
	}

	return routes.URL(RouteInvoiceDownload, map[string]string{
		"type":        billable.BillableType(),
		"id":          strconv.FormatInt(billable.BillableID(), 10),
		"transaction": id,
	})
}
