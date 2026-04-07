package spark

import "context"

// ReconcileAfterCheckoutListener handles post-checkout reconciliation when
// a new subscription is created on the provider.
type ReconcileAfterCheckoutListener struct {
	billing *BillingWorkflow
}

// NewReconcileAfterCheckoutListener creates the listener.
func NewReconcileAfterCheckoutListener(billing *BillingWorkflow) *ReconcileAfterCheckoutListener {
	return &ReconcileAfterCheckoutListener{billing: billing}
}

// Handle processes a subscription created event.
func (l *ReconcileAfterCheckoutListener) Handle(ctx context.Context, sub *Subscription, billableID int64) error {
	return l.billing.ReconcileCreated(ctx, sub, billableID)
}
