package support

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

// NumberFormat formats a number with grouped thousands and fixed precision.
func NumberFormat(value float64, precision ...int) string {
	digits := 0
	if len(precision) > 0 {
		digits = precision[0]
	}

	sign := ""
	if value < 0 {
		sign = "-"
		value = -value
	}

	formatted := fmt.Sprintf("%.*f", digits, value)
	parts := strings.SplitN(formatted, ".", 2)

	integer := groupThousands(parts[0])
	if len(parts) == 1 {
		return sign + integer
	}

	return sign + integer + "." + parts[1]
}

// NumberPercent formats a number as a percentage.
func NumberPercent(value float64, precision ...int) string {
	return NumberFormat(value, precision...) + "%"
}

// NumberCurrency formats a number as currency using a small deterministic
// currency-symbol map.
func NumberCurrency(value float64, currency ...string) string {
	code := "USD"
	if len(currency) > 0 && currency[0] != "" {
		code = strings.ToUpper(currency[0])
	}

	symbol := map[string]string{
		"EUR": "€",
		"GBP": "£",
		"JPY": "¥",
		"USD": "$",
	}[code]
	if symbol == "" {
		symbol = code + " "
	}

	return symbol + NumberFormat(value, 2)
}

// NumberClamp clamps value to the inclusive min/max bounds.
func NumberClamp[T ~int | ~int64 | ~float64](value, min, max T) T {
	if value < min {
		return min
	}
	if value > max {
		return max
	}

	return value
}

// NumberBytesToHuman formats bytes as a binary human-readable size.
func NumberBytesToHuman(bytes int64, precision ...int) string {
	digits := 2
	if len(precision) > 0 {
		digits = precision[0]
	}

	units := []string{"B", "KB", "MB", "GB", "TB", "PB"}
	value := float64(bytes)
	unit := 0

	for math.Abs(value) >= 1024 && unit < len(units)-1 {
		value /= 1024
		unit++
	}

	if unit == 0 {
		return fmt.Sprintf("%d %s", bytes, units[unit])
	}

	return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%.*f", digits, value), "0"), ".") + " " + units[unit]
}

// NumberToHuman formats large numbers using short English suffixes.
func NumberToHuman(value float64, precision ...int) string {
	digits := 1
	if len(precision) > 0 {
		digits = precision[0]
	}

	units := []struct {
		value  float64
		suffix string
	}{
		{1_000_000_000_000, "T"},
		{1_000_000_000, "B"},
		{1_000_000, "M"},
		{1_000, "K"},
	}

	for _, unit := range units {
		if math.Abs(value) >= unit.value {
			return NumberTrim(fmt.Sprintf("%.*f", digits, value/unit.value)) + unit.suffix
		}
	}

	return NumberTrim(fmt.Sprintf("%.*f", digits, value))
}

// NumberSummarize is an alias for NumberToHuman.
func NumberSummarize(value float64, precision ...int) string {
	return NumberToHuman(value, precision...)
}

// NumberPairs splits a value into digit pairs from right to left.
func NumberPairs(value int64) []string {
	raw := strconv.FormatInt(value, 10)
	if strings.HasPrefix(raw, "-") {
		raw = raw[1:]
	}

	var pairs []string
	for len(raw) > 2 {
		pairs = append([]string{raw[len(raw)-2:]}, pairs...)
		raw = raw[:len(raw)-2]
	}
	if raw != "" {
		pairs = append([]string{raw}, pairs...)
	}

	return pairs
}

// NumberTrim removes insignificant trailing decimal zeroes.
func NumberTrim(value string) string {
	if !strings.Contains(value, ".") {
		return value
	}

	return strings.TrimRight(strings.TrimRight(value, "0"), ".")
}

// NumberParse strips common separators and parses a decimal number.
func NumberParse(value string) (float64, error) {
	normalized := strings.ReplaceAll(value, ",", "")
	normalized = strings.TrimSpace(normalized)

	return strconv.ParseFloat(normalized, 64)
}

// NumberParseInt strips common separators and parses an integer.
func NumberParseInt(value string) (int64, error) {
	parsed, err := NumberParse(value)
	if err != nil {
		return 0, err
	}

	return int64(parsed), nil
}

// NumberParseFloat strips common separators and parses a float.
func NumberParseFloat(value string) (float64, error) {
	return NumberParse(value)
}

// NumberOrdinal formats an integer as an ordinal.
func NumberOrdinal(value int) string {
	abs := value
	if abs < 0 {
		abs = -abs
	}

	suffix := "th"
	if abs%100 < 11 || abs%100 > 13 {
		switch abs % 10 {
		case 1:
			suffix = "st"
		case 2:
			suffix = "nd"
		case 3:
			suffix = "rd"
		}
	}

	return strconv.Itoa(value) + suffix
}

// NumberSpellout spells small non-negative integers in English.
func NumberSpellout(value int) string {
	if value == 0 {
		return "zero"
	}
	if value < 0 {
		return "minus " + NumberSpellout(-value)
	}

	ones := []string{"", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen"}
	tens := []string{"", "", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"}

	if value < 20 {
		return ones[value]
	}
	if value < 100 {
		if value%10 == 0 {
			return tens[value/10]
		}

		return tens[value/10] + "-" + ones[value%10]
	}
	if value < 1000 {
		if value%100 == 0 {
			return ones[value/100] + " hundred"
		}

		return ones[value/100] + " hundred " + NumberSpellout(value%100)
	}

	if value < 1_000_000 {
		if value%1000 == 0 {
			return NumberSpellout(value/1000) + " thousand"
		}

		return NumberSpellout(value/1000) + " thousand " + NumberSpellout(value%1000)
	}

	return strconv.Itoa(value)
}

// NumberSpellOrdinal spells small non-negative integer ordinals in English.
func NumberSpellOrdinal(value int) string {
	ordinals := map[int]string{
		0: "zeroth", 1: "first", 2: "second", 3: "third", 4: "fourth", 5: "fifth",
		6: "sixth", 7: "seventh", 8: "eighth", 9: "ninth", 10: "tenth", 11: "eleventh",
		12: "twelfth", 13: "thirteenth", 14: "fourteenth", 15: "fifteenth", 16: "sixteenth",
		17: "seventeenth", 18: "eighteenth", 19: "nineteenth", 20: "twentieth",
	}
	if value, ok := ordinals[value]; ok {
		return value
	}

	return NumberSpellout(value) + "th"
}

func groupThousands(value string) string {
	if len(value) <= 3 {
		return value
	}

	var out []byte
	remainder := len(value) % 3
	if remainder == 0 {
		remainder = 3
	}

	out = append(out, value[:remainder]...)
	for i := remainder; i < len(value); i += 3 {
		out = append(out, ',')
		out = append(out, value[i:i+3]...)
	}

	return string(out)
}
