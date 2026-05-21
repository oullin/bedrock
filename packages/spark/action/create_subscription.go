// Package action provides business logic actions for billing operations.
package action

import (
	"context"

	"github.com/bedrock/packages/spark"
)

// SubscriptionCreator handles creating new subscriptions.
type SubscriptionCreator struct {
	subscriptions spark.SubscriptionStore
	manager       *spark.Manager
}

// NewSubscriptionCreator creates a SubscriptionCreator.
func NewSubscriptionCreator(subs spark.SubscriptionStore, mgr *spark.Manager) *SubscriptionCreator {
	return &SubscriptionCreator{subscriptions: subs, manager: mgr}
}

// Execute creates a new subscription checkout for the billable.
func (c *SubscriptionCreator) Execute(ctx context.Context, billable spark.Billable, plan *spark.Plan, opts map[string]any) (*spark.Checkout, error) {
	if err := c.manager.EnsurePlanEligibility(billable, *plan); err != nil {
		return nil, err
	}

	items := []spark.CheckoutItem{
		{PriceID: plan.ID, Quantity: 1},
	}

	// Apply per-seat quantity.
	if c.manager.ChargesPerSeat(billable.BillableType()) {
		items[0].Quantity = c.manager.SeatCount(billable.BillableType(), billable)
	}

	checkout := spark.GuestCheckout(items)

	if returnURL, ok := opts["return_url"].(string); ok {
		checkout.WithReturnURL(returnURL)
	}

	return checkout, nil
}
