package spark

import "context"

// ReconcileSubscriptionsCommand reconciles local subscription records against
// the payment provider.
type ReconcileSubscriptionsCommand struct {
	billing *BillingWorkflow
}

// NewReconcileSubscriptionsCommand creates the command.
func NewReconcileSubscriptionsCommand(billing *BillingWorkflow) *ReconcileSubscriptionsCommand {
	return &ReconcileSubscriptionsCommand{billing: billing}
}

// Name returns the command name.
func (c *ReconcileSubscriptionsCommand) Name() string {
	return "subscriptions:reconcile"
}

// Description returns a human-readable description.
func (c *ReconcileSubscriptionsCommand) Description() string {
	return "Reconcile local subscription records against the payment provider."
}

// Handle executes the command.
func (c *ReconcileSubscriptionsCommand) Handle(ctx context.Context, log Logger) error {
	log.Warn("This command mutates local billing state.")

	if err := c.billing.ReconcilePendingSubscriptions(ctx); err != nil {
		return err
	}

	log.Info("Subscription reconciliation completed.")

	return nil
}
