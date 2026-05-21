package events

// WorkerStarting is dispatched once when the worker daemon boots, before
// it enters the main loop. Mirrors @bedrock\Queue\Events\WorkerStarting.
type WorkerStarting struct {
	ConnectionName string
	Queue          string
	WorkerName     string
}

// WorkerStopping is dispatched once when the worker daemon exits. Status
// carries the machine-readable stop reason; see WorkerStopReason in the
// top-level queue package for the enum values.
// Mirrors @bedrock\Queue\Events\WorkerStopping.
type WorkerStopping struct {
	Status     int
	WorkerName string
}

// WorkerPausing is dispatched when the worker daemon receives a pause
// signal (SIGUSR2 on Unix) and is about to suspend the run loop. The
// event carries the connection and queue the worker is bound to so
// listeners can correlate the pause with the affected backend.
// Mirrors upstream @bedrock\Queue\Events\WorkerPausing.
type WorkerPausing struct {
	ConnectionName string
	Queue          string
	WorkerName     string
}

// WorkerResuming is dispatched when a paused worker daemon receives a
// resume signal (SIGCONT on Unix) and is about to re-enter the run
// loop. Mirrors upstream @bedrock\Queue\Events\WorkerResuming.
type WorkerResuming struct {
	ConnectionName string
	Queue          string
	WorkerName     string
}
