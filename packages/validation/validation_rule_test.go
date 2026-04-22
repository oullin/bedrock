package validation_test

import (
	"testing"

	"github.com/bedrock/packages/validation"
)

// ValidationInRuleTest::testItCorrectlyFormatsAStringVersionOfTheRule
func TestInRuleStringAndValidation(t *testing.T) {
	t.Parallel()

	rule := validation.Rule.In("draft", "published")

	if got := rule.String(); got != "in:draft,published" {
		t.Fatalf("String() = %q, want %q", got, "in:draft,published")
	}

	v := validation.NewFactory().Make(
		map[string]any{"status": "published"},
		map[string]any{"status": []any{rule}},
		nil,
		nil,
	)

	assertPasses(t, v)

	v2 := validation.NewFactory().Make(
		map[string]any{"status": "archived"},
		map[string]any{"status": []any{rule}},
		nil,
		nil,
	)

	assertFails(t, v2)
}

// ValidationNotInRuleTest::testItCorrectlyFormatsAStringVersionOfTheRule
func TestNotInRuleStringAndValidation(t *testing.T) {
	t.Parallel()

	rule := validation.Rule.NotIn("draft", "published")

	if got := rule.String(); got != "not_in:draft,published" {
		t.Fatalf("String() = %q, want %q", got, "not_in:draft,published")
	}

	v := validation.NewFactory().Make(
		map[string]any{"status": "archived"},
		map[string]any{"status": []any{rule}},
		nil,
		nil,
	)

	assertPasses(t, v)

	v2 := validation.NewFactory().Make(
		map[string]any{"status": "published"},
		map[string]any{"status": []any{rule}},
		nil,
		nil,
	)

	assertFails(t, v2)
}
