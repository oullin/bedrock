package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/spark/action"
)

// NewSubscriptionHandler handles creating new subscriptions.
// Mirrors Spark\Http\Controllers\NewSubscriptionController.
type NewSubscriptionHandler struct {
	creator  *action.SubscriptionCreator
	manager  *spark.Manager
	resolver spark.ResolverFunc
}

// NewNewSubscriptionHandler creates a NewSubscriptionHandler.
func NewNewSubscriptionHandler(c *action.SubscriptionCreator, mgr *spark.Manager, resolver spark.ResolverFunc) *NewSubscriptionHandler {
	return &NewSubscriptionHandler{creator: c, manager: mgr, resolver: resolver}
}

// Create handles subscription creation requests.
func (h *NewSubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	billable, err := h.resolver(r)

	if err != nil {
		http.Error(w, spark.ErrBillableRequired.Error(), http.StatusBadRequest)

		return
	}

	var input struct {
		Plan string `json:"plan"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)

		return
	}

	if !spark.ValidPlan(h.manager, billable.BillableType(), input.Plan) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": map[string][]string{
				"plan": {"The selected plan is invalid."},
			},
		})

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

	checkout, err := h.creator.Execute(r.Context(), billable, plan, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(checkout.ToMap())
}
