package testing_test

import (
	"testing"

	httptesting "github.com/bedrock/packages/httpx/testing"
)

func TestFakeFile(t *testing.T) {
	t.Parallel()

	file := httptesting.FakeFile("document.pdf", 1)

	if file == nil {
		t.Fatal("expected non-nil file")
	}

	if file.ClientOriginalName() != "document.pdf" {
		t.Fatalf("expected document.pdf, got %s", file.ClientOriginalName())
	}

	if !file.IsValid() {
		t.Fatal("expected fake file to be valid")
	}

	data, err := file.Get()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("expected non-empty content")
	}
}

func TestFakeImage(t *testing.T) {
	t.Parallel()

	file := httptesting.FakeImage("photo.png", 100, 100)

	if file == nil {
		t.Fatal("expected non-nil file")
	}

	if file.ClientOriginalName() != "photo.png" {
		t.Fatalf("expected photo.png, got %s", file.ClientOriginalName())
	}

	data, err := file.Get()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// PNG magic bytes.
	if len(data) < 8 || data[0] != 0x89 || data[1] != 'P' {
		t.Fatal("expected PNG content")
	}
}

func TestFakeFileWithContent(t *testing.T) {
	t.Parallel()

	content := []byte("custom content here")
	file := httptesting.FakeFileWithContent("test.txt", content)

	data, err := file.Get()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if string(data) != "custom content here" {
		t.Fatalf("expected 'custom content here', got %s", string(data))
	}
}
