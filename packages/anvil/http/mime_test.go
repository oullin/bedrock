package http_test

import (
	"testing"

	bedhttp "github.com/bedrock/packages/anvil/http"
)

// Laravel: testMimeTypeFromFileNameExistsTrue
func TestMimeFromFileNameExists(t *testing.T) {
	t.Parallel()

	if got := bedhttp.MimeFrom("photo.jpg"); got != "image/jpeg" {
		t.Fatalf("expected 'image/jpeg', got %q", got)
	}
	if got := bedhttp.MimeFrom("document.pdf"); got != "application/pdf" {
		t.Fatalf("expected 'application/pdf', got %q", got)
	}
	if got := bedhttp.MimeFrom("style.css"); got != "text/css" {
		t.Fatalf("expected 'text/css', got %q", got)
	}
}

// Laravel: testMimeTypeFromFileNameExistsFalse
func TestMimeFromFileNameUnknown(t *testing.T) {
	t.Parallel()

	if got := bedhttp.MimeFrom("file.xyz123"); got != "application/octet-stream" {
		t.Fatalf("expected 'application/octet-stream', got %q", got)
	}
	if got := bedhttp.MimeFrom("noextension"); got != "application/octet-stream" {
		t.Fatalf("expected 'application/octet-stream' for no extension, got %q", got)
	}
}

// Laravel: testMimeTypeFromExtensionExistsTrue
func TestMimeGetExtensionExists(t *testing.T) {
	t.Parallel()

	if got := bedhttp.MimeGet("jpg"); got != "image/jpeg" {
		t.Fatalf("expected 'image/jpeg', got %q", got)
	}
	if got := bedhttp.MimeGet("png"); got != "image/png" {
		t.Fatalf("expected 'image/png', got %q", got)
	}
	if got := bedhttp.MimeGet("json"); got != "application/json" {
		t.Fatalf("expected 'application/json', got %q", got)
	}
}

// Laravel: testMimeTypeFromExtensionExistsFalse
func TestMimeGetExtensionUnknown(t *testing.T) {
	t.Parallel()

	if got := bedhttp.MimeGet("xyz123"); got != "application/octet-stream" {
		t.Fatalf("expected 'application/octet-stream', got %q", got)
	}
}

// Laravel: testSearchExtensionFromMimeType
func TestMimeSearchFromMimeType(t *testing.T) {
	t.Parallel()

	got := bedhttp.MimeSearch("image/jpeg")
	if got != "jpg" && got != "jpeg" {
		t.Fatalf("expected 'jpg' or 'jpeg', got %q", got)
	}

	if got := bedhttp.MimeSearch("application/pdf"); got != "pdf" {
		t.Fatalf("expected 'pdf', got %q", got)
	}

	if got := bedhttp.MimeSearch("application/totally-made-up"); got != "" {
		t.Fatalf("expected empty string for unknown MIME, got %q", got)
	}
}
