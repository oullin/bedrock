package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bedrock/packages/billing"
)

// PaymentHandler handles payment method updates and pending checkouts.
type PaymentHandler struct {
	manager       *billing.Manager
	subscriptions billing.SubscriptionStore
	customers     billing.CustomerStore
	provider      billing.ProviderSubscriptionManager
}

// NewPaymentHandler creates a PaymentHandler.
func NewPaymentHandler(
	manager *billing.Manager,
	subscriptions billing.SubscriptionStore,
	customers billing.CustomerStore,
	provider billing.ProviderSubscriptionManager,
) *PaymentHandler {
	return &PaymentHandler{
		manager:       manager,
		subscriptions: subscriptions,
		customers:     customers,
		provider:      provider,
	}
}

// UpdatePaymentMethod returns a transaction ID for updating the payment method.
func (h *PaymentHandler) UpdatePaymentMethod(w http.ResponseWriter, r *http.Request) {
	billableType := r.FormValue("billableType")

	if billableType == "" {
		billableType = h.manager.DefaultBillableType()
	}

	billable, err := h.manager.ResolveBillable(billableType, r)

	if err != nil {
		http.Error(w, billing.ErrBillableRequired.Error(), http.StatusBadRequest)

		return
	}

	ctx := r.Context()
	sub, err := h.subscriptions.CurrentForBillable(ctx, billable.BillableType(), billable.BillableID())

	if err != nil || sub == nil {
		http.Error(w, billing.ErrNotSubscribed.Error(), http.StatusBadRequest)

		return
	}

	txnID, err := h.provider.PaymentMethodUpdateTransaction(ctx, sub.ProviderID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"transaction_id": txnID})
}

// NewPendingCheckout marks a pending checkout on the customer record.
func (h *PaymentHandler) NewPendingCheckout(w http.ResponseWriter, r *http.Request) {
	billableType := r.FormValue("billableType")

	if billableType == "" {
		billableType = h.manager.DefaultBillableType()
	}

	billable, err := h.manager.ResolveBillable(billableType, r)

	if err != nil {
		http.Error(w, billing.ErrBillableRequired.Error(), http.StatusBadRequest)

		return
	}

	checkoutID := r.FormValue("checkout_id")

	if checkoutID == "" {
		w.WriteHeader(http.StatusOK)

		return
	}

	ctx := r.Context()
	customer, err := h.customers.FindByBillable(ctx, billable.BillableType(), billable.BillableID())

	if err != nil || customer == nil {
		w.WriteHeader(http.StatusOK)

		return
	}

	customer.PendingCheckout = &checkoutID
	_ = h.customers.Save(ctx, customer)

	w.WriteHeader(http.StatusOK)
}
