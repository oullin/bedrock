package action

import (
	"context"

	"github.com/bedrock/packages/billing"
)

// SubscriptionUpdater handles plan changes on existing subscriptions.
type SubscriptionUpdater struct {
	subscriptions billing.SubscriptionStore
	manager       *billing.Manager
}

// NewSubscriptionUpdater creates a SubscriptionUpdater.
func NewSubscriptionUpdater(subs billing.SubscriptionStore, mgr *billing.Manager) *SubscriptionUpdater {
	return &SubscriptionUpdater{subscriptions: subs, manager: mgr}
}

// Execute updates an existing subscription to a new plan.
func (u *SubscriptionUpdater) Execute(ctx context.Context, sub *billing.Subscription, plan *billing.Plan) error {
	// Update the subscription items to reflect the new plan price.
	sub.Items = []billing.SubscriptionItem{
		{
			SubscriptionID: sub.ID,
			PriceID:        plan.ID,
			Quantity:       1,
			Status:         string(billing.StatusActive),
		},
	}

	return u.subscriptions.Save(ctx, sub)
}
