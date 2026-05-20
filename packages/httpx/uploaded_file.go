package httpx

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

// FileStore abstracts file storage for uploaded files. Implementations may
// write to local disk, S3, GCS, or any other backend.
type FileStore interface {
	Put(path string, contents io.Reader) error
}

// UploadedFile represents a file received via a multipart upload.
type UploadedFile struct {
	header *multipart.FileHeader

	// test indicates this is a fake file created for testing.
	test bool
	// testContent holds fake file content for test files.
	testContent []byte
}

// NewUploadedFile wraps a multipart.FileHeader.
func NewUploadedFile(header *multipart.FileHeader) *UploadedFile {
	return &UploadedFile{header: header}
}

// ClientOriginalName returns the filename provided by the client.
func (f *UploadedFile) ClientOriginalName() string {
	return f.header.Filename
}

// ClientMimeType returns the MIME type reported by the client.
func (f *UploadedFile) ClientMimeType() string {
	ct := f.header.Header.Get("Content-Type")

	if ct == "" {
		ct = mime.TypeByExtension(filepath.Ext(f.header.Filename))
	}

	return ct
}

// ClientExtension returns the file extension based on the client-reported name.
func (f *UploadedFile) ClientExtension() string {
	return strings.TrimPrefix(filepath.Ext(f.header.Filename), ".")
}

// Size returns the file size in bytes.
func (f *UploadedFile) Size() int64 {
	return f.header.Size
}

// IsValid returns true when the file was uploaded without errors.
func (f *UploadedFile) IsValid() bool {
	if f.test {
		return true
	}

	file, err := f.header.Open()

	if err != nil {
		return false
	}

	file.Close()

	return true
}

// Open returns a reader for the file contents.
func (f *UploadedFile) Open() (io.ReadCloser, error) {
	if f.test && f.testContent != nil {
		return io.NopCloser(bytes.NewReader(f.testContent)), nil
	}

	return f.header.Open()
}

// Get reads and returns the entire file content as bytes.
func (f *UploadedFile) Get() ([]byte, error) {
	rc, err := f.Open()

	if err != nil {
		return nil, err
	}

	defer rc.Close()

	return io.ReadAll(rc)
}

// Store saves the file to the given directory using a random hash name.
func (f *UploadedFile) Store(directory string, fs FileStore) (string, error) {
	return f.StoreAs(directory, f.HashName(), fs)
}

// StoreAs saves the file to the given directory with a specific name.
func (f *UploadedFile) StoreAs(directory, name string, fs FileStore) (string, error) {
	rc, err := f.Open()

	if err != nil {
		return "", err
	}

	defer rc.Close()

	path := filepath.Join(directory, name)

	if err := fs.Put(path, rc); err != nil {
		return "", err
	}

	return path, nil
}

// StorePublicly stores the file with public visibility. In Go the visibility
// semantics depend on the FileStore implementation; this is a convenience alias
// that matches the upstream API.
func (f *UploadedFile) StorePublicly(directory string, fs FileStore) (string, error) {
	return f.Store(directory, fs)
}

// StorePubliclyAs stores the file with a specific name and public visibility.
func (f *UploadedFile) StorePubliclyAs(directory, name string, fs FileStore) (string, error) {
	return f.StoreAs(directory, name, fs)
}

// HashName generates a random hex name preserving the client extension.
func (f *UploadedFile) HashName(directory ...string) string {
	b := make([]byte, 20)
	_, _ = rand.Read(b)

	ext := f.ClientExtension()
	name := hex.EncodeToString(b)

	if ext != "" {
		name += "." + ext
	}

	if len(directory) > 0 && directory[0] != "" {
		return filepath.Join(directory[0], name)
	}

	return name
}

// Path returns the temporary path where the file is stored, if available.
func (f *UploadedFile) Path() string {
	if f.test {
		return f.ClientOriginalName()
	}

	return ""
}

// Extension returns the guessed extension based on MIME type.
func (f *UploadedFile) Extension() string {
	mimeType := f.ClientMimeType()

	if mimeType == "" {
		return f.ClientExtension()
	}

	exts, err := mime.ExtensionsByType(mimeType)

	if err != nil || len(exts) == 0 {
		return f.ClientExtension()
	}

	return strings.TrimPrefix(exts[0], ".")
}

// Header returns the underlying multipart.FileHeader.
func (f *UploadedFile) Header() *multipart.FileHeader {
	return f.header
}

// CreateFromBase64 constructs an UploadedFile from a base64-encoded string.
func CreateFromBase64(encoded, name string) (*UploadedFile, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)

	if err != nil {
		return nil, err
	}

	ct := http.DetectContentType(data)

	h := &multipart.FileHeader{
		Filename: name,
		Size:     int64(len(data)),
		Header:   make(map[string][]string),
	}
	h.Header.Set("Content-Type", ct)

	return &UploadedFile{
		header:      h,
		test:        true,
		testContent: data,
	}, nil
}
