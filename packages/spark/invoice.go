package spark

import "time"

// Invoice represents a billing invoice linked to a team.
type Invoice struct {
	ID                   int64
	UUID                 string
	TeamID               int64
	UserID               *int64
	ProviderTransactionID string
	InvoiceNumber        string
	Total                int64  // Amount in minor units.
	Currency             string // ISO 4217 code.
	Status               string
	HostedInvoiceURL     string
	InvoicePDFURL        string
	PaidAt               *time.Time
	Payload              map[string]any
	CreatedAt            time.Time
	UpdatedAt            time.Time
}
