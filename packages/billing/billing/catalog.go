package billing

import (
	"context"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/catalog"
)

// CatalogWorkflow is the facade for plan catalog operations.
type CatalogWorkflow struct {
	presentation *catalog.PlanPresentation
	registry     *catalog.BillingPlanRegistry
}

// NewCatalogWorkflow creates a CatalogWorkflow.
func NewCatalogWorkflow(presentation *catalog.PlanPresentation, registry *catalog.BillingPlanRegistry) *CatalogWorkflow {
	return &CatalogWorkflow{presentation: presentation, registry: registry}
}

// LandingPlans returns plans formatted for the landing page.
func (w *CatalogWorkflow) LandingPlans(ctx context.Context) ([]billing.LandingPlanView, error) {
	return w.presentation.LandingPlans(ctx)
}

// RegisterBillingPlans registers database plans with the Billing manager.
func (w *CatalogWorkflow) RegisterBillingPlans(ctx context.Context, billableType string) error {
	return w.registry.Register(ctx, billableType)
}
