package currency

import (
	"database/sql/driver"
	"fmt"

	"github.com/bedrock/packages/money/exception"
	"github.com/bedrock/packages/money/format"
)

// Currency represents the formatting rules for a specific currency.
type Currency struct {
	// Decimal is the decimal separator (e.g., ".").
	Decimal string
	// A Thousand is the thousand separators (e.g., ",").
	Thousand string
	// Code is the ISO 4217 currency code (e.g. "SGD").
	Code string
	// Fraction is the number of fraction digits (e.g. 2 for cents).
	Fraction int
	// NumericCode is the numeric code of the currency (e.g., "702" for SGD).
	NumericCode string
	// Grapheme is the currency symbol (e.g., "$").
	Grapheme string
	// Template is the formatting template (e.g., "$1").
	Template string
}

// Formatter returns currency formatter representing
func (c *Currency) Formatter() *format.Formatter {
	return format.NewFormatter(
		c.Fraction,
		c.Decimal,
		c.Thousand,
		c.Grapheme,
		c.Template,
	)
}

// Equals checks if two currencies are equal based on their code.
func (c *Currency) Equals(oc *Currency) bool {
	if c == nil || oc == nil {
		return c == oc
	}

	return c.Code == oc.Code
}

// Get returns the currency if it exists, otherwise returns an error.
func (c *Currency) Get() (*Currency, error) {
	if c == nil {
		return nil, exception.ErrCurrencyNotFound
	}

	return c, nil
}

// DbValue implements driver.Valuer to serialize a Currency code into a string for saving to a database
func (c *Currency) DbValue() (driver.Value, error) {
	if c == nil {
		return nil, exception.ErrCurrencyNotFound
	}

	return c.Code, nil
}

// DbScan implements sql.Scanner to deserialize a Currency from a string value read from a database
func (c *Currency) DbScan(src interface{}) error {
	if c == nil {
		return exception.ErrCurrencyNotFound
	}

	var val *Currency

	var code string

	switch v := src.(type) {
	case string:
		code = v
	case []byte:
		code = string(v)
	default:
		return fmt.Errorf("%T is not a supported type for a Currency (store the Currency.Code value as a string only)", src)
	}

	currencies := NewCurrenciesMap()

	if val = currencies.FindByCode(code); val == nil {
		return fmt.Errorf("currency(%#v) returned nil", code)
	}

	// copy the value
	*c = *val

	return nil
}
