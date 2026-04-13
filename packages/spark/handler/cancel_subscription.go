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
		http.Error(w, spark.ErrBillableRequired.Error(), http.StatusBadRequest)

		return
	}

	if err := h.billing.CancelSubscription(r.Context(), billable.BillableType(), billable.BillableID()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
