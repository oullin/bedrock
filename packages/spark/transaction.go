package spark

import (
	"fmt"
	"time"
)

// Transaction represents a billing transaction record.
// Mirrors Laravel\Paddle\Transaction.
type Transaction struct {
	ID                   int64
	BillableType         string
	BillableID           int64
	PaddleSubscriptionID string
	PaddleID             string
	InvoiceNumber        string
	Status               TransactionStatus
	Total                int64
	Tax                  int64
	Currency             string
	BilledAt             *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

// TotalFormatted returns the total as a formatted currency string.
func (t *Transaction) TotalFormatted() string {
	return formatMinorUnits(t.Total, t.Currency)
}

// TaxFormatted returns the tax as a formatted currency string.
func (t *Transaction) TaxFormatted() string {
	return formatMinorUnits(t.Tax, t.Currency)
}

func formatMinorUnits(amount int64, currency string) string {
	major := amount / 100
	minor := amount % 100

	if minor < 0 {
		minor = -minor
	}

	return fmt.Sprintf("%s %d.%02d", currency, major, minor)
}
