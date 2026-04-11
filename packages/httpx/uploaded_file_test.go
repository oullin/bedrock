package httpx_test

import (
	"bytes"
	"encoding/base64"
	"io"
	"mime/multipart"
	"strings"
	"testing"

	"github.com/bedrock/packages/httpx"
)

func createTestUploadedFile(t *testing.T, fieldName, fileName, content string) *httpx.UploadedFile {
	t.Helper()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile(fieldName, fileName)

	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}

	part.Write([]byte(content))
	writer.Close()

	reader := multipart.NewReader(&buf, writer.Boundary())
	form, err := reader.ReadForm(32 << 20)

	if err != nil {
		t.Fatalf("failed to read form: %v", err)
	}

	files := form.File[fieldName]

	if len(files) == 0 {
		t.Fatal("no files found")
	}

	return httpx.NewUploadedFile(files[0])
}

func TestUploadedFileClientOriginalName(t *testing.T) {
	t.Parallel()

	file := createTestUploadedFile(t, "doc", "report.pdf", "pdf content")

	if file.ClientOriginalName() != "report.pdf" {
		t.Fatalf("expected report.pdf, got %s", file.ClientOriginalName())
	}
}

func TestUploadedFileClientExtension(t *testing.T) {
	t.Parallel()

	file := createTestUploadedFile(t, "doc", "report.pdf", "pdf content")

	if file.ClientExtension() != "pdf" {
		t.Fatalf("expected pdf, got %s", file.ClientExtension())
	}
}

func TestUploadedFileSize(t *testing.T) {
	t.Parallel()

	content := "hello world"
	file := createTestUploadedFile(t, "doc", "test.txt", content)

	if file.Size() != int64(len(content)) {
		t.Fatalf("expected size %d, got %d", len(content), file.Size())
	}
}

func TestUploadedFileIsValid(t *testing.T) {
	t.Parallel()

	file := createTestUploadedFile(t, "doc", "test.txt", "content")

	if !file.IsValid() {
		t.Fatal("expected file to be valid")
	}
}

func TestUploadedFileGet(t *testing.T) {
	t.Parallel()

	file := createTestUploadedFile(t, "doc", "test.txt", "file content here")

	data, err := file.Get()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != "file content here" {
		t.Fatalf("expected 'file content here', got %s", string(data))
	}
}

func TestUploadedFileOpen(t *testing.T) {
	t.Parallel()

	file := createTestUploadedFile(t, "doc", "test.txt", "readable")

	rc, err := file.Open()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	defer rc.Close()

	data, _ := io.ReadAll(rc)

	if string(data) != "readable" {
		t.Fatalf("expected 'readable', got %s", string(data))
	}
}

func TestUploadedFileHashName(t *testing.T) {
	t.Parallel()

	file := createTestUploadedFile(t, "doc", "photo.jpg", "image data")

	hash := file.HashName()

	if !strings.HasSuffix(hash, ".jpg") {
		t.Fatalf("expected .jpg suffix, got %s", hash)
	}

	hash2 := file.HashName("uploads")

	if !strings.HasPrefix(hash2, "uploads/") {
		t.Fatalf("expected uploads/ prefix, got %s", hash2)
	}
}

func TestUploadedFileStore(t *testing.T) {
	t.Parallel()

	file := createTestUploadedFile(t, "doc", "test.txt", "store me")

	store := &memoryFileStore{files: make(map[string][]byte)}

	path, err := file.Store("uploads", store)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.HasPrefix(path, "uploads/") {
		t.Fatalf("expected uploads/ prefix, got %s", path)
	}

	if len(store.files) != 1 {
		t.Fatalf("expected 1 file stored, got %d", len(store.files))
	}
}

func TestUploadedFileStoreAs(t *testing.T) {
	t.Parallel()

	file := createTestUploadedFile(t, "doc", "test.txt", "store me")

	store := &memoryFileStore{files: make(map[string][]byte)}

	path, err := file.StoreAs("uploads", "custom.txt", store)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if path != "uploads/custom.txt" {
		t.Fatalf("expected uploads/custom.txt, got %s", path)
	}

	if string(store.files["uploads/custom.txt"]) != "store me" {
		t.Fatal("expected stored content to match")
	}
}

func TestUploadedFileHeader(t *testing.T) {
	t.Parallel()

	file := createTestUploadedFile(t, "doc", "test.txt", "content")

	if file.Header() == nil {
		t.Fatal("expected non-nil multipart.FileHeader")
	}
}

func TestCreateFromBase64(t *testing.T) {
	t.Parallel()

	content := "Hello, World!"
	encoded := base64.StdEncoding.EncodeToString([]byte(content))

	file, err := httpx.CreateFromBase64(encoded, "test.txt")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if file.ClientOriginalName() != "test.txt" {
		t.Fatalf("expected test.txt, got %s", file.ClientOriginalName())
	}

	if file.Size() != int64(len(content)) {
		t.Fatalf("expected size %d, got %d", len(content), file.Size())
	}

	if !file.IsValid() {
		t.Fatal("expected base64 file to be valid")
	}

	data, err := file.Get()

	if err != nil {
		t.Fatalf("unexpected error reading base64 file: %v", err)
	}

	if string(data) != content {
		t.Fatalf("expected '%s', got '%s'", content, string(data))
	}
}

func TestCreateFromBase64Invalid(t *testing.T) {
	t.Parallel()

	_, err := httpx.CreateFromBase64("not-valid-base64!!!", "test.txt")

	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
}

// memoryFileStore is a test double for httpx.FileStore.
type memoryFileStore struct {
	files map[string][]byte
}

func (s *memoryFileStore) Put(path string, contents io.Reader) error {
	data, err := io.ReadAll(contents)

	if err != nil {
		return err
	}

	s.files[path] = data

	return nil
}
