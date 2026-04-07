package spark

import "context"

// EntitlementSynchronizer syncs subscription features from plan definitions.
type EntitlementSynchronizer struct {
	seeder   *PlanSeeder
	plans    PlanStore
	features SubscriptionFeatureStore
	txn      TransactionManager
	clock    Clock
}

// NewEntitlementSynchronizer creates an EntitlementSynchronizer.
func NewEntitlementSynchronizer(
	seeder *PlanSeeder,
	plans PlanStore,
	features SubscriptionFeatureStore,
	txn TransactionManager,
	clock Clock,
) *EntitlementSynchronizer {
	return &EntitlementSynchronizer{
		seeder:   seeder,
		plans:    plans,
		features: features,
		txn:      txn,
		clock:    clock,
	}
}

// Sync performs a full sync of subscription features from the plan definition.
func (s *EntitlementSynchronizer) Sync(ctx context.Context, sub *Subscription) error {
	return s.txn.RunInTransaction(ctx, func(ctx context.Context) error {
		plan, err := s.plans.FindByCode(ctx, sub.Plan)
		if err != nil {
			return err
		}

		if plan == nil {
			return s.features.DeleteAll(ctx, sub.ID)
		}

		var syncFeatures []SubscriptionFeature

		for _, pf := range plan.PlanFeatures {
			syncFeatures = append(syncFeatures, SubscriptionFeature{
				SubscriptionID: sub.ID,
				FeatureID:      pf.FeatureID,
				Value:          pf.Value,
			})
		}

		if err := s.features.Sync(ctx, sub.ID, syncFeatures); err != nil {
			return err
		}

		return s.syncActivationState(ctx, sub)
	})
}

// SyncActivationState updates the activation state of all features based on
// the subscription status.
func (s *EntitlementSynchronizer) SyncActivationState(ctx context.Context, sub *Subscription) error {
	return s.syncActivationState(ctx, sub)
}

func (s *EntitlementSynchronizer) syncActivationState(ctx context.Context, sub *Subscription) error {
	isActive := sub.Status.GrantsAccess()
	now := s.clock.Now()

	return s.features.SetActivationState(ctx, sub.ID, isActive, now)
}
