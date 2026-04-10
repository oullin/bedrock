package command

import (
	"context"

	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/spark/billing"
	"github.com/bedrock/packages/spark/checkout"
)

// ReconcileSubscriptions reconciles local subscription records against
// the payment provider.
type ReconcileSubscriptions struct {
	billing *billing.Workflow
}

// NewReconcileSubscriptions creates the command.

// Name returns the command name.

// Description returns a human-readable description.

// Handle executes the command.

// ExpirePendingSubscriptions expires local subscriptions whose
// time-limited access window has elapsed.
type ExpirePendingSubscriptions struct {
	billing *billing.Workflow
}

// NewExpirePendingSubscriptions creates the command.

// Name returns the command name.

// Description returns a human-readable description.

// Handle executes the command.

// RecoverExplicitSubscription recovers a pending local subscription
// by attaching an explicit remote provider subscription.
type RecoverExplicitSubscription struct {
	recovery *checkout.Recovery
}

func NewReconcileSubscriptions(b *billing.Workflow) *ReconcileSubscriptions {
	return &ReconcileSubscriptions{billing: b}
}

func (c *ReconcileSubscriptions) Name() string {
	return "subscriptions:reconcile"
}

func (c *ReconcileSubscriptions) Description() string {
	return "Reconcile local subscription records against the payment provider."
}

func (c *ReconcileSubscriptions) Handle(ctx context.Context, log spark.Logger) error {
	log.Warn("This command mutates local billing state.")

	if err := c.billing.ReconcilePendingSubscriptions(ctx); err != nil {
		return err
	}

	log.Info("Subscription reconciliation completed.")

	return nil
}

func NewExpirePendingSubscriptions(b *billing.Workflow) *ExpirePendingSubscriptions {
	return &ExpirePendingSubscriptions{billing: b}
}

func (c *ExpirePendingSubscriptions) Name() string {
	return "subscriptions:expire-pending"
}

func (c *ExpirePendingSubscriptions) Description() string {
	return "Expire and remove local subscriptions whose time-limited access window has elapsed."
}

func (c *ExpirePendingSubscriptions) Handle(ctx context.Context, log spark.Logger) error {
	expired, err := c.billing.ExpireStaleSubscriptions(ctx)

	if err != nil {
		return err
	}

	log.Info("Expired %d local subscriptions.", expired)

	return nil
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
func (c *RecoverExplicitSubscription) Handle(ctx context.Context, log spark.Logger, billable spark.Billable, providerSubID string) error {
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
