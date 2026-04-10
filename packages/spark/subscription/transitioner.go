package subscription

import (
	"context"

	"github.com/bedrock/packages/spark"
	"github.com/bedrock/packages/spark/entitlement"
)

// Transitioner handles subscription state machine transitions.
type Transitioner struct {
	subscriptions spark.SubscriptionStore
	entitlements  *entitlement.Synchronizer
	events        spark.EventDispatcher
	clock         spark.Clock
}

// NewTransitioner creates a Transitioner.
func NewTransitioner(
	subscriptions spark.SubscriptionStore,
	entitlements *entitlement.Synchronizer,
	events spark.EventDispatcher,
	clock spark.Clock,
) *Transitioner {
	return &Transitioner{
		subscriptions: subscriptions,
		entitlements:  entitlements,
		events:        events,
		clock:         clock,
	}
}

// MarkPaymentReady transitions a pending subscription to awaiting payment.
func (t *Transitioner) MarkPaymentReady(ctx context.Context, sub *spark.Subscription) (*spark.Subscription, error) {
	if sub.Status != spark.StatusPending {
		return sub, nil
	}

	now := t.clock.Now()
	sub.Status = spark.StatusAwaitingPayment
	sub.PaymentReadyAt = &now
	sub.UpdatedAt = now

	if err := t.subscriptions.Save(ctx, sub); err != nil {
		return nil, err
	}

	return sub, nil
}

// Activate transitions a subscription to active.
func (t *Transitioner) Activate(ctx context.Context, sub *spark.Subscription) (*spark.Subscription, error) {
	if sub.Status == spark.StatusActive {
		return sub, nil
	}

	sub.Status = spark.StatusActive
	sub.PendingExpiresAt = nil
	sub.UpdatedAt = t.clock.Now()

	if err := t.subscriptions.Save(ctx, sub); err != nil {
		return nil, err
	}

	_ = t.entitlements.SyncActivationState(ctx, sub)

	return sub, nil
}

// MarkPastDue transitions a subscription to past due.
func (t *Transitioner) MarkPastDue(ctx context.Context, sub *spark.Subscription) (*spark.Subscription, error) {
	if sub.Status == spark.StatusPastDue {
		return sub, nil
	}

	sub.Status = spark.StatusPastDue
	sub.UpdatedAt = t.clock.Now()

	if err := t.subscriptions.Save(ctx, sub); err != nil {
		return nil, err
	}

	_ = t.entitlements.SyncActivationState(ctx, sub)

	return sub, nil
}

// Pause transitions a subscription to paused.
func (t *Transitioner) Pause(ctx context.Context, sub *spark.Subscription) (*spark.Subscription, error) {
	if sub.Status == spark.StatusPaused {
		return sub, nil
	}

	now := t.clock.Now()
	sub.Status = spark.StatusPaused
	sub.PausedAt = &now
	sub.UpdatedAt = now

	if err := t.subscriptions.Save(ctx, sub); err != nil {
		return nil, err
	}

	_ = t.entitlements.SyncActivationState(ctx, sub)

	return sub, nil
}

// Cancel transitions a subscription to canceled.
func (t *Transitioner) Cancel(ctx context.Context, sub *spark.Subscription) (*spark.Subscription, error) {
	if sub.Status == spark.StatusCanceled {
		return sub, nil
	}

	sub.Status = spark.StatusCanceled
	sub.UpdatedAt = t.clock.Now()

	if err := t.subscriptions.Save(ctx, sub); err != nil {
		return nil, err
	}

	_ = t.entitlements.SyncActivationState(ctx, sub)

	return sub, nil
}

// Expire transitions a subscription to expired.
func (t *Transitioner) Expire(ctx context.Context, sub *spark.Subscription) (*spark.Subscription, error) {
	if sub.Status == spark.StatusExpired {
		return sub, nil
	}

	sub.Status = spark.StatusExpired
	sub.UpdatedAt = t.clock.Now()

	if err := t.subscriptions.Save(ctx, sub); err != nil {
		return nil, err
	}

	_ = t.entitlements.SyncActivationState(ctx, sub)

	return sub, nil
}

// ExpireStaleSubscriptions finds and expires subscriptions whose pending
// hold window has elapsed.
func (t *Transitioner) ExpireStaleSubscriptions(ctx context.Context) (int, error) {
	subs, err := t.subscriptions.FindExpirable(ctx, t.clock.Now())

	if err != nil {
		return 0, err
	}

	count := 0

	for _, sub := range subs {
		if _, err := t.Expire(ctx, sub); err == nil {
			count++
		}
	}

	return count, nil
}
