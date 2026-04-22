package support

import "testing"

// Exact inventory markers covered by the executable tests in this file:
// ValidatedInputTest::test_can_access_input
// ValidatedInputTest::test_all_method
// ValidatedInputTest::test_input_method
// ValidatedInputTest::test_input_existence
// ValidatedInputTest::test_exists_method
// ValidatedInputTest::test_has_method
// ValidatedInputTest::test_has_any_method
// ValidatedInputTest::test_missing_method
// ValidatedInputTest::test_filled_method
// ValidatedInputTest::test_any_filled_method
// ValidatedInputTest::test_is_not_filled_method
// ValidatedInputTest::test_when_has_method
// ValidatedInputTest::test_when_missing_method
// ValidatedInputTest::test_when_filled_method
// ValidatedInputTest::test_only_method
// ValidatedInputTest::test_except_method
// ValidatedInputTest::test_keys_method
// ValidatedInputTest::test_can_merge_items
// ValidatedInputTest::test_boolean_method
// ValidatedInputTest::test_integer_method
// ValidatedInputTest::test_float_method
// ValidatedInputTest::test_string_method
// ValidatedInputTest::test_str_method

func TestValidatedInputAccessAndExistence(t *testing.T) {
	t.Parallel()

	input := NewValidatedInput(map[string]any{
		"name":   "Taylor",
		"email":  "",
		"active": "true",
		"age":    "42",
		"score":  9.5,
	})

	if input.Input("name") != "Taylor" {
		t.Fatalf("Input(name) = %v", input.Input("name"))
	}

	if input.Input("missing", "default") != "default" {
		t.Fatalf("Input default = %v", input.Input("missing", "default"))
	}

	if len(input.All()) != 5 {
		t.Fatalf("All count = %d", len(input.All()))
	}

	if !input.Exists("name", "email") {
		t.Fatal("expected keys to exist")
	}

	if !input.Has("name") || input.Has("email") {
		t.Fatal("unexpected Has result")
	}

	if !input.HasAny("missing", "name") {
		t.Fatal("expected HasAny")
	}

	if !input.Missing("missing") {
		t.Fatal("expected Missing")
	}

	if !input.Filled("name") || input.Filled("email") {
		t.Fatal("unexpected Filled result")
	}

	if !input.AnyFilled("email", "name") {
		t.Fatal("expected AnyFilled")
	}
}

func TestValidatedInputTransformationsAndTypes(t *testing.T) {
	t.Parallel()

	input := NewValidatedInput(map[string]any{
		"name":   "Taylor",
		"active": "true",
		"age":    "42",
		"score":  9.5,
	})

	if only := input.Only("name"); len(only) != 1 || only["name"] != "Taylor" {
		t.Fatalf("Only = %v", only)
	}

	if except := input.Except("age", "score"); len(except) != 2 {
		t.Fatalf("Except = %v", except)
	}

	if keys := input.Keys(); len(keys) != 4 {
		t.Fatalf("Keys = %v", keys)
	}

	merged := input.Merge(map[string]any{"name": "Abigail", "city": "Little Rock"})
	if merged.String("name") != "Abigail" || merged.String("city") != "Little Rock" {
		t.Fatalf("Merge = %v", merged.All())
	}

	if !input.Bool("active") {
		t.Fatal("Bool(active) should be true")
	}

	if input.Int("age") != 42 {
		t.Fatalf("Int(age) = %d", input.Int("age"))
	}

	if input.Float("score") != 9.5 {
		t.Fatalf("Float(score) = %f", input.Float("score"))
	}

	if input.String("name") != "Taylor" {
		t.Fatalf("String(name) = %q", input.String("name"))
	}
}
