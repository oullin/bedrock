package handler

import (
	"net/http"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/service"
)

// ResumeSubscriptionHandler handles subscription resumption.
// Mirrors Billing\Http\Controllers\ResumeSubscriptionController.
type ResumeSubscriptionHandler struct {
	billing  *service.BillingService
	resolver billing.ResolverFunc
}

// NewResumeSubscriptionHandler creates a ResumeSubscriptionHandler.
func NewResumeSubscriptionHandler(b *service.BillingService, resolver billing.ResolverFunc) *ResumeSubscriptionHandler {
	return &ResumeSubscriptionHandler{billing: b, resolver: resolver}
}

// Resume handles resumption requests.
func (h *ResumeSubscriptionHandler) Resume(w http.ResponseWriter, r *http.Request) {
	billable, err := h.resolver(r)

	if err != nil {
		http.Error(w, billing.ErrBillableRequired.Error(), http.StatusBadRequest)

		return
	}

	if err := h.billing.ResumeSubscription(r.Context(), billable.BillableType(), billable.BillableID()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
