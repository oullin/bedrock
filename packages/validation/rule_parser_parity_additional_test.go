package validation_test

import (
	"reflect"
	"testing"

	"github.com/bedrock/packages/validation"
)

// ValidationRuleParserTest::testEmptyRulesArePreserved
func TestRuleParser_EmptyRulesArePreserved(t *testing.T) {
	t.Parallel()

	got := validation.Parse("")

	if got.Name != "" {
		t.Fatalf("Parse(\"\").Name = %q, want empty", got.Name)
	}
}

// ValidationRuleParserTest::testEmptyRulesCanBeExploded
func TestRuleParser_EmptyRulesCanBeExploded(t *testing.T) {
	t.Parallel()

	got := validation.Explode("")

	if len(got) != 0 {
		t.Fatalf("Explode(\"\") = %v, want no rules", got)
	}
}

// ValidationRuleParserTest::testExplodeProperlyParsesSingleRegexRuleNotContainingPipe
func TestRuleParser_ExplodeProperlyParsesSingleRegexRuleNotContainingPipe(t *testing.T) {
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

// ValidationRuleParserTest::testExplodeProperlyParsesRegexWithArrayOfRules
func TestRuleParser_ExplodeProperlyParsesRegexWithArrayOfRules(t *testing.T) {
	t.Parallel()

	got := validation.Explode([]any{"required", "regex:/^[a-z]+$/", "email"})

	if len(got) != 3 {
		t.Fatalf("Explode([]any) = %v, want 3 rules", got)
	}

	if got[1].Name != "Regex" {
		t.Fatalf("Explode([]any)[1].Name = %q, want Regex", got[1].Name)
	}
}

// ValidationRuleParserTest::testExplodeHandlesNumericStringRule
func TestRuleParser_ExplodeHandlesNumericStringRule(t *testing.T) {
	t.Parallel()

	got := validation.Parse("numeric")

	if got.Name != "Numeric" {
		t.Fatalf("Parse(numeric).Name = %q, want Numeric", got.Name)
	}
}

// ValidationRuleParserTest::testExplodeHandlesStringDateRule
func TestRuleParser_ExplodeHandlesStringDateRule(t *testing.T) {
	t.Parallel()

	got := validation.Parse("date")

	if got.Name != "Date" {
		t.Fatalf("Parse(date).Name = %q, want Date", got.Name)
	}
}

// ValidationRuleParserTest::testConditionalRulesAreProperlyExpandedAndFiltered
func TestRuleParser_ConditionalRulesAreProperlyExpandedAndFiltered(t *testing.T) {
	t.Parallel()

	whenRules := validation.Rule.When(true, "required|min:2").ActiveRules()

	if len(whenRules) != 2 || whenRules[0].Name != "Required" || whenRules[1].Name != "Min" {
		t.Fatalf("Rule.When(true, ...).ActiveRules() = %v, want Required and Min", whenRules)
	}

	unlessRules := validation.Rule.Unless(true, "required|min:2", "string|max:10").ActiveRules()

	if len(unlessRules) != 2 || unlessRules[0].Name != "String" || unlessRules[1].Name != "Max" {
		t.Fatalf("Rule.Unless(true, ...).ActiveRules() = %v, want String and Max", unlessRules)
	}
}

// ValidationRuleParserTest::testConditionalRulesWithDefault
func TestRuleParser_ConditionalRulesWithDefault(t *testing.T) {
	t.Parallel()

	rules := validation.Rule.When(false, "required|min:2", "string|max:10").ActiveRules()

	if len(rules) != 2 || rules[0].Name != "String" || rules[1].Name != "Max" {
		t.Fatalf("Rule.When(false, ...).ActiveRules() = %v, want String and Max", rules)
	}
}

// ValidationRuleParserTest::testExplodeHandlesDateRule
func TestRuleParser_ExplodeHandlesDateRule(t *testing.T) {
	t.Parallel()

	got := validation.Parse("date")

	if got.Name != "Date" {
		t.Fatalf("Parse(date).Name = %q, want Date", got.Name)
	}
}

// ValidationRuleParserTest::testExplodeHandlesNumericRule
func TestRuleParser_ExplodeHandlesNumericRule(t *testing.T) {
	t.Parallel()

	got := validation.Parse("numeric")

	if got.Name != "Numeric" {
		t.Fatalf("Parse(numeric).Name = %q, want Numeric", got.Name)
	}
}
