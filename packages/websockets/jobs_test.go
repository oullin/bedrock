package reverb_test

import (
	"context"
	"testing"
	"time"

	"github.com/bedrock/packages/websockets"
)

func TestStartJobLoop_StopsOnContextCancel(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	apps := websockets.NewAppManager(nil)
	conns := websockets.NewConnectionManager()

	done := make(chan struct{})
	go func() {
		websockets.StartJobLoop(ctx, conns, apps, 10*time.Millisecond)
		close(done)
	}()

	cancel()

	select {
	case <-done:
		// pass
	case <-time.After(time.Second):
		t.Fatal("StartJobLoop did not stop after context cancellation")
	}
}
