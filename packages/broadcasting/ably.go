package broadcasting

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// AblyPublisher is the minimal Ably publishing backend.
type AblyPublisher interface {
	Publish(channel string, message AblyMessage) error
}

// AblyMessage is the message sent to Ably channels.
type AblyMessage struct {
	Name          string
	Data          map[string]any
	ConnectionKey any
}

// AblyBroadcaster implements Ably-compatible broadcasting.
type AblyBroadcaster struct {
	*BaseBroadcaster
	key       string
	publisher AblyPublisher
}

// NewAblyBroadcaster creates an AblyBroadcaster.
func NewAblyBroadcaster(key string, publisher AblyPublisher) *AblyBroadcaster {
	return &AblyBroadcaster{BaseBroadcaster: NewBaseBroadcaster(), key: key, publisher: publisher}
}

// Auth authenticates a private or presence channel request.
func (b *AblyBroadcaster) Auth(request AuthRequest) (any, error) {
	channel := NormalizeChannelName(request.ChannelName)

	if request.ChannelName == "" || (IsGuardedChannel(request.ChannelName) && b.RetrieveUser(request, channel) == nil) {
		return nil, ErrAccessDenied
	}

	result, err := b.VerifyUserCanAccessChannel(request, channel)

	if err != nil {
		return nil, err
	}

	return b.ValidAuthenticationResponse(request, result)
}

// ValidAuthenticationResponse returns an Ably-compatible auth payload.
func (b *AblyBroadcaster) ValidAuthenticationResponse(request AuthRequest, result any) (any, error) {
	if IsPrivateChannel(request.ChannelName) {
		return map[string]any{
			"auth": b.publicToken() + ":" + b.GenerateSignature(request.ChannelName, request.SocketID, nil),
		}, nil
	}

	channel := NormalizeChannelName(request.ChannelName)
	user := b.RetrieveUser(request, channel)
	userData := map[string]any{
		"user_id":   broadcastingIdentifier(user),
		"user_info": result,
	}

	encoded, err := json.Marshal(userData)

	if err != nil {
		return nil, err
	}

	return map[string]any{
		"auth":         b.publicToken() + ":" + b.GenerateSignature(request.ChannelName, request.SocketID, userData),
		"channel_data": string(encoded),
	}, nil
}

// GenerateSignature computes Ably's channel auth signature.
func (b *AblyBroadcaster) GenerateSignature(channelName, socketID string, userData map[string]any) string {
	message := socketID + ":" + channelName

	if userData != nil {
		encoded, _ := json.Marshal(userData)
		message += ":" + string(encoded)
	}

	return hmacSHA256(b.privateToken(), message)
}

// Broadcast publishes an Ably message to each formatted channel.
func (b *AblyBroadcaster) Broadcast(_ context.Context, channels []string, event string, payload map[string]any) error {
	for _, channel := range channels {
		message := AblyMessage{Name: event, Data: clonePayload(payload), ConnectionKey: payload["socket"]}

		if err := b.publisher.Publish(formatAblyChannel(channel), message); err != nil {
			return fmt.Errorf("%w: %v", ErrBroadcast, err)
		}
	}

	return nil
}

func (b *AblyBroadcaster) publicToken() string {
	token, _, _ := strings.Cut(b.key, ":")

	return token
}

func (b *AblyBroadcaster) privateToken() string {
	_, token, _ := strings.Cut(b.key, ":")

	return token
}

func formatAblyChannel(channel string) string {
	switch {
	case strings.HasPrefix(channel, "private-"):
		return "private:" + strings.TrimPrefix(channel, "private-")
	case strings.HasPrefix(channel, "presence-"):
		return "presence:" + strings.TrimPrefix(channel, "presence-")
	default:
		return "public:" + channel
	}
}
