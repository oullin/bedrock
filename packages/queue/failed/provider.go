package failed

import "time"

// FailedJob is the decoded record returned by All and Find. It is the
// Go analogue of the stdClass rows the upstream providers yield: the public
// fields mirror failed_jobs columns, plus an ID that maps to either the
// integer primary key (database) or the job UUID (uuid/file/dynamodb).
type FailedJob struct {
	ID         string
	UUID       string
	Connection string
	Queue      string
	Payload    string
	Exception  string
	FailedAt   time.Time
}

// Ref: @bedrock/code-0258
// persistence layer the worker uses to record permanently-failed jobs.
//
// Differences from PHP:
//
//   - IDs returns []string; callers coerce where they need ints.
//   - Find returns (nil, nil) when the id does not exist (upstream returns null).
//   - Flush takes `hours int` where 0 means "flush everything". upstream
//     accepts int|null with null meaning the same.
type FailedJobProvider interface {
	// Log persists a failed job record and returns its id (the uuid
	// for uuid-keyed providers, the integer for database providers
	// stringified). An error is returned only on storage failure.
	Log(connection, queue, payload string, exception error) (string, error)

	// IDs returns every failed job id; if queueFilter is non-empty the
	// list is limited to that queue.
	IDs(queueFilter string) ([]string, error)

	// All returns every failed job ordered newest-first.
	All() ([]FailedJob, error)

	// Find loads a single failed job by id, or (nil, nil) if missing.
	Find(id string) (*FailedJob, error)

	// Forget deletes a failed job by id. It returns true when a row
	// was actually removed.
	Forget(id string) (bool, error)

	// Flush deletes all (hours == 0) or all failed jobs older than the
	// given number of hours.
	Flush(hours int) error
}

// Providers that
// support efficient counting implement this optional interface.
// A nil or empty string acts as "any" for the corresponding dimension.
type Countable interface {
	Count(connection, queueFilter string) (int64, error)
}

// Implementations
// return the number of rows removed.
type Prunable interface {
	Prune(before time.Time) (int64, error)
}
