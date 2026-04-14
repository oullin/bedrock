package events

// WorkerStarting is dispatched once when the worker daemon boots, before
// it enters the main loop. Mirrors Framework\Queue\Events\WorkerStarting.
type WorkerStarting struct {
	ConnectionName string
	Queue          string
	WorkerName     string
}

// WorkerStopping is dispatched once when the worker daemon exits. Status
// carries the machine-readable stop reason; see WorkerStopReason in the
// top-level queue package for the enum values.
// Mirrors Framework\Queue\Events\WorkerStopping.
type WorkerStopping struct {
	Status     int
	WorkerName string
}
