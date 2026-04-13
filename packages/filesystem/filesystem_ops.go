package filesystem

import (
	"io"
	"os"
	"path/filepath"
)

// Delete removes one or more files. Returns an error if any file cannot
// be removed.
func (f *Filesystem) Delete(paths ...string) error {
	for _, path := range paths {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	return nil
}

// Move moves a file from path to target.
func (f *Filesystem) Move(path, target string) error {
	return os.Rename(path, target)
}

// Copy copies a file from path to target.
func (f *Filesystem) Copy(path, target string) error {
	src, err := os.Open(path)

	if err != nil {
		return err
	}

	defer src.Close()

	info, err := src.Stat()

	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}

	dst, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, info.Mode())

	if err != nil {
		return err
	}

	defer dst.Close()

	_, err = io.Copy(dst, src)

	return err
}

// Link creates a symbolic link.
func (f *Filesystem) Link(target, link string) error {
	return os.Symlink(target, link)
}

// RelativeLink creates a relative symbolic link.
func (f *Filesystem) RelativeLink(target, link string) error {
	linkDir := filepath.Dir(link)

	rel, err := filepath.Rel(linkDir, target)

	if err != nil {
		return err
	}

	return os.Symlink(rel, link)
}
