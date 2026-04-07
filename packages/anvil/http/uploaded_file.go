package http

import (
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
)

// UploadedFile wraps a multipart file header providing convenient
// accessors for uploaded file metadata and content. It mirrors
// Upstream's Framework\Http\UploadedFile.
type UploadedFile struct {
	header *multipart.FileHeader
}

// NewUploadedFile creates an UploadedFile from a multipart file header.
func NewUploadedFile(header *multipart.FileHeader) *UploadedFile {
	return &UploadedFile{header: header}
}

// Name returns the base filename without the extension.
func (uf *UploadedFile) Name() string {
	name := uf.header.Filename
	ext := filepath.Ext(name)
	return strings.TrimSuffix(name, ext)
}

// OriginalName returns the original filename as provided by the client.
func (uf *UploadedFile) OriginalName() string {
	return uf.header.Filename
}

// Extension returns the file extension (without the leading dot).
func (uf *UploadedFile) Extension() string {
	ext := filepath.Ext(uf.header.Filename)
	if ext == "" {
		return ""
	}
	return ext[1:]
}

// MimeType returns the MIME type from the Content-Type header of the
// uploaded file part. Falls back to extension-based detection.
func (uf *UploadedFile) MimeType() string {
	ct := uf.header.Header.Get("Content-Type")
	if ct != "" {
		if i := strings.IndexByte(ct, ';'); i >= 0 {
			ct = strings.TrimSpace(ct[:i])
		}
		return ct
	}

	return MimeFrom(uf.header.Filename)
}

// Size returns the file size in bytes.
func (uf *UploadedFile) Size() int64 {
	return uf.header.Size
}

// Content reads and returns the entire file content.
func (uf *UploadedFile) Content() ([]byte, error) {
	f, err := uf.header.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return io.ReadAll(f)
}

// FileHeader returns the underlying multipart.FileHeader.
func (uf *UploadedFile) FileHeader() *multipart.FileHeader {
	return uf.header
}
