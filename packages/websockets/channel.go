package websockets

import (
	"context"
	"encoding/json"
	"strings"
	"sync"

	contractsWebSockets "github.com/bedrock/packages/contracts/websockets"
)

// TypeOf returns the channel type string for the given channel name, using
// prefix matching in priority order.

// channel is the unexported base struct shared by all channel types.
type channel struct {
	name  string
	app   *App
	mu    sync.RWMutex
	conns map[string]contractsWebSockets.Connection // socketID → Connection
}

// Name returns the channel name.

// Connections returns a snapshot of all subscribed connections.

// HasConnection reports whether the given socket ID is subscribed.

// broadcast sends event to all connections except the one identified by except.
// Pass nil to broadcast to everyone.

// broadcastAll sends event to every subscriber without exclusions.

// sendSubscriptionSucceeded sends the subscription_succeeded event to conn.

// ─────────────────────────────────────────────
// Channel – public channel (no auth required)
// ─────────────────────────────────────────────

// Channel is a public (unauthenticated) Pusher channel.
type Channel struct {
	channel
}

// Subscribe adds conn to the channel and sends subscription_succeeded.

// Unsubscribe removes conn from the channel.

// Broadcast sends the event to all subscribers except the sender.

// BroadcastToAll sends the event to every subscriber.

// ─────────────────────────────────────────────
// PrivateChannel – HMAC-authenticated channel
// ─────────────────────────────────────────────

// PrivateChannel is a private Pusher channel that requires HMAC auth.
type PrivateChannel struct {
	channel
}

func TypeOf(name string) string {
	switch {
	case strings.HasPrefix(name, "private-cache-"):
		return "private-cache"
	case strings.HasPrefix(name, "presence-cache-"):
		return "presence-cache"
	case strings.HasPrefix(name, "cache-"):
		return "cache"
	case strings.HasPrefix(name, "private-"):
		return "private"
	case strings.HasPrefix(name, "presence-"):
		return "presence"
	default:
		return "public"
	}
}

func (c *channel) Name() string {
	return c.name
}

func (c *channel) Connections() []contractsWebSockets.Connection {
	c.mu.RLock()

	defer c.mu.RUnlock()

	out := make([]contractsWebSockets.Connection, 0, len(c.conns))

	for _, conn := range c.conns {
		out = append(out, conn)
	}

	return out
}

func (c *channel) HasConnection(socketID string) bool {
	c.mu.RLock()

	defer c.mu.RUnlock()

	_, ok := c.conns[socketID]

	return ok
}

func (c *channel) broadcast(ctx context.Context, event contractsWebSockets.Event, except *string) error {
	b, err := json.Marshal(event)

	if err != nil {
		return err
	}

	c.mu.RLock()

	defer c.mu.RUnlock()

	for id, conn := range c.conns {
		if except != nil && id == *except {
			continue
		}

		if err := conn.Send(ctx, b); err != nil {
			return err
		}
	}

	return nil
}

func (c *channel) broadcastAll(ctx context.Context, event contractsWebSockets.Event) error {
	return c.broadcast(ctx, event, nil)
}

func sendSubscriptionSucceeded(ctx context.Context, conn contractsWebSockets.Connection, name string) error {
	b, err := json.Marshal(map[string]string{
		"event":   "pusher_internal:subscription_succeeded",
		"data":    "{}",
		"channel": name,
	})

	if err != nil {
		return err
	}

	return conn.Send(ctx, b)
}

var _ contractsWebSockets.Channel = (*Channel)(nil)

func (ch *Channel) Subscribe(ctx context.Context, conn contractsWebSockets.Connection, auth, data string) error {
	ch.mu.Lock()
	ch.conns[conn.SocketID()] = conn
	ch.mu.Unlock()

	return sendSubscriptionSucceeded(ctx, conn, ch.name)
}

func (ch *Channel) Unsubscribe(_ context.Context, conn contractsWebSockets.Connection) {
	ch.mu.Lock()

	defer ch.mu.Unlock()

	delete(ch.conns, conn.SocketID())
}

func (ch *Channel) Broadcast(ctx context.Context, event contractsWebSockets.Event, except *string) error {
	return ch.broadcast(ctx, event, except)
}

func (ch *Channel) BroadcastToAll(ctx context.Context, event contractsWebSockets.Event) error {
	return ch.broadcastAll(ctx, event)
}

var _ contractsWebSockets.Channel = (*PrivateChannel)(nil)

// Subscribe verifies the auth token and, on success, adds conn to the channel.
func (ch *PrivateChannel) Subscribe(ctx context.Context, conn contractsWebSockets.Connection, auth, data string) error {
	if !VerifyChannelAuth(ch.app.Secret(), conn.SocketID(), ch.name, "", auth) {
		return ErrUnauthorized
	}

	ch.mu.Lock()
	ch.conns[conn.SocketID()] = conn
	ch.mu.Unlock()

	return sendSubscriptionSucceeded(ctx, conn, ch.name)
}

// Unsubscribe removes conn from the channel.
func (ch *PrivateChannel) Unsubscribe(_ context.Context, conn contractsWebSockets.Connection) {
	ch.mu.Lock()

	defer ch.mu.Unlock()

	delete(ch.conns, conn.SocketID())
}

// Broadcast sends the event to all subscribers except the sender.
func (ch *PrivateChannel) Broadcast(ctx context.Context, event contractsWebSockets.Event, except *string) error {
	return ch.broadcast(ctx, event, except)
}

// BroadcastToAll sends the event to every subscriber.
func (ch *PrivateChannel) BroadcastToAll(ctx context.Context, event contractsWebSockets.Event) error {
	return ch.broadcastAll(ctx, event)
}

// ─────────────────────────────────────────────
// NewChannel factory
// ─────────────────────────────────────────────

// NewChannel returns the appropriate Channel implementation for name.
func NewChannel(name string, app *App) contractsWebSockets.Channel {
	switch {
	case strings.HasPrefix(name, "private-cache-"):
		return NewPrivateCacheChannel(name, app)
	case strings.HasPrefix(name, "presence-cache-"):
		return NewPresenceCacheChannel(name, app)
	case strings.HasPrefix(name, "cache-"):
		return NewCacheChannel(name, app)
	case strings.HasPrefix(name, "private-"):
		return NewPrivateChannel(name, app)
	case strings.HasPrefix(name, "presence-"):
		return NewPresenceChannel(name, app)
	default:
		return &Channel{
			channel: channel{
				name:  name,
				app:   app,
				conns: make(map[string]contractsWebSockets.Connection),
			},
		}
	}
}

// NewPrivateChannel constructs a PrivateChannel.
func NewPrivateChannel(name string, app *App) *PrivateChannel {
	return &PrivateChannel{
		channel: channel{
			name:  name,
			app:   app,
			conns: make(map[string]contractsWebSockets.Connection),
		},
	}
}
