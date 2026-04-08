package billing

import "time"

// Payment represents a payment event (last payment made, next payment due).
type Payment struct {
	Amount   int64  // Amount in minor units.
	Currency string // ISO 4217 code.
	Date     time.Time
}

// RawAmount returns the amount in minor units.
func (p Payment) RawAmount() int64 { return p.Amount }
