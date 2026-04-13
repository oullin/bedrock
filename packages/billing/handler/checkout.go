package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/service"
)

// CheckoutHandler handles product checkout flows for both subscription
// and one-time purchases. Mirrors app/Http/Controllers/CheckoutController.
type CheckoutHandler struct {
	billing  *service.BillingService
	products billing.ProductStore
	resolver billing.ResolverFunc
}

// NewCheckoutHandler creates a CheckoutHandler.
func NewCheckoutHandler(b *service.BillingService, products billing.ProductStore, resolver billing.ResolverFunc) *CheckoutHandler {
	return &CheckoutHandler{billing: b, products: products, resolver: resolver}
}

// Subscribe handles subscription checkout requests.
func (h *CheckoutHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	billable, err := h.resolver(r)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

		return
	}

	var input struct {
		ProductID int64  `json:"product_id"`
		Provider  string `json:"provider"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)

		return
	}

	product, err := h.products.FindByID(r.Context(), input.ProductID)

	if err != nil || product == nil {
		http.Error(w, "Product not found", http.StatusNotFound)

		return
	}

	priceID := product.PriceIDForProvider(input.Provider)

	if priceID == "" {
		http.Error(w, "No price configured for provider", http.StatusBadRequest)

		return
	}

	_ = billable // checkout creation would use billable + provider SDK
	_ = priceID

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Purchase handles one-time purchase checkout requests.
func (h *CheckoutHandler) Purchase(w http.ResponseWriter, r *http.Request) {
	billable, err := h.resolver(r)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

		return
	}

	var input struct {
		ProductID int64  `json:"product_id"`
		Provider  string `json:"provider"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)

		return
	}

	product, err := h.products.FindByID(r.Context(), input.ProductID)

	if err != nil || product == nil {
		http.Error(w, "Product not found", http.StatusNotFound)

		return
	}

	_, err = h.billing.CreateOneTimeOrder(r.Context(), billable.BillableID(), product, input.Provider)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Success renders the checkout success page.
func (h *CheckoutHandler) Success(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":   "success",
		"provider": r.URL.Query().Get("provider"),
		"order":    r.URL.Query().Get("order"),
	})
}

// Cancel renders the checkout cancellation page.
func (h *CheckoutHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "cancelled"})
}
