package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/action"
)

// UpdateSubscriptionHandler handles subscription plan changes.
// Mirrors Billing\Http\Controllers\UpdateSubscriptionController.
type UpdateSubscriptionHandler struct {
	updater       *action.SubscriptionUpdater
	manager       *billing.Manager
	subscriptions billing.SubscriptionStore
	resolver      billing.ResolverFunc
}

// NewUpdateSubscriptionHandler creates an UpdateSubscriptionHandler.
func NewUpdateSubscriptionHandler(
	u *action.SubscriptionUpdater,
	mgr *billing.Manager,
	subs billing.SubscriptionStore,
	resolver billing.ResolverFunc,
) *UpdateSubscriptionHandler {
	return &UpdateSubscriptionHandler{updater: u, manager: mgr, subscriptions: subs, resolver: resolver}
}

// Update handles plan change requests.
func (h *UpdateSubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	billable, err := h.resolver(r)

	if err != nil {
		http.Error(w, billing.ErrBillableRequired.Error(), http.StatusBadRequest)

		return
	}

	var input struct {
		Plan string `json:"plan"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)

		return
	}

	sub, err := h.subscriptions.CurrentForBillable(r.Context(), billable.BillableType(), billable.BillableID())

	if err != nil || sub == nil {
		http.Error(w, billing.ErrNotSubscribed.Error(), http.StatusBadRequest)

		return
	}

	plans := h.manager.Plans(billable.BillableType())

	var plan *billing.Plan

	for _, p := range plans {
		if p.ID == input.Plan {
			plan = p

			break
		}
	}

	if plan == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": map[string][]string{
				"plan": {"The selected plan is invalid."},
			},
		})

		return
	}

	if err := h.updater.Execute(r.Context(), sub, plan); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
