package filesystem_test

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/bedrock/packages/filesystem"
)

func TestGet(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, path, "hello world")

	data, err := fs.Get(path)

	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "hello world" {
		t.Fatalf("expected 'hello world', got %q", string(data))
	}
}

func TestGetNonexistentFile(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()

	_, err := fs.Get(filepath.Join(dir, "nonexistent.txt"))

	if !errors.Is(err, filesystem.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestJSON(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "data.json")

	writeFile(t, path, `{"name":"john","age":30}`)

	var result map[string]any

	if err := fs.JSON(path, &result); err != nil {
		t.Fatal(err)
	}

	if result["name"] != "john" {
		t.Fatalf("unexpected name: %v", result["name"])
	}

	if result["age"] != float64(30) {
		t.Fatalf("unexpected age: %v", result["age"])
	}
}

func TestJSONInvalid(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "invalid.json")

	writeFile(t, path, "not valid json{{{")

	var result map[string]any

	if err := fs.JSON(path, &result); err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestJSONNonexistentFile(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()

	var result map[string]any
	err := fs.JSON(filepath.Join(dir, "missing.json"), &result)

	if !errors.Is(err, filesystem.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestSharedGet(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	writeFile(t, path, "shared content")

	data, err := fs.SharedGet(path)

	if err != nil {
		t.Fatal(err)
	}

	if string(data) != "shared content" {
		t.Fatalf("expected 'shared content', got %q", string(data))
	}
}

func TestSharedGetNonexistentFile(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()

	_, err := fs.SharedGet(filepath.Join(dir, "missing.txt"))

	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestLines(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "lines.txt")

	writeFile(t, path, "line1\nline2\nline3")

	seq, err := fs.Lines(path)

	if err != nil {
		t.Fatal(err)
	}

	var lines []string

	for line := range seq {
		lines = append(lines, line)
	}

	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}

	if lines[0] != "line1" || lines[1] != "line2" || lines[2] != "line3" {
		t.Fatalf("unexpected lines: %v", lines)
	}
}

func TestLinesNonexistentFile(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()

	_, err := fs.Lines(filepath.Join(dir, "missing.txt"))

	if !errors.Is(err, filesystem.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestLinesEmptyFile(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "empty.txt")

	writeFile(t, path, "")

	seq, err := fs.Lines(path)

	if err != nil {
		t.Fatal(err)
	}

	var lines []string

	for line := range seq {
		lines = append(lines, line)
	}

	if len(lines) != 0 {
		t.Fatalf("expected 0 lines, got %d", len(lines))
	}
}

func TestLinesBreakEarly(t *testing.T) {
	t.Parallel()

	fs := newFilesystem()
	dir := t.TempDir()
	path := filepath.Join(dir, "lines.txt")

	writeFile(t, path, "line1\nline2\nline3\nline4\nline5")

	seq, err := fs.Lines(path)

	if err != nil {
		t.Fatal(err)
	}

	var lines []string

	for line := range seq {
		lines = append(lines, line)

		if len(lines) == 2 {
			break
		}
	}

	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}
