package queue

import "github.com/bedrock/packages/queue/events"

// This file re-exports the event types defined in the events subpackage at
// the root of the queue package. The aliases exist for two reasons:
//
//  1. Callers (bus, worker, tests) can keep writing queue.JobProcessing
//     rather than events.JobProcessing — a shorter, more familiar name.
//  2. Worker emission sites stay compact.
//
// The subpackage is the source of truth; changing a field there changes
// it here. New Laravel events should be added to events/ first, then
// re-exported here only if callers need the unqualified name.
//
// The full set of Laravel 13.x Illuminate\Queue\Events\* types is
// re-exported below so a migration to qualified names is a pure find &
// replace.

// --- Job queueing (push side) -----------------------------------------

type JobQueueing = events.JobQueueing

type JobQueued = events.JobQueued

// --- Job processing (worker side) -------------------------------------

type JobPopping = events.JobPopping

type JobPopped = events.JobPopped

type JobProcessing = events.JobProcessing

type JobAttempted = events.JobAttempted

type JobProcessed = events.JobProcessed

type JobFailed = events.JobFailed

type JobExceptionOccurred = events.JobExceptionOccurred

type JobReleasedAfterException = events.JobReleasedAfterException

type JobTimedOut = events.JobTimedOut

type JobRetryRequested = events.JobRetryRequested

// --- Queue state ------------------------------------------------------

type Looping = events.Looping

type QueueBusy = events.QueueBusy

type QueuePaused = events.QueuePaused

type QueueResumed = events.QueueResumed

type QueueFailedOver = events.QueueFailedOver

// --- Worker lifecycle -------------------------------------------------

type WorkerStarting = events.WorkerStarting

type WorkerStopping = events.WorkerStopping

type WorkerPausing = events.WorkerPausing

type WorkerResuming = events.WorkerResuming
