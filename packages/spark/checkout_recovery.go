package spark

import "context"

// ExplicitSubscriptionRecovery handles manual recovery of subscriptions by
// attaching an explicit provider subscription ID.
type ExplicitSubscriptionRecovery struct {
	subscriptions SubscriptionStore
	fetcher       ProviderSubscriptionFetcher
	transitioner  *SubscriptionTransitioner
	clock         Clock
}

// NewExplicitSubscriptionRecovery creates an ExplicitSubscriptionRecovery.
func NewExplicitSubscriptionRecovery(
	subscriptions SubscriptionStore,
	fetcher ProviderSubscriptionFetcher,
	transitioner *SubscriptionTransitioner,
	clock Clock,
) *ExplicitSubscriptionRecovery {
	return &ExplicitSubscriptionRecovery{
		subscriptions: subscriptions,
		fetcher:       fetcher,
		transitioner:  transitioner,
		clock:         clock,
	}
}

// Recover attaches a known provider subscription ID to the billable's pending
// local subscription.
func (r *ExplicitSubscriptionRecovery) Recover(ctx context.Context, billable Billable, providerSubID string) (*Subscription, error) {
	providerSub, err := r.fetcher.FetchExplicit(ctx, providerSubID)
	if err != nil {
		return nil, err
	}

	if providerSub == nil {
		return nil, ErrRecoveryFailed
	}

	localSub, err := r.subscriptions.CurrentForBillable(ctx, billable.BillableType(), billable.BillableID())
	if err != nil || localSub == nil {
		return nil, ErrRecoveryFailed
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
