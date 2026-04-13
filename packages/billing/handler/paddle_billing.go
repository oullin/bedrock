package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/service"
)

// PaddleBillingHandler handles the Paddle-specific billing portal.
// Mirrors app/Http/Controllers/PaddleBillingController.
type PaddleBillingHandler struct {
	billing  *service.BillingService
	resolver billing.ResolverFunc
}

// NewPaddleBillingHandler creates a PaddleBillingHandler.
func NewPaddleBillingHandler(b *service.BillingService, resolver billing.ResolverFunc) *PaddleBillingHandler {
	return &PaddleBillingHandler{billing: b, resolver: resolver}
}

// Index returns the Paddle billing portal state as JSON.
func (h *PaddleBillingHandler) Index(w http.ResponseWriter, r *http.Request) {
	billable, err := h.resolver(r)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

		return
	}

	ctx := r.Context()
	active, _ := h.billing.GetActiveSubscription(ctx, billable.BillableType(), billable.BillableID())

	var subData map[string]any

	if active != nil {
		subData = map[string]any{
			"provider":      active.Provider,
			"status":        active.Subscription.Status,
			"onGracePeriod": active.Subscription.OnGracePeriod(),
		}

		if active.Subscription.EndsAt != nil {
			subData["ends_at"] = active.Subscription.EndsAt.Format("2006-01-02")
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"subscription": subData,
	})
}

// Cancel cancels the Paddle subscription.
func (h *PaddleBillingHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	billable, err := h.resolver(r)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

		return
	}

	if err := h.billing.CancelSubscription(r.Context(), billable.BillableType(), billable.BillableID()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "cancelled"})
}

// Resume resumes a canceled Paddle subscription on grace period.
func (h *PaddleBillingHandler) Resume(w http.ResponseWriter, r *http.Request) {
	billable, err := h.resolver(r)

	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)

		return
	}

	if err := h.billing.ResumeSubscription(r.Context(), billable.BillableType(), billable.BillableID()); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "resumed"})
}
