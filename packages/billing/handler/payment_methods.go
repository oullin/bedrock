package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bedrock/packages/billing"
)

// PaymentMethodsHandler handles payment method operations.
// Mirrors Billing\Http\Controllers\PaymentMethodsController.
type PaymentMethodsHandler struct {
	resolver billing.ResolverFunc
}

// NewPaymentMethodsHandler creates a PaymentMethodsHandler.
func NewPaymentMethodsHandler(resolver billing.ResolverFunc) *PaymentMethodsHandler {
	return &PaymentMethodsHandler{resolver: resolver}
}

// Setup creates a checkout session for adding a payment method.
func (h *PaymentMethodsHandler) Setup(w http.ResponseWriter, r *http.Request) {
	_, err := h.resolver(r)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

		return
	}

	// Provider-specific setup would go here.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
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
