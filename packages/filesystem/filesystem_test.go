package filesystem_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bedrock/packages/filesystem"
)

func TestExists(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, path, "hello")

	if !fs.Exists(path) {
		t.Fatal("expected file to exist")
	}

	if fs.Exists(filepath.Join(dir, "nonexistent.txt")) {
		t.Fatal("expected file to not exist")
	}
}

func TestMissing(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()

	if !fs.Missing(filepath.Join(dir, "nonexistent.txt")) {
		t.Fatal("expected file to be missing")
	}

	path := filepath.Join(dir, "file.txt")
	writeFile(t, path, "hello")

	if fs.Missing(path) {
		t.Fatal("expected file to not be missing")
	}
}

func TestIsFile(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, path, "hello")

	if !fs.IsFile(path) {
		t.Fatal("expected path to be a file")
	}

	if fs.IsFile(dir) {
		t.Fatal("expected directory to not be a file")
	}

	if fs.IsFile(filepath.Join(dir, "nonexistent.txt")) {
		t.Fatal("expected nonexistent path to not be a file")
	}
}

func TestIsDirectory(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, path, "hello")

	if !fs.IsDirectory(dir) {
		t.Fatal("expected path to be a directory")
	}

	if fs.IsDirectory(path) {
		t.Fatal("expected file to not be a directory")
	}
}

func TestIsEmptyDirectory(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()

	empty, err := fs.IsEmptyDirectory(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	if !empty {
		t.Fatal("expected directory to be empty")
	}

	writeFile(t, filepath.Join(dir, "file.txt"), "hello")

	empty, err = fs.IsEmptyDirectory(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	if empty {
		t.Fatal("expected directory to not be empty")
	}
}

func TestIsEmptyDirectoryIgnoreDotFiles(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, ".hidden"), "hidden")

	empty, err := fs.IsEmptyDirectory(dir, true)
	if err != nil {
		t.Fatal(err)
	}
	if !empty {
		t.Fatal("expected directory to be empty when ignoring dot files")
	}

	empty, err = fs.IsEmptyDirectory(dir, false)
	if err != nil {
		t.Fatal(err)
	}
	if empty {
		t.Fatal("expected directory to not be empty when not ignoring dot files")
	}
}

func TestIsReadable(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, path, "hello")

	if !fs.IsReadable(path) {
		t.Fatal("expected file to be readable")
	}

	if fs.IsReadable(filepath.Join(dir, "nonexistent.txt")) {
		t.Fatal("expected nonexistent file to not be readable")
	}
}

func TestIsWritable(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, path, "hello")

	if !fs.IsWritable(path) {
		t.Fatal("expected file to be writable")
	}

	if !fs.IsWritable(dir) {
		t.Fatal("expected directory to be writable")
	}

	if fs.IsWritable(filepath.Join(dir, "nonexistent.txt")) {
		t.Fatal("expected nonexistent file to not be writable")
	}
}

func TestHash(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, path, "hello")

	h, err := fs.Hash(path)
	if err != nil {
		t.Fatal(err)
	}
	if h == "" {
		t.Fatal("expected non-empty hash")
	}

	// MD5 of "hello" is 5d41402abc4b2a76b9719d911017c592
	if h != "5d41402abc4b2a76b9719d911017c592" {
		t.Fatalf("unexpected md5 hash: %s", h)
	}
}

func TestHashSHA1(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, path, "hello")

	h, err := fs.Hash(path, "sha1")
	if err != nil {
		t.Fatal(err)
	}
	// SHA1 of "hello" is aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d
	if h != "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d" {
		t.Fatalf("unexpected sha1 hash: %s", h)
	}
}

func TestHashSHA256(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, path, "hello")

	h, err := fs.Hash(path, "sha256")
	if err != nil {
		t.Fatal(err)
	}
	if h != "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824" {
		t.Fatalf("unexpected sha256 hash: %s", h)
	}
}

func TestHashUnsupportedAlgorithm(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, path, "hello")

	_, err := fs.Hash(path, "blake2b")
	if err != filesystem.ErrHashAlgorithm {
		t.Fatalf("expected ErrHashAlgorithm, got %v", err)
	}
}

func TestHasSameHash(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path1 := filepath.Join(dir, "file1.txt")
	path2 := filepath.Join(dir, "file2.txt")
	path3 := filepath.Join(dir, "file3.txt")

	writeFile(t, path1, "hello")
	writeFile(t, path2, "hello")
	writeFile(t, path3, "world")

	same, err := fs.HasSameHash(path1, path2)
	if err != nil {
		t.Fatal(err)
	}
	if !same {
		t.Fatal("expected files with same content to have same hash")
	}

	same, err = fs.HasSameHash(path1, path3)
	if err != nil {
		t.Fatal(err)
	}
	if same {
		t.Fatal("expected files with different content to have different hash")
	}
}

func TestType(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, path, "hello")

	typ, err := fs.Type(path)
	if err != nil {
		t.Fatal(err)
	}
	if typ != "file" {
		t.Fatalf("expected 'file', got %q", typ)
	}

	typ, err = fs.Type(dir)
	if err != nil {
		t.Fatal(err)
	}
	if typ != "dir" {
		t.Fatalf("expected 'dir', got %q", typ)
	}
}

func TestMimeType(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()

	htmlPath := filepath.Join(dir, "page.html")
	writeFile(t, htmlPath, "<html><body>Hello</body></html>")

	mtype, err := fs.MimeType(htmlPath)
	if err != nil {
		t.Fatal(err)
	}
	if mtype != "text/html; charset=utf-8" {
		t.Fatalf("unexpected MIME type: %s", mtype)
	}
}

func TestSize(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, path, "hello")

	size, err := fs.Size(path)
	if err != nil {
		t.Fatal(err)
	}
	if size != 5 {
		t.Fatalf("expected size 5, got %d", size)
	}
}

func TestLastModified(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, path, "hello")

	ts, err := fs.LastModified(path)
	if err != nil {
		t.Fatal(err)
	}
	if ts <= 0 {
		t.Fatal("expected positive timestamp")
	}
}

func TestChmod(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, path, "hello")

	if err := fs.Chmod(path, 0o755); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	if info.Mode().Perm() != 0o755 {
		t.Fatalf("expected 0755 permissions, got %o", info.Mode().Perm())
	}
}

func TestName(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()

	if fs.Name("/path/to/file.txt") != "file" {
		t.Fatalf("unexpected name: %s", fs.Name("/path/to/file.txt"))
	}

	if fs.Name("/path/to/file.tar.gz") != "file.tar" {
		t.Fatalf("unexpected name: %s", fs.Name("/path/to/file.tar.gz"))
	}
}

func TestBasename(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()

	if fs.Basename("/path/to/file.txt") != "file.txt" {
		t.Fatalf("unexpected basename: %s", fs.Basename("/path/to/file.txt"))
	}
}

func TestDirname(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()

	if fs.Dirname("/path/to/file.txt") != "/path/to" {
		t.Fatalf("unexpected dirname: %s", fs.Dirname("/path/to/file.txt"))
	}
}

func TestExtension(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()

	if fs.Extension("/path/to/file.txt") != "txt" {
		t.Fatalf("unexpected extension: %s", fs.Extension("/path/to/file.txt"))
	}

	if fs.Extension("/path/to/file") != "" {
		t.Fatalf("unexpected extension for file without ext: %s", fs.Extension("/path/to/file"))
	}
}

func TestGlob(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "a.txt"), "a")
	writeFile(t, filepath.Join(dir, "b.txt"), "b")
	writeFile(t, filepath.Join(dir, "c.go"), "c")

	matches, err := fs.Glob(filepath.Join(dir, "*.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}
}

func TestGuessExtension(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()

	txtPath := filepath.Join(dir, "file.txt")
	writeFile(t, txtPath, "plain text content")

	ext, err := fs.GuessExtension(txtPath)
	if err != nil {
		t.Fatal(err)
	}
	if ext == "" {
		t.Fatal("expected non-empty extension")
	}
}
