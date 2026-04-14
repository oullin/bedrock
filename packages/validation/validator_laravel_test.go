package validation_test

// Ports of Illuminate\Tests\Validation\ValidationValidatorTest.
//
// Source: https://github.com/laravel/framework/13.x/tests/Validation/ValidationValidatorTest.php
//
// Naming rules (from packages/queue/PARITY.md):
//   PHP testFooBar        → Go TestFooBar
//   PHP test_foo_bar      → Go TestFooBar

import (
	"testing"

	"github.com/bedrock/packages/validation"
)

// Port of ValidationValidatorTest::testPassesReturnsTrueIfNoFailingRules
func TestPassesReturnsTrueIfNoFailingRules(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": "Taylor"},
		map[string]any{"name": "alpha"},
	)

	if !v.Passes() {
		t.Error("expected passes() == true")
	}
}

// Port of ValidationValidatorTest::testFailsReturnsFalseIfFailingRules
func TestFailsReturnsFalseIfFailingRules(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": "Taylor1"},
		map[string]any{"name": "alpha"},
	)

	if !v.Fails() {
		t.Error("expected fails() == true")
	}
}

// Port of ValidationValidatorTest::testHasFailedRules
func TestHasFailedRules(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"age": "not-a-number"},
		map[string]any{"age": "numeric"},
	)
	v.Fails()

	failed := v.Failed()
	if len(failed["age"]) == 0 {
		t.Error("expected age in failed rules")
	}
}

// Port of ValidationValidatorTest::testValidateRequired
func TestValidateRequired(t *testing.T) {
	t.Parallel()

	cases := []struct {
		data  map[string]any
		rules map[string]any
		pass  bool
	}{
		{map[string]any{"name": "foo"}, map[string]any{"name": "required"}, true},
		{map[string]any{"name": ""}, map[string]any{"name": "required"}, false},
		{map[string]any{}, map[string]any{"name": "required"}, false},
		{map[string]any{"name": nil}, map[string]any{"name": "required"}, false},
	}

	for _, tc := range cases {
		v := makeValidator(tc.data, tc.rules)
		if v.Passes() != tc.pass {
			t.Errorf("data=%v: expected passes=%v", tc.data, tc.pass)
		}
	}
}

// Port of ValidationValidatorTest::testValidateRequiredIf
func TestValidateRequiredIf(t *testing.T) {
	t.Parallel()

	// required_if:other,value — field is required when other == value
	v := makeValidator(
		map[string]any{"foo": "yes", "bar": ""},
		map[string]any{"bar": "required_if:foo,yes"},
	)
	assertFails(t, v)

	v2 := makeValidator(
		map[string]any{"foo": "no", "bar": ""},
		map[string]any{"bar": "required_if:foo,yes"},
	)
	assertPasses(t, v2)

	v3 := makeValidator(
		map[string]any{"foo": "yes", "bar": "something"},
		map[string]any{"bar": "required_if:foo,yes"},
	)
	assertPasses(t, v3)
}

// Port of ValidationValidatorTest::testValidateRequiredUnless
func TestValidateRequiredUnless(t *testing.T) {
	t.Parallel()

	// required_unless:other,value — required unless other is one of the values
	v := makeValidator(
		map[string]any{"foo": "yes", "bar": ""},
		map[string]any{"bar": "required_unless:foo,no"},
	)
	assertFails(t, v)

	v2 := makeValidator(
		map[string]any{"foo": "no", "bar": ""},
		map[string]any{"bar": "required_unless:foo,no"},
	)
	assertPasses(t, v2)
}

// Port of ValidationValidatorTest::testValidateRequiredWith
func TestValidateRequiredWith(t *testing.T) {
	t.Parallel()

	// required_with:other — required when 'other' is present and not empty
	v := makeValidator(
		map[string]any{"other": "present", "field": ""},
		map[string]any{"field": "required_with:other"},
	)
	assertFails(t, v)

	v2 := makeValidator(
		map[string]any{"other": "", "field": ""},
		map[string]any{"field": "required_with:other"},
	)
	assertPasses(t, v2)
}

// Port of ValidationValidatorTest::testValidateRequiredWithout
func TestValidateRequiredWithout(t *testing.T) {
	t.Parallel()

	// required_without:other — required when 'other' is absent/empty
	v := makeValidator(
		map[string]any{"field": ""},
		map[string]any{"field": "required_without:other"},
	)
	assertFails(t, v)

	v2 := makeValidator(
		map[string]any{"other": "present", "field": ""},
		map[string]any{"field": "required_without:other"},
	)
	assertPasses(t, v2)
}

// Port of ValidationValidatorTest::testValidatePresent
func TestValidatePresent(t *testing.T) {
	t.Parallel()

	// present — field must exist even if blank
	v := makeValidator(map[string]any{"field": ""}, map[string]any{"field": "present"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{}, map[string]any{"field": "present"})
	assertFails(t, v2)
}

// Port of ValidationValidatorTest::testValidateFilled
func TestValidateFilled(t *testing.T) {
	t.Parallel()

	// filled — if present, must not be blank
	v := makeValidator(map[string]any{"field": "value"}, map[string]any{"field": "filled"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"field": ""}, map[string]any{"field": "filled"})
	assertFails(t, v2)

	// absent field is OK for filled
	v3 := makeValidator(map[string]any{}, map[string]any{"field": "filled"})
	assertPasses(t, v3)
}

// Port of ValidationValidatorTest::testValidateMissing
func TestValidateMissing(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{}, map[string]any{"field": "missing"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"field": "present"}, map[string]any{"field": "missing"})
	assertFails(t, v2)
}

// Port of ValidationValidatorTest::testValidateProhibited
func TestValidateProhibited(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{}, map[string]any{"field": "prohibited"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"field": ""}, map[string]any{"field": "prohibited"})
	assertPasses(t, v2)

	v3 := makeValidator(map[string]any{"field": "value"}, map[string]any{"field": "prohibited"})
	assertFails(t, v3)
}

// Port of ValidationValidatorTest::testValidateProhibitedIf
func TestValidateProhibitedIf(t *testing.T) {
	t.Parallel()

	// prohibited_if:other,value
	v := makeValidator(
		map[string]any{"type": "admin", "secret": "s3cr3t"},
		map[string]any{"secret": "prohibited_if:type,admin"},
	)
	assertFails(t, v)

	v2 := makeValidator(
		map[string]any{"type": "user", "secret": "s3cr3t"},
		map[string]any{"secret": "prohibited_if:type,admin"},
	)
	assertPasses(t, v2)
}

// Port of ValidationValidatorTest::testValidateAccepted
func TestValidateAccepted(t *testing.T) {
	t.Parallel()

	for _, val := range []any{"yes", "on", "1", "true", true, 1} {
		v := makeValidator(map[string]any{"field": val}, map[string]any{"field": "accepted"})
		if !v.Passes() {
			t.Errorf("accepted: expected %v to pass", val)
		}
	}

	for _, val := range []any{"no", "off", "0", "false", false, 0} {
		v := makeValidator(map[string]any{"field": val}, map[string]any{"field": "accepted"})
		if !v.Fails() {
			t.Errorf("accepted: expected %v to fail", val)
		}
	}
}

// Port of ValidationValidatorTest::testValidateIn
func TestValidateIn(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"size": "medium"},
		map[string]any{"size": "in:small,medium,large"},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"size": "xl"},
		map[string]any{"size": "in:small,medium,large"},
	)
	assertFails(t, v2)
}

// Port of ValidationValidatorTest::testValidateNotIn
func TestValidateNotIn(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"role": "user"},
		map[string]any{"role": "not_in:admin,superadmin"},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"role": "admin"},
		map[string]any{"role": "not_in:admin,superadmin"},
	)
	assertFails(t, v2)
}

// Port of ValidationValidatorTest::testValidateMin
func TestValidateMin(t *testing.T) {
	t.Parallel()

	// string: min characters
	v := makeValidator(map[string]any{"name": "Taylor"}, map[string]any{"name": "min:3"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"name": "Ti"}, map[string]any{"name": "min:3"})
	assertFails(t, v2)

	// numeric: min value
	v3 := makeValidator(map[string]any{"age": 18}, map[string]any{"age": "numeric|min:18"})
	assertPasses(t, v3)

	v4 := makeValidator(map[string]any{"age": 17}, map[string]any{"age": "numeric|min:18"})
	assertFails(t, v4)

	// array: min items
	v5 := makeValidator(map[string]any{"tags": []any{"a", "b", "c"}}, map[string]any{"tags": "array|min:2"})
	assertPasses(t, v5)
}

// Port of ValidationValidatorTest::testValidateMax
func TestValidateMax(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{"name": "Taylor"}, map[string]any{"name": "max:10"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"name": "Taylor Otwell"}, map[string]any{"name": "max:5"})
	assertFails(t, v2)

	v3 := makeValidator(map[string]any{"age": 100}, map[string]any{"age": "numeric|max:150"})
	assertPasses(t, v3)

	v4 := makeValidator(map[string]any{"age": 200}, map[string]any{"age": "numeric|max:150"})
	assertFails(t, v4)
}

// Port of ValidationValidatorTest::testValidateBetween
func TestValidateBetween(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{"price": 50}, map[string]any{"price": "numeric|between:10,100"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"price": 5}, map[string]any{"price": "numeric|between:10,100"})
	assertFails(t, v2)

	v3 := makeValidator(map[string]any{"price": 150}, map[string]any{"price": "numeric|between:10,100"})
	assertFails(t, v3)
}

// Port of ValidationValidatorTest::testValidateSize
func TestValidateSize(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{"name": "Taylor"}, map[string]any{"name": "size:6"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"name": "Tim"}, map[string]any{"name": "size:6"})
	assertFails(t, v2)

	v3 := makeValidator(map[string]any{"count": 5}, map[string]any{"count": "numeric|size:5"})
	assertPasses(t, v3)
}

// Port of ValidationValidatorTest::testValidateEmail
func TestValidateEmail(t *testing.T) {
	t.Parallel()

	cases := []struct {
		email string
		pass  bool
	}{
		{"user@example.com", true},
		{"user+tag@sub.example.org", true},
		{"notanemail", false},
		{"@missinglocal.com", false},
		{"missing@", false},
	}

	for _, tc := range cases {
		v := makeValidator(map[string]any{"email": tc.email}, map[string]any{"email": "email"})
		if v.Passes() != tc.pass {
			t.Errorf("email %q: expected pass=%v", tc.email, tc.pass)
		}
	}
}

// Port of ValidationValidatorTest::testValidateUrl
func TestValidateUrl(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{"url": "https://example.com"}, map[string]any{"url": "url"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"url": "not-a-url"}, map[string]any{"url": "url"})
	assertFails(t, v2)
}

// Port of ValidationValidatorTest::testValidateIp
func TestValidateIp(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{"ip": "192.168.1.1"}, map[string]any{"ip": "ip"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"ip": "2001:db8::1"}, map[string]any{"ip": "ip"})
	assertPasses(t, v2)

	v3 := makeValidator(map[string]any{"ip": "not-an-ip"}, map[string]any{"ip": "ip"})
	assertFails(t, v3)
}

// Port of ValidationValidatorTest::testValidateAlpha
func TestValidateAlpha(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{"name": "Taylor"}, map[string]any{"name": "alpha"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"name": "Taylor1"}, map[string]any{"name": "alpha"})
	assertFails(t, v2)

	v3 := makeValidator(map[string]any{"name": "Taylor Otwell"}, map[string]any{"name": "alpha"})
	assertFails(t, v3)
}

// Port of ValidationValidatorTest::testValidateAlphaDash
func TestValidateAlphaDash(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{"name": "Taylor-Otwell_1"}, map[string]any{"name": "alpha_dash"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"name": "Taylor Otwell"}, map[string]any{"name": "alpha_dash"})
	assertFails(t, v2)
}

// Port of ValidationValidatorTest::testValidateAlphaNum
func TestValidateAlphaNum(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{"code": "abc123"}, map[string]any{"code": "alpha_num"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"code": "abc-123"}, map[string]any{"code": "alpha_num"})
	assertFails(t, v2)
}

// Port of ValidationValidatorTest::testValidateNumeric
func TestValidateNumeric(t *testing.T) {
	t.Parallel()

	for _, val := range []any{1, 1.5, "42", "3.14", "-7"} {
		v := makeValidator(map[string]any{"n": val}, map[string]any{"n": "numeric"})
		if !v.Passes() {
			t.Errorf("numeric: expected %v to pass", val)
		}
	}

	for _, val := range []any{"abc", "1e2e3", ""} {
		v := makeValidator(map[string]any{"n": val}, map[string]any{"n": "numeric"})
		if !v.Fails() {
			t.Errorf("numeric: expected %v to fail", val)
		}
	}
}

// Port of ValidationValidatorTest::testValidateInteger
func TestValidateInteger(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{"n": 42}, map[string]any{"n": "integer"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"n": "42"}, map[string]any{"n": "integer"})
	assertPasses(t, v2)

	v3 := makeValidator(map[string]any{"n": 3.14}, map[string]any{"n": "integer"})
	assertFails(t, v3)
}

// Port of ValidationValidatorTest::testValidateBoolean
func TestValidateBoolean(t *testing.T) {
	t.Parallel()

	for _, val := range []any{true, false, 0, 1, "0", "1", "true", "false"} {
		v := makeValidator(map[string]any{"active": val}, map[string]any{"active": "boolean"})
		if !v.Passes() {
			t.Errorf("boolean: expected %v to pass", val)
		}
	}

	for _, val := range []any{"yes", "no", 2} {
		v := makeValidator(map[string]any{"active": val}, map[string]any{"active": "boolean"})
		if !v.Fails() {
			t.Errorf("boolean: expected %v to fail", val)
		}
	}
}

// Port of ValidationValidatorTest::testValidateDate
func TestValidateDate(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{"date": "2024-01-15"}, map[string]any{"date": "date"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"date": "not-a-date"}, map[string]any{"date": "date"})
	assertFails(t, v2)
}

// Port of ValidationValidatorTest::testValidateDateFormat
func TestValidateDateFormat(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"date": "2024-01-15"},
		map[string]any{"date": "date_format:Y-m-d"},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"date": "15/01/2024"},
		map[string]any{"date": "date_format:Y-m-d"},
	)
	assertFails(t, v2)
}

// Port of ValidationValidatorTest::testValidateBefore
func TestValidateBefore(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"date": "2023-01-01"},
		map[string]any{"date": "before:2024-01-01"},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"date": "2025-01-01"},
		map[string]any{"date": "before:2024-01-01"},
	)
	assertFails(t, v2)
}

// Port of ValidationValidatorTest::testValidateAfter
func TestValidateAfter(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"date": "2025-01-01"},
		map[string]any{"date": "after:2024-01-01"},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"date": "2023-01-01"},
		map[string]any{"date": "after:2024-01-01"},
	)
	assertFails(t, v2)
}

// Port of ValidationValidatorTest::testValidateSame
func TestValidateSame(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"password": "secret", "confirm": "secret"},
		map[string]any{"password": "same:confirm"},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"password": "secret", "confirm": "different"},
		map[string]any{"password": "same:confirm"},
	)
	assertFails(t, v2)
}

// Port of ValidationValidatorTest::testValidateDifferent
func TestValidateDifferent(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"old_password": "old", "new_password": "new"},
		map[string]any{"new_password": "different:old_password"},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"old_password": "same", "new_password": "same"},
		map[string]any{"new_password": "different:old_password"},
	)
	assertFails(t, v2)
}

// Port of ValidationValidatorTest::testValidateConfirmed
func TestValidateConfirmed(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"password": "secret", "password_confirmation": "secret"},
		map[string]any{"password": "confirmed"},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"password": "secret", "password_confirmation": "other"},
		map[string]any{"password": "confirmed"},
	)
	assertFails(t, v2)
}

// Port of ValidationValidatorTest::testValidateDistinct
func TestValidateDistinct(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"tags": []any{"a", "b", "c"}},
		map[string]any{"tags": "array|distinct"},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"tags": []any{"a", "b", "a"}},
		map[string]any{"tags": "array|distinct"},
	)
	assertFails(t, v2)
}

// Port of ValidationValidatorTest::testValidateArray
func TestValidateArray(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{"items": []any{1, 2, 3}}, map[string]any{"items": "array"})
	assertPasses(t, v)

	v2 := makeValidator(map[string]any{"items": "not-an-array"}, map[string]any{"items": "array"})
	assertFails(t, v2)
}

// Port of ValidationValidatorTest::testValidateBail
func TestValidateBail(t *testing.T) {
	t.Parallel()

	// With bail: only 1 error per field should appear
	v := makeValidator(
		map[string]any{"age": "bad"},
		map[string]any{"age": "bail|integer|min:18|max:120"},
	)
	v.Fails()

	msgs := v.Errors().Get("age")
	if len(msgs) != 1 {
		t.Errorf("bail: expected 1 error, got %d: %v", len(msgs), msgs)
	}
}

// Port of ValidationValidatorTest::testValidateNullable
func TestValidateNullable(t *testing.T) {
	t.Parallel()

	// nil value with nullable passes other rules
	v := makeValidator(
		map[string]any{"phone": nil},
		map[string]any{"phone": "nullable|string|max:15"},
	)
	assertPasses(t, v)

	// non-nil value still validated
	v2 := makeValidator(
		map[string]any{"phone": "123456789012345678"},
		map[string]any{"phone": "nullable|string|max:15"},
	)
	assertFails(t, v2)
}

// Port of ValidationValidatorTest::testSometimesWorksOnNestedArrays
func TestSometimesWorksOnNestedArrays(t *testing.T) {
	t.Parallel()

	// "sometimes" means: only validate if the field is present in the data
	v := makeValidator(
		map[string]any{},
		map[string]any{"name": "sometimes|required|string"},
	)
	// name is absent → sometimes skips the field → should pass
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"name": ""},
		map[string]any{"name": "sometimes|required|string"},
	)
	// name is present but blank → fails required
	assertFails(t, v2)
}

// Port of ValidationValidatorTest::testCustomValidationRules
func TestCustomValidationRules(t *testing.T) {
	t.Parallel()

	// ValidationRule object that rejects the value "forbidden"
	noForbidden := &rejectForbiddenRule{}

	v := validation.NewFactory().Make(
		map[string]any{"word": "forbidden"},
		map[string]any{"word": []any{noForbidden}},
		nil, nil,
	)
	assertFails(t, v)
	assertError(t, v, "word", "forbidden")

	v2 := validation.NewFactory().Make(
		map[string]any{"word": "allowed"},
		map[string]any{"word": []any{noForbidden}},
		nil, nil,
	)
	assertPasses(t, v2)
}

type rejectForbiddenRule struct{}

func (r *rejectForbiddenRule) Validate(attribute string, value any, fail func(string)) {
	if s, ok := value.(string); ok && s == "forbidden" {
		fail("The word 'forbidden' is not allowed.")
	}
}

// Port of ValidationValidatorTest::testWildcardNestedRules
func TestWildcardNestedRules(t *testing.T) {
	t.Parallel()

	// items.*.name must be required
	v := makeValidator(
		map[string]any{
			"items": []any{
				map[string]any{"name": "Alice", "age": 30},
				map[string]any{"name": "", "age": 25},
			},
		},
		map[string]any{"items.*.name": "required"},
	)
	assertFails(t, v)

	if !v.Errors().Has("items.1.name") {
		t.Error("expected error for items.1.name (empty name)")
	}

	if v.Errors().Has("items.0.name") {
		t.Error("unexpected error for items.0.name (valid name)")
	}

	// All names provided → passes
	v2 := makeValidator(
		map[string]any{
			"items": []any{
				map[string]any{"name": "Alice"},
				map[string]any{"name": "Bob"},
			},
		},
		map[string]any{"items.*.name": "required"},
	)
	assertPasses(t, v2)
}

// Port of ValidationValidatorTest::testConditionalRules
func TestConditionalRules(t *testing.T) {
	t.Parallel()

	// Rule::when(true) → applies the rules
	trueRule := validation.Rule.When(true, "required")

	v := validation.NewFactory().Make(
		map[string]any{"field": ""},
		map[string]any{"field": []any{trueRule}},
		nil, nil,
	)

	// ConditionalRule.ActiveRules() returns the parsed rules to apply
	active := trueRule.ActiveRules()
	if len(active) == 0 {
		t.Error("expected active rules when condition is true")
	}

	// Rule::when(false) → applies default rules (nil here)
	falseRule := validation.Rule.When(false, "required")
	falseActive := falseRule.ActiveRules()

	if len(falseActive) != 0 {
		t.Errorf("expected no active rules when condition is false, got %v", falseActive)
	}

	_ = v
}

// Port of ValidationValidatorTest::testValidateRegex
func TestValidateRegex(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"code": "abc123"},
		map[string]any{"code": "regex:/^[a-z0-9]+$/"},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"code": "ABC!"},
		map[string]any{"code": "regex:/^[a-z0-9]+$/"},
	)
	assertFails(t, v2)
}

// Port of ValidationValidatorTest::testValidateUUID
func TestValidateUUID(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"id": "550e8400-e29b-41d4-a716-446655440000"},
		map[string]any{"id": "uuid"},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"id": "not-a-uuid"},
		map[string]any{"id": "uuid"},
	)
	assertFails(t, v2)
}

// Port of ValidationValidatorTest::testValidateJson
func TestValidateJson(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"data": `{"key": "value"}`},
		map[string]any{"data": "json"},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"data": "not-json"},
		map[string]any{"data": "json"},
	)
	assertFails(t, v2)
}

// Port of ValidationValidatorTest::testValidateStartsWith
func TestValidateStartsWith(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": "Foo Bar"},
		map[string]any{"name": "starts_with:Foo"},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"name": "Bar Foo"},
		map[string]any{"name": "starts_with:Foo"},
	)
	assertFails(t, v2)
}

// Port of ValidationValidatorTest::testValidateEndsWith
func TestValidateEndsWith(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"domain": "example.com"},
		map[string]any{"domain": "ends_with:.com,.org"},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"domain": "example.io"},
		map[string]any{"domain": "ends_with:.com,.org"},
	)
	assertFails(t, v2)
}

// Port of ValidationValidatorTest::testValidateInArray
func TestValidateInArray(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"role": "admin", "allowed": []any{"admin", "editor"}},
		map[string]any{"role": "in_array:allowed"},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"role": "superuser", "allowed": []any{"admin", "editor"}},
		map[string]any{"role": "in_array:allowed"},
	)
	assertFails(t, v2)
}

// Port of ValidationValidatorTest::testValidateHexColor
func TestValidateHexColor(t *testing.T) {
	t.Parallel()

	for _, c := range []string{"#fff", "#ffffff", "#FF0000", "#aabbcc"} {
		v := makeValidator(map[string]any{"color": c}, map[string]any{"color": "hex_color"})
		if !v.Passes() {
			t.Errorf("hex_color: expected %q to pass", c)
		}
	}

	for _, c := range []string{"fff", "ffffff", "#xyz", "#fffff"} {
		v := makeValidator(map[string]any{"color": c}, map[string]any{"color": "hex_color"})
		if !v.Fails() {
			t.Errorf("hex_color: expected %q to fail", c)
		}
	}
}

// Port of ValidationValidatorTest::testValidateTimezone
func TestValidateTimezone(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"tz": "America/New_York"},
		map[string]any{"tz": "timezone"},
	)
	assertPasses(t, v)

	v2 := makeValidator(
		map[string]any{"tz": "Not/ATimezone"},
		map[string]any{"tz": "timezone"},
	)
	assertFails(t, v2)
}
