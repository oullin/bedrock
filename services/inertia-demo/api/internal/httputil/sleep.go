package httputil

import (
	"context"
	"time"
)

// SleepCtx blocks for duration d or until ctx is cancelled,
// whichever comes first. Returns ctx.Err() on cancellation.
func SleepCtx(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)

	defer t.Stop()

	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
