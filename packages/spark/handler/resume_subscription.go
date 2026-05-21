package handler

import (
	"net/http"

	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/spark/service"
)

// ResumeSubscriptionHandler handles subscription resumption.
type ResumeSubscriptionHandler struct {
	billing  *service.BillingService
	resolver spark.ResolverFunc
}

// NewResumeSubscriptionHandler creates a ResumeSubscriptionHandler.
func NewResumeSubscriptionHandler(b *service.BillingService, resolver spark.ResolverFunc) *ResumeSubscriptionHandler {
	return &ResumeSubscriptionHandler{billing: b, resolver: resolver}
}

// Resume handles resumption requests.
func (h *ResumeSubscriptionHandler) Resume(w http.ResponseWriter, r *http.Request) {
	billable, err := h.resolver(r)

	if err != nil {
		errorResponse(w, http.StatusBadRequest, spark.ErrBillableRequired.Error())

		return
	}

	if err := h.billing.ResumeSubscription(r.Context(), billable.BillableType(), billable.BillableID()); err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())

		return
	}

	noContent(w)
}
