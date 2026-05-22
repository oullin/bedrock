package events

// JobPopping is dispatched before the worker asks the driver for the next
// Ref: @bedrock/code-0237
type JobPopping struct {
	ConnectionName string
}

// JobPopped is dispatched after the worker has retrieved a job from the
// driver but before it begins processing.
// Ref: @bedrock/code-0236
type JobPopped struct {
	ConnectionName string
	Job            any
}

// JobProcessing is dispatched immediately before a job is fired.
// Ref: @bedrock/code-0239
type JobProcessing struct {
	ConnectionName string
	Job            any
}

// JobAttempted is dispatched after every attempt to run a job, regardless
// of whether it succeeded or threw.
// Ref: @bedrock/code-0233
type JobAttempted struct {
	ConnectionName string
	Job            any
}

// JobProcessed is dispatched after a job has been processed successfully.
// Ref: @bedrock/code-0238
type JobProcessed struct {
	ConnectionName string
	Job            any
}

// JobFailed is dispatched when a job has exhausted its retry budget and
// is being marked as permanently failed.
// Ref: @bedrock/code-0235
type JobFailed struct {
	ConnectionName string
	Job            any
	Err            error
}

// JobExceptionOccurred is dispatched every time a job attempt throws an
// exception, regardless of whether the job will be retried or failed.
// Ref: @bedrock/code-0234
type JobExceptionOccurred struct {
	ConnectionName string
	Job            any
	Err            error
}

// JobReleasedAfterException is dispatched when the worker releases a job
// back onto the queue following a caught exception.
// Ref: @bedrock/code-0242
type JobReleasedAfterException struct {
	ConnectionName string
	Job            any
}

// JobTimedOut is dispatched when a job exceeds its configured timeout.
// Ref: @bedrock/code-0244
type JobTimedOut struct {
	ConnectionName string
	Job            any
}

// JobRetryRequested is dispatched when an operator re-queues a failed job
// via the queue:retry command.
// Ref: @bedrock/code-0243
type JobRetryRequested struct {
	// Payload is the failed-job record being retried.
	Payload any
}
