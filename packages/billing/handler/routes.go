package handler

import "net/http"

// Handlers bundles all HTTP handler instances for route registration.
type Handlers struct {
	Gateway       *GatewayHandler
	Portal        *PortalHandler
	NewSub        *NewSubscriptionHandler
	UpdateSub     *UpdateSubscriptionHandler
	CancelSub     *CancelSubscriptionHandler
	ResumeSub     *ResumeSubscriptionHandler
	Pending       *PendingCheckoutHandler
	Checkout      *CheckoutHandler
	PaddleBilling *PaddleBillingHandler
	Order         *OrderHandler
	Invoice       *DownloadInvoiceHandler
	PaymentMethod *PaymentMethodsHandler
}

// RegisterRoutes registers all Billing routes on the given ServeMux.
func RegisterRoutes(mux *http.ServeMux, h *Handlers) {
	// Billing subscription API.
	mux.HandleFunc("POST /billing/subscription", h.NewSub.Create)
	mux.HandleFunc("PUT /billing/subscription", h.UpdateSub.Update)
	mux.HandleFunc("PUT /billing/subscription/cancel", h.CancelSub.Cancel)
	mux.HandleFunc("PUT /billing/subscription/resume", h.ResumeSub.Resume)
	mux.HandleFunc("PUT /billing/subscription/payment-method", h.PaymentMethod.Setup)
	mux.HandleFunc("POST /billing/pending-checkout", h.Pending.Create)

	// Invoices.
	mux.HandleFunc("GET /billing/{type}/{id}/invoices/{transaction}/download", h.Invoice.Download)

	// Billing portal.
	mux.HandleFunc("GET /billing", h.Portal.Show)
	mux.HandleFunc("GET /billing/{type}", h.Portal.Show)
	mux.HandleFunc("GET /billing/{type}/{id}", h.Portal.Show)
	mux.HandleFunc("GET /billing/state", h.Portal.State)
}
