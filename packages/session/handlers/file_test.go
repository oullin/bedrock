package handlers_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bedrock/packages/session/handlers"
)

func TestFileHandlerOpen(t *testing.T) {
	dir := t.TempDir()
	sub := filepath.Join(dir, "sessions")
	h := handlers.NewFileHandler(sub)

	if err := h.Open(context.Background(), "", "test"); err != nil {
		t.Fatalf("Open failed: %v", err)
	}

	info, err := os.Stat(sub)

	if err != nil {
		t.Fatalf("directory not created: %v", err)
	}

	if !info.IsDir() {
		t.Error("expected directory")
	}
}

func TestFileHandlerClose(t *testing.T) {
	h := handlers.NewFileHandler(t.TempDir())

	if err := h.Close(context.Background()); err != nil {
		t.Errorf("Close returned error: %v", err)
	}
}

func TestFileHandlerReadMissing(t *testing.T) {
	h := handlers.NewFileHandler(t.TempDir())
	ctx := context.Background()
	_ = h.Open(ctx, "", "test")

	data, err := h.Read(ctx, "nonexistent")

	if err != nil {
		t.Errorf("Read returned error: %v", err)
	}

	if data != "" {
		t.Errorf("expected empty string, got %q", data)
	}
}

func TestFileHandlerWriteAndRead(t *testing.T) {
	h := handlers.NewFileHandler(t.TempDir())
	ctx := context.Background()
	_ = h.Open(ctx, "", "test")

	payload := `{"foo":"bar"}`

	if err := h.Write(ctx, "sess1", payload); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	data, err := h.Read(ctx, "sess1")

	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if data != payload {
		t.Errorf("expected %q, got %q", payload, data)
	}
}

func TestFileHandlerDestroy(t *testing.T) {
	h := handlers.NewFileHandler(t.TempDir())
	ctx := context.Background()
	_ = h.Open(ctx, "", "test")

	_ = h.Write(ctx, "sess1", "data")

	if err := h.Destroy(ctx, "sess1"); err != nil {
		t.Fatalf("Destroy failed: %v", err)
	}

	data, _ := h.Read(ctx, "sess1")

	if data != "" {
		t.Error("session should be empty after Destroy")
	}
}

func TestFileHandlerDestroyNonExistent(t *testing.T) {
	h := handlers.NewFileHandler(t.TempDir())
	ctx := context.Background()
	_ = h.Open(ctx, "", "test")

	if err := h.Destroy(ctx, "nope"); err != nil {
		t.Errorf("Destroy non-existent should not error: %v", err)
	}
}

func TestFileHandlerGC(t *testing.T) {
	dir := t.TempDir()
	h := handlers.NewFileHandler(dir)
	ctx := context.Background()
	_ = h.Open(ctx, "", "test")

	// Write two sessions.
	_ = h.Write(ctx, "old", "data-old")
	_ = h.Write(ctx, "new", "data-new")

	// Set "old" file's mod time to the past.
	oldPath := filepath.Join(dir, "old")
	past := time.Now().Add(-2 * time.Hour)
	_ = os.Chtimes(oldPath, past, past)

	// GC with 1-hour max lifetime should remove "old".
	if err := h.GC(ctx, 3600); err != nil {
		t.Fatalf("GC failed: %v", err)
	}

	data, _ := h.Read(ctx, "old")

	if data != "" {
		t.Error("old session should be removed by GC")
	}

	data, _ = h.Read(ctx, "new")

	if data != "data-new" {
		t.Error("new session should survive GC")
	}
}

func TestFileHandlerOverwrite(t *testing.T) {
	h := handlers.NewFileHandler(t.TempDir())
	ctx := context.Background()
	_ = h.Open(ctx, "", "test")

	_ = h.Write(ctx, "sess1", "first")
	_ = h.Write(ctx, "sess1", "second")

	data, _ := h.Read(ctx, "sess1")

	if data != "second" {
		t.Errorf("expected second, got %q", data)
	}
}
