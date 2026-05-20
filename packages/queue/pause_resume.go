package queue

import (
	"time"

	"github.com/bedrock/packages/queue/events"
)

// PauseResumer coordinates the pause/resume lifecycle across a
// PauseStore and an EventEmitter. It is the shared primitive that the
// QueueManager (Step 6) will embed to expose pause/pauseFor/resume/
// isPaused at the application level, and that the test port consumes
// directly.
//
// Mirrors the pause/resume methods on the upstream QueueManager.
type PauseResumer struct {
	store   PauseStore
	emitter EventEmitter
}

// NewPauseResumer constructs a PauseResumer. The emitter may be nil —
// events are simply dropped when no emitter is wired up.
func NewPauseResumer(store PauseStore, emitter EventEmitter) *PauseResumer {
	return &PauseResumer{store: store, emitter: emitter}
}

// Pause marks (connection, queue) as paused indefinitely and emits a
// QueuePaused event with a nil TTL.
func (p *PauseResumer) Pause(connection, queue string) error {
	if err := p.store.Pause(pauseKey(connection, queue)); err != nil {
		return err
	}

	p.emit(events.QueuePaused{ConnectionName: connection, Queue: queue, TTL: nil})

	return nil
}

// PauseFor marks (connection, queue) as paused until now+ttl and emits
// a QueuePaused event whose TTL is non-nil and equal to ttl.
func (p *PauseResumer) PauseFor(connection, queue string, ttl time.Duration) error {
	if err := p.store.PauseFor(pauseKey(connection, queue), ttl); err != nil {
		return err
	}

	d := ttl

	p.emit(events.QueuePaused{ConnectionName: connection, Queue: queue, TTL: &d})

	return nil
}

// Resume removes any pause state for (connection, queue) and emits a
// QueueResumed event.
func (p *PauseResumer) Resume(connection, queue string) error {
	if err := p.store.Resume(pauseKey(connection, queue)); err != nil {
		return err
	}

	p.emit(events.QueueResumed{ConnectionName: connection, Queue: queue})

	return nil
}

// IsPaused reports whether (connection, queue) is currently paused.
// Any underlying store error collapses to false — callers that care
// about the distinction can query the store directly.
func (p *PauseResumer) IsPaused(connection, queue string) bool {
	paused, err := p.store.IsPaused(pauseKey(connection, queue))

	if err != nil {
		return false
	}

	return paused
}

// emit sends event through the configured emitter, if any.
func (p *PauseResumer) emit(event any) {
	if p.emitter == nil {
		return
	}

	p.emitter.Emit(event)
}

// pauseKey builds the opaque storage key used by PauseStore.
func pauseKey(connection, queue string) string {
	return connection + ":" + queue
}
