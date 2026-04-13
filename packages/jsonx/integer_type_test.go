package jsonx_test

import (
	"testing"

	"github.com/bedrock/packages/jsonx"
)

func TestIntegerMinValue(t *testing.T) {
	t.Parallel()

	schema := jsonx.Integer().Title("Age").Min(5)
	result := schema.ToMap()

	assertEqual(t, result["type"], "integer")
	assertEqual(t, result["title"], "Age")
	assertEqual(t, result["minimum"], 5)
}

func TestIntegerMaxValue(t *testing.T) {
	t.Parallel()

	schema := jsonx.Integer().Description("Max age").Max(10)
	result := schema.ToMap()

	assertEqual(t, result["type"], "integer")
	assertEqual(t, result["description"], "Max age")
	assertEqual(t, result["maximum"], 10)
}

func TestIntegerDefaultValue(t *testing.T) {
	t.Parallel()

	schema := jsonx.Integer().Default(18)
	result := schema.ToMap()

	assertEqual(t, result["type"], "integer")
	assertEqual(t, result["default"], 18)
}

func TestIntegerMultipleOf(t *testing.T) {
	t.Parallel()

	schema := jsonx.Integer().MultipleOf(5)
	result := schema.ToMap()

	assertEqual(t, result["type"], "integer")
	assertEqual(t, result["multipleOf"], 5)
}

func TestIntegerCombineMultipleOfWithMinAndMax(t *testing.T) {
	t.Parallel()

	schema := jsonx.Integer().Min(0).Max(100).MultipleOf(10)
	result := schema.ToMap()

	assertEqual(t, result["type"], "integer")
	assertEqual(t, result["minimum"], 0)
	assertEqual(t, result["maximum"], 100)
	assertEqual(t, result["multipleOf"], 10)
}

func TestIntegerEnum(t *testing.T) {
	t.Parallel()

	schema := jsonx.Integer().Enum([]any{1, 2, 3})
	result := schema.ToMap()

	assertEqual(t, result["type"], "integer")
	assertSliceEqual(t, result["enum"], []any{1, 2, 3})
}
