package filesystem_test

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestFiles(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "a.txt"), "a")
	writeFile(t, filepath.Join(dir, "b.txt"), "b")
	makeDir(t, filepath.Join(dir, "subdir"))

	files, err := fs.Files(dir)

	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}
}

func TestFilesExcludesHidden(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "a.txt"), "a")
	writeFile(t, filepath.Join(dir, ".hidden"), "h")

	files, err := fs.Files(dir)

	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 file (hidden excluded), got %d", len(files))
	}
}

func TestFilesIncludesHidden(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "a.txt"), "a")
	writeFile(t, filepath.Join(dir, ".hidden"), "h")

	files, err := fs.Files(dir, true)

	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 2 {
		t.Fatalf("expected 2 files (hidden included), got %d", len(files))
	}
}

func TestAllFiles(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "a.txt"), "a")
	writeFile(t, filepath.Join(dir, "sub", "b.txt"), "b")
	writeFile(t, filepath.Join(dir, "sub", "deep", "c.txt"), "c")

	files, err := fs.AllFiles(dir)

	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 3 {
		t.Fatalf("expected 3 files, got %d: %v", len(files), files)
	}
}

func TestAllFilesExcludesHidden(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "a.txt"), "a")
	writeFile(t, filepath.Join(dir, ".hidden", "b.txt"), "b")
	writeFile(t, filepath.Join(dir, ".secret"), "s")

	files, err := fs.AllFiles(dir)

	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 file (hidden excluded), got %d: %v", len(files), files)
	}
}

func TestAllFilesIncludesHidden(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "a.txt"), "a")
	writeFile(t, filepath.Join(dir, ".hidden", "b.txt"), "b")

	files, err := fs.AllFiles(dir, true)

	if err != nil {
		t.Fatal(err)
	}

	if len(files) != 2 {
		t.Fatalf("expected 2 files (hidden included), got %d: %v", len(files), files)
	}
}

func TestDirectories(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()

	makeDir(t, filepath.Join(dir, "sub1"))
	makeDir(t, filepath.Join(dir, "sub2"))
	writeFile(t, filepath.Join(dir, "file.txt"), "f")

	dirs, err := fs.Directories(dir)

	if err != nil {
		t.Fatal(err)
	}

	if len(dirs) != 2 {
		t.Fatalf("expected 2 directories, got %d", len(dirs))
	}
}

func TestAllDirectories(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()

	makeDir(t, filepath.Join(dir, "sub1"))
	makeDir(t, filepath.Join(dir, "sub1", "nested"))
	makeDir(t, filepath.Join(dir, "sub2"))

	dirs, err := fs.AllDirectories(dir)

	if err != nil {
		t.Fatal(err)
	}

	if len(dirs) != 3 {
		t.Fatalf("expected 3 directories, got %d: %v", len(dirs), dirs)
	}
}

func TestEnsureDirectoryExists(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "new", "deep", "dir")

	if err := fs.EnsureDirectoryExists(path); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(path)

	if err != nil {
		t.Fatal(err)
	}

	if !info.IsDir() {
		t.Fatal("expected directory to exist")
	}
}

func TestEnsureDirectoryExistsAlreadyExists(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()

	if err := fs.EnsureDirectoryExists(dir); err != nil {
		t.Fatal(err)
	}
}

func TestMakeDirectory(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "new_dir")

	if err := fs.MakeDirectory(path); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(path)

	if err != nil {
		t.Fatal(err)
	}

	if !info.IsDir() {
		t.Fatal("expected path to be a directory")
	}
}

func TestMakeDirectoryWithMode(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "mode_dir")

	if err := fs.MakeDirectory(path, 0o700); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(path)

	if err != nil {
		t.Fatal(err)
	}

	if info.Mode().Perm() != 0o700 {
		t.Fatalf("expected 0700, got %o", info.Mode().Perm())
	}
}

func TestDeleteDirectory(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "todelete")

	makeDir(t, path)
	writeFile(t, filepath.Join(path, "file.txt"), "content")
	makeDir(t, filepath.Join(path, "sub"))
	writeFile(t, filepath.Join(path, "sub", "nested.txt"), "nested")

	if err := fs.DeleteDirectory(path); err != nil {
		t.Fatal(err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatal("expected directory to be deleted")
	}
}

func TestDeleteDirectoryReturnsFalseWhenNotADirectory(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, path, "content")

	err := fs.DeleteDirectory(path)

	if err == nil {
		t.Fatal("expected error when path is a file")
	}
}

func TestDeleteDirectoryPreserve(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "preserved")

	makeDir(t, path)
	writeFile(t, filepath.Join(path, "file.txt"), "content")
	makeDir(t, filepath.Join(path, "sub"))

	if err := fs.DeleteDirectory(path, true); err != nil {
		t.Fatal(err)
	}

	// Directory should still exist.
	info, err := os.Stat(path)

	if err != nil {
		t.Fatal(err)
	}

	if !info.IsDir() {
		t.Fatal("expected directory to still exist")
	}

	// But contents should be gone.
	entries, err := os.ReadDir(path)

	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 0 {
		t.Fatalf("expected empty directory, got %d entries", len(entries))
	}
}

func TestDeleteDirectories(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()

	makeDir(t, filepath.Join(dir, "sub1"))
	makeDir(t, filepath.Join(dir, "sub2"))
	writeFile(t, filepath.Join(dir, "file.txt"), "keep")

	if err := fs.DeleteDirectories(dir); err != nil {
		t.Fatal(err)
	}

	entries, err := os.ReadDir(dir)

	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 1 {
		t.Fatalf("expected 1 entry (file.txt), got %d", len(entries))
	}

	if entries[0].Name() != "file.txt" {
		t.Fatalf("expected file.txt, got %s", entries[0].Name())
	}
}

func TestCleanDirectory(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "file.txt"), "content")
	makeDir(t, filepath.Join(dir, "sub"))
	writeFile(t, filepath.Join(dir, "sub", "nested.txt"), "nested")

	if err := fs.CleanDirectory(dir); err != nil {
		t.Fatal(err)
	}

	// Directory should still exist.
	info, err := os.Stat(dir)

	if err != nil {
		t.Fatal(err)
	}

	if !info.IsDir() {
		t.Fatal("expected directory to still exist")
	}

	// But contents should be empty.
	entries, err := os.ReadDir(dir)

	if err != nil {
		t.Fatal(err)
	}

	if len(entries) != 0 {
		t.Fatalf("expected empty directory, got %d entries", len(entries))
	}
}

func TestCopyDirectory(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	dst := filepath.Join(dir, "dst")

	writeFile(t, filepath.Join(src, "a.txt"), "a")
	writeFile(t, filepath.Join(src, "sub", "b.txt"), "b")

	if err := fs.CopyDirectory(src, dst); err != nil {
		t.Fatal(err)
	}

	// Verify copied files.
	data, err := os.ReadFile(filepath.Join(dst, "a.txt"))

	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "a" {
		t.Fatalf("expected 'a', got %q", string(data))
	}

	data, err = os.ReadFile(filepath.Join(dst, "sub", "b.txt"))

	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "b" {
		t.Fatalf("expected 'b', got %q", string(data))
	}

	// Source should still exist.
	if _, err := os.Stat(src); err != nil {
		t.Fatal("expected source to still exist")
	}
}

func TestCopyDirectoryNotADirectory(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")
	dst := filepath.Join(dir, "dst")

	writeFile(t, path, "content")

	err := fs.CopyDirectory(path, dst)

	if err == nil {
		t.Fatal("expected error when source is not a directory")
	}
}

func TestMoveDirectory(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	dst := filepath.Join(dir, "dst")

	writeFile(t, filepath.Join(src, "a.txt"), "a")
	writeFile(t, filepath.Join(src, "sub", "b.txt"), "b")

	if err := fs.MoveDirectory(src, dst); err != nil {
		t.Fatal(err)
	}

	// Source should be gone.
	if _, err := os.Stat(src); !os.IsNotExist(err) {
		t.Fatal("expected source to be gone")
	}

	// Destination should have the files.
	data, err := os.ReadFile(filepath.Join(dst, "a.txt"))

	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "a" {
		t.Fatalf("expected 'a', got %q", string(data))
	}
}

func TestMoveDirectoryOverwrite(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	src := filepath.Join(dir, "src")
	dst := filepath.Join(dir, "dst")

	writeFile(t, filepath.Join(src, "new.txt"), "new")
	writeFile(t, filepath.Join(dst, "old.txt"), "old")

	if err := fs.MoveDirectory(src, dst, true); err != nil {
		t.Fatal(err)
	}

	// Old file should be gone.
	if _, err := os.Stat(filepath.Join(dst, "old.txt")); !os.IsNotExist(err) {
		t.Fatal("expected old file to be gone")
	}

	data, err := os.ReadFile(filepath.Join(dst, "new.txt"))

	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "new" {
		t.Fatalf("expected 'new', got %q", string(data))
	}
}

func TestFilesReturnsFullPaths(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()

	writeFile(t, filepath.Join(dir, "a.txt"), "a")

	files, err := fs.Files(dir)

	if err != nil {
		t.Fatal(err)
	}

	expected := filepath.Join(dir, "a.txt")

	if len(files) != 1 || files[0] != expected {
		t.Fatalf("expected [%s], got %v", expected, files)
	}
}

func TestDirectoriesReturnsFullPaths(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()

	makeDir(t, filepath.Join(dir, "sub"))

	dirs, err := fs.Directories(dir)

	if err != nil {
		t.Fatal(err)
	}

	expected := filepath.Join(dir, "sub")

	if len(dirs) != 1 || dirs[0] != expected {
		t.Fatalf("expected [%s], got %v", expected, dirs)
	}
}

func TestAllDirectoriesSorted(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()

	makeDir(t, filepath.Join(dir, "c"))
	makeDir(t, filepath.Join(dir, "a"))
	makeDir(t, filepath.Join(dir, "b"))

	dirs, err := fs.AllDirectories(dir)

	if err != nil {
		t.Fatal(err)
	}

	sorted := make([]string, len(dirs))
	copy(sorted, dirs)

	sort.Strings(sorted)

	for i := range dirs {
		if dirs[i] != sorted[i] {
			t.Fatalf("expected sorted directories, got %v", dirs)
		}
	}
}
