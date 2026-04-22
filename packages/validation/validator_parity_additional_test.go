package validation_test

import (
	"errors"
	"testing"

	"github.com/bedrock/packages/validation"
)

// ValidationValidatorTest::testValidateThrowsOnFail
func TestValidator_ValidateThrowsOnFail(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": ""},
		map[string]any{"name": "required"},
	)

	err := v.Validate()

	if err == nil {
		t.Fatal("expected Validate to return an error")
	}

	if !errors.Is(err, validation.ErrValidationFailed) {
		t.Fatalf("errors.Is(err, ErrValidationFailed) = false, want true")
	}
}

// ValidationValidatorTest::testValidateDoesntThrowOnPass
func TestValidator_ValidateDoesntThrowOnPass(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": "Taylor"},
		map[string]any{"name": "required|string"},
	)

	if err := v.Validate(); err != nil {
		t.Fatalf("Validate() = %v, want nil", err)
	}
}

// ValidationValidatorTest::testValidatedThrowsOnFail
func TestValidator_ValidatedThrowsOnFail(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": ""},
		map[string]any{"name": "required"},
	)

	if _, err := v.Validated(); err == nil {
		t.Fatal("expected Validated to return an error")
	}
}

// ValidationValidatorTest::testValidatedDoesntThrowOnPass
func TestValidator_ValidatedDoesntThrowOnPass(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"email": "user@example.com", "extra": "ignored"},
		map[string]any{"email": "required|email"},
	)

	got, err := v.Validated()

	if err != nil {
		t.Fatalf("Validated() = %v", err)
	}

	if got["email"] != "user@example.com" {
		t.Fatalf("Validated() = %v, want email to be preserved", got)
	}
}

// ValidationValidatorTest::testHasFailedValidationRules
func TestValidator_HasFailedValidationRules(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"age": "not-a-number"},
		map[string]any{"age": "numeric"},
	)
	v.Fails()

	failed := v.Failed()

	if len(failed["age"]) == 0 {
		t.Fatal("expected failed rules for age")
	}
}

// ValidationValidatorTest::testHasNotFailedValidationRules
func TestValidator_HasNotFailedValidationRules(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"age": 42},
		map[string]any{"age": "numeric"},
	)
	v.Passes()

	if got := len(v.Failed()); got != 0 {
		t.Fatalf("Failed() = %d entries, want 0", got)
	}
}

// ValidationValidatorTest::testSometimesCanSkipRequiredRules
func TestValidator_SometimesCanSkipRequiredRules(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{},
		map[string]any{"name": "sometimes|required|string"},
	)

	if !v.Passes() {
		t.Fatalf("expected sometimes to skip an absent field, got %v", v.Errors().All())
	}

	v2 := makeValidator(
		map[string]any{"name": ""},
		map[string]any{"name": "sometimes|required|string"},
	)

	if !v2.Fails() {
		t.Fatal("expected present blank field to fail required")
	}
}

// ValidationValidatorTest::testInValidatableRulesReturnsValid
func TestValidator_InValidatableRulesReturnsValid(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": "Taylor"},
		map[string]any{"name": "made_up_rule"},
	)

	if !v.Passes() {
		t.Fatalf("expected unknown rules to be ignored, got %v", v.Errors().All())
	}
}

// ValidationValidatorTest::testValidateUsingNestedValidationRulesPasses
func TestValidator_ValidateUsingNestedValidationRulesPasses(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{
			"items": []any{
				map[string]any{"name": "Alice"},
				map[string]any{"name": "Bob"},
			},
		},
		map[string]any{"items.*.name": "required"},
	)

	if !v.Passes() {
		t.Fatalf("expected nested wildcard rules to pass, got %v", v.Errors().All())
	}
}

// ValidationValidatorTest::testValidateArrayKeys
func TestValidator_ValidateArrayKeys(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"profile": map[string]any{"name": "Taylor", "email": "user@example.com"}},
		map[string]any{"profile": "array:name,email"},
	)

	if !v.Passes() {
		t.Fatalf("expected allowed keys to pass, got %v", v.Errors().All())
	}

	v2 := makeValidator(
		map[string]any{"profile": map[string]any{"name": "Taylor", "email": "user@example.com", "admin": true}},
		map[string]any{"profile": "array:name,email"},
	)

	if !v2.Fails() {
		t.Fatal("expected disallowed keys to fail")
	}
}

// ValidationValidatorTest::testValidateList
func TestValidator_ValidateList(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"items": []any{"a", "b", "c"}},
		map[string]any{"items": "list"},
	)

	if !v.Passes() {
		t.Fatalf("expected list value to pass, got %v", v.Errors().All())
	}

	v2 := makeValidator(
		map[string]any{"items": map[string]any{"0": "a", "1": "b"}},
		map[string]any{"items": "list"},
	)

	if !v2.Fails() {
		t.Fatal("expected map value to fail list validation")
	}
}

// ValidationValidatorTest::testValidateAcceptedIf
func TestValidator_ValidateAcceptedIf(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"status": "published", "publish": "yes"},
		map[string]any{"publish": "accepted_if:status,published"},
	)

	if !v.Passes() {
		t.Fatalf("expected accepted_if to pass, got %v", v.Errors().All())
	}

	v2 := makeValidator(
		map[string]any{"status": "published", "publish": "no"},
		map[string]any{"publish": "accepted_if:status,published"},
	)

	if !v2.Fails() {
		t.Fatal("expected accepted_if to fail when the other field matches")
	}
}

// ValidationValidatorTest::testValidateRequiredAcceptedIf
func TestValidator_ValidateRequiredAcceptedIf(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"notify": "yes", "email": ""},
		map[string]any{"email": "required_if_accepted:notify"},
	)

	if !v.Fails() {
		t.Fatal("expected required_if_accepted to fail when the driver field is accepted")
	}

	v2 := makeValidator(
		map[string]any{"notify": "no", "email": ""},
		map[string]any{"email": "required_if_accepted:notify"},
	)

	if !v2.Passes() {
		t.Fatalf("expected required_if_accepted to pass when the driver field is not accepted, got %v", v2.Errors().All())
	}
}

// ValidationValidatorTest::testValidateRequiredIfDeclined
func TestValidator_ValidateRequiredIfDeclined(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"notify": "no", "email": ""},
		map[string]any{"email": "required_if_declined:notify"},
	)

	if !v.Fails() {
		t.Fatal("expected required_if_declined to fail when the driver field is declined")
	}

	v2 := makeValidator(
		map[string]any{"notify": "yes", "email": ""},
		map[string]any{"email": "required_if_declined:notify"},
	)

	if !v2.Passes() {
		t.Fatalf("expected required_if_declined to pass when the driver field is not declined, got %v", v2.Errors().All())
	}
}

// ValidationValidatorTest::testValidateDeclined
func TestValidator_ValidateDeclined(t *testing.T) {
	t.Parallel()

	for _, val := range []any{"no", "off", "0", "false", false, 0} {
		v := makeValidator(map[string]any{"publish": val}, map[string]any{"publish": "declined"})

		if !v.Passes() {
			t.Fatalf("expected declined to accept %v", val)
		}
	}

	for _, val := range []any{"yes", "on", "1", "true", true, 1} {
		v := makeValidator(map[string]any{"publish": val}, map[string]any{"publish": "declined"})

		if !v.Fails() {
			t.Fatalf("expected declined to reject %v", val)
		}
	}
}

// ValidationValidatorTest::testValidateMissingIf
func TestValidator_ValidateMissingIf(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"type": "internal", "token": "present"},
		map[string]any{"token": "missing_if:type,internal"},
	)

	if !v.Fails() {
		t.Fatal("expected missing_if to fail when the other field matches")
	}

	v2 := makeValidator(
		map[string]any{"type": "external", "token": "present"},
		map[string]any{"token": "missing_if:type,internal"},
	)

	if !v2.Passes() {
		t.Fatalf("expected missing_if to pass when the other field does not match, got %v", v2.Errors().All())
	}
}

// ValidationValidatorTest::testValidateMissingUnless
func TestValidator_ValidateMissingUnless(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"type": "internal", "token": "present"},
		map[string]any{"token": "missing_unless:type,external"},
	)

	if !v.Fails() {
		t.Fatal("expected missing_unless to fail when the other field is not in the allowed set")
	}

	v2 := makeValidator(
		map[string]any{"type": "external", "token": "present"},
		map[string]any{"token": "missing_unless:type,external"},
	)

	if !v2.Passes() {
		t.Fatalf("expected missing_unless to pass when the other field matches, got %v", v2.Errors().All())
	}
}

// ValidationValidatorTest::testValidateMissingWith
func TestValidator_ValidateMissingWith(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": "Taylor", "token": "present"},
		map[string]any{"token": "missing_with:name"},
	)

	if !v.Fails() {
		t.Fatal("expected missing_with to fail when the companion field is present")
	}

	v2 := makeValidator(
		map[string]any{"token": "present"},
		map[string]any{"token": "missing_with:name"},
	)

	if !v2.Passes() {
		t.Fatalf("expected missing_with to pass when the companion field is absent, got %v", v2.Errors().All())
	}
}

// ValidationValidatorTest::testValidateMissingWithAll
func TestValidator_ValidateMissingWithAll(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": "Taylor", "email": "user@example.com", "token": "present"},
		map[string]any{"token": "missing_with_all:name,email"},
	)

	if !v.Fails() {
		t.Fatal("expected missing_with_all to fail when all companion fields are present")
	}

	v2 := makeValidator(
		map[string]any{"name": "Taylor", "token": "present"},
		map[string]any{"token": "missing_with_all:name,email"},
	)

	if !v2.Passes() {
		t.Fatalf("expected missing_with_all to pass when one companion field is absent, got %v", v2.Errors().All())
	}
}

// ValidationValidatorTest::testValidateProhibitedAcceptedIf
func TestValidator_ValidateProhibitedAcceptedIf(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"notify": "yes", "token": "present"},
		map[string]any{"token": "prohibited_if_accepted:notify"},
	)

	if !v.Fails() {
		t.Fatal("expected prohibited_if_accepted to fail when the driver field is accepted")
	}

	v2 := makeValidator(
		map[string]any{"notify": "no", "token": "present"},
		map[string]any{"token": "prohibited_if_accepted:notify"},
	)

	if !v2.Passes() {
		t.Fatalf("expected prohibited_if_accepted to pass when the driver field is not accepted, got %v", v2.Errors().All())
	}
}

// ValidationValidatorTest::testValidateProhibitedDeclinedIf
func TestValidator_ValidateProhibitedDeclinedIf(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"notify": "no", "token": "present"},
		map[string]any{"token": "prohibited_if_declined:notify"},
	)

	if !v.Fails() {
		t.Fatal("expected prohibited_if_declined to fail when the driver field is declined")
	}

	v2 := makeValidator(
		map[string]any{"notify": "yes", "token": "present"},
		map[string]any{"token": "prohibited_if_declined:notify"},
	)

	if !v2.Passes() {
		t.Fatalf("expected prohibited_if_declined to pass when the driver field is not declined, got %v", v2.Errors().All())
	}
}

// ValidationValidatorTest::testValidateProhibitedUnless
func TestValidator_ValidateProhibitedUnless(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"role": "admin", "token": "present"},
		map[string]any{"token": "prohibited_unless:role,guest"},
	)

	if !v.Fails() {
		t.Fatal("expected prohibited_unless to fail when the other field is outside the allowed set")
	}

	v2 := makeValidator(
		map[string]any{"role": "guest", "token": "present"},
		map[string]any{"token": "prohibited_unless:role,guest"},
	)

	if !v2.Passes() {
		t.Fatalf("expected prohibited_unless to pass when the other field matches, got %v", v2.Errors().All())
	}
}

// ValidationValidatorTest::testProhibits
func TestValidator_Prohibits(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"token": "present", "secret": "also-present"},
		map[string]any{"token": "prohibits:secret"},
	)

	if !v.Fails() {
		t.Fatal("expected prohibits to fail when the forbidden field is present")
	}

	v2 := makeValidator(
		map[string]any{"token": "present"},
		map[string]any{"token": "prohibits:secret"},
	)

	if !v2.Passes() {
		t.Fatalf("expected prohibits to pass when the forbidden field is absent, got %v", v2.Errors().All())
	}
}

// ValidationValidatorTest::testProperMessagesAreReturnedForSizes
func TestValidator_ProperMessagesAreReturnedForSizes(t *testing.T) {
	t.Parallel()

	stringValidator := makeValidator(
		map[string]any{"name": "Taylor"},
		map[string]any{"name": "max:3"},
	)

	if !stringValidator.Fails() {
		t.Fatal("expected string max validation to fail")
	}

	if got := stringValidator.Errors().First("name"); got != "The name field must not be greater than 3 characters." {
		t.Fatalf("string max message = %q, want %q", got, "The name field must not be greater than 3 characters.")
	}

	numericValidator := makeValidator(
		map[string]any{"count": 9},
		map[string]any{"count": "max:3"},
	)

	if !numericValidator.Fails() {
		t.Fatal("expected numeric max validation to fail")
	}

	if got := numericValidator.Errors().First("count"); got != "The count field must not be greater than 3." {
		t.Fatalf("numeric max message = %q, want %q", got, "The count field must not be greater than 3.")
	}
}

// ValidationValidatorTest::testValidateGtPlaceHolderIsReplacedProperly
func TestValidator_ValidateGtPlaceHolderIsReplacedProperly(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"score": 5},
		map[string]any{"score": "gt:10"},
	)

	if !v.Fails() {
		t.Fatal("expected gt validation to fail")
	}

	if got := v.Errors().First("score"); got != "The score field must be greater than 10." {
		t.Fatalf("gt message = %q, want %q", got, "The score field must be greater than 10.")
	}
}
