package format

import (
	"math"
	"strconv"
	"strings"
)

// Formatter handles the formatting of monetary amounts.
type Formatter struct {
	Fraction int    // Number of decimal places (e.g., 2 for SGD, 0 for JPY)
	Decimal  string // The character separating whole numbers from decimals (e.g., ".")
	Thousand string // The character separating thousands (e.g., ",")
	Grapheme string // The currency symbol (e.g. "$", "£", "€")
	Template string // The pattern for the output (e.g. "$1")
}

// NewFormatter creates a new Formatter instance with the specified configuration.
func NewFormatter(fraction int, decimal, thousand, grapheme, template string) *Formatter {
	return &Formatter{
		Fraction: fraction,
		Decimal:  decimal,
		Thousand: thousand,
		Grapheme: grapheme,
		Template: template,
	}
}

// Format returns string of formatted integer using given currency template.
func (f *Formatter) Format(amount int64) string {
	// Work with absolute amount value
	sa := strconv.FormatInt(f.abs(amount), 10)

	if len(sa) <= f.Fraction {
		sa = strings.Repeat("0", f.Fraction-len(sa)+1) + sa
	}

	if f.Thousand != "" {
		for i := len(sa) - f.Fraction - 3; i > 0; i -= 3 {
			sa = sa[:i] + f.Thousand + sa[i:]
		}
	}

	if f.Fraction > 0 {
		sa = sa[:len(sa)-f.Fraction] + f.Decimal + sa[len(sa)-f.Fraction:]
	}

	sa = strings.Replace(f.Template, "1", sa, 1)
	sa = strings.Replace(sa, "$", f.Grapheme, 1)

	// Add minus sign for negative amount.
	if amount < 0 {
		sa = "-" + sa
	}

	return sa
}

// ToMajorUnits converts the integer amount to a float representing the major currency units.
// e.g., 100 cents becomes 1.00 dollars if fraction is 2.
func (f *Formatter) ToMajorUnits(amount int64) float64 {
	if f.Fraction == 0 {
		return float64(amount)
	}

	return float64(amount) / math.Pow10(f.Fraction)
}

func (f *Formatter) abs(amount int64) int64 {
	if amount < 0 {
		return -amount
	}

	return amount
}
