package testing_test

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"testing"

	httptesting "github.com/bedrock/packages/anvil/http/testing"
)

// Laravel: testImagePng
func TestFileFactoryImagePng(t *testing.T) {
	t.Parallel()

	factory := httptesting.NewFileFactory()
	uf, err := factory.Image("test.png", 10, 10)
	if err != nil {
		t.Fatalf("Image: %v", err)
	}

	if got := uf.MimeType(); got != "image/png" {
		t.Fatalf("expected 'image/png', got %q", got)
	}

	data, err := uf.Content()
	if err != nil {
		t.Fatalf("Content: %v", err)
	}

	assertImageDimensions(t, data, 10, 10)
}

// Laravel: testImageJpeg
func TestFileFactoryImageJpeg(t *testing.T) {
	t.Parallel()

	factory := httptesting.NewFileFactory()
	uf, err := factory.Image("photo.jpg", 16, 16)
	if err != nil {
		t.Fatalf("Image: %v", err)
	}

	if got := uf.MimeType(); got != "image/jpeg" {
		t.Fatalf("expected 'image/jpeg', got %q", got)
	}

	data, err := uf.Content()
	if err != nil {
		t.Fatalf("Content: %v", err)
	}

	assertImageDimensions(t, data, 16, 16)
}

// Laravel: testImageGif
func TestFileFactoryImageGif(t *testing.T) {
	t.Parallel()

	factory := httptesting.NewFileFactory()
	uf, err := factory.Image("anim.gif", 8, 8)
	if err != nil {
		t.Fatalf("Image: %v", err)
	}

	if got := uf.MimeType(); got != "image/gif" {
		t.Fatalf("expected 'image/gif', got %q", got)
	}

	data, err := uf.Content()
	if err != nil {
		t.Fatalf("Content: %v", err)
	}

	assertImageDimensions(t, data, 8, 8)
}

// Laravel: testCreateWithMimeType
func TestFileFactoryCreateWithMimeType(t *testing.T) {
	t.Parallel()

	factory := httptesting.NewFileFactory()
	uf, err := factory.Create("document.dat", 100, "application/pdf")
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if got := uf.MimeType(); got != "application/pdf" {
		t.Fatalf("expected 'application/pdf', got %q", got)
	}

	if got := uf.OriginalName(); got != "document.dat" {
		t.Fatalf("expected 'document.dat', got %q", got)
	}

	data, err := uf.Content()
	if err != nil {
		t.Fatalf("Content: %v", err)
	}

	// 100 KB.
	if len(data) != 100*1024 {
		t.Fatalf("expected %d bytes, got %d", 100*1024, len(data))
	}
}

// Laravel: testCreateWithoutMimeType
func TestFileFactoryCreateWithoutMimeType(t *testing.T) {
	t.Parallel()

	factory := httptesting.NewFileFactory()
	uf, err := factory.Create("data.json", 1)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}

	if got := uf.MimeType(); got != "application/json" {
		t.Fatalf("expected 'application/json', got %q", got)
	}
}

// Laravel: testImageWebp — adapted to test unsupported format error
func TestFileFactoryUnsupportedImageFormat(t *testing.T) {
	t.Parallel()

	factory := httptesting.NewFileFactory()
	_, err := factory.Image("test.webp", 10, 10)
	if err == nil {
		t.Fatal("expected error for unsupported image format 'webp'")
	}
}

// Laravel: testImageBmp — adapted to test unsupported format error
func TestFileFactoryUnsupportedBmpFormat(t *testing.T) {
	t.Parallel()

	factory := httptesting.NewFileFactory()
	_, err := factory.Image("test.bmp", 10, 10)
	if err == nil {
		t.Fatal("expected error for unsupported image format 'bmp'")
	}
}

func assertImageDimensions(t *testing.T, data []byte, expectedWidth, expectedHeight int) {
	t.Helper()

	cfg, _, err := image.DecodeConfig(byteReader(data))
	if err != nil {
		t.Fatalf("DecodeConfig: %v", err)
	}

	if cfg.Width != expectedWidth || cfg.Height != expectedHeight {
		t.Fatalf("expected %dx%d, got %dx%d", expectedWidth, expectedHeight, cfg.Width, cfg.Height)
	}
}

type byteReaderWrapper struct {
	data []byte
	pos  int
}

func byteReader(data []byte) *byteReaderWrapper {
	return &byteReaderWrapper{data: data}
}

func (r *byteReaderWrapper) Read(p []byte) (int, error) {
	if r.pos >= len(r.data) {
		return 0, nil
	}
	n := copy(p, r.data[r.pos:])
	r.pos += n
	if r.pos >= len(r.data) {
		return n, nil
	}
	return n, nil
}
