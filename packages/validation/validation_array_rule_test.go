package validation_test

import (
	"testing"

	"github.com/bedrock/packages/validation"
)

// ValidationArrayRuleTest::testItCorrectlyFormatsAStringVersionOfTheRule
func TestArrayRule_StringAndValidation(t *testing.T) {
	t.Parallel()

	rule := validation.Rule.Array("name", "email")

	if got := rule.String(); got != "array:name,email" {
		t.Fatalf("String() = %q, want %q", got, "array:name,email")
	}

	v := validation.NewFactory().Make(
		map[string]any{"profile": map[string]any{"name": "Taylor", "email": "user@example.com"}},
		map[string]any{"profile": []any{rule}},
		nil,
		nil,
	)
	assertPasses(t, v)

	v2 := validation.NewFactory().Make(
		map[string]any{"profile": map[string]any{"name": "Taylor", "email": "user@example.com", "admin": true}},
		map[string]any{"profile": []any{rule}},
		nil,
		nil,
	)
	assertFails(t, v2)
}
