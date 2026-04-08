// Package testing provides test utilities for the HTTP package.
// It mirrors Laravel's Illuminate\Http\Testing namespace.
package testing

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	bedhttp "github.com/bedrock/packages/anvil/http"
)

// FileFactory creates fake uploaded files for testing purposes.
// It mirrors Laravel's Illuminate\Http\Testing\FileFactory.
type FileFactory struct{}

// NewFileFactory creates a new FileFactory.
func NewFileFactory() *FileFactory {
	return &FileFactory{}
}

// Image creates a fake uploaded image file with the given dimensions.
// Supported formats: png, jpg/jpeg, gif. The format is inferred from the
// file extension.
func (f *FileFactory) Image(name string, width, height int) (*bedhttp.UploadedFile, error) {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(name), "."))
	if ext == "" {
		return nil, fmt.Errorf("file_factory: image name %q has no extension", name)
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Fill with a solid colour so the image is not empty.
	for y := range height {
		for x := range width {
			img.Set(x, y, color.RGBA{R: 0, G: 128, B: 255, A: 255})
		}
	}

	var buf bytes.Buffer
	var mimeType string

	switch ext {
	case "png":
		if err := png.Encode(&buf, img); err != nil {
			return nil, fmt.Errorf("file_factory: png encode: %w", err)
		}
		mimeType = "image/png"
	case "jpg", "jpeg":
		if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}); err != nil {
			return nil, fmt.Errorf("file_factory: jpeg encode: %w", err)
		}
		mimeType = "image/jpeg"
	case "gif":
		if err := gif.Encode(&buf, img, nil); err != nil {
			return nil, fmt.Errorf("file_factory: gif encode: %w", err)
		}
		mimeType = "image/gif"
	default:
		return nil, fmt.Errorf("file_factory: unsupported image format %q", ext)
	}

	return f.createUploadedFile(name, mimeType, buf.Bytes())
}

// Create creates a fake uploaded file with the given size in kilobytes.
// If no MIME type is provided, it is inferred from the file extension.
func (f *FileFactory) Create(name string, sizeKB int, mimeType ...string) (*bedhttp.UploadedFile, error) {
	mt := ""
	if len(mimeType) > 0 {
		mt = mimeType[0]
	}
	if mt == "" {
		mt = bedhttp.MimeFrom(name)
	}

	data := make([]byte, sizeKB*1024)

	return f.createUploadedFile(name, mt, data)
}

func (f *FileFactory) createUploadedFile(name, mimeType string, data []byte) (*bedhttp.UploadedFile, error) {
	// Write to a temp file so the multipart.FileHeader has real content.
	tmp, err := os.CreateTemp("", "bedrock-test-*-"+name)
	if err != nil {
		return nil, fmt.Errorf("file_factory: temp file: %w", err)
	}

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return nil, fmt.Errorf("file_factory: write: %w", err)
	}
	tmp.Close()

	// Build a multipart form with the file.
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", name)
	if err != nil {
		os.Remove(tmp.Name())
		return nil, fmt.Errorf("file_factory: create form file: %w", err)
	}

	if _, err := part.Write(data); err != nil {
		os.Remove(tmp.Name())
		return nil, fmt.Errorf("file_factory: write part: %w", err)
	}
	writer.Close()

	reader := multipart.NewReader(&buf, writer.Boundary())
	form, err := reader.ReadForm(int64(len(data)) + 1<<20)
	if err != nil {
		os.Remove(tmp.Name())
		return nil, fmt.Errorf("file_factory: read form: %w", err)
	}

	files := form.File["file"]
	if len(files) == 0 {
		os.Remove(tmp.Name())
		return nil, fmt.Errorf("file_factory: no file in form")
	}

	header := files[0]
	header.Header.Set("Content-Type", mimeType)

	os.Remove(tmp.Name())

	return bedhttp.NewUploadedFile(header), nil
}
