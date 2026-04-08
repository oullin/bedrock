package listener

import (
	"context"

	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/spark/billing"
)

// ReconcileAfterCheckout handles post-checkout reconciliation when
// a new subscription is created on the provider.
type ReconcileAfterCheckout struct {
	billing *billing.Workflow
}

// NewReconcileAfterCheckout creates the listener.
func NewReconcileAfterCheckout(b *billing.Workflow) *ReconcileAfterCheckout {
	return &ReconcileAfterCheckout{billing: b}
}

// Handle processes a subscription created event.
func (l *ReconcileAfterCheckout) Handle(ctx context.Context, sub *spark.Subscription, billableID int64) error {
	return l.billing.ReconcileCreated(ctx, sub, billableID)
}

// SubscriptionCreated handles new subscriptions from the provider,
// clearing the pending checkout and canceling other active subscriptions.
type SubscriptionCreated struct {
	customers     spark.CustomerStore
	subscriptions spark.SubscriptionStore
	provider      spark.ProviderSubscriptionManager
}

// NewSubscriptionCreated creates the listener.
func NewSubscriptionCreated(
	customers spark.CustomerStore,
	subscriptions spark.SubscriptionStore,
	provider spark.ProviderSubscriptionManager,
) *SubscriptionCreated {
	return &SubscriptionCreated{
		customers:     customers,
		subscriptions: subscriptions,
		provider:      provider,
	}
}

// Handle clears the pending checkout and cancels other active subscriptions.
func (l *SubscriptionCreated) Handle(ctx context.Context, sub *spark.Subscription, billableID int64) error {
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

// SyncEntitlements syncs entitlements after a subscription state transition.
type SyncEntitlements struct {
	billing *billing.Workflow
	events  spark.EventDispatcher
}

// NewSyncEntitlements creates the listener.
func NewSyncEntitlements(b *billing.Workflow, events spark.EventDispatcher) *SyncEntitlements {
	return &SyncEntitlements{billing: b, events: events}
}

// Handle processes a subscription transition and syncs entitlements.
func (l *SyncEntitlements) Handle(ctx context.Context, sub *spark.Subscription) error {
	if err := l.billing.SyncEntitlements(ctx, sub); err != nil {
		return err
	}

	return l.events.Dispatch(ctx, spark.SubscriptionChangedEvent{
		TeamID: sub.BillableID,
		Plan:   sub.Plan,
		Status: string(sub.Status),
	})
}

// SyncAfterUpdate syncs local subscription data after a provider update event.
type SyncAfterUpdate struct {
	billing *billing.Workflow
}

// NewSyncAfterUpdate creates the listener.
func NewSyncAfterUpdate(b *billing.Workflow) *SyncAfterUpdate {
	return &SyncAfterUpdate{billing: b}
}

// Handle processes a subscription updated event.
func (l *SyncAfterUpdate) Handle(ctx context.Context, sub *spark.Subscription) error {
	return l.billing.SyncUpdated(ctx, sub)
}
