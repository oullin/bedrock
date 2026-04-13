package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bedrock/packages/billing"
)

// OrderHandler handles order listing.
// Mirrors app/Http/Controllers/OrderController.
type OrderHandler struct {
	orders   billing.OrderStore
	resolver billing.ResolverFunc
}

// NewOrderHandler creates an OrderHandler.
func NewOrderHandler(orders billing.OrderStore, resolver billing.ResolverFunc) *OrderHandler {
	return &OrderHandler{orders: orders, resolver: resolver}
}

// Index returns the billable's orders as JSON.
func (h *OrderHandler) Index(w http.ResponseWriter, r *http.Request) {
	billable, err := h.resolver(r)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

		return
	}

	orders, err := h.orders.FindByBillable(r.Context(), billable.BillableID(), 15)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"orders": map[string]any{
			"data": orders,
		},
	})
}
