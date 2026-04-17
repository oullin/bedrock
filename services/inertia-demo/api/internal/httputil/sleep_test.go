package httputil

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestSleepCtx_Completes(t *testing.T) {
	t.Parallel()

	start := time.Now()

	if err := SleepCtx(context.Background(), 5*time.Millisecond); err != nil {
		t.Fatalf("SleepCtx() error = %v", err)
	}

	if elapsed := time.Since(start); elapsed < 5*time.Millisecond {
		t.Fatalf("SleepCtx() elapsed = %v, want at least %v", elapsed, 5*time.Millisecond)
	}
}

func TestSleepCtx_Cancelled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())

	cancel()

	err := SleepCtx(ctx, time.Second)

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("SleepCtx() error = %v, want %v", err, context.Canceled)
	}
}
