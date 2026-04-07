package command

import (
	"context"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/billing"
	"github.com/bedrock/packages/billing/checkout"
)

// ReconcileSubscriptions reconciles local subscription records against
// the payment provider.
type ReconcileSubscriptions struct {
	billing *billing.Workflow
}

// NewReconcileSubscriptions creates the command.
func NewReconcileSubscriptions(b *billing.Workflow) *ReconcileSubscriptions {
	return &ReconcileSubscriptions{billing: b}
}

// Name returns the command name.
func (c *ReconcileSubscriptions) Name() string {
	return "subscriptions:reconcile"
}

// Description returns a human-readable description.
func (c *ReconcileSubscriptions) Description() string {
	return "Reconcile local subscription records against the payment provider."
}

// Handle executes the command.
func (c *ReconcileSubscriptions) Handle(ctx context.Context, log billing.Logger) error {
	log.Warn("This command mutates local billing state.")

	if err := c.billing.ReconcilePendingSubscriptions(ctx); err != nil {
		return err
	}

	log.Info("Subscription reconciliation completed.")

	return nil
}

// ExpirePendingSubscriptions expires local subscriptions whose
// time-limited access window has elapsed.
type ExpirePendingSubscriptions struct {
	billing *billing.Workflow
}

// NewExpirePendingSubscriptions creates the command.
func NewExpirePendingSubscriptions(b *billing.Workflow) *ExpirePendingSubscriptions {
	return &ExpirePendingSubscriptions{billing: b}
}

// Name returns the command name.
func (c *ExpirePendingSubscriptions) Name() string {
	return "subscriptions:expire-pending"
}

// Description returns a human-readable description.
func (c *ExpirePendingSubscriptions) Description() string {
	return "Expire and remove local subscriptions whose time-limited access window has elapsed."
}

// Handle executes the command.
func (c *ExpirePendingSubscriptions) Handle(ctx context.Context, log billing.Logger) error {
	expired, err := c.billing.ExpireStaleSubscriptions(ctx)
	if err != nil {
		return err
	}

	log.Info("Expired %d local subscriptions.", expired)

	return nil
}

// RecoverExplicitSubscription recovers a pending local subscription
// by attaching an explicit remote provider subscription.
type RecoverExplicitSubscription struct {
	recovery *checkout.Recovery
}

// NewRecoverExplicitSubscription creates the command.
func NewRecoverExplicitSubscription(recovery *checkout.Recovery) *RecoverExplicitSubscription {
	return &RecoverExplicitSubscription{recovery: recovery}
}

// Name returns the command name.
func (c *RecoverExplicitSubscription) Name() string {
	return "subscriptions:recover-explicit"
}

// Description returns a human-readable description.
func (c *RecoverExplicitSubscription) Description() string {
	return "Recover a pending local subscription by attaching an explicit remote provider subscription."
}

// Handle executes the command. The billable and provider subscription ID must
// be provided by the caller.
func (c *RecoverExplicitSubscription) Handle(ctx context.Context, log billing.Logger, billable billing.Billable, providerSubID string) error {
	log.Warn("This command mutates local billing state for the specified billable.")

	sub, err := c.recovery.Recover(ctx, billable, providerSubID)
	if err != nil {
		log.Error("Recovery failed: %s", err.Error())

		return err
	}

	log.Info("Recovered subscription %s for billable %s (%s).",
		sub.ProviderID, billable.BillableName(), billable.BillableUUID())

	return nil
}
