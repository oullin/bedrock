package filesystem

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Files returns the files in the given directory (non-recursive by default).
// When hidden is true (or omitted), hidden files (starting with ".") are included.
// When hidden is explicitly set to false, hidden files are excluded.
func (f *Filesystem) Files(directory string, hidden ...bool) ([]string, error) {
	includeHidden := false

	if len(hidden) > 0 {
		includeHidden = hidden[0]
	}

	entries, err := os.ReadDir(directory)

	if err != nil {
		return nil, err
	}

	var files []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if !includeHidden && strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		files = append(files, filepath.Join(directory, entry.Name()))
	}

	return files, nil
}

// AllFiles returns all files in the directory tree recursively.
func (f *Filesystem) AllFiles(directory string, hidden ...bool) ([]string, error) {
	includeHidden := false

	if len(hidden) > 0 {
		includeHidden = hidden[0]
	}

	var files []string

	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if !includeHidden && path != directory && strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}

			return nil
		}

		if !includeHidden && strings.HasPrefix(d.Name(), ".") {
			return nil
		}

		files = append(files, path)

		return nil
	})

	return files, err
}

// Directories returns the directories in the given directory (non-recursive).
func (f *Filesystem) Directories(directory string) ([]string, error) {
	entries, err := os.ReadDir(directory)

	if err != nil {
		return nil, err
	}

	var dirs []string

	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, filepath.Join(directory, entry.Name()))
		}
	}

	return dirs, nil
}

// AllDirectories returns all directories in the directory tree recursively.
func (f *Filesystem) AllDirectories(directory string) ([]string, error) {
	var dirs []string

	err := filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && path != directory {
			dirs = append(dirs, path)
		}

		return nil
	})

	return dirs, err
}

// EnsureDirectoryExists creates the directory if it does not already exist.
// The default mode is 0755.
func (f *Filesystem) EnsureDirectoryExists(path string, mode ...fs.FileMode) error {
	perm := fs.FileMode(0o755)

	if len(mode) > 0 {
		perm = mode[0]
	}

	return os.MkdirAll(path, perm)
}

// MakeDirectory creates a directory. The default mode is 0755.
// By default it creates parent directories recursively.
func (f *Filesystem) MakeDirectory(path string, mode ...fs.FileMode) error {
	perm := fs.FileMode(0o755)

	if len(mode) > 0 {
		perm = mode[0]
	}

	return os.MkdirAll(path, perm)
}

// MoveDirectory moves a directory from one location to another.
// When overwrite is true, any existing destination directory is removed first.
func (f *Filesystem) MoveDirectory(from, to string, overwrite ...bool) error {
	shouldOverwrite := false

	if len(overwrite) > 0 {
		shouldOverwrite = overwrite[0]
	}

	if shouldOverwrite {
		if err := os.RemoveAll(to); err != nil {
			return err
		}
	}

	// Try atomic rename first.
	if err := os.Rename(from, to); err == nil {
		return nil
	}

	// Fallback: copy + delete (handles cross-device moves).
	if err := f.CopyDirectory(from, to); err != nil {
		return err
	}

	return os.RemoveAll(from)
}

// CopyDirectory recursively copies a directory and its contents.
func (f *Filesystem) CopyDirectory(directory, destination string) error {
	info, err := os.Stat(directory)

	if err != nil {
		return err
	}

	if !info.IsDir() {
		return ErrNotDirectory
	}

	return filepath.WalkDir(directory, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(directory, path)

		if err != nil {
			return err
		}

		target := filepath.Join(destination, rel)

		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		return f.Copy(path, target)
	})
}

// DeleteDirectory removes a directory and all of its contents.
// When preserve is true, the directory itself is kept but its contents are removed.
func (f *Filesystem) DeleteDirectory(directory string, preserve ...bool) error {
	info, err := os.Stat(directory)

	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	if !info.IsDir() {
		return ErrNotDirectory
	}

	shouldPreserve := false

	if len(preserve) > 0 {
		shouldPreserve = preserve[0]
	}

	if shouldPreserve {
		return f.cleanDir(directory)
	}

	return os.RemoveAll(directory)
}

// DeleteDirectories removes all subdirectories within the given directory,
// leaving files intact.
func (f *Filesystem) DeleteDirectories(directory string) error {
	entries, err := os.ReadDir(directory)

	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			if err := os.RemoveAll(filepath.Join(directory, entry.Name())); err != nil {
				return err
			}
		}
	}

	return nil
}

// CleanDirectory removes all contents of a directory but keeps the directory.
func (f *Filesystem) CleanDirectory(directory string) error {
	return f.cleanDir(directory)
}

func (f *Filesystem) cleanDir(directory string) error {
	entries, err := os.ReadDir(directory)

	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := filepath.Join(directory, entry.Name())

		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}

	return nil
}
