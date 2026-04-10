package catalog

import (
	"context"

	"github.com/bedrock/packages/spark"
)

// PlanPresentation builds display-ready plan data for landing pages and
// the billing portal.
type PlanPresentation struct {
	catalog *PlanCatalog
}

// NewPlanPresentation creates a PlanPresentation service.

// LandingPlans returns plans formatted for the landing/marketing page.

// SparkPlanRegistry registers plans from the database with the Spark Manager.
type SparkPlanRegistry struct {
	catalog *PlanCatalog
	manager *spark.Manager
}

func NewPlanPresentation(catalog *PlanCatalog) *PlanPresentation {
	return &PlanPresentation{catalog: catalog}
}

func (p *PlanPresentation) LandingPlans(ctx context.Context) ([]spark.LandingPlanView, error) {
	plans, err := p.catalog.Plans(ctx)

	if err != nil {
		return nil, err
	}

	var views []spark.LandingPlanView

	for _, plan := range plans {
		if !plan.IsActive {
			continue
		}

		var periods []spark.PlanPeriodView

		for _, pp := range plan.PlanPeriods {
			for _, price := range pp.Prices {
				if !price.IsActive {
					continue
				}

				periods = append(periods, spark.PlanPeriodView{
					Period:        string(pp.PeriodCode),
					DisplayLabel:  pp.PeriodCode.DisplayLabel(),
					DisplayAmount: price.DisplayAmount(nil),
					PricingMode:   string(price.PricingMode),
				})
			}
		}

		views = append(views, spark.LandingPlanView{
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

// NewSparkPlanRegistry creates a SparkPlanRegistry.
func NewSparkPlanRegistry(catalog *PlanCatalog, manager *spark.Manager) *SparkPlanRegistry {
	return &SparkPlanRegistry{catalog: catalog, manager: manager}
}

// Register reads active plans from the database and registers them with the
// Spark Manager for the given billable type.
func (r *SparkPlanRegistry) Register(ctx context.Context, billableType string) error {
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

				r.manager.AddPlan(billableType, spark.SparkPlan{
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
