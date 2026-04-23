package validation_test

import (
	"testing"
	"time"

	"github.com/bedrock/packages/validation"
)

type presenceCountCall struct {
	table     string
	column    string
	value     string
	excludeID *string
	idColumn  *string
	extras    map[string]any
}

type presenceVerifierFake struct {
	count      int
	multiCount int
	lastCount  presenceCountCall
}

func (f *presenceVerifierFake) GetCount(table, column, value string, excludeID *string, idColumn *string, extras map[string]any) int {
	f.lastCount = presenceCountCall{
		table:     table,
		column:    column,
		value:     value,
		excludeID: excludeID,
		idColumn:  idColumn,
		extras:    extras,
	}

	return f.count
}

func (f *presenceVerifierFake) GetMultiCount(table, column string, values []string, extras map[string]any) int {
	return f.multiCount
}

func makeValidatorWithPresenceVerifier(pv validation.PresenceVerifier, data, rules map[string]any) *validation.Validator {
	f := validation.NewFactory()
	f.SetPresenceVerifier(pv)

	return f.Make(data, rules, nil, nil)
}

func utcDay(offsetDays int) string {
	return time.Now().UTC().AddDate(0, 0, offsetDays).Format("2006-01-02")
}

// ValidationNumericRuleTest::testDefaultNumericRule
// ValidationValidatorTest::testValidateNumeric
func TestNumericRule_DefaultNumeric(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{"amount": 42}, map[string]any{"amount": "numeric"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"amount": "not-a-number"}, map[string]any{"amount": "numeric"})
	assertFails(t, v2)
}

// ValidationNumericRuleTest::testDigitsRule
// ValidationValidatorTest::testValidateDigits
func TestNumericRule_Digits(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{"code": "12345"}, map[string]any{"code": "digits:5"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"code": "1234a"}, map[string]any{"code": "digits:5"})
	assertFails(t, v2)
	assertError(t, v2, "code", "digits")
}

// ValidationNumericRuleTest::testDigitsBetweenRule
func TestNumericRule_DigitsBetween(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{"code": "12345"}, map[string]any{"code": "digits_between:4,6"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"code": "123"}, map[string]any{"code": "digits_between:4,6"})
	assertFails(t, v2)
}

// ValidationNumericRuleTest::testDecimalRule
// ValidationValidatorTest::testValidateDecimal
func TestNumericRule_Decimal(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{"price": "12.34"}, map[string]any{"price": "decimal:2"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"price": "12.3"}, map[string]any{"price": "decimal:2"})
	assertFails(t, v2)
	assertError(t, v2, "price", "decimal")
}

// ValidationNumericRuleTest::testMinDigitsRule
func TestNumericRule_MinDigits(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{"code": "12345"}, map[string]any{"code": "min_digits:4"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"code": "123"}, map[string]any{"code": "min_digits:4"})
	assertFails(t, v2)
}

// ValidationNumericRuleTest::testMaxDigitsRule
func TestNumericRule_MaxDigits(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{"code": "1234"}, map[string]any{"code": "max_digits:4"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"code": "12345"}, map[string]any{"code": "max_digits:4"})
	assertFails(t, v2)
}

// ValidationNumericRuleTest::testMultipleOfRule
// ValidationValidatorTest::testValidateMultipleOf
func TestNumericRule_MultipleOf(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{"amount": "1.5"}, map[string]any{"amount": "multiple_of:0.5"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"amount": "1.2"}, map[string]any{"amount": "multiple_of:0.5"})
	assertFails(t, v2)
	assertError(t, v2, "amount", "multiple of")
}

// ValidationRuleContainsTest::testContainsValidation
// ValidationValidatorTest::testValidateContains
func TestArrayRule_Contains(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{"tags": []any{"alpha", "beta"}}, map[string]any{"tags": "contains:alpha,beta"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"tags": []any{"alpha"}}, map[string]any{"tags": "contains:alpha,beta"})
	assertFails(t, v2)
	assertError(t, v2, "tags", "contain")
}

// ValidationRuleDoesntContainTest::testDoesntContainValidation
func TestArrayRule_DoesntContain(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{"tags": []any{"alpha", "beta"}}, map[string]any{"tags": "doesnt_contain:secret,classified"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"tags": []any{"alpha", "secret"}}, map[string]any{"tags": "doesnt_contain:secret,classified"})
	assertFails(t, v2)
}

// ValidationValidatorTest::testLowercase
func TestValidator_Lowercase(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{"handle": "lowercase"}, map[string]any{"handle": "lowercase"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"handle": "CamelCase"}, map[string]any{"handle": "lowercase"})
	assertFails(t, v2)
	assertError(t, v2, "handle", "lowercase")
}

// ValidationValidatorTest::testUppercase
func TestValidator_Uppercase(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{"handle": "UPPERCASE"}, map[string]any{"handle": "uppercase"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"handle": "MixedCase"}, map[string]any{"handle": "uppercase"})
	assertFails(t, v2)
	assertError(t, v2, "handle", "uppercase")
}

// ValidationValidatorTest::testDateEquals
// ValidationValidatorTest::testDateEqualsRespectsCarbonTestNowWhenParameterIsRelative
func TestValidator_DateEquals(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{"published_at": utcDay(0)}, map[string]any{"published_at": "date_equals:today"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"published_at": utcDay(1)}, map[string]any{"published_at": "date_equals:today"})
	assertFails(t, v2)

	v3 := makeValidator(map[string]any{"published_at": utcDay(1)}, map[string]any{"published_at": "date_equals:tomorrow"})
	assertPasses(t, v3)
}

// ValidationDatabasePresenceVerifierTest::testBasicCount
func TestDatabasePresenceVerifier_BasicCount(t *testing.T) {
	t.Parallel()

	fake := &presenceVerifierFake{count: 1}
	v := makeValidatorWithPresenceVerifier(
		fake,
		map[string]any{"email": "user@example.com"},
		map[string]any{"email": "exists:users,email"},
	)

	assertPasses(t, v)

	if fake.lastCount.table != "users" || fake.lastCount.column != "email" || fake.lastCount.value != "user@example.com" {
		t.Fatalf("GetCount() captured %+v, want users/email/user@example.com", fake.lastCount)
	}

	if fake.lastCount.excludeID != nil || fake.lastCount.idColumn != nil || fake.lastCount.extras != nil {
		t.Fatalf("GetCount() captured unexpected optional args: %+v", fake.lastCount)
	}

	fake = &presenceVerifierFake{count: 0}
	v2 := makeValidatorWithPresenceVerifier(
		fake,
		map[string]any{"email": "user@example.com"},
		map[string]any{"email": "exists:users,email"},
	)

	assertFails(t, v2)
}

// ValidationDatabasePresenceVerifierTest::testBasicCountWithClosures
func TestDatabasePresenceVerifier_BasicCountWithClosures(t *testing.T) {
	t.Parallel()

	fake := &presenceVerifierFake{count: 1}
	v := makeValidatorWithPresenceVerifier(
		fake,
		map[string]any{"email": "user@example.com", "status": "active"},
		map[string]any{"email": "exists:users,email,status,active"},
	)

	assertPasses(t, v)

	if fake.lastCount.extras["status"] != "active" {
		t.Fatalf("GetCount() extras = %v, want status=active", fake.lastCount.extras)
	}
}

// ValidationDatabasePresenceVerifierTest::testGetCountWithValidExcludeId
func TestDatabasePresenceVerifier_GetCountWithValidExcludeID(t *testing.T) {
	t.Parallel()

	excludeID := "123"
	fake := &presenceVerifierFake{count: 0}
	v := makeValidatorWithPresenceVerifier(
		fake,
		map[string]any{"email": "user@example.com"},
		map[string]any{"email": "unique:users,email,123,id"},
	)

	assertPasses(t, v)

	if fake.lastCount.excludeID == nil || *fake.lastCount.excludeID != excludeID {
		t.Fatalf("GetCount() excludeID = %v, want %q", fake.lastCount.excludeID, excludeID)
	}

	if fake.lastCount.idColumn == nil || *fake.lastCount.idColumn != "id" {
		t.Fatalf("GetCount() idColumn = %v, want id", fake.lastCount.idColumn)
	}
}

// ValidationExcludeUnlessRuleTest::testFieldIsExcludedWhenConditionFalse
// ValidationExcludeUnlessRuleTest::testFieldIsKeptWhenConditionTrue
func TestValidator_ExcludeUnless(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"type": "user", "secret": "classified"},
		map[string]any{"secret": "exclude_unless:type,admin|string"},
	)

	got, err := v.Validated()

	if err != nil {
		t.Fatalf("Validated() returned error: %v", err)
	}

	if _, ok := got["secret"]; ok {
		t.Fatalf("Validated() = %v, want secret excluded", got)
	}

	v2 := makeValidator(
		map[string]any{"type": "admin", "secret": "classified"},
		map[string]any{"secret": "exclude_unless:type,admin|string"},
	)

	got2, err := v2.Validated()

	if err != nil {
		t.Fatalf("Validated() returned error: %v", err)
	}

	if got2["secret"] != "classified" {
		t.Fatalf("Validated() = %v, want secret retained", got2)
	}
}
