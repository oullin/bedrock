package spark

import "net/http"

// Handlers bundles all HTTP handler instances for route registration.
type Handlers struct {
	Portal       *BillingPortalHandler
	Subscription *SubscriptionHandler
	Payment      *PaymentHandler
	Invoice      *InvoiceHandler
	Billing      *BillingHandler
	Inquiry      *CustomPlanInquiryHandler
	Webhook      *WebhookHandler
}

// RegisterRoutes registers all Spark routes on the given ServeMux.
func RegisterRoutes(mux *http.ServeMux, h *Handlers, cfg *Config) {
	// Spark API routes.
	mux.HandleFunc("POST /spark/subscription", h.Subscription.Create)
	mux.HandleFunc("PUT /spark/subscription", h.Subscription.Update)
	mux.HandleFunc("PUT /spark/subscription/cancel", h.Subscription.Cancel)
	mux.HandleFunc("PUT /spark/subscription/resume", h.Subscription.Resume)
	mux.HandleFunc("PUT /spark/subscription/payment-method", h.Payment.UpdatePaymentMethod)
	mux.HandleFunc("POST /spark/pending-checkout", h.Payment.NewPendingCheckout)
	mux.HandleFunc("GET /spark/{type}/{id}/invoices/{transaction}/download", h.Invoice.Download)

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
