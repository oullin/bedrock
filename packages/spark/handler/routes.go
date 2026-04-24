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

// RegisterRoutes registers all Spark routes on the given ServeMux.
func RegisterRoutes(mux *http.ServeMux, h *Handlers) {
	// Provider selection.
	mux.HandleFunc("GET /billing/choose-provider", h.Gateway.Show)
	mux.HandleFunc("POST /billing/choose-provider", h.Gateway.Store)

	// Spark subscription API.
	mux.HandleFunc("POST /spark/subscription", h.NewSub.Create)
	mux.HandleFunc("PUT /spark/subscription", h.UpdateSub.Update)
	mux.HandleFunc("PUT /spark/subscription/cancel", h.CancelSub.Cancel)
	mux.HandleFunc("PUT /spark/subscription/resume", h.ResumeSub.Resume)
	mux.HandleFunc("POST /spark/pending-checkout", h.Pending.Create)

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
	mux.HandleFunc("GET /spark/{type}/{id}/invoices/{transaction}/download", h.Invoice.Download)

	// Payment methods.
	mux.HandleFunc("PUT /spark/subscription/payment-method", h.PaymentMethod.Setup)
	mux.HandleFunc("POST /spark/payment-method/setup", h.PaymentMethod.Setup)
	mux.HandleFunc("PUT /spark/payment-method/default", h.PaymentMethod.SetDefault)
	mux.HandleFunc("DELETE /spark/payment-method", h.PaymentMethod.Delete)
}
