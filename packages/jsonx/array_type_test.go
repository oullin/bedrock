package jsonx_test

import (
	"testing"

	"github.com/bedrock/packages/jsonx"
)

func TestArrayMinItems(t *testing.T) {
	t.Parallel()

	schema := jsonx.Array().Title("Tags").Min(1)
	result := schema.ToMap()

	assertEqual(t, result["type"], "array")
	assertEqual(t, result["title"], "Tags")
	assertEqual(t, result["minItems"], 1)
}

func TestArrayMaxItems(t *testing.T) {
	t.Parallel()

	schema := jsonx.Array().Description("A list of tags").Max(10)
	result := schema.ToMap()

	assertEqual(t, result["type"], "array")
	assertEqual(t, result["description"], "A list of tags")
	assertEqual(t, result["maxItems"], 10)
}

func TestArrayItemsType(t *testing.T) {
	t.Parallel()

	schema := jsonx.Array().Items(jsonx.String().Max(20))
	result := schema.ToMap()

	assertEqual(t, result["type"], "array")

	items, ok := result["items"].(map[string]any)

	if !ok {
		t.Fatal("expected items to be map[string]any")
	}

	assertEqual(t, items["type"], "string")
	assertEqual(t, items["maxLength"], 20)
}

func TestArrayDefaultValue(t *testing.T) {
	t.Parallel()

	schema := jsonx.Array().Default([]any{"a", "b"})
	result := schema.ToMap()

	assertEqual(t, result["type"], "array")
	assertSliceEqual(t, result["default"], []any{"a", "b"})
}

func TestArrayUniqueItems(t *testing.T) {
	t.Parallel()

	schema := jsonx.Array().Items(jsonx.String()).Unique()
	result := schema.ToMap()

	assertEqual(t, result["type"], "array")
	assertEqual(t, result["uniqueItems"], true)

	items, ok := result["items"].(map[string]any)

	if !ok {
		t.Fatal("expected items to be map[string]any")
	}

	assertEqual(t, items["type"], "string")
}

func TestArrayCombineUniqueWithMinAndMax(t *testing.T) {
	t.Parallel()

	schema := jsonx.Array().Min(1).Max(5).Unique()
	result := schema.ToMap()

	assertEqual(t, result["type"], "array")
	assertEqual(t, result["minItems"], 1)
	assertEqual(t, result["maxItems"], 5)
	assertEqual(t, result["uniqueItems"], true)
}

func TestArrayEnum(t *testing.T) {
	t.Parallel()

	schema := jsonx.Array().Enum([]any{
		[]any{"a"},
		[]any{"b", "c"},
	})
	result := schema.ToMap()

	assertEqual(t, result["type"], "array")

	enum, ok := result["enum"].([]any)

	if !ok {
		t.Fatal("expected enum to be []any")
	}

	if len(enum) != 2 {
		t.Fatalf("expected 2 enum values, got %d", len(enum))
	}
}
