package queue

// JobProcessing is emitted before a job is processed.
type JobProcessing struct {
	ConnectionName string
	Job            Job
}

// JobProcessed is emitted after a job is successfully processed.
type JobProcessed struct {
	ConnectionName string
	Job            Job
}

// JobFailed is emitted when a job fails all its attempts.
type JobFailed struct {
	ConnectionName string
	Job            Job
	Err            error
}

// JobAttempted is emitted each time a job is attempted.
type JobAttempted struct {
	ConnectionName string
	Job            Job
}

// JobExceptionOccurred is emitted when a job throws an exception during processing.
type JobExceptionOccurred struct {
	ConnectionName string
	Job            Job
	Err            error
}
