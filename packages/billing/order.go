package billing

import "time"

// OrderStatus constants.

// Order represents a one-time purchase record linked to a team and
// product. Mirrors app/Models/Order.php.
type Order struct {
	ID                    int64
	TeamID                int64
	ProductID             int64
	PaymentProvider       string
	ProviderTransactionID string
	Amount                int64 // minor units
	Currency              string
	Status                string
	Metadata              map[string]any
	CompletedAt           *time.Time
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

const (
	OrderStatusPending   = "pending"
	OrderStatusCompleted = "completed"
)

// IsCompleted reports whether the order has been fulfilled.
func (o *Order) IsCompleted() bool {
	return o.Status == OrderStatusCompleted
}

// IsPending reports whether the order is awaiting payment confirmation.
func (o *Order) IsPending() bool {
	return o.Status == OrderStatusPending
}

// MarkAsCompleted transitions the order to completed status with the
// provider's transaction identifier.
func (o *Order) MarkAsCompleted(providerTransactionID string, now time.Time) {
	o.Status = OrderStatusCompleted
	o.ProviderTransactionID = providerTransactionID
	o.CompletedAt = &now
	o.UpdatedAt = now
}
