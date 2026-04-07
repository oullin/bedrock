package billing

import (
	"context"

	"github.com/bedrock/packages/billing"
	"github.com/bedrock/packages/billing/catalog"
	"github.com/bedrock/packages/billing/checkout"
	"github.com/bedrock/packages/billing/entitlement"
	"github.com/bedrock/packages/billing/subscription"
)

// Workflow is the primary facade for billing operations. Application
// code calls this rather than individual services directly.
type Workflow struct {
	catalog      *catalog.PlanCatalog
	stateRes     *subscription.StateResolver
	subs         *subscription.Repository
	lifecycle    *subscription.Transitioner
	starter      *checkout.Starter
	reconciler   *checkout.Reconciler
	syncer       *checkout.Syncer
	entitlements *entitlement.Synchronizer
}

// NewWorkflow creates a Workflow.
func NewWorkflow(
	cat *catalog.PlanCatalog,
	stateRes *subscription.StateResolver,
	subs *subscription.Repository,
	lifecycle *subscription.Transitioner,
	starter *checkout.Starter,
	reconciler *checkout.Reconciler,
	syncer *checkout.Syncer,
	entitlements *entitlement.Synchronizer,
) *Workflow {
	return &Workflow{
		catalog:      cat,
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
func (w *Workflow) ReadBillingState(ctx context.Context, billableType string, billableID int64) billing.BillingStateSnapshot {
	return w.stateRes.Read(ctx, billableType, billableID)
}

// StartPendingSubscription creates a pending subscription.
func (w *Workflow) StartPendingSubscription(ctx context.Context, billable billing.Billable, plan billing.SubscriptionPlan, period billing.BillingPeriod) (*billing.Subscription, error) {
	return w.subs.StartPendingSubscription(ctx, billable, plan, period)
}

// StartStarterTrial creates a starter trial subscription.
func (w *Workflow) StartStarterTrial(ctx context.Context, billable billing.Billable) (*billing.Subscription, error) {
	return w.subs.StartStarterTrial(ctx, billable)
}

// StartCheckout initiates a checkout for the plan and period.
func (w *Workflow) StartCheckout(ctx context.Context, billable billing.Billable, plan billing.SubscriptionPlan, period billing.BillingPeriod) (string, error) {
	return w.starter.Start(ctx, billable, plan, period)
}

// MarkPaymentReady transitions a subscription to awaiting payment.
func (w *Workflow) MarkPaymentReady(ctx context.Context, sub *billing.Subscription) (*billing.Subscription, error) {
	return w.lifecycle.MarkPaymentReady(ctx, sub)
}

// Activate transitions a subscription to active.
func (w *Workflow) Activate(ctx context.Context, sub *billing.Subscription) (*billing.Subscription, error) {
	return w.lifecycle.Activate(ctx, sub)
}

// MarkPastDue transitions a subscription to past due.
func (w *Workflow) MarkPastDue(ctx context.Context, sub *billing.Subscription) (*billing.Subscription, error) {
	return w.lifecycle.MarkPastDue(ctx, sub)
}

// Pause transitions a subscription to paused.
func (w *Workflow) Pause(ctx context.Context, sub *billing.Subscription) (*billing.Subscription, error) {
	return w.lifecycle.Pause(ctx, sub)
}

// Cancel transitions a subscription to canceled.
func (w *Workflow) Cancel(ctx context.Context, sub *billing.Subscription) (*billing.Subscription, error) {
	return w.lifecycle.Cancel(ctx, sub)
}

// Expire transitions a subscription to expired.
func (w *Workflow) Expire(ctx context.Context, sub *billing.Subscription) (*billing.Subscription, error) {
	return w.lifecycle.Expire(ctx, sub)
}

// ExpireStaleSubscriptions expires all stale pending subscriptions.
func (w *Workflow) ExpireStaleSubscriptions(ctx context.Context) (int, error) {
	return w.lifecycle.ExpireStaleSubscriptions(ctx)
}

// CurrentSubscription returns the current subscription for the billable.
func (w *Workflow) CurrentSubscription(ctx context.Context, billableType string, billableID int64) (*billing.Subscription, error) {
	return w.subs.CurrentSubscription(ctx, billableType, billableID)
}

// ReconcileCreated handles post-checkout reconciliation.
func (w *Workflow) ReconcileCreated(ctx context.Context, sub *billing.Subscription, billableID int64) error {
	return w.reconciler.ReconcileCreated(ctx, sub, billableID)
}

// SyncUpdated syncs a subscription after a provider update event.
func (w *Workflow) SyncUpdated(ctx context.Context, sub *billing.Subscription) error {
	return w.syncer.SyncUpdated(ctx, sub)
}

// ReconcilePendingSubscriptions batch-reconciles pending subscriptions.
func (w *Workflow) ReconcilePendingSubscriptions(ctx context.Context) error {
	return w.reconciler.ReconcilePending(ctx)
}

// SyncEntitlements performs a full sync of subscription features.
func (w *Workflow) SyncEntitlements(ctx context.Context, sub *billing.Subscription) error {
	return w.entitlements.Sync(ctx, sub)
}

// PendingExpiryDays returns the hold period for pending subscriptions.
func (w *Workflow) PendingExpiryDays() int {
	return w.subs.PendingExpiryDays()
}
