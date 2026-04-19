package websockets

import (
	"context"
	"encoding/json"
	"fmt"

	contractsWebSockets "github.com/bedrock/packages/contracts/websockets"
	"github.com/bedrock/packages/redis"
)

// RedisDispatcher distributes events across multiple WebSockets server instances
// using Redis pub/sub. Each event is published to a key named:
//
//	{prefix}:{appID}:{channelName}
//
// All server instances subscribe to patterns for their applications and
// locally broadcast received events to their connected clients.
type RedisDispatcher struct {
	channels *ChannelManager
	redis    *redis.Manager
	prefix   string
}

var _ contractsWebSockets.Dispatcher = (*RedisDispatcher)(nil)

// NewRedisDispatcher constructs a RedisDispatcher that uses the given Redis
// manager for pub/sub and prefixes all channel keys with prefix.
func NewRedisDispatcher(channels *ChannelManager, mgr *redis.Manager, prefix string) *RedisDispatcher {
	return &RedisDispatcher{
		channels: channels,
		redis:    mgr,
		prefix:   prefix,
	}
}

// redisKey returns the Redis pub/sub key for an app + channel combination.
func (d *RedisDispatcher) redisKey(appID, channel string) string {
	return fmt.Sprintf("%s:%s:%s", d.prefix, appID, channel)
}

// Dispatch serialises event to JSON and publishes it to the Redis pub/sub
// key for the event's channel. All server instances subscribed to the pattern
// for appID will receive the message and broadcast it locally.
func (d *RedisDispatcher) Dispatch(ctx context.Context, appID string, event contractsWebSockets.Event) error {
	conn, err := d.redis.Connection("")

	if err != nil {
		return err
	}

	payload, err := json.Marshal(event)

	if err != nil {
		return err
	}

	_, err = conn.Command(ctx, "PUBLISH", d.redisKey(appID, event.Channel), string(payload))

	return err
}

// Subscribe starts a background goroutine that listens on the Redis pattern
// "{prefix}:{appID}:*" and locally broadcasts any received events to the
// matching channel's subscribers. It returns immediately; the goroutine runs
// until ctx is cancelled.
func (d *RedisDispatcher) Subscribe(ctx context.Context, appID string) error {
	conn, err := d.redis.Connection("")

	if err != nil {
		return err
	}

	pattern := fmt.Sprintf("%s:%s:*", d.prefix, appID)

	go func() {
		_ = conn.PSubscribe(ctx, []string{pattern}, func(_, _, payload string) {
			var event contractsWebSockets.Event

			if err := json.Unmarshal([]byte(payload), &event); err != nil {
				return
			}

			ch, ok := d.channels.Get(appID, event.Channel)

			if !ok {
				return
			}

			_ = ch.BroadcastToAll(ctx, event)
		})
	}()

	return nil
}
