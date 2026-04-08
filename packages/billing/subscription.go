package billing

import "time"

// Subscription represents a billing subscription linking a billable entity to
// a plan through a payment provider.
type Subscription struct {
	ID                int64
	UUID              string
	BillableType      string
	BillableID        int64
	Type              string // Subscription type (default: "default").
	ProviderID        string // Payment provider subscription ID.
	Status            SubscriptionStatus
	Plan              string // Plan code (e.g. "pro").
	BillingPeriod     BillingPeriod
	PlanPeriodPriceID *int64
	PendingExpiresAt  *time.Time
	PaymentReadyAt    *time.Time
	TrialEndsAt       *time.Time
	PausedAt          *time.Time
	EndsAt            *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time

	// Loaded relations.
	Items    []SubscriptionItem
	Features []SubscriptionFeature

	// Transient state for proration control.
	proration ProrationBehavior
}

// IsActive reports whether the subscription grants access to features.
func (s *Subscription) IsActive() bool {
	return s.Status.GrantsAccess()
}

// IsNew reports whether the subscription has not been persisted yet.
func (s *Subscription) IsNew() bool {
	return s.ID == 0
}

// Valid reports whether the subscription should be considered valid for access.
// When keepPastDueActive is true, past-due subscriptions count as valid.
func (s *Subscription) Valid(clock Clock, keepPastDueActive bool) bool {
	if s.OnTrial(clock) {
		return true
	}

	if s.Active() {
		return true
	}

	if keepPastDueActive && s.PastDue() {
		return true
	}

	return false
}

// OnTrial reports whether the subscription is currently within its trial period.
func (s *Subscription) OnTrial(clock Clock) bool {
	if s.TrialEndsAt == nil {
		return false
	}

	return clock.Now().Before(*s.TrialEndsAt)
}

// HasExpiredTrial reports whether the subscription had a trial that has elapsed.
func (s *Subscription) HasExpiredTrial(clock Clock) bool {
	if s.TrialEndsAt == nil {
		return false
	}

	now := clock.Now()

	return now.After(*s.TrialEndsAt) || now.Equal(*s.TrialEndsAt)
}

// Active reports whether the status is active.
func (s *Subscription) Active() bool {
	return s.Status == StatusActive
}

// Recurring reports whether the subscription is active and not on trial.
func (s *Subscription) Recurring(clock Clock) bool {
	return s.Active() && !s.OnTrial(clock)
}

// PastDue reports whether the subscription has a past-due payment.
func (s *Subscription) PastDue() bool {
	return s.Status == StatusPastDue
}

// Paused reports whether the subscription is paused.
func (s *Subscription) Paused() bool {
	return s.Status == StatusPaused
}

// Canceled reports whether the subscription has been canceled.
func (s *Subscription) Canceled() bool {
	return s.Status == StatusCanceled
}

// OnGracePeriod reports whether the subscription is canceled but still within
// its paid period.
func (s *Subscription) OnGracePeriod(clock Clock) bool {
	if s.EndsAt == nil {
		return false
	}

	return s.Canceled() && clock.Now().Before(*s.EndsAt)
}

// OnPausedGracePeriod reports whether the subscription is paused but the pause
// has not yet taken effect.
func (s *Subscription) OnPausedGracePeriod(clock Clock) bool {
	if s.PausedAt == nil {
		return false
	}

	return s.Paused() && clock.Now().Before(*s.PausedAt)
}

// HasProduct reports whether any item in the subscription matches the given
// product ID.
func (s *Subscription) HasProduct(productID string) bool {
	for _, item := range s.Items {
		if item.ProductID == productID {
			return true
		}
	}

	return false
}

// HasPrice reports whether any item in the subscription matches the given
// price ID.
func (s *Subscription) HasPrice(priceID string) bool {
	for _, item := range s.Items {
		if item.PriceID == priceID {
			return true
		}
	}

	return false
}

// HasMultiplePrices reports whether the subscription contains more than one
// price item.
func (s *Subscription) HasMultiplePrices() bool {
	return len(s.Items) > 1
}

// HasSinglePrice reports whether the subscription contains exactly one price
// item.
func (s *Subscription) HasSinglePrice() bool {
	return len(s.Items) == 1
}

// ProrationBehavior returns the proration setting for this subscription.
func (s *Subscription) ProrationBehavior() ProrationBehavior {
	if s.proration == "" {
		return ProratedNextBillingPeriod
	}

	return s.proration
}

// SetProrationBehavior sets the proration behaviour for subsequent operations.
func (s *Subscription) SetProrationBehavior(b ProrationBehavior) {
	s.proration = b
}

// Prorate sets the proration to prorated next billing period.
func (s *Subscription) Prorate() { s.proration = ProratedNextBillingPeriod }

// NoProrate sets the proration to full next billing period.
func (s *Subscription) NoProrate() { s.proration = FullNextBillingPeriod }

// ProrateImmediately sets the proration to prorated immediately.
func (s *Subscription) ProrateImmediately() { s.proration = ProratedImmediately }

// ImmediatelyWithoutProrate sets the proration to full immediately.
func (s *Subscription) ImmediatelyWithoutProrate() { s.proration = FullImmediately }

// DoNotBill sets the proration to do not bill.
func (s *Subscription) DoNotBill() { s.proration = DoNotBill }
