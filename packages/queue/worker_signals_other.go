//go:build !unix

package queue

import (
	"context"
	"os"
	"os/signal"
)

// installSignalHandlers on non-Unix platforms only wires graceful
// shutdown via os.Interrupt. SIGUSR2 / SIGCONT do not exist outside of
// Unix, so the pause/resume signal path is a no-op. Callers can still
// drive WorkerPausing / WorkerResuming through external means (e.g. by
// flipping the paused flag through their own control channel) without
// the runtime panicking on undefined syscalls.
func (w *Worker) installSignalHandlers(ctx context.Context, cancel context.CancelFunc, _ string) func() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt)

	go func() {
		select {
		case <-shutdown:
			cancel()
		case <-ctx.Done():
		}
	}()

	return func() {
		signal.Stop(shutdown)
	}
}
