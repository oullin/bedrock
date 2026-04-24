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
		errorResponse(w, http.StatusBadRequest, billing.ErrBillableRequired.Error())

		return
	}

	if err := h.billing.ResumeSubscription(r.Context(), billable.BillableType(), billable.BillableID()); err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())

		return
	}

	noContent(w)
}
