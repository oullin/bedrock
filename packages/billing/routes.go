package billing

import (
	"strconv"

	"github.com/bedrock/packages/routegen"
)

const (
	RoutePortal                    = "billing.portal"
	RoutePortalForType             = "billing.portal.type"
	RoutePortalForBillable         = "billing.portal.billable"
	RouteState                     = "billing.state"
	RouteRouteGen                  = "billing.routegen"
	RouteSubscriptionStore         = "billing.subscription.store"
	RouteSubscriptionUpdate        = "billing.subscription.update"
	RouteSubscriptionCancel        = "billing.subscription.cancel"
	RouteSubscriptionResume        = "billing.subscription.resume"
	RouteSubscriptionPaymentMethod = "billing.subscription.payment-method"
	RoutePendingCheckout           = "billing.pending-checkout"
	RouteInvoiceDownload           = "billing.invoices.download"
)

// NewRouteRegistry returns the canonical Billing route registry shared by the
// backend router, generated frontend helpers, and portal shell.
func NewRouteRegistry() *routegen.Registry {
	routes := routegen.New()
	RegisterRouteGenRoutes(routes)

	return routes
}

// RegisterRouteGenRoutes registers canonical Billing routes on registry.
func RegisterRouteGenRoutes(routes *routegen.Registry) {
	routes.Add(RoutePortal, "GET", "/billing")

	routes.Add(RoutePortalForType, "GET", "/billing/{type}")

	routes.Add(RoutePortalForBillable, "GET", "/billing/{type}/{id}")

	routes.Add(RouteState, "GET", "/billing/state")

	routes.Add(RouteRouteGen, "GET", "/billing/routegen")

	routes.Add(RouteSubscriptionStore, "POST", "/billing/subscription")

	routes.Add(RouteSubscriptionUpdate, "PUT", "/billing/subscription")

	routes.Add(RouteSubscriptionCancel, "PUT", "/billing/subscription/cancel")

	routes.Add(RouteSubscriptionResume, "PUT", "/billing/subscription/resume")

	routes.Add(RouteSubscriptionPaymentMethod, "PUT", "/billing/subscription/payment-method")

	routes.Add(RoutePendingCheckout, "POST", "/billing/pending-checkout")

	routes.Add(RouteInvoiceDownload, "GET", "/billing/{type}/{id}/invoices/{transaction}/download")
}

// InvoiceDownloadURL resolves the canonical invoice download URL.
func InvoiceDownloadURL(routes *routegen.Registry, billable Billable, transaction Transaction) string {
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
