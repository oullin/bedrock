package money

import (
	"encoding/json"
	"errors"
	"math"
	"testing"

	"github.com/gollin/packages/support/money/currency"
	"github.com/gollin/packages/support/money/exception"
	"github.com/gollin/packages/support/money/exchange"
	testutil "github.com/gollin/packages/support/money/tests"
)

var testManager = NewManager()

func TestNewAndAccessors(t *testing.T) {
	m := testManager.Create(1500, currency.SGD)

	amount := testutil.TestRequire(t, m.Amount)

	if amount != 1500 {
		t.Fatalf("Amount() = %d, want %d", amount, 1500)
	}

	curr := testutil.TestRequire(t, m.Currency)

	if curr.Code != currency.SGD {
		t.Fatalf("Currency().Code = %s, want %s", curr.Code, currency.SGD)
	}
}

func TestNewFromFloatTruncatesTowardZero(t *testing.T) {
	tests := []struct {
		name   string
		amount float64
		code   string
		want   int64
	}{
		{name: "usd trailing decimals", amount: 12.349, code: currency.SGD, want: 1235},
		{name: "no fraction currency", amount: 98.75, code: currency.JPY, want: 99},
		{name: "negative usd trailing decimals", amount: -12.349, code: currency.SGD, want: -1235},
		{name: "negative no fraction currency", amount: -98.75, code: currency.JPY, want: -99},
		{name: "negative with more decimals", amount: -1.567, code: currency.SGD, want: -157},
		{name: "positive with more decimals", amount: 1.567, code: currency.SGD, want: 157},
	}

	manager := NewManager()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := manager.CreateFromFloat(tt.amount, tt.code)
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

func TestNewFromString(t *testing.T) {
	tests := []struct {
		name      string
		amount    string
		code      string
		want      int64
		wantError bool
	}{
		{name: "positive usd with decimals", amount: "12.34", code: currency.SGD, want: 1234, wantError: false},
		{name: "negative usd with decimals", amount: "-99.99", code: currency.SGD, want: -9999, wantError: false},
		{name: "positive integer usd", amount: "100", code: currency.SGD, want: 10000, wantError: false},
		{name: "negative integer usd", amount: "-50", code: currency.SGD, want: -5000, wantError: false},
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
		{name: "too many decimal places usd", amount: "12.345", code: currency.SGD, want: 0, wantError: true},
		{name: "jpy with decimals not allowed", amount: "100.5", code: currency.JPY, want: 0, wantError: true},
		{name: "bhd 3 decimals", amount: "12.345", code: currency.BHD, want: 12345, wantError: false},
		{name: "bhd negative 3 decimals", amount: "-99.999", code: currency.BHD, want: -99999, wantError: false},
		{name: "bhd too many decimals", amount: "12.3456", code: currency.BHD, want: 0, wantError: true},
		{name: "clf 4 decimals", amount: "99.9999", code: currency.CLF, want: 999999, wantError: false},
		{name: "clf negative 4 decimals", amount: "-12.3456", code: currency.CLF, want: -123456, wantError: false},
		{
			name:      "clf too many decimals",
			amount:    "99.99999",
			code:      currency.CLF,
			want:      0,
			wantError: true,
		},
		{
			name:      "overflow int64",
			amount:    "9223372036854775808", // MaxInt64 + 1
			code:      currency.SGD,
			want:      0,
			wantError: true,
		},
	}

	manager := NewManager()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := manager.CreateFromString(tt.amount, tt.code)

			if tt.wantError {
				if err == nil {
					t.Fatalf("CreateFromString() expected error but got nil")
				}

				return
			}

			if err != nil {
				t.Fatalf("CreateFromString() unexpected error = %v", err)
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

func TestCurrencyChecks(t *testing.T) {
	sgd := testManager.Create(100, currency.SGD)
	eur := testManager.Create(50, currency.EUR)

	sameCurrency, err := sgd.SameCurrency(eur)

	if err != nil {
		t.Fatalf("SameCurrency() unexpected error: %v", err)
	}

	if sameCurrency {
		t.Fatal("SameCurrency() returned true for different currencies")
	}

	if err := sgd.AssertSameCurrency(eur); !errors.Is(err, exception.ErrCurrencyMismatch) {
		t.Fatalf("AssertSameCurrency() error = %v, want %v", err, exception.ErrCurrencyMismatch)
	}
}

func TestComparisons(t *testing.T) {
	left := testManager.Create(200, currency.SGD)
	right := testManager.Create(150, currency.SGD)

	if cmp := testutil.TestRequire(t, func() (int, error) { return left.CompareAmount(right) }); cmp != 1 {
		t.Fatalf("CompareAmount() = %d, want 1", cmp)
	}

	equal := testManager.Create(200, currency.SGD)
	ok, err := left.Equals(equal)

	if err != nil || !ok {
		t.Fatalf("Equals() = (%v, %v), want (true, nil)", ok, err)
	}

	gt, err := left.GreaterThan(right)

	if err != nil || !gt {
		t.Fatalf("GreaterThan() = (%v, %v), want (true, nil)", gt, err)
	}

	gte, err := right.GreaterThanOrEqual(right)

	if err != nil || !gte {
		t.Fatalf("GreaterThanOrEqual() = (%v, %v), want (true, nil)", gte, err)
	}

	lt, err := right.LessThan(left)

	if err != nil || !lt {
		t.Fatalf("LessThan() = (%v, %v), want (true, nil)", lt, err)
	}

	lte, err := right.LessThanOrEqual(left)

	if err != nil || !lte {
		t.Fatalf("LessThanOrEqual() = (%v, %v), want (true, nil)", lte, err)
	}
}

func TestCompareCurrencyMismatch(t *testing.T) {
	sgd := testManager.Create(120, currency.SGD)
	eur := testManager.Create(80, currency.EUR)

	result, err := sgd.Compare(eur)

	if !errors.Is(err, exception.ErrCurrencyMismatch) {
		t.Fatalf("Compare() error = %v, want %v", err, exception.ErrCurrencyMismatch)
	}

	if result != int(testutil.TestRequire(t, sgd.Amount)) {
		t.Fatalf("Compare() result = %d, want %d", result, testutil.TestRequire(t, sgd.Amount))
	}
}

func TestSignChecks(t *testing.T) {
	zero := testManager.Create(0, currency.SGD)
	positive := testManager.Create(250, currency.SGD)
	negative := testManager.Create(-10, currency.SGD)

	if !testutil.TestRequire(t, zero.IsZero) {
		t.Fatal("IsZero() = false, want true")
	}

	if !testutil.TestRequire(t, positive.IsPositive) || testutil.TestRequire(t, positive.IsNegative) {
		t.Fatalf("IsPositive() = %v, IsNegative() = %v, want true/false", testutil.TestRequire(t, positive.IsPositive), testutil.TestRequire(t, positive.IsNegative))
	}

	if !testutil.TestRequire(t, negative.IsNegative) || testutil.TestRequire(t, negative.IsPositive) {
		t.Fatalf("IsNegative() = %v, IsPositive() = %v, want true/false", testutil.TestRequire(t, negative.IsNegative), testutil.TestRequire(t, negative.IsPositive))
	}
}

func TestAbsoluteAndNegative(t *testing.T) {
	m := testManager.Create(-300, currency.SGD)
	manager := NewManager()
	abs, err := manager.Absolute(m)
	testutil.TestRequireNoErr(t, err)

	if testutil.TestRequire(t, abs.Amount) != 300 || testutil.TestRequire(t, abs.Currency).Code != currency.SGD {
		t.Fatalf("Absolute() = (%d, %s), want (300, SGD)", testutil.TestRequire(t, abs.Amount), testutil.TestRequire(t, abs.Currency).Code)
	}

	neg, err := manager.Negative(abs)
	testutil.TestRequireNoErr(t, err)

	if testutil.TestRequire(t, neg.Amount) != -300 || testutil.TestRequire(t, neg.Currency).Code != currency.SGD {
		t.Fatalf("Negative() = (%d, %s), want (-300, SGD)", testutil.TestRequire(t, neg.Amount), testutil.TestRequire(t, neg.Currency).Code)
	}
}

func TestAddAndSubtract(t *testing.T) {
	base := testManager.Create(100, currency.SGD)
	addOne := testManager.Create(40, currency.SGD)
	addTwo := testManager.Create(10, currency.SGD)
	manager := NewManager()

	unchanged, err := manager.Add(base)

	if err != nil {
		t.Fatalf("Add() unexpected error: %v", err)
	}

	if unchanged == base {
		t.Fatal("Add() with no operands should not return the same pointer")
	}

	sum, err := manager.Add(base, addOne, addTwo)

	if err != nil {
		t.Fatalf("Add() unexpected error: %v", err)
	}

	if testutil.TestRequire(t, sum.Amount) != 150 {
		t.Fatalf("Add() amount = %d, want %d", testutil.TestRequire(t, sum.Amount), 150)
	}

	subtractUnchanged, err := manager.Subtract(base)

	if err != nil {
		t.Fatalf("Subtract() unexpected error: %v", err)
	}

	if subtractUnchanged == base {
		t.Fatal("Subtract() with no operands should not return the same pointer")
	}

	diff, err := manager.Subtract(base, addOne)

	if err != nil {
		t.Fatalf("Subtract() unexpected error: %v", err)
	}

	if testutil.TestRequire(t, diff.Amount) != 60 {
		t.Fatalf("Subtract() amount = %d, want %d", testutil.TestRequire(t, diff.Amount), 60)
	}

	_, err = manager.Add(base, testManager.Create(1, currency.EUR))

	if !errors.Is(err, exception.ErrCurrencyMismatch) {
		t.Fatalf("Add() error = %v, want %v", err, exception.ErrCurrencyMismatch)
	}

	_, err = manager.Subtract(base, testManager.Create(1, currency.EUR))

	if !errors.Is(err, exception.ErrCurrencyMismatch) {
		t.Fatalf("Subtract() error = %v, want %v", err, exception.ErrCurrencyMismatch)
	}
}

func TestMultiply(t *testing.T) {
	m := testManager.Create(10, currency.SGD)
	manager := NewManager()

	if _, err := manager.Multiply(m); err == nil {
		t.Fatal("Multiply() expected error when no multipliers are provided")
	}
}

func TestMultiplyWithValues(t *testing.T) {
	m := testManager.Create(10, currency.SGD)
	manager := NewManager()
	result, err := manager.Multiply(m, 2, 3)

	if err != nil {
		t.Fatalf("Multiply() unexpected error: %v", err)
	}

	if testutil.TestRequire(t, result.Amount) != 60 {
		t.Fatalf("Multiply() amount = %d, want %d", testutil.TestRequire(t, result.Amount), 60)
	}

	if testutil.TestRequire(t, result.Currency).Code != currency.SGD {
		t.Fatalf("Multiply() currency = %s, want SGD", testutil.TestRequire(t, result.Currency).Code)
	}
}

func TestRound(t *testing.T) {
	tests := []struct {
		name   string
		amount int64
		want   int64
	}{
		{name: "rounds half to even at midpoint", amount: 150, want: 100},
		{name: "rounds up above midpoint", amount: 155, want: 200},
		{name: "negative rounds", amount: -155, want: -200},
	}

	manager := NewManager()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := testManager.Create(tt.amount, currency.SGD)
			rounded, err := manager.Round(m)
			testutil.TestRequireNoErr(t, err)

			if testutil.TestRequire(t, rounded.Amount) != tt.want {
				t.Fatalf("Round() amount = %d, want %d", testutil.TestRequire(t, rounded.Amount), tt.want)
			}
		})
	}
}

func TestSplit(t *testing.T) {
	positive := testManager.Create(100, currency.SGD)
	manager := NewManager()
	parts, err := manager.Split(positive, 3)

	if err != nil {
		t.Fatalf("Split() unexpected error: %v", err)
	}

	want := []int64{34, 33, 33}

	for i, p := range parts {
		if testutil.TestRequire(t, p.Amount) != want[i] || testutil.TestRequire(t, p.Currency).Code != currency.SGD {
			t.Fatalf("Split() part %d = (%d, %s), want (%d, %s)", i, testutil.TestRequire(t, p.Amount), testutil.TestRequire(t, p.Currency).Code, want[i], currency.SGD)
		}
	}

	negative := testManager.Create(-100, currency.SGD)
	negParts, err := manager.Split(negative, 3)

	if err != nil {
		t.Fatalf("Split() unexpected error: %v", err)
	}

	negWant := []int64{-34, -33, -33}

	for i, p := range negParts {
		if testutil.TestRequire(t, p.Amount) != negWant[i] || testutil.TestRequire(t, p.Currency).Code != currency.SGD {
			t.Fatalf("Split() negative part %d = (%d, %s), want (%d, %s)", i, testutil.TestRequire(t, p.Amount), testutil.TestRequire(t, p.Currency).Code, negWant[i], currency.SGD)
		}
	}

	if _, err := manager.Split(positive, 0); err == nil {
		t.Fatal("Split() expected error when divisor is zero")
	}
}

func TestAllocate(t *testing.T) {
	base := testManager.Create(100, currency.SGD)
	manager := NewManager()

	if _, err := manager.Allocate(base); err == nil {
		t.Fatal("Allocate() expected error when no ratios provided")
	}

	if _, err := manager.Allocate(base, 1, -1); err == nil {
		t.Fatal("Allocate() expected error for negative ratios")
	}

	if _, err := manager.Allocate(base, int(math.MaxInt64), 1); err == nil {
		t.Fatal("Allocate() expected overflow error")
	}

	splits, err := manager.Allocate(base, 1, 1, 1)

	if err != nil {
		t.Fatalf("Allocate() unexpected error: %v", err)
	}

	want := []int64{34, 33, 33}

	for i, p := range splits {
		if testutil.TestRequire(t, p.Amount) != want[i] {
			t.Fatalf("Allocate() part %d = %d, want %d", i, testutil.TestRequire(t, p.Amount), want[i])
		}
	}

	negMoney := testManager.Create(-100, currency.SGD)
	negative, err := manager.Allocate(negMoney, 1, 2)

	if err != nil {
		t.Fatalf("Allocate() unexpected error: %v", err)
	}

	negWant := []int64{-34, -66}

	for i, p := range negative {
		if testutil.TestRequire(t, p.Amount) != negWant[i] {
			t.Fatalf("Allocate() negative part %d = %d, want %d", i, testutil.TestRequire(t, p.Amount), negWant[i])
		}
	}

	zeros, err := manager.Allocate(base, 0, 0, 0)

	if err != nil {
		t.Fatalf("Allocate() unexpected error when ratios sum to zero: %v", err)
	}

	for i, p := range zeros {
		if testutil.TestRequire(t, p.Amount) != 0 {
			t.Fatalf("Allocate() zero-ratio part %d = %d, want 0", i, testutil.TestRequire(t, p.Amount))
		}
	}
}

func TestDisplayAndMajorUnits(t *testing.T) {
	m := testManager.Create(12345, currency.SGD)

	if got := testutil.TestRequire(t, m.Display); got != "S$123.45" {
		t.Fatalf("Display() = %s, want S$123.45", got)
	}

	if major := testutil.TestRequire(t, m.AsMajorUnits); math.Abs(major-123.45) > 0.000001 {
		t.Fatalf("AsMajorUnits() = %f, want 123.45", major)
	}
}

func TestMarshalJSON(t *testing.T) {
	m := testManager.Create(1234, currency.SGD)
	data, err := m.MarshalJSON()

	if err != nil {
		t.Fatalf("MarshalJSON() unexpected error: %v", err)
	}

	if string(data) != `{"amount":1234,"currency":"SGD"}` {
		t.Fatalf("MarshalJSON() = %s, want %s", string(data), `{"amount":1234,"currency":"SGD"}`)
	}

	var zero Money
	zeroData, err := zero.MarshalJSON()

	if err != nil {
		t.Fatalf("MarshalJSON() zero value unexpected error: %v", err)
	}

	if string(zeroData) != `{"amount":0,"currency":"SGD"}` {
		t.Fatalf("MarshalJSON() zero value = %s, want %s", string(zeroData), `{"amount":0,"currency":"SGD"}`)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	payload := []byte(`{"amount": 1234, "currency": "SGD"}`)

	var m Money

	if err := json.Unmarshal(payload, &m); err != nil {
		t.Fatalf("json.Unmarshal() unexpected error: %v", err)
	}

	if testutil.TestRequire(t, m.Amount) != 1234 || testutil.TestRequire(t, m.Currency).Code != currency.SGD {
		t.Fatalf("UnmarshalJSON() = (%d, %s), want (1234, SGD)", testutil.TestRequire(t, m.Amount), testutil.TestRequire(t, m.Currency).Code)
	}

	var empty Money

	if err := json.Unmarshal([]byte(`{}`), &empty); err != nil {
		t.Fatalf("json.Unmarshal() empty unexpected error: %v", err)
	}

	if testutil.TestRequire(t, empty.Amount) != 0 {
		t.Fatalf("UnmarshalJSON() empty amount = %d, want 0", testutil.TestRequire(t, empty.Amount))
	}

	if err := json.Unmarshal([]byte(`{"amount": "oops", "currency": "SGD"}`), &m); !errors.Is(err, exception.ErrInvalidJSONUnmarshal) {
		t.Fatalf("UnmarshalJSON() error = %v, want %v", err, exception.ErrInvalidJSONUnmarshal)
	}

	if err := json.Unmarshal([]byte(`{"amount": 10, "currency": 1}`), &m); !errors.Is(err, exception.ErrInvalidJSONUnmarshal) {
		t.Fatalf("UnmarshalJSON() error = %v, want %v", err, exception.ErrInvalidJSONUnmarshal)
	}
}

func TestCompareSameCurrency(t *testing.T) {
	sgd100 := testManager.Create(100, currency.SGD)
	sgd200 := testManager.Create(200, currency.SGD)
	sgd100Clone := testManager.Create(100, currency.SGD)

	// Test <
	cmp, err := sgd100.Compare(sgd200)

	if err != nil {
		t.Fatalf("Compare() unexpected error: %v", err)
	}

	if cmp != -1 {
		t.Errorf("Compare(100, 200) = %d, want -1", cmp)
	}

	// Test >
	cmp, err = sgd200.Compare(sgd100)

	if err != nil {
		t.Fatalf("Compare() unexpected error: %v", err)
	}

	if cmp != 1 {
		t.Errorf("Compare(200, 100) = %d, want 1", cmp)
	}

	// Test ==
	cmp, err = sgd100.Compare(sgd100Clone)

	if err != nil {
		t.Fatalf("Compare() unexpected error: %v", err)
	}

	if cmp != 0 {
		t.Errorf("Compare(100, 100) = %d, want 0", cmp)
	}
}

func TestUnmarshalJSONInvalidSyntax(t *testing.T) {
	var m Money
	// Invalid JSON syntax
	err := m.UnmarshalJSON([]byte(`{"amount": 123`))

	if err == nil {
		t.Fatal("UnmarshalJSON() expected error on invalid JSON, got nil")
	}

	if !errors.Is(err, exception.ErrInvalidJSONUnmarshal) {
		t.Errorf("UnmarshalJSON() error = %v, want %v", err, exception.ErrInvalidJSONUnmarshal)
	}
}

func TestUnmarshalJSON_EdgeCases(t *testing.T) {
	var m Money

	// Float amount logic (rounds to integer)
	// 12.50 -> 13
	err := json.Unmarshal([]byte(`{"amount": 12.50, "currency": "SGD"}`), &m)

	if err != nil {
		t.Fatalf("UnmarshalJSON() float error: %v", err)
	}

	if testutil.TestRequire(t, m.Amount) != 13 {
		t.Errorf("Amount() = %d, want 13", testutil.TestRequire(t, m.Amount))
	}

	// Negative float
	// -12.50 -> -13
	err = json.Unmarshal([]byte(`{"amount": -12.50, "currency": "SGD"}`), &m)

	if err != nil {
		t.Fatalf("UnmarshalJSON() negative float error: %v", err)
	}

	if testutil.TestRequire(t, m.Amount) != -13 {
		t.Errorf("Amount() = %d, want -13", testutil.TestRequire(t, m.Amount))
	}

	// Type error (amount is non-numeric string)
	// json.Number might accept string tokens, but ParseInt/ParseFloat will fail.
	err = json.Unmarshal([]byte(`{"amount": "foo", "currency": "SGD"}`), &m)

	if !errors.Is(err, exception.ErrInvalidJSONUnmarshal) {
		t.Errorf("UnmarshalJSON() string amount error = %v, want ErrInvalidJSONUnmarshal", err)
	}

	// Type error (amount is boolean)
	// Triggers UnmarshalTypeError
	err = json.Unmarshal([]byte(`{"amount": true, "currency": "SGD"}`), &m)

	if !errors.Is(err, exception.ErrInvalidJSONUnmarshal) {
		t.Errorf("UnmarshalJSON() boolean amount error = %v, want ErrInvalidJSONUnmarshal", err)
	}

	// Type error (currency is number)
	err = json.Unmarshal([]byte(`{"amount": 100, "currency": 123}`), &m)

	if !errors.Is(err, exception.ErrInvalidJSONUnmarshal) {
		t.Errorf("UnmarshalJSON() number currency error = %v, want ErrInvalidJSONUnmarshal", err)
	}

	// Float overflow
	// 1e1000 overflows float64, ParseFloat returns error (ErrRange)
	err = json.Unmarshal([]byte(`{"amount": 1e1000, "currency": "SGD"}`), &m)

	if !errors.Is(err, exception.ErrInvalidJSONUnmarshal) {
		t.Errorf("UnmarshalJSON() float overflow error = %v, want ErrInvalidJSONUnmarshal", err)
	}
}

func TestUnmarshalJSON_LargePrecision(t *testing.T) {
	// Test that large numbers beyond float64's 2^53 precision limit are handled correctly
	// 9,007,199,254,740,993 is 2^53 + 1, which cannot be represented precisely in float64
	tests := []struct {
		name     string
		json     string
		expected int64
	}{
		{
			name:     "large integer within int64 range",
			json:     `{"amount": 9007199254740993, "currency": "SGD"}`,
			expected: 9007199254740993,
		},
		{
			name:     "large decimal within int64 range",
			json:     `{"amount": 9007199254740993.4, "currency": "SGD"}`,
			expected: 9007199254740993,
		},
		{
			name:     "large decimal with rounding",
			json:     `{"amount": 9007199254740993.5, "currency": "SGD"}`,
			expected: 9007199254740994,
		},
		{
			name:     "exponent notation",
			json:     `{"amount": 1.23e10, "currency": "SGD"}`,
			expected: 12300000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m Money
			err := json.Unmarshal([]byte(tt.json), &m)

			if err != nil {
				t.Fatalf("UnmarshalJSON() error: %v", err)
			}

			if testutil.TestRequire(t, m.Amount) != tt.expected {
				t.Errorf("Amount() = %d, want %d", testutil.TestRequire(t, m.Amount), tt.expected)
			}
		})
	}
}

func TestComparisons_Mismatch(t *testing.T) {
	sgd := testManager.Create(100, currency.SGD)
	eur := testManager.Create(100, currency.EUR)

	if ok, err := sgd.Equals(eur); err == nil || ok {
		t.Error("Equals() expected mismatch error")
	}

	if ok, err := sgd.GreaterThan(eur); err == nil || ok {
		t.Error("GreaterThan() expected mismatch error")
	}

	if ok, err := sgd.GreaterThanOrEqual(eur); err == nil || ok {
		t.Error("GreaterThanOrEqual() expected mismatch error")
	}

	if ok, err := sgd.LessThan(eur); err == nil || ok {
		t.Error("LessThan() expected mismatch error")
	}

	if ok, err := sgd.LessThanOrEqual(eur); err == nil || ok {
		t.Error("LessThanOrEqual() expected mismatch error")
	}
}

func TestMoneyNilReceiverErrors(t *testing.T) {
	var m *Money
	manager := NewManager()
	tests := []struct {
		name string
		fn   func() error
	}{

		{name: "Currency", fn: func() error { _, err := m.Currency(); return err }},

		{name: "Amount", fn: func() error { _, err := m.Amount(); return err }},

		{name: "IsZero", fn: func() error { _, err := m.IsZero(); return err }},

		{name: "IsPositive", fn: func() error { _, err := m.IsPositive(); return err }},

		{name: "IsNegative", fn: func() error { _, err := m.IsNegative(); return err }},

		{name: "Absolute", fn: func() error { _, err := manager.Absolute(m); return err }},

		{name: "Negative", fn: func() error { _, err := manager.Negative(m); return err }},

		{name: "Add", fn: func() error { _, err := manager.Add(m); return err }},

		{name: "Subtract", fn: func() error { _, err := manager.Subtract(m); return err }},

		{name: "Multiply", fn: func() error { _, err := manager.Multiply(m, 2); return err }},

		{name: "Round", fn: func() error { _, err := manager.Round(m); return err }},

		{name: "Split", fn: func() error { _, err := manager.Split(m, 2); return err }},

		{name: "Allocate", fn: func() error { _, err := manager.Allocate(m, 1); return err }},

		{name: "Display", fn: func() error { _, err := m.Display(); return err }},

		{name: "AsMajorUnits", fn: func() error { _, err := m.AsMajorUnits(); return err }},

		{name: "MarshalJSON", fn: func() error { _, err := m.MarshalJSON(); return err }},
		{name: "UnmarshalJSON", fn: func() error { return m.UnmarshalJSON([]byte("{}")) }},

		{name: "Value", fn: func() error { _, err := m.Value(); return err }},
		{name: "AssertSameCurrency", fn: func() error { return m.AssertSameCurrency(testManager.Create(1, currency.SGD)) }},

		{name: "SameCurrency", fn: func() error { _, err := m.SameCurrency(testManager.Create(1, currency.SGD)); return err }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.fn(); !errors.Is(err, exception.ErrNoMoneyProvided) {
				t.Fatalf("%s error = %v, want ErrNoMoneyProvided", tt.name, err)
			}
		})
	}
}

func TestMoneyNilOther(t *testing.T) {
	base := testManager.Create(10, currency.SGD)

	if err := base.AssertSameCurrency(nil); !errors.Is(err, exception.ErrNoMoneyProvided) {
		t.Fatalf("AssertSameCurrency(nil) error = %v, want ErrNoMoneyProvided", err)
	}

	if ok, err := base.SameCurrency(nil); err == nil || ok {
		t.Fatalf("SameCurrency(nil) = (%v, %v), want (false, error)", ok, err)
	}
}

func TestMoneyScanNilReceiver(t *testing.T) {
	var m *Money

	if err := m.Scan("10|SGD"); err == nil {
		t.Fatal("expected error when scanning into nil Money")
	}
}

func TestMoneyCompareNilReceiver(t *testing.T) {
	var m *Money

	if _, err := m.Compare(testManager.Create(1, currency.SGD)); !errors.Is(err, exception.ErrNoMoneyProvided) {
		t.Fatalf("Compare on nil receiver error = %v, want ErrNoMoneyProvided", err)
	}
}

func TestMoneyMarshalUnmarshalErrorPaths(t *testing.T) {
	m := testManager.Create(10, currency.SGD)

	if err := m.UnmarshalJSON([]byte(`{"amount":10,"currency":"INVALID"}`)); err == nil {
		t.Fatal("expected error for invalid currency code during unmarshal")
	}
}

func TestMoneyConverterErrorPaths(t *testing.T) {
	ex := exchange.NewExchange() // no rates -> invalid
	converter := &Converter{currencies: currency.NewManager(), exchange: ex}

	if _, err := converter.Convert(testManager.Create(10, currency.SGD), "EUR"); !errors.Is(err, exception.ErrInvalidExchangeRate) {
		t.Fatalf("Convert with invalid exchange error = %v, want ErrInvalidExchangeRate", err)
	}

	if _, err := converter.ConvertWithRate(testManager.Create(10, currency.SGD), "EUR", -1); !errors.Is(err, exception.ErrInvalidExchangeRate) {
		t.Fatalf("ConvertWithRate with invalid exchange error = %v, want ErrInvalidExchangeRate", err)
	}

	var nilConverter *Converter

	if _, err := nilConverter.Convert(testManager.Create(1, currency.SGD), currency.EUR); !errors.Is(err, exception.ErrInvalidExchangeRate) {
		t.Fatalf("nil converter Convert error = %v, want ErrInvalidExchangeRate", err)
	}

	if _, err := nilConverter.ConvertWithRate(testManager.Create(1, currency.SGD), currency.EUR, 1); !errors.Is(err, exception.ErrInvalidExchangeRate) {
		t.Fatalf("nil converter ConvertWithRate error = %v, want ErrInvalidExchangeRate", err)
	}

	if _, err := NewConverter(nil, exchange.NewExchange()); err == nil {
		t.Fatal("expected error when currencies manager is nil")
	}

	invalidExchange := (*exchange.Exchange)(nil)

	if _, err := NewConverter(currency.NewManager(), invalidExchange); !errors.Is(err, exception.ErrInvalidExchangeRate) {
		t.Fatalf("NewConverter with invalid exchange error = %v, want ErrInvalidExchangeRate", err)
	}

	exWithRate := exchange.NewExchange()
	testutil.TestRequireNoErr(t, exWithRate.AddRate(currency.SGD, currency.EUR, 0.5))
	validConverter, err := NewConverter(currency.NewManager(), exWithRate)
	testutil.TestRequireNoErr(t, err)

	if _, err := validConverter.Convert(testManager.Create(5, currency.SGD), "ZZZ"); err == nil {
		t.Fatal("expected error for unknown target currency")
	}
}

func TestJSONNilParserAndSetters(t *testing.T) {
	var parser *JSON

	if _, err := parser.Marshal(Money{}); !errors.Is(err, exception.ErrNoJSONParserProvided) {
		t.Fatalf("Marshal nil parser error = %v, want ErrNoJSONParserProvided", err)
	}

	if err := parser.Unmarshal(&Money{}, []byte(`{}`)); !errors.Is(err, exception.ErrNoJSONParserProvided) {
		t.Fatalf("Unmarshal nil parser error = %v, want ErrNoJSONParserProvided", err)
	}

	parser = NewJson()

	if err := parser.SetMarshal(nil); err == nil {
		t.Fatal("expected error when setting nil marshal function")
	}

	if err := parser.SetUnmarshal(nil); err == nil {
		t.Fatal("expected error when setting nil unmarshal function")
	}

	var nilParser *JSON

	if err := nilParser.SetUnmarshal(func(*Money, []byte) error { return nil }); !errors.Is(err, exception.ErrNoJSONParserProvided) {
		t.Fatalf("SetUnmarshal nil parser error = %v, want ErrNoJSONParserProvided", err)
	}

	if err := nilParser.SetMarshal(func(Money) ([]byte, error) { return nil, nil }); !errors.Is(err, exception.ErrNoJSONParserProvided) {
		t.Fatalf("SetMarshal nil parser error = %v, want ErrNoJSONParserProvided", err)
	}

	if err := nilParser.SetCurrency(func() (*currency.Currency, error) { return nil, nil }); !errors.Is(err, exception.ErrNoJSONParserProvided) {
		t.Fatalf("SetCurrency nil parser error = %v, want ErrNoJSONParserProvided", err)
	}
}

func TestJSONDefaultCurrencyAndInvalidCode(t *testing.T) {
	parser := NewJson()

	if _, err := parser.getUnmarshalJSONMoney(10, JSONRawData{Currency: "XXX_INVALID"}); err == nil {
		t.Fatal("expected error for invalid currency code")
	}

	var nilParser *JSON

	if _, err := nilParser.defaultJSONCurrency(); !errors.Is(err, exception.ErrNoJSONParserProvided) {
		t.Fatalf("defaultJSONCurrency nil parser error = %v, want ErrNoJSONParserProvided", err)
	}

	if _, err := nilParser.getUnmarshalJSONMoney(0, JSONRawData{}); !errors.Is(err, exception.ErrNoJSONParserProvided) {
		t.Fatalf("getUnmarshalJSONMoney nil parser error = %v, want ErrNoJSONParserProvided", err)
	}
}

func TestMoneyMarshalJSONErrorPropagation(t *testing.T) {
	parser := NewJsonWithParser(
		nil,
		func(Money) ([]byte, error) { return nil, errors.New("boom") },
		nil,
	)

	m := testManager.Create(1, currency.SGD)
	parser.SetMarshal(parser.defaultMarshalJSON) // ensure parser is usable

	if err := parser.Unmarshal(m, []byte(`{"amount": 1, "currency": "SGD"}`)); err != nil {
		t.Fatalf("unexpected setup error: %v", err)
	}

	parser.SetMarshal(func(Money) ([]byte, error) { return nil, errors.New("marshal failure") })

	if _, err := parser.Marshal(*m); err == nil {
		t.Fatal("expected marshal error to propagate")
	}
}

func TestMoneyUnmarshalJSONInvalidNumberLiteral(t *testing.T) {
	var m Money
	err := m.UnmarshalJSON([]byte(`{"amount": abc, "currency": "SGD"}`))

	if !errors.Is(err, exception.ErrInvalidJSONUnmarshal) {
		t.Fatalf("expected ErrInvalidJSONUnmarshal, got %v", err)
	}
}

func TestJSONDefaultMarshalUnmarshalPaths(t *testing.T) {
	var parser *JSON

	if _, err := parser.defaultMarshalJSON(Money{}); !errors.Is(err, exception.ErrNoJSONParserProvided) {
		t.Fatalf("defaultMarshalJSON nil parser error = %v, want ErrNoJSONParserProvided", err)
	}

	parser = NewJson()

	var m Money

	if err := parser.defaultUnmarshalJSON(&m, []byte(`{"amount": 1.2.3, "currency": "SGD"}`)); !errors.Is(err, exception.ErrInvalidJSONUnmarshal) {
		t.Fatalf("defaultUnmarshalJSON invalid number literal error = %v, want ErrInvalidJSONUnmarshal", err)
	}

	if err := parser.defaultUnmarshalJSON(&m, []byte(`{"amount": NaN, "currency": "SGD"}`)); !errors.Is(err, exception.ErrInvalidJSONUnmarshal) {
		t.Fatalf("defaultUnmarshalJSON unsupported value error = %v, want ErrInvalidJSONUnmarshal", err)
	}

	if err := parser.defaultUnmarshalJSON(&m, []byte(`{"amount": 10, "currency": "SGD"}`)); err != nil {
		t.Fatalf("defaultUnmarshalJSON valid path error = %v", err)
	}

	var nilParser *JSON

	if err := nilParser.defaultUnmarshalJSON(&m, []byte(`{"amount": 10, "currency": "SGD"}`)); !errors.Is(err, exception.ErrNoJSONParserProvided) {
		t.Fatalf("defaultUnmarshalJSON nil parser error = %v, want ErrNoJSONParserProvided", err)
	}
}

func TestConvertWithRateUnknownCurrency(t *testing.T) {
	ex := exchange.NewExchange()
	testutil.TestRequireNoErr(t, ex.AddRate(currency.SGD, currency.EUR, 1))

	converter := newTestConverter(t, currency.NewManager(), ex)

	if _, err := converter.ConvertWithRate(testManager.Create(10, currency.SGD), "ZZZ", 1); err == nil {
		t.Fatal("expected error for unknown target currency")
	}
}

func TestAssertSameCurrencySuccess(t *testing.T) {
	a := testManager.Create(10, currency.SGD)
	b := testManager.Create(20, currency.SGD)

	if err := a.AssertSameCurrency(b); err != nil {
		t.Fatalf("expected nil error for same currency, got %v", err)
	}
}

func TestJSONCurrencyFunctionError(t *testing.T) {
	parser := NewJson()

	testutil.TestRequireNoErr(t, parser.SetCurrency(func() (*currency.Currency, error) {
		return nil, errors.New("boom")
	}))

	var m Money

	if err := parser.defaultUnmarshalJSON(&m, []byte(`{"amount": 10, "currency": "SGD"}`)); err == nil {
		t.Fatal("expected error from currency function")
	}
}

func TestMoney_Display_Coverage(t *testing.T) {
	t.Parallel()

	m := &Money{amount: 100, currency: nil}

	if got, err := m.Display(); got != "" || err == nil {
		t.Errorf("Display() with nil currency = (%q, %v), want empty string and error", got, err)
	}
}

func TestMoney_AsMajorUnits_Coverage(t *testing.T) {
	t.Parallel()

	m := &Money{amount: 100, currency: nil}

	if got, err := m.AsMajorUnits(); got != 0 || err == nil {
		t.Errorf("AsMajorUnits() with nil currency = (%f, %v), want 0 and error", got, err)
	}
}

func TestMoney_AssertSameCurrency_Coverage(t *testing.T) {
	t.Parallel()

	m1 := testManager.Create(100, currency.SGD)

	// Test with nil other money
	if err := m1.AssertSameCurrency(nil); err == nil {
		t.Error("AssertSameCurrency(nil) should return error")
	}

	// Test with other money having nil currency
	m2 := &Money{amount: 100, currency: nil}

	if err := m1.AssertSameCurrency(m2); err == nil {
		t.Error("AssertSameCurrency(money with nil currency) should return error")
	}
}
