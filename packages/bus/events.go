package bus

// EventFunc is a callback invoked when batch lifecycle events occur.
type EventFunc func(event any)

// BatchDispatched is fired after a batch is dispatched.
type BatchDispatched struct{ Batch *Batch }

// BatchStarted is fired when the first batch job is recorded.
type BatchStarted struct{ Batch *Batch }

// BatchFinished is fired when all batch jobs have completed successfully.
type BatchFinished struct{ Batch *Batch }

// BatchCanceled is fired when a batch is cancelled.
type BatchCanceled struct{ Batch *Batch }
