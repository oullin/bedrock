package events

// JobPopping is dispatched before the worker asks the driver for the next
// available job. Mirrors Illuminate\Queue\Events\JobPopping.
type JobPopping struct {
	ConnectionName string
}

// JobPopped is dispatched after the worker has retrieved a job from the
// driver but before it begins processing.
// Mirrors Illuminate\Queue\Events\JobPopped.
type JobPopped struct {
	ConnectionName string
	Job            any
}

// JobProcessing is dispatched immediately before a job is fired.
// Mirrors Illuminate\Queue\Events\JobProcessing.
type JobProcessing struct {
	ConnectionName string
	Job            any
}

// JobAttempted is dispatched after every attempt to run a job, regardless
// of whether it succeeded or threw.
// Mirrors Illuminate\Queue\Events\JobAttempted.
type JobAttempted struct {
	ConnectionName string
	Job            any
}

// JobProcessed is dispatched after a job has been processed successfully.
// Mirrors Illuminate\Queue\Events\JobProcessed.
type JobProcessed struct {
	ConnectionName string
	Job            any
}

// JobFailed is dispatched when a job has exhausted its retry budget and
// is being marked as permanently failed.
// Mirrors Illuminate\Queue\Events\JobFailed.
type JobFailed struct {
	ConnectionName string
	Job            any
	Err            error
}

// JobExceptionOccurred is dispatched every time a job attempt throws an
// exception, regardless of whether the job will be retried or failed.
// Mirrors Illuminate\Queue\Events\JobExceptionOccurred.
type JobExceptionOccurred struct {
	ConnectionName string
	Job            any
	Err            error
}

// JobReleasedAfterException is dispatched when the worker releases a job
// back onto the queue following a caught exception.
// Mirrors Illuminate\Queue\Events\JobReleasedAfterException.
type JobReleasedAfterException struct {
	ConnectionName string
	Job            any
}

// JobTimedOut is dispatched when a job exceeds its configured timeout.
// Mirrors Illuminate\Queue\Events\JobTimedOut.
type JobTimedOut struct {
	ConnectionName string
	Job            any
}

// JobRetryRequested is dispatched when an operator re-queues a failed job
// via the queue:retry command.
// Mirrors Illuminate\Queue\Events\JobRetryRequested.
type JobRetryRequested struct {
	// Payload is the failed-job record being retried.
	Payload any
}
