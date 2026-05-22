package events

import "time"

// JobQueueing is dispatched immediately before a job is pushed onto the
// Ref: @bedrock/code-0241
type JobQueueing struct {
	ConnectionName string
	Queue          string
	// Job is the original job value handed to the dispatcher — can be a
	// struct, a string class name, or whatever the caller pushed.
	Job     any
	Payload string
	Delay   time.Duration
}

// JobQueued is dispatched immediately after a job has been successfully
// Ref: @bedrock/code-0240
type JobQueued struct {
	ConnectionName string
	Queue          string
	// ID is the backend-specific identifier returned by the driver.
	ID      any
	Job     any
	Payload string
}
