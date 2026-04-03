package exception

import (
	"errors"
	"fmt"
)

// Sentinel errors for the money package.
var (
	ErrCurrencyMismatch           = errors.New("currencies don't match")
	ErrCurrencyNotFound           = errors.New("currency not found")
	ErrCurrencyConversionNotFound = errors.New("currency conversion rate not found")
	ErrNoCurrencyInstance         = errors.New("money instance has no currency")
	ErrNoCurrencyManager          = errors.New("currency manager cannot be nil")
	ErrNoCurrencyMapDataset       = errors.New("currency map dataset cannot be nil or empty")

	ErrInvalidJSONUnmarshal = errors.New("invalid json unmarshal")
	ErrInvalidExchangeRate  = errors.New("invalid exchange rate")
	ErrNoJSONParserProvided = errors.New("no json parser provided")
	ErrJSONUnmarshalFuncNil = errors.New("money.JSON: unmarshal function cannot be nil")
	ErrJSONMarshalFuncNil   = errors.New("money.JSON: marshal function cannot be nil")

	ErrNoMoneyProvided      = errors.New("no money objects provided")
	ErrInvalidMoneyString   = errors.New("invalid money string format")
	ErrCurrencyNotSpecified = errors.New("currency not specified or detected")

	ErrNoMultipliersProvided = errors.New("no multipliers provided")

	ErrNoConverterProvided = errors.New("no converter provided")

	ErrEmptyAmountString     = errors.New("amount string cannot be empty")
	ErrInvalidAmountMultiple = errors.New("invalid amount: multiple decimal points")
	ErrInvalidAmountFraction = errors.New("too many decimal places for curr")
	ErrInvalidAmount         = errors.New("invalid amount")
	ErrInvalidSplit          = errors.New("split must be higher than zero")
	ErrNoRatiosProvided      = errors.New("no ratios specified")
	ErrNegativeRatios        = errors.New("negative ratios not allowed")
	ErrRatiosExceedMaxInt    = errors.New("sum of given ratios exceeds max int")

	ErrInvalidAggregatorProvider = errors.New("invalid aggregator: nil manager")

	ErrOverflow = errors.New("arithmetic operation resulted in overflow")

	ErrParserNotProvided  = errors.New("parser was not provided")
	ErrParserInvalidState = errors.New("parser is nil or iso is nil")
)

// NewErrInvalidJSONUnmarshalFrom creates a new ErrInvalidJSONUnmarshal error wrapping the original error.
func NewErrInvalidJSONUnmarshalFrom(err error) error {
	if err == nil {
		return ErrInvalidJSONUnmarshal
	}

	return fmt.Errorf("%w: %w", ErrInvalidJSONUnmarshal, err)
}

// All returns a slice with all sentinel errors defined in this package.
func All() []error {
	return []error{
		ErrCurrencyMismatch,
		ErrInvalidJSONUnmarshal,
		ErrInvalidExchangeRate,
		ErrCurrencyConversionNotFound,
		ErrNoMoneyProvided,
		ErrInvalidMoneyString,
		ErrCurrencyNotSpecified,
		ErrNoMultipliersProvided,
		ErrNoJSONParserProvided,
		ErrNoConverterProvided,
		ErrCurrencyNotFound,
		ErrJSONUnmarshalFuncNil,
		ErrJSONMarshalFuncNil,
		ErrEmptyAmountString,
		ErrInvalidAmountMultiple,
		ErrInvalidAmountFraction,
		ErrInvalidAmount,
		ErrInvalidSplit,
		ErrNoRatiosProvided,
		ErrNegativeRatios,
		ErrRatiosExceedMaxInt,
		ErrNoCurrencyInstance,
		ErrNoCurrencyManager,
		ErrNoCurrencyMapDataset,
		ErrInvalidAggregatorProvider,
		ErrOverflow,
	}
}
