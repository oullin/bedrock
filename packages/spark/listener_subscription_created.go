package spark

import "context"

// SubscriptionCreatedListener handles new subscriptions from the provider,
// clearing the pending checkout and canceling other active subscriptions.
type SubscriptionCreatedListener struct {
	customers     CustomerStore
	subscriptions SubscriptionStore
	provider      ProviderSubscriptionManager
}

// NewSubscriptionCreatedListener creates the listener.
func NewSubscriptionCreatedListener(
	customers CustomerStore,
	subscriptions SubscriptionStore,
	provider ProviderSubscriptionManager,
) *SubscriptionCreatedListener {
	return &SubscriptionCreatedListener{
		customers:     customers,
		subscriptions: subscriptions,
		provider:      provider,
	}
}

// Handle clears the pending checkout and cancels other active subscriptions.
func (l *SubscriptionCreatedListener) Handle(ctx context.Context, sub *Subscription, billableID int64) error {
	// Clear pending checkout.
	customer, err := l.customers.FindByBillable(ctx, sub.BillableType, billableID)
	if err == nil && customer != nil && customer.PendingCheckout != nil {
		customer.PendingCheckout = nil
		_ = l.customers.Save(ctx, customer)
	}

	// Cancel other active subscriptions.
	active, err := l.subscriptions.ActiveForBillable(ctx, sub.BillableType, billableID)
	if err != nil {
		return err
	}

	for _, other := range active {
		if other.ID == sub.ID {
			continue
		}

		if other.ProviderID != "" {
			_ = l.provider.Cancel(ctx, other.ProviderID, true)
		}
	}

	return nil
}
