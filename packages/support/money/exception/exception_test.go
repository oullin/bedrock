package exception

import (
	"errors"
	"testing"
)

var sentinelErrors = []struct {
	name string
	err  error
	msg  string
}{
	{name: "ErrCurrencyMismatch", err: ErrCurrencyMismatch, msg: "currencies don't match"},
	{name: "ErrInvalidJSONUnmarshal", err: ErrInvalidJSONUnmarshal, msg: "invalid json unmarshal"},
	{name: "ErrInvalidExchangeRate", err: ErrInvalidExchangeRate, msg: "invalid exchange rate"},
	{name: "ErrCurrencyConversionNotFound", err: ErrCurrencyConversionNotFound, msg: "currency conversion rate not found"},
	{name: "ErrNoMoneyProvided", err: ErrNoMoneyProvided, msg: "no money objects provided"},
	{name: "ErrInvalidMoneyString", err: ErrInvalidMoneyString, msg: "invalid money string format"},
	{name: "ErrCurrencyNotSpecified", err: ErrCurrencyNotSpecified, msg: "currency not specified or detected"},
	{name: "ErrNoMultipliersProvided", err: ErrNoMultipliersProvided, msg: "no multipliers provided"},
	{name: "ErrNoJSONParserProvided", err: ErrNoJSONParserProvided, msg: "no json parser provided"},
	{name: "ErrNoConverterProvided", err: ErrNoConverterProvided, msg: "no converter provided"},
	{name: "ErrCurrencyNotFound", err: ErrCurrencyNotFound, msg: "currency not found"},
	{name: "ErrJSONUnmarshalFuncNil", err: ErrJSONUnmarshalFuncNil, msg: "money.JSON: unmarshal function cannot be nil"},
	{name: "ErrJSONMarshalFuncNil", err: ErrJSONMarshalFuncNil, msg: "money.JSON: marshal function cannot be nil"},
	{name: "ErrEmptyAmountString", err: ErrEmptyAmountString, msg: "amount string cannot be empty"},
	{name: "ErrInvalidAmountMultiple", err: ErrInvalidAmountMultiple, msg: "invalid amount: multiple decimal points"},
	{name: "ErrInvalidAmountFraction", err: ErrInvalidAmountFraction, msg: "too many decimal places for curr"},
	{name: "ErrInvalidAmount", err: ErrInvalidAmount, msg: "invalid amount"},
	{name: "ErrInvalidSplit", err: ErrInvalidSplit, msg: "split must be higher than zero"},
	{name: "ErrNoRatiosProvided", err: ErrNoRatiosProvided, msg: "no ratios specified"},
	{name: "ErrNegativeRatios", err: ErrNegativeRatios, msg: "negative ratios not allowed"},
	{name: "ErrRatiosExceedMaxInt", err: ErrRatiosExceedMaxInt, msg: "sum of given ratios exceeds max int"},
	{name: "ErrNoCurrencyInstance", err: ErrNoCurrencyInstance, msg: "money instance has no currency"},
	{name: "ErrNoCurrencyManager", err: ErrNoCurrencyManager, msg: "currency manager cannot be nil"},
	{name: "ErrNoCurrencyMapDataset", err: ErrNoCurrencyMapDataset, msg: "currency map dataset cannot be nil or empty"},
	{name: "ErrInvalidAggregatorProvider", err: ErrInvalidAggregatorProvider, msg: "invalid aggregator: nil manager"},
	{name: "ErrOverflow", err: ErrOverflow, msg: "arithmetic operation resulted in overflow"},
}

func TestErrorsAreNotNil(t *testing.T) {
	for _, tt := range sentinelErrors {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Fatalf("%s should not be nil", tt.name)
			}
		})
	}
}

func TestErrorMessages(t *testing.T) {
	for _, tt := range sentinelErrors {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.msg {
				t.Fatalf("%s.Error() = %q, want %q", tt.name, got, tt.msg)
			}
		})
	}
}

func TestErrorsComparison(t *testing.T) {
	for _, tt := range sentinelErrors {
		t.Run(tt.name, func(t *testing.T) {
			// Test that "errors.Is works" correctly
			if !errors.Is(tt.err, tt.err) {
				t.Fatalf("errors.Is(%s, %s) should return true", tt.name, tt.name)
			}

			// Test that the error doesn't match a different error
			differentErr := errors.New("different error")

			if errors.Is(tt.err, differentErr) {
				t.Fatalf("errors.Is(%s, differentErr) should return false", tt.name)
			}
		})
	}
}

func TestErrorsAreUnique(t *testing.T) {
	allErrors := All()

	if got, want := len(allErrors), len(sentinelErrors); got != want {
		t.Fatalf("All() returned %d errors, want %d", got, want)
	}

	for _, tt := range sentinelErrors {
		found := false

		for _, err := range allErrors {
			if errors.Is(err, tt.err) {
				found = true

				break
			}
		}

		if !found {
			t.Fatalf("All() missing error %s", tt.name)
		}
	}

	// Verify each error is unique (different instance as another)
	for i, err1 := range allErrors {
		for j, err2 := range allErrors {
			if i != j && errors.Is(err2, err1) {
				t.Fatalf("Errors at index %d and %d are the same instance: %v", i, j, err1)
			}
		}
	}
}

func TestErrorWrapping(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wrapMsg string
	}{
		{name: "ErrCurrencyMismatch", err: ErrCurrencyMismatch, wrapMsg: "operation failed"},
		{name: "ErrInvalidJSONUnmarshal", err: ErrInvalidJSONUnmarshal, wrapMsg: "parsing failed"},
		{name: "ErrInvalidExchangeRate", err: ErrInvalidExchangeRate, wrapMsg: "rate error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapped := errors.Join(errors.New(tt.wrapMsg), tt.err)

			if !errors.Is(wrapped, tt.err) {
				t.Fatalf("wrapped error should contain %s", tt.name)
			}
		})
	}
}

func TestNewErrInvalidJSONUnmarshalFrom(t *testing.T) {
	orig := errors.New("boom")

	wrapped := NewErrInvalidJSONUnmarshalFrom(orig)

	if !errors.Is(wrapped, ErrInvalidJSONUnmarshal) {
		t.Fatalf("expected wrapped error to include ErrInvalidJSONUnmarshal, got %v", wrapped)
	}

	if !errors.Is(wrapped, orig) {
		t.Fatalf("expected wrapped error to include original error, got %v", wrapped)
	}
}
