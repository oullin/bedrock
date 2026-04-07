package handler

import (
	"net/http"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/webhook"
)

// Handlers bundles all HTTP handler instances for route registration.
type Handlers struct {
	Portal       *PortalHandler
	Subscription *SubscriptionHandler
	Payment      *PaymentHandler
	Invoice      *InvoiceHandler
	Billing      *BillingHandler
	Inquiry      *InquiryHandler
	Webhook      *webhook.Handler
}

// RegisterRoutes registers all Billing routes on the given ServeMux.
func RegisterRoutes(mux *http.ServeMux, h *Handlers, cfg *billing.Config) {
	// Billing API routes.
	mux.HandleFunc("POST /billing/subscription", h.Subscription.Create)
	mux.HandleFunc("PUT /billing/subscription", h.Subscription.Update)
	mux.HandleFunc("PUT /billing/subscription/cancel", h.Subscription.Cancel)
	mux.HandleFunc("PUT /billing/subscription/resume", h.Subscription.Resume)
	mux.HandleFunc("PUT /billing/subscription/payment-method", h.Payment.UpdatePaymentMethod)
	mux.HandleFunc("POST /billing/pending-checkout", h.Payment.NewPendingCheckout)
	mux.HandleFunc("GET /billing/{type}/{id}/invoices/{transaction}/download", h.Invoice.Download)

	// Portal route.
	portalPattern := "GET /" + cfg.Path + "/{type}/{id}"
	mux.HandleFunc(portalPattern, h.Portal.Show)
	mux.HandleFunc("GET /"+cfg.Path, h.Portal.Show)

	// Webhook route.
	mux.HandleFunc("POST /"+cfg.WebhookPath, h.Webhook.Handle)

	// Madora billing routes.
	mux.HandleFunc("GET /subscription", h.Billing.Show)
	mux.HandleFunc("POST /billing/checkout", h.Billing.StartCheckout)
	mux.HandleFunc("POST /billing/inquiry", h.Inquiry.Store)
}
