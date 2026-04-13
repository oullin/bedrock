package filesystem_test

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPut(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	if err := fs.Put(path, []byte("hello")); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)

	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "hello" {
		t.Fatalf("expected 'hello', got %q", string(data))
	}
}

func TestPutCreatesDirectories(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "deep", "file.txt")

	if err := fs.Put(path, []byte("nested")); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)

	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "nested" {
		t.Fatalf("expected 'nested', got %q", string(data))
	}
}

func TestPutWithMode(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	if err := fs.Put(path, []byte("hello"), 0o755); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(path)

	if err != nil {
		t.Fatal(err)
	}

	if info.Mode().Perm() != 0o755 {
		t.Fatalf("expected 0755, got %o", info.Mode().Perm())
	}
}

func TestReplace(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, path, "original")

	if err := fs.Replace(path, []byte("replaced")); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)

	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "replaced" {
		t.Fatalf("expected 'replaced', got %q", string(data))
	}
}

func TestReplaceCreatesNewFile(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "new.txt")

	if err := fs.Replace(path, []byte("new content")); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)

	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "new content" {
		t.Fatalf("expected 'new content', got %q", string(data))
	}
}

func TestReplaceWithMode(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	if err := fs.Replace(path, []byte("content"), 0o755); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(path)

	if err != nil {
		t.Fatal(err)
	}

	if info.Mode().Perm() != 0o755 {
		t.Fatalf("expected 0755, got %o", info.Mode().Perm())
	}
}

func TestReplaceSymlink(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	target := filepath.Join(dir, "target.txt")
	link := filepath.Join(dir, "link.txt")

	writeFile(t, target, "original")

	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}

	if err := fs.Replace(link, []byte("via symlink")); err != nil {
		t.Fatal(err)
	}

	// The symlink now points to a new file with the replaced content.
	data, err := os.ReadFile(link)

	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "via symlink" {
		t.Fatalf("expected 'via symlink', got %q", string(data))
	}
}

func TestReplaceInFile(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, path, "Hello World! Hello Go!")

	if err := fs.ReplaceInFile("Hello", "Hi", path); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)

	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "Hi World! Hi Go!" {
		t.Fatalf("expected 'Hi World! Hi Go!', got %q", string(data))
	}
}

func TestPrepend(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, path, "World")

	if err := fs.Prepend(path, []byte("Hello ")); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)

	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "Hello World" {
		t.Fatalf("expected 'Hello World', got %q", string(data))
	}
}

func TestPrependNewFile(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "new.txt")

	if err := fs.Prepend(path, []byte("first")); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)

	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "first" {
		t.Fatalf("expected 'first', got %q", string(data))
	}
}

func TestAppend(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, path, "Hello")

	if err := fs.Append(path, []byte(" World")); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)

	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "Hello World" {
		t.Fatalf("expected 'Hello World', got %q", string(data))
	}
}

func TestAppendNewFile(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "new.txt")

	if err := fs.Append(path, []byte("appended")); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)

	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "appended" {
		t.Fatalf("expected 'appended', got %q", string(data))
	}
}

func TestAppendMultiple(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, path, "a")

	if err := fs.Append(path, []byte("b")); err != nil {
		t.Fatal(err)
	}

	if err := fs.Append(path, []byte("c")); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)

	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "abc" {
		t.Fatalf("expected 'abc', got %q", string(data))
	}
}
