package billing

import "context"

// SubscriptionUpdateSyncer syncs plan/period changes after a provider update event.
type SubscriptionUpdateSyncer struct {
	subscriptions SubscriptionStore
	catalog       *PlanCatalog
	clock         Clock
}

// NewSubscriptionUpdateSyncer creates a SubscriptionUpdateSyncer.
func NewSubscriptionUpdateSyncer(
	subscriptions SubscriptionStore,
	catalog *PlanCatalog,
	clock Clock,
) *SubscriptionUpdateSyncer {
	return &SubscriptionUpdateSyncer{
		subscriptions: subscriptions,
		catalog:       catalog,
		clock:         clock,
	}
}

// SyncUpdated resolves the new price from provider items and updates the
// local subscription fields.
func (s *SubscriptionUpdateSyncer) SyncUpdated(ctx context.Context, sub *Subscription) error {
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
