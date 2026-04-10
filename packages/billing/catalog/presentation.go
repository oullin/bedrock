package catalog

import (
	"context"

	"github.com/bedrock/packages/billing"
)

// PlanPresentation builds display-ready plan data for landing pages and
// the billing portal.
type PlanPresentation struct {
	catalog *PlanCatalog
}

// NewPlanPresentation creates a PlanPresentation service.

// LandingPlans returns plans formatted for the landing/marketing page.

// BillingPlanRegistry registers plans from the database with the Billing Manager.
type BillingPlanRegistry struct {
	catalog *PlanCatalog
	manager *billing.Manager
}

func NewPlanPresentation(catalog *PlanCatalog) *PlanPresentation {
	return &PlanPresentation{catalog: catalog}
}

func (p *PlanPresentation) LandingPlans(ctx context.Context) ([]billing.LandingPlanView, error) {
	plans, err := p.catalog.Plans(ctx)

	if err != nil {
		return nil, err
	}

	var views []billing.LandingPlanView

	for _, plan := range plans {
		if !plan.IsActive {
			continue
		}

		var periods []billing.PlanPeriodView

		for _, pp := range plan.PlanPeriods {
			for _, price := range pp.Prices {
				if !price.IsActive {
					continue
				}

				periods = append(periods, billing.PlanPeriodView{
					Period:        string(pp.PeriodCode),
					DisplayLabel:  pp.PeriodCode.DisplayLabel(),
					DisplayAmount: price.DisplayAmount(nil),
					PricingMode:   string(price.PricingMode),
				})
			}
		}

		views = append(views, billing.LandingPlanView{
			Code:              plan.Code,
			Name:              plan.Name,
			Description:       plan.Description,
			MarketingFeatures: plan.MarketingFeatures,
			CTALabel:          plan.CTALabel,
			CTAStyle:          plan.CTAStyle,
			Featured:          plan.Featured,
			Badge:             plan.Badge,
			Periods:           periods,
		})
	}

	return views, nil
}

// NewBillingPlanRegistry creates a BillingPlanRegistry.
func NewBillingPlanRegistry(catalog *PlanCatalog, manager *billing.Manager) *BillingPlanRegistry {
	return &BillingPlanRegistry{catalog: catalog, manager: manager}
}

// Register reads active plans from the database and registers them with the
// Billing Manager for the given billable type.
func (r *BillingPlanRegistry) Register(ctx context.Context, billableType string) error {
	plans, err := r.catalog.Plans(ctx)

	if err != nil {
		return err
	}

	for _, plan := range plans {
		if !plan.IsActive {
			continue
		}

		for _, pp := range plan.PlanPeriods {
			for _, price := range pp.Prices {
				if !price.IsActive {
					continue
				}

				r.manager.AddPlan(billableType, billing.BillingPlan{
					ID:       price.ProviderPriceID,
					Name:     plan.Name,
					Interval: string(pp.PeriodCode),
					Active:   true,
				})
			}
		}
	}

	return nil
}
