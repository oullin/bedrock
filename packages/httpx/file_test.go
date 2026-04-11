package httpx_test

import (
	"strings"
	"testing"

	"github.com/bedrock/packages/httpx"
)

func TestFilePath(t *testing.T) {
	t.Parallel()

	f := httpx.NewFile("/tmp/test.txt")

	if f.Path() != "/tmp/test.txt" {
		t.Fatalf("expected /tmp/test.txt, got %s", f.Path())
	}
}

func TestFileBasename(t *testing.T) {
	t.Parallel()

	f := httpx.NewFile("/tmp/dir/photo.jpg")

	if f.Basename() != "photo.jpg" {
		t.Fatalf("expected photo.jpg, got %s", f.Basename())
	}
}

func TestFileExtension(t *testing.T) {
	t.Parallel()

	f := httpx.NewFile("/tmp/document.pdf")

	if f.Extension() != "pdf" {
		t.Fatalf("expected pdf, got %s", f.Extension())
	}
}

func TestFileExtensionEmpty(t *testing.T) {
	t.Parallel()

	f := httpx.NewFile("/tmp/Makefile")

	if f.Extension() != "" {
		t.Fatalf("expected empty extension, got %s", f.Extension())
	}
}

func TestFileHashName(t *testing.T) {
	t.Parallel()

	f := httpx.NewFile("/tmp/photo.jpg")
	hash := f.HashName()

	if !strings.HasSuffix(hash, ".jpg") {
		t.Fatalf("expected .jpg suffix, got %s", hash)
	}

	if len(hash) < 10 {
		t.Fatal("hash name seems too short")
	}

	// Two calls should produce different names.
	hash2 := f.HashName()

	if hash == hash2 {
		t.Fatal("expected different hash names on subsequent calls")
	}
}

func TestFileHashNameWithDirectory(t *testing.T) {
	t.Parallel()

	f := httpx.NewFile("/tmp/photo.jpg")
	hash := f.HashName("uploads")

	if !strings.HasPrefix(hash, "uploads/") {
		t.Fatalf("expected uploads/ prefix, got %s", hash)
	}

	if !strings.HasSuffix(hash, ".jpg") {
		t.Fatalf("expected .jpg suffix, got %s", hash)
	}
}
