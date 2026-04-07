package billing

import "context"

// SubscriptionReconciler matches provider-created subscriptions with local
// pending records.
type SubscriptionReconciler struct {
	subscriptions SubscriptionStore
	catalog       *PlanCatalog
	fetcher       ProviderSubscriptionFetcher
	transitioner  *SubscriptionTransitioner
	clock         Clock
}

// NewSubscriptionReconciler creates a SubscriptionReconciler.
func NewSubscriptionReconciler(
	subscriptions SubscriptionStore,
	catalog *PlanCatalog,
	fetcher ProviderSubscriptionFetcher,
	transitioner *SubscriptionTransitioner,
	clock Clock,
) *SubscriptionReconciler {
	return &SubscriptionReconciler{
		subscriptions: subscriptions,
		catalog:       catalog,
		fetcher:       fetcher,
		transitioner:  transitioner,
		clock:         clock,
	}
}

// ReconcileCreated handles post-checkout reconciliation. It matches a
// provider subscription with the local pending record.
func (r *SubscriptionReconciler) ReconcileCreated(ctx context.Context, providerSub *Subscription, billableID int64) error {
	if providerSub == nil || providerSub.ProviderID == "" {
		return nil
	}

	localSub, err := r.subscriptions.CurrentForBillable(ctx, providerSub.BillableType, billableID)
	if err != nil || localSub == nil {
		return err
	}

	if localSub.Status != StatusPending && localSub.Status != StatusAwaitingPayment {
		return nil
	}

	localSub.ProviderID = providerSub.ProviderID
	localSub.UpdatedAt = r.clock.Now()

	if err := r.subscriptions.Save(ctx, localSub); err != nil {
		return err
	}

	// Resolve price from provider items.
	if len(providerSub.Items) > 0 {
		price, _ := r.catalog.PriceForProviderID(ctx, providerSub.Items[0].PriceID)
		if price != nil {
			localSub.PlanPeriodPriceID = &price.ID
			_ = r.subscriptions.Save(ctx, localSub)
		}
	}

	_, err = r.transitioner.Activate(ctx, localSub)

	return err
}

// ReconcilePending batch-reconciles all pending subscriptions without a
// provider ID by querying the payment provider.
func (r *SubscriptionReconciler) ReconcilePending(ctx context.Context) error {
	pending, err := r.subscriptions.FindPendingWithoutProvider(ctx)
	if err != nil {
		return err
	}

	for _, sub := range pending {
		providerSub, err := r.fetcher.FetchForBillable(ctx, &billableRef{
			id:    sub.BillableID,
			btype: sub.BillableType,
		})
		if err != nil || providerSub == nil {
			continue
		}

		sub.ProviderID = providerSub.ID
		sub.UpdatedAt = r.clock.Now()

		if err := r.subscriptions.Save(ctx, sub); err != nil {
			continue
		}

		_, _ = r.transitioner.Activate(ctx, sub)
	}

	return nil
}

// billableRef is a minimal Billable implementation for reconciliation.
type billableRef struct {
	id    int64
	btype string
}

func (b *billableRef) BillableID() int64     { return b.id }
func (b *billableRef) BillableType() string  { return b.btype }
func (b *billableRef) BillableUUID() string  { return "" }
func (b *billableRef) BillableName() string  { return "" }
func (b *billableRef) BillableEmail() string { return "" }
