package money

import (
	"testing"

	"github.com/gollin/packages/support/money/currency"
	testutil "github.com/gollin/packages/support/money/tests"
)

func TestMoney_Value(t *testing.T) {
	m := NewManager().Create(1234, "SGD")

	val, err := m.Value()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := "1234|SGD"
	if val != expected {
		t.Errorf("expected %q, got %q", expected, val)
	}
}

func TestMoney_Scan(t *testing.T) {
	tests := []struct {
		name       string
		input      interface{}
		wantAmount int64
		wantCode   string
		wantErr    bool
	}{
		{
			name:       "valid money string",
			input:      "1234|SGD",
			wantAmount: 1234,
			wantCode:   "SGD",
			wantErr:    false,
		},
		{
			name:       "negative amount",
			input:      "-5000|EUR",
			wantAmount: -5000,
			wantCode:   "EUR",
			wantErr:    false,
		},
		{
			name:       "valid money bytes",
			input:      []byte("1234|SGD"),
			wantAmount: 1234,
			wantCode:   "SGD",
			wantErr:    false,
		},
		{
			name:       "negative amount bytes",
			input:      []byte("-5000|EUR"),
			wantAmount: -5000,
			wantCode:   "EUR",
			wantErr:    false,
		},
		{
			name:    "invalid format - missing separator",
			input:   "1234SGD",
			wantErr: true,
		},
		{
			name:    "invalid format - empty amount",
			input:   "|SGD",
			wantErr: true,
		},
		{
			name:    "invalid format - empty currency",
			input:   "1234|",
			wantErr: true,
		},
		{
			name:    "invalid type",
			input:   123,
			wantErr: true,
		},
		{
			name:    "invalid amount parsing",
			input:   "ABC|SGD",
			wantErr: true,
		},
		{
			name:    "invalid currency code scanning",
			input:   "123|INVALID",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var m Money
			err := m.Scan(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("DbScan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if testutil.TestRequire(t, m.Amount) != tt.wantAmount {
					t.Errorf("Amount() = %v, want %v", testutil.TestRequire(t, m.Amount), tt.wantAmount)
				}

				if testutil.TestRequire(t, m.Currency).Code != tt.wantCode {
					t.Errorf("Currency().Code = %v, want %v", testutil.TestRequire(t, m.Currency).Code, tt.wantCode)
				}
			}
		})
	}
}

func TestCurrency_Value(t *testing.T) {
	t.Parallel()
	c := currency.NewManager().Resolve("SGD")

	val, err := c.DbValue()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if val != "SGD" {
		t.Errorf("expected %q, got %q", "SGD", val)
	}
}

func TestCurrency_Scan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    interface{}
		wantCode string
		wantErr  bool
	}{
		{
			name:     "valid currency code",
			input:    "SGD",
			wantCode: "SGD",
			wantErr:  false,
		},
		{
			name:     "lowercase currency code",
			input:    "eur",
			wantCode: "EUR",
			wantErr:  false,
		},
		{
			name:     "valid currency bytes",
			input:    []byte("SGD"),
			wantCode: "SGD",
			wantErr:  false,
		},
		{
			name:     "lowercase currency bytes",
			input:    []byte("gbp"),
			wantCode: "GBP",
			wantErr:  false,
		},
		{
			name:    "invalid type",
			input:   123,
			wantErr: true,
		},
		{
			name:    "invalid currency code",
			input:   "INVALID",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var c currency.Currency
			err := c.DbScan(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("DbScan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && c.Code != tt.wantCode {
				t.Errorf("Code = %v, want %v", c.Code, tt.wantCode)
			}
		})
	}
}

func TestSetDBMoneyValueSeparator(t *testing.T) {
	// Save the original separator to restore after a test
	originalSep := GetDBMoneyValueSeparator()
	defer func(separator string) {
		err := SetDBMoneyValueSeparator(separator)
		if err != nil {
			t.Fatalf("Failed to restore original separator: %v", err)
		}
	}(originalSep)

	tests := []struct {
		name       string
		separator  string
		shouldFail bool
	}{
		{
			name:       "comma separator",
			separator:  ",",
			shouldFail: false,
		},
		{
			name:       "colon separator",
			separator:  ":",
			shouldFail: false,
		},
		{
			name:       "double dash separator",
			separator:  "--",
			shouldFail: false,
		},
		{
			name:       "space separator (should fail)",
			separator:  " ",
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := SetDBMoneyValueSeparator(tt.separator)

			if tt.shouldFail {
				if err == nil {
					t.Fatalf("SetDBMoneyValueSeparator() expected an error, but got none")
				}
				return // Don't proceed with success checks if an error is expected
			}

			if err != nil {
				t.Fatalf("SetDBMoneyValueSeparator() error: %v", err)
			}

			// Verify getter returns the new separator
			if got := GetDBMoneyValueSeparator(); got != tt.separator {
				t.Errorf("GetDBMoneyValueSeparator() = %q, want %q", got, tt.separator)
			}

			// Verify Value() uses the new separator
			m := NewManager().Create(1234, currency.SGD)
			val, err := m.Value()
			if err != nil {
				t.Fatalf("Value() unexpected error: %v", err)
			}

			expected := "1234" + tt.separator + "SGD"
			if val != expected {
				t.Errorf("Value() = %q, want %q", val, expected)
			}

			// Verify Scan() works with the new separator
			var scanned Money
			if err := scanned.Scan(expected); err != nil {
				t.Fatalf("Scan() unexpected error: %v", err)
			}

			if testutil.TestRequire(t, scanned.Amount) != 1234 {
				t.Errorf("Scanned amount = %d, want 1234", testutil.TestRequire(t, scanned.Amount))
			}

			if testutil.TestRequire(t, scanned.Currency).Code != currency.SGD {
				t.Errorf("Scanned currency = %s, want SGD", testutil.TestRequire(t, scanned.Currency).Code)
			}
		})
	}
}

func TestSetDBMoneyValueSeparator_ThreadSafety(t *testing.T) {
	// Save the original separator to restore after a test
	originalSep := GetDBMoneyValueSeparator()
	defer func(separator string) {
		err := SetDBMoneyValueSeparator(separator)
		if err != nil {
			t.Fatalf("Failed to restore original separator: %v", err)
		}
	}(originalSep)

	done := make(chan bool)
	iterations := 1000

	// Concurrent readers
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < iterations; j++ {
				sep := GetDBMoneyValueSeparator()
				// Separator should always be a valid value (not empty)
				if sep == "" {
					t.Error("GetDBMoneyValueSeparator() returned empty string")
				}

				// Test Value() during concurrent access
				m := NewManager().Create(100, currency.SGD)
				if _, err := m.Value(); err != nil {
					t.Errorf("Value() error during concurrent access: %v", err)
				}
			}
			done <- true
		}()
	}

	// Concurrent writers
	separators := []string{"|", ",", ":", ";", "-"}
	for i := 0; i < 5; i++ {
		go func(sep string) {
			for j := 0; j < iterations; j++ {
				err := SetDBMoneyValueSeparator(sep)

				if err != nil {
					t.Errorf("SetDBMoneyValueSeparator() error during concurrent access: %v", err)
				}
			}
			done <- true
		}(separators[i])
	}

	// Wait for all goroutines to complete
	for i := 0; i < 15; i++ {
		<-done
	}

	// Verify we can still get a valid separator after concurrent access
	finalSep := GetDBMoneyValueSeparator()
	if finalSep == "" {
		t.Error("Final separator is empty after concurrent access")
	}
}

func TestSetDBMoneyValueSeparator_RoundTrip(t *testing.T) {
	// Save the original separator to restore after a test
	originalSep := GetDBMoneyValueSeparator()
	defer func(separator string) {
		err := SetDBMoneyValueSeparator(separator)
		if err != nil {
			t.Fatalf("Failed to restore original separator: %v", err)
		}
	}(originalSep)

	testCases := []struct {
		separator string
		amount    int64
		currency  string
	}{
		{separator: "||", amount: 12345, currency: currency.EUR},
		{separator: "::", amount: -500, currency: currency.GBP},
		{separator: "##", amount: 0, currency: currency.JPY},
	}

	for _, tc := range testCases {
		t.Run(tc.separator, func(t *testing.T) {
			err := SetDBMoneyValueSeparator(tc.separator)

			if err != nil {
				t.Fatalf("SetDBMoneyValueSeparator() error: %v", err)
			}

			// Create money and convert to value
			original := NewManager().Create(tc.amount, tc.currency)
			val, err := original.Value()
			if err != nil {
				t.Fatalf("Value() error: %v", err)
			}

			// Scan back from value
			var restored Money
			if err := restored.Scan(val); err != nil {
				t.Fatalf("Scan() error: %v", err)
			}

			// Verify round-trip preserves data
			if testutil.TestRequire(t, restored.Amount) != testutil.TestRequire(t, original.Amount) {
				t.Errorf("Amount after round-trip = %d, want %d", testutil.TestRequire(t, restored.Amount), testutil.TestRequire(t, original.Amount))
			}

			if testutil.TestRequire(t, restored.Currency).Code != testutil.TestRequire(t, original.Currency).Code {
				t.Errorf("Currency after round-trip = %s, want %s", testutil.TestRequire(t, restored.Currency).Code, testutil.TestRequire(t, original.Currency).Code)
			}
		})
	}
}

func TestMoney_Value_Coverage(t *testing.T) {
	t.Parallel()

	// Test nil Money
	var m *Money
	_, err := m.Value()
	if err == nil {
		t.Error("Value() called on nil Money should return error")
	}

	// Test Money with nil currency
	m2 := &Money{amount: 100, currency: nil}
	_, err = m2.Value()
	if err == nil {
		t.Error("Value() called on Money with nil currency should return error")
	}
}
