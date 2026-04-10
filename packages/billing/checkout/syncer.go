package checkout

import (
	"context"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/catalog"
)

// Syncer syncs plan/period changes after a provider update event.
type Syncer struct {
	subscriptions billing.SubscriptionStore
	catalog       *catalog.PlanCatalog
	clock         billing.Clock
}

// NewSyncer creates a Syncer.
func NewSyncer(
	subscriptions billing.SubscriptionStore,
	cat *catalog.PlanCatalog,
	clock billing.Clock,
) *Syncer {
	return &Syncer{
		subscriptions: subscriptions,
		catalog:       cat,
		clock:         clock,
	}
}

// SyncUpdated resolves the new price from provider items and updates the
// local subscription fields.
func (s *Syncer) SyncUpdated(ctx context.Context, sub *billing.Subscription) error {
	if sub == nil || len(sub.Items) == 0 {
		return nil
	}

	price, err := s.catalog.PriceForProviderID(ctx, sub.Items[0].PriceID)

	if err != nil || price == nil {
		return err
	}

	sub.PlanPeriodPriceID = &price.ID
	sub.UpdatedAt = s.clock.Now()

	return s.subscriptions.Save(ctx, sub)
}
