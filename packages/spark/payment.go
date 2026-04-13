package spark

import "time"

// Payment represents a single payment amount with currency and date.
// Mirrors Laravel\Paddle\Payment.
type Payment struct {
	Amount   int64
	Currency string
	Date     time.Time
}

// RawAmount returns the payment amount in minor units.
func (p *Payment) RawAmount() int64 {
	return p.Amount
}

// FormattedAmount returns the payment amount as a formatted string.
func (p *Payment) FormattedAmount() string {
	return formatMinorUnits(p.Amount, p.Currency)
}
