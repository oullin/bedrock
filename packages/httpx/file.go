package httpx

import (
	"crypto/rand"
	"encoding/hex"
	"path/filepath"
	"strings"
)

// File represents a file on disk with convenience path helpers.
type File struct {
	path string
}

// NewFile creates a File from the given filesystem path.
func NewFile(path string) *File {
	return &File{path: path}
}

// Path returns the full filesystem path.
func (f *File) Path() string {
	return f.path
}

// Basename returns the file name without its directory.
func (f *File) Basename() string {
	return filepath.Base(f.path)
}

// Extension returns the file extension without the leading dot.
func (f *File) Extension() string {
	ext := filepath.Ext(f.path)

	return strings.TrimPrefix(ext, ".")
}

// HashName generates a random hex name preserving the original extension.
// An optional directory prefix is prepended to the result.
func (f *File) HashName(directory ...string) string {
	b := make([]byte, 20)
	_, _ = rand.Read(b)

	name := hex.EncodeToString(b) + "." + f.Extension()

	if len(directory) > 0 && directory[0] != "" {
		return filepath.Join(directory[0], name)
	}

	return name
}
