package catalog

import (
	"context"

	"github.com/bedrock/packages/billing"
)

// PlanPeriodPriceManager handles price rotation — deactivating old prices
// and creating new ones when prices change.
type PlanPeriodPriceManager struct {
	prices billing.PriceStore
	txn    billing.TransactionManager
}

// NewPlanPeriodPriceManager creates a PlanPeriodPriceManager.
func NewPlanPeriodPriceManager(prices billing.PriceStore, txn billing.TransactionManager) *PlanPeriodPriceManager {
	return &PlanPeriodPriceManager{prices: prices, txn: txn}
}

// RotatePrice deactivates the current active price for a plan period and
// creates a new one if the provider price ID has changed.
func (m *PlanPeriodPriceManager) RotatePrice(ctx context.Context, planPeriodID int64, newPrice billing.PlanPeriodPrice) error {
	return m.txn.RunInTransaction(ctx, func(ctx context.Context) error {
		if err := m.prices.DeactivateForPeriod(ctx, planPeriodID); err != nil {
			return err
		}

		newPrice.PlanPeriodID = planPeriodID
		newPrice.IsActive = true

		return m.prices.Create(ctx, &newPrice)
	})
}
