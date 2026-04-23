package validation_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/bedrock/packages/validation"
	"github.com/bedrock/packages/validation/rules"
)

// ValidationRuleParserTest::testEmptyConditionalRulesArePreserved
func TestInventoryRuleParser_EmptyConditionalRulesArePreserved(t *testing.T) {
	t.Parallel()

	whenRules := validation.Rule.When(true, "", "string|max:10").ActiveRules()
	unlessRules := validation.Rule.Unless(true, "required|min:2").ActiveRules()

	if len(whenRules) != 0 {
		t.Fatalf("Rule.When(true, empty).ActiveRules() = %v, want no rules", whenRules)
	}

	if len(unlessRules) != 0 {
		t.Fatalf("Rule.Unless(true, ...).ActiveRules() = %v, want no rules", unlessRules)
	}
}

// ValidationRuleParserTest::testExplodeFailsParsingSingleRegexRuleContainingPipe
func TestInventoryRuleParser_ExplodeFailsParsingSingleRegexRuleContainingPipe(t *testing.T) {
	t.Parallel()

	got := validation.Explode("regex:/^(foo|bar)$/i")

	if len(got) != 2 {
		t.Fatalf("Explode(regex with pipe) = %v, want 2 split rules", got)
	}

	if got[0].Name != "Regex" || !reflect.DeepEqual(got[0].Parameters, []string{"/^(foo"}) {
		t.Fatalf("Explode(regex with pipe)[0] = %#v, want Regex with truncated parameter", got[0])
	}

	if got[1].Name != "Bar)$/i" {
		t.Fatalf("Explode(regex with pipe)[1].Name = %q, want Bar)$/i", got[1].Name)
	}
}

// ValidationRuleParserTest::testExplodeFailsParsingRegexWithOtherRulesInSingleString
func TestInventoryRuleParser_ExplodeFailsParsingRegexWithOtherRulesInSingleString(t *testing.T) {
	t.Parallel()

	got := validation.Explode("in:foo|regex:/^(foo|bar)$/i")

	if len(got) != 3 {
		t.Fatalf("Explode(in plus regex with pipe) = %v, want 3 split rules", got)
	}

	if got[0].Name != "In" || !reflect.DeepEqual(got[0].Parameters, []string{"foo"}) {
		t.Fatalf("Explode(...)[0] = %#v, want In:foo", got[0])
	}

	if got[1].Name != "Regex" || !reflect.DeepEqual(got[1].Parameters, []string{"/^(foo"}) {
		t.Fatalf("Explode(...)[1] = %#v, want Regex with truncated parameter", got[1])
	}

	if got[2].Name != "Bar)$/i" {
		t.Fatalf("Explode(...)[2].Name = %q, want Bar)$/i", got[2].Name)
	}
}

// ValidationRuleParserTest::testExplodeHandlesForwardSlashesInWildcardRule
func TestInventoryRuleParser_ExplodeHandlesForwardSlashesInWildcardRule(t *testing.T) {
	t.Parallel()

	flat := validation.FlattenData(map[string]any{
		"redirects": map[string]any{
			"directory/subdirectory/file": []any{"directory/subdirectory/redirectedfile"},
		},
	})

	got := validation.ExpandWildcards("redirects.directory/subdirectory/file.*", flat)
	want := []string{"redirects.directory/subdirectory/file.0"}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("ExpandWildcards(...) = %v, want %v", got, want)
	}
}

// ValidationFactoryTest::testValidateCallsValidateOnTheValidator
func TestInventoryFactory_ValidateCallsValidateOnTheValidator(t *testing.T) {
	t.Parallel()

	f := validation.NewFactory()

	got, err := f.Validate(
		map[string]any{"foo": "bar", "baz": "boom"},
		map[string]any{"foo": "required"},
		nil,
		nil,
	)

	if err != nil {
		t.Fatalf("Validate() returned error: %v", err)
	}

	if !reflect.DeepEqual(got, map[string]any{"foo": "bar"}) {
		t.Fatalf("Validate() = %v, want only validated data", got)
	}

	_, err = f.Validate(
		map[string]any{"foo": ""},
		map[string]any{"foo": "required"},
		nil,
		nil,
	)

	if !errors.Is(err, validation.ErrValidationFailed) {
		t.Fatalf("Validate() error = %v, want ErrValidationFailed", err)
	}
}

// ValidationValidatorTest::testValidatedThrowsOnFailEvenAfterPassesCall
func TestInventoryValidator_ValidatedThrowsOnFailEvenAfterPassesCall(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": ""},
		map[string]any{"name": "required"},
	)

	if v.Passes() {
		t.Fatal("Passes() = true, want false")
	}

	if _, err := v.Validated(); !errors.Is(err, validation.ErrValidationFailed) {
		t.Fatalf("Validated() error = %v, want ErrValidationFailed", err)
	}
}

// ValidationValidatorTest::testEmptyExistingAttributesAreValidated
func TestInventoryValidator_EmptyExistingAttributesAreValidated(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": ""},
		map[string]any{"name": "min:1"},
	)

	if !v.Fails() {
		t.Fatal("expected present empty string to be validated and fail min")
	}
}

// ValidationValidatorTest::testNullable
func TestInventoryValidator_Nullable(t *testing.T) {
	t.Parallel()

	nilValue := makeValidator(
		map[string]any{"name": nil},
		map[string]any{"name": "nullable|string|min:3"},
	)
	assertPasses(t, nilValue)

	nonNilValue := makeValidator(
		map[string]any{"name": 123},
		map[string]any{"name": "nullable|string"},
	)
	assertFails(t, nonNilValue)
}

// ValidationValidatorTest::testIfRulesAreSuccessfullyAdded
func TestInventoryValidator_IfRulesAreSuccessfullyAdded(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": "Taylor"},
		map[string]any{"name": "required"},
	)

	v.AddRules(map[string]any{"name": "email"})

	if !v.Fails() {
		t.Fatal("expected dynamically added email rule to fail")
	}

	if !reflect.DeepEqual(v.Failed()["name"], []string{"Email"}) {
		t.Fatalf("Failed()[name] = %v, want Email", v.Failed()["name"])
	}
}

// ValidationValidatorTest::testValidatedNotValidateTwiceData
func TestInventoryValidator_ValidatedNotValidateTwiceData(t *testing.T) {
	t.Parallel()

	calls := 0
	v := makeValidator(
		map[string]any{"name": "Taylor"},
		map[string]any{"name": "counted"},
	)
	v.AddExtension("counted", func(_ string, _ any, _ []string, _ rules.RuleContext) bool {
		calls++

		return true
	})

	if !v.Passes() {
		t.Fatalf("Passes() = false, errors: %v", v.Errors().All())
	}

	if _, err := v.Validated(); err != nil {
		t.Fatalf("Validated() returned error: %v", err)
	}

	if calls != 1 {
		t.Fatalf("custom rule calls = %d, want 1", calls)
	}
}

// ValidationValidatorTest::testMultiplePassesCalls
func TestInventoryValidator_MultiplePassesCalls(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": ""},
		map[string]any{"name": "required"},
	)

	if v.Passes() {
		t.Fatal("first Passes() = true, want false")
	}

	v.SetData(map[string]any{"name": "Taylor"})

	if !v.Passes() {
		t.Fatalf("second Passes() = false, errors: %v", v.Errors().All())
	}

	if !v.Errors().IsEmpty() {
		t.Fatalf("Errors() = %v, want empty after second pass", v.Errors().All())
	}
}

// ValidationValidatorTest::testExcludeBeforeADependentRule
func TestInventoryValidator_ExcludeBeforeADependentRule(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"type": "admin"},
		map[string]any{"secret": "exclude_if:type,admin|required"},
	)

	if !v.Passes() {
		t.Fatalf("expected excluded missing field not to run required, got %v", v.Errors().All())
	}
}

// ValidationValidatorTest::testExcludeValuesAreReallyRemoved
func TestInventoryValidator_ExcludeValuesAreReallyRemoved(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"type": "admin", "secret": "classified"},
		map[string]any{"type": "required", "secret": "exclude_if:type,admin|string"},
	)

	got, err := v.Validated()

	if err != nil {
		t.Fatalf("Validated() returned error: %v", err)
	}

	if _, ok := got["secret"]; ok {
		t.Fatalf("Validated() = %v, want secret excluded", got)
	}
}

// ValidationValidatorTest::testExcludeWithValuesAreReallyRemoved
func TestInventoryValidator_ExcludeWithValuesAreReallyRemoved(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"profile": "present", "avatar": "avatar.png"},
		map[string]any{"profile": "required", "avatar": "exclude_with:profile|string"},
	)

	got, err := v.Validated()

	if err != nil {
		t.Fatalf("Validated() returned error: %v", err)
	}

	if _, ok := got["avatar"]; ok {
		t.Fatalf("Validated() = %v, want avatar excluded", got)
	}
}

// ValidationExcludeIfTest::testExcludeIfRuleValidation
func TestInventoryExcludeIfRule_Validation(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"secret": "classified", "name": "Taylor"},
		map[string]any{"secret": []any{validation.Rule.ExcludeIf(func() bool { return true }), "required"}, "name": "required"},
	)

	got, err := v.Validated()

	if err != nil {
		t.Fatalf("Validated() returned error: %v", err)
	}

	if _, ok := got["secret"]; ok {
		t.Fatalf("Validated() = %v, want callback-excluded field removed", got)
	}
}
