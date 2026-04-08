package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/subscription"
)

// PortalHandler renders the billing portal page.
type PortalHandler struct {
	manager       *billing.Manager
	frontendState *subscription.FrontendStateBuilder
}

// NewPortalHandler creates a PortalHandler.
func NewPortalHandler(manager *billing.Manager, frontendState *subscription.FrontendStateBuilder) *PortalHandler {
	return &PortalHandler{manager: manager, frontendState: frontendState}
}

// Show renders the billing portal with full frontend state.
func (h *PortalHandler) Show(w http.ResponseWriter, r *http.Request) {
	billableType := r.PathValue("type")
	if billableType == "" {
		billableType = h.manager.DefaultBillableType()
	}

	billable, err := h.manager.ResolveBillable(billableType, r)
	if err != nil {
		http.Error(w, billing.ErrBillableRequired.Error(), http.StatusBadRequest)

		return
	}

	if !h.manager.IsAuthorized(billable, r) {
		http.Error(w, billing.ErrUnauthorized.Error(), http.StatusForbidden)

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
