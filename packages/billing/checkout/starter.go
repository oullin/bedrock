package checkout

import (
	"context"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/catalog"
	"github.com/bedrock/packages/billing/subscription"
)

// Starter initiates the checkout flow for a billable.
type Starter struct {
	catalog       *catalog.PlanCatalog
	subscriptions *subscription.Repository
}

// NewStarter creates a Starter.
func NewStarter(cat *catalog.PlanCatalog, subscriptions *subscription.Repository) *Starter {
	return &Starter{catalog: cat, subscriptions: subscriptions}
}

// Start initiates a checkout for the given plan and period. For the starter
// plan, it creates a trial directly. For paid plans, it validates pricing
// and delegates to the subscription repository.
func (s *Starter) Start(
	ctx context.Context,
	billable billing.Billable,
	plan billing.SubscriptionPlan,
	period billing.BillingPeriod,
) (string, error) {
	if plan == billing.PlanStarter {
		_, err := s.subscriptions.StartStarterTrial(ctx, billable)
		if err != nil {
			return "", err
		}

		return "", nil // Redirect handled by caller.
	}

	price, err := s.catalog.ActivePriceFor(ctx, string(plan), period)
	if err != nil {
		return "", err
	}

	if price == nil {
		return "", billing.ErrNoActivePrice
	}

	return s.subscriptions.BeginCheckout(ctx, billable, plan, period, price)
}
