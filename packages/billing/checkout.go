package billing

import "context"

// CheckoutStarter initiates the checkout flow for a billable.
type CheckoutStarter struct {
	catalog       *PlanCatalog
	subscriptions *SubscriptionRepository
}

// NewCheckoutStarter creates a CheckoutStarter.
func NewCheckoutStarter(catalog *PlanCatalog, subscriptions *SubscriptionRepository) *CheckoutStarter {
	return &CheckoutStarter{catalog: catalog, subscriptions: subscriptions}
}

// Start initiates a checkout for the given plan and period. For the starter
// plan, it creates a trial directly. For paid plans, it validates pricing
// and delegates to the subscription repository.
func (s *CheckoutStarter) Start(
	ctx context.Context,
	billable Billable,
	plan SubscriptionPlan,
	period BillingPeriod,
) (string, error) {
	if plan == PlanStarter {
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
		return "", ErrNoActivePrice
	}

	return s.subscriptions.BeginCheckout(ctx, billable, plan, period, price)
}
