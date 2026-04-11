package testing_test

import (
	"testing"

	httptesting "github.com/bedrock/packages/httpx/testing"
)

func TestMimeTypeByExtension(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    string
		expected string
	}{
		{"jpg", "image/jpeg"},
		{"png", "image/png"},
		{"pdf", "application/pdf"},
		{"json", "application/json"},
		{"txt", "text/plain"},
		{"html", "text/html"},
	}

	for _, tt := range tests {
		if got := httptesting.MimeType(tt.input); got != tt.expected {
			t.Errorf("MimeType(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestMimeTypeByFilename(t *testing.T) {
	t.Parallel()

	if got := httptesting.MimeType("photo.jpg"); got != "image/jpeg" {
		t.Fatalf("expected image/jpeg, got %s", got)
	}
}

func TestMimeTypeUnknown(t *testing.T) {
	t.Parallel()

	if got := httptesting.MimeType("xyz"); got != "application/octet-stream" {
		t.Fatalf("expected application/octet-stream, got %s", got)
	}
}

func TestExtensionForMime(t *testing.T) {
	t.Parallel()

	ext := httptesting.ExtensionForMime("image/png")

	if ext != "png" {
		t.Fatalf("expected png, got %s", ext)
	}
}

func TestExtensionForMimeUnknown(t *testing.T) {
	t.Parallel()

	ext := httptesting.ExtensionForMime("application/x-unknown")

	if ext != "" {
		t.Fatalf("expected empty, got %s", ext)
	}
}
