package validation_test

import (
	"errors"
	"testing"

	"github.com/bedrock/packages/validation"
)

// ValidationExceptionTest::testGetExceptionClassFromValidator
func TestValidationExceptionFromValidator(t *testing.T) {
	t.Parallel()

	v := makeValidator(
		map[string]any{"name": ""},
		map[string]any{"name": "required"},
	)

	err := v.Validate()

	if err == nil {
		t.Fatal("expected validation error")
	}

	if !errors.Is(err, validation.ErrValidationFailed) {
		t.Fatalf("errors.Is(err, ErrValidationFailed) = false, want true")
	}

	var ve *validation.ValidationException

	if !errors.As(err, &ve) {
		t.Fatalf("errors.As(err, *ValidationException) = false, got %T", err)
	}

	if ve.Bag == nil || ve.Bag.IsEmpty() {
		t.Fatal("expected validation exception to carry error messages")
	}
}
