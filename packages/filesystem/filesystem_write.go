package filesystem

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Put writes contents to a file, creating it if necessary.
// An optional file mode can be provided; defaults to 0644.
func (f *Filesystem) Put(path string, contents []byte, mode ...fs.FileMode) error {
	perm := fs.FileMode(0o644)

	if len(mode) > 0 {
		perm = mode[0]
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(path, contents, perm)
}

// Replace atomically writes content to a file using a temporary file
// and rename. An optional file mode can be provided.
func (f *Filesystem) Replace(path string, content []byte, mode ...fs.FileMode) error {
	perm := fs.FileMode(0o644)

	if len(mode) > 0 {
		perm = mode[0]
	}

	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(dir, ".tmp_*")

	if err != nil {
		return err
	}

	tmpName := tmp.Name()

	if _, err := tmp.Write(content); err != nil {
		tmp.Close()
		os.Remove(tmpName)

		return err
	}

	if err := tmp.Close(); err != nil {
		os.Remove(tmpName)

		return err
	}

	if err := os.Chmod(tmpName, perm); err != nil {
		os.Remove(tmpName)

		return err
	}

	return os.Rename(tmpName, path)
}

// ReplaceInFile replaces all occurrences of search with replace in the file.
func (f *Filesystem) ReplaceInFile(search, replace, path string) error {
	data, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	content := strings.ReplaceAll(string(data), search, replace)

	return os.WriteFile(path, []byte(content), 0o644)
}

// Prepend prepends data to the beginning of a file. If the file does not
// exist, it is created with only the given data.
func (f *Filesystem) Prepend(path string, data []byte) error {
	existing, err := os.ReadFile(path)

	if err != nil {
		if os.IsNotExist(err) {
			return f.Put(path, data)
		}

		return err
	}

	combined := make([]byte, 0, len(data)+len(existing))
	combined = append(combined, data...)
	combined = append(combined, existing...)

	return os.WriteFile(path, combined, 0o644)
}

// Append appends data to the end of a file.
func (f *Filesystem) Append(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(data)

	return err
}
