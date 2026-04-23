package billing

import (
	"fmt"
	"strings"
)

// PlanPricingMode describes how a plan period price should be displayed.
type PlanPricingMode string

// PlanPeriodPrice represents a versioned catalog price for a plan period.
type PlanPeriodPrice struct {
	ID              int64
	AmountMinor     *int64
	Currency        string
	PricingMode     PlanPricingMode
	Period          string
	ProviderPriceID string
	Active          bool
}

const (
	PlanPricingModeMoney  PlanPricingMode = "money"
	PlanPricingModeFree   PlanPricingMode = "free"
	PlanPricingModeCustom PlanPricingMode = "custom"
)

// DisplayAmount returns the user-facing amount for the price.
func (p PlanPeriodPrice) DisplayAmount() string {
	switch p.PricingMode {
	case PlanPricingModeFree:
		return "Free"
	case PlanPricingModeCustom:
		return "Custom"
	default:
		if p.AmountMinor == nil {
			return ""
		}

		return formatMinorUnits(*p.AmountMinor, p.Currency)
	}
}

// Validate reports whether the price has a coherent canonical money shape.
func (p PlanPeriodPrice) Validate() error {
	if p.PricingMode != PlanPricingModeMoney {
		return nil
	}

	var errs ValidationErrors

	if p.AmountMinor == nil {
		errs = append(errs, ValidationError{Field: "amount_minor", Message: "is required"})
	}

	if strings.TrimSpace(p.Currency) == "" {
		errs = append(errs, ValidationError{Field: "currency", Message: "is required"})
	}

	if errs.HasErrors() {
		return errs
	}

	return ValidatePriceMoney(*p.AmountMinor, p.Currency)
}

// RotatePlanPrice appends a new active price while preserving older versions
// as inactive history for the same plan period.
func RotatePlanPrice(history []PlanPeriodPrice, next PlanPeriodPrice) ([]PlanPeriodPrice, error) {
	next.ProviderPriceID = strings.TrimSpace(next.ProviderPriceID)
	next.Active = true

	if err := next.Validate(); err != nil {
		return history, fmt.Errorf("invalid plan price %q: %w", next.ProviderPriceID, err)
	}

	for _, price := range history {
		if price.Active && strings.TrimSpace(price.ProviderPriceID) == next.ProviderPriceID && next.ProviderPriceID != "" {
			return history, fmt.Errorf("active provider price id %q already exists", next.ProviderPriceID)
		}
	}

	rotated := make([]PlanPeriodPrice, len(history), len(history)+1)
	copy(rotated, history)

	for i := range rotated {
		rotated[i].Active = false
	}

	rotated = append(rotated, next)

	return rotated, nil
}
