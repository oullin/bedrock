package spark

import "context"

// RecoverExplicitSubscriptionCommand recovers a pending local subscription
// by attaching an explicit remote provider subscription.
type RecoverExplicitSubscriptionCommand struct {
	recovery *ExplicitSubscriptionRecovery
}

// NewRecoverExplicitSubscriptionCommand creates the command.
func NewRecoverExplicitSubscriptionCommand(recovery *ExplicitSubscriptionRecovery) *RecoverExplicitSubscriptionCommand {
	return &RecoverExplicitSubscriptionCommand{recovery: recovery}
}

// Name returns the command name.
func (c *RecoverExplicitSubscriptionCommand) Name() string {
	return "subscriptions:recover-explicit"
}

// Description returns a human-readable description.
func (c *RecoverExplicitSubscriptionCommand) Description() string {
	return "Recover a pending local subscription by attaching an explicit remote provider subscription."
}

// Handle executes the command. The billable and provider subscription ID must
// be provided by the caller.
func (c *RecoverExplicitSubscriptionCommand) Handle(ctx context.Context, log Logger, billable Billable, providerSubID string) error {
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
