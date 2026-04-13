package billing

import (
	"fmt"
	"time"
)

// ProductType distinguishes subscription products from one-time purchases.

// Product represents a billable product that can be purchased via
// subscription or one-time payment. Supports multiple payment providers.
// Mirrors app/Models/Product.php.
type Product struct {
	ID            int64
	Name          string
	Description   string
	Type          string // ProductTypeSubscription or ProductTypeOneTime
	StripePriceID string
	PaddlePriceID string
	PriceAmount   int64 // in minor units (cents)
	Currency      string
	Interval      string // "monthly", "yearly", etc.
	Active        bool
	Features      []string
	Metadata      map[string]any
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

const (
	ProductTypeSubscription = "subscription"
	ProductTypeOneTime      = "one_time"
)

// PriceIDForProvider returns the provider-specific price ID for the
// given payment provider. Returns an empty string for unknown providers.
func (p *Product) PriceIDForProvider(provider string) string {
	switch provider {
	case "stripe":
		return p.StripePriceID
	case "paddle":
		return p.PaddlePriceID
	default:
		return ""
	}
}

// IsSubscription reports whether the product is a subscription type.
func (p *Product) IsSubscription() bool {
	return p.Type == ProductTypeSubscription
}

// IsOneTime reports whether the product is a one-time purchase type.
func (p *Product) IsOneTime() bool {
	return p.Type == ProductTypeOneTime
}

// FormattedPrice returns the price formatted as a currency string
// (e.g. "9.99").
func (p *Product) FormattedPrice() string {
	major := p.PriceAmount / 100
	minor := p.PriceAmount % 100

	return fmt.Sprintf("%d.%02d", major, minor)
}
