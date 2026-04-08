package handler

import (
	"encoding/json"
	"net/http"

	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/spark/billing"
)

// BillingHandler handles the Madora billing page and checkout endpoints.
type BillingHandler struct {
	billing  *billing.Workflow
	resolver spark.BillableResolver
}

// NewBillingHandler creates a BillingHandler.
func NewBillingHandler(b *billing.Workflow, resolver spark.BillableResolver) *BillingHandler {
	return &BillingHandler{billing: b, resolver: resolver}
}

// Show renders the billing page data as JSON.
func (h *BillingHandler) Show(w http.ResponseWriter, r *http.Request) {
	billable, err := h.resolver.Resolve(r)
	if err != nil {
		http.Error(w, spark.ErrBillableRequired.Error(), http.StatusBadRequest)

		return
	}

	ctx := r.Context()
	state := h.billing.ReadBillingState(ctx, billable.BillableType(), billable.BillableID())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(state)
}

// StartCheckout initiates a checkout from the billing page.
func (h *BillingHandler) StartCheckout(w http.ResponseWriter, r *http.Request) {
	billable, err := h.resolver.Resolve(r)
	if err != nil {
		http.Error(w, spark.ErrBillableRequired.Error(), http.StatusBadRequest)

		return
	}

	var input CheckoutInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)

		return
	}

	if errs := input.Validate(); errs != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]any{"errors": errs})

		return
	}

	ctx := r.Context()

	if input.Plan == spark.PlanStarter {
		if _, err := h.billing.StartStarterTrial(ctx, billable); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"redirect": "/subscription"})

		return
	}

	url, err := h.billing.StartCheckout(ctx, billable, input.Plan, input.Period)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]any{
			"errors": []spark.ValidationError{{Field: "plan", Message: "The selected plan pricing is currently unavailable."}},
		})

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"redirect": url})
}
