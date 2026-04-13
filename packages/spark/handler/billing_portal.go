package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/spark/state"
)

// PortalHandler renders the billing portal frontend state.
// Mirrors Spark\Http\Controllers\BillingPortalController.
type PortalHandler struct {
	manager  *spark.Manager
	frontend *state.FrontendState
	resolver spark.ResolverFunc
}

// NewPortalHandler creates a PortalHandler.
func NewPortalHandler(mgr *spark.Manager, fs *state.FrontendState, resolver spark.ResolverFunc) *PortalHandler {
	return &PortalHandler{manager: mgr, frontend: fs, resolver: resolver}
}

// Show returns the billing portal state as JSON.
func (h *PortalHandler) Show(w http.ResponseWriter, r *http.Request) {
	billableType := h.manager.DefaultBillableType()

	billable, err := h.resolver(r)

	if err != nil {
		http.Error(w, spark.ErrBillableRequired.Error(), http.StatusBadRequest)

		return
	}

	if !h.manager.IsAuthorized(billable, r) {
		http.Error(w, "Forbidden", http.StatusForbidden)

		return
	}

	data, err := h.frontend.Current(r.Context(), billableType, billable)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
