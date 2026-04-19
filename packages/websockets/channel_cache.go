package websockets

import (
	"context"
	"encoding/json"
	"sync"

	contractsWebSockets "github.com/bedrock/packages/contracts/websockets"
)

// ─────────────────────────────────────────────
// cacheStore – thread-safe last-event store
// ─────────────────────────────────────────────

type cacheStore struct {
	mu    sync.RWMutex
	event *contractsWebSockets.Event
}

// Set stores ev as the most recently broadcast event.

// Get returns the most recently broadcast event, or nil if none.

// sendCacheMiss sends the pusher:cache_miss event to conn.

// ─────────────────────────────────────────────
// CacheChannel – public channel with caching
// ─────────────────────────────────────────────

// CacheChannel extends Channel with last-event caching.
// New subscribers receive the most recently broadcast event, or
// pusher:cache_miss if no event has been broadcast yet.
type CacheChannel struct {
	Channel
	cache cacheStore
}

// NewCacheChannel constructs a CacheChannel.

// Subscribe subscribes conn and delivers the last cached event (or cache_miss).

// Broadcast caches the event and then broadcasts it.

// BroadcastToAll caches the event and then broadcasts it to all subscribers.

// LastEvent returns the most recently broadcast event, or nil if none.

// CacheEvent stores ev as the last broadcast event without broadcasting.

// ─────────────────────────────────────────────
// PrivateCacheChannel
// ─────────────────────────────────────────────

// PrivateCacheChannel extends PrivateChannel with last-event caching.
type PrivateCacheChannel struct {
	PrivateChannel
	cache cacheStore
}

// NewPrivateCacheChannel constructs a PrivateCacheChannel.

// Subscribe verifies auth, subscribes conn, and delivers the cached event (or cache_miss).

// Broadcast caches the event and then broadcasts it.

// BroadcastToAll caches the event and then broadcasts it to all subscribers.

// LastEvent returns the most recently broadcast event, or nil if none.

// CacheEvent stores ev as the last broadcast event without broadcasting.

// ─────────────────────────────────────────────
// PresenceCacheChannel
// ─────────────────────────────────────────────

// PresenceCacheChannel extends PresenceChannel with last-event caching.
type PresenceCacheChannel struct {
	PresenceChannel
	cache cacheStore
}

func (s *cacheStore) Set(ev contractsWebSockets.Event) {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.event = &ev
}

func (s *cacheStore) Get() *contractsWebSockets.Event {
	s.mu.RLock()

	defer s.mu.RUnlock()

	return s.event
}

func sendCacheMiss(ctx context.Context, conn contractsWebSockets.Connection, name string) error {
	b, err := json.Marshal(map[string]string{
		"event":   "pusher:cache_miss",
		"data":    "{}",
		"channel": name,
	})

	if err != nil {
		return err
	}

	return conn.Send(ctx, b)
}

var _ contractsWebSockets.Channel = (*CacheChannel)(nil)
var _ contractsWebSockets.CacheableChannel = (*CacheChannel)(nil)

func NewCacheChannel(name string, app *App) *CacheChannel {
	return &CacheChannel{
		Channel: Channel{
			channel: channel{
				name:  name,
				app:   app,
				conns: make(map[string]contractsWebSockets.Connection),
			},
		},
	}
}

func (ch *CacheChannel) Subscribe(ctx context.Context, conn contractsWebSockets.Connection, auth, data string) error {
	if err := ch.Channel.Subscribe(ctx, conn, auth, data); err != nil {
		return err
	}

	if ev := ch.cache.Get(); ev != nil {
		b, err := json.Marshal(ev)

		if err != nil {
			return err
		}

		return conn.Send(ctx, b)
	}

	return sendCacheMiss(ctx, conn, ch.name)
}

func (ch *CacheChannel) Broadcast(ctx context.Context, event contractsWebSockets.Event, except *string) error {
	ch.cache.Set(event)

	return ch.channel.broadcast(ctx, event, except)
}

func (ch *CacheChannel) BroadcastToAll(ctx context.Context, event contractsWebSockets.Event) error {
	ch.cache.Set(event)

	return ch.channel.broadcastAll(ctx, event)
}

func (ch *CacheChannel) LastEvent() *contractsWebSockets.Event {
	return ch.cache.Get()
}

func (ch *CacheChannel) CacheEvent(ev contractsWebSockets.Event) {
	ch.cache.Set(ev)
}

var _ contractsWebSockets.Channel = (*PrivateCacheChannel)(nil)
var _ contractsWebSockets.CacheableChannel = (*PrivateCacheChannel)(nil)

func NewPrivateCacheChannel(name string, app *App) *PrivateCacheChannel {
	return &PrivateCacheChannel{
		PrivateChannel: PrivateChannel{
			channel: channel{
				name:  name,
				app:   app,
				conns: make(map[string]contractsWebSockets.Connection),
			},
		},
	}
}

func (ch *PrivateCacheChannel) Subscribe(ctx context.Context, conn contractsWebSockets.Connection, auth, data string) error {
	if err := ch.PrivateChannel.Subscribe(ctx, conn, auth, data); err != nil {
		return err
	}

	if ev := ch.cache.Get(); ev != nil {
		b, err := json.Marshal(ev)

		if err != nil {
			return err
		}

		return conn.Send(ctx, b)
	}

	return sendCacheMiss(ctx, conn, ch.name)
}

func (ch *PrivateCacheChannel) Broadcast(ctx context.Context, event contractsWebSockets.Event, except *string) error {
	ch.cache.Set(event)

	return ch.channel.broadcast(ctx, event, except)
}

func (ch *PrivateCacheChannel) BroadcastToAll(ctx context.Context, event contractsWebSockets.Event) error {
	ch.cache.Set(event)

	return ch.channel.broadcastAll(ctx, event)
}

func (ch *PrivateCacheChannel) LastEvent() *contractsWebSockets.Event {
	return ch.cache.Get()
}

func (ch *PrivateCacheChannel) CacheEvent(ev contractsWebSockets.Event) {
	ch.cache.Set(ev)
}

var _ contractsWebSockets.Channel = (*PresenceCacheChannel)(nil)
var _ contractsWebSockets.CacheableChannel = (*PresenceCacheChannel)(nil)
var _ contractsWebSockets.PresenceChanneler = (*PresenceCacheChannel)(nil)

// NewPresenceCacheChannel constructs a PresenceCacheChannel.
func NewPresenceCacheChannel(name string, app *App) *PresenceCacheChannel {
	return &PresenceCacheChannel{
		PresenceChannel: PresenceChannel{
			channel: channel{
				name:  name,
				app:   app,
				conns: make(map[string]contractsWebSockets.Connection),
			},
			app:     app,
			members: make(map[string]presenceMember),
		},
	}
}

// Subscribe verifies auth, subscribes conn, and delivers the cached event (or cache_miss).
func (ch *PresenceCacheChannel) Subscribe(ctx context.Context, conn contractsWebSockets.Connection, auth, data string) error {
	if err := ch.PresenceChannel.Subscribe(ctx, conn, auth, data); err != nil {
		return err
	}

	if ev := ch.cache.Get(); ev != nil {
		b, err := json.Marshal(ev)

		if err != nil {
			return err
		}

		return conn.Send(ctx, b)
	}

	return sendCacheMiss(ctx, conn, ch.name)
}

// Broadcast caches the event and then broadcasts it.
func (ch *PresenceCacheChannel) Broadcast(ctx context.Context, event contractsWebSockets.Event, except *string) error {
	ch.cache.Set(event)

	return ch.channel.broadcast(ctx, event, except)
}

// BroadcastToAll caches the event and then broadcasts it to all subscribers.
func (ch *PresenceCacheChannel) BroadcastToAll(ctx context.Context, event contractsWebSockets.Event) error {
	ch.cache.Set(event)

	return ch.channel.broadcastAll(ctx, event)
}

// LastEvent returns the most recently broadcast event, or nil if none.
func (ch *PresenceCacheChannel) LastEvent() *contractsWebSockets.Event {
	return ch.cache.Get()
}

// CacheEvent stores ev as the last broadcast event without broadcasting.
func (ch *PresenceCacheChannel) CacheEvent(ev contractsWebSockets.Event) {
	ch.cache.Set(ev)
}
