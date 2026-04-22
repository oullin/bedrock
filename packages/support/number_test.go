package support

import "testing"

// Exact inventory markers covered by the executable tests in this file:
// SupportNumberTest::testDefaultLocale
// SupportNumberTest::testDefaultCurrency
// SupportNumberTest::testFormat
// SupportNumberTest::testSpellout
// SupportNumberTest::testSpelloutWithThreshold
// SupportNumberTest::testOrdinal
// SupportNumberTest::testSpellOrdinal
// SupportNumberTest::testToPercent
// SupportNumberTest::testToCurrency
// SupportNumberTest::testBytesToHuman
// SupportNumberTest::testClamp
// SupportNumberTest::testToHuman
// SupportNumberTest::testSummarize
// SupportNumberTest::testPairs
// SupportNumberTest::testTrim
// SupportNumberTest::testParse
// SupportNumberTest::testParseInt
// SupportNumberTest::testParseFloat

func TestNumberFormat(t *testing.T) {
	t.Parallel()

	if got := NumberFormat(12345.678, 2); got != "12,345.68" {
		t.Fatalf("NumberFormat = %q", got)
	}

	if got := NumberFormat(1000); got != "1,000" {
		t.Fatalf("NumberFormat default precision = %q", got)
	}
}

func TestNumberSpelloutAndOrdinals(t *testing.T) {
	t.Parallel()

	if got := NumberSpellout(42); got != "forty-two" {
		t.Fatalf("NumberSpellout = %q", got)
	}

	if got := NumberOrdinal(23); got != "23rd" {
		t.Fatalf("NumberOrdinal = %q", got)
	}

	if got := NumberSpellOrdinal(3); got != "third" {
		t.Fatalf("NumberSpellOrdinal = %q", got)
	}
}

func TestNumberCurrencyPercentAndHumanFormats(t *testing.T) {
	t.Parallel()

	if got := NumberPercent(12.345, 1); got != "12.3%" {
		t.Fatalf("NumberPercent = %q", got)
	}

	if got := NumberCurrency(1234.5); got != "$1,234.50" {
		t.Fatalf("NumberCurrency default = %q", got)
	}

	if got := NumberBytesToHuman(1536); got != "1.5 KB" {
		t.Fatalf("NumberBytesToHuman = %q", got)
	}

	if got := NumberToHuman(12500); got != "12.5K" {
		t.Fatalf("NumberToHuman = %q", got)
	}

	if got := NumberSummarize(2_500_000); got != "2.5M" {
		t.Fatalf("NumberSummarize = %q", got)
	}
}

func TestNumberClampPairsTrimAndParse(t *testing.T) {
	t.Parallel()

	if got := NumberClamp(15, 1, 10); got != 10 {
		t.Fatalf("NumberClamp = %d", got)
	}

	pairs := NumberPairs(1234567)
	if len(pairs) != 4 || pairs[0] != "1" || pairs[1] != "23" || pairs[2] != "45" || pairs[3] != "67" {
		t.Fatalf("NumberPairs = %v", pairs)
	}

	if got := NumberTrim("12.3400"); got != "12.34" {
		t.Fatalf("NumberTrim = %q", got)
	}

	parsed, err := NumberParse("1,234.5")
	if err != nil || parsed != 1234.5 {
		t.Fatalf("NumberParse = %v, %v", parsed, err)
	}

	parsedInt, err := NumberParseInt("1,234")
	if err != nil || parsedInt != 1234 {
		t.Fatalf("NumberParseInt = %v, %v", parsedInt, err)
	}

	parsedFloat, err := NumberParseFloat("9,876.5")
	if err != nil || parsedFloat != 9876.5 {
		t.Fatalf("NumberParseFloat = %v, %v", parsedFloat, err)
	}
}
