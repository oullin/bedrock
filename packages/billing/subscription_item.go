package billing

import "time"

// SubscriptionItem represents a line item within a subscription, linking it
// to a specific product and price on the payment provider.
type SubscriptionItem struct {
	ID             int64
	UUID           string
	SubscriptionID int64
	ProductID      string
	PriceID        string
	Status         string
	Quantity       int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
