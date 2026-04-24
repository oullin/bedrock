package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bedrock/packages/billing"
)

// PaymentMethodsHandler handles payment method operations.
// Mirrors Billing\Http\Controllers\PaymentMethodsController.
type PaymentMethodsHandler struct {
	resolver      billing.ResolverFunc
	subscriptions billing.SubscriptionStore
	provider      billing.PaymentMethodUpdater
	manager       *billing.Manager
}

// NewPaymentMethodsHandler creates a PaymentMethodsHandler.
func NewPaymentMethodsHandler(
	resolver billing.ResolverFunc,
	subscriptions billing.SubscriptionStore,
	provider billing.PaymentMethodUpdater,
	manager *billing.Manager,
) *PaymentMethodsHandler {
	return &PaymentMethodsHandler{
		resolver:      resolver,
		subscriptions: subscriptions,
		provider:      provider,
		manager:       manager,
	}
}

// Setup creates a provider transaction for updating the current payment method.
func (h *PaymentMethodsHandler) Setup(w http.ResponseWriter, r *http.Request) {
	billable, err := h.resolver(r)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

		return
	}

	subscription, err := h.subscriptions.CurrentForBillable(r.Context(), billable.BillableType(), billable.BillableID())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if subscription == nil {
		http.Error(w, billing.ErrNotSubscribed.Error(), http.StatusBadRequest)

		return
	}

	options := map[string]any{}

	if h.manager != nil {
		if callback := h.manager.PaymentMethodSessionOptions(billable.BillableType()); callback != nil {
			options = callback(billable)
		}
	}

	transaction, err := h.provider.CreatePaymentMethodUpdateTransaction(r.Context(), billable, subscription, options)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"transaction_id": transaction.ID,
		"transaction":    transaction.Data,
	})
}

// SetDefault sets the default payment method.
func (h *PaymentMethodsHandler) SetDefault(w http.ResponseWriter, r *http.Request) {
	_, err := h.resolver(r)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Delete removes a payment method.
func (h *PaymentMethodsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	_, err := h.resolver(r)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
