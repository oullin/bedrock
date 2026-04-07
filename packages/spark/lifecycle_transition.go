package spark

import "context"

// SubscriptionTransitioner handles subscription state machine transitions.
type SubscriptionTransitioner struct {
	subscriptions SubscriptionStore
	entitlements  *EntitlementSynchronizer
	events        EventDispatcher
	clock         Clock
}

// NewSubscriptionTransitioner creates a SubscriptionTransitioner.
func NewSubscriptionTransitioner(
	subscriptions SubscriptionStore,
	entitlements *EntitlementSynchronizer,
	events EventDispatcher,
	clock Clock,
) *SubscriptionTransitioner {
	return &SubscriptionTransitioner{
		subscriptions: subscriptions,
		entitlements:  entitlements,
		events:        events,
		clock:         clock,
	}
}

// MarkPaymentReady transitions a pending subscription to awaiting payment.
func (t *SubscriptionTransitioner) MarkPaymentReady(ctx context.Context, sub *Subscription) (*Subscription, error) {
	if sub.Status != StatusPending {
		return sub, nil
	}

	now := t.clock.Now()
	sub.Status = StatusAwaitingPayment
	sub.PaymentReadyAt = &now
	sub.UpdatedAt = now

	if err := t.subscriptions.Save(ctx, sub); err != nil {
		return nil, err
	}

	return sub, nil
}

// Activate transitions a subscription to active.
func (t *SubscriptionTransitioner) Activate(ctx context.Context, sub *Subscription) (*Subscription, error) {
	if sub.Status == StatusActive {
		return sub, nil
	}

	sub.Status = StatusActive
	sub.PendingExpiresAt = nil
	sub.UpdatedAt = t.clock.Now()

	if err := t.subscriptions.Save(ctx, sub); err != nil {
		return nil, err
	}

	_ = t.entitlements.SyncActivationState(ctx, sub)

	return sub, nil
}

// MarkPastDue transitions a subscription to past due.
func (t *SubscriptionTransitioner) MarkPastDue(ctx context.Context, sub *Subscription) (*Subscription, error) {
	if sub.Status == StatusPastDue {
		return sub, nil
	}

	sub.Status = StatusPastDue
	sub.UpdatedAt = t.clock.Now()

	if err := t.subscriptions.Save(ctx, sub); err != nil {
		return nil, err
	}

	_ = t.entitlements.SyncActivationState(ctx, sub)

	return sub, nil
}

// Pause transitions a subscription to paused.
func (t *SubscriptionTransitioner) Pause(ctx context.Context, sub *Subscription) (*Subscription, error) {
	if sub.Status == StatusPaused {
		return sub, nil
	}

	now := t.clock.Now()
	sub.Status = StatusPaused
	sub.PausedAt = &now
	sub.UpdatedAt = now

	if err := t.subscriptions.Save(ctx, sub); err != nil {
		return nil, err
	}

	_ = t.entitlements.SyncActivationState(ctx, sub)

	return sub, nil
}

// Cancel transitions a subscription to canceled.
func (t *SubscriptionTransitioner) Cancel(ctx context.Context, sub *Subscription) (*Subscription, error) {
	if sub.Status == StatusCanceled {
		return sub, nil
	}

	sub.Status = StatusCanceled
	sub.UpdatedAt = t.clock.Now()

	if err := t.subscriptions.Save(ctx, sub); err != nil {
		return nil, err
	}

	_ = t.entitlements.SyncActivationState(ctx, sub)

	return sub, nil
}

// Expire transitions a subscription to expired.
func (t *SubscriptionTransitioner) Expire(ctx context.Context, sub *Subscription) (*Subscription, error) {
	if sub.Status == StatusExpired {
		return sub, nil
	}

	sub.Status = StatusExpired
	sub.UpdatedAt = t.clock.Now()

	if err := t.subscriptions.Save(ctx, sub); err != nil {
		return nil, err
	}

	_ = t.entitlements.SyncActivationState(ctx, sub)

	return sub, nil
}

// ExpireStaleSubscriptions finds and expires subscriptions whose pending
// hold window has elapsed.
func (t *SubscriptionTransitioner) ExpireStaleSubscriptions(ctx context.Context) (int, error) {
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
