package jsonx_test

import (
	"testing"

	"github.com/bedrock/packages/jsonx"
)

func TestNumberMinFloat(t *testing.T) {
	t.Parallel()

	schema := jsonx.Number().Title("Price").Min(5.5)
	result := schema.ToMap()

	assertEqual(t, result["type"], "number")
	assertEqual(t, result["title"], "Price")
	assertEqual(t, result["minimum"], 5.5)
}

func TestNumberMinInt(t *testing.T) {
	t.Parallel()

	schema := jsonx.Number().Min(5)
	result := schema.ToMap()

	assertEqual(t, result["type"], "number")
	assertEqual(t, result["minimum"], 5.0)
}

func TestNumberMaxFloat(t *testing.T) {
	t.Parallel()

	schema := jsonx.Number().Description("Max price").Max(99.9)
	result := schema.ToMap()

	assertEqual(t, result["type"], "number")
	assertEqual(t, result["description"], "Max price")
	assertEqual(t, result["maximum"], 99.9)
}

func TestNumberMaxInt(t *testing.T) {
	t.Parallel()

	schema := jsonx.Number().Max(100)
	result := schema.ToMap()

	assertEqual(t, result["type"], "number")
	assertEqual(t, result["maximum"], 100.0)
}

func TestNumberDefaultValue(t *testing.T) {
	t.Parallel()

	schema := jsonx.Number().Default(9.99)
	result := schema.ToMap()

	assertEqual(t, result["type"], "number")
	assertEqual(t, result["default"], 9.99)
}

func TestNumberMultipleOfFloat(t *testing.T) {
	t.Parallel()

	schema := jsonx.Number().MultipleOf(0.5)
	result := schema.ToMap()

	assertEqual(t, result["type"], "number")
	assertEqual(t, result["multipleOf"], 0.5)
}

func TestNumberMultipleOfInt(t *testing.T) {
	t.Parallel()

	schema := jsonx.Number().MultipleOf(3)
	result := schema.ToMap()

	assertEqual(t, result["type"], "number")
	assertEqual(t, result["multipleOf"], 3.0)
}

func TestNumberCombineConstraints(t *testing.T) {
	t.Parallel()

	schema := jsonx.Number().Min(0).Max(100).MultipleOf(0.5)
	result := schema.ToMap()

	assertEqual(t, result["type"], "number")
	assertEqual(t, result["minimum"], 0.0)
	assertEqual(t, result["maximum"], 100.0)
	assertEqual(t, result["multipleOf"], 0.5)
}

func TestNumberEnum(t *testing.T) {
	t.Parallel()

	schema := jsonx.Number().Enum([]any{1, 2.5, 3})
	result := schema.ToMap()

	assertEqual(t, result["type"], "number")
	assertSliceEqual(t, result["enum"], []any{1, 2.5, 3})
}
