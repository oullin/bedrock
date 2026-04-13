package spark

import "time"

// Subscription represents a billable's subscription record.
// Mirrors Laravel\Paddle\Subscription.
type Subscription struct {
	ID           int64
	BillableType string
	BillableID   int64
	Type         string
	PaddleID     string
	Status       SubscriptionStatus
	TrialEndsAt  *time.Time
	PausedAt     *time.Time
	EndsAt       *time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Items        []SubscriptionItem

	prorationBehavior ProrationBehavior
}

// Active reports whether the subscription status is active.
func (s *Subscription) Active() bool {
	return s.Status == StatusActive
}

// OnTrial reports whether the subscription is currently trialing.
func (s *Subscription) OnTrial() bool {
	if s.Status != StatusTrialing {
		return false
	}

	if s.TrialEndsAt == nil {
		return false
	}

	return s.TrialEndsAt.After(time.Now())
}

// PastDue reports whether the subscription is past due.
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

// OnGracePeriod reports whether the subscription has been canceled
// but is still within its paid period.
func (s *Subscription) OnGracePeriod() bool {
	if !s.Canceled() {
		return false
	}

	if s.EndsAt == nil {
		return false
	}

	return s.EndsAt.After(time.Now())
}

// Valid reports whether the subscription currently grants access.
// This includes active, trialing, past-due, and grace-period states.
func (s *Subscription) Valid() bool {
	return s.Active() || s.OnTrial() || s.PastDue() || s.OnGracePeriod()
}

// Recurring reports whether the subscription is active and not on trial.
func (s *Subscription) Recurring() bool {
	return s.Active() && !s.OnTrial()
}

// HasProduct reports whether any subscription item matches the product ID.
func (s *Subscription) HasProduct(productID string) bool {
	for _, item := range s.Items {
		if item.ProductID == productID {
			return true
		}
	}

	return false
}

// HasPrice reports whether any subscription item matches the price ID.
func (s *Subscription) HasPrice(priceID string) bool {
	for _, item := range s.Items {
		if item.PriceID == priceID {
			return true
		}
	}

	return false
}

// FindItemByPrice returns the subscription item matching the given
// price ID, or nil if not found.
func (s *Subscription) FindItemByPrice(priceID string) *SubscriptionItem {
	for i := range s.Items {
		if s.Items[i].PriceID == priceID {
			return &s.Items[i]
		}
	}

	return nil
}

// Prorate sets the proration behavior to prorate on the next billing period.
func (s *Subscription) Prorate() *Subscription {
	s.prorationBehavior = ProrateNextBilling

	return s
}

// NoProrate sets the proration behavior to charge the full amount on the next billing period.
func (s *Subscription) NoProrate() *Subscription {
	s.prorationBehavior = FullNextBilling

	return s
}

// ProrateImmediately sets the proration behavior to charge prorated immediately.
func (s *Subscription) ProrateImmediately() *Subscription {
	s.prorationBehavior = ProrateImmediately

	return s
}

// ImmediatelyWithoutProrate sets the proration behavior to charge full amount immediately.
func (s *Subscription) ImmediatelyWithoutProrate() *Subscription {
	s.prorationBehavior = FullImmediately

	return s
}

// DoNotBill sets the proration behavior to not bill the customer.
func (s *Subscription) DoNotBill() *Subscription {
	s.prorationBehavior = DoNotBill

	return s
}

// ProrationBehavior returns the current proration behavior setting.
func (s *Subscription) ProrationBehavior() ProrationBehavior {
	if s.prorationBehavior == "" {
		return ProrateNextBilling
	}

	return s.prorationBehavior
}
