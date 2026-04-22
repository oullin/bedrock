package cache

import (
	"context"
	"time"
)

// Event is the marker interface for all cache events.
type Event interface {
	CacheEvent()
}

// EventDispatcher dispatches cache events to registered listeners.
type EventDispatcher interface {
	Dispatch(ctx context.Context, event Event)
}

// CacheHit is dispatched when a cache key is found.
type CacheHit struct {
	StoreName string
	Key       string
	Value     any
	Tags      []string
}

// CacheMissed is dispatched when a cache key is not found.
type CacheMissed struct {
	StoreName string
	Key       string
	Tags      []string
}

// RetrievingKey is dispatched before a single cache key is read.
type RetrievingKey struct {
	StoreName string
	Key       string
	Tags      []string
}

// RetrievingManyKeys is dispatched before multiple cache keys are read.
type RetrievingManyKeys struct {
	StoreName string
	Keys      []string
	Tags      []string
}

// WritingKey is dispatched before a value is written to the cache.
type WritingKey struct {
	StoreName string
	Key       string
	Value     any
	TTL       time.Duration
	Tags      []string
}

// WritingManyKeys is dispatched before multiple cache values are written.
type WritingManyKeys struct {
	StoreName string
	Keys      []string
	Values    map[string]any
	TTL       time.Duration
	Tags      []string
}

// KeyWritten is dispatched after a value is written to the cache.
type KeyWritten struct {
	StoreName string
	Key       string
	Value     any
	TTL       time.Duration
	Tags      []string
}

// ForgettingKey is dispatched before a key is removed from the cache.
type ForgettingKey struct {
	StoreName string
	Key       string
	Tags      []string
}

// KeyForgotten is dispatched after a key is removed from the cache.
type KeyForgotten struct {
	StoreName string
	Key       string
	Tags      []string
}

// KeyForgetFailed is dispatched when a cache key removal fails.
type KeyForgetFailed struct {
	StoreName string
	Key       string
	Err       error
	Tags      []string
}

// CacheFlushing is dispatched before a cache store is flushed.
type CacheFlushing struct {
	StoreName string
	Tags      []string
}

// CacheFlushed is dispatched after a cache store is flushed.
type CacheFlushed struct {
	StoreName string
	Tags      []string
}

// CacheFlushFailed is dispatched when a cache store flush fails.
type CacheFlushFailed struct {
	StoreName string
	Err       error
	Tags      []string
}

// CacheLocksFlushing is dispatched before cache locks are flushed.
type CacheLocksFlushing struct {
	StoreName string
	Tags      []string
}

// CacheLocksFlushed is dispatched after cache locks are flushed.
type CacheLocksFlushed struct {
	StoreName string
	Tags      []string
}

// CacheLocksFlushFailed is dispatched when cache lock flushing fails.
type CacheLocksFlushFailed struct {
	StoreName string
	Err       error
	Tags      []string
}

func (CacheHit) CacheEvent()              {}
func (CacheMissed) CacheEvent()           {}
func (RetrievingKey) CacheEvent()         {}
func (RetrievingManyKeys) CacheEvent()    {}
func (WritingKey) CacheEvent()            {}
func (WritingManyKeys) CacheEvent()       {}
func (KeyWritten) CacheEvent()            {}
func (ForgettingKey) CacheEvent()         {}
func (KeyForgotten) CacheEvent()          {}
func (KeyForgetFailed) CacheEvent()       {}
func (CacheFlushing) CacheEvent()         {}
func (CacheFlushed) CacheEvent()          {}
func (CacheFlushFailed) CacheEvent()      {}
func (CacheLocksFlushing) CacheEvent()    {}
func (CacheLocksFlushed) CacheEvent()     {}
func (CacheLocksFlushFailed) CacheEvent() {}

var _ Event = CacheHit{}
var _ Event = CacheMissed{}
var _ Event = RetrievingKey{}
var _ Event = RetrievingManyKeys{}
var _ Event = WritingKey{}
var _ Event = WritingManyKeys{}
var _ Event = KeyWritten{}
var _ Event = ForgettingKey{}
var _ Event = KeyForgotten{}
var _ Event = KeyForgetFailed{}
var _ Event = CacheFlushing{}
var _ Event = CacheFlushed{}
var _ Event = CacheFlushFailed{}
var _ Event = CacheLocksFlushing{}
var _ Event = CacheLocksFlushed{}
var _ Event = CacheLocksFlushFailed{}
