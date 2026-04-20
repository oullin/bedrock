package exchange

import (
	"math"
	"sync"

	"github.com/bedrock/packages/money/exception"
)

// Exchange provides currency conversion functionality.
// Safe for concurrent use by multiple goroutines.
type Exchange struct {
	mu    sync.RWMutex
	rates map[string]map[string]float64
}

// NewExchange creates a new Exchange instance.
func NewExchange() *Exchange {
	return &Exchange{
		rates: make(map[string]map[string]float64),
	}
}

// IsValid checks if the Exchange is properly initialised with a valid rates map
func (e *Exchange) IsValid() bool {
	if e == nil {
		return false
	}

	e.mu.RLock()

	defer e.mu.RUnlock()

	return e.rates != nil && len(e.rates) > 0
}

// IsInvalid checks if the Exchange is not properly initialised or has no rates
func (e *Exchange) IsInvalid() bool {
	return !e.IsValid()
}

// AddRate adds or updates a conversion rate between two tables.
func (e *Exchange) AddRate(baseCurrency, counterCurrency string, rate float64) error {
	if e == nil || rate <= 0 {
		return exception.ErrInvalidExchangeRate
	}

	e.mu.Lock()

	defer e.mu.Unlock()

	if e.rates[baseCurrency] == nil {
		e.rates[baseCurrency] = make(map[string]float64)
	}

	e.rates[baseCurrency][counterCurrency] = rate

	return nil
}

// GetRate retrieves the conversion rate between two tables.
// If a direct rate is not found, it attempts to use the inverse rate.
func (e *Exchange) GetRate(baseCurrency, counterCurrency string) (float64, error) {
	if e == nil {
		return -1, exception.ErrInvalidExchangeRate
	}

	if baseCurrency == counterCurrency {
		return 1.0, nil
	}

	e.mu.RLock()

	defer e.mu.RUnlock()

	if e.rates[baseCurrency] != nil {
		if rate, ok := e.rates[baseCurrency][counterCurrency]; ok {
			return rate, nil
		}
	}

	// Try inverse rate
	if e.rates[counterCurrency] != nil {
		if rate, ok := e.rates[counterCurrency][baseCurrency]; ok {
			return 1.0 / rate, nil
		}
	}

	return 0, exception.ErrCurrencyConversionNotFound
}

// ConvertAmount converts an amount from one currency to another using the stored exchange rates.
func (e *Exchange) ConvertAmount(amount int64, fromCurrencyCode string, fromFraction int, toCurrencyCode string, toFraction int) (int64, error) {
	if e == nil {
		return -1, exception.ErrInvalidExchangeRate
	}

	if fromCurrencyCode == toCurrencyCode {
		return amount, nil
	}

	rate, err := e.GetRate(fromCurrencyCode, toCurrencyCode)

	if err != nil {
		return 0, err
	}

	// Convert the amount
	// First convert to major units (float), apply rate, then back to minor units
	fromFractionPow := math.Pow10(fromFraction)
	toFractionPow := math.Pow10(toFraction)

	majorUnits := float64(amount) / fromFractionPow
	convertedMajorUnits := majorUnits * rate

	// Use proper rounding to handle precision correctly.
	// For financial calculations, this ensures fractional values are rounded appropriately.
	convertedAmount := int64(math.Round(convertedMajorUnits * toFractionPow))

	return convertedAmount, nil
}

// ConvertAmountWithRate converts an amount from one currency to another using a provided exchange rate.
func (e *Exchange) ConvertAmountWithRate(amount int64, fromFraction int, toFraction int, rate float64) (int64, error) {
	if e == nil || rate <= 0 {
		return 0, exception.ErrInvalidExchangeRate
	}

	// Convert the amount
	fromFractionPow := math.Pow10(fromFraction)
	toFractionPow := math.Pow10(toFraction)

	majorUnits := float64(amount) / fromFractionPow
	convertedMajorUnits := majorUnits * rate

	// Use proper rounding to handle precision correctly.
	// For financial calculations, this ensures fractional values are rounded appropriately.
	convertedAmount := int64(math.Round(convertedMajorUnits * toFractionPow))

	return convertedAmount, nil
}
