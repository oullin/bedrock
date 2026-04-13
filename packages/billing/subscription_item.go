package billing

import "time"

// SubscriptionItem represents a single line item within a subscription.
// Mirrors Upstream\Paddle\SubscriptionItem.
type SubscriptionItem struct {
	ID             int64
	SubscriptionID int64
	ProductID      string
	PriceID        string
	Status         string
	Quantity       int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
