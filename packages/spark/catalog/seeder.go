package catalog

import (
	"context"

	"github.com/bedrock/packages/spark"
)

// SeedData holds the complete plan catalog seed data.
type SeedData struct {
	Features []spark.Feature
	Plans    []spark.Plan
}

// PlanSeeder ensures the plan catalog data is persisted in the database.
type PlanSeeder struct {
	plans  spark.PlanStore
	prices spark.PriceStore
	txn    spark.TransactionManager
}

// NewPlanSeeder creates a PlanSeeder.
func NewPlanSeeder(plans spark.PlanStore, prices spark.PriceStore, txn spark.TransactionManager) *PlanSeeder {
	return &PlanSeeder{plans: plans, prices: prices, txn: txn}
}

// EnsurePersisted upserts all seed data (features, plans, periods, prices,
// plan-feature links) within a transaction.
func (s *PlanSeeder) EnsurePersisted(ctx context.Context, seed SeedData) error {
	return s.txn.RunInTransaction(ctx, func(ctx context.Context) error {
		if err := s.plans.UpsertFeatures(ctx, seed.Features); err != nil {
			return err
		}

		if err := s.plans.UpsertPlans(ctx, seed.Plans); err != nil {
			return err
		}

		for _, plan := range seed.Plans {
			if err := s.plans.UpsertPlanPeriods(ctx, plan.PlanPeriods); err != nil {
				return err
			}

			if err := s.plans.SyncPlanFeatures(ctx, plan.ID, plan.PlanFeatures); err != nil {
				return err
			}
		}

		return nil
	})
}
