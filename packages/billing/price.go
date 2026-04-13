package billing

// Price represents a provider price with billing interval information.
// Mirrors Upstream\Paddle\Price.
type Price struct {
	ID        string
	Amount    int64
	Currency  string
	Interval  SubscriptionInterval
	Frequency int
}

// RawAmount returns the price amount in minor units.
func (p *Price) RawAmount() int64 {
	return p.Amount
}

// FormattedAmount returns the price as a formatted currency string.
func (p *Price) FormattedAmount() string {
	return formatMinorUnits(p.Amount, p.Currency)
}
