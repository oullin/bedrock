package precognition_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/precognition"
)

func TestWriteSuccessResponse(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	precognition.WriteSuccessResponse(rec)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", rec.Code)
	}

	if got := rec.Header().Get("Precognition-Success"); got != "true" {
		t.Fatalf("expected Precognition-Success: true, got %q", got)
	}
}

func TestAddVaryHeader(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	precognition.AddVaryHeader(rec)

	if got := rec.Header().Get("Vary"); got != "Precognition" {
		t.Fatalf("expected Vary: Precognition, got %q", got)
	}
}

func TestAddVaryHeaderExisting(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	rec.Header().Set("Vary", "Accept-Encoding")
	precognition.AddVaryHeader(rec)

	values := rec.Header().Values("Vary")

	if len(values) != 2 {
		t.Fatalf("expected 2 Vary values, got %d: %v", len(values), values)
	}

	if values[0] != "Accept-Encoding" {
		t.Fatalf("expected first Vary to be Accept-Encoding, got %q", values[0])
	}

	if values[1] != "Precognition" {
		t.Fatalf("expected second Vary to be Precognition, got %q", values[1])
	}
}

func TestAddPrecognitionHeader(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	precognition.AddPrecognitionHeader(rec)

	if got := rec.Header().Get("Precognition"); got != "true" {
		t.Fatalf("expected Precognition: true, got %q", got)
	}
}
