package filesystem

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// LockableFile provides file operations with advisory locking via flock(2).
type LockableFile struct {
	path string
	file *os.File
}

// NewLockableFile opens or creates a file at path for reading and writing.
func NewLockableFile(path string, mode fs.FileMode) (*LockableFile, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, mode)

	if err != nil {
		return nil, err
	}

	return &LockableFile{path: path, file: file}, nil
}

// SharedLock acquires a shared (read) lock on the file.
func (lf *LockableFile) SharedLock() error {
	if err := lockShared(lf.file); err != nil {
		return ErrLockFailed
	}

	return nil
}

// ExclusiveLock acquires an exclusive (write) lock on the file.
func (lf *LockableFile) ExclusiveLock() error {
	if err := lockExclusive(lf.file); err != nil {
		return ErrLockFailed
	}

	return nil
}

// Unlock releases the advisory lock on the file.
func (lf *LockableFile) Unlock() error {
	return unlockFile(lf.file)
}

// Read reads up to size bytes from the file. If no size is specified, reads
// the entire file.
func (lf *LockableFile) Read(size ...int) ([]byte, error) {
	if _, err := lf.file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	if len(size) > 0 && size[0] > 0 {
		buf := make([]byte, size[0])
		n, err := lf.file.Read(buf)

		if err != nil && err != io.EOF {
			return nil, err
		}

		return buf[:n], nil
	}

	return io.ReadAll(lf.file)
}

// Write writes contents to the file, starting from the beginning.
func (lf *LockableFile) Write(contents []byte) (int, error) {
	if _, err := lf.file.Seek(0, io.SeekStart); err != nil {
		return 0, err
	}

	return lf.file.Write(contents)
}

// Truncate truncates the file to zero length.
func (lf *LockableFile) Truncate() error {
	return lf.file.Truncate(0)
}

// Close releases any lock and closes the file.
func (lf *LockableFile) Close() error {
	_ = lf.Unlock()

	return lf.file.Close()
}

// Path returns the file path.
func (lf *LockableFile) Path() string {
	return lf.path
}

// Size returns the current file size.
func (lf *LockableFile) Size() (int64, error) {
	info, err := lf.file.Stat()

	if err != nil {
		return 0, err
	}

	return info.Size(), nil
}

// Chmod sets the file permissions.
func (lf *LockableFile) Chmod(mode fs.FileMode) error {
	return lf.file.Chmod(mode)
}
