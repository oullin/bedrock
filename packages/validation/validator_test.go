package validation_test

import (
	"testing"

	"github.com/bedrock/packages/validation"
)

// ─── helpers ──────────────────────────────────────────────────────────────────

func makeValidator(data, rules map[string]any) *validation.Validator {
	return validation.NewFactory().Make(data, rules, nil, nil)
}

func assertPasses(t *testing.T, v *validation.Validator) {
	t.Helper()

	if !v.Passes() {
		t.Errorf("expected validator to pass, got errors: %v", v.Errors().All())
	}
}

func assertFails(t *testing.T, v *validation.Validator) {
	t.Helper()

	if !v.Fails() {
		t.Error("expected validator to fail, but it passed")
	}
}

func assertError(t *testing.T, v *validation.Validator, attribute, ruleContains string) {
	t.Helper()

	msgs := v.Errors().Get(attribute)
	if len(msgs) == 0 {
		t.Errorf("expected error for %q, got none", attribute)
		return
	}

	for _, m := range msgs {
		if containsStr(m, ruleContains) {
			return
		}
	}

	t.Errorf("expected error for %q to contain %q, got %v", attribute, ruleContains, msgs)
}

func containsStr(s, sub string) bool {
	return len(s) >= len(sub) && (sub == "" || findSub(s, sub))
}

func findSub(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}

	return false
}

// ─── basic tests ──────────────────────────────────────────────────────────────

func TestValidator_PassesWhenNoRules(t *testing.T) {
	t.Parallel()

	v := makeValidator(map[string]any{"name": "Alice"}, nil)
	assertPasses(t, v)
}

func TestValidator_PassesWhenAllRulesPass(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"email": "user@example.com", "age": 25},
		map[string]any{"email": "required|email", "age": "required|integer|min:18"},
	)
	assertPasses(t, v)
}

func TestValidator_FailsWhenAnyRuleFails(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"email": "not-an-email"},
		map[string]any{"email": "required|email"},
	)
	assertFails(t, v)
}

func TestValidator_Errors(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"email": "not-an-email"},
		map[string]any{"email": "required|email"},
	)
	v.Fails()

	if v.Errors().IsEmpty() {
		t.Error("expected errors, got empty bag")
	}

	msgs := v.Errors().Get("email")
	if len(msgs) == 0 {
		t.Error("expected error for 'email'")
	}
}

func TestValidator_Validate_ReturnsNilOnSuccess(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": "Alice"},
		map[string]any{"name": "required|string"},
	)

	if err := v.Validate(); err != nil {
		t.Errorf("Validate: expected nil, got %v", err)
	}
}

func TestValidator_Validate_ReturnsErrorOnFailure(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": ""},
		map[string]any{"name": "required"},
	)

	err := v.Validate()
	if err == nil {
		t.Error("Validate: expected error, got nil")
	}

	var ve *validation.ValidationException
	if !isValidationException(err, &ve) {
		t.Errorf("Validate: expected *ValidationException, got %T", err)
	}
}

func isValidationException(err error, ve **validation.ValidationException) bool {
	if e, ok := err.(*validation.ValidationException); ok {
		*ve = e
		return true
	}

	return false
}

func TestValidator_Validated_ReturnsSubset(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"email": "user@example.com", "extra": "ignored"},
		map[string]any{"email": "required|email"},
	)

	vd, err := v.Validated()
	if err != nil {
		t.Fatalf("Validated: %v", err)
	}

	if _, ok := vd["email"]; !ok {
		t.Error("Validated: expected 'email' in result")
	}
}

func TestValidator_Failed(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"email": "bad", "age": "not-a-number"},
		map[string]any{"email": "email", "age": "integer"},
	)
	v.Fails()

	failed := v.Failed()
	if len(failed["email"]) == 0 {
		t.Error("Failed: expected 'email' in failed map")
	}

	if len(failed["age"]) == 0 {
		t.Error("Failed: expected 'age' in failed map")
	}
}

func TestValidator_HasRule(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{},
		map[string]any{"email": "required|email|max:255"},
	)

	if !v.HasRule("email", "required") {
		t.Error("HasRule: expected 'required' to be present")
	}

	if !v.HasRule("email", "Max") {
		t.Error("HasRule: expected 'Max' to be present")
	}

	if v.HasRule("email", "integer") {
		t.Error("HasRule: expected 'integer' to be absent")
	}
}

func TestValidator_SetData(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": ""},
		map[string]any{"name": "required"},
	)
	assertFails(t, v)

	v.SetData(map[string]any{"name": "Alice"})
	assertPasses(t, v)
}

func TestValidator_CustomMessages(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": ""},
		map[string]any{"name": "required"},
	)
	v.SetCustomMessages(map[string]string{
		"required": "This field cannot be empty.",
	})
	v.Fails()

	msgs := v.Errors().Get("name")
	if len(msgs) == 0 || msgs[0] != "This field cannot be empty." {
		t.Errorf("custom message: got %v", msgs)
	}
}

func TestValidator_AttributeNames(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"user_email": ""},
		map[string]any{"user_email": "required"},
	)
	v.SetAttributeNames(map[string]string{"user_email": "Email Address"})
	v.Fails()

	msgs := v.Errors().Get("user_email")
	if len(msgs) == 0 {
		t.Fatal("expected error for user_email")
	}

	if !containsStr(msgs[0], "Email Address") {
		t.Errorf("expected attribute name in message, got: %s", msgs[0])
	}
}

func TestValidator_WildcardRules(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{
			"items": []any{
				map[string]any{"name": "A"},
				map[string]any{"name": ""},
			},
		},
		map[string]any{"items.*.name": "required"},
	)
	assertFails(t, v)

	// items.0.name passes, items.1.name fails
	if v.Errors().Has("items.1.name") == false {
		t.Error("expected error for items.1.name")
	}
}

func TestValidator_BailStopsOnFirstFailure(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"age": "not-a-number"},
		map[string]any{"age": "bail|integer|min:18"},
	)
	v.Fails()

	// With bail, only one error should be reported for the field
	msgs := v.Errors().Get("age")
	if len(msgs) != 1 {
		t.Errorf("bail: expected 1 error, got %d: %v", len(msgs), msgs)
	}
}

func TestValidator_NullableSkipsRules(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"phone": nil},
		map[string]any{"phone": "nullable|string|max:15"},
	)
	assertPasses(t, v)
}

func TestValidator_CustomRuleObject(t *testing.T) {
	t.Parallel()

	always := &alwaysFailRule{msg: "custom always fails"}

	v := validation.NewFactory().Make(
		map[string]any{"field": "anything"},
		map[string]any{"field": []any{always}},
		nil, nil,
	)
	assertFails(t, v)
	assertError(t, v, "field", "custom always fails")
}

type alwaysFailRule struct{ msg string }

func (r *alwaysFailRule) Validate(_ string, _ any, fail func(string)) {
	fail(r.msg)
}

func TestValidator_AddExtension(t *testing.T) {
	t.Parallel()

	magicCode := validation.RuleFunc(func(attr string, value any, params []string, ctx validation.RuleContext) bool {
		s, ok := value.(string)
		return ok && s == "SECRET"
	})

	v := makeValidator(
		map[string]any{"code": "INVALID"},
		map[string]any{"code": "required|magic_code"},
	)
	v.AddExtension("magic_code", magicCode)
	assertFails(t, v)

	v2 := makeValidator(
		map[string]any{"code": "SECRET"},
		map[string]any{"code": "required|magic_code"},
	)
	v2.AddExtension("magic_code", magicCode)
	assertPasses(t, v2)
}
