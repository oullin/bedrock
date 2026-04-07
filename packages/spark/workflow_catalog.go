package spark

import "context"

// CatalogWorkflow is the facade for plan catalog operations.
type CatalogWorkflow struct {
	presentation *PlanPresentation
	registry     *SparkPlanRegistry
}

// NewCatalogWorkflow creates a CatalogWorkflow.
func NewCatalogWorkflow(presentation *PlanPresentation, registry *SparkPlanRegistry) *CatalogWorkflow {
	return &CatalogWorkflow{presentation: presentation, registry: registry}
}

// LandingPlans returns plans formatted for the landing page.
func (w *CatalogWorkflow) LandingPlans(ctx context.Context) ([]LandingPlanView, error) {
	return w.presentation.LandingPlans(ctx)
}

// RegisterSparkPlans registers database plans with the Spark manager.
func (w *CatalogWorkflow) RegisterSparkPlans(ctx context.Context, billableType string) error {
	return w.registry.Register(ctx, billableType)
}
