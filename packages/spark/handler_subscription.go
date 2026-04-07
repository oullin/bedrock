package spark

import (
	"context"
	"encoding/json"
	"net/http"
)

// SubscriptionHandler handles Spark vendor subscription operations (create,
// update, cancel, resume).
type SubscriptionHandler struct {
	manager       *Manager
	subscriptions SubscriptionStore
	provider      ProviderSubscriptionManager
	checkout      ProviderCheckoutGenerator
}

// NewSubscriptionHandler creates a SubscriptionHandler.
func NewSubscriptionHandler(
	manager *Manager,
	subscriptions SubscriptionStore,
	provider ProviderSubscriptionManager,
	checkout ProviderCheckoutGenerator,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		manager:       manager,
		subscriptions: subscriptions,
		provider:      provider,
		checkout:      checkout,
	}
}

// Create handles new subscription creation.
func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	billableType := r.FormValue("billableType")
	if billableType == "" {
		billableType = h.manager.DefaultBillableType()
	}

	billable, err := h.manager.ResolveBillable(billableType, r)
	if err != nil {
		http.Error(w, ErrBillableRequired.Error(), http.StatusBadRequest)

		return
	}

	planID := r.FormValue("plan")
	if !ValidPlan(h.manager, billableType, planID) {
		http.Error(w, "invalid plan", http.StatusUnprocessableEntity)

		return
	}

	// Find the plan.
	var plan SparkPlan
	for _, p := range h.manager.Plans(billableType) {
		if p.ID == planID {
			plan = p

			break
		}
	}

	if err := h.manager.EnsurePlanEligibility(billable, plan); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)

		return
	}

	ctx := r.Context()
	customer, err := findOrCreateCustomer(ctx, billable, h.subscriptions)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	items := []CheckoutItem{{PriceID: planID, Quantity: 1}}
	if h.manager.ChargesPerSeat(billableType) {
		items[0].Quantity = h.manager.SeatCount(billableType, billable)
	}

	session, err := h.checkout.Generate(ctx, customer, items, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(session)
}

// Update handles subscription plan changes.
func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	billableType := r.FormValue("billableType")
	if billableType == "" {
		billableType = h.manager.DefaultBillableType()
	}

	billable, err := h.manager.ResolveBillable(billableType, r)
	if err != nil {
		http.Error(w, ErrBillableRequired.Error(), http.StatusBadRequest)

		return
	}

	planID := r.FormValue("plan")
	if !ValidPlan(h.manager, billableType, planID) {
		http.Error(w, "invalid plan", http.StatusUnprocessableEntity)

		return
	}

	var plan SparkPlan
	for _, p := range h.manager.Plans(billableType) {
		if p.ID == planID {
			plan = p

			break
		}
	}

	if err := h.manager.EnsurePlanEligibility(billable, plan); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)

		return
	}

	ctx := r.Context()
	sub, err := h.subscriptions.CurrentForBillable(ctx, billable.BillableType(), billable.BillableID())
	if err != nil || sub == nil {
		http.Error(w, ErrNotSubscribed.Error(), http.StatusBadRequest)

		return
	}

	items := []CheckoutItem{{PriceID: planID, Quantity: 1}}
	if h.manager.ChargesPerSeat(billableType) {
		items[0].Quantity = h.manager.SeatCount(billableType, billable)
	}

	proration := ProratedNextBillingPeriod
	if h.manager.Prorates() {
		proration = ProratedImmediately
	}

	if err := h.provider.Swap(ctx, sub.ProviderID, items, proration); err != nil {
		http.Error(w, err.Error(), http.StatusUnprocessableEntity)

		return
	}

	w.WriteHeader(http.StatusOK)
}

// Cancel handles subscription cancellation.
func (h *SubscriptionHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	billableType := r.FormValue("billableType")
	if billableType == "" {
		billableType = h.manager.DefaultBillableType()
	}

	billable, err := h.manager.ResolveBillable(billableType, r)
	if err != nil {
		http.Error(w, ErrBillableRequired.Error(), http.StatusBadRequest)

		return
	}

	ctx := r.Context()
	sub, err := h.subscriptions.CurrentForBillable(ctx, billable.BillableType(), billable.BillableID())
	if err != nil || sub == nil {
		http.Error(w, ErrNotSubscribed.Error(), http.StatusBadRequest)

		return
	}

	if err := h.provider.Cancel(ctx, sub.ProviderID, false); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

// Resume handles subscription resumption within a grace period.
func (h *SubscriptionHandler) Resume(w http.ResponseWriter, r *http.Request) {
	billableType := r.FormValue("billableType")
	if billableType == "" {
		billableType = h.manager.DefaultBillableType()
	}

	billable, err := h.manager.ResolveBillable(billableType, r)
	if err != nil {
		http.Error(w, ErrBillableRequired.Error(), http.StatusBadRequest)

		return
	}

	ctx := r.Context()
	sub, err := h.subscriptions.CurrentForBillable(ctx, billable.BillableType(), billable.BillableID())
	if err != nil || sub == nil {
		http.Error(w, ErrNotSubscribed.Error(), http.StatusBadRequest)

		return
	}

	if err := h.provider.StopCancellation(ctx, sub.ProviderID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if h.manager.ChargesPerSeat(billableType) {
		quantity := h.manager.SeatCount(billableType, billable)
		priceID := ""
		if len(sub.Items) > 0 {
			priceID = sub.Items[0].PriceID
		}

		if priceID != "" {
			_ = h.provider.UpdateQuantity(ctx, sub.ProviderID, priceID, quantity, ProratedNextBillingPeriod)
		}
	}

	w.WriteHeader(http.StatusOK)
}

// findOrCreateCustomer is a helper that resolves the customer for a billable.
func findOrCreateCustomer(ctx context.Context, billable Billable, subs SubscriptionStore) (*Customer, error) {
	_ = ctx
	_ = subs

	return &Customer{
		BillableType: billable.BillableType(),
		BillableID:   billable.BillableID(),
		Name:         billable.BillableName(),
		Email:        billable.BillableEmail(),
	}, nil
}
