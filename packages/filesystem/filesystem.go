package filesystem

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"hash"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	contract "github.com/bedrock/packages/contracts/filesystem"
)

// Filesystem provides local filesystem operations.
type Filesystem struct{}

var _ contract.Filesystem = (*Filesystem)(nil)

// New creates a Filesystem instance.
func New() *Filesystem {
	return &Filesystem{}
}

// Exists determines if a file or directory exists at the given path.
func (f *Filesystem) Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Missing determines if a file or directory is missing at the given path.
func (f *Filesystem) Missing(path string) bool {
	return !f.Exists(path)
}

// IsFile determines if the given path is a regular file.
func (f *Filesystem) IsFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

// IsDirectory determines if the given path is a directory.
func (f *Filesystem) IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// IsEmptyDirectory determines if the given directory is empty.
// When ignoreDotFiles is true, files starting with a dot are excluded.
func (f *Filesystem) IsEmptyDirectory(directory string, ignoreDotFiles bool) (bool, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return false, err
	}

	for _, entry := range entries {
		if ignoreDotFiles && strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		return false, nil
	}

	return true, nil
}

// IsReadable determines if the given path is readable.
func (f *Filesystem) IsReadable(path string) bool {
	file, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return false
	}
	file.Close()
	return true
}

// IsWritable determines if the given path is writable.
func (f *Filesystem) IsWritable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	// For directories, try to create a temp file.
	if info.IsDir() {
		tmp, err := os.CreateTemp(path, ".writable_check_*")
		if err != nil {
			return false
		}
		name := tmp.Name()
		tmp.Close()
		os.Remove(name)
		return true
	}

	// For files, try to open for writing.
	file, err := os.OpenFile(path, os.O_WRONLY, 0)
	if err != nil {
		return false
	}
	file.Close()
	return true
}

// Hash calculates the hash of a file. The algorithm defaults to "md5".
// Supported algorithms: "md5", "sha1", "sha256", "sha512".
func (f *Filesystem) Hash(path string, algorithm ...string) (string, error) {
	algo := "md5"
	if len(algorithm) > 0 {
		algo = algorithm[0]
	}

	var h hash.Hash
	switch strings.ToLower(algo) {
	case "md5":
		h = md5.New()
	case "sha1":
		h = sha1.New()
	case "sha256":
		h = sha256.New()
	case "sha512":
		h = sha512.New()
	default:
		return "", ErrHashAlgorithm
	}

	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// HasSameHash determines if two files have the same hash.
func (f *Filesystem) HasSameHash(firstFile, secondFile string) (bool, error) {
	h1, err := f.Hash(firstFile)
	if err != nil {
		return false, err
	}

	h2, err := f.Hash(secondFile)
	if err != nil {
		return false, err
	}

	return h1 == h2, nil
}

// Type returns the file type: "file" or "dir".
func (f *Filesystem) Type(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if info.IsDir() {
		return "dir", nil
	}

	return "file", nil
}

// MimeType returns the MIME type of a file.
func (f *Filesystem) MimeType(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Read up to 512 bytes for content sniffing.
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return "", err
	}

	detected := http.DetectContentType(buf[:n])

	// If detection returned the generic fallback, try extension-based lookup.
	if detected == "application/octet-stream" {
		ext := filepath.Ext(path)
		if ext != "" {
			if mtype := mime.TypeByExtension(ext); mtype != "" {
				return mtype, nil
			}
		}
	}

	return detected, nil
}

// GuessExtension returns the extension for a file based on its MIME type.
func (f *Filesystem) GuessExtension(path string) (string, error) {
	mtype, err := f.MimeType(path)
	if err != nil {
		return "", err
	}

	exts, err := mime.ExtensionsByType(mtype)
	if err != nil || len(exts) == 0 {
		return "", nil
	}

	// Return without leading dot.
	return strings.TrimPrefix(exts[0], "."), nil
}

// Size returns the file size in bytes.
func (f *Filesystem) Size(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}

// LastModified returns the last modification time as a Unix timestamp.
func (f *Filesystem) LastModified(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	return info.ModTime().Unix(), nil
}

// Chmod sets the permission mode of a file or directory.
func (f *Filesystem) Chmod(path string, mode fs.FileMode) error {
	return os.Chmod(path, mode)
}

// Name returns the filename without the extension.
func (f *Filesystem) Name(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}

// Basename returns the trailing name component of the path.
func (f *Filesystem) Basename(path string) string {
	return filepath.Base(path)
}

// Dirname returns the parent directory of the path.
func (f *Filesystem) Dirname(path string) string {
	return filepath.Dir(path)
}

// Extension returns the file extension without the leading dot.
func (f *Filesystem) Extension(path string) string {
	return strings.TrimPrefix(filepath.Ext(path), ".")
}

// Glob finds path names matching a pattern.
func (f *Filesystem) Glob(pattern string) ([]string, error) {
	return filepath.Glob(pattern)
}
