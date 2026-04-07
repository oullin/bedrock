package spark

import "time"

// Customer represents a billable entity's payment provider record.
type Customer struct {
	ID              int64
	BillableType    string
	BillableID      int64
	ProviderID      string // Paddle customer ID.
	Name            string
	Email           string
	TrialEndsAt     *time.Time
	PendingCheckout *string // Provider checkout ID while awaiting webhook.
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// OnGenericTrial reports whether the customer is on a generic (non-subscription)
// trial that has not yet expired.
func (c *Customer) OnGenericTrial(clock Clock) bool {
	if c.TrialEndsAt == nil {
		return false
	}

	return clock.Now().Before(*c.TrialEndsAt)
}

// HasExpiredGenericTrial reports whether the customer had a generic trial that
// has since elapsed.
func (c *Customer) HasExpiredGenericTrial(clock Clock) bool {
	if c.TrialEndsAt == nil {
		return false
	}

	return clock.Now().After(*c.TrialEndsAt) || clock.Now().Equal(*c.TrialEndsAt)
}
