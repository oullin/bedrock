package http_test

import (
	"bytes"
	"mime/multipart"
	"testing"

	bedhttp "github.com/bedrock/packages/anvil/http"
)

// Upstream: testUploadedFileCanRetrieveContentsFromTextFile
func TestUploadedFileContent(t *testing.T) {
	t.Parallel()

	content := "Hello, uploaded world!"
	header := createFileHeader(t, "hello.txt", "text/plain", []byte(content))

	uf := bedhttp.NewUploadedFile(header)

	got, err := uf.Content()
	if err != nil {
		t.Fatalf("Content: %v", err)
	}
	if string(got) != content {
		t.Fatalf("expected %q, got %q", content, string(got))
	}
}

// Upstream: testUploadedFileInRequestContainsOriginalPathAndName
func TestUploadedFileOriginalName(t *testing.T) {
	t.Parallel()

	header := createFileHeader(t, "document.pdf", "application/pdf", []byte("fake pdf"))

	uf := bedhttp.NewUploadedFile(header)

	if got := uf.OriginalName(); got != "document.pdf" {
		t.Fatalf("expected 'document.pdf', got %q", got)
	}
	if got := uf.Name(); got != "document" {
		t.Fatalf("expected 'document', got %q", got)
	}
	if got := uf.Extension(); got != "pdf" {
		t.Fatalf("expected 'pdf', got %q", got)
	}
	if got := uf.MimeType(); got != "application/pdf" {
		t.Fatalf("expected 'application/pdf', got %q", got)
	}
}

// createFileHeader builds a multipart.FileHeader for testing.
func createFileHeader(t *testing.T, filename, contentType string, data []byte) *multipart.FileHeader {
	t.Helper()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}

	if _, err := part.Write(data); err != nil {
		t.Fatalf("Write: %v", err)
	}

	if err := writer.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	reader := multipart.NewReader(&buf, writer.Boundary())
	form, err := reader.ReadForm(1 << 20)
	if err != nil {
		t.Fatalf("ReadForm: %v", err)
	}

	files := form.File["file"]
	if len(files) == 0 {
		t.Fatal("expected at least one file")
	}

	header := files[0]
	// Set explicit Content-Type if the multipart writer didn't.
	if contentType != "" {
		header.Header.Set("Content-Type", contentType)
	}

	return header
}
