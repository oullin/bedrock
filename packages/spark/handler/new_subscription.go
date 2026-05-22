package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/spark/action"
)

// NewSubscriptionHandler handles creating new subscriptions.
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

	if !spark.ValidPlan(h.manager, billable.BillableType(), input.Plan) {
		validationErrors(w, map[string][]string{"plan": {"The selected plan is invalid."}})

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
		errorResponse(w, http.StatusInternalServerError, err.Error())

		return
	}

	jsonResponse(w, http.StatusOK, checkout.ToMap())
}
