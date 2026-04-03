package parser

import (
	"errors"
	"math"
	"reflect"
	"regexp"
	"testing"
	"unsafe"

	"github.com/gollin/packages/support/money/currency"
	"github.com/gollin/packages/support/money/exception"
)

func TestParser_ParseAmount(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		defaultCurrency string
		expectedVal     float64
		expectedCurr    string
		expectError     bool
	}{
		// Basic US Format
		{"Simple SGD", "S$100.50", "", 100.50, currency.SGD, false},
		{"Suffix SGD", "100.50 SGD", "", 100.50, currency.SGD, false},
		{"Prefix SGD", "SGD 100.50", "", 100.50, currency.SGD, false},
		{"Comma Thousand", "S$1,234.56", "", 1234.56, currency.SGD, false},

		// European Format (Mixed Separators)
		{"Euro Symbol", "€100.50", "", 100.50, currency.EUR, false},
		{"Euro Mixed Sep", "€1.234,56", "", 1234.56, currency.EUR, false},
		{"Euro Code", "1.234,56 EUR", "", 1234.56, currency.EUR, false},

		// Complex Symbols (The "C$" vs "$" Bug)
		{"Canadian Dollar", "C$150.00", "", 150.00, currency.CAD, false},
		{"Australian Dollar", "A$200", "", 200.00, currency.AUD, false},
		{"Brazilian Real", "R$ 50,000.00", "", 50000.00, currency.BRL, false},
		{"Taiwan Dollar", "NT$100", "", 100.00, currency.TWD, false},
		{"CFP Franc", "CFP 100", "", 100.00, currency.XPF, false},

		// Negatives
		{"Negative Symbol", "-S$50.00", "", -50.00, currency.SGD, false},
		{"Negative Code", "-100 SGD", "", -100.00, currency.SGD, false},

		// Defaults
		{"Use Provider", "10.00", currency.GBP, 10.00, currency.GBP, false},
		{"Override Provider", "S$10.00", "GBP", 10.00, currency.SGD, false},

		// Edge Cases
		{"Whitespace", "  S$  100.00  ", "", 100.00, currency.SGD, false},
		{"Sentence With Uppercase Word", "I PAY you 100 SGD", "", 0, "", true},
		{"No Currency", "100.00", "", 0, "", true},
		{"Invalid Number", "$abc", "", 0, "", true},
		{"Empty String", "", "", 0, "", true},
	}

	parser := NewParser()
	const epsilon = 1e-9

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, curr, err := parser.ParseAmount(tt.input, tt.defaultCurrency)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}

				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)

				return
			}

			if curr != tt.expectedCurr {
				t.Errorf("Expected currency %s, got %s", tt.expectedCurr, curr)
			}

			if math.Abs(val-tt.expectedVal) > epsilon {
				t.Errorf("Expected value %f, got %f", tt.expectedVal, val)
			}
		})
	}
}

func TestParseDecimal(t *testing.T) {
	tests := []struct {
		input       string
		expected    float64
		expectError bool
	}{
		{"100", 100, false},
		{"100.50", 100.50, false},
		{"1,000.50", 1000.50, false},
		{"1.000,50", 1000.50, false}, // Euro check
		{"-50.00", -50.00, false},
		{"abc", 0, true},   // Malformed number case
		{"1.2,3", 0, true}, // Malformed mixed European grouping
		{"1,2.3", 0, true}, // Malformed mixed US grouping
	}

	const epsilon = 1e-9

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			p := NewParser()
			val, err := p.ParseDecimal(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for %s but got nil", tt.input)
				}

				return
			}

			if err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.input, err)

				return
			}

			if math.Abs(val-tt.expected) > epsilon {
				t.Errorf("ParseDecimal(%s): expected %f, got %f", tt.input, tt.expected, val)
			}
		})
	}
}

// TestCommaOnlyAmbiguity documents the default comma-only parsing behavior
// Comma-only inputs are treated as thousands separators by default
func TestCommaOnlyAmbiguity(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"European decimal treated as thousands", "10,50", 1050.00},
		{"Thousands separator", "1,000", 1000.00},
		{"Multiple thousands", "1,000,000", 1000000.00},
		{"Small European decimal", "5,99", 599.00},
	}

	const epsilon = 1e-9

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			val, err := p.ParseDecimal(tt.input)

			if err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.input, err)

				return
			}

			if math.Abs(val-tt.expected) > epsilon {
				t.Errorf("ParseDecimal(%s): expected %f, got %f (default treats comma as thousands separator)",
					tt.input, tt.expected, val)
			}
		})
	}
}

// TestParseDecimalWithComma tests European decimal comma parsing
func TestParseDecimalWithComma(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    float64
		expectError bool
	}{
		{"European decimal", "10,50", 10.50, false},
		{"European decimal small", "5,99", 5.99, false},
		{"European decimal many decimals", "100,123", 100.123, false},
		{"Whole number with comma", "1000,00", 1000.00, false},
		{"Just decimals", "0,50", 0.50, false},
		{"Negative European", "-10,50", -10.50, false},
		{"Mixed separators European", "1.234,56", 1234.56, false}, // Should still work
		{"Mixed separators US", "1,234.56", 1234.56, false},       // Should still detect correctly
		{"No separator", "100", 100.00, false},
		{"Dot only", "100.50", 100.50, false},
		{"Invalid", "abc", 0, true},
	}

	const epsilon = 1e-9

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser()
			val, err := p.ParseDecimalWithComma(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for %s but got nil", tt.input)
				}

				return
			}

			if err != nil {
				t.Errorf("Unexpected error for %s: %v", tt.input, err)

				return
			}

			if math.Abs(val-tt.expected) > epsilon {
				t.Errorf("ParseDecimalWithComma(%s): expected %f, got %f", tt.input, tt.expected, val)
			}
		})
	}
}

// TestParseAmountWithDecimalComma tests European decimal comma parsing with currency
func TestParseAmountWithDecimalComma(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		defaultCurrency string
		expectedVal     float64
		expectedCurr    string
		expectError     bool
	}{
		// European decimal comma cases
		{"Euro with comma decimal", "€10,50", "", 10.50, currency.EUR, false},
		{"Euro with comma and code", "10,50 EUR", "", 10.50, currency.EUR, false},
		{"SGD with comma decimal", "S$99,99", "", 99.99, currency.SGD, false},
		{"GBP with comma decimal", "£5,50", "", 5.50, currency.GBP, false},

		// Mixed separators should still work
		{"Euro mixed format", "€1.234,56", "", 1234.56, currency.EUR, false},
		{"SGD standard format", "S$1,234.56", "", 1234.56, currency.SGD, false},

		// Provider currency
		{"With default currency", "10,50", currency.GBP, 10.50, currency.GBP, false},

		// Negative
		{"Negative with comma", "-€10,50", "", -10.50, currency.EUR, false},

		// Edge cases
		{"Whitespace", "  € 10,50  ", "", 10.50, currency.EUR, false},
		{"No currency", "10,50", "", 0, "", true},
		{"Invalid number", "€abc", "", 0, "", true},
		{"Empty", "", "", 0, "", true},
	}

	parser := NewParser()
	const epsilon = 1e-9

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, curr, err := parser.ParseAmountWithDecimalComma(tt.input, tt.defaultCurrency)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}

				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)

				return
			}

			if curr != tt.expectedCurr {
				t.Errorf("Expected currency %s, got %s", tt.expectedCurr, curr)
			}

			if math.Abs(val-tt.expectedVal) > epsilon {
				t.Errorf("Expected value %f, got %f", tt.expectedVal, val)
			}
		})
	}
}

func TestValidThousandsGrouping(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		sep         string
		decimal     string
		expectValid bool
	}{
		{name: "no separator", input: "123", sep: ",", decimal: ".", expectValid: true},
		{name: "valid grouping", input: "1,234,567.89", sep: ",", decimal: ".", expectValid: true},
		{name: "invalid leading group", input: "1234,567", sep: ",", decimal: ".", expectValid: false},
		{name: "invalid middle group size", input: "12,34,567", sep: ",", decimal: ".", expectValid: false},
		{name: "empty integer part", input: ".99", sep: ",", decimal: ".", expectValid: false},
		{name: "signed valid grouping", input: "-1,234", sep: ",", decimal: ".", expectValid: true},
		{name: "decimal with invalid grouping", input: "1.23,456", sep: ".", decimal: ",", expectValid: false},
	}

	p := NewParser()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.validThousandsGrouping(tt.input, tt.sep, tt.decimal)

			if err != nil {
				t.Fatalf("validThousandsGrouping(%q) unexpected error: %v", tt.input, err)
			}

			if got != tt.expectValid {
				t.Fatalf("validThousandsGrouping(%q) = %v, want %v", tt.input, got, tt.expectValid)
			}
		})
	}
}

func TestParseStringSign(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantAmount   string
		wantNegative bool
	}{
		{
			name:         "positive number",
			input:        "123.45",
			wantAmount:   "123.45",
			wantNegative: false,
		},
		{
			name:         "negative number",
			input:        "-123.45",
			wantAmount:   "123.45",
			wantNegative: true,
		},
		{
			name:         "explicit positive",
			input:        "+123.45",
			wantAmount:   "123.45",
			wantNegative: false,
		},
		{
			name:         "no sign",
			input:        "999",
			wantAmount:   "999",
			wantNegative: false,
		},
		{
			name:         "negative zero",
			input:        "-0",
			wantAmount:   "0",
			wantNegative: true,
		},
		{
			name:         "positive zero",
			input:        "+0",
			wantAmount:   "0",
			wantNegative: false,
		},
	}

	p := NewParser()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAmount, gotNegative, err := p.ParseStringSign(tt.input)

			if err != nil {
				t.Fatalf("ParseStringSign(%q) unexpected error: %v", tt.input, err)
			}

			if gotAmount != tt.wantAmount {
				t.Errorf("ParseStringSign(%q) amount = %q, want %q", tt.input, gotAmount, tt.wantAmount)
			}

			if gotNegative != tt.wantNegative {
				t.Errorf("ParseStringSign(%q) negative = %v, want %v", tt.input, gotNegative, tt.wantNegative)
			}
		})
	}
}

func TestParseDecimalParts(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantInteger string
		wantDecimal string
		wantErr     error
	}{
		{
			name:        "whole number",
			input:       "123",
			wantInteger: "123",
			wantDecimal: "",
			wantErr:     nil,
		},
		{
			name:        "decimal number",
			input:       "123.45",
			wantInteger: "123",
			wantDecimal: "45",
			wantErr:     nil,
		},
		{
			name:        "zero with decimal",
			input:       "0.99",
			wantInteger: "0",
			wantDecimal: "99",
			wantErr:     nil,
		},
		{
			name:        "empty integer part",
			input:       ".99",
			wantInteger: "0",
			wantDecimal: "99",
			wantErr:     nil,
		},
		{
			name:        "trailing decimal point",
			input:       "123.",
			wantInteger: "123",
			wantDecimal: "",
			wantErr:     nil,
		},
		{
			name:        "multiple decimal points",
			input:       "123.45.67",
			wantInteger: "",
			wantDecimal: "",
			wantErr:     exception.ErrInvalidAmountMultiple,
		},
		{
			name:        "single digit",
			input:       "5",
			wantInteger: "5",
			wantDecimal: "",
			wantErr:     nil,
		},
		{
			name:        "large number",
			input:       "1234567.89",
			wantInteger: "1234567",
			wantDecimal: "89",
			wantErr:     nil,
		},
	}

	p := NewParser()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotInteger, gotDecimal, gotErr := p.ParseDecimalParts(tt.input)

			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("ParseDecimalParts(%q) error = %v, want %v", tt.input, gotErr, tt.wantErr)

				return
			}

			if gotErr == nil {
				if gotInteger != tt.wantInteger {
					t.Errorf("ParseDecimalParts(%q) integer = %q, want %q", tt.input, gotInteger, tt.wantInteger)
				}

				if gotDecimal != tt.wantDecimal {
					t.Errorf("ParseDecimalParts(%q) decimal = %q, want %q", tt.input, gotDecimal, tt.wantDecimal)
				}
			}
		})
	}
}

func TestValidateAndPadDecimal(t *testing.T) {
	tests := []struct {
		name        string
		decimalPart string
		fraction    int
		want        string
		wantErr     error
	}{
		{
			name:        "exact match",
			decimalPart: "99",
			fraction:    2,
			want:        "99",
			wantErr:     nil,
		},
		{
			name:        "needs padding",
			decimalPart: "5",
			fraction:    2,
			want:        "50",
			wantErr:     nil,
		},
		{
			name:        "empty needs padding",
			decimalPart: "",
			fraction:    2,
			want:        "00",
			wantErr:     nil,
		},
		{
			name:        "too many decimals",
			decimalPart: "999",
			fraction:    2,
			want:        "",
			wantErr:     exception.ErrInvalidAmountFraction,
		},
		{
			name:        "no fraction currency",
			decimalPart: "",
			fraction:    0,
			want:        "",
			wantErr:     nil,
		},
		{
			name:        "decimal on no-fraction currency",
			decimalPart: "5",
			fraction:    0,
			want:        "",
			wantErr:     exception.ErrInvalidAmountFraction,
		},
		{
			name:        "three decimal places",
			decimalPart: "123",
			fraction:    3,
			want:        "123",
			wantErr:     nil,
		},
		{
			name:        "pad to three decimals",
			decimalPart: "12",
			fraction:    3,
			want:        "120",
			wantErr:     nil,
		},
	}

	p := NewParser()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.ValidateAndPadDecimal(tt.decimalPart, tt.fraction)

			if err != tt.wantErr {
				t.Errorf("ValidateAndPadDecimal(%q, %d) error = %v, want %v", tt.decimalPart, tt.fraction, err, tt.wantErr)

				return
			}

			if got != tt.want {
				t.Errorf("ValidateAndPadDecimal(%q, %d) = %q, want %q", tt.decimalPart, tt.fraction, got, tt.want)
			}
		})
	}
}

func TestParseAmountString(t *testing.T) {
	tests := []struct {
		name     string
		amount   string
		fraction int
		negative bool
		want     int64
		wantErr  error
	}{
		{
			name:     "positive SGD",
			amount:   "12.34",
			fraction: 2,
			negative: false,
			want:     1234,
			wantErr:  nil,
		},
		{
			name:     "negative SGD",
			amount:   "12.34",
			fraction: 2,
			negative: true,
			want:     -1234,
			wantErr:  nil,
		},
		{
			name:     "whole number SGD",
			amount:   "100",
			fraction: 2,
			negative: false,
			want:     10000,
			wantErr:  nil,
		},
		{
			name:     "single decimal SGD",
			amount:   "10.5",
			fraction: 2,
			negative: false,
			want:     1050,
			wantErr:  nil,
		},
		{
			name:     "JPY no fraction",
			amount:   "100",
			fraction: 0,
			negative: false,
			want:     100,
			wantErr:  nil,
		},
		{
			name:     "BHD three decimals",
			amount:   "10.123",
			fraction: 3,
			negative: false,
			want:     10123,
			wantErr:  nil,
		},
		{
			name:     "zero amount",
			amount:   "0.00",
			fraction: 2,
			negative: false,
			want:     0,
			wantErr:  nil,
		},
		{
			name:     "negative zero",
			amount:   "0.00",
			fraction: 2,
			negative: true,
			want:     0,
			wantErr:  nil,
		},
		{
			name:     "too many decimals",
			amount:   "10.999",
			fraction: 2,
			negative: false,
			want:     0,
			wantErr:  exception.ErrInvalidAmountFraction,
		},
		{
			name:     "multiple decimal points",
			amount:   "10.12.34",
			fraction: 2,
			negative: false,
			want:     0,
			wantErr:  exception.ErrInvalidAmountMultiple,
		},
		{
			name:     "empty integer part",
			amount:   ".99",
			fraction: 2,
			negative: false,
			want:     99,
			wantErr:  nil,
		},
		{
			name:     "large amount",
			amount:   "999999.99",
			fraction: 2,
			negative: false,
			want:     99999999,
			wantErr:  nil,
		},
		{
			name:     "invalid characters",
			amount:   "12.ab",
			fraction: 2,
			negative: false,
			want:     0,
			wantErr:  exception.ErrInvalidAmount,
		},
	}

	p := NewParser()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := p.ParseAmountString(tt.amount, tt.fraction, tt.negative)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ParseAmountString(%q, %d, %v) error = %v, want %v", tt.amount, tt.fraction, tt.negative, err, tt.wantErr)

				return
			}

			if got != tt.want {
				t.Errorf("ParseAmountString(%q, %d, %v) = %d, want %d", tt.amount, tt.fraction, tt.negative, got, tt.want)
			}
		})
	}
}

func TestParser_NewParserWithAndNilReceiverPaths(t *testing.T) {
	p := NewParserWith(currency.NewISOCodePattern())

	if p == nil {
		t.Fatal("NewParserWith() returned nil")
	}

	_, _, err := p.ParseAmount("SGD 10.00")

	if err != nil {
		t.Fatalf("ParseAmount() unexpected error: %v", err)
	}

	var nilParser *Parser

	if _, _, err := nilParser.ParseAmount("SGD 10"); !errors.Is(err, exception.ErrParserNotProvided) {
		t.Fatalf("ParseAmount(nil) error = %v, want ErrParserNotProvided", err)
	}

	if _, _, err := nilParser.ParseAmountWithDecimalComma("SGD 10"); !errors.Is(err, exception.ErrParserNotProvided) {
		t.Fatalf("ParseAmountWithDecimalComma(nil) error = %v, want ErrParserNotProvided", err)
	}
}

func TestParser_HelperCoveragePaths(t *testing.T) {
	p := NewParser()

	t.Run("validThousandsGrouping nil receiver", func(t *testing.T) {
		var nilParser *Parser
		_, err := nilParser.validThousandsGrouping("1,234.56", ",", ".")

		if !errors.Is(err, exception.ErrParserNotProvided) {
			t.Fatalf("validThousandsGrouping(nil) error = %v, want ErrParserNotProvided", err)
		}
	})

	t.Run("parseNumericString nil receiver", func(t *testing.T) {
		var nilParser *Parser
		_, err := nilParser.parseNumericString("10", false)

		if !errors.Is(err, exception.ErrParserNotProvided) {
			t.Fatalf("parseNumericString(nil) error = %v, want ErrParserNotProvided", err)
		}
	})

	t.Run("validThousandsGrouping empty integer part", func(t *testing.T) {
		ok, err := p.validThousandsGrouping(".12", ",", ".")

		if err != nil || ok {
			t.Fatalf("validThousandsGrouping(.12) = (%v,%v), want (false,nil)", ok, err)
		}
	})

	t.Run("extractCurrency invalid state", func(t *testing.T) {
		bad := &Parser{iso: nil}
		_, _, err := bad.extractCurrency("$10")

		if !errors.Is(err, exception.ErrParserInvalidState) {
			t.Fatalf("extractCurrency() error = %v, want ErrParserInvalidState", err)
		}
	})

	t.Run("extractCurrency len(match) < 4", func(t *testing.T) {
		iso := currency.NewISOCodePattern()

		if iso == nil {
			t.Fatal("currency.NewISOCodePattern() returned nil")
		}

		_ = iso.GetPattern() // initialize once

		// Replace the cached regexp with one that has no capture groups so
		// FindAllStringSubmatchIndex returns a match slice of length 2.
		re := regexp.MustCompile("\\bUSD\\b")
		rv := reflect.ValueOf(iso).Elem()
		f := rv.FieldByName("isoCodePattern")
		reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem().Set(reflect.ValueOf(re))

		pp := NewParserWith(iso)
		remaining, curr, err := pp.extractCurrency("USD 10", currency.SGD)

		if err != nil {
			t.Fatalf("extractCurrency() unexpected error: %v", err)
		}

		if curr != currency.SGD {
			t.Fatalf("currency = %q, want %q", curr, currency.SGD)
		}

		if remaining != "USD 10" {
			t.Fatalf("remaining = %q, want %q", remaining, "USD 10")
		}
	})

	t.Run("parseNumericString comma-only modes", func(t *testing.T) {
		val, err := p.parseNumericString("10,50", false)

		if err != nil || val != 1050 {
			t.Fatalf("parseNumericString provider mode = (%v,%v), want (1050,nil)", val, err)
		}

		val, err = p.parseNumericString("10,50", true)

		if err != nil || val != 10.50 {
			t.Fatalf("parseNumericString decimal comma mode = (%v,%v), want (10.50,nil)", val, err)
		}
	})

	t.Run("ParseStringSign plus", func(t *testing.T) {
		amount, negative, err := p.ParseStringSign("+12.34")

		if err != nil || negative || amount != "12.34" {
			t.Fatalf("ParseStringSign(+12.34) = (%q,%v,%v), want (\"12.34\",false,nil)", amount, negative, err)
		}
	})

	t.Run("ParseStringSign nil receiver", func(t *testing.T) {
		var nilParser *Parser
		_, _, err := nilParser.ParseStringSign("1")

		if !errors.Is(err, exception.ErrParserNotProvided) {
			t.Fatalf("ParseStringSign(nil) error = %v, want ErrParserNotProvided", err)
		}
	})

	t.Run("ParseDecimalParts errors and defaults", func(t *testing.T) {
		if _, _, err := p.ParseDecimalParts("1.2.3"); !errors.Is(err, exception.ErrInvalidAmountMultiple) {
			t.Fatalf("ParseDecimalParts error = %v, want ErrInvalidAmountMultiple", err)
		}

		ip, dp, err := p.ParseDecimalParts(".50")

		if err != nil || ip != "0" || dp != "50" {
			t.Fatalf("ParseDecimalParts(.50) = (%q,%q,%v), want (\"0\",\"50\",nil)", ip, dp, err)
		}
	})

	t.Run("ParseDecimalParts nil receiver", func(t *testing.T) {
		var nilParser *Parser
		_, _, err := nilParser.ParseDecimalParts("1")

		if !errors.Is(err, exception.ErrParserNotProvided) {
			t.Fatalf("ParseDecimalParts(nil) error = %v, want ErrParserNotProvided", err)
		}
	})

	t.Run("ValidateAndPadDecimal errors and padding", func(t *testing.T) {
		if _, err := p.ValidateAndPadDecimal("123", 2); !errors.Is(err, exception.ErrInvalidAmountFraction) {
			t.Fatalf("ValidateAndPadDecimal error = %v, want ErrInvalidAmountFraction", err)
		}

		got, err := p.ValidateAndPadDecimal("5", 2)

		if err != nil || got != "50" {
			t.Fatalf("ValidateAndPadDecimal(5,2) = (%q,%v), want (\"50\",nil)", got, err)
		}
	})

	t.Run("ValidateAndPadDecimal nil receiver", func(t *testing.T) {
		var nilParser *Parser
		_, err := nilParser.ValidateAndPadDecimal("1", 2)

		if !errors.Is(err, exception.ErrParserNotProvided) {
			t.Fatalf("ValidateAndPadDecimal(nil) error = %v, want ErrParserNotProvided", err)
		}
	})

	t.Run("ParseAmountString parse-int error", func(t *testing.T) {
		if _, err := p.ParseAmountString("12.a", 2, false); !errors.Is(err, exception.ErrInvalidAmount) {
			t.Fatalf("ParseAmountString error = %v, want ErrInvalidAmount", err)
		}
	})

	t.Run("ParseAmountString nil receiver", func(t *testing.T) {
		var nilParser *Parser
		_, err := nilParser.ParseAmountString("1", 2, false)

		if !errors.Is(err, exception.ErrParserNotProvided) {
			t.Fatalf("ParseAmountString(nil) error = %v, want ErrParserNotProvided", err)
		}
	})
}
