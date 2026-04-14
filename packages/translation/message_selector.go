package translation

import (
	"math"
	"strconv"
	"strings"
)

// MessageSelector selects the appropriate plural form from a translation
// string, mirroring Laravel's Illuminate\Translation\MessageSelector.
//
// Translation strings use pipe-delimited segments with optional bracket
// conditions:
//
//	"one item|many items"
//	"{0} no items|{1} one item|[2,*] :count items"
//	"[*,-1] negative|{0} zero|[1,*] positive"
type MessageSelector struct{}

// NewMessageSelector returns a new *MessageSelector.
func NewMessageSelector() *MessageSelector {
	return &MessageSelector{}
}

// Choose selects the appropriate plural segment from line for the given
// number and locale.  Pipe characters (|) separate segments; each segment
// may carry an optional bracket condition prefix.
func (s *MessageSelector) Choose(line string, number float64, locale string) string {
	segments := splitSegments(line)

	// First pass: match explicit bracket conditions.
	for _, seg := range segments {
		text, matched := extractFromString(seg, number)

		if matched {
			return strings.TrimSpace(text)
		}
	}

	// Second pass: strip conditions and pick by plural index.
	stripped := stripConditions(segments)

	if len(stripped) == 0 {
		return ""
	}

	idx := getPluralIndex(locale, math.Abs(number))

	if idx >= len(stripped) {
		idx = len(stripped) - 1
	}

	return strings.TrimSpace(stripped[idx])
}

// splitSegments splits a translation line on "|" preserving empty segments.
func splitSegments(line string) []string {
	return strings.Split(line, "|")
}

// extractFromString checks whether segment has an explicit bracket condition
// ({N} or [M,N]) that matches number.  Returns the content and true on match.
func extractFromString(segment string, number float64) (string, bool) {
	segment = strings.TrimSpace(segment)

	if len(segment) == 0 {
		return "", false
	}

	switch segment[0] {
	case '{':
		end := strings.Index(segment, "}")

		if end < 0 {
			return "", false
		}

		raw := segment[1:end]
		n, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)

		if err != nil {
			return "", false
		}

		if number == n {
			return segment[end+1:], true
		}

		return "", false

	case '[', '(':
		end := strings.IndexAny(segment, "])")

		if end < 0 {
			return "", false
		}

		leftInclusive := segment[0] == '['
		rightInclusive := segment[end] == ']'
		inner := segment[1:end]
		text := segment[end+1:]

		parts := strings.SplitN(inner, ",", 2)

		if len(parts) != 2 {
			return "", false
		}

		lo, lInf := parseRangeBound(strings.TrimSpace(parts[0]))
		hi, hInf := parseRangeBound(strings.TrimSpace(parts[1]))

		var loMatch, hiMatch bool

		if lInf {
			loMatch = true
		} else if leftInclusive {
			loMatch = number >= lo
		} else {
			loMatch = number > lo
		}

		if hInf {
			hiMatch = true
		} else if rightInclusive {
			hiMatch = number <= hi
		} else {
			hiMatch = number < hi
		}

		if loMatch && hiMatch {
			return text, true
		}

		return "", false
	}

	return "", false
}

// parseRangeBound parses a range bound value.  "*" means infinity (returns
// 0, true).
func parseRangeBound(s string) (float64, bool) {
	if s == "*" {
		return 0, true
	}

	n, err := strconv.ParseFloat(s, 64)

	if err != nil {
		return 0, false
	}

	return n, false
}

// stripConditions removes bracket/curly-brace condition prefixes from each
// segment, returning only the text portions.
func stripConditions(segments []string) []string {
	out := make([]string, 0, len(segments))

	for _, seg := range segments {
		seg = strings.TrimSpace(seg)

		if len(seg) == 0 {
			out = append(out, seg)

			continue
		}

		switch seg[0] {
		case '{':
			end := strings.Index(seg, "}")

			if end >= 0 {
				out = append(out, seg[end+1:])
			} else {
				out = append(out, seg)
			}
		case '[', '(':
			end := strings.IndexAny(seg, "])")

			if end >= 0 {
				out = append(out, seg[end+1:])
			} else {
				out = append(out, seg)
			}
		default:
			out = append(out, seg)
		}
	}

	return out
}

// getPluralIndex returns the CLDR plural-form index for locale and n.
// n is the absolute value of the original number as a float so that 1.2
// does not satisfy n==1 (mirroring PHP's loose comparison behaviour).
// Rules ported from Zend Framework / Laravel's MessageSelector::getPluralIndex.
//
//nolint:cyclop,gocyclo
func getPluralIndex(locale string, n float64) int {
	// Integer value used for modulo-based rules (Russian, Polish, etc.).
	ni := int(n)
	// Normalise locale: strip region tag for matching (e.g. "zh_CN" → "zh").
	base := locale

	if idx := strings.IndexAny(locale, "-_"); idx >= 0 {
		base = locale[:idx]
	}

	switch base {
	// Always singular (index 0).
	case "az", "bo", "cnh", "dz", "fa", "id", "ja", "ka", "km", "kn", "ko",
		"lo", "ms", "my", "sah", "th", "tk", "tr", "ug", "vi", "wo", "zh":
		return 0

	// Binary: 0 = singular (n==1), 1 = plural.
	case "af", "bn", "bg", "ca", "da", "de", "el", "en", "eo", "es", "et",
		"eu", "fi", "fo", "fy", "gl", "gu", "he", "hu", "is", "it",
		"ku", "lb", "ml", "mn", "mr", "nb", "ne", "nl", "nn", "no", "or",
		"pa", "pt", "sd", "si", "sq", "sv", "sw", "ta", "te", "ur", "yo",
		"zu":
		if n == 1 {
			return 0
		}

		return 1

	// French / Portuguese (Brazil): 0 if n == 0 or 1, else 1.
	case "am", "bh", "fil", "fr", "gun", "hi", "hy", "ln", "mg", "nso",
		"xbr", "ti", "wa":
		if n == 0 || n == 1 {
			return 0
		}

		return 1

	// Latvian.
	case "lv":
		if n == 0 {
			return 0
		}

		if ni%10 == 1 && ni%100 != 11 {
			return 1
		}

		return 2

	// Cornish / Gaelic (Scottish).
	case "ga":
		if n == 1 {
			return 0
		}

		if n == 2 {
			return 1
		}

		return 2

	// Romanian.
	case "ro":
		if n == 1 {
			return 0
		}

		if n == 0 || (ni%100 > 0 && ni%100 < 20) {
			return 1
		}

		return 2

	// Lithuanian.
	case "lt":
		if ni%10 == 1 && ni%100 != 11 {
			return 0
		}

		if ni%10 >= 2 && (ni%100 < 10 || ni%100 >= 20) {
			return 1
		}

		return 2

	// Russian, Ukrainian, Belarusian, Serbian.
	case "ru", "uk", "be", "sr":
		if ni%10 == 1 && ni%100 != 11 {
			return 0
		}

		if ni%10 >= 2 && ni%10 <= 4 && (ni%100 < 10 || ni%100 >= 20) {
			return 1
		}

		return 2

	// Czech, Slovak.
	case "cs", "sk":
		if n == 1 {
			return 0
		}

		if n >= 2 && n <= 4 {
			return 1
		}

		return 2

	// Polish.
	case "pl":
		if n == 1 {
			return 0
		}

		if ni%10 >= 2 && ni%10 <= 4 && (ni%100 < 12 || ni%100 > 14) {
			return 1
		}

		return 2

	// Slovenian.
	case "sl":
		if ni%100 == 1 {
			return 0
		}

		if ni%100 == 2 {
			return 1
		}

		if ni%100 == 3 || ni%100 == 4 {
			return 2
		}

		return 3

	// Macedonian.
	case "mk":
		if ni%10 == 1 {
			return 0
		}

		if ni%10 == 2 {
			return 1
		}

		return 2

	// Maltese.
	case "mt":
		if n == 1 {
			return 0
		}

		if n == 0 || (ni%100 > 1 && ni%100 < 11) {
			return 1
		}

		if ni%100 > 10 && ni%100 < 20 {
			return 2
		}

		return 3

	// Welsh.
	case "cy":
		if n == 1 {
			return 0
		}

		if n == 2 {
			return 1
		}

		if n == 8 || n == 11 {
			return 2
		}

		return 3

	// Irish.
	case "gd":
		if n == 1 || n == 11 {
			return 0
		}

		if n == 2 || n == 12 {
			return 1
		}

		if (n >= 3 && n <= 10) || (n >= 13 && n <= 19) {
			return 2
		}

		return 3

	// Breton.
	case "br":
		if ni%10 == 1 && ni%100 != 11 && ni%100 != 71 && ni%100 != 91 {
			return 0
		}

		if ni%10 == 2 && ni%100 != 12 && ni%100 != 72 && ni%100 != 92 {
			return 1
		}

		if (ni%10 == 3 || ni%10 == 4 || ni%10 == 9) &&
			(ni%100 < 10 || ni%100 > 19) && (ni%100 < 70 || ni%100 > 79) &&
			(ni%100 < 90 || ni%100 > 99) {
			return 2
		}

		if n != 0 && ni%1000000 == 0 {
			return 3
		}

		return 4

	// Arabic.
	case "ar":
		if n == 0 {
			return 0
		}

		if n == 1 {
			return 1
		}

		if n == 2 {
			return 2
		}

		if ni%100 >= 3 && ni%100 <= 10 {
			return 3
		}

		if ni%100 >= 11 && ni%100 <= 99 {
			return 4
		}

		return 5

	// Upper Sorbian.
	case "hsb":
		if ni%100 == 1 {
			return 0
		}

		if ni%100 == 2 {
			return 1
		}

		if ni%100 == 3 || ni%100 == 4 {
			return 2
		}

		return 3

	// Tagalog / Filipino (extended).
	case "tl":
		if n == 1 || n == 2 || n == 3 {
			return 0
		}

		if ni%10 == 4 || ni%10 == 6 || ni%10 == 9 {
			return 0
		}

		return 1

	default:
		return 0
	}
}
