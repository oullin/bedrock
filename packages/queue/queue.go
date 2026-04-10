package queue

import (
	"context"
	"time"
)

// Queue defines the interface for a queue backend.
type Queue interface {
	// Push adds a job payload to the queue.
	Push(ctx context.Context, queue string, payload []byte) (string, error)
	// PushDelayed adds a delayed job to the queue.
	PushDelayed(ctx context.Context, queue string, payload []byte, delay time.Duration) (string, error)
	// PushMultiple adds multiple payloads to the queue.
	PushMultiple(ctx context.Context, queue string, payloads [][]byte) ([]string, error)
	// Pop retrieves and reserves the next available job.
	Pop(ctx context.Context, queue string) (Job, error)
	// Size returns the number of jobs ready for processing.
	Size(ctx context.Context, queue string) (int64, error)
	// PendingSize returns the total number of pending (not reserved) jobs.
	PendingSize(ctx context.Context, queue string) (int64, error)
	// DelayedSize returns the number of delayed jobs.
	DelayedSize(ctx context.Context, queue string) (int64, error)
	// ReservedSize returns the number of reserved (in-flight) jobs.
	ReservedSize(ctx context.Context, queue string) (int64, error)
	// ConnectionName returns the connection name for this queue.
	ConnectionName() string
}

// Connector creates a Queue from a configuration map.
type Connector interface {
	Connect(config map[string]any) (Queue, error)
}
