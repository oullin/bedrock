package spark

import "context"

// SyncSubscriptionEntitlementsListener syncs entitlements after a subscription
// state transition.
type SyncSubscriptionEntitlementsListener struct {
	billing *BillingWorkflow
	events  EventDispatcher
}

// NewSyncSubscriptionEntitlementsListener creates the listener.
func NewSyncSubscriptionEntitlementsListener(billing *BillingWorkflow, events EventDispatcher) *SyncSubscriptionEntitlementsListener {
	return &SyncSubscriptionEntitlementsListener{billing: billing, events: events}
}

// Handle processes a subscription transition and syncs entitlements.
func (l *SyncSubscriptionEntitlementsListener) Handle(ctx context.Context, sub *Subscription) error {
	if err := l.billing.SyncEntitlements(ctx, sub); err != nil {
		return err
	}

	return l.events.Dispatch(ctx, SubscriptionChangedEvent{
		TeamID: sub.BillableID,
		Plan:   sub.Plan,
		Status: string(sub.Status),
	})
}
