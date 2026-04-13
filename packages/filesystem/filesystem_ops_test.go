package filesystem_test

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, path, "hello")

	if err := fs.Delete(path); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("expected file to be deleted")
	}
}

func TestDeleteMultiple(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path1 := filepath.Join(dir, "file1.txt")
	path2 := filepath.Join(dir, "file2.txt")

	writeFile(t, path1, "a")
	writeFile(t, path2, "b")

	if err := fs.Delete(path1, path2); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(path1); !os.IsNotExist(err) {
		t.Fatal("expected file1 to be deleted")
	}

	if _, err := os.Stat(path2); !os.IsNotExist(err) {
		t.Fatal("expected file2 to be deleted")
	}
}

func TestDeleteNonexistentFile(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()

	if err := fs.Delete(filepath.Join(dir, "nonexistent.txt")); err != nil {
		t.Fatalf("delete of nonexistent file should not error, got %v", err)
	}
}

func TestMove(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	src := filepath.Join(dir, "source.txt")
	dst := filepath.Join(dir, "dest.txt")

	writeFile(t, src, "content")

	if err := fs.Move(src, dst); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Fatal("expected source to be gone")
	}

	data, err := os.ReadFile(dst)

	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "content" {
		t.Fatalf("expected 'content', got %q", string(data))
	}
}

func TestCopy(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	src := filepath.Join(dir, "source.txt")
	dst := filepath.Join(dir, "dest.txt")

	writeFile(t, src, "content")

	if err := fs.Copy(src, dst); err != nil {
		t.Fatal(err)
	}

	// Source should still exist.
	srcData, err := os.ReadFile(src)

	if err != nil {
		t.Fatal(err)
	}

	if string(srcData) != "content" {
		t.Fatalf("expected source to be unchanged, got %q", string(srcData))
	}

	dstData, err := os.ReadFile(dst)

	if err != nil {
		t.Fatal(err)
	}

	if string(dstData) != "content" {
		t.Fatalf("expected 'content', got %q", string(dstData))
	}
}

func TestCopyPreservesPermissions(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	src := filepath.Join(dir, "source.txt")
	dst := filepath.Join(dir, "dest.txt")

	writeFile(t, src, "content")
	os.Chmod(src, 0o755)

	if err := fs.Copy(src, dst); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(dst)

	if err != nil {
		t.Fatal(err)
	}

	if info.Mode().Perm() != 0o755 {
		t.Fatalf("expected 0755, got %o", info.Mode().Perm())
	}
}

func TestLink(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	target := filepath.Join(dir, "target.txt")
	link := filepath.Join(dir, "link.txt")

	writeFile(t, target, "linked")

	if err := fs.Link(target, link); err != nil {
		t.Fatal(err)
	}

	resolved, err := os.Readlink(link)

	if err != nil {
		t.Fatal(err)
	}

	if resolved != target {
		t.Fatalf("expected link to point to %q, got %q", target, resolved)
	}

	data, err := os.ReadFile(link)

	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "linked" {
		t.Fatalf("expected 'linked', got %q", string(data))
	}
}

func TestRelativeLink(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	target := filepath.Join(dir, "sub", "target.txt")
	link := filepath.Join(dir, "link.txt")

	writeFile(t, target, "relative")

	if err := fs.RelativeLink(target, link); err != nil {
		t.Fatal(err)
	}

	resolved, err := os.Readlink(link)

	if err != nil {
		t.Fatal(err)
	}

	// Should be a relative path.
	if filepath.IsAbs(resolved) {
		t.Fatalf("expected relative symlink, got absolute: %q", resolved)
	}

	data, err := os.ReadFile(link)

	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "relative" {
		t.Fatalf("expected 'relative', got %q", string(data))
	}
}
