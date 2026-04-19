package reverb

import (
	"context"

	contractsReverb "github.com/bedrock/packages/contracts/reverb"
)

// SyncDispatcher is the default in-process event dispatcher.
// It routes events directly to channel subscribers without any external
// coordination. Use this when all Reverb server instances share the same
// process; for horizontal scaling across multiple processes use RedisDispatcher.
type SyncDispatcher struct {
	channels *ChannelManager
}

var _ contractsReverb.Dispatcher = (*SyncDispatcher)(nil)

// NewSyncDispatcher constructs a SyncDispatcher backed by channels.
func NewSyncDispatcher(channels *ChannelManager) *SyncDispatcher {
	return &SyncDispatcher{channels: channels}
}

// Dispatch broadcasts event to every subscriber of event.Channel within the
// given application. If event.Channel does not exist the call is a no-op,
// mirroring Laravel Reverb's behaviour of silently ignoring missing channels.
//
// Sender exclusion (socket_id) is handled at the server layer: callers that
// need to exclude a sender should call the channel's Broadcast method directly
// with a non-nil except pointer. The Dispatcher interface does not carry
// socket-ID semantics — it is used for HTTP API triggers and Redis fan-out,
// both of which broadcast to all subscribers.
func (d *SyncDispatcher) Dispatch(ctx context.Context, appID string, event contractsReverb.Event) error {
	ch, ok := d.channels.Get(appID, event.Channel)

	if !ok {
		return nil
	}

	return ch.BroadcastToAll(ctx, event)
}

// Subscribe is a no-op for the synchronous dispatcher; cross-server
// coordination is only needed for the Redis-backed implementation.
func (d *SyncDispatcher) Subscribe(_ context.Context, _ string) error {
	return nil
}
