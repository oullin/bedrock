package billing

import (
	"fmt"
	"strconv"
	"strings"
)

const DefaultCurrency = "USD"

var supportedCurrencies = map[string]struct{}{
	"EUR": {},
	"GBP": {},
	"SGD": {},
	"USD": {},
}

// ParseMinorAmount normalises integer-like provider payload values into
// canonical minor currency units.
func ParseMinorAmount(value any) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		if uint64(v) > uint64(^uint64(0)>>1) {
			return 0, fmt.Errorf("minor amount overflows int64")
		}

		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		if v > uint64(^uint64(0)>>1) {
			return 0, fmt.Errorf("minor amount overflows int64")
		}

		return int64(v), nil
	case string:
		return strconv.ParseInt(strings.TrimSpace(v), 10, 64)
	default:
		return 0, fmt.Errorf("unsupported minor amount type %T", value)
	}
}

// ValidateTransactionMoney validates canonical transaction money fields.
func ValidateTransactionMoney(total, tax int64, currency string) error {
	var errs ValidationErrors

	if total < 0 {
		errs = append(errs, ValidationError{Field: "total", Message: "must be zero or greater"})
	}

	if tax < 0 {
		errs = append(errs, ValidationError{Field: "tax", Message: "must be zero or greater"})
	}

	if !IsSupportedCurrency(currency) {
		errs = append(errs, ValidationError{Field: "currency", Message: "is not supported"})
	}

	if errs.HasErrors() {
		return errs
	}

	return nil
}

// ValidatePriceMoney validates canonical plan-period price money fields.
func ValidatePriceMoney(amount int64, currency string) error {
	var errs ValidationErrors

	if amount < 0 {
		errs = append(errs, ValidationError{Field: "amount", Message: "must be zero or greater"})
	}

	if !IsSupportedCurrency(currency) {
		errs = append(errs, ValidationError{Field: "currency", Message: "is not supported"})
	}

	if errs.HasErrors() {
		return errs
	}

	return nil
}

// IsSupportedCurrency reports whether the ISO currency code is known to
// Billing's canonical billing schema.
func IsSupportedCurrency(currency string) bool {
	_, ok := supportedCurrencies[strings.ToUpper(strings.TrimSpace(currency))]

	return ok
}
