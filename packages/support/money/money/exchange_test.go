package money

import (
	"errors"
	"strings"
	"testing"

	"github.com/gollin/packages/support/money/currency"
	"github.com/gollin/packages/support/money/exception"
	"github.com/gollin/packages/support/money/exchange"
	testutil "github.com/gollin/packages/support/money/tests"
)

func TestConvert(t *testing.T) {
	currencies := currency.NewManager()
	ex := exchange.NewExchange()
	err := ex.AddRate(currency.USD, currency.EUR, 0.85)

	if err != nil {
		t.Fatalf("AddRate() unexpected error: %v", err)
	}

	converter := newTestConverter(t, currencies, ex)

	usd := NewManager().Create(500, currency.USD)

	// Test same-currency conversion
	same, err := converter.Convert(usd, currency.USD)

	if err != nil {
		t.Fatalf("Convert() unexpected error: %v", err)
	}

	if testutil.TestRequire(t, same.Amount) != testutil.TestRequire(t, usd.Amount) || testutil.TestRequire(t, same.Currency).Code != currency.USD {
		t.Fatalf("Convert() same currency = (%d, %s), want (%d, USD)", testutil.TestRequire(t, same.Amount), testutil.TestRequire(t, same.Currency).Code, testutil.TestRequire(t, usd.Amount))
	}

	// Test cross-currency conversion with rate
	eur, err := converter.Convert(usd, currency.EUR)
	if err != nil {
		t.Fatalf("Convert() unexpected error: %v", err)
	}

	expectedAmount := int64(425) // 500 * 0.85 = 425
	if testutil.TestRequire(t, eur.Amount) != expectedAmount || testutil.TestRequire(t, eur.Currency).Code != currency.EUR {
		t.Fatalf("Convert() = (%d, %s), want (%d, EUR)", testutil.TestRequire(t, eur.Amount), testutil.TestRequire(t, eur.Currency).Code, expectedAmount)
	}

	// Test conversion without rate (should fail)
	if _, err := converter.Convert(usd, currency.GBP); !errors.Is(err, exception.ErrCurrencyConversionNotFound) {
		t.Fatalf("Convert() error = %v, want %v", err, exception.ErrCurrencyConversionNotFound)
	}
}

func TestConvertNilMoney(t *testing.T) {
	currencies := currency.NewManager()
	ex := exchange.NewExchange()
	err := ex.AddRate(currency.USD, currency.EUR, 0.85)
	if err != nil {
		t.Fatalf("AddRate() unexpected error: %v", err)
	}

	converter := newTestConverter(t, currencies, ex)

	if _, err := converter.Convert(nil, currency.EUR); !errors.Is(err, exception.ErrNoMoneyProvided) {
		t.Fatalf("Convert(nil) error = %v, want %v", err, exception.ErrNoMoneyProvided)
	}
}

func TestConvertWithRate(t *testing.T) {
	currencies := currency.NewManager()
	ex := exchange.NewExchange()
	// Add a rate to make the exchange valid
	err := ex.AddRate(currency.USD, currency.EUR, 0.85)

	if err != nil {
		t.Fatalf("AddRate() unexpected error: %v", err)
	}

	converter := newTestConverter(t, currencies, ex)

	usd := NewManager().Create(100, currency.USD)

	converted, err := converter.ConvertWithRate(usd, currency.EUR, 2.0)

	if err != nil {
		t.Fatalf("ConvertWithRate() unexpected error: %v", err)
	}

	if testutil.TestRequire(t, converted.Amount) != 200 || testutil.TestRequire(t, converted.Currency).Code != currency.EUR {
		t.Fatalf("ConvertWithRate() = (%d, %s), want (200, EUR)", testutil.TestRequire(t, converted.Amount), testutil.TestRequire(t, converted.Currency).Code)
	}

	if _, err := converter.ConvertWithRate(usd, currency.EUR, 0); !errors.Is(err, exception.ErrInvalidExchangeRate) {
		t.Fatalf("ConvertWithRate() error = %v, want %v", err, exception.ErrInvalidExchangeRate)
	}
}

func TestConvertWithRateNilMoney(t *testing.T) {
	currencies := currency.NewManager()
	ex := exchange.NewExchange()
	err := ex.AddRate(currency.USD, currency.EUR, 0.85)
	if err != nil {
		t.Fatalf("AddRate() unexpected error: %v", err)
	}

	converter := newTestConverter(t, currencies, ex)

	if _, err := converter.ConvertWithRate(nil, currency.EUR, 1.0); !errors.Is(err, exception.ErrNoMoneyProvided) {
		t.Fatalf("ConvertWithRate(nil) error = %v, want %v", err, exception.ErrNoMoneyProvided)
	}
}

func TestNewConverterNilExchange(t *testing.T) {
	currencies := currency.NewManager()

	_, err := NewConverter(currencies, nil)
	if !errors.Is(err, exception.ErrInvalidExchangeRate) {
		t.Fatalf("NewConverter(nil) error = %v, want %v", err, exception.ErrInvalidExchangeRate)
	}
}

func TestConverterNilCheck(t *testing.T) {
	// Create a converter with nil exchange by bypassing NewConverter (for testing only)
	var converter *Converter

	usd := NewManager().Create(100, currency.USD)

	// Test Convert with nil converter
	_, err := converter.Convert(usd, currency.EUR)
	if !errors.Is(err, exception.ErrInvalidExchangeRate) {
		t.Fatalf("Convert() with nil converter error = %v, want %v", err, exception.ErrInvalidExchangeRate)
	}

	// Test ConvertWithRate with nil converter
	_, err = converter.ConvertWithRate(usd, currency.EUR, 0.85)
	if !errors.Is(err, exception.ErrInvalidExchangeRate) {
		t.Fatalf("ConvertWithRate() with nil converter error = %v, want %v", err, exception.ErrInvalidExchangeRate)
	}
}

func TestNewConverterInvalidExchange(t *testing.T) {
	currencies := currency.NewManager()
	// Test with Exchange that has nil rates map
	invalidEx := &exchange.Exchange{}

	_, err := NewConverter(currencies, invalidEx)
	if !errors.Is(err, exception.ErrInvalidExchangeRate) {
		t.Fatalf("NewConverter(invalid exchange with nil rates) error = %v, want %v", err, exception.ErrInvalidExchangeRate)
	}

	// Test with Exchange that has empty rates map
	emptyEx := exchange.NewExchange()

	_, err = NewConverter(currencies, emptyEx)
	if !errors.Is(err, exception.ErrInvalidExchangeRate) {
		t.Fatalf("NewConverter(invalid exchange with empty rates) error = %v, want %v", err, exception.ErrInvalidExchangeRate)
	}
}

func TestRoundWithZeroFractionCurrency(t *testing.T) {
	jpy := NewManager().Create(155, currency.JPY)
	manager := NewManager()
	rounded, err := manager.Round(jpy)
	testutil.TestRequireNoErr(t, err)

	if testutil.TestRequire(t, rounded.Amount) != 155 {
		t.Fatalf("Round() JPY amount = %d, want 155", testutil.TestRequire(t, rounded.Amount))
	}
}

func TestMultiplyWithNegativeValues(t *testing.T) {
	m := NewManager().Create(10, currency.USD)
	manager := NewManager()
	result, err := manager.Multiply(m, -1, 2, -3) // 10 * -1 * 2 * -3 = 60

	if err != nil {
		t.Fatalf("Multiply() unexpected error: %v", err)
	}

	if testutil.TestRequire(t, result.Amount) != 60 {
		t.Fatalf("Multiply() with negatives = %d, want 60", testutil.TestRequire(t, result.Amount))
	}

	if testutil.TestRequire(t, result.Currency).Code != currency.USD {
		t.Fatalf("Multiply() currency = %s, want USD", testutil.TestRequire(t, result.Currency).Code)
	}
}

func TestSplitSinglePart(t *testing.T) {
	m := NewManager().Create(999, currency.USD)
	manager := NewManager()
	parts, err := manager.Split(m, 1)
	if err != nil {
		t.Fatalf("Split() unexpected error: %v", err)
	}

	if len(parts) != 1 || testutil.TestRequire(t, parts[0].Amount) != 999 {
		t.Fatalf("Split() single part = (%d, len=%d), want (999, len=1)", testutil.TestRequire(t, parts[0].Amount), len(parts))
	}
}

func TestAllocateLeftoverDistributionOrder(t *testing.T) {
	m := NewManager().Create(2, currency.USD)
	manager := NewManager()
	parts, err := manager.Allocate(m, 1, 3)

	if err != nil {
		t.Fatalf("Allocate() unexpected error: %v", err)
	}

	if testutil.TestRequire(t, parts[0].Amount) != 1 || testutil.TestRequire(t, parts[1].Amount) != 1 {
		t.Fatalf("Allocate() amounts = (%d, %d), want (1, 1)", testutil.TestRequire(t, parts[0].Amount), testutil.TestRequire(t, parts[1].Amount))
	}
}

func TestConvertWithRateFractionLoss(t *testing.T) {
	currencies := currency.NewManager()
	ex := exchange.NewExchange()
	// Add a dummy rate to make the exchange valid
	err := ex.AddRate(currency.USD, currency.EUR, 1.0)

	if err != nil {
		t.Fatalf("AddRate() unexpected error: %v", err)
	}

	converter := newTestConverter(t, currencies, ex)

	m := NewManager().Create(123, currency.USD) // 1.23 major units (123 cents)
	result, err := converter.ConvertWithRate(m, currency.BHD, 1.5)

	if err != nil {
		t.Fatalf("ConvertWithRate() unexpected error: %v", err)
	}

	// 1.23 USD * 1.5 = 1.845 BHD -> 1845 fils (BHD minor units with 3 decimals)
	if testutil.TestRequire(t, result.Amount) != 1845 {
		t.Fatalf("ConvertWithRate() amount = %d, want 1845", testutil.TestRequire(t, result.Amount))
	}

	if testutil.TestRequire(t, result.Currency).Code != currency.BHD {
		t.Fatalf("ConvertWithRate() currency = %s, want BHD", testutil.TestRequire(t, result.Currency).Code)
	}
}

func TestDBValueCustomSeparator(t *testing.T) {
	m := NewManager().Create(500, currency.EUR)
	sep := GetDBMoneyValueSeparator()

	val, err := m.Value()
	if err != nil {
		t.Fatalf("Value() unexpected error: %v", err)
	}

	expected := "500" + sep + "EUR"
	if val != expected {
		t.Fatalf("Value() with separator '%s' = %s, want %s", sep, val, expected)
	}
}

// TestConvertJPYToHigherPrecision tests converting from JPY (0 decimals) to currencies with higher precision
func TestConvertJPYToHigherPrecision(t *testing.T) {
	currencies := currency.NewManager()
	ex := exchange.NewExchange()
	// Add a dummy rate to make the exchange valid
	err := ex.AddRate(currency.USD, currency.EUR, 1.0)

	if err != nil {
		t.Fatalf("AddRate() unexpected error: %v", err)
	}

	converter := newTestConverter(t, currencies, ex)

	tests := []struct {
		name           string
		amount         int64
		toCurrency     string
		rate           float64
		expectedAmount int64
		description    string
	}{
		{
			name:           "JPY to USD with rounding down",
			amount:         100,
			toCurrency:     currency.USD,
			rate:           0.0091,
			expectedAmount: 91, // 100 * 0.0091 = 0.91 USD = 91 cents
			description:    "100 JPY at rate 0.0091 = 0.91 USD (91 cents)",
		},
		{
			name:           "JPY to USD with rounding up",
			amount:         100,
			toCurrency:     currency.USD,
			rate:           0.0095,
			expectedAmount: 95, // 100 * 0.0095 = 0.95 USD = 95 cents
			description:    "100 JPY at rate 0.0095 = 0.95 USD (95 cents)",
		},
		{
			name:           "JPY to BHD with high precision",
			amount:         1000,
			toCurrency:     currency.BHD,
			rate:           0.00343,
			expectedAmount: 3430, // 1000 * 0.00343 = 3.43 BHD = 3430 fils
			description:    "1000 JPY at rate 0.00343 = 3.43 BHD (3430 fils)",
		},
		{
			name:           "Small JPY amount rounds down to zero",
			amount:         1,
			toCurrency:     currency.USD,
			rate:           0.0001,
			expectedAmount: 0, // 1 * 0.0001 = 0.0001 USD = 0.01 cents, rounds to 0
			description:    "1 JPY at very small rate rounds to 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager().Create(tt.amount, currency.JPY)
			result, err := converter.ConvertWithRate(m, tt.toCurrency, tt.rate)

			if err != nil {
				t.Fatalf("ConvertWithRate() unexpected error: %v", err)
			}

			if testutil.TestRequire(t, result.Amount) != tt.expectedAmount {
				t.Errorf("%s: got %d, want %d", tt.description, testutil.TestRequire(t, result.Amount), tt.expectedAmount)
			}

			if testutil.TestRequire(t, result.Currency).Code != tt.toCurrency {
				t.Errorf("Currency mismatch: got %s, want %s", testutil.TestRequire(t, result.Currency).Code, tt.toCurrency)
			}
		})
	}
}

// TestConvertToJPYFromHigherPrecision tests converting to JPY (0 decimals) from currencies with higher precision
func TestConvertToJPYFromHigherPrecision(t *testing.T) {
	currencies := currency.NewManager()
	ex := exchange.NewExchange()
	// Add a dummy rate to make the exchange valid
	err := ex.AddRate(currency.USD, currency.EUR, 1.0)

	if err != nil {
		t.Fatalf("AddRate() unexpected error: %v", err)
	}

	converter := newTestConverter(t, currencies, ex)

	tests := []struct {
		name           string
		amount         int64
		fromCurrency   string
		rate           float64
		expectedAmount int64
		description    string
	}{
		{
			name:           "USD to JPY with rounding down",
			amount:         100, // $1.00
			fromCurrency:   currency.USD,
			rate:           149.3,
			expectedAmount: 149, // 1.00 * 149.3 = 149.3, rounds to 149
			description:    "$1.00 at rate 149.3 = 149.3 JPY, rounds to 149",
		},
		{
			name:           "USD to JPY with rounding up",
			amount:         100, // $1.00
			fromCurrency:   currency.USD,
			rate:           149.7,
			expectedAmount: 150, // 1.00 * 149.7 = 149.7, rounds to 150
			description:    "$1.00 at rate 149.7 = 149.7 JPY, rounds to 150",
		},
		{
			name:           "USD to JPY exact midpoint rounds up",
			amount:         100, // $1.00
			fromCurrency:   currency.USD,
			rate:           149.5,
			expectedAmount: 150, // 1.00 * 149.5 = 149.5, rounds to 150
			description:    "$1.00 at rate 149.5 = 149.5 JPY, rounds to 150 (banker's rounding)",
		},
		{
			name:           "BHD to JPY high precision to zero decimals",
			amount:         1000, // 1.000 BHD
			fromCurrency:   currency.BHD,
			rate:           291.545,
			expectedAmount: 292, // 1.000 * 291.545 = 291.545, rounds to 292
			description:    "1.000 BHD at rate 291.545 = 291.545 JPY, rounds to 292",
		},
		{
			name:           "Small USD cents to JPY rounds to zero",
			amount:         1, // $0.01
			fromCurrency:   currency.USD,
			rate:           0.3,
			expectedAmount: 0, // 0.01 * 0.3 = 0.003 JPY, rounds to 0
			description:    "$0.01 at rate 0.3 = 0.003 JPY, rounds to 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager().Create(tt.amount, tt.fromCurrency)
			result, err := converter.ConvertWithRate(m, currency.JPY, tt.rate)

			if err != nil {
				t.Fatalf("ConvertWithRate() unexpected error: %v", err)
			}

			if testutil.TestRequire(t, result.Amount) != tt.expectedAmount {
				t.Errorf("%s: got %d, want %d", tt.description, testutil.TestRequire(t, result.Amount), tt.expectedAmount)
			}

			if testutil.TestRequire(t, result.Currency).Code != currency.JPY {
				t.Errorf("Currency mismatch: got %s, want JPY", testutil.TestRequire(t, result.Currency).Code)
			}
		})
	}
}

// TestConvertExtremeAmounts tests conversion with very large and very small amounts
func TestConvertExtremeAmounts(t *testing.T) {
	currencies := currency.NewManager()
	ex := exchange.NewExchange()
	// Add a dummy rate to make the exchange valid
	err := ex.AddRate(currency.USD, currency.EUR, 1.0)

	if err != nil {
		t.Fatalf("AddRate() unexpected error: %v", err)
	}

	converter := newTestConverter(t, currencies, ex)

	tests := []struct {
		name           string
		amount         int64
		fromCurrency   string
		toCurrency     string
		rate           float64
		expectedAmount int64
		description    string
	}{
		{
			name:           "Very large USD amount",
			amount:         999999999, // $9,999,999.99
			fromCurrency:   currency.USD,
			toCurrency:     currency.JPY,
			rate:           150.0,
			expectedAmount: 1499999999, // ~15 billion JPY
			description:    "Large USD amount converts correctly to JPY",
		},
		{
			name:           "Very small rate produces zero",
			amount:         100, // $1.00
			fromCurrency:   currency.USD,
			toCurrency:     currency.JPY,
			rate:           0.001,
			expectedAmount: 0, // 1.00 * 0.001 = 0.001, rounds to 0
			description:    "Very small rate rounds to 0",
		},
		{
			name:           "Very large rate",
			amount:         100, // $1.00
			fromCurrency:   currency.USD,
			toCurrency:     currency.JPY,
			rate:           10000.0,
			expectedAmount: 10000, // 1.00 * 10000 = 10000 JPY
			description:    "Very large rate converts correctly",
		},
		{
			name:           "Precision preservation with BHD",
			amount:         12345, // $123.45
			fromCurrency:   currency.USD,
			toCurrency:     currency.BHD,
			rate:           0.377,
			expectedAmount: 46541, // 123.45 * 0.377 = 46.54065 BHD = 46540.65 fils, rounds to 46541
			description:    "High precision preserved in BHD conversion",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager().Create(tt.amount, tt.fromCurrency)
			result, err := converter.ConvertWithRate(m, tt.toCurrency, tt.rate)

			if err != nil {
				t.Fatalf("ConvertWithRate() unexpected error: %v", err)
			}

			if testutil.TestRequire(t, result.Amount) != tt.expectedAmount {
				t.Errorf("%s: got %d, want %d", tt.description, testutil.TestRequire(t, result.Amount), tt.expectedAmount)
			}
		})
	}
}

// TestConvertRoundingBehavior tests that rounding works consistently across different scenarios
func TestConvertRoundingBehavior(t *testing.T) {
	currencies := currency.NewManager()
	ex := exchange.NewExchange()
	// Add a dummy rate to make the exchange valid
	err := ex.AddRate(currency.USD, currency.EUR, 1.0)

	if err != nil {
		t.Fatalf("AddRate() unexpected error: %v", err)
	}

	converter := newTestConverter(t, currencies, ex)

	tests := []struct {
		name           string
		amount         int64
		fromCurrency   string
		toCurrency     string
		rate           float64
		expectedAmount int64
		description    string
	}{
		{
			name:           "Round down at 0.4",
			amount:         100,
			fromCurrency:   currency.USD,
			toCurrency:     currency.JPY,
			rate:           10.4,
			expectedAmount: 10, // 1.00 * 10.4 = 10.4, rounds to 10
			description:    "0.4 rounds down to nearest integer",
		},
		{
			name:           "Round up at 0.6",
			amount:         100,
			fromCurrency:   currency.USD,
			toCurrency:     currency.JPY,
			rate:           10.6,
			expectedAmount: 11, // 1.00 * 10.6 = 10.6, rounds to 11
			description:    "0.6 rounds up to nearest integer",
		},
		{
			name:           "Round at exact 0.5",
			amount:         100,
			fromCurrency:   currency.USD,
			toCurrency:     currency.JPY,
			rate:           10.5,
			expectedAmount: 11, // 1.00 * 10.5 = 10.5, rounds to 11 (Go's math.Round uses half-away-from-zero rounding)
			description:    "0.5 rounds using banker's rounding",
		},
		{
			name:           "Multiple decimals round correctly",
			amount:         12345, // $123.45
			fromCurrency:   currency.USD,
			toCurrency:     currency.JPY,
			rate:           149.999,
			expectedAmount: 18517, // 123.45 * 149.999 = 18517.38255, rounds to 18517
			description:    "Complex decimal rounds correctly",
		},
		{
			name:           "Fractional cent preservation in BHD",
			amount:         333, // $3.33
			fromCurrency:   currency.USD,
			toCurrency:     currency.BHD,
			rate:           0.377,
			expectedAmount: 1255, // 3.33 * 0.377 = 1.25541 BHD = 1255.41 fils, rounds to 1255
			description:    "BHD preserves fractional precision",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager().Create(tt.amount, tt.fromCurrency)
			result, err := converter.ConvertWithRate(m, tt.toCurrency, tt.rate)

			if err != nil {
				t.Fatalf("ConvertWithRate() unexpected error: %v", err)
			}

			if testutil.TestRequire(t, result.Amount) != tt.expectedAmount {
				t.Errorf("%s: got %d, want %d (calculation shows expected behavior)",
					tt.description, testutil.TestRequire(t, result.Amount), tt.expectedAmount)
			}
		})
	}
}

func newTestConverter(t *testing.T, currencies *currency.Manager, ex *exchange.Exchange) *Converter {
	t.Helper()

	converter, err := NewConverter(currencies, ex)
	if err != nil {
		t.Fatalf("AddRate() unexpected error: %v", err)
	}

	return converter
}

func TestConverter_Convert_Coverage(t *testing.T) {
	t.Parallel()

	currencies := currency.NewManager()
	ex := exchange.NewExchange()
	_ = ex.AddRate(currency.USD, currency.EUR, 0.85)
	converter := newTestConverter(t, currencies, ex)
	usd := NewManager().Create(100, currency.USD)

	// Test invalid toCurrency
	_, err := converter.Convert(usd, "INVALID")
	if err == nil {
		t.Fatal("Convert() with invalid currency code should return error")
	}
	if !errors.Is(err, exception.ErrCurrencyNotFound) {
		t.Errorf("Convert() error = %v, want %v", err, exception.ErrCurrencyNotFound)
	}
	if !strings.Contains(err.Error(), "INVALID") {
		t.Errorf("Convert() error message %q does not contain currency code", err.Error())
	}

	// Test ConvertWithRate invalid toCurrency
	_, err = converter.ConvertWithRate(usd, "INVALID", 1.0)
	if err == nil {
		t.Fatal("ConvertWithRate() with invalid currency code should return error")
	}
	if !errors.Is(err, exception.ErrCurrencyNotFound) {
		t.Errorf("ConvertWithRate() error = %v, want %v", err, exception.ErrCurrencyNotFound)
	}
	if !strings.Contains(err.Error(), "INVALID") {
		t.Errorf("ConvertWithRate() error message %q does not contain currency code", err.Error())
	}

	// Test Convert with Money having nil currency
	badMoney := &Money{amount: 100, currency: nil}
	_, err = converter.Convert(badMoney, currency.EUR)
	if err == nil {
		t.Fatal("Convert() with nil currency money should return error")
	}
	if !errors.Is(err, exception.ErrNoMoneyProvided) && err.Error() != "money instance has no currency" {
		t.Errorf("Convert() with nil currency money error = %v", err)
	}
}

func TestConverter_RealWorld_CommonCurrencies(t *testing.T) {
	t.Parallel()

	cm := currency.NewManager()
	ex := exchange.NewExchange()

	// Fixed fixture rates (not market data).
	testutil.TestRequireNoErr(t, ex.AddRate(currency.USD, currency.EUR, testutil.TestMustParseFloatMoney(t, "0.92")))
	testutil.TestRequireNoErr(t, ex.AddRate(currency.USD, currency.GBP, testutil.TestMustParseFloatMoney(t, "0.79")))
	testutil.TestRequireNoErr(t, ex.AddRate(currency.USD, currency.JPY, testutil.TestMustParseFloatMoney(t, "150")))
	testutil.TestRequireNoErr(t, ex.AddRate(currency.USD, currency.CNY, testutil.TestMustParseFloatMoney(t, "7.2")))
	testutil.TestRequireNoErr(t, ex.AddRate(currency.USD, currency.CAD, testutil.TestMustParseFloatMoney(t, "1.36")))
	testutil.TestRequireNoErr(t, ex.AddRate(currency.USD, currency.AUD, testutil.TestMustParseFloatMoney(t, "1.52")))
	testutil.TestRequireNoErr(t, ex.AddRate(currency.USD, currency.CHF, testutil.TestMustParseFloatMoney(t, "0.88")))

	converter := newTestConverter(t, cm, ex)
	mm := NewManager()

	tests := []struct {
		name         string
		from         *Money
		toCurrency   string
		rate         string
		wantTol      int64
		wantCode     string
		fromFraction int
		toFraction   int
	}{
		{
			name:         "Invoice total USD to EUR",
			from:         mm.Create(123456, currency.USD), // $1,234.56
			toCurrency:   currency.EUR,
			rate:         "0.92",
			wantTol:      0,
			wantCode:     currency.EUR,
			fromFraction: 2,
			toFraction:   2,
		},
		{
			name:         "Card purchase USD to JPY",
			from:         mm.Create(1999, currency.USD), // $19.99
			toCurrency:   currency.JPY,
			rate:         "150",
			wantTol:      1,
			wantCode:     currency.JPY,
			fromFraction: 2,
			toFraction:   0,
		},
		{
			name:         "Salary USD to GBP",
			from:         mm.Create(250000, currency.USD), // $2,500.00
			toCurrency:   currency.GBP,
			rate:         "0.79",
			wantTol:      0,
			wantCode:     currency.GBP,
			fromFraction: 2,
			toFraction:   2,
		},
		{
			name:         "Refund negative USD to EUR stays negative",
			from:         mm.Create(-1099, currency.USD), // -$10.99
			toCurrency:   currency.EUR,
			rate:         "0.92",
			wantTol:      0,
			wantCode:     currency.EUR,
			fromFraction: 2,
			toFraction:   2,
		},
		{
			name:         "JPY to USD via inverse rate",
			from:         mm.Create(10000, currency.JPY), // ¥10,000
			toCurrency:   currency.USD,
			rate:         "1/150",
			wantTol:      1,
			wantCode:     currency.USD,
			fromFraction: 0,
			toFraction:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := converter.Convert(tt.from, tt.toCurrency)
			if err != nil {
				t.Fatalf("Convert() unexpected error: %v", err)
			}

			want := testutil.TestExpectedConvertAmountMoney(t, testutil.TestRequire(t, tt.from.Amount), tt.fromFraction, tt.toFraction, tt.rate)
			if diff := testutil.TestAbs64Money(t, testutil.TestRequire(t, got.Amount)-want); diff > tt.wantTol {
				t.Fatalf("Convert() amount = %d, want %d (tol=%d; diff=%d)", testutil.TestRequire(t, got.Amount), want, tt.wantTol, diff)
			}

			if testutil.TestRequire(t, got.Currency).Code != tt.wantCode {
				t.Fatalf("Convert() currency = %s, want %s", testutil.TestRequire(t, got.Currency).Code, tt.wantCode)
			}
		})
	}
}

func TestConverter_RealWorld_NoImplicitCrossRate(t *testing.T) {
	t.Parallel()

	cm := currency.NewManager()
	ex := exchange.NewExchange()
	testutil.TestRequireNoErr(t, ex.AddRate(currency.USD, currency.EUR, testutil.TestMustParseFloatMoney(t, "0.92")))
	testutil.TestRequireNoErr(t, ex.AddRate(currency.USD, currency.JPY, testutil.TestMustParseFloatMoney(t, "150")))

	converter := newTestConverter(t, cm, ex)
	eur := NewManager().Create(10000, currency.EUR) // €100.00

	_, err := converter.Convert(eur, currency.JPY)
	if !errors.Is(err, exception.ErrCurrencyConversionNotFound) {
		t.Fatalf("Convert(EUR->JPY) error = %v, want %v", err, exception.ErrCurrencyConversionNotFound)
	}
}
