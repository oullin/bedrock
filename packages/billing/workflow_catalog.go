package billing

import "context"

// CatalogWorkflow is the facade for plan catalog operations.
type CatalogWorkflow struct {
	presentation *PlanPresentation
	registry     *BillingPlanRegistry
}

// NewCatalogWorkflow creates a CatalogWorkflow.
func NewCatalogWorkflow(presentation *PlanPresentation, registry *BillingPlanRegistry) *CatalogWorkflow {
	return &CatalogWorkflow{presentation: presentation, registry: registry}
}

// LandingPlans returns plans formatted for the landing page.
func (w *CatalogWorkflow) LandingPlans(ctx context.Context) ([]LandingPlanView, error) {
	return w.presentation.LandingPlans(ctx)
}

// RegisterBillingPlans registers database plans with the Billing manager.
func (w *CatalogWorkflow) RegisterBillingPlans(ctx context.Context, billableType string) error {
	return w.registry.Register(ctx, billableType)
}
