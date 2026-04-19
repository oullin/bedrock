package precognition_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/httppreview"
)

func TestMarkPrecognitive(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r = httppreview.MarkPrecognitive(r)

	if !httppreview.IsPrecognitive(r) {
		t.Fatal("expected IsPrecognitive to return true after MarkPrecognitive")
	}
}

func TestIsPrecognitiveFromContext(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r = httppreview.MarkPrecognitive(r)

	if !httppreview.IsPrecognitive(r) {
		t.Fatal("expected IsPrecognitive to return true")
	}
}

func TestIsPrecognitiveWithoutContext(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodPost, "/", nil)

	if httppreview.IsPrecognitive(r) {
		t.Fatal("expected IsPrecognitive to return false without context")
	}
}

func TestIsPrecognitiveIgnoresHeader(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("HTTPPreview", "true")

	if httppreview.IsPrecognitive(r) {
		t.Fatal("IsPrecognitive should not check the header, only the context")
	}
}

func TestIsAttemptingHTTPPreview(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodPost, "/", nil)
	r.Header.Set("HTTPPreview", "true")

	if !httppreview.IsAttemptingHTTPPreview(r) {
		t.Fatal("expected IsAttemptingHTTPPreview to return true")
	}
}

func TestIsAttemptingHTTPPreviewWithoutHeader(t *testing.T) {
	t.Parallel()

	r := httptest.NewRequest(http.MethodPost, "/", nil)

	if httppreview.IsAttemptingHTTPPreview(r) {
		t.Fatal("expected IsAttemptingHTTPPreview to return false without header")
	}
}

func TestIsAttemptingHTTPPreviewRequiresExactTrue(t *testing.T) {
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
			r.Header.Set("HTTPPreview", tt.value)
		}

		got := httppreview.IsAttemptingHTTPPreview(r)

		if got != tt.want {
			t.Errorf("IsAttemptingHTTPPreview(%q) = %v, want %v", tt.value, got, tt.want)
		}
	}
}
