package httppreview_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/httppreview"
)

func TestWriteSuccessResponse(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	httppreview.WriteSuccessResponse(rec)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}

	if got := rec.Header().Get("HTTPPreview-Success"); got != "true" {
		t.Fatalf("expected HTTPPreview-Success: true, got %q", got)
	}
}

func TestAddVaryHeader(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	httppreview.AddVaryHeader(rec)

	if got := rec.Header().Get("Vary"); got != "HTTPPreview" {
		t.Fatalf("expected Vary: HTTPPreview, got %q", got)
	}
}

func TestAddVaryHeaderExisting(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	rec.Header().Set("Vary", "Accept-Encoding")
	httppreview.AddVaryHeader(rec)

	values := rec.Header().Values("Vary")

	if len(values) != 2 {
		t.Fatalf("expected 2 Vary values, got %d: %v", len(values), values)
	}

	if values[0] != "Accept-Encoding" {
		t.Fatalf("expected first Vary to be Accept-Encoding, got %q", values[0])
	}

	if values[1] != "HTTPPreview" {
		t.Fatalf("expected second Vary to be HTTPPreview, got %q", values[1])
	}
}

func TestAddHTTPPreviewHeader(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	httppreview.AddHTTPPreviewHeader(rec)

	if got := rec.Header().Get("HTTPPreview"); got != "true" {
		t.Fatalf("expected HTTPPreview: true, got %q", got)
	}
}
