package billing

import "context"

// BillingWorkflow is the primary facade for billing operations. Application
// code calls this rather than individual services directly.
type BillingWorkflow struct {
	catalog      *PlanCatalog
	stateRes     *BillingStateResolver
	subs         *SubscriptionRepository
	lifecycle    *SubscriptionTransitioner
	starter      *CheckoutStarter
	reconciler   *SubscriptionReconciler
	syncer       *SubscriptionUpdateSyncer
	entitlements *EntitlementSynchronizer
}

// NewBillingWorkflow creates a BillingWorkflow.
func NewBillingWorkflow(
	catalog *PlanCatalog,
	stateRes *BillingStateResolver,
	subs *SubscriptionRepository,
	lifecycle *SubscriptionTransitioner,
	starter *CheckoutStarter,
	reconciler *SubscriptionReconciler,
	syncer *SubscriptionUpdateSyncer,
	entitlements *EntitlementSynchronizer,
) *BillingWorkflow {
	return &BillingWorkflow{
		catalog:      catalog,
		stateRes:     stateRes,
		subs:         subs,
		lifecycle:    lifecycle,
		starter:      starter,
		reconciler:   reconciler,
		syncer:       syncer,
		entitlements: entitlements,
	}
}

// ReadBillingState returns the current billing state for the billable.
func (w *BillingWorkflow) ReadBillingState(ctx context.Context, billableType string, billableID int64) BillingStateSnapshot {
	return w.stateRes.Read(ctx, billableType, billableID)
}

// StartPendingSubscription creates a pending subscription.
func (w *BillingWorkflow) StartPendingSubscription(ctx context.Context, billable Billable, plan SubscriptionPlan, period BillingPeriod) (*Subscription, error) {
	return w.subs.StartPendingSubscription(ctx, billable, plan, period)
}

// StartStarterTrial creates a starter trial subscription.
func (w *BillingWorkflow) StartStarterTrial(ctx context.Context, billable Billable) (*Subscription, error) {
	return w.subs.StartStarterTrial(ctx, billable)
}

// StartCheckout initiates a checkout for the plan and period.
func (w *BillingWorkflow) StartCheckout(ctx context.Context, billable Billable, plan SubscriptionPlan, period BillingPeriod) (string, error) {
	return w.starter.Start(ctx, billable, plan, period)
}

// MarkPaymentReady transitions a subscription to awaiting payment.
func (w *BillingWorkflow) MarkPaymentReady(ctx context.Context, sub *Subscription) (*Subscription, error) {
	return w.lifecycle.MarkPaymentReady(ctx, sub)
}

// Activate transitions a subscription to active.
func (w *BillingWorkflow) Activate(ctx context.Context, sub *Subscription) (*Subscription, error) {
	return w.lifecycle.Activate(ctx, sub)
}

// MarkPastDue transitions a subscription to past due.
func (w *BillingWorkflow) MarkPastDue(ctx context.Context, sub *Subscription) (*Subscription, error) {
	return w.lifecycle.MarkPastDue(ctx, sub)
}

// Pause transitions a subscription to paused.
func (w *BillingWorkflow) Pause(ctx context.Context, sub *Subscription) (*Subscription, error) {
	return w.lifecycle.Pause(ctx, sub)
}

// Cancel transitions a subscription to canceled.
func (w *BillingWorkflow) Cancel(ctx context.Context, sub *Subscription) (*Subscription, error) {
	return w.lifecycle.Cancel(ctx, sub)
}

// Expire transitions a subscription to expired.
func (w *BillingWorkflow) Expire(ctx context.Context, sub *Subscription) (*Subscription, error) {
	return w.lifecycle.Expire(ctx, sub)
}

// ExpireStaleSubscriptions expires all stale pending subscriptions.
func (w *BillingWorkflow) ExpireStaleSubscriptions(ctx context.Context) (int, error) {
	return w.lifecycle.ExpireStaleSubscriptions(ctx)
}

// CurrentSubscription returns the current subscription for the billable.
func (w *BillingWorkflow) CurrentSubscription(ctx context.Context, billableType string, billableID int64) (*Subscription, error) {
	return w.subs.CurrentSubscription(ctx, billableType, billableID)
}

// ReconcileCreated handles post-checkout reconciliation.
func (w *BillingWorkflow) ReconcileCreated(ctx context.Context, sub *Subscription, billableID int64) error {
	return w.reconciler.ReconcileCreated(ctx, sub, billableID)
}

// SyncUpdated syncs a subscription after a provider update event.
func (w *BillingWorkflow) SyncUpdated(ctx context.Context, sub *Subscription) error {
	return w.syncer.SyncUpdated(ctx, sub)
}

// ReconcilePendingSubscriptions batch-reconciles pending subscriptions.
func (w *BillingWorkflow) ReconcilePendingSubscriptions(ctx context.Context) error {
	return w.reconciler.ReconcilePending(ctx)
}

// SyncEntitlements performs a full sync of subscription features.
func (w *BillingWorkflow) SyncEntitlements(ctx context.Context, sub *Subscription) error {
	return w.entitlements.Sync(ctx, sub)
}

// PendingExpiryDays returns the hold period for pending subscriptions.
func (w *BillingWorkflow) PendingExpiryDays() int {
	return w.subs.PendingExpiryDays()
}
