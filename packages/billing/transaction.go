package billing

import "time"

// Transaction represents a payment transaction recorded by the provider.
type Transaction struct {
	ID                     int64
	BillableType           string
	BillableID             int64
	ProviderID             string // Paddle transaction ID.
	ProviderSubscriptionID string // Linked subscription on the provider side.
	InvoiceNumber          string
	Status                 TransactionStatus
	Total                  int64  // Amount in minor units (cents).
	Tax                    int64  // Tax in minor units.
	Currency               string // ISO 4217 code.
	BilledAt               time.Time
	CreatedAt              time.Time
	UpdatedAt              time.Time
}
