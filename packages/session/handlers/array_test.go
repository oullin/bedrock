package handlers_test

import (
	"context"
	"sync"
	"testing"

	"github.com/bedrock/packages/session/handlers"
)

func TestArrayHandlerOpen(t *testing.T) {
	h := handlers.NewArrayHandler()

	if err := h.Open(context.Background(), "/tmp", "test"); err != nil {
		t.Errorf("Open returned error: %v", err)
	}
}

func TestArrayHandlerClose(t *testing.T) {
	h := handlers.NewArrayHandler()

	if err := h.Close(context.Background()); err != nil {
		t.Errorf("Close returned error: %v", err)
	}
}

func TestArrayHandlerReadMissing(t *testing.T) {
	h := handlers.NewArrayHandler()
	ctx := context.Background()

	data, err := h.Read(ctx, "nonexistent")

	if err != nil {
		t.Errorf("Read returned error: %v", err)
	}

	if data != "" {
		t.Errorf("expected empty string, got %q", data)
	}
}

func TestArrayHandlerWriteAndRead(t *testing.T) {
	h := handlers.NewArrayHandler()
	ctx := context.Background()

	if err := h.Write(ctx, "sess1", `{"foo":"bar"}`); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	data, err := h.Read(ctx, "sess1")

	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	if data != `{"foo":"bar"}` {
		t.Errorf("expected {\"foo\":\"bar\"}, got %q", data)
	}
}

func TestArrayHandlerDestroy(t *testing.T) {
	h := handlers.NewArrayHandler()
	ctx := context.Background()

	_ = h.Write(ctx, "sess1", "data")

	if err := h.Destroy(ctx, "sess1"); err != nil {
		t.Fatalf("Destroy failed: %v", err)
	}

	data, _ := h.Read(ctx, "sess1")

	if data != "" {
		t.Error("session should be empty after Destroy")
	}
}

func TestArrayHandlerGC(t *testing.T) {
	h := handlers.NewArrayHandler()

	if err := h.GC(context.Background(), 3600); err != nil {
		t.Errorf("GC returned error: %v", err)
	}
}

func TestArrayHandlerOverwrite(t *testing.T) {
	h := handlers.NewArrayHandler()
	ctx := context.Background()

	_ = h.Write(ctx, "sess1", "first")
	_ = h.Write(ctx, "sess1", "second")

	data, _ := h.Read(ctx, "sess1")

	if data != "second" {
		t.Errorf("expected second, got %q", data)
	}
}

func TestArrayHandlerMultipleSessions(t *testing.T) {
	h := handlers.NewArrayHandler()
	ctx := context.Background()

	_ = h.Write(ctx, "a", "data-a")
	_ = h.Write(ctx, "b", "data-b")

	da, _ := h.Read(ctx, "a")
	db, _ := h.Read(ctx, "b")

	if da != "data-a" {
		t.Errorf("expected data-a, got %q", da)
	}

	if db != "data-b" {
		t.Errorf("expected data-b, got %q", db)
	}
}

func TestArrayHandlerDestroyNonExistent(t *testing.T) {
	h := handlers.NewArrayHandler()

	if err := h.Destroy(context.Background(), "nope"); err != nil {
		t.Errorf("Destroy non-existent should not error: %v", err)
	}
}

func TestArrayHandlerConcurrency(t *testing.T) {
	h := handlers.NewArrayHandler()
	ctx := context.Background()

	var wg sync.WaitGroup

	for i := range 50 {
		wg.Add(1)

		go func(n int) {
			defer wg.Done()

			id := "sess"

			_ = h.Write(ctx, id, "data")
			_, _ = h.Read(ctx, id)

			if n%3 == 0 {
				_ = h.Destroy(ctx, id)
			}
		}(i)
	}

	wg.Wait()
}
