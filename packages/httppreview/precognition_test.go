package precognition_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/httppreview"
)

// emptyMessages simulates a validator with no errors.
type emptyMessages struct{}

// nonEmptyMessages simulates a validator with errors.
type nonEmptyMessages struct{}

func (emptyMessages) IsEmpty() bool { return true }

func (nonEmptyMessages) IsEmpty() bool { return false }

func TestAfterValidationHookPanicsWhenValidationPassesAndValidateOnlyPresent(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("HTTPPreview", "true")
	r.Header.Set("HTTPPreview-Validate-Only", "name,email")

	hook := httppreview.AfterValidationHook(r)

	defer func() {
		v := recover()

		if v == nil {
			t.Fatal("expected panic from AfterValidationHook")
		}

		if _, ok := v.(httppreview.SuccessResponse); !ok {
			t.Fatalf("expected SuccessResponse panic, got %T", v)
		}
	}()

	hook(emptyMessages{})
}

func TestAfterValidationHookNoOpWhenValidationFails(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("HTTPPreview", "true")
	r.Header.Set("HTTPPreview-Validate-Only", "name,email")

	hook := httppreview.AfterValidationHook(r)

	// Should not panic when validation has errors.
	hook(nonEmptyMessages{})
}

func TestAfterValidationHookNoOpWithoutValidateOnly(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("HTTPPreview", "true")

	hook := httppreview.AfterValidationHook(r)

	// Should not panic when no HTTPPreview-Validate-Only header.
	hook(emptyMessages{})
}

func TestAfterValidationHookNoOpForNonPrecognitive(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodPost, "/", nil)

	hook := httppreview.AfterValidationHook(r)

	// Should not panic for non-precognitive requests.
	hook(emptyMessages{})
}
