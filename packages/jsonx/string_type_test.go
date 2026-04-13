package jsonx_test

import (
	"testing"

	"github.com/bedrock/packages/jsonx"
)

func TestStringMinLength(t *testing.T) {
	t.Parallel()

	schema := jsonx.String().Min(5)
	result := schema.ToMap()

	assertEqual(t, result["type"], "string")
	assertEqual(t, result["minLength"], 5)
}

func TestStringMaxLength(t *testing.T) {
	t.Parallel()

	schema := jsonx.String().Description("User handle").Max(10)
	result := schema.ToMap()

	assertEqual(t, result["type"], "string")
	assertEqual(t, result["description"], "User handle")
	assertEqual(t, result["maxLength"], 10)
}

func TestStringPattern(t *testing.T) {
	t.Parallel()

	schema := jsonx.String().Default("foo").Pattern("^foo.*$")
	result := schema.ToMap()

	assertEqual(t, result["type"], "string")
	assertEqual(t, result["default"], "foo")
	assertEqual(t, result["pattern"], "^foo.*$")
}

func TestStringFormat(t *testing.T) {
	t.Parallel()

	schema := jsonx.String().Default("foo").Format("date")
	result := schema.ToMap()

	assertEqual(t, result["type"], "string")
	assertEqual(t, result["default"], "foo")
	assertEqual(t, result["format"], "date")
}

func TestStringEnum(t *testing.T) {
	t.Parallel()

	schema := jsonx.String().Enum([]any{"draft", "published"})
	result := schema.ToMap()

	assertEqual(t, result["type"], "string")
	assertSliceEqual(t, result["enum"], []any{"draft", "published"})
}
