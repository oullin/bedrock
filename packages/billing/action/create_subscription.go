// Package action provides business logic actions for billing operations.
package action

import (
	"context"

	"github.com/bedrock/packages/billing"
)

// SubscriptionCreator handles creating new subscriptions.
// Mirrors Billing\Actions\CreateSubscription.
type SubscriptionCreator struct {
	subscriptions billing.SubscriptionStore
	manager       *billing.Manager
}

// NewSubscriptionCreator creates a SubscriptionCreator.
func NewSubscriptionCreator(subs billing.SubscriptionStore, mgr *billing.Manager) *SubscriptionCreator {
	return &SubscriptionCreator{subscriptions: subs, manager: mgr}
}

// Execute creates a new subscription checkout for the billable.
func (c *SubscriptionCreator) Execute(ctx context.Context, billable billing.Billable, plan *billing.Plan, opts map[string]any) (*billing.Checkout, error) {
	if err := c.manager.EnsurePlanEligibility(billable, *plan); err != nil {
		return nil, err
	}

	items := []billing.CheckoutItem{
		{PriceID: plan.ID, Quantity: 1},
	}

	// Apply per-seat quantity.
	if c.manager.ChargesPerSeat(billable.BillableType()) {
		items[0].Quantity = c.manager.SeatCount(billable.BillableType(), billable)
	}

	checkout := billing.GuestCheckout(items)

	if returnURL, ok := opts["return_url"].(string); ok {
		checkout.WithReturnURL(returnURL)
	}

	return checkout, nil
}
