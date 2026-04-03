package parser

import (
	"strconv"
	"strings"

	"github.com/gollin/packages/support/money/currency"
	"github.com/gollin/packages/support/money/exception"
)

type Parser struct {
	iso *currency.ISOCodePattern
}

// NewParser creates a new Parser with the default ISOCodePattern.
func NewParser() *Parser {
	return &Parser{
		iso: currency.NewISOCodePattern(),
	}
}

// NewParserWith creates a new Parser with a custom ISOCodePattern.
func NewParserWith(iso *currency.ISOCodePattern) *Parser {
	return &Parser{
		iso: iso,
	}
}

// ParseAmount parses a string amount and returns the amount value and currency code.
// It intelligently handles mixed separators (e.g. "1.200,50" vs "1,200.50").
//
// Note: Comma-only inputs like "10,50" are treated as thousands of separators (parses as 1050).
// For European decimal comma notation, use ParseAmountWithDecimalComma instead.
func (p *Parser) ParseAmount(input string, defaultCurrency ...string) (float64, string, error) {
	if p == nil {
		return 0, "", exception.ErrParserNotProvided
	}

	input = strings.TrimSpace(input)

	if input == "" {
		return 0, "", exception.ErrInvalidMoneyString
	}

	// 1. Determine Currency
	var curr string
	input, curr, err := p.extractCurrency(input, defaultCurrency...)

	if err != nil {
		return 0, "", err
	}

	// 2. Parse Number
	val, err := p.parseNumericString(input, false)

	if err != nil {
		return 0, "", err
	}

	return val, curr, nil
}

// ParseAmountWithDecimalComma parses a string amount treating commas as decimal separators.
// This is useful for European-style decimal notation (e.g. "10,50" parses as 10.50).
//
// Note: This method treats comma-only inputs as decimal separators, not thousands of separators.
// Mixed separator inputs (e.g. "1.234,56") are still handled intelligently.
func (p *Parser) ParseAmountWithDecimalComma(input string, defaultCurrency ...string) (float64, string, error) {
	if p == nil {
		return 0, "", exception.ErrParserNotProvided
	}

	input = strings.TrimSpace(input)

	if input == "" {
		return 0, "", exception.ErrInvalidMoneyString
	}

	// 1. Determine Currency (same logic as ParseAmount)
	var curr string
	input, curr, err := p.extractCurrency(input, defaultCurrency...)

	if err != nil {
		return 0, "", err
	}

	// 2. Parse Number with decimal comma mode
	val, err := p.parseNumericString(input, true)

	if err != nil {
		return 0, "", err
	}

	return val, curr, nil
}

// ParseDecimal parses just a numeric amount string and returns it as a float64
// Comma-only inputs are treated as thousands of separators (e.g., "1,000" = 1000).
func (p *Parser) ParseDecimal(amount string) (float64, error) {
	return p.parseNumericString(amount, false)
}

// ParseDecimalWithComma parses a numeric amount string treating commas as decimal separators.
// This is useful for European-style decimal notation (e.g. "10,50" parses as 10.50).
func (p *Parser) ParseDecimalWithComma(amount string) (float64, error) {
	return p.parseNumericString(amount, true)
}

// parseNumericString handles standardising "1,000.00" and "1.000,00" into "1000.00"
//
// Decimal separator handling:
//   - Mixed separators (e.g. "1,234.56" or "1.234,56"): Automatically detects a format based on last separator position
//   - Comma-only (e.g. "1,000" or "10,50"): Behaviour depends on useDecimalComma parameter
//   - useDecimalComma=false: Treated as thousands separator (default) - "10,50" = 1050
//   - useDecimalComma=true: Treated as decimal separator - "10,50" = 10.50
func (p *Parser) parseNumericString(input string, useDecimalComma bool) (float64, error) {
	if p == nil {
		return 0, exception.ErrParserNotProvided
	}

	input = strings.TrimSpace(input)
	input = strings.ReplaceAll(input, " ", "")

	// Handle Mixed Separators (e.g. 1.234,56 vs 1,234.56)
	hasDot := strings.Contains(input, ".")
	hasComma := strings.Contains(input, ",")

	if hasDot && hasComma {
		lastDot := strings.LastIndex(input, ".")
		lastComma := strings.LastIndex(input, ",")

		if lastDot < lastComma {
			if valid, err := p.validThousandsGrouping(input, ".", ","); err != nil || !valid {
				return 0, exception.ErrInvalidMoneyString
			}

			// Format is 1.234,56 (European) -> Convert to 1234.56
			input = strings.ReplaceAll(input, ".", "")  // Remove thousands' separator
			input = strings.ReplaceAll(input, ",", ".") // Make comma decimal
		} else {
			if valid, err := p.validThousandsGrouping(input, ",", "."); err != nil || !valid {
				return 0, exception.ErrInvalidMoneyString
			}

			// Format is 1,234.56 (US/Standard)
			input = strings.ReplaceAll(input, ",", "")
		}
	} else if hasComma {
		// Only commas present. Behaviour depends on useDecimalComma parameter.
		if useDecimalComma {
			// European decimal mode: treat comma as decimal separator
			input = strings.ReplaceAll(input, ",", ".")
		} else {
			// Provider mode: treat comma as thousands separator
			input = strings.ReplaceAll(input, ",", "")
		}
	}

	// Parse
	amount, err := strconv.ParseFloat(input, 64)

	if err != nil {
		return 0, exception.ErrInvalidMoneyString
	}

	return amount, nil
}

// validThousandsGrouping ensures the integer portion uses 1-3 leading digits followed by 3-digit groups.
func (p *Parser) validThousandsGrouping(input, thousandsSep, decimalSep string) (bool, error) {
	if p == nil {
		return false, exception.ErrParserNotProvided
	}

	decimalIdx := strings.LastIndex(input, decimalSep)
	intPart := input

	if decimalIdx != -1 {
		intPart = input[:decimalIdx]
	}

	if len(intPart) == 0 {
		return false, nil
	}

	// Trim sign if present
	if intPart[0] == '-' || intPart[0] == '+' {
		intPart = intPart[1:]
	}

	groups := strings.Split(intPart, thousandsSep)

	if len(groups) == 1 {
		// No thousands of separators to validate
		return true, nil
	}

	if len(groups[0]) == 0 || len(groups[0]) > 3 {
		return false, nil
	}

	for _, grp := range groups[1:] {
		if len(grp) != 3 {
			return false, nil
		}
	}

	return true, nil
}

// extractCurrency encapsulates the shared logic for detecting and extracting
// currency symbols or ISO codes from the input string.
// It returns the modified input string (with currency removed), the detected currency,
// and an error if no currency is specified.
func (p *Parser) extractCurrency(input string, defaultCurrency ...string) (string, string, error) {
	if p == nil || p.iso == nil {
		return "", "", exception.ErrParserInvalidState
	}

	curr := ""

	if len(defaultCurrency) > 0 {
		curr = defaultCurrency[0]
	}

	// Check symbols (Longest first)
	for _, item := range *p.iso.GetSymbolsLongestFirst() {
		if strings.Contains(input, item.Id) {
			curr = item.Currency
			input = strings.ReplaceAll(input, item.Id, "")

			break
		}
	}

	input = strings.TrimSpace(input)

	// Check for ISO Codes (e.g. "SGD")
	// Uses a pattern that matches only valid ISO 4217 currency codes
	pattern := p.iso.GetPattern()

	if matches := pattern.FindAllStringSubmatchIndex(input, -1); len(matches) > 0 {
		for _, match := range matches {
			if len(match) < 4 {
				continue
			}

			// match[0:2] is the full match, match[2:4] is capture group 1 ("SGD")
			code := strings.TrimSpace(input[match[2]:match[3]])

			// If we found an explicit code, it overrides the symbol or default
			curr = code
			input = input[:match[0]] + input[match[1]:]

			break
		}
	}

	if curr == "" {
		return "", "", exception.ErrCurrencyNotSpecified
	}

	return input, curr, nil
}

// ParseStringSign extracts the sign from an amount string and returns the cleaned string and sign indicator.
func (p *Parser) ParseStringSign(amount string) (string, bool, error) {
	if p == nil {
		return "", false, exception.ErrParserNotProvided
	}

	negative := false

	if strings.HasPrefix(amount, "-") {
		negative = true
		amount = amount[1:]
	} else if strings.HasPrefix(amount, "+") {
		amount = amount[1:]
	}

	return amount, negative, nil
}

// ParseDecimalParts splits a decimal string into integer and decimal parts.
// Returns the integer part, decimal part, and any error.
func (p *Parser) ParseDecimalParts(amount string) (string, string, error) {
	if p == nil {
		return "", "", exception.ErrParserNotProvided
	}

	parts := strings.Split(amount, ".")

	if len(parts) > 2 {
		return "", "", exception.ErrInvalidAmountMultiple
	}

	integerPart := parts[0]

	if integerPart == "" {
		integerPart = "0"
	}

	decimalPart := ""

	if len(parts) == 2 {
		decimalPart = parts[1]
	}

	return integerPart, decimalPart, nil
}

// ValidateAndPadDecimal validates the decimal part doesn't exceed a fraction and pads it to match a currency fraction.
func (p *Parser) ValidateAndPadDecimal(decimalPart string, fraction int) (string, error) {
	if p == nil {
		return "", exception.ErrParserNotProvided
	}

	if len(decimalPart) > fraction {
		return "", exception.ErrInvalidAmountFraction
	}

	// Pad decimal part to match currency fraction
	for len(decimalPart) < fraction {
		decimalPart += "0"
	}

	return decimalPart, nil
}

// ParseAmountString parses a string amount into an int64 value according to a currency fraction.
func (p *Parser) ParseAmountString(amount string, fraction int, negative bool) (int64, error) {
	if p == nil {
		return 0, exception.ErrParserNotProvided
	}

	integerPart, decimalPart, err := p.ParseDecimalParts(amount)

	if err != nil {
		return 0, err
	}

	decimalPart, err = p.ValidateAndPadDecimal(decimalPart, fraction)

	if err != nil {
		return 0, err
	}

	// Combine integer and decimal parts
	combinedStr := integerPart + decimalPart

	// Parse as int64
	value, err := strconv.ParseInt(combinedStr, 10, 64)

	if err != nil {
		return 0, exception.ErrInvalidAmount
	}

	if negative {
		value = -value
	}

	return value, nil
}
