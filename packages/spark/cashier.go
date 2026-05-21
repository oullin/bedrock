package spark

import (
	"fmt"
	"strings"
)

// FormatAmount formats a minor-unit amount as a human-readable currency
// string. Mirrors the underlying behavior.
func FormatAmount(amount int64, currency string) string {
	major := amount / 100
	minor := amount % 100

	if minor < 0 {
		minor = -minor
	}

	formatted := fmt.Sprintf("%d.%02d", major, minor)

	// Trim trailing ".00" like upstream does.
	formatted = strings.TrimSuffix(formatted, ".00")
	formatted = strings.TrimSuffix(formatted, ".0")

	return fmt.Sprintf("%s %s", strings.ToUpper(currency), formatted)
}

// NormalizeItems converts a mixed set of price/quantity pairs into a
// uniform slice of CheckoutItems.
func NormalizeItems(items []CheckoutItem) []CheckoutItem {
	normalised := make([]CheckoutItem, len(items))

	for i, item := range items {
		q := item.Quantity

		if q < 1 {
			q = 1
		}

		normalised[i] = CheckoutItem{PriceID: item.PriceID, Quantity: q}
	}

	return normalised
}

// SupportedProviders lists all known payment providers.
var SupportedProviders = []string{"stripe", "paddle"}

// IsValidProvider reports whether the given string is a supported
// payment provider.
func IsValidProvider(provider string) bool {
	for _, p := range SupportedProviders {
		if p == provider {
			return true
		}
	}

	return false
}
