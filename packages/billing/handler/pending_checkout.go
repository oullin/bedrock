package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bedrock/packages/billing"
)

// PendingCheckoutHandler records the provider checkout currently in flight.
// Mirrors Billing\Http\Controllers\NewPendingCheckoutController.
type PendingCheckoutHandler struct {
	customers billing.CustomerStore
	resolver  billing.ResolverFunc
}

// NewPendingCheckoutHandler creates a PendingCheckoutHandler.
func NewPendingCheckoutHandler(customers billing.CustomerStore, resolver billing.ResolverFunc) *PendingCheckoutHandler {
	return &PendingCheckoutHandler{customers: customers, resolver: resolver}
}

// Create persists the checkout ID for the resolved billable customer record.
func (h *PendingCheckoutHandler) Create(w http.ResponseWriter, r *http.Request) {
	billable, err := h.resolver(r)

	if err != nil {
		http.Error(w, billing.ErrBillableRequired.Error(), http.StatusBadRequest)

		return
	}

	var input struct {
		CheckoutID string `json:"checkout_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)

		return
	}

	customer, err := h.customers.FindByBillable(r.Context(), billable.BillableType(), billable.BillableID())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if customer == nil {
		customer = &billing.Customer{
			BillableType: billable.BillableType(),
			BillableID:   billable.BillableID(),
			Name:         billable.BillableName(),
			Email:        billable.BillableEmail(),
		}
		customer.PendingCheckoutID = input.CheckoutID

		if err := h.customers.Create(r.Context(), customer); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusNoContent)

		return
	}

	customer.PendingCheckoutID = input.CheckoutID

	if err := h.customers.Save(r.Context(), customer); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
