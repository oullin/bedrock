package action

import (
	"context"

	"github.com/bedrock/packages/spark"
)

// SubscriptionUpdater handles plan changes on existing subscriptions.
type SubscriptionUpdater struct {
	subscriptions spark.SubscriptionStore
	manager       *spark.Manager
}

// NewSubscriptionUpdater creates a SubscriptionUpdater.
func NewSubscriptionUpdater(subs spark.SubscriptionStore, mgr *spark.Manager) *SubscriptionUpdater {
	return &SubscriptionUpdater{subscriptions: subs, manager: mgr}
}

// Execute updates an existing subscription to a new plan.
func (u *SubscriptionUpdater) Execute(ctx context.Context, sub *spark.Subscription, plan *spark.Plan) error {
	// Update the subscription items to reflect the new plan price.
	sub.Items = []spark.SubscriptionItem{
		{
			SubscriptionID: sub.ID,
			PriceID:        plan.ID,
			Quantity:       1,
			Status:         string(spark.StatusActive),
		},
	}

	return u.subscriptions.Save(ctx, sub)
}
