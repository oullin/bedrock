package spark

import (
	"encoding/json"
	"net/http"
)

// BillingPortalHandler renders the billing portal page.
type BillingPortalHandler struct {
	manager       *Manager
	frontendState *FrontendStateBuilder
}

// NewBillingPortalHandler creates a BillingPortalHandler.
func NewBillingPortalHandler(manager *Manager, frontendState *FrontendStateBuilder) *BillingPortalHandler {
	return &BillingPortalHandler{manager: manager, frontendState: frontendState}
}

// Show renders the billing portal with full frontend state.
func (h *BillingPortalHandler) Show(w http.ResponseWriter, r *http.Request) {
	billableType := r.PathValue("type")
	if billableType == "" {
		billableType = h.manager.DefaultBillableType()
	}

	billable, err := h.manager.ResolveBillable(billableType, r)
	if err != nil {
		http.Error(w, ErrBillableRequired.Error(), http.StatusBadRequest)

		return
	}

	if !h.manager.IsAuthorized(billable, r) {
		http.Error(w, ErrUnauthorized.Error(), http.StatusForbidden)

		return
	}

	ctx := r.Context()
	state, err := h.frontendState.Current(ctx, billableType, billable)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}
