package events

// WorkerStarting is dispatched once when the worker daemon boots, before
// Ref: @bedrock/code-0252
type WorkerStarting struct {
	ConnectionName string
	Queue          string
	WorkerName     string
}

// WorkerStopping is dispatched once when the worker daemon exits. Status
// carries the machine-readable stop reason; see WorkerStopReason in the
// top-level queue package for the enum values.
// Ref: @bedrock/code-0253
type WorkerStopping struct {
	Status     int
	WorkerName string
}

// WorkerPausing is dispatched when the worker daemon receives a pause
// signal (SIGUSR2 on Unix) and is about to suspend the run loop. The
// event carries the connection and queue the worker is bound to so
// listeners can correlate the pause with the affected backend.
// Ref: @bedrock/code-0250
type WorkerPausing struct {
	ConnectionName string
	Queue          string
	WorkerName     string
}

// WorkerResuming is dispatched when a paused worker daemon receives a
// resume signal (SIGCONT on Unix) and is about to re-enter the run
// Ref: @bedrock/code-0251
type WorkerResuming struct {
	ConnectionName string
	Queue          string
	WorkerName     string
}
