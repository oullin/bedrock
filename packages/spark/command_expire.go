package spark

import "context"

// ExpirePendingSubscriptionsCommand expires local subscriptions whose
// time-limited access window has elapsed.
type ExpirePendingSubscriptionsCommand struct {
	billing *BillingWorkflow
}

// NewExpirePendingSubscriptionsCommand creates the command.
func NewExpirePendingSubscriptionsCommand(billing *BillingWorkflow) *ExpirePendingSubscriptionsCommand {
	return &ExpirePendingSubscriptionsCommand{billing: billing}
}

// Name returns the command name.
func (c *ExpirePendingSubscriptionsCommand) Name() string {
	return "subscriptions:expire-pending"
}

// Description returns a human-readable description.
func (c *ExpirePendingSubscriptionsCommand) Description() string {
	return "Expire and remove local subscriptions whose time-limited access window has elapsed."
}

// Handle executes the command.
func (c *ExpirePendingSubscriptionsCommand) Handle(ctx context.Context, log Logger) error {
	expired, err := c.billing.ExpireStaleSubscriptions(ctx)
	if err != nil {
		return err
	}

	log.Info("Expired %d local subscriptions.", expired)

	return nil
}
