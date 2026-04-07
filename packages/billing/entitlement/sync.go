package entitlement

import (
	"context"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/catalog"
)

// Synchronizer syncs subscription features from plan definitions.
type Synchronizer struct {
	seeder   *catalog.PlanSeeder
	plans    billing.PlanStore
	features billing.SubscriptionFeatureStore
	txn      billing.TransactionManager
	clock    billing.Clock
}

// NewSynchronizer creates a Synchronizer.
func NewSynchronizer(
	seeder *catalog.PlanSeeder,
	plans billing.PlanStore,
	features billing.SubscriptionFeatureStore,
	txn billing.TransactionManager,
	clock billing.Clock,
) *Synchronizer {
	return &Synchronizer{
		seeder:   seeder,
		plans:    plans,
		features: features,
		txn:      txn,
		clock:    clock,
	}
}

// Sync performs a full sync of subscription features from the plan definition.
func (s *Synchronizer) Sync(ctx context.Context, sub *billing.Subscription) error {
	return s.txn.RunInTransaction(ctx, func(ctx context.Context) error {
		plan, err := s.plans.FindByCode(ctx, sub.Plan)
		if err != nil {
			return err
		}

		if plan == nil {
			return s.features.DeleteAll(ctx, sub.ID)
		}

		var syncFeatures []billing.SubscriptionFeature

		for _, pf := range plan.PlanFeatures {
			syncFeatures = append(syncFeatures, billing.SubscriptionFeature{
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
func (s *Synchronizer) SyncActivationState(ctx context.Context, sub *billing.Subscription) error {
	return s.syncActivationState(ctx, sub)
}

func (s *Synchronizer) syncActivationState(ctx context.Context, sub *billing.Subscription) error {
	isActive := sub.Status.GrantsAccess()
	now := s.clock.Now()

	return s.features.SetActivationState(ctx, sub.ID, isActive, now)
}
