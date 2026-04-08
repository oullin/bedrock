package spark

// Price represents a pricing entry returned by the payment provider.
type Price struct {
	ID        string // Provider price ID.
	Amount    int64  // Amount in minor units.
	Currency  string // ISO 4217 code.
	Interval  SubscriptionInterval
	Frequency int // Number of intervals per billing cycle (e.g. 1 for monthly).
}

// RawAmount returns the amount in minor units.
func (p Price) RawAmount() int64 {
	return p.Amount
}
