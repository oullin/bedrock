package filesystem_test

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/bedrock/packages/filesystem"
)

func TestLockableFileCreateAndClose(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "lockable.txt")

	lf, err := filesystem.NewLockableFile(path, 0o644)

	if err != nil {
		t.Fatal(err)
	}

	if lf.Path() != path {
		t.Fatalf("expected path %q, got %q", path, lf.Path())
	}

	if err := lf.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestLockableFileWriteAndRead(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "lockable.txt")

	lf, err := filesystem.NewLockableFile(path, 0o644)

	if err != nil {
		t.Fatal(err)
	}

	defer lf.Close()

	n, err := lf.Write([]byte("hello world"))

	if err != nil {
		t.Fatal(err)
	}

	if n != 11 {
		t.Fatalf("expected 11 bytes written, got %d", n)
	}

	data, err := lf.Read()

	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "hello world" {
		t.Fatalf("expected 'hello world', got %q", string(data))
	}
}

func TestLockableFileReadPartial(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "lockable.txt")

	lf, err := filesystem.NewLockableFile(path, 0o644)

	if err != nil {
		t.Fatal(err)
	}

	defer lf.Close()

	lf.Write([]byte("hello world"))

	data, err := lf.Read(5)

	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "hello" {
		t.Fatalf("expected 'hello', got %q", string(data))
	}
}

func TestLockableFileTruncate(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "lockable.txt")

	lf, err := filesystem.NewLockableFile(path, 0o644)

	if err != nil {
		t.Fatal(err)
	}

	defer lf.Close()

	lf.Write([]byte("hello world"))

	if err := lf.Truncate(); err != nil {
		t.Fatal(err)
	}

	size, err := lf.Size()

	if err != nil {
		t.Fatal(err)
	}

	if size != 0 {
		t.Fatalf("expected size 0 after truncate, got %d", size)
	}
}

func TestLockableFileSharedLock(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "lockable.txt")

	lf, err := filesystem.NewLockableFile(path, 0o644)

	if err != nil {
		t.Fatal(err)
	}

	defer lf.Close()

	if err := lf.SharedLock(); err != nil {
		t.Fatal(err)
	}

	if err := lf.Unlock(); err != nil {
		t.Fatal(err)
	}
}

func TestLockableFileExclusiveLock(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "lockable.txt")

	lf, err := filesystem.NewLockableFile(path, 0o644)

	if err != nil {
		t.Fatal(err)
	}

	defer lf.Close()

	if err := lf.ExclusiveLock(); err != nil {
		t.Fatal(err)
	}

	// Write while locked.
	_, err = lf.Write([]byte("exclusive"))

	if err != nil {
		t.Fatal(err)
	}

	if err := lf.Unlock(); err != nil {
		t.Fatal(err)
	}

	data, err := lf.Read()

	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "exclusive" {
		t.Fatalf("expected 'exclusive', got %q", string(data))
	}
}

func TestLockableFileMultipleSharedLocks(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "shared.txt")

	writeFile(t, path, "shared content")

	lf1, err := filesystem.NewLockableFile(path, 0o644)

	if err != nil {
		t.Fatal(err)
	}

	defer lf1.Close()

	lf2, err := filesystem.NewLockableFile(path, 0o644)

	if err != nil {
		t.Fatal(err)
	}

	defer lf2.Close()

	// Both should be able to acquire shared locks.
	if err := lf1.SharedLock(); err != nil {
		t.Fatal(err)
	}

	if err := lf2.SharedLock(); err != nil {
		t.Fatal(err)
	}

	if err := lf1.Unlock(); err != nil {
		t.Fatal(err)
	}

	if err := lf2.Unlock(); err != nil {
		t.Fatal(err)
	}
}

func TestLockableFileExclusiveLockWaitsForSharedLockRelease(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "contended.txt")

	writeFile(t, path, "shared content")

	lf1, err := filesystem.NewLockableFile(path, 0o644)

	if err != nil {
		t.Fatal(err)
	}

	defer lf1.Close()

	lf2, err := filesystem.NewLockableFile(path, 0o644)

	if err != nil {
		t.Fatal(err)
	}

	defer lf2.Close()

	if err := lf1.SharedLock(); err != nil {
		t.Fatal(err)
	}

	locked := make(chan error, 1)

	go func() {
		locked <- lf2.ExclusiveLock()
	}()

	select {
	case err := <-locked:
		t.Fatalf("expected exclusive lock to block until shared lock is released, got %v", err)
	case <-time.After(100 * time.Millisecond):
	}

	if err := lf1.Unlock(); err != nil {
		t.Fatal(err)
	}

	select {
	case err := <-locked:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for exclusive lock after unlock")
	}

	if err := lf2.Unlock(); err != nil {
		t.Fatal(err)
	}
}

func TestLockableFileChmod(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "lockable.txt")

	lf, err := filesystem.NewLockableFile(path, 0o644)

	if err != nil {
		t.Fatal(err)
	}

	defer lf.Close()

	if err := lf.Chmod(0o755); err != nil {
		t.Fatal(err)
	}

	size, err := lf.Size()

	if err != nil {
		t.Fatal(err)
	}
	// New empty file.
	if size != 0 {
		t.Fatalf("expected size 0, got %d", size)
	}
}

func TestLockableFileCreatesParentDirectories(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "deep", "lockable.txt")

	lf, err := filesystem.NewLockableFile(path, 0o644)

	if err != nil {
		t.Fatal(err)
	}

	defer lf.Close()

	if lf.Path() != path {
		t.Fatalf("expected path %q, got %q", path, lf.Path())
	}
}

func TestLockableFileWriteOverwrite(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "lockable.txt")

	lf, err := filesystem.NewLockableFile(path, 0o644)

	if err != nil {
		t.Fatal(err)
	}

	defer lf.Close()

	lf.Write([]byte("first"))
	lf.Truncate()
	lf.Write([]byte("second"))

	data, err := lf.Read()

	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "second" {
		t.Fatalf("expected 'second', got %q", string(data))
	}
}
