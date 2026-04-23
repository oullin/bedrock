package spark

import (
	"context"
	"time"
)

// Subscription represents a billable's subscription record.
// Mirrors Laravel\Paddle\Subscription.
type Subscription struct {
	ID               int64
	BillableType     string
	BillableID       int64
	Type             string
	Plan             string
	PaddleID         string
	Status           SubscriptionStatus
	PendingExpiresAt *time.Time
	PaymentReadyAt   *time.Time
	TrialEndsAt      *time.Time
	PausedAt         *time.Time
	EndsAt           *time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Items            []SubscriptionItem

	prorationBehavior ProrationBehavior
}

// NewPendingSubscription creates a local pending subscription record using
// Spark's default hold window.

// NewTrialSubscription creates a local trial subscription record.

// Active reports whether the subscription status is active.

// OnTrial reports whether the subscription is currently trialing.

// PastDue reports whether the subscription is past due.

// Paused reports whether the subscription is paused.

// Canceled reports whether the subscription has been canceled.

// OnGracePeriod reports whether the subscription has been canceled
// but is still within its paid period.

// Valid reports whether the subscription currently grants access.
// This includes active, trialing, past-due, and grace-period states.

// Recurring reports whether the subscription is active and not on trial.

// HasProduct reports whether any subscription item matches the product ID.

// HasPrice reports whether any subscription item matches the price ID.

// FindItemByPrice returns the subscription item matching the given
// price ID, or nil if not found.

// Prorate sets the proration behavior to prorate on the next billing period.

// NoProrate sets the proration behavior to charge the full amount on the next billing period.

// ProrateImmediately sets the proration behavior to charge prorated immediately.

// ImmediatelyWithoutProrate sets the proration behavior to charge full amount immediately.

// DoNotBill sets the proration behavior to not bill the customer.

// ProrationBehavior returns the current proration behavior setting.

// MarkPaymentReady transitions a pending subscription into awaiting payment.

// Activate transitions the subscription to active.

// Expire transitions the subscription to expired.

// MarkPastDue transitions the subscription to past due.

// Pause transitions the subscription to paused.

// Cancel transitions the subscription to canceled and starts its grace period.

// Resume reactivates a canceled subscription that is still on its grace period.

// ExpirableSubscriptionStore exposes pending-like subscriptions that may have
// outlived their local checkout hold window.
type ExpirableSubscriptionStore interface {
	ExpirableSubscriptions(ctx context.Context, now time.Time) ([]*Subscription, error)
	Delete(ctx context.Context, id int64) error
}

const DefaultPendingExpiryDays = 14

func NewPendingSubscription(billable Billable, plan string, now time.Time) *Subscription {
	pendingExpiresAt := now.AddDate(0, 0, DefaultPendingExpiryDays)

	return &Subscription{
		BillableType:     billable.BillableType(),
		BillableID:       billable.BillableID(),
		Type:             DefaultSubscriptionType,
		Plan:             plan,
		Status:           StatusPending,
		PendingExpiresAt: &pendingExpiresAt,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
}

func NewTrialSubscription(billable Billable, plan string, trialDays int, now time.Time) *Subscription {
	sub := NewPendingSubscription(billable, plan, now)
	trialEndsAt := now.AddDate(0, 0, trialDays)
	sub.Status = StatusTrialing
	sub.TrialEndsAt = &trialEndsAt

	return sub
}

func (s *Subscription) Active() bool {
	return s.Status == StatusActive
}

func (s *Subscription) OnTrial() bool {
	if s.Status != StatusTrialing {
		return false
	}

	if s.TrialEndsAt == nil {
		return false
	}

	return s.TrialEndsAt.After(time.Now())
}

func (s *Subscription) PastDue() bool {
	return s.Status == StatusPastDue
}

func (s *Subscription) Paused() bool {
	return s.Status == StatusPaused
}

func (s *Subscription) Canceled() bool {
	return s.Status == StatusCanceled
}

func (s *Subscription) OnGracePeriod() bool {
	if !s.Canceled() {
		return false
	}

	if s.EndsAt == nil {
		return false
	}

	return s.EndsAt.After(time.Now())
}

func (s *Subscription) Valid() bool {
	return s.Active() || s.OnTrial() || s.PastDue() || s.OnGracePeriod()
}

func (s *Subscription) Recurring() bool {
	return s.Active() && !s.OnTrial()
}

func (s *Subscription) HasProduct(productID string) bool {
	for _, item := range s.Items {
		if item.ProductID == productID {
			return true
		}
	}

	return false
}

func (s *Subscription) HasPrice(priceID string) bool {
	for _, item := range s.Items {
		if item.PriceID == priceID {
			return true
		}
	}

	return false
}

func (s *Subscription) FindItemByPrice(priceID string) *SubscriptionItem {
	for i := range s.Items {
		if s.Items[i].PriceID == priceID {
			return &s.Items[i]
		}
	}

	return nil
}

func (s *Subscription) Prorate() *Subscription {
	s.prorationBehavior = ProrateNextBilling

	return s
}

func (s *Subscription) NoProrate() *Subscription {
	s.prorationBehavior = FullNextBilling

	return s
}

func (s *Subscription) ProrateImmediately() *Subscription {
	s.prorationBehavior = ProrateImmediately

	return s
}

func (s *Subscription) ImmediatelyWithoutProrate() *Subscription {
	s.prorationBehavior = FullImmediately

	return s
}

func (s *Subscription) DoNotBill() *Subscription {
	s.prorationBehavior = DoNotBill

	return s
}

func (s *Subscription) ProrationBehavior() ProrationBehavior {
	if s.prorationBehavior == "" {
		return ProrateNextBilling
	}

	return s.prorationBehavior
}

func (s *Subscription) MarkPaymentReady(now time.Time) bool {
	if s.Status == StatusAwaitingPayment {
		return false
	}

	s.Status = StatusAwaitingPayment

	if s.PaymentReadyAt == nil {
		readyAt := now
		s.PaymentReadyAt = &readyAt
	}

	s.UpdatedAt = now

	return true
}

func (s *Subscription) Activate(now time.Time) bool {
	if s.Status == StatusActive {
		return false
	}

	s.Status = StatusActive
	s.UpdatedAt = now

	return true
}

func (s *Subscription) Expire(now time.Time) bool {
	if s.Status == StatusExpired {
		return false
	}

	s.Status = StatusExpired
	s.UpdatedAt = now

	return true
}

func (s *Subscription) MarkPastDue(now time.Time) bool {
	if s.Status == StatusPastDue {
		return false
	}

	s.Status = StatusPastDue
	s.UpdatedAt = now

	return true
}

func (s *Subscription) Pause(now time.Time) bool {
	if s.Status == StatusPaused {
		return false
	}

	s.Status = StatusPaused

	if s.PausedAt == nil {
		pausedAt := now
		s.PausedAt = &pausedAt
	}

	s.UpdatedAt = now

	return true
}

func (s *Subscription) Cancel(now time.Time) bool {
	if s.Status == StatusCanceled {
		return false
	}

	s.Status = StatusCanceled

	if s.EndsAt == nil {
		endsAt := now
		s.EndsAt = &endsAt
	}

	s.UpdatedAt = now

	return true
}

func (s *Subscription) Resume(now time.Time) bool {
	if !s.OnGracePeriod() {
		return false
	}

	s.Status = StatusActive
	s.EndsAt = nil
	s.UpdatedAt = now

	return true
}

// Expirable reports whether the subscription should be removed by the stale
// pending-subscription batch.
func (s *Subscription) Expirable(now time.Time) bool {
	if s.PendingExpiresAt == nil || s.PendingExpiresAt.After(now) {
		return false
	}

	switch s.Status {
	case StatusPending, StatusAwaitingPayment, StatusTrialing:
		return true
	default:
		return false
	}
}

// ExpireStaleSubscriptions removes pending-like subscriptions whose checkout
// hold window has elapsed.
func ExpireStaleSubscriptions(ctx context.Context, store ExpirableSubscriptionStore, now time.Time) (int, error) {
	subscriptions, err := store.ExpirableSubscriptions(ctx, now)

	if err != nil {
		return 0, err
	}

	expired := 0

	for _, subscription := range append([]*Subscription(nil), subscriptions...) {
		if subscription == nil || !subscription.Expirable(now) {
			continue
		}

		subscription.Expire(now)

		if err := store.Delete(ctx, subscription.ID); err != nil {
			return expired, err
		}

		expired++
	}

	return expired, nil
}
