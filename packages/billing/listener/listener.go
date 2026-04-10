package listener

import (
	"context"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/billing"
)

// ReconcileAfterCheckout handles post-checkout reconciliation when
// a new subscription is created on the provider.
type ReconcileAfterCheckout struct {
	billing *billing.Workflow
}

// NewReconcileAfterCheckout creates the listener.

// Handle processes a subscription created event.

// SubscriptionCreated handles new subscriptions from the provider,
// clearing the pending checkout and canceling other active subscriptions.
type SubscriptionCreated struct {
	customers     billing.CustomerStore
	subscriptions billing.SubscriptionStore
	provider      billing.ProviderSubscriptionManager
}

// NewSubscriptionCreated creates the listener.

// Handle clears the pending checkout and cancels other active subscriptions.

// Clear pending checkout.

// Cancel other active subscriptions.

// SyncEntitlements syncs entitlements after a subscription state transition.
type SyncEntitlements struct {
	billing *billing.Workflow
	events  billing.EventDispatcher
}

// NewSyncEntitlements creates the listener.

// Handle processes a subscription transition and syncs entitlements.

// SyncAfterUpdate syncs local subscription data after a provider update event.
type SyncAfterUpdate struct {
	billing *billing.Workflow
}

func NewReconcileAfterCheckout(b *billing.Workflow) *ReconcileAfterCheckout {
	return &ReconcileAfterCheckout{billing: b}
}

func (l *ReconcileAfterCheckout) Handle(ctx context.Context, sub *billing.Subscription, billableID int64) error {
	return l.billing.ReconcileCreated(ctx, sub, billableID)
}

func NewSubscriptionCreated(
	customers billing.CustomerStore,
	subscriptions billing.SubscriptionStore,
	provider billing.ProviderSubscriptionManager,
) *SubscriptionCreated {
	return &SubscriptionCreated{
		customers:     customers,
		subscriptions: subscriptions,
		provider:      provider,
	}
}

func (l *SubscriptionCreated) Handle(ctx context.Context, sub *billing.Subscription, billableID int64) error {

	customer, err := l.customers.FindByBillable(ctx, sub.BillableType, billableID)

	if err == nil && customer != nil && customer.PendingCheckout != nil {
		customer.PendingCheckout = nil
		_ = l.customers.Save(ctx, customer)
	}

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

func NewSyncEntitlements(b *billing.Workflow, events billing.EventDispatcher) *SyncEntitlements {
	return &SyncEntitlements{billing: b, events: events}
}

func (l *SyncEntitlements) Handle(ctx context.Context, sub *billing.Subscription) error {
	if err := l.billing.SyncEntitlements(ctx, sub); err != nil {
		return err
	}

	return l.events.Dispatch(ctx, billing.SubscriptionChangedEvent{
		TeamID: sub.BillableID,
		Plan:   sub.Plan,
		Status: string(sub.Status),
	})
}

// NewSyncAfterUpdate creates the listener.
func NewSyncAfterUpdate(b *billing.Workflow) *SyncAfterUpdate {
	return &SyncAfterUpdate{billing: b}
}

// Handle processes a subscription updated event.
func (l *SyncAfterUpdate) Handle(ctx context.Context, sub *billing.Subscription) error {
	return l.billing.SyncUpdated(ctx, sub)
}
