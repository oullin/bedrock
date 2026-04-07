package session

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// FileHandler stores sessions as individual files on disk. It is safe for
// concurrent use via the filesystem.
type FileHandler struct {
	path    string
	minutes int
}

// NewFileHandler creates a file-based session handler. Sessions are stored
// as files in path, and expire after minutes of inactivity.
func NewFileHandler(path string, minutes int) *FileHandler {
	return &FileHandler{
		path:    path,
		minutes: minutes,
	}
}

func (h *FileHandler) Open(_ context.Context, _ string, _ string) error { return nil }
func (h *FileHandler) Close(_ context.Context) error                    { return nil }

// Read returns the session data stored in the file identified by id. It
// returns an empty string if the file does not exist or has expired.
func (h *FileHandler) Read(_ context.Context, id string) (string, error) {
	name := filepath.Join(h.path, id)

	info, err := os.Stat(name)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", nil
		}

		return "", err
	}

	if time.Since(info.ModTime()) > time.Duration(h.minutes)*time.Minute {
		return "", nil
	}

	data, err := os.ReadFile(name)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Write stores session data in a file identified by id.
func (h *FileHandler) Write(_ context.Context, id string, data string) error {
	return os.WriteFile(filepath.Join(h.path, id), []byte(data), 0600)
}

// Destroy removes the session file for the given id.
func (h *FileHandler) Destroy(_ context.Context, id string) error {
	err := os.Remove(filepath.Join(h.path, id))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	return err
}

// GC removes session files that have not been modified within maxLifetime
// seconds.
func (h *FileHandler) GC(_ context.Context, maxLifetime int) error {
	entries, err := os.ReadDir(h.path)
	if err != nil {
		return err
	}

	cutoff := time.Duration(maxLifetime) * time.Second

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		if time.Since(info.ModTime()) > cutoff {
			_ = os.Remove(filepath.Join(h.path, entry.Name()))
		}
	}

	return nil
}
