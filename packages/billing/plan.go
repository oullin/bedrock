package billing

import (
	"fmt"
	"time"
)

// Plan represents a subscription plan stored in the application database.
type Plan struct {
	ID                int64
	UUID              string
	Code              string // Unique plan code (e.g. "pro").
	Name              string
	Slug              string
	ProviderProductID string // Product ID on the payment provider.
	Description       string
	MarketingFeatures []string
	CTALabel          string
	CTAStyle          string
	Featured          bool
	Badge             string
	IsActive          bool
	CreatedAt         time.Time
	UpdatedAt         time.Time

	// Loaded relations.
	PlanPeriods []PlanPeriod
	PlanFeatures []PlanFeature
}

// PlanPeriod represents a billing period variant for a plan (e.g. monthly, yearly).
type PlanPeriod struct {
	ID         int64
	UUID       string
	PlanID     int64
	PeriodCode BillingPeriod
	CreatedAt  time.Time
	UpdatedAt  time.Time

	// Loaded relations.
	Prices []PlanPeriodPrice
}

// PlanPeriodPrice represents the pricing for a specific plan period in a
// specific currency.
type PlanPeriodPrice struct {
	ID                int64
	UUID              string
	PlanPeriodID      int64
	AmountMinor       *int64 // Price in minor units (cents); nil for non-money pricing modes.
	Currency          string // ISO 4217 code.
	PricingMode       PlanPricingMode
	Period            string // Denormalised period label.
	ProviderPriceID   string // Price ID on the payment provider.
	IsActive          bool
	IsCheckoutEnabled bool
	ActiveGuard       string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

// DisplayAmount returns a human-readable price string based on the pricing mode.
func (p *PlanPeriodPrice) DisplayAmount(formatter CurrencyFormatter) string {
	switch p.PricingMode {
	case PricingModeFree:
		return "Free"
	case PricingModeCustom:
		return "Custom"
	case PricingModeMoney:
		if p.AmountMinor == nil || p.Currency == "" {
			return "N/A"
		}

		if formatter != nil {
			return formatter.FormatAmount(*p.AmountMinor, p.Currency, "")
		}

		return fmt.Sprintf("%d %s", *p.AmountMinor, p.Currency)
	default:
		return string(p.PricingMode)
	}
}

// Feature represents a billable feature definition.
type Feature struct {
	ID          int64
	UUID        string
	Code        string
	Name        string
	Description string
	ValueType   string // Type of value (e.g. "integer", "boolean").
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// PlanFeature links a plan to a feature with an optional value.
type PlanFeature struct {
	ID        int64
	PlanID    int64
	FeatureID int64
	Value     *string // The feature value for this plan (e.g. "10" for capacity).
	CreatedAt time.Time
	UpdatedAt time.Time

	// Loaded relations.
	Feature *Feature
}

// BillingPlan represents a plan as registered with the Billing manager. This is
// the provider-facing plan definition used by the billing portal.
type BillingPlan struct {
	ID               string // Provider price ID.
	Name             string
	Interval         string // "monthly" or "yearly".
	ShortDescription string
	Features         []string
	Options          map[string]any
	Active           bool
	MonthlyIncentive string
	YearlyIncentive  string
	PriceIncludesVAT bool
	Price            float64
	Currency         string
}
