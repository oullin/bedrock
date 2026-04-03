package money

import (
	"encoding/json"
	"errors"
	"strconv"
	"sync"
	"testing"

	"github.com/gollin/packages/support/money/currency"
	"github.com/gollin/packages/support/money/exception"
	testutil "github.com/gollin/packages/support/money/tests"
)

func TestJSONSetMarshal(t *testing.T) {
	parser := NewJson()

	customMarshal := func(m Money) ([]byte, error) {
		amount, err := m.Amount()
		if err != nil {
			return nil, err
		}

		curr, err := m.Currency()
		if err != nil {
			return nil, err
		}

		return []byte(`{"amount":` + strconv.FormatInt(amount, 10) + `,"currency":"` + curr.Code + `","custom":true}`), nil
	}

	testutil.TestRequireNoErr(t, parser.SetMarshal(customMarshal))

	m := NewManager().Create(1000, currency.SGD)
	data, err := parser.Marshal(*m)

	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Unmarshal result failed: %v", err)
	}

	if custom, ok := result["custom"]; !ok || custom != true {
		t.Errorf("Custom marshaller was not used, got: %v", result)
	}
}

func TestJSONSetUnmarshal(t *testing.T) {
	parser := NewJson()

	customUnmarshal := func(m *Money, b []byte) error {
		*m = Money{
			amount:   9999,
			currency: currency.NewManager().Resolve(currency.SGD),
		}
		return nil
	}

	testutil.TestRequireNoErr(t, parser.SetUnmarshal(customUnmarshal))

	jsonData := []byte(`{"amount": 1000, "currency": "EUR"}`)
	var m Money
	err := parser.Unmarshal(&m, jsonData)

	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if testutil.TestRequire(t, m.Amount) != 9999 {
		t.Errorf("Custom unmarshaller was not used, got amount: %d, want: 9999", testutil.TestRequire(t, m.Amount))
	}

	if testutil.TestRequire(t, m.Currency).Code != currency.SGD {
		t.Errorf("Custom unmarshaller was not used, got currency: %s, want: %s", testutil.TestRequire(t, m.Currency).Code, currency.SGD)
	}
}

func TestJSONSetCurrency(t *testing.T) {
	parser := NewJson()

	testutil.TestRequireNoErr(t, parser.SetCurrency(func() (*currency.Currency, error) {
		return currency.NewManager().Resolve(currency.GBP), nil
	}))

	jsonData := []byte(`{"amount": 5000}`)
	var m Money
	err := parser.Unmarshal(&m, jsonData)

	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if testutil.TestRequire(t, m.Currency).Code != currency.GBP {
		t.Errorf("Custom currency was not used, got: %s, want: %s", testutil.TestRequire(t, m.Currency).Code, currency.GBP)
	}

	if testutil.TestRequire(t, m.Amount) != 5000 {
		t.Errorf("Amount incorrect, got: %d, want: 5000", testutil.TestRequire(t, m.Amount))
	}
}

func TestJSONParserThreadSafety(t *testing.T) {
	parser := NewJson()

	var wg sync.WaitGroup
	iterations := 100

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := parser.SetMarshal(parser.defaultMarshalJSON); err != nil {
				t.Errorf("SetMarshal() unexpected error: %v", err)
			}
		}()
	}

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := parser.SetUnmarshal(parser.defaultUnmarshalJSON); err != nil {
				t.Errorf("SetUnmarshal() unexpected error: %v", err)
			}
		}()
	}

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := parser.SetCurrency(parser.defaultJSONCurrency); err != nil {
				t.Errorf("SetCurrency() unexpected error: %v", err)
			}
		}()
	}

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m := NewManager().Create(1000, currency.SGD)
			if _, err := parser.Marshal(*m); err != nil {
				t.Errorf("Marshal() unexpected error: %v", err)
			}
		}()
	}

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var m Money
			if err := parser.Unmarshal(&m, []byte(`{"amount": 1000, "currency": "SGD"}`)); err != nil {
				t.Errorf("Unmarshal() unexpected error: %v", err)
			}
		}()
	}

	wg.Wait()
}

func TestJSONParserRestore(t *testing.T) {
	parser := NewJson()

	originalMarshal := parser.defaultMarshalJSON

	customMarshal := func(m Money) ([]byte, error) {
		return []byte(`{"custom":true}`), nil
	}
	testutil.TestRequireNoErr(t, parser.SetMarshal(customMarshal))

	m := NewManager().Create(1000, currency.SGD)
	data, err := parser.Marshal(*m)
	testutil.TestRequireNoErr(t, err)

	if string(data) != `{"custom":true}` {
		t.Errorf("Custom marshaller not active, got: %s", string(data))
	}

	testutil.TestRequireNoErr(t, parser.SetMarshal(originalMarshal))

	data, err = parser.Marshal(*m)
	testutil.TestRequireNoErr(t, err)

	var result map[string]interface{}

	if err = json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Unmarshal result failed: %v", err)
	}

	if _, hasCustom := result["custom"]; hasCustom {
		t.Errorf("Provider marshaller was not restored")
	}

	if result["amount"] != float64(1000) {
		t.Errorf("Provider marshaller not working correctly, got amount: %v", result["amount"])
	}
}

func TestCustomJSONCurrencyInUnmarshal(t *testing.T) {
	parser := NewJson()

	testutil.TestRequireNoErr(t, parser.SetCurrency(func() (*currency.Currency, error) {
		return currency.NewManager().Resolve(currency.EUR), nil
	}))

	tests := []struct {
		name     string
		json     string
		wantCode string
		wantAmt  int64
	}{
		{
			name:     "missing currency uses custom default",
			json:     `{"amount": 1000}`,
			wantCode: currency.EUR,
			wantAmt:  1000,
		},
		{
			name:     "empty currency uses custom default",
			json:     `{"amount": 2000, "currency": ""}`,
			wantCode: currency.EUR,
			wantAmt:  2000,
		},
		{
			name:     "whitespace currency uses custom default",
			json:     `{"amount": 2500, "currency": "   "}`,
			wantCode: currency.EUR,
			wantAmt:  2500,
		},
		{
			name:     "explicit currency overrides default",
			json:     `{"amount": 3000, "currency": "GBP"}`,
			wantCode: currency.GBP,
			wantAmt:  3000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m Money
			err := parser.Unmarshal(&m, []byte(tt.json))

			if err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}

			if testutil.TestRequire(t, m.Currency).Code != tt.wantCode {
				t.Errorf("Currency = %s, want %s", testutil.TestRequire(t, m.Currency).Code, tt.wantCode)
			}

			if testutil.TestRequire(t, m.Amount) != tt.wantAmt {
				t.Errorf("Amount = %d, want %d", testutil.TestRequire(t, m.Amount), tt.wantAmt)
			}
		})
	}
}

func TestNewJsonWithParser(t *testing.T) {
	// Custom functions for testing
	customMarshal := func(m Money) ([]byte, error) {
		amount, err := m.Amount()
		if err != nil {
			return nil, err
		}

		curr, err := m.Currency()
		if err != nil {
			return nil, err
		}

		return []byte(`{"amount":` + strconv.FormatInt(amount, 10) + `,"currency":"` + curr.Code + `","new_parser_marshal":true}`), nil
	}

	customUnmarshal := func(m *Money, b []byte) error {
		*m = Money{
			amount:   8888,
			currency: currency.NewManager().Resolve(currency.JPY),
		}
		return nil
	}

	customCurrency := func() (*currency.Currency, error) {
		return currency.NewManager().Resolve(currency.AUD), nil
	}

	t.Run("with all custom functions", func(t *testing.T) {
		parser := NewJsonWithParser(customUnmarshal, customMarshal, customCurrency)

		// Test Marshal
		m := NewManager().Create(1234, currency.SGD)
		data, err := parser.Marshal(*m)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		var result map[string]interface{}

		err = json.Unmarshal(data, &result)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if _, ok := result["new_parser_marshal"]; !ok {
			t.Errorf("Custom marshal function was not used")
		}

		// Test Unmarshal
		jsonData := []byte(`{"amount": 1000, "currency": "EUR"}`)
		var unmarshalledMoney Money
		err = parser.Unmarshal(&unmarshalledMoney, jsonData)

		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if testutil.TestRequire(t, unmarshalledMoney.Amount) != 8888 || testutil.TestRequire(t, unmarshalledMoney.Currency).Code != currency.JPY {
			t.Errorf("Custom unmarshal function was not used, got amount: %d, currency: %s", testutil.TestRequire(t, unmarshalledMoney.Amount), testutil.TestRequire(t, unmarshalledMoney.Currency).Code)
		}

		// Test default currency when currency is missing in JSON
		jsonDataNoCurrency := []byte(`{"amount": 5000}`)
		var moneyWithDefaultCurrency Money
		err = parser.Unmarshal(&moneyWithDefaultCurrency, jsonDataNoCurrency)

		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		// The custom unmarshal completely overrides parsing, so the currency will be JPY
		if testutil.TestRequire(t, moneyWithDefaultCurrency.Currency).Code != currency.JPY {
			t.Errorf("Expected custom unmarshal to set currency to JPY, got %s", testutil.TestRequire(t, moneyWithDefaultCurrency.Currency).Code)
		}
	})

	t.Run("with nil unmarshal", func(t *testing.T) {
		parser := NewJsonWithParser(nil, customMarshal, customCurrency)

		// Test Unmarshal uses default
		jsonData := []byte(`{"amount": 1000, "currency": "EUR"}`)
		var unmarshalledMoney Money
		err := parser.Unmarshal(&unmarshalledMoney, jsonData)

		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if testutil.TestRequire(t, unmarshalledMoney.Amount) != 1000 || testutil.TestRequire(t, unmarshalledMoney.Currency).Code != currency.EUR {
			t.Errorf("Provider unmarshal function was not used")
		}
	})

	t.Run("with nil marshal", func(t *testing.T) {
		parser := NewJsonWithParser(customUnmarshal, nil, customCurrency)

		// Test Marshal uses default
		m := NewManager().Create(1234, currency.SGD)
		data, err := parser.Marshal(*m)

		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		var result map[string]interface{}
		err = json.Unmarshal(data, &result)
		if err != nil {
			t.Errorf("Unmarshal failed: %v", err)
		}

		if _, ok := result["new_parser_marshal"]; ok {
			t.Errorf("Provider marshal function was not used")
		}

		if result["amount"] != float64(1234) || result["currency"] != currency.SGD {
			t.Errorf("Provider marshal function did not work as expected")
		}
	})

	t.Run("with nil currency", func(t *testing.T) {
		parser := NewJsonWithParser(nil, nil, nil) // All nil should result in default behavior

		// Test default currency when currency is missing in JSON
		jsonDataNoCurrency := []byte(`{"amount": 5000}`)
		var moneyWithDefaultCurrency Money
		err := parser.Unmarshal(&moneyWithDefaultCurrency, jsonDataNoCurrency)

		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		// Provider currency should be SGD as per defaultJSONCurrency
		if testutil.TestRequire(t, moneyWithDefaultCurrency.Currency).Code != currency.SGD {
			t.Errorf("Expected default currency to be SGD, got %s", testutil.TestRequire(t, moneyWithDefaultCurrency.Currency).Code)
		}
	})
}

func TestJSON_Coverage(t *testing.T) {
	t.Parallel()

	// Test defaultJSONCurrency with nil receiver
	var j *JSON
	_, err := j.defaultJSONCurrency()
	if err == nil {
		t.Error("defaultJSONCurrency(nil) should return error")
	}
}

func TestJSONDefaultMarshalJSON_CoveragePaths(t *testing.T) {
	t.Run("nil receiver", func(t *testing.T) {
		var j *JSON
		_, err := j.defaultMarshalJSON(Money{})
		if !errors.Is(err, exception.ErrNoJSONParserProvided) {
			t.Fatalf("defaultMarshalJSON error = %v, want ErrNoJSONParserProvided", err)
		}
	})

	t.Run("Money{} uses default", func(t *testing.T) {
		j := NewJson()

		b, err := j.defaultMarshalJSON(Money{})
		if err != nil {
			t.Fatalf("defaultMarshalJSON unexpected error: %v", err)
		}

		var decoded struct {
			Amount   int64  `json:"amount"`
			Currency string `json:"currency"`
		}
		if err := json.Unmarshal(b, &decoded); err != nil {
			t.Fatalf("json.Unmarshal() unexpected error: %v", err)
		}

		if decoded.Currency == "" || decoded.Amount != 0 {
			t.Fatalf("marshal result = (%d,%q), want (0, non-empty currency)", decoded.Amount, decoded.Currency)
		}
	})

	t.Run("Money with nil currency errors", func(t *testing.T) {
		j := NewJson()

		_, err := j.defaultMarshalJSON(Money{amount: 1, currency: nil})
		if !errors.Is(err, exception.ErrNoCurrencyInstance) {
			t.Fatalf("defaultMarshalJSON error = %v, want ErrNoCurrencyInstance", err)
		}
	})
}

func TestJSONDefaultUnmarshalJSON_CoveragePaths(t *testing.T) {
	t.Run("nil receiver", func(t *testing.T) {
		var j *JSON
		var m Money
		err := j.defaultUnmarshalJSON(&m, []byte(`{"amount": 10, "currency": "SGD"}`))
		if !errors.Is(err, exception.ErrNoJSONParserProvided) {
			t.Fatalf("defaultUnmarshalJSON error = %v, want ErrNoJSONParserProvided", err)
		}
	})

	t.Run("type error amount", func(t *testing.T) {
		j := NewJson()
		var m Money
		err := j.defaultUnmarshalJSON(&m, []byte(`{"amount": "nope", "currency": "SGD"}`))
		if !errors.Is(err, exception.ErrInvalidJSONUnmarshal) {
			t.Fatalf("defaultUnmarshalJSON error = %v, want ErrInvalidJSONUnmarshal", err)
		}
	})

	t.Run("type error currency", func(t *testing.T) {
		j := NewJson()
		var m Money
		err := j.defaultUnmarshalJSON(&m, []byte(`{"amount": 10, "currency": 123}`))
		if !errors.Is(err, exception.ErrInvalidJSONUnmarshal) {
			t.Fatalf("defaultUnmarshalJSON error = %v, want ErrInvalidJSONUnmarshal", err)
		}
	})

	t.Run("syntax error", func(t *testing.T) {
		j := NewJson()
		var m Money
		err := j.defaultUnmarshalJSON(&m, []byte(`{"amount":`))
		if !errors.Is(err, exception.ErrInvalidJSONUnmarshal) {
			t.Fatalf("defaultUnmarshalJSON error = %v, want ErrInvalidJSONUnmarshal", err)
		}
	})

	t.Run("decimal rounding positive and negative", func(t *testing.T) {
		j := NewJson()

		var pos Money
		testutil.TestRequireNoErr(t, j.defaultUnmarshalJSON(&pos, []byte(`{"amount": 12.50, "currency": "SGD"}`)))
		if testutil.TestRequire(t, pos.Amount) != 13 {
			t.Fatalf("amount = %d, want 13", testutil.TestRequire(t, pos.Amount))
		}

		var neg Money
		testutil.TestRequireNoErr(t, j.defaultUnmarshalJSON(&neg, []byte(`{"amount": -12.50, "currency": "SGD"}`)))
		if testutil.TestRequire(t, neg.Amount) != -13 {
			t.Fatalf("amount = %d, want -13", testutil.TestRequire(t, neg.Amount))
		}
	})

	t.Run("amount overflow beyond int64", func(t *testing.T) {
		j := NewJson()
		var m Money
		err := j.defaultUnmarshalJSON(&m, []byte(`{"amount": 9223372036854775808, "currency": "SGD"}`))
		if !errors.Is(err, exception.ErrInvalidJSONUnmarshal) {
			t.Fatalf("defaultUnmarshalJSON error = %v, want ErrInvalidJSONUnmarshal", err)
		}
	})

	t.Run("unknown currency", func(t *testing.T) {
		j := NewJson()
		var m Money
		err := j.defaultUnmarshalJSON(&m, []byte(`{"amount": 10, "currency": "ZZZ"}`))
		if !errors.Is(err, exception.ErrCurrencyNotFound) {
			t.Fatalf("defaultUnmarshalJSON error = %v, want ErrCurrencyNotFound", err)
		}
	})

	t.Run("default currency used when currency missing", func(t *testing.T) {
		j := NewJson()
		var m Money
		testutil.TestRequireNoErr(t, j.defaultUnmarshalJSON(&m, []byte(`{"amount": 10}`)))
		if testutil.TestRequire(t, m.Currency).Code != currency.SGD {
			t.Fatalf("currency = %s, want SGD", testutil.TestRequire(t, m.Currency).Code)
		}
	})
}
