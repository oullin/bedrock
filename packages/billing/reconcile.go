package billing

import "context"

// CheckoutReconciliationStore exposes the subscription operations needed to
// attach a real Paddle subscription to a pre-checkout local subscription.
type CheckoutReconciliationStore interface {
	ActiveForBillable(ctx context.Context, billableType string, billableID int64) ([]*Subscription, error)
	Save(ctx context.Context, subscription *Subscription) error
	Delete(ctx context.Context, id int64) error
}

// ReconcileSubscriptionAfterCheckout copies local pre-checkout metadata onto
// the Paddle-created subscription, removes the pre-checkout row, and dispatches
// an update event.
func ReconcileSubscriptionAfterCheckout(
	ctx context.Context,
	store CheckoutReconciliationStore,
	billableType string,
	billableID int64,
	cashier *Subscription,
	events EventDispatcher,
) (bool, error) {
	if cashier == nil || cashier.PaddleID == "" {
		return false, nil
	}

	subscriptions, err := store.ActiveForBillable(ctx, billableType, billableID)
	if err != nil {
		return false, err
	}

	pending := matchingPrePaddleSubscription(subscriptions, cashier)
	if pending == nil {
		return false, nil
	}

	if cashier.Type == "" {
		cashier.Type = pending.Type
	}
	if cashier.Plan == "" {
		cashier.Plan = pending.Plan
	}

	if err := store.Save(ctx, cashier); err != nil {
		return false, err
	}

	if err := store.Delete(ctx, pending.ID); err != nil {
		return false, err
	}

	if events != nil {
		_ = events.Dispatch(SubscriptionUpdatedEvent{Subscription: cashier})
	}

	return true, nil
}

func matchingPrePaddleSubscription(subscriptions []*Subscription, cashier *Subscription) *Subscription {
	for _, subscription := range subscriptions {
		if subscription == nil || subscription.ID == cashier.ID || subscription.PaddleID != "" {
			continue
		}

		switch subscription.Status {
		case StatusPending, StatusAwaitingPayment:
		default:
			continue
		}

		if subscriptionMatchesCheckout(subscription, cashier) {
			return subscription
		}
	}

	return nil
}

func subscriptionMatchesCheckout(pending, cashier *Subscription) bool {
	for _, item := range pending.Items {
		if item.PriceID == "" {
			continue
		}

		if cashier.HasPrice(item.PriceID) {
			return true
		}
	}

	return pending.Plan != "" && pending.Plan == cashier.Plan
}
