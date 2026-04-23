package validation_test

import (
	"reflect"
	"testing"

	"github.com/bedrock/packages/validation"
)

// ValidationValidatorTest::testValidateUrlWithProtocols
func TestValidateUrlWithProtocols(t *testing.T) {
	t.Parallel()

	for _, url := range []string{"http://example.com", "https://example.com", "ftp://example.com"} {
		v := makeValidator(map[string]any{"site": url}, map[string]any{"site": "url"})

		if !v.Passes() {
			t.Fatalf("expected %q to pass url validation, got %v", url, v.Errors().All())
		}
	}
}

// ValidationValidatorTest::testValidateUrlWithValidUrls
func TestValidateUrlWithValidUrls(t *testing.T) {
	t.Parallel()

	for _, url := range []string{
		"https://example.com",
		"https://example.com/path",
		"https://sub.example.co.uk?query=1",
	} {
		v := makeValidator(map[string]any{"site": url}, map[string]any{"site": "url"})

		if !v.Passes() {
			t.Fatalf("expected %q to pass url validation, got %v", url, v.Errors().All())
		}
	}
}

// ValidationValidatorTest::testValidateUrlWithInvalidUrls
func TestValidateUrlWithInvalidUrls(t *testing.T) {
	t.Parallel()

	for _, url := range []string{"example.com", "http://", "://broken"} {
		v := makeValidator(map[string]any{"site": url}, map[string]any{"site": "url"})

		if !v.Fails() {
			t.Fatalf("expected %q to fail url validation", url)
		}
	}
}

// ValidationValidatorTest::testValidateWithValidUuid
func TestValidateWithValidUuid(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"id": "550e8400-e29b-41d4-a716-446655440000"},
		map[string]any{"id": "uuid"},
	)

	if !v.Passes() {
		t.Fatalf("expected uuid validation to pass, got %v", v.Errors().All())
	}
}

// ValidationValidatorTest::testValidateWithInvalidUuid
func TestValidateWithInvalidUuid(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"id": "not-a-uuid"},
		map[string]any{"id": "uuid"},
	)

	if !v.Fails() {
		t.Fatal("expected invalid uuid to fail validation")
	}
}

// ValidationValidatorTest::testValidateNotRegex
func TestValidateNotRegex(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"code": "abc123"},
		map[string]any{"code": "not_regex:/^[0-9]+$/"},
	)

	if !v.Passes() {
		t.Fatalf("expected not_regex to pass, got %v", v.Errors().All())
	}

	v2 := makeValidator(
		map[string]any{"code": "12345"},
		map[string]any{"code": "not_regex:/^[0-9]+$/"},
	)

	if !v2.Fails() {
		t.Fatal("expected matching value to fail not_regex validation")
	}
}

// ValidationValidatorTest::testValidateStartsWithDoesNotThrowOnNonStringValue
func TestValidateStartsWithDoesNotThrowOnNonStringValue(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": 12345},
		map[string]any{"name": "starts_with:12"},
	)

	if !v.Fails() {
		t.Fatal("expected non-string value to fail starts_with validation")
	}
}

// ValidationValidatorTest::testValidateEndsWithDoesNotThrowOnNonStringValue
func TestValidateEndsWithDoesNotThrowOnNonStringValue(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": 12345},
		map[string]any{"name": "ends_with:45"},
	)

	if !v.Fails() {
		t.Fatal("expected non-string value to fail ends_with validation")
	}
}

// ValidationValidatorTest::testValidateHexColorDoesNotThrowOnNonStringValue
func TestValidateHexColorDoesNotThrowOnNonStringValue(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"color": 42},
		map[string]any{"color": "hex_color"},
	)

	if !v.Fails() {
		t.Fatal("expected non-string value to fail hex_color validation")
	}
}

// ValidationValidatorTest::testValidateMacAddress
func TestValidateMacAddress(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"mac": "aa:bb:cc:dd:ee:ff"},
		map[string]any{"mac": "mac_address"},
	)

	if !v.Passes() {
		t.Fatalf("expected mac address validation to pass, got %v", v.Errors().All())
	}

	v2 := makeValidator(
		map[string]any{"mac": "not-a-mac"},
		map[string]any{"mac": "mac_address"},
	)

	if !v2.Fails() {
		t.Fatal("expected invalid mac address to fail validation")
	}
}

// ValidationStringRuleTest::testDefaultStringRule
func TestStringRule_DefaultStringRule(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": "Taylor"},
		map[string]any{"name": "string"},
	)

	if !v.Passes() {
		t.Fatalf("expected string validation to pass, got %v", v.Errors().All())
	}
}

// ValidationStringRuleTest::testStartsWithRule
func TestStringRule_StartsWithRule(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": "Taylor Swift"},
		map[string]any{"name": "starts_with:Taylor"},
	)

	if !v.Passes() {
		t.Fatalf("expected starts_with validation to pass, got %v", v.Errors().All())
	}
}

// ValidationStringRuleTest::testEndsWithRule
func TestStringRule_EndsWithRule(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": "Taylor Swift"},
		map[string]any{"name": "ends_with:Swift"},
	)

	if !v.Passes() {
		t.Fatalf("expected ends_with validation to pass, got %v", v.Errors().All())
	}
}

// ValidationStringRuleTest::testDoesntStartWithRule
func TestStringRule_DoesntStartWithRule(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": "Taylor Swift"},
		map[string]any{"name": "doesnt_start_with:Dr.,Mr."},
	)

	if !v.Passes() {
		t.Fatalf("expected doesnt_start_with validation to pass, got %v", v.Errors().All())
	}
}

// ValidationStringRuleTest::testDoesntEndWithRule
func TestStringRule_DoesntEndWithRule(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": "Taylor Swift"},
		map[string]any{"name": "doesnt_end_with:.jpg,.png"},
	)

	if !v.Passes() {
		t.Fatalf("expected doesnt_end_with validation to pass, got %v", v.Errors().All())
	}
}

// ValidationStringRuleTest::testMinRule
func TestStringRule_MinRule(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": "Taylor"},
		map[string]any{"name": "min:3"},
	)

	if !v.Passes() {
		t.Fatalf("expected min validation to pass, got %v", v.Errors().All())
	}
}

// ValidationStringRuleTest::testMaxRule
func TestStringRule_MaxRule(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": "Taylor"},
		map[string]any{"name": "max:10"},
	)

	if !v.Passes() {
		t.Fatalf("expected max validation to pass, got %v", v.Errors().All())
	}
}

// ValidationStringRuleTest::testBetweenRule
func TestStringRule_BetweenRule(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": "Taylor"},
		map[string]any{"name": "between:3,10"},
	)

	if !v.Passes() {
		t.Fatalf("expected between validation to pass, got %v", v.Errors().All())
	}
}

// ValidationStringRuleTest::testExactlyRule
func TestStringRule_ExactlyRule(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": "Taylor"},
		map[string]any{"name": "size:6"},
	)

	if !v.Passes() {
		t.Fatalf("expected size validation to pass, got %v", v.Errors().All())
	}
}

// ValidationStringRuleTest::testAlphaRule
func TestStringRule_AlphaRule(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": "Taylor"},
		map[string]any{"name": "alpha"},
	)

	if !v.Passes() {
		t.Fatalf("expected alpha validation to pass, got %v", v.Errors().All())
	}

	v2 := makeValidator(
		map[string]any{"name": "Taylor123"},
		map[string]any{"name": "alpha"},
	)

	if !v2.Fails() {
		t.Fatal("expected alpha validation to fail on digits")
	}
}

// ValidationStringRuleTest::testAlphaNumericRule
func TestStringRule_AlphaNumericRule(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": "Taylor123"},
		map[string]any{"name": "alpha_num"},
	)

	if !v.Passes() {
		t.Fatalf("expected alpha_num validation to pass, got %v", v.Errors().All())
	}

	v2 := makeValidator(
		map[string]any{"name": "Taylor-123"},
		map[string]any{"name": "alpha_num"},
	)

	if !v2.Fails() {
		t.Fatal("expected alpha_num validation to fail on punctuation")
	}
}

// ValidationStringRuleTest::testAlphaDashRule
func TestStringRule_AlphaDashRule(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": "Taylor-123_abc"},
		map[string]any{"name": "alpha_dash"},
	)

	if !v.Passes() {
		t.Fatalf("expected alpha_dash validation to pass, got %v", v.Errors().All())
	}

	v2 := makeValidator(
		map[string]any{"name": "Taylor 123"},
		map[string]any{"name": "alpha_dash"},
	)

	if !v2.Fails() {
		t.Fatal("expected alpha_dash validation to fail on spaces")
	}
}

// ValidationStringRuleTest::testAsciiRule
func TestStringRule_AsciiRule(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": "Plain ASCII"},
		map[string]any{"name": "ascii"},
	)

	if !v.Passes() {
		t.Fatalf("expected ascii validation to pass, got %v", v.Errors().All())
	}

	v2 := makeValidator(
		map[string]any{"name": "mañana"},
		map[string]any{"name": "ascii"},
	)

	if !v2.Fails() {
		t.Fatal("expected ascii validation to fail on non-ascii text")
	}
}

// ValidationStringRuleTest::testUppercaseRule
func TestStringRule_UppercaseRule(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": "TAYLOR"},
		map[string]any{"name": "uppercase"},
	)

	if !v.Passes() {
		t.Fatalf("expected uppercase validation to pass, got %v", v.Errors().All())
	}
}

// ValidationStringRuleTest::testLowercaseRule
func TestStringRule_LowercaseRule(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": "taylor"},
		map[string]any{"name": "lowercase"},
	)

	if !v.Passes() {
		t.Fatalf("expected lowercase validation to pass, got %v", v.Errors().All())
	}
}

// ValidationRuleParserTest::testExplodeProperlyParsesRegexThatDoesNotContainPipe
func TestRuleParser_ExplodeProperlyParsesRegexThatDoesNotContainPipeExact(t *testing.T) {
	t.Parallel()

	got := validation.Explode("regex:/^[a-z]+$/")

	if len(got) != 1 {
		t.Fatalf("Explode(regex) = %v, want one rule", got)
	}

	if got[0].Name != "Regex" {
		t.Fatalf("Explode(regex)[0].Name = %q, want Regex", got[0].Name)
	}

	if !reflect.DeepEqual(got[0].Parameters, []string{"/^[a-z]+$/"}) {
		t.Fatalf("Explode(regex)[0].Parameters = %v, want a preserved regex parameter", got[0].Parameters)
	}
}

// ValidationRuleParserTest::testExplodeHandlesDateRuleWithAdditionalRules
func TestRuleParser_ExplodeHandlesDateRuleWithAdditionalRules(t *testing.T) {
	t.Parallel()

	got := validation.Explode("date|before:2026-04-23")

	if len(got) != 2 {
		t.Fatalf("Explode(date|before) = %v, want two rules", got)
	}

	if got[0].Name != "Date" || got[1].Name != "Before" {
		t.Fatalf("Explode(date|before) = %v, want Date then Before", got)
	}
}

// ValidationRuleParserTest::testExplodeHandlesNumericRuleWithAdditionalRules
func TestRuleParser_ExplodeHandlesNumericRuleWithAdditionalRules(t *testing.T) {
	t.Parallel()

	got := validation.Explode("numeric|min:10")

	if len(got) != 2 {
		t.Fatalf("Explode(numeric|min) = %v, want two rules", got)
	}

	if got[0].Name != "Numeric" || got[1].Name != "Min" {
		t.Fatalf("Explode(numeric|min) = %v, want Numeric then Min", got)
	}
}
