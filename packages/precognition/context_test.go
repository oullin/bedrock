package precognition_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/precognition"
)

func TestMarkPrecognitive(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r = precognition.MarkPrecognitive(r)

	if !precognition.IsPrecognitive(r) {
		t.Fatal("expected IsPrecognitive to return true after MarkPrecognitive")
	}
}

func TestIsPrecognitiveFromContext(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r = precognition.MarkPrecognitive(r)

	if !precognition.IsPrecognitive(r) {
		t.Fatal("expected IsPrecognitive to return true")
	}
}

func TestIsPrecognitiveWithoutContext(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodPost, "/", nil)

	if precognition.IsPrecognitive(r) {
		t.Fatal("expected IsPrecognitive to return false without context")
	}
}

func TestIsPrecognitiveIgnoresHeader(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("Precognition", "true")

	if precognition.IsPrecognitive(r) {
		t.Fatal("IsPrecognitive should not check the header, only the context")
	}
}

func TestIsAttemptingPrecognition(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("Precognition", "true")

	if !precognition.IsAttemptingPrecognition(r) {
		t.Fatal("expected IsAttemptingPrecognition to return true")
	}
}

func TestIsAttemptingPrecognitionWithoutHeader(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodPost, "/", nil)

	if precognition.IsAttemptingPrecognition(r) {
		t.Fatal("expected IsAttemptingPrecognition to return false without header")
	}
}

func TestIsAttemptingPrecognitionRequiresExactTrue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value string
		want  bool
	}{
		{"true", true},
		{"1", false},
		{"yes", false},
		{"on", false},
		{"True", false},
		{"TRUE", false},
		{"false", false},
		{"", false},
	}

	for _, tt := range tests {
		r := httptest.NewRequest(http.MethodPost, "/", nil)

		if tt.value != "" {
			r.Header.Set("Precognition", tt.value)
		}

		got := precognition.IsAttemptingPrecognition(r)

		if got != tt.want {
			t.Errorf("IsAttemptingPrecognition(%q) = %v, want %v", tt.value, got, tt.want)
		}
	}
}
