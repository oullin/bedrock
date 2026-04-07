package session

import (
	"context"
	"sync"
	"testing"
)

func TestArrayHandlerReadWrite(t *testing.T) {
	t.Parallel()

	h := NewArrayHandler()
	ctx := context.Background()

	_ = h.Write(ctx, "id1", `{"key":"value"}`)

	data, err := h.Read(ctx, "id1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data != `{"key":"value"}` {
		t.Fatalf("want stored data, got %q", data)
	}
}

func TestArrayHandlerReadMissing(t *testing.T) {
	t.Parallel()

	h := NewArrayHandler()

	data, err := h.Read(context.Background(), "missing")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data != "" {
		t.Fatalf("want empty string for missing session, got %q", data)
	}
}

func TestArrayHandlerDestroy(t *testing.T) {
	t.Parallel()

	h := NewArrayHandler()
	ctx := context.Background()

	_ = h.Write(ctx, "id1", "data")
	_ = h.Destroy(ctx, "id1")

	data, _ := h.Read(ctx, "id1")
	if data != "" {
		t.Fatal("destroyed session should return empty string")
	}
}

func TestArrayHandlerOpenClose(t *testing.T) {
	t.Parallel()

	h := NewArrayHandler()
	ctx := context.Background()

	if err := h.Open(ctx, "/tmp", "sess"); err != nil {
		t.Fatalf("Open should be a no-op: %v", err)
	}

	if err := h.Close(ctx); err != nil {
		t.Fatalf("Close should be a no-op: %v", err)
	}
}

func TestArrayHandlerGC(t *testing.T) {
	t.Parallel()

	h := NewArrayHandler()

	if err := h.GC(context.Background(), 3600); err != nil {
		t.Fatalf("GC should be a no-op: %v", err)
	}
}

func TestArrayHandlerConcurrent(t *testing.T) {
	t.Parallel()

	h := NewArrayHandler()
	ctx := context.Background()

	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			_ = h.Write(ctx, "shared", "data")
		}(i)
	}

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = h.Read(ctx, "shared")
		}()
	}

	wg.Wait()
}
