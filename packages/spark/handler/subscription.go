package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/bedrock/packages/spark"
)

// SubscriptionHandler handles Spark vendor subscription operations (create,
// update, cancel, resume).
type SubscriptionHandler struct {
	manager       *spark.Manager
	subscriptions spark.SubscriptionStore
	provider      spark.ProviderSubscriptionManager
	checkout      spark.ProviderCheckoutGenerator
}

// NewSubscriptionHandler creates a SubscriptionHandler.
func NewSubscriptionHandler(
	manager *spark.Manager,
	subscriptions spark.SubscriptionStore,
	provider spark.ProviderSubscriptionManager,
	checkout spark.ProviderCheckoutGenerator,
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
		http.Error(w, spark.ErrBillableRequired.Error(), http.StatusBadRequest)

		return
	}

	planID := r.FormValue("plan")
	if !spark.ValidPlan(h.manager, billableType, planID) {
		http.Error(w, "invalid plan", http.StatusUnprocessableEntity)

		return
	}

	// Find the plan.
	var plan spark.SparkPlan
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

	items := []spark.CheckoutItem{{PriceID: planID, Quantity: 1}}
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
		http.Error(w, spark.ErrBillableRequired.Error(), http.StatusBadRequest)

		return
	}

	planID := r.FormValue("plan")
	if !spark.ValidPlan(h.manager, billableType, planID) {
		http.Error(w, "invalid plan", http.StatusUnprocessableEntity)

		return
	}

	var plan spark.SparkPlan
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
		http.Error(w, spark.ErrNotSubscribed.Error(), http.StatusBadRequest)

		return
	}

	items := []spark.CheckoutItem{{PriceID: planID, Quantity: 1}}
	if h.manager.ChargesPerSeat(billableType) {
		items[0].Quantity = h.manager.SeatCount(billableType, billable)
	}

	proration := spark.ProratedNextBillingPeriod
	if h.manager.Prorates() {
		proration = spark.ProratedImmediately
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
		http.Error(w, spark.ErrBillableRequired.Error(), http.StatusBadRequest)

		return
	}

	ctx := r.Context()
	sub, err := h.subscriptions.CurrentForBillable(ctx, billable.BillableType(), billable.BillableID())
	if err != nil || sub == nil {
		http.Error(w, spark.ErrNotSubscribed.Error(), http.StatusBadRequest)

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
		http.Error(w, spark.ErrBillableRequired.Error(), http.StatusBadRequest)

		return
	}

	ctx := r.Context()
	sub, err := h.subscriptions.CurrentForBillable(ctx, billable.BillableType(), billable.BillableID())
	if err != nil || sub == nil {
		http.Error(w, spark.ErrNotSubscribed.Error(), http.StatusBadRequest)

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
			_ = h.provider.UpdateQuantity(ctx, sub.ProviderID, priceID, quantity, spark.ProratedNextBillingPeriod)
		}
	}

	w.WriteHeader(http.StatusOK)
}

// findOrCreateCustomer is a helper that resolves the customer for a billable.
func findOrCreateCustomer(ctx context.Context, billable spark.Billable, subs spark.SubscriptionStore) (*spark.Customer, error) {
	_ = ctx
	_ = subs

	return &spark.Customer{
		BillableType: billable.BillableType(),
		BillableID:   billable.BillableID(),
		Name:         billable.BillableName(),
		Email:        billable.BillableEmail(),
	}, nil
}
