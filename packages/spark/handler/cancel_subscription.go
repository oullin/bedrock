package handler

import (
	"net/http"

	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/spark/service"
)

// CancelSubscriptionHandler handles subscription cancellation.
// Mirrors Spark\Http\Controllers\CancelSubscriptionController.
type CancelSubscriptionHandler struct {
	billing  *service.BillingService
	resolver spark.ResolverFunc
}

// NewCancelSubscriptionHandler creates a CancelSubscriptionHandler.
func NewCancelSubscriptionHandler(b *service.BillingService, resolver spark.ResolverFunc) *CancelSubscriptionHandler {
	return &CancelSubscriptionHandler{billing: b, resolver: resolver}
}

// Cancel handles cancellation requests.
func (h *CancelSubscriptionHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	billable, err := h.resolver(r)

	if err != nil {
		errorResponse(w, http.StatusBadRequest, spark.ErrBillableRequired.Error())

		return
	}

	if err := h.billing.CancelSubscription(r.Context(), billable.BillableType(), billable.BillableID()); err != nil {
		errorResponse(w, http.StatusBadRequest, err.Error())

		return
	}

	noContent(w)
}
