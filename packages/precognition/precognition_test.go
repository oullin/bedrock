package precognition_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/precognition"
)

// emptyMessages simulates a validator with no errors.
type emptyMessages struct{}

func (emptyMessages) IsEmpty() bool { return true }

// nonEmptyMessages simulates a validator with errors.
type nonEmptyMessages struct{}

func (nonEmptyMessages) IsEmpty() bool { return false }

func TestAfterValidationHookPanicsWhenValidationPassesAndValidateOnlyPresent(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("Precognition", "true")
	r.Header.Set("Precognition-Validate-Only", "name,email")

	hook := precognition.AfterValidationHook(r)

	defer func() {
		v := recover()

		if v == nil {
			t.Fatal("expected panic from AfterValidationHook")
		}

		if _, ok := v.(precognition.SuccessResponse); !ok {
			t.Fatalf("expected SuccessResponse panic, got %T", v)
		}
	}()

	hook(emptyMessages{})
}

func TestAfterValidationHookNoOpWhenValidationFails(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("Precognition", "true")
	r.Header.Set("Precognition-Validate-Only", "name,email")

	hook := precognition.AfterValidationHook(r)

	// Should not panic when validation has errors.
	hook(nonEmptyMessages{})
}

func TestAfterValidationHookNoOpWithoutValidateOnly(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("Precognition", "true")

	hook := precognition.AfterValidationHook(r)

	// Should not panic when no Precognition-Validate-Only header.
	hook(emptyMessages{})
}

func TestAfterValidationHookNoOpForNonPrecognitive(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodPost, "/", nil)

	hook := precognition.AfterValidationHook(r)

	// Should not panic for non-precognitive requests.
	hook(emptyMessages{})
}
