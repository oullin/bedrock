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
	Checkout      *CheckoutHandler
	PaddleBilling *PaddleBillingHandler
	Order         *OrderHandler
	Invoice       *DownloadInvoiceHandler
	PaymentMethod *PaymentMethodsHandler
}

// RegisterRoutes registers all Billing routes on the given ServeMux.
func RegisterRoutes(mux *http.ServeMux, h *Handlers) {
	// Provider selection.
	mux.HandleFunc("GET /billing/choose-provider", h.Gateway.Show)
	mux.HandleFunc("POST /billing/choose-provider", h.Gateway.Store)

	// Billing subscription API.
	mux.HandleFunc("POST /billing/subscription", h.NewSub.Create)
	mux.HandleFunc("PUT /billing/subscription", h.UpdateSub.Update)
	mux.HandleFunc("PUT /billing/subscription/cancel", h.CancelSub.Cancel)
	mux.HandleFunc("PUT /billing/subscription/resume", h.ResumeSub.Resume)

	// Billing portal.
	mux.HandleFunc("GET /billing", h.Portal.Show)

	// Checkout routes.
	mux.HandleFunc("POST /checkout/subscribe", h.Checkout.Subscribe)
	mux.HandleFunc("POST /checkout/purchase", h.Checkout.Purchase)
	mux.HandleFunc("GET /billing/success", h.Checkout.Success)
	mux.HandleFunc("GET /billing/cancel", h.Checkout.Cancel)

	// Paddle billing portal.
	mux.HandleFunc("GET /billing/paddle", h.PaddleBilling.Index)
	mux.HandleFunc("POST /billing/paddle/cancel", h.PaddleBilling.Cancel)
	mux.HandleFunc("POST /billing/paddle/resume", h.PaddleBilling.Resume)

	// Orders.
	mux.HandleFunc("GET /billing/orders", h.Order.Index)

	// Invoices.
	mux.HandleFunc("GET /billing/{type}/{id}/invoices/{transaction}/download", h.Invoice.Download)

	// Payment methods.
	mux.HandleFunc("POST /billing/payment-method/setup", h.PaymentMethod.Setup)
	mux.HandleFunc("PUT /billing/payment-method/default", h.PaymentMethod.SetDefault)
	mux.HandleFunc("DELETE /billing/payment-method", h.PaymentMethod.Delete)
}
