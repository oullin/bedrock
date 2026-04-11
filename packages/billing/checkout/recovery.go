package checkout

import (
	"context"

	"github.com/bedrock/packages/contracts"
	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/subscription"
)

// Recovery handles manual recovery of subscriptions by
// attaching an explicit provider subscription ID.
type Recovery struct {
	subscriptions billing.SubscriptionStore
	fetcher       billing.ProviderSubscriptionFetcher
	transitioner  *subscription.Transitioner
	clock         contracts.Clock
}

// NewRecovery creates a Recovery.
func NewRecovery(
	subscriptions billing.SubscriptionStore,
	fetcher billing.ProviderSubscriptionFetcher,
	transitioner *subscription.Transitioner,
	clock contracts.Clock,
) *Recovery {
	return &Recovery{
		subscriptions: subscriptions,
		fetcher:       fetcher,
		transitioner:  transitioner,
		clock:         clock,
	}
}

// Recover attaches a known provider subscription ID to the billable's pending
// local subscription.
func (r *Recovery) Recover(ctx context.Context, billable billing.Billable, providerSubID string) (*billing.Subscription, error) {
	providerSub, err := r.fetcher.FetchExplicit(ctx, providerSubID)

	if err != nil {
		return nil, err
	}

	if providerSub == nil {
		return nil, billing.ErrRecoveryFailed
	}

	localSub, err := r.subscriptions.CurrentForBillable(ctx, billable.BillableType(), billable.BillableID())

	if err != nil || localSub == nil {
		return nil, billing.ErrRecoveryFailed
	}

	localSub.ProviderID = providerSub.ID
	localSub.UpdatedAt = r.clock.Now()

	if err := r.subscriptions.Save(ctx, localSub); err != nil {
		return nil, err
	}

	activated, err := r.transitioner.Activate(ctx, localSub)

	if err != nil {
		return nil, err
	}

	return activated, nil
}
