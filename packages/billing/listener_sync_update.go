package billing

import "context"

// SyncSubscriptionAfterUpdateListener syncs local subscription data after a
// provider update event.
type SyncSubscriptionAfterUpdateListener struct {
	billing *BillingWorkflow
}

// NewSyncSubscriptionAfterUpdateListener creates the listener.
func NewSyncSubscriptionAfterUpdateListener(billing *BillingWorkflow) *SyncSubscriptionAfterUpdateListener {
	return &SyncSubscriptionAfterUpdateListener{billing: billing}
}

// Handle processes a subscription updated event.
func (l *SyncSubscriptionAfterUpdateListener) Handle(ctx context.Context, sub *Subscription) error {
	return l.billing.SyncUpdated(ctx, sub)
}
