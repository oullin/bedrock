package broadcasting

import (
	"context"
	"encoding/json"
	"fmt"
)

// RedisPublisher is the minimal Redis publishing backend.
type RedisPublisher interface {
	Publish(channel string, payload []byte) error
}

// RedisBroadcaster implements Redis pub/sub broadcasting.
type RedisBroadcaster struct {
	*BaseBroadcaster
	publisher RedisPublisher
	prefix    string
}

// NewRedisBroadcaster creates a RedisBroadcaster.
func NewRedisBroadcaster(publisher RedisPublisher, prefix string) *RedisBroadcaster {
	return &RedisBroadcaster{BaseBroadcaster: NewBaseBroadcaster(), publisher: publisher, prefix: prefix}
}

// Auth authenticates a private or presence channel request.
func (b *RedisBroadcaster) Auth(request AuthRequest) (any, error) {
	channel := NormalizeChannelName(trimPrefix(request.ChannelName, b.prefix))

	if request.ChannelName == "" || (IsGuardedChannel(request.ChannelName) && b.RetrieveUser(request, channel) == nil) {
		return nil, ErrAccessDenied
	}

	result, err := b.VerifyUserCanAccessChannel(request, channel)

	if err != nil {
		return nil, err
	}

	return b.ValidAuthenticationResponse(request, result)
}

// ValidAuthenticationResponse returns Redis broadcaster auth JSON.
func (b *RedisBroadcaster) ValidAuthenticationResponse(request AuthRequest, result any) (any, error) {
	if _, ok := result.(bool); ok {
		data, err := json.Marshal(result)

		return string(data), err
	}

	channel := NormalizeChannelName(trimPrefix(request.ChannelName, b.prefix))
	user := b.RetrieveUser(request, channel)

	data, err := json.Marshal(map[string]any{
		"channel_data": map[string]any{
			"user_id":   broadcastingIdentifier(user),
			"user_info": result,
		},
	})

	if err != nil {
		return nil, err
	}

	return string(data), nil
}

// Broadcast publishes the event payload to each Redis channel.
func (b *RedisBroadcaster) Broadcast(_ context.Context, channels []string, event string, payload map[string]any) error {
	if len(channels) == 0 {
		return nil
	}

	nextPayload, _ := payloadWithoutSocket(payload)
	socket := payload["socket"]

	data, err := json.Marshal(map[string]any{
		"event":  event,
		"data":   nextPayload,
		"socket": socket,
	})

	if err != nil {
		return err
	}

	for _, channel := range channels {
		if err := b.publisher.Publish(b.prefix+channel, data); err != nil {
			return fmt.Errorf("%w: %v", ErrBroadcast, err)
		}
	}

	return nil
}

func trimPrefix(value, prefix string) string {
	if prefix == "" {
		return value
	}

	if len(value) >= len(prefix) && value[:len(prefix)] == prefix {
		return value[len(prefix):]
	}

	return value
}
