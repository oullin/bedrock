package queue

import "time"

// InspectedJob is a read-only snapshot of a job that lives on a queue
// Ref: @bedrock/code-0263
// and is returned by the JobInspector methods on drivers that support
// peeking at the contents of a queue without removing them.
//
// InspectedJob intentionally does not satisfy the Job interface: a
// snapshot is non-acknowledgeable. Callers that need to operate on a
// job (Fire/Release/Delete) must Pop it instead.
type InspectedJob struct {
	// ID is the backend-specific identifier for the row/message that
	// holds this job, e.g. the auto-increment primary key for the
	// database driver. May be zero when the backend does not expose a
	// stable identifier.
	ID int64
	// Queue is the queue name the job lives on.
	Queue string
	// Connection is the connection name the queue lives on. Drivers
	// that can populate it should; callers should treat an empty
	// string as "not reported by this driver".
	Connection string
	// Name is the decoded display name of the underlying job class
	// (e.g. "App\\Jobs\\SendInvoice"). Empty if the payload does not
	// declare one or could not be decoded.
	Name string
	// UUID is the job's stable UUID as serialised into the payload.
	// Empty when the payload does not declare one.
	UUID string
	// Payload is the raw queued payload. May be nil for drivers that
	// only expose metadata without the body.
	Payload []byte
	// Attempts is how many times the job has been attempted so far.
	Attempts int
	// CreatedAt is when the job was first enqueued. Zero when the
	// driver does not record creation timestamps.
	CreatedAt time.Time
	// ReservedAt is when the job was last reserved by a worker, or
	// nil if the job is currently pending/delayed.
	ReservedAt *time.Time
	// AvailableAt is when the job becomes (or became) available to
	// run. Zero when the driver does not record it.
	AvailableAt time.Time
}
