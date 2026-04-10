package catalog

import (
	"context"

	"github.com/bedrock/packages/billing"
)

// PlanCatalog provides access to the plan catalog.
type PlanCatalog struct {
	plans  billing.PlanStore
	prices billing.PriceStore
}

// NewPlanCatalog creates a PlanCatalog service.
func NewPlanCatalog(plans billing.PlanStore, prices billing.PriceStore) *PlanCatalog {
	return &PlanCatalog{plans: plans, prices: prices}
}

// Plans returns all active plans.
func (c *PlanCatalog) Plans(ctx context.Context) ([]billing.Plan, error) {
	return c.plans.ActivePlans(ctx)
}

// PlanByCode returns the plan matching the given code.
func (c *PlanCatalog) PlanByCode(ctx context.Context, code string) (*billing.Plan, error) {
	return c.plans.FindByCode(ctx, code)
}

// PlanBySlug returns the plan matching the given slug.
func (c *PlanCatalog) PlanBySlug(ctx context.Context, slug string) (*billing.Plan, error) {
	return c.plans.FindBySlug(ctx, slug)
}

// MarketingName returns the human-readable name for a plan code.
func (c *PlanCatalog) MarketingName(ctx context.Context, code string) string {
	plan, err := c.plans.FindByCode(ctx, code)

	if err != nil || plan == nil {
		return code
	}

	return plan.Name
}

// DefaultPeriodForPlan returns the default billing period for a plan.
// If the plan has a yearly price, it defaults to yearly; otherwise monthly.
func (c *PlanCatalog) DefaultPeriodForPlan(ctx context.Context, planCode string) billing.BillingPeriod {
	plan, err := c.plans.FindByCode(ctx, planCode)

	if err != nil || plan == nil {
		return billing.PeriodMonthly
	}

	for _, pp := range plan.PlanPeriods {
		if pp.PeriodCode == billing.PeriodYearly {
			price, _ := c.prices.ActiveForPeriod(ctx, pp.ID)

			if price != nil {
				return billing.PeriodYearly
			}
		}
	}

	return billing.PeriodMonthly
}

// ActivePriceFor returns the active price for a plan and billing period.
func (c *PlanCatalog) ActivePriceFor(ctx context.Context, planCode string, period billing.BillingPeriod) (*billing.PlanPeriodPrice, error) {
	plan, err := c.plans.FindByCode(ctx, planCode)

	if err != nil {
		return nil, err
	}

	if plan == nil {
		return nil, billing.ErrNotFound
	}

	pp, err := c.plans.FindPlanPeriod(ctx, plan.ID, period)

	if err != nil {
		return nil, err
	}

	if pp == nil {
		return nil, nil
	}

	return c.prices.ActiveForPeriod(ctx, pp.ID)
}

// PriceForProviderID returns the price matching a provider price ID.
func (c *PlanCatalog) PriceForProviderID(ctx context.Context, providerPriceID string) (*billing.PlanPeriodPrice, error) {
	return c.prices.FindByProviderPriceID(ctx, providerPriceID)
}
