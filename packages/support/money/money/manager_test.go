package money

import (
	"errors"
	"math"
	"sync"
	"testing"

	"github.com/gollin/packages/support/money/currency"
	"github.com/gollin/packages/support/money/exception"
	testutil "github.com/gollin/packages/support/money/tests"
)

func TestNewManager(t *testing.T) {
	t.Run("creates manager with default currency manager", func(t *testing.T) {
		mm := NewManager()
		if mm == nil {
			t.Fatal("NewManager() returned nil")
		}
		if mm.currencyManager == nil {
			t.Fatal("currencyManager should not be nil")
		}
	})
}

func TestNewManagerWith(t *testing.T) {
	t.Run("creates manager with custom currency manager", func(t *testing.T) {
		cm := currency.NewManager()
		mm, err := NewManagerWith(cm)
		if err != nil {
			t.Fatalf("NewManagerWith() unexpected error: %v", err)
		}
		if mm == nil {
			t.Fatal("NewManagerWith() returned nil")
		}
		if mm.currencyManager != cm {
			t.Fatal("currencyManager should be the same instance")
		}
	})

	t.Run("returns error with nil currency manager", func(t *testing.T) {
		mm, err := NewManagerWith(nil)
		if err == nil {
			t.Fatal("NewManagerWith(nil) expected error, got nil")
		}
		if mm != nil {
			t.Fatal("NewManagerWith(nil) should return nil manager")
		}
		if err != exception.ErrNoCurrencyManager {
			t.Fatalf("NewManagerWith(nil) error = %v, want %v", err, exception.ErrNoCurrencyManager)
		}
	})
}

func TestManagerCreate(t *testing.T) {
	mm := NewManager()

	t.Run("creates money with valid currency", func(t *testing.T) {
		m := mm.Create(1500, currency.SGD)
		if m == nil {
			t.Fatal("Create() returned nil")
		}

		amount := testutil.TestRequire(t, m.Amount)
		if amount != 1500 {
			t.Fatalf("Amount() = %d, want %d", amount, 1500)
		}

		curr := testutil.TestRequire(t, m.Currency)
		if curr.Code != currency.SGD {
			t.Fatalf("Currency().Code = %s, want %s", curr.Code, currency.SGD)
		}
	})

	t.Run("creates money with different currencies", func(t *testing.T) {
		tests := []struct {
			name   string
			amount int64
			code   string
		}{
			{name: "SGD", amount: 1000, code: currency.SGD},
			{name: "EUR", amount: 2500, code: currency.EUR},
			{name: "GBP", amount: 750, code: currency.GBP},
			{name: "JPY", amount: 10000, code: currency.JPY},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				m := mm.Create(tt.amount, tt.code)
				amount := testutil.TestRequire(t, m.Amount)
				if amount != tt.amount {
					t.Fatalf("Amount() = %d, want %d", amount, tt.amount)
				}

				curr := testutil.TestRequire(t, m.Currency)
				if curr.Code != tt.code {
					t.Fatalf("Currency().Code = %s, want %s", curr.Code, tt.code)
				}
			})
		}
	})

	t.Run("handles unknown currency code", func(t *testing.T) {
		m := mm.Create(1000, "UNKNOWN")
		if m == nil {
			t.Fatal("Create() returned nil")
		}
		// Should fall back to default currency
		curr := testutil.TestRequire(t, m.Currency)
		if curr == nil {
			t.Fatal("Currency() returned nil")
		}
	})
}

func TestManagerCreateFromFloat(t *testing.T) {
	mm := NewManager()

	tests := []struct {
		name   string
		amount float64
		code   string
		want   int64
	}{
		{name: "sgd trailing decimals", amount: 12.349, code: currency.SGD, want: 1235},
		{name: "no fraction currency", amount: 98.75, code: currency.JPY, want: 99},
		{name: "negative sgd trailing decimals", amount: -12.349, code: currency.SGD, want: -1235},
		{name: "negative no fraction currency", amount: -98.75, code: currency.JPY, want: -99},
		{name: "negative with more decimals", amount: -1.567, code: currency.SGD, want: -157},
		{name: "positive with more decimals", amount: 1.567, code: currency.SGD, want: 157},
		{name: "zero", amount: 0.0, code: currency.SGD, want: 0},
		{name: "large amount", amount: 123456.78, code: currency.SGD, want: 12345678},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := mm.CreateFromFloat(tt.amount, tt.code)
			amount := testutil.TestRequire(t, m.Amount)
			if amount != tt.want {
				t.Fatalf("Amount() = %d, want %d", amount, tt.want)
			}

			curr := testutil.TestRequire(t, m.Currency)
			if curr.Code != tt.code {
				t.Fatalf("Currency().Code = %s, want %s", curr.Code, tt.code)
			}
		})
	}
}

func TestManagerCreateFromString(t *testing.T) {
	mm := NewManager()

	tests := []struct {
		name      string
		amount    string
		code      string
		want      int64
		wantError bool
	}{
		{name: "positive sgd with decimals", amount: "12.34", code: currency.SGD, want: 1234, wantError: false},
		{name: "negative sgd with decimals", amount: "-99.99", code: currency.SGD, want: -9999, wantError: false},
		{name: "positive integer sgd", amount: "100", code: currency.SGD, want: 10000, wantError: false},
		{name: "negative integer sgd", amount: "-50", code: currency.SGD, want: -5000, wantError: false},
		{name: "zero", amount: "0", code: currency.SGD, want: 0, wantError: false},
		{name: "zero with decimals", amount: "0.00", code: currency.SGD, want: 0, wantError: false},
		{name: "jpy no decimals", amount: "1000", code: currency.JPY, want: 1000, wantError: false},
		{name: "negative jpy", amount: "-500", code: currency.JPY, want: -500, wantError: false},
		{name: "one decimal place", amount: "12.5", code: currency.SGD, want: 1250, wantError: false},
		{name: "leading plus sign", amount: "+12.34", code: currency.SGD, want: 1234, wantError: false},
		{name: "with spaces", amount: "  12.34  ", code: currency.SGD, want: 1234, wantError: false},
		{name: "decimal only", amount: ".99", code: currency.SGD, want: 99, wantError: false},
		{name: "negative decimal only", amount: "-.50", code: currency.SGD, want: -50, wantError: false},
		{name: "empty string", amount: "", code: currency.SGD, want: 0, wantError: true},
		{name: "invalid format", amount: "abc", code: currency.SGD, want: 0, wantError: true},
		{name: "multiple decimals", amount: "12.34.56", code: currency.SGD, want: 0, wantError: true},
		{name: "too many decimal places sgd", amount: "12.345", code: currency.SGD, want: 0, wantError: true},
		{name: "jpy with decimals not allowed", amount: "100.5", code: currency.JPY, want: 0, wantError: true},
		{name: "bhd 3 decimals", amount: "12.345", code: currency.BHD, want: 12345, wantError: false},
		{name: "bhd negative 3 decimals", amount: "-99.999", code: currency.BHD, want: -99999, wantError: false},
		{name: "bhd too many decimals", amount: "12.3456", code: currency.BHD, want: 0, wantError: true},
		{name: "clf 4 decimals", amount: "99.9999", code: currency.CLF, want: 999999, wantError: false},
		{name: "clf negative 4 decimals", amount: "-12.3456", code: currency.CLF, want: -123456, wantError: false},
		{name: "clf too many decimals", amount: "99.99999", code: currency.CLF, want: 0, wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := mm.CreateFromString(tt.amount, tt.code)
			if tt.wantError {
				if err == nil {
					t.Fatal("CreateFromString() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("CreateFromString() unexpected error: %v", err)
			}

			amount := testutil.TestRequire(t, m.Amount)
			if amount != tt.want {
				t.Fatalf("Amount() = %d, want %d", amount, tt.want)
			}

			curr := testutil.TestRequire(t, m.Currency)
			if curr.Code != tt.code {
				t.Fatalf("Currency().Code = %s, want %s", curr.Code, tt.code)
			}
		})
	}
}

func TestManagerCreateFromStringErrors(t *testing.T) {
	mm := NewManager()

	tests := []struct {
		name      string
		amount    string
		code      string
		wantError error
	}{
		{
			name:      "empty string",
			amount:    "",
			code:      currency.SGD,
			wantError: exception.ErrEmptyAmountString,
		},
		{
			name:      "multiple decimal points",
			amount:    "12.34.56",
			code:      currency.SGD,
			wantError: exception.ErrInvalidAmountMultiple,
		},
		{
			name:      "too many decimal places",
			amount:    "12.345",
			code:      currency.SGD,
			wantError: exception.ErrInvalidAmountFraction,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := mm.CreateFromString(tt.amount, tt.code)
			if err == nil {
				t.Fatal("CreateFromString() expected error, got nil")
			}
		})
	}
}

func TestManagerGetCurrencyManager(t *testing.T) {
	cm := currency.NewManager()
	mm, err := NewManagerWith(cm)
	if err != nil {
		t.Fatalf("NewManagerWith() unexpected error: %v", err)
	}

	got := mm.GetCurrencyManager()
	if got != cm {
		t.Fatal("GetCurrencyManager() returned different instance")
	}
}

func TestManagerWithCustomCurrencyManager(t *testing.T) {
	// Create a custom currency manager with a custom dataset
	customCurrency := &currency.Currency{
		Code:        "TST",
		Grapheme:    "T$",
		Template:    "1 $",
		Decimal:     ".",
		Thousand:    ",",
		NumericCode: "999",
		Fraction:    2,
	}

	dataset := map[string]*currency.Currency{
		"TST": customCurrency,
	}

	cm, err := currency.NewManagerFor(nil, &dataset)
	if err != nil {
		t.Fatalf("failed to create custom currency manager: %v", err)
	}

	mm, err := NewManagerWith(cm)
	if err != nil {
		t.Fatalf("NewManagerWith() unexpected error: %v", err)
	}

	t.Run("creates money with custom currency", func(t *testing.T) {
		m := mm.Create(1500, "TST")
		curr := testutil.TestRequire(t, m.Currency)
		if curr.Code != "TST" {
			t.Fatalf("Currency().Code = %s, want TST", curr.Code)
		}
		if curr.Grapheme != "T$" {
			t.Fatalf("Currency().Grapheme = %s, want T$", curr.Grapheme)
		}
	})

	t.Run("creates money from string with custom currency", func(t *testing.T) {
		m, err := mm.CreateFromString("10.50", "TST")
		if err != nil {
			t.Fatalf("CreateFromString() unexpected error: %v", err)
		}
		amount := testutil.TestRequire(t, m.Amount)
		if amount != 1050 {
			t.Fatalf("Amount() = %d, want 1050", amount)
		}
	})
}

func TestManagerConcurrency(t *testing.T) {
	t.Parallel()

	mm := NewManager()

	// Test that multiple goroutines can use the same manager safely
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		amount := int64(i * 100)
		wg.Add(1)
		go func(amount int64) {
			defer wg.Done()
			m := mm.Create(amount, currency.SGD)

			gotAmount, err := m.Amount()
			if err != nil {
				t.Errorf("Amount() unexpected error: %v", err)
				return
			}
			if gotAmount != amount {
				t.Errorf("Amount() = %d, want %d", gotAmount, amount)
			}
		}(amount)
	}

	wg.Wait()
}

func TestManagerAdd(t *testing.T) {
	mm := NewManager()

	t.Run("adds two money values", func(t *testing.T) {
		m1 := mm.Create(1000, currency.SGD)
		m2 := mm.Create(500, currency.SGD)

		result, err := mm.Add(m1, m2)
		testutil.TestRequireNoErr(t, err)

		amount := testutil.TestRequire(t, result.Amount)
		if amount != 1500 {
			t.Fatalf("Amount() = %d, want 1500", amount)
		}
	})

	t.Run("adds multiple money values", func(t *testing.T) {
		m1 := mm.Create(1000, currency.SGD)
		m2 := mm.Create(500, currency.SGD)
		m3 := mm.Create(250, currency.SGD)

		result, err := mm.Add(m1, m2, m3)
		testutil.TestRequireNoErr(t, err)

		amount := testutil.TestRequire(t, result.Amount)
		if amount != 1750 {
			t.Fatalf("Amount() = %d, want 1750", amount)
		}
	})

	t.Run("returns same money when no values to add", func(t *testing.T) {
		m := mm.Create(1000, currency.SGD)
		result, err := mm.Add(m)
		testutil.TestRequireNoErr(t, err)

		amount := testutil.TestRequire(t, result.Amount)
		if amount != 1000 {
			t.Fatalf("Amount() = %d, want 1000", amount)
		}
	})

	t.Run("returns error for different currencies", func(t *testing.T) {
		m1 := mm.Create(1000, currency.USD)
		m2 := mm.Create(500, currency.EUR)

		_, err := mm.Add(m1, m2)
		if err == nil {
			t.Fatal("Add() expected error for different currencies, got nil")
		}
	})
}

func TestManagerSubtract(t *testing.T) {
	mm := NewManager()

	t.Run("subtracts two money values", func(t *testing.T) {
		m1 := mm.Create(1000, currency.SGD)
		m2 := mm.Create(300, currency.SGD)

		result, err := mm.Subtract(m1, m2)
		testutil.TestRequireNoErr(t, err)

		amount := testutil.TestRequire(t, result.Amount)
		if amount != 700 {
			t.Fatalf("Amount() = %d, want 700", amount)
		}
	})

	t.Run("subtracts multiple money values", func(t *testing.T) {
		m1 := mm.Create(1000, currency.SGD)
		m2 := mm.Create(200, currency.SGD)
		m3 := mm.Create(300, currency.SGD)

		result, err := mm.Subtract(m1, m2, m3)
		testutil.TestRequireNoErr(t, err)

		amount := testutil.TestRequire(t, result.Amount)
		if amount != 500 {
			t.Fatalf("Amount() = %d, want 500", amount)
		}
	})

	t.Run("returns error for different currencies", func(t *testing.T) {
		m1 := mm.Create(1000, currency.USD)
		m2 := mm.Create(500, currency.EUR)

		_, err := mm.Subtract(m1, m2)
		if err == nil {
			t.Fatal("Subtract() expected error for different currencies, got nil")
		}
	})
}

func TestManagerMultiply(t *testing.T) {
	mm := NewManager()

	t.Run("multiplies by single value", func(t *testing.T) {
		m := mm.Create(100, currency.SGD)

		result, err := mm.Multiply(m, 3)
		testutil.TestRequireNoErr(t, err)

		amount := testutil.TestRequire(t, result.Amount)
		if amount != 300 {
			t.Fatalf("Amount() = %d, want 300", amount)
		}
	})

	t.Run("multiplies by multiple values", func(t *testing.T) {
		m := mm.Create(10, currency.SGD)

		result, err := mm.Multiply(m, 2, 3, 5)
		testutil.TestRequireNoErr(t, err)

		amount := testutil.TestRequire(t, result.Amount)
		if amount != 300 {
			t.Fatalf("Amount() = %d, want 300", amount)
		}
	})

	t.Run("returns error when no multipliers provided", func(t *testing.T) {
		m := mm.Create(100, currency.SGD)

		_, err := mm.Multiply(m)
		if err == nil {
			t.Fatal("Multiply() expected error for no multipliers, got nil")
		}
	})
}

func TestManagerAbsolute(t *testing.T) {
	mm := NewManager()

	t.Run("returns absolute value of negative amount", func(t *testing.T) {
		m := mm.Create(-1000, currency.SGD)

		result, err := mm.Absolute(m)
		testutil.TestRequireNoErr(t, err)

		amount := testutil.TestRequire(t, result.Amount)
		if amount != 1000 {
			t.Fatalf("Amount() = %d, want 1000", amount)
		}
	})

	t.Run("returns same value for positive amount", func(t *testing.T) {
		m := mm.Create(1000, currency.SGD)

		result, err := mm.Absolute(m)
		testutil.TestRequireNoErr(t, err)

		amount := testutil.TestRequire(t, result.Amount)
		if amount != 1000 {
			t.Fatalf("Amount() = %d, want 1000", amount)
		}
	})
}

func TestManagerNegative(t *testing.T) {
	mm := NewManager()

	t.Run("returns negative of positive amount", func(t *testing.T) {
		m := mm.Create(1000, currency.SGD)

		result, err := mm.Negative(m)
		testutil.TestRequireNoErr(t, err)

		amount := testutil.TestRequire(t, result.Amount)
		if amount != -1000 {
			t.Fatalf("Amount() = %d, want -1000", amount)
		}
	})

	t.Run("returns positive of negative amount", func(t *testing.T) {
		m := mm.Create(-1000, currency.SGD)

		result, err := mm.Negative(m)
		testutil.TestRequireNoErr(t, err)

		amount := testutil.TestRequire(t, result.Amount)
		if amount != 1000 {
			t.Fatalf("Amount() = %d, want 1000", amount)
		}
	})
}

func TestManagerRound(t *testing.T) {
	mm := NewManager()

	t.Run("rounds SGD amount", func(t *testing.T) {
		m := mm.Create(1234, currency.SGD)

		result, err := mm.Round(m)
		testutil.TestRequireNoErr(t, err)

		amount := testutil.TestRequire(t, result.Amount)
		if amount != 1200 {
			t.Fatalf("Amount() = %d, want 1200", amount)
		}
	})

	t.Run("rounds up when closer to next value", func(t *testing.T) {
		m := mm.Create(1256, currency.SGD)

		result, err := mm.Round(m)
		testutil.TestRequireNoErr(t, err)

		amount := testutil.TestRequire(t, result.Amount)
		if amount != 1300 {
			t.Fatalf("Amount() = %d, want 1300", amount)
		}
	})
}

func TestManagerSplit(t *testing.T) {
	mm := NewManager()

	t.Run("splits money evenly", func(t *testing.T) {
		m := mm.Create(1000, currency.SGD)

		results, err := mm.Split(m, 2)
		testutil.TestRequireNoErr(t, err)

		if len(results) != 2 {
			t.Fatalf("Split() returned %d results, want 2", len(results))
		}

		for i, result := range results {
			amount := testutil.TestRequire(t, result.Amount)
			if amount != 500 {
				t.Fatalf("results[%d].Amount() = %d, want 500", i, amount)
			}
		}
	})

	t.Run("distributes remainder to first parties", func(t *testing.T) {
		m := mm.Create(1000, currency.SGD)

		results, err := mm.Split(m, 3)
		testutil.TestRequireNoErr(t, err)

		if len(results) != 3 {
			t.Fatalf("Split() returned %d results, want 3", len(results))
		}

		expected := []int64{334, 333, 333}
		for i, result := range results {
			amount := testutil.TestRequire(t, result.Amount)
			if amount != expected[i] {
				t.Fatalf("results[%d].Amount() = %d, want %d", i, amount, expected[i])
			}
		}
	})

	t.Run("returns error for invalid split", func(t *testing.T) {
		m := mm.Create(1000, currency.SGD)

		_, err := mm.Split(m, 0)
		if err == nil {
			t.Fatal("Split() expected error for n=0, got nil")
		}
	})

	t.Run("handles negative amounts", func(t *testing.T) {
		m := mm.Create(-1000, currency.SGD)

		results, err := mm.Split(m, 3)
		testutil.TestRequireNoErr(t, err)

		expected := []int64{-334, -333, -333}
		for i, result := range results {
			amount := testutil.TestRequire(t, result.Amount)
			if amount != expected[i] {
				t.Fatalf("results[%d].Amount() = %d, want %d", i, amount, expected[i])
			}
		}
	})
}

func TestManagerAllocate(t *testing.T) {
	mm := NewManager()

	t.Run("allocates money by ratios", func(t *testing.T) {
		m := mm.Create(1000, currency.SGD)

		results, err := mm.Allocate(m, 1, 1, 1)
		testutil.TestRequireNoErr(t, err)

		if len(results) != 3 {
			t.Fatalf("Allocate() returned %d results, want 3", len(results))
		}

		expected := []int64{334, 333, 333}
		for i, result := range results {
			amount := testutil.TestRequire(t, result.Amount)
			if amount != expected[i] {
				t.Fatalf("results[%d].Amount() = %d, want %d", i, amount, expected[i])
			}
		}
	})

	t.Run("allocates with different ratios", func(t *testing.T) {
		m := mm.Create(1000, currency.SGD)

		results, err := mm.Allocate(m, 1, 2, 3)
		testutil.TestRequireNoErr(t, err)

		if len(results) != 3 {
			t.Fatalf("Allocate() returned %d results, want 3", len(results))
		}

		// 1000 / 6 = 166.666...
		// 1/6 = 166, 2/6 = 333, 3/6 = 500
		// Leftover = 1 (distributed to first party)
		expected := []int64{167, 333, 500}
		for i, result := range results {
			amount := testutil.TestRequire(t, result.Amount)
			if amount != expected[i] {
				t.Fatalf("results[%d].Amount() = %d, want %d", i, amount, expected[i])
			}
		}
	})

	t.Run("returns error for negative ratios", func(t *testing.T) {
		m := mm.Create(1000, currency.SGD)

		_, err := mm.Allocate(m, 1, -2, 3)
		if err == nil {
			t.Fatal("Allocate() expected error for negative ratios, got nil")
		}
	})

	t.Run("returns error when no ratios provided", func(t *testing.T) {
		m := mm.Create(1000, currency.SGD)

		_, err := mm.Allocate(m)
		if err == nil {
			t.Fatal("Allocate() expected error for no ratios, got nil")
		}
	})

	t.Run("handles zero ratios", func(t *testing.T) {
		m := mm.Create(1000, currency.SGD)

		results, err := mm.Allocate(m, 0, 0, 0)
		testutil.TestRequireNoErr(t, err)

		if len(results) != 3 {
			t.Fatalf("Allocate() returned %d results, want 3", len(results))
		}

		for i, result := range results {
			amount := testutil.TestRequire(t, result.Amount)
			if amount != 0 {
				t.Fatalf("results[%d].Amount() = %d, want 0", i, amount)
			}
		}
	})
}

func TestManagerCreateFromString_ErrorPaths(t *testing.T) {
	mm := NewManager()

	if _, err := mm.CreateFromString("   ", currency.SGD); !errors.Is(err, exception.ErrEmptyAmountString) {
		t.Fatalf("CreateFromString(empty) error = %v, want ErrEmptyAmountString", err)
	}

	if _, err := mm.CreateFromString("12.345", currency.SGD); !errors.Is(err, exception.ErrInvalidAmountFraction) {
		t.Fatalf("CreateFromString(fraction too long) error = %v, want ErrInvalidAmountFraction", err)
	}

	if _, err := mm.CreateFromString("12.a", currency.SGD); !errors.Is(err, exception.ErrInvalidAmount) {
		t.Fatalf("CreateFromString(invalid digits) error = %v, want ErrInvalidAmount", err)
	}

	mm.parser = nil
	if _, err := mm.CreateFromString("+1.00", currency.SGD); !errors.Is(err, exception.ErrParserNotProvided) {
		t.Fatalf("CreateFromString(nil parser) error = %v, want ErrParserNotProvided", err)
	}
}

func TestManagerMultiply_ErrorPaths(t *testing.T) {
	mm := NewManager()
	m := mm.Create(10, currency.SGD)

	if _, err := mm.Multiply(m); !errors.Is(err, exception.ErrNoMultipliersProvided) {
		t.Fatalf("Multiply(no values) error = %v, want ErrNoMultipliersProvided", err)
	}

	huge := mm.Create(math.MaxInt64, currency.SGD)
	if _, err := mm.Multiply(huge, 2); err == nil {
		t.Fatal("Multiply() expected overflow error")
	}
}
