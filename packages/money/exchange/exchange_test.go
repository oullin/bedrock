package exchange

import (
	"errors"
	"math"
	"sync"
	"testing"

	"github.com/bedrock/packages/money/exception"
	testutil "github.com/bedrock/packages/money/tests"
)

func TestNewExchange(t *testing.T) {
	e := NewExchange()

	if e == nil {
		t.Fatal("NewExchange() returned nil")
	}

	// A newly created Exchange has no rates, so it should be invalid
	if e.IsValid() {
		t.Fatal("NewExchange() returned valid Exchange without any rates")
	}

	err := e.AddRate("SGD", "EUR", 0.85)

	if err != nil {
		t.Fatalf("Failed to add rate: %v", err)
	}

	if e.IsInvalid() {
		t.Fatal("Exchange with rates should be valid")
	}
}

func TestExchangeIsValid(t *testing.T) {
	// Test valid exchange (with rates)
	ex := NewExchange()
	err := ex.AddRate("SGD", "EUR", 0.85)

	if err != nil {
		t.Fatalf("Failed to add rate: %v", err)
	}

	if ex.IsInvalid() {
		t.Fatal("IsInvalid() = true for properly initialized Exchange with rates, want false")
	}

	// Test nil exchange
	var nilEx *Exchange

	if nilEx.IsValid() {
		t.Fatal("IsValid() = true for nil Exchange, want false")
	}

	if !nilEx.IsInvalid() {
		t.Fatal("IsInvalid() = false for nil Exchange, want true")
	}

	// Test exchange with nil rates map
	invalidEx := &Exchange{rates: nil}

	if invalidEx.IsValid() {
		t.Fatal("IsValid() = true for Exchange with nil rates map, want false")
	}

	if !invalidEx.IsInvalid() {
		t.Fatal("IsInvalid() = false for Exchange with nil rates map, want true")
	}

	// Test exchange with empty rates map
	emptyEx := NewExchange()

	if emptyEx.IsValid() {
		t.Fatal("IsValid() = true for Exchange with empty rates map, want false")
	}

	if !emptyEx.IsInvalid() {
		t.Fatal("IsInvalid() = false for Exchange with empty rates map, want true")
	}
}

func TestAddRate(t *testing.T) {
	e := NewExchange()

	err := e.AddRate("SGD", "EUR", 0.85)

	if err != nil {
		t.Errorf("AddRate failed: %v", err)
	}

	err = e.AddRate("SGD", "GBP", -1.0)

	if !errors.Is(err, exception.ErrInvalidExchangeRate) {
		t.Errorf("Expected ErrInvalidExchangeRate, got %v", err)
	}

	var nilExchange *Exchange

	err = nilExchange.AddRate("SGD", "EUR", 0.85)

	if !errors.Is(err, exception.ErrInvalidExchangeRate) {
		t.Errorf("Expected ErrInvalidExchangeRate for nil exchange, got %v", err)
	}
}

func TestGetRate(t *testing.T) {
	const epsilon = 1e-9

	e := NewExchange()
	err := e.AddRate("SGD", "EUR", 0.85)

	if err != nil {
		t.Errorf("AddRate error: %v", err)
	}

	// Test direct rate
	rate, err := e.GetRate("SGD", "EUR")

	if err != nil {
		t.Errorf("GetRate failed: %v", err)
	}

	if rate != 0.85 {
		t.Errorf("Expected rate 0.85, got %f", rate)
	}

	// Test inverse rate
	rate, err = e.GetRate("EUR", "SGD")

	if err != nil {
		t.Errorf("GetRate inverse failed: %v", err)
	}

	expectedInverse := 1.0 / 0.85

	if math.Abs(rate-expectedInverse) > epsilon {
		t.Errorf("Expected inverse rate %f, got %f (diff: %e)", expectedInverse, rate, math.Abs(rate-expectedInverse))
	}

	// Test the same currency
	rate, err = e.GetRate("SGD", "SGD")

	if err != nil {
		t.Errorf("GetRate same currency failed: %v", err)
	}

	if rate != 1.0 {
		t.Errorf("Expected rate 1.0, got %f", rate)
	}

	// Test non-existent rate
	_, err = e.GetRate("SGD", "GBP")

	if !errors.Is(err, exception.ErrCurrencyConversionNotFound) {
		t.Errorf("Expected ErrCurrencyConversionNotFound, got %v", err)
	}

	// Test nil exchange
	var nilExchange *Exchange

	_, err = nilExchange.GetRate("SGD", "EUR")

	if !errors.Is(err, exception.ErrInvalidExchangeRate) {
		t.Errorf("Expected ErrInvalidExchangeRate for nil exchange, got %v", err)
	}
}

func TestConvertAmount(t *testing.T) {
	e := NewExchange()
	err := e.AddRate("SGD", "EUR", 0.85)

	if err != nil {
		t.Errorf("AddRate error: %v", err)
	}

	// 10.00 SGD (fraction 2) -> 8.50 EUR (fraction 2)
	// Amount 1000 -> 850
	amount, err := e.ConvertAmount(1000, "SGD", 2, "EUR", 2)

	if err != nil {
		t.Errorf("ConvertAmount failed: %v", err)
	}

	if amount != 850 {
		t.Errorf("Expected amount 850, got %d", amount)
	}

	// Test the same currency
	amount, err = e.ConvertAmount(1000, "SGD", 2, "SGD", 2)

	if err != nil {
		t.Errorf("ConvertAmount same currency failed: %v", err)
	}

	if amount != 1000 {
		t.Errorf("Expected amount 1000, got %d", amount)
	}

	// Test missing rate
	_, err = e.ConvertAmount(1000, "SGD", 2, "GBP", 2)

	if err == nil {
		t.Errorf("Expected ErrCurrencyConversionNotFound, got nil")
	} else if !errors.Is(err, exception.ErrCurrencyConversionNotFound) {
		t.Errorf("Expected ErrCurrencyConversionNotFound, got %v", err)
	}

	// Test nil exchange
	var nilExchange *Exchange

	_, err = nilExchange.ConvertAmount(1000, "SGD", 2, "EUR", 2)

	if !errors.Is(err, exception.ErrInvalidExchangeRate) {
		t.Errorf("Expected ErrInvalidExchangeRate for nil exchange, got %v", err)
	}

	// Test different fractions
	// 10 SGD (fraction 0) -> EUR (fraction 2)
	// 10 -> 8.50 EUR -> 850
	amount, err = e.ConvertAmount(10, "SGD", 0, "EUR", 2)

	if err != nil {
		t.Errorf("ConvertAmount diff fractions failed: %v", err)
	}

	if amount != 850 {
		t.Errorf("Expected amount 850, got %d", amount)
	}
}

func TestConvertAmountWithRate(t *testing.T) {
	e := NewExchange()

	// 10.00 (fraction 2) * 1.5 -> 15.00 (fraction 2)
	// 1000 -> 1500
	amount, err := e.ConvertAmountWithRate(1000, 2, 2, 1.5)

	if err != nil {
		t.Errorf("ConvertAmountWithRate failed: %v", err)
	}

	if amount != 1500 {
		t.Errorf("Expected amount 1500, got %d", amount)
	}

	// Invalid rate
	_, err = e.ConvertAmountWithRate(1000, 2, 2, -1.0)

	if !errors.Is(err, exception.ErrInvalidExchangeRate) {
		t.Errorf("Expected ErrInvalidExchangeRate, got %v", err)
	}
}

func TestConvertAmountDifferentFractions(t *testing.T) {
	e := NewExchange()

	// Setup exchange rates
	// USD to JPY rate: 110.0
	// USD to BHD rate: 0.377
	// JPY to BHD rate: 0.377/110 = 0.00342727...
	err := e.AddRate("USD", "JPY", 110.0)

	if err != nil {
		t.Fatalf("Failed to add USD/JPY rate: %v", err)
	}

	err = e.AddRate("USD", "BHD", 0.377)

	if err != nil {
		t.Fatalf("Failed to add USD/BHD rate: %v", err)
	}

	err = e.AddRate("JPY", "BHD", 0.377/110.0)

	if err != nil {
		t.Fatalf("Failed to add JPY/BHD rate: %v", err)
	}

	tests := []struct {
		name           string
		amount         int64
		fromCurrency   string
		fromFraction   int
		toCurrency     string
		toFraction     int
		expectedAmount int64
		description    string
	}{
		{
			name:           "USD to JPY - 2 decimals to 0 decimals",
			amount:         1000,
			fromCurrency:   "USD",
			fromFraction:   2,
			toCurrency:     "JPY",
			toFraction:     0,
			expectedAmount: 1100,
			description:    "USD ($10.00) stored as 1000 -> JPY (¥1100) stored as 1100",
		},
		{
			name:           "USD to BHD - 2 decimals to 3 decimals",
			amount:         1000,
			fromCurrency:   "USD",
			fromFraction:   2,
			toCurrency:     "BHD",
			toFraction:     3,
			expectedAmount: 3770,
			description:    "USD ($10.00) stored as 1000 -> BHD (BD 3.770) stored as 3770",
		},
		{
			name:           "JPY to USD - 0 decimals to 2 decimals",
			amount:         1000,
			fromCurrency:   "JPY",
			fromFraction:   0,
			toCurrency:     "USD",
			toFraction:     2,
			expectedAmount: 909,
			description:    "JPY (¥1000) stored as 1000 -> USD ($9.09) stored as 909",
		},
		{
			name:           "JPY to BHD - 0 decimals to 3 decimals",
			amount:         1000,
			fromCurrency:   "JPY",
			fromFraction:   0,
			toCurrency:     "BHD",
			toFraction:     3,
			expectedAmount: 3427,
			description:    "JPY (¥1000) stored as 1000 -> BHD (BD 3.427) stored as 3427",
		},
		{
			name:           "BHD to USD - 3 decimals to 2 decimals",
			amount:         1000,
			fromCurrency:   "BHD",
			fromFraction:   3,
			toCurrency:     "USD",
			toFraction:     2,
			expectedAmount: 265,
			description:    "BHD (BD 1.000) stored as 1000 -> USD ($2.65) stored as 265",
		},
		{
			name:           "BHD to JPY - 3 decimals to 0 decimals",
			amount:         1000,
			fromCurrency:   "BHD",
			fromFraction:   3,
			toCurrency:     "JPY",
			toFraction:     0,
			expectedAmount: 292,
			description:    "BHD (BD 1.000) stored as 1000 -> JPY (¥292) stored as 292",
		},
		{
			name:           "Same currency SGD",
			amount:         1000,
			fromCurrency:   "SGD",
			fromFraction:   2,
			toCurrency:     "SGD",
			toFraction:     2,
			expectedAmount: 1000,
			description:    "SGD (S$10.00) stored as 1000 -> SGD (S$10.00) stored as 1000",
		},
		{
			name:           "Same currency JPY",
			amount:         1000,
			fromCurrency:   "JPY",
			fromFraction:   0,
			toCurrency:     "JPY",
			toFraction:     0,
			expectedAmount: 1000,
			description:    "JPY (¥1000) stored as 1000 -> JPY (¥1000) stored as 1000",
		},
		{
			name:           "Same currency BHD",
			amount:         1000,
			fromCurrency:   "BHD",
			fromFraction:   3,
			toCurrency:     "BHD",
			toFraction:     3,
			expectedAmount: 1000,
			description:    "BHD (BD 1.000) stored as 1000 -> BHD (BD 1.000) stored as 1000",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			amount, err := e.ConvertAmount(tt.amount, tt.fromCurrency, tt.fromFraction, tt.toCurrency, tt.toFraction)

			if err != nil {
				t.Errorf("ConvertAmount failed: %v", err)
			}

			if amount != tt.expectedAmount {
				t.Errorf("%s\nExpected amount %d, got %d", tt.description, tt.expectedAmount, amount)
			}
		})
	}
}

// TestConcurrentAccess verifies that Exchange is safe for concurrent use
// by multiple goroutines performing reads and writes simultaneously.
func TestConcurrentAccess(t *testing.T) {
	e := NewExchange()

	// Pre-populate with some rates
	currencies := []string{"SGD", "EUR", "GBP", "JPY", "CAD"}

	for i, base := range currencies {
		for j, counter := range currencies {
			if i != j {
				if err := e.AddRate(base, counter, float64(i+1)/float64(j+1)); err != nil {
					t.Fatalf("Failed to add initial rate %s->%s: %v", base, counter, err)
				}
			}
		}
	}

	const (
		numWriters = 10
		numReaders = 50
		iterations = 1000
	)

	var wg sync.WaitGroup
	errChan := make(chan error, numWriters+numReaders)

	// Start concurrent writers
	for i := 0; i < numWriters; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < iterations; j++ {
				base := currencies[j%len(currencies)]
				counter := currencies[(j+1)%len(currencies)]
				rate := float64(id+1) / float64(j+1)

				if err := e.AddRate(base, counter, rate); err != nil {
					errChan <- err

					return
				}
			}
		}(i)
	}

	// Start concurrent readers
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < iterations; j++ {
				base := currencies[j%len(currencies)]
				counter := currencies[(j+1)%len(currencies)]

				_, err := e.GetRate(base, counter)

				if err != nil && !errors.Is(err, exception.ErrCurrencyConversionNotFound) {
					errChan <- err

					return
				}

				// Also test conversion
				_, err = e.ConvertAmount(1000, base, 2, counter, 2)

				if err != nil && !errors.Is(err, exception.ErrCurrencyConversionNotFound) {
					errChan <- err

					return
				}
			}
		}(i)
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		t.Errorf("Concurrent operation failed: %v", err)
	}
}

// TestConcurrentSameInnerMap tests the specific case where multiple goroutines
// access the same inner map (same base currency) to ensure the nested map
// structure is properly protected.
func TestConcurrentSameInnerMap(t *testing.T) {
	e := NewExchange()

	const (
		numGoroutines = 50
		iterations    = 1000
		baseCurrency  = "USD"
	)

	currencies := []string{"EUR", "GBP", "JPY", "CAD", "AUD", "CHF", "CNY"}

	var wg sync.WaitGroup
	errChan := make(chan error, numGoroutines*2)

	// Half writers, half readers - all targeting the same base currency's inner map
	for i := 0; i < numGoroutines/2; i++ {
		// Writers
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < iterations; j++ {
				counter := currencies[j%len(currencies)]
				rate := float64(id+1) / float64(j+1)

				if err := e.AddRate(baseCurrency, counter, rate); err != nil {
					errChan <- err

					return
				}
			}
		}(i)

		// Readers
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			for j := 0; j < iterations; j++ {
				counter := currencies[j%len(currencies)]

				_, err := e.GetRate(baseCurrency, counter)

				if err != nil && !errors.Is(err, exception.ErrCurrencyConversionNotFound) {
					errChan <- err

					return
				}
			}
		}(i)
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		t.Errorf("Concurrent same-inner-map operation failed: %v", err)
	}
}

func TestConverterGetExchange_CoveragePaths(t *testing.T) {
	t.Run("nil receiver", func(t *testing.T) {
		var c *Converter
		_, err := c.GetExchange()

		if !errors.Is(err, exception.ErrNoConverterProvided) {
			t.Fatalf("GetExchange() error = %v, want ErrNoConverterProvided", err)
		}
	})

	t.Run("nil exchange", func(t *testing.T) {
		c := &Converter{exchange: nil}
		_, err := c.GetExchange()

		if !errors.Is(err, exception.ErrInvalidExchangeRate) {
			t.Fatalf("GetExchange() error = %v, want ErrInvalidExchangeRate", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		ex := &Exchange{}
		c, err := NewConverter(ex)

		if err != nil {
			t.Fatalf("NewConverter() unexpected error: %v", err)
		}

		got, err := c.GetExchange()

		if err != nil {
			t.Fatalf("GetExchange() unexpected error: %v", err)
		}

		if got != ex {
			t.Fatal("GetExchange() did not return the same exchange pointer")
		}
	})
}

func TestExchange_ConvertAmount_RealWorldCommonCurrencies(t *testing.T) {
	t.Parallel()

	// Commonly used global currencies for real-world scenarios.
	const (
		USD = "USD"
		EUR = "EUR"
		GBP = "GBP"
		JPY = "JPY"
		CNY = "CNY"
		CAD = "CAD"
		AUD = "AUD"
		CHF = "CHF"
	)

	ex := NewExchange()
	// Rates are fixed test fixtures (not market data). Use decimal strings so the test's
	// expected values can be computed deterministically.
	testutil.TestRequireNoErr(t, ex.AddRate(USD, EUR, testutil.TestMustParseFloatMoney(t, "0.92")))
	testutil.TestRequireNoErr(t, ex.AddRate(USD, GBP, testutil.TestMustParseFloatMoney(t, "0.79")))
	testutil.TestRequireNoErr(t, ex.AddRate(USD, JPY, testutil.TestMustParseFloatMoney(t, "150")))
	testutil.TestRequireNoErr(t, ex.AddRate(USD, CNY, testutil.TestMustParseFloatMoney(t, "7.2")))
	testutil.TestRequireNoErr(t, ex.AddRate(USD, CAD, testutil.TestMustParseFloatMoney(t, "1.36")))
	testutil.TestRequireNoErr(t, ex.AddRate(USD, AUD, testutil.TestMustParseFloatMoney(t, "1.52")))
	testutil.TestRequireNoErr(t, ex.AddRate(USD, CHF, testutil.TestMustParseFloatMoney(t, "0.88")))

	tests := []struct {
		name         string
		amount       int64
		fromCode     string
		fromFraction int
		toCode       string
		toFraction   int
		rate         string
		wantTol      int64
	}{
		{
			name:         "Invoice total USD to EUR",
			amount:       123456, // $1,234.56
			fromCode:     USD,
			fromFraction: 2,
			toCode:       EUR,
			toFraction:   2,
			rate:         "0.92",
			wantTol:      0,
		},
		{
			name:         "Salary USD to GBP",
			amount:       250000, // $2,500.00
			fromCode:     USD,
			fromFraction: 2,
			toCode:       GBP,
			toFraction:   2,
			rate:         "0.79",
			wantTol:      0,
		},
		{
			name:         "Card purchase USD to JPY (0-decimal)",
			amount:       1999, // $19.99
			fromCode:     USD,
			fromFraction: 2,
			toCode:       JPY,
			toFraction:   0,
			rate:         "150",
			// This case lands exactly on an x.5 boundary in exact arithmetic; since the implementation
			// uses float64, allow +/- 1 JPY to keep the test stable while still bounding error.
			wantTol: 1,
		},
		{
			name:         "Refund negative USD to EUR stays negative",
			amount:       -1099, // -$10.99
			fromCode:     USD,
			fromFraction: 2,
			toCode:       EUR,
			toFraction:   2,
			rate:         "0.92",
			wantTol:      0,
		},
		{
			name:         "Round amount USD to CNY",
			amount:       5000, // $50.00
			fromCode:     USD,
			fromFraction: 2,
			toCode:       CNY,
			toFraction:   2,
			rate:         "7.2",
			wantTol:      0,
		},
		{
			name:         "Round amount USD to CAD",
			amount:       10000, // $100.00
			fromCode:     USD,
			fromFraction: 2,
			toCode:       CAD,
			toFraction:   2,
			rate:         "1.36",
			wantTol:      0,
		},
		{
			name:         "Round amount USD to AUD",
			amount:       10000, // $100.00
			fromCode:     USD,
			fromFraction: 2,
			toCode:       AUD,
			toFraction:   2,
			rate:         "1.52",
			wantTol:      0,
		},
		{
			name:         "Round amount USD to CHF",
			amount:       10000, // $100.00
			fromCode:     USD,
			fromFraction: 2,
			toCode:       CHF,
			toFraction:   2,
			rate:         "0.88",
			wantTol:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := ex.ConvertAmount(tt.amount, tt.fromCode, tt.fromFraction, tt.toCode, tt.toFraction)

			if err != nil {
				t.Fatalf("ConvertAmount() unexpected error: %v", err)
			}

			want := testutil.TestExpectedConvertAmountMoney(t, tt.amount, tt.fromFraction, tt.toFraction, tt.rate)

			if diff := testutil.TestAbs64Money(t, got-want); diff > tt.wantTol {
				t.Fatalf("ConvertAmount() = %d, want %d (tol=%d; diff=%d)", got, want, tt.wantTol, diff)
			}
		})
	}
}

func TestExchange_ConvertAmount_UsesInverseRatePath_RoundTripBounded(t *testing.T) {
	t.Parallel()

	ex := NewExchange()
	testutil.TestRequireNoErr(t, ex.AddRate("USD", "EUR", testutil.TestMustParseFloatMoney(t, "0.92")))

	eurCents, err := ex.ConvertAmount(10000, "USD", 2, "EUR", 2)

	if err != nil {
		t.Fatalf("ConvertAmount(USD->EUR) unexpected error: %v", err)
	}

	usdCents, err := ex.ConvertAmount(eurCents, "EUR", 2, "USD", 2)

	if err != nil {
		t.Fatalf("ConvertAmount(EUR->USD) unexpected error: %v", err)
	}

	if diff := testutil.TestAbs64Money(t, usdCents-10000); diff > 1 {
		t.Fatalf("round-trip USD cents = %d, want 10000 (diff=%d, want <= 1)", usdCents, diff)
	}
}

func TestExchange_ConvertAmount_DoesNotAssumeChainedCrossRates(t *testing.T) {
	t.Parallel()

	ex := NewExchange()
	testutil.TestRequireNoErr(t, ex.AddRate("USD", "EUR", testutil.TestMustParseFloatMoney(t, "0.92")))
	testutil.TestRequireNoErr(t, ex.AddRate("USD", "JPY", testutil.TestMustParseFloatMoney(t, "150")))

	_, err := ex.ConvertAmount(10000, "EUR", 2, "JPY", 0)

	if !errors.Is(err, exception.ErrCurrencyConversionNotFound) {
		t.Fatalf("ConvertAmount(EUR->JPY) error = %v, want %v", err, exception.ErrCurrencyConversionNotFound)
	}
}
