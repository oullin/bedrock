package events

import "time"

// Looping is dispatched once per iteration of the worker's main loop,
// before it attempts to pop a job. Mirrors Illuminate\Queue\Events\Looping.
type Looping struct {
	ConnectionName string
	Queue          string
}

// QueueBusy is dispatched by the queue:monitor command when a queue has
// more pending work than its configured threshold.
// Mirrors Illuminate\Queue\Events\QueueBusy.
type QueueBusy struct {
	ConnectionName string
	Queue          string
	Size           int64
}

// QueuePaused is dispatched when an operator pauses a queue via the
// queue:pause command. Mirrors Illuminate\Queue\Events\QueuePaused.
//
// the upstream event carries an integer `ttl` (seconds, null for indefinite).
// The Go port uses a *time.Duration so callers can inspect both the
// nullability and the magnitude without unit conversion. A nil TTL means
// the pause is indefinite.
type QueuePaused struct {
	ConnectionName string
	Queue          string
	TTL            *time.Duration
}

// QueueResumed is dispatched when an operator resumes a paused queue via
// the queue:resume command. Mirrors Illuminate\Queue\Events\QueueResumed.
type QueueResumed struct {
	ConnectionName string
	Queue          string
}

// QueueFailedOver is dispatched by the failover driver when it abandons
// one backend and switches to the next. Mirrors
// Illuminate\Queue\Events\QueueFailedOver.
type QueueFailedOver struct {
	From string
	To   string
	Err  error
}
