package validation_test

import (
	"testing"

	"github.com/bedrock/packages/validation"
)

// ValidationFactoryTest::testMakeMethodCreatesValidValidator
func TestFactoryMakeCreatesValidator(t *testing.T) {
	t.Parallel()

	f := validation.NewFactory()
	v := f.Make(
		map[string]any{"email": "user@example.com"},
		map[string]any{"email": "required|email"},
		nil,
		nil,
	)

	assertPasses(t, v)
}

// ValidationFactoryTest::testValidateMethodCanBeCalledPublicly
func TestFactoryValidateCanBeCalledPublicly(t *testing.T) {
	t.Parallel()

	f := validation.NewFactory()

	got, err := f.Validate(
		map[string]any{"name": "Taylor"},
		map[string]any{"name": "required|string"},
		nil,
		nil,
	)

	if err != nil {
		t.Fatalf("Validate() returned error: %v", err)
	}

	if got["name"] != "Taylor" {
		t.Fatalf("Validate() = %v, want name to be preserved", got)
	}
}
