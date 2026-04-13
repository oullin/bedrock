package filesystem

import (
	"io/fs"
	"iter"
)

// Filesystem defines the contract for local filesystem operations.
type Filesystem interface {
	// Existence & detection
	Exists(path string) bool
	Missing(path string) bool
	IsFile(path string) bool
	IsDirectory(path string) bool
	IsEmptyDirectory(directory string, ignoreDotFiles bool) (bool, error)
	IsReadable(path string) bool
	IsWritable(path string) bool

	// Reading
	Get(path string) ([]byte, error)
	JSON(path string, v any) error
	SharedGet(path string) ([]byte, error)
	Lines(path string) (iter.Seq[string], error)

	// Writing
	Put(path string, contents []byte, mode ...fs.FileMode) error
	Replace(path string, content []byte, mode ...fs.FileMode) error
	ReplaceInFile(search, replace, path string) error
	Prepend(path string, data []byte) error
	Append(path string, data []byte) error

	// Metadata
	Hash(path string, algorithm ...string) (string, error)
	HasSameHash(firstFile, secondFile string) (bool, error)
	Type(path string) (string, error)
	MimeType(path string) (string, error)
	GuessExtension(path string) (string, error)
	Size(path string) (int64, error)
	LastModified(path string) (int64, error)
	Chmod(path string, mode fs.FileMode) error

	// Path operations
	Name(path string) string
	Basename(path string) string
	Dirname(path string) string
	Extension(path string) string
	Glob(pattern string) ([]string, error)

	// File operations
	Delete(paths ...string) error
	Move(path, target string) error
	Copy(path, target string) error
	Link(target, link string) error
	RelativeLink(target, link string) error

	// Directory operations
	Files(directory string, hidden ...bool) ([]string, error)
	AllFiles(directory string, hidden ...bool) ([]string, error)
	Directories(directory string) ([]string, error)
	AllDirectories(directory string) ([]string, error)
	EnsureDirectoryExists(path string, mode ...fs.FileMode) error
	MakeDirectory(path string, mode ...fs.FileMode) error
	MoveDirectory(from, to string, overwrite ...bool) error
	CopyDirectory(directory, destination string) error
	DeleteDirectory(directory string, preserve ...bool) error
	DeleteDirectories(directory string) error
	CleanDirectory(directory string) error
}
