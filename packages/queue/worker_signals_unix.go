//go:build unix

package queue

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// installSignalHandlers wires the worker's process-level signal
// handling. SIGTERM/SIGQUIT/SIGINT cancel the run context for graceful
// shutdown. SIGUSR2 toggles the worker into the paused state (emitting
// WorkerPausing); SIGCONT lifts it (emitting WorkerResuming). The
// returned function stops the signal subscription and must be deferred
// by the caller.
//
// This file is Unix-only because SIGUSR2 and SIGCONT do not exist on
// Windows. A no-op implementation lives in worker_signals_other.go.
func (w *Worker) installSignalHandlers(ctx context.Context, cancel context.CancelFunc, queueName string) func() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)

	pause := make(chan os.Signal, 1)
	signal.Notify(pause, syscall.SIGUSR2)

	resume := make(chan os.Signal, 1)
	signal.Notify(resume, syscall.SIGCONT)

	connectionName := w.queue.ConnectionName()

	go func() {
		for {
			select {
			case <-shutdown:
				cancel()

				return
			case <-pause:
				if w.paused.CompareAndSwap(false, true) {
					w.emit(WorkerPausing{
						ConnectionName: connectionName,
						Queue:          queueName,
						WorkerName:     w.opts.Name,
					})
				}
			case <-resume:
				if w.paused.CompareAndSwap(true, false) {
					w.emit(WorkerResuming{
						ConnectionName: connectionName,
						Queue:          queueName,
						WorkerName:     w.opts.Name,
					})
				}
			case <-ctx.Done():
				return
			}
		}
	}()

	return func() {
		signal.Stop(shutdown)
		signal.Stop(pause)
		signal.Stop(resume)
	}
}
