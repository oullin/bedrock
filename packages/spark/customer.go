package spark

import "time"

// Customer represents a payment-provider customer record linked to a
// billable entity. Mirrors Laravel\Paddle\Customer.
type Customer struct {
	ID                int64
	BillableType      string
	BillableID        int64
	PaddleID          string
	PendingCheckoutID string
	Name              string
	Email             string
	TrialEndsAt       *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// OnGenericTrial reports whether the customer is currently on a
// generic (non-subscription) trial period.
func (c *Customer) OnGenericTrial() bool {
	if c.TrialEndsAt == nil {
		return false
	}

	return c.TrialEndsAt.After(time.Now())
}

// HasExpiredGenericTrial reports whether the customer had a generic
// trial that has now ended.
func (c *Customer) HasExpiredGenericTrial() bool {
	if c.TrialEndsAt == nil {
		return false
	}

	return !c.TrialEndsAt.After(time.Now())
}
