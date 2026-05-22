package broadcasting

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// PusherSettings are the credentials needed for Pusher-compatible signing.
type PusherSettings struct {
	AuthKey string
	Secret  string
}

// PusherClient is the minimal Pusher-compatible backend used by PusherBroadcaster.
type PusherClient interface {
	SocketAuth(channel, socketID string) (string, error)
	PresenceAuth(channel, socketID, userID string, userInfo any) (string, error)
	Trigger(channels []string, event string, payload map[string]any, params map[string]string) error
	Settings() PusherSettings
}

// PusherUserAuthenticator can provide native Pusher user-auth responses.
type PusherUserAuthenticator interface {
	AuthenticateUser(socketID string, user any) (map[string]any, error)
}

// PusherBroadcaster implements Pusher/WebSockets-compatible broadcasting.
type PusherBroadcaster struct {
	*BaseBroadcaster
	client PusherClient
}

// NewPusherBroadcaster creates a PusherBroadcaster.
func NewPusherBroadcaster(client PusherClient) *PusherBroadcaster {
	return &PusherBroadcaster{BaseBroadcaster: NewBaseBroadcaster(), client: client}
}

// Auth authenticates a private or presence channel request.
func (b *PusherBroadcaster) Auth(request AuthRequest) (any, error) {
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

// ValidAuthenticationResponse returns a Pusher-compatible auth payload.
func (b *PusherBroadcaster) ValidAuthenticationResponse(request AuthRequest, result any) (any, error) {
	if IsPrivateChannel(request.ChannelName) {
		response, err := b.client.SocketAuth(request.ChannelName, request.SocketID)

		if err != nil {
			return nil, err
		}

		return decodeJSONMap(response)
	}

	channel := NormalizeChannelName(request.ChannelName)
	user := b.RetrieveUser(request, channel)
	userID := broadcastingIdentifier(user)

	response, err := b.client.PresenceAuth(request.ChannelName, request.SocketID, userID, result)

	if err != nil {
		return nil, err
	}

	return decodeJSONMap(response)
}

// ResolveAuthenticatedUser returns a Pusher-compatible user auth response.
func (b *PusherBroadcaster) ResolveAuthenticatedUser(request AuthRequest) (map[string]any, error) {
	user := b.BaseBroadcaster.ResolveAuthenticatedUser(request)

	if user == nil {
		return nil, nil
	}

	if native, ok := b.client.(PusherUserAuthenticator); ok {
		return native.AuthenticateUser(request.SocketID, user)
	}

	encoded, err := json.Marshal(user)

	if err != nil {
		return nil, err
	}

	settings := b.client.Settings()
	message := fmt.Sprintf("%s::user::%s", request.SocketID, string(encoded))
	signature := hmacSHA256(settings.Secret, message)

	return map[string]any{
		"auth":      settings.AuthKey + ":" + signature,
		"user_data": string(encoded),
	}, nil
}

// Broadcast triggers an event on the configured Pusher-compatible backend.
func (b *PusherBroadcaster) Broadcast(_ context.Context, channels []string, event string, payload map[string]any) error {
	if len(channels) == 0 {
		return nil
	}

	nextPayload, params := payloadWithoutSocket(payload)

	for start := 0; start < len(channels); start += 100 {
		end := start + 100

		if end > len(channels) {
			end = len(channels)
		}

		if err := b.client.Trigger(formatChannels(channels[start:end]), event, nextPayload, params); err != nil {
			return fmt.Errorf("%w: %v", ErrBroadcast, err)
		}
	}

	return nil
}

// IsPrivateChannel reports whether channel uses private channel auth.
func IsPrivateChannel(channel string) bool {
	return stringsHasAnyPrefix(channel, "private-", "private-encrypted-")
}

func decodeJSONMap(value string) (map[string]any, error) {
	var decoded map[string]any

	if err := json.Unmarshal([]byte(value), &decoded); err != nil {
		return nil, err
	}

	return decoded, nil
}

func hmacSHA256(secret, message string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write([]byte(message))

	return hex.EncodeToString(mac.Sum(nil))
}

func stringsHasAnyPrefix(value string, prefixes ...string) bool {
	for _, prefix := range prefixes {
		if len(value) >= len(prefix) && value[:len(prefix)] == prefix {
			return true
		}
	}

	return false
}
