package handlers

import (
	"context"
	"os"
	"path/filepath"
	"time"
)

// FileHandler stores session data as files on the filesystem.
type FileHandler struct {
	path        string
	maxLifetime int
}

// NewFileHandler creates a FileHandler that stores sessions under path.
func NewFileHandler(path string) *FileHandler {
	return &FileHandler{path: path}
}

func (h *FileHandler) Open(_ context.Context, path, _ string) error {
	if path != "" {
		h.path = path
	}

	return os.MkdirAll(h.path, 0o755)
}

func (h *FileHandler) Close(_ context.Context) error { return nil }

func (h *FileHandler) Read(_ context.Context, id string) (string, error) {
	path := h.filePath(id)
	info, err := os.Stat(path)

	if os.IsNotExist(err) {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	if h.maxLifetime > 0 {
		cutoff := time.Now().Add(-time.Duration(h.maxLifetime) * time.Second)

		if info.ModTime().Before(cutoff) {
			return "", nil
		}
	}

	data, err := os.ReadFile(path)

	if os.IsNotExist(err) {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (h *FileHandler) Write(_ context.Context, id, data string) error {
	tmp, err := os.CreateTemp(h.path, "sess-*")

	if err != nil {
		return err
	}

	if _, err = tmp.WriteString(data); err != nil {
		tmp.Close()
		os.Remove(tmp.Name())

		return err
	}

	if err = tmp.Close(); err != nil {
		os.Remove(tmp.Name())

		return err
	}

	return os.Rename(tmp.Name(), h.filePath(id))
}

func (h *FileHandler) Destroy(_ context.Context, id string) error {
	err := os.Remove(h.filePath(id))

	if os.IsNotExist(err) {
		return nil
	}

	return err
}

func (h *FileHandler) GC(_ context.Context, maxLifetime int) error {
	lifetime := maxLifetime

	if lifetime <= 0 {
		lifetime = h.maxLifetime
	}

	if lifetime <= 0 {
		return nil
	}

	cutoff := time.Now().Add(-time.Duration(lifetime) * time.Second)

	entries, err := os.ReadDir(h.path)

	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()

		if err != nil {
			continue
		}

		if info.ModTime().Before(cutoff) {
			os.Remove(filepath.Join(h.path, entry.Name()))
		}
	}

	return nil
}

func (h *FileHandler) filePath(id string) string {
	return filepath.Join(h.path, id)
}
