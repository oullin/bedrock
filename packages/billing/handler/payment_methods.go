package handler

import (
	"net/http"

	"github.com/bedrock/packages/billing"
)

// PaymentMethodsHandler handles payment method operations.
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
		errorResponse(w, http.StatusUnauthorized, "Unauthorized")

		return
	}

	subscription, err := h.subscriptions.CurrentForBillable(r.Context(), billable.BillableType(), billable.BillableID())

	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())

		return
	}

	if subscription == nil {
		errorResponse(w, http.StatusBadRequest, billing.ErrNotSubscribed.Error())

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
		errorResponse(w, http.StatusInternalServerError, err.Error())

		return
	}

	jsonResponse(w, http.StatusOK, map[string]any{
		"transaction_id": transaction.ID,
		"transaction":    transaction.Data,
	})
}

// SetDefault sets the default payment method.
func (h *PaymentMethodsHandler) SetDefault(w http.ResponseWriter, r *http.Request) {
	_, err := h.resolver(r)

	if err != nil {
		errorResponse(w, http.StatusUnauthorized, "Unauthorized")

		return
	}

	noContent(w)
}

// Delete removes a payment method.
func (h *PaymentMethodsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	_, err := h.resolver(r)

	if err != nil {
		errorResponse(w, http.StatusUnauthorized, "Unauthorized")

		return
	}

	noContent(w)
}
