package session

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFileHandlerOpen(t *testing.T) {
	t.Parallel()

	h := NewFileHandler(t.TempDir(), 30)

	if err := h.Open(context.Background(), "/tmp", "sess"); err != nil {
		t.Fatalf("Open should be a no-op: %v", err)
	}
}

func TestFileHandlerClose(t *testing.T) {
	t.Parallel()

	h := NewFileHandler(t.TempDir(), 30)

	if err := h.Close(context.Background()); err != nil {
		t.Fatalf("Close should be a no-op: %v", err)
	}
}

func TestFileHandlerReadValid(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	h := NewFileHandler(dir, 30)
	ctx := context.Background()

	_ = os.WriteFile(filepath.Join(dir, "session-id"), []byte(`{"foo":"bar"}`), 0600)

	data, err := h.Read(ctx, "session-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data != `{"foo":"bar"}` {
		t.Fatalf("want stored data, got %q", data)
	}
}

func TestFileHandlerReadExpired(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	h := NewFileHandler(dir, 30)
	ctx := context.Background()

	name := filepath.Join(dir, "session-id")
	_ = os.WriteFile(name, []byte(`{"foo":"bar"}`), 0600)

	past := time.Now().Add(-31 * time.Minute)
	_ = os.Chtimes(name, past, past)

	data, err := h.Read(ctx, "session-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data != "" {
		t.Fatalf("expired session should return empty string, got %q", data)
	}
}

func TestFileHandlerReadMissing(t *testing.T) {
	t.Parallel()

	h := NewFileHandler(t.TempDir(), 30)

	data, err := h.Read(context.Background(), "nonexistent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data != "" {
		t.Fatalf("missing session should return empty string, got %q", data)
	}
}

func TestFileHandlerWrite(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	h := NewFileHandler(dir, 30)
	ctx := context.Background()

	if err := h.Write(ctx, "session-id", `{"key":"value"}`); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(dir, "session-id"))
	if err != nil {
		t.Fatalf("file should exist: %v", err)
	}

	if string(content) != `{"key":"value"}` {
		t.Fatalf("want written data, got %q", string(content))
	}
}

func TestFileHandlerDestroy(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	h := NewFileHandler(dir, 30)
	ctx := context.Background()

	_ = os.WriteFile(filepath.Join(dir, "session-id"), []byte("data"), 0600)

	if err := h.Destroy(ctx, "session-id"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "session-id")); !os.IsNotExist(err) {
		t.Fatal("session file should be deleted")
	}

	// Destroying a nonexistent session should not error.
	if err := h.Destroy(ctx, "nonexistent"); err != nil {
		t.Fatalf("destroying nonexistent should not error: %v", err)
	}
}

func TestFileHandlerGC(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	h := NewFileHandler(dir, 30)
	ctx := context.Background()

	// Create two session files.
	fresh := filepath.Join(dir, "fresh")
	stale := filepath.Join(dir, "stale")

	_ = os.WriteFile(fresh, []byte("data"), 0600)
	_ = os.WriteFile(stale, []byte("data"), 0600)

	// Backdate the stale file beyond the 5-second lifetime.
	past := time.Now().Add(-10 * time.Second)
	_ = os.Chtimes(stale, past, past)

	if err := h.GC(ctx, 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(fresh); os.IsNotExist(err) {
		t.Fatal("fresh session should still exist")
	}

	if _, err := os.Stat(stale); !os.IsNotExist(err) {
		t.Fatal("stale session should be deleted")
	}
}
