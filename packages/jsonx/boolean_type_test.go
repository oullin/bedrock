package jsonx_test

import (
	"testing"

	"github.com/bedrock/packages/jsonx"
)

func TestBooleanSerializesWithMetadata(t *testing.T) {
	t.Parallel()

	schema := jsonx.Boolean().Title("Enabled").Description("Feature flag")
	result := schema.ToMap()

	assertEqual(t, result["type"], "boolean")
	assertEqual(t, result["title"], "Enabled")
	assertEqual(t, result["description"], "Feature flag")
}

func TestBooleanDefaultTrue(t *testing.T) {
	t.Parallel()

	schema := jsonx.Boolean().Default(true)
	result := schema.ToMap()

	assertEqual(t, result["type"], "boolean")
	assertEqual(t, result["default"], true)
}

func TestBooleanDefaultFalse(t *testing.T) {
	t.Parallel()

	schema := jsonx.Boolean().Default(false)
	result := schema.ToMap()

	assertEqual(t, result["type"], "boolean")
	assertEqual(t, result["default"], false)
}

func TestBooleanEnum(t *testing.T) {
	t.Parallel()

	schema := jsonx.Boolean().Enum([]any{true, false})
	result := schema.ToMap()

	assertEqual(t, result["type"], "boolean")
	assertSliceEqual(t, result["enum"], []any{true, false})
}
