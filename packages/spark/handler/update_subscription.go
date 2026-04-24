package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/spark/action"
)

// UpdateSubscriptionHandler handles subscription plan changes.
// Mirrors Spark\Http\Controllers\UpdateSubscriptionController.
type UpdateSubscriptionHandler struct {
	updater       *action.SubscriptionUpdater
	manager       *spark.Manager
	subscriptions spark.SubscriptionStore
	resolver      spark.ResolverFunc
}

// NewUpdateSubscriptionHandler creates an UpdateSubscriptionHandler.
func NewUpdateSubscriptionHandler(
	u *action.SubscriptionUpdater,
	mgr *spark.Manager,
	subs spark.SubscriptionStore,
	resolver spark.ResolverFunc,
) *UpdateSubscriptionHandler {
	return &UpdateSubscriptionHandler{updater: u, manager: mgr, subscriptions: subs, resolver: resolver}
}

// Update handles plan change requests.
func (h *UpdateSubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	billable, err := h.resolver(r)

	if err != nil {
		errorResponse(w, http.StatusBadRequest, spark.ErrBillableRequired.Error())

		return
	}

	var input struct {
		Plan string `json:"plan"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		errorResponse(w, http.StatusBadRequest, "invalid request body")

		return
	}

	sub, err := h.subscriptions.CurrentForBillable(r.Context(), billable.BillableType(), billable.BillableID())

	if err != nil || sub == nil {
		errorResponse(w, http.StatusBadRequest, spark.ErrNotSubscribed.Error())

		return
	}

	plans := h.manager.Plans(billable.BillableType())

	var plan *spark.Plan

	for _, p := range plans {
		if p.ID == input.Plan {
			plan = p

			break
		}
	}

	if plan == nil {
		validationErrors(w, map[string][]string{"plan": {"The selected plan is invalid."}})

		return
	}

	if err := h.updater.Execute(r.Context(), sub, plan); err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())

		return
	}

	noContent(w)
}
