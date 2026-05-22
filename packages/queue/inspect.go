package queue

import (
	"context"
	"errors"
)

// ErrNotSupported is returned by drivers that cannot satisfy a queue
// inspection or enumeration request. Backends like SQS cannot peek at
// reserved or delayed messages non-destructively, and Redis does not
// track reserved jobs in its default layout — those calls surface this
// sentinel so callers can branch on it with errors.Is.
//
// LogicException("Not supported") for the same scenarios.

// QueueNamer is the optional contract a Queue implementation satisfies
// when it can enumerate the queue names that currently exist on its
// connection. The manager-level cross-queue inspection helpers
// (AllPendingJobs / AllDelayedJobs / AllReservedJobs) require this
// contract so they know which named queues to fan out to.
//
// Drivers that cannot enumerate their queues (sync, null) need not
// implement this interface — the manager treats a missing implementation
// as "no queues to inspect" and returns an empty result, not an error.
//
// added alongside allReservedJobs/allDelayedJobs/allPendingJobs in
// Upstream 13.8.0.
type QueueNamer interface {
	QueueNames(ctx context.Context) ([]string, error)
}

// JobInspector is the optional contract a Queue implementation
// satisfies when it can return read-only snapshots of jobs currently
// sitting on a named queue without consuming them. The three methods
// mirror the per-state inspection methods upstream exposes on the queue
// contract: pending (unreserved, ready), delayed (unreserved, not yet
// due), and reserved (in-flight).
//
// Drivers that cannot peek at one or more states should still implement
// the interface and return ErrNotSupported for the unsupported state
// rather than dropping the method entirely — that keeps the
// manager-level fan-out predictable and lets callers distinguish "I
// asked, the driver can't" from "no rows".
type JobInspector interface {
	PendingJobs(ctx context.Context, queue string) ([]InspectedJob, error)
	DelayedJobs(ctx context.Context, queue string) ([]InspectedJob, error)
	ReservedJobs(ctx context.Context, queue string) ([]InspectedJob, error)
}

var ErrNotSupported = errors.New("queue: operation not supported by this driver")
