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

// WritingKey is dispatched before a value is written to the cache.
type WritingKey struct {
	StoreName string
	Key       string
	Value     any
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

func (CacheHit) CacheEvent()      {}
func (CacheMissed) CacheEvent()   {}
func (WritingKey) CacheEvent()    {}
func (KeyWritten) CacheEvent()    {}
func (ForgettingKey) CacheEvent() {}
func (KeyForgotten) CacheEvent()  {}
func (CacheFlushing) CacheEvent() {}
func (CacheFlushed) CacheEvent()  {}

var _ Event = CacheHit{}
var _ Event = CacheMissed{}
var _ Event = WritingKey{}
var _ Event = KeyWritten{}
var _ Event = ForgettingKey{}
var _ Event = KeyForgotten{}
var _ Event = CacheFlushing{}
var _ Event = CacheFlushed{}
