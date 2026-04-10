package checkout

import (
	"context"

	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/spark/subscription"
)

// Recovery handles manual recovery of subscriptions by
// attaching an explicit provider subscription ID.
type Recovery struct {
	subscriptions spark.SubscriptionStore
	fetcher       spark.ProviderSubscriptionFetcher
	transitioner  *subscription.Transitioner
	clock         spark.Clock
}

// NewRecovery creates a Recovery.
func NewRecovery(
	subscriptions spark.SubscriptionStore,
	fetcher spark.ProviderSubscriptionFetcher,
	transitioner *subscription.Transitioner,
	clock spark.Clock,
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
func (r *Recovery) Recover(ctx context.Context, billable spark.Billable, providerSubID string) (*spark.Subscription, error) {
	providerSub, err := r.fetcher.FetchExplicit(ctx, providerSubID)

	if err != nil {
		return nil, err
	}

	if providerSub == nil {
		return nil, spark.ErrRecoveryFailed
	}

	localSub, err := r.subscriptions.CurrentForBillable(ctx, billable.BillableType(), billable.BillableID())

	if err != nil || localSub == nil {
		return nil, spark.ErrRecoveryFailed
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
