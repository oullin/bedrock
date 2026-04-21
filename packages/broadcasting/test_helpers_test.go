package broadcasting_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/bedrock/packages/broadcasting"
	contractsauth "github.com/bedrock/packages/contracts/auth"
)

type testUser struct {
	id          string
	broadcastID string
}

type userResolver struct {
	users map[string]contractsauth.Authenticatable
	calls []string
}

type fakePusherClient struct {
	socketAuthCalls   int
	presenceAuthCalls int
	triggers          []pusherTrigger
	settings          broadcasting.PusherSettings
}

type pusherTrigger struct {
	channels []string
	event    string
	payload  map[string]any
	params   map[string]string
}

type fakeRedisPublisher struct {
	messages map[string][]byte
}

type fakeAblyPublisher struct {
	messages map[string]broadcasting.AblyMessage
}

type recordingBroadcaster struct {
	calls []broadcastCall
}

type broadcastCall struct {
	channels []string
	event    string
	payload  map[string]any
}

var _ contractsauth.BroadcastingAuthenticatable = (*testUser)(nil)

func (u *testUser) GetAuthIdentifierName() string            { return "id" }
func (u *testUser) GetAuthIdentifier() string                { return u.id }
func (u *testUser) GetAuthPasswordName() string              { return "password" }
func (u *testUser) GetAuthPassword() string                  { return "" }
func (u *testUser) SetAuthPassword(_ string)                 {}
func (u *testUser) GetRememberToken() string                 { return "" }
func (u *testUser) SetRememberToken(_ string)                {}
func (u *testUser) GetRememberTokenName() string             { return "remember_token" }
func (u *testUser) GetAuthIdentifierForBroadcasting() string { return u.broadcastID }

func (r *userResolver) User(guard string) contractsauth.Authenticatable {
	r.calls = append(r.calls, guard)

	return r.users[guard]
}

func requestWithUser(channel string) broadcasting.AuthRequest {
	user := &testUser{id: "42", broadcastID: "42"}

	return broadcasting.AuthRequest{
		ChannelName:  channel,
		SocketID:     "abcd.1234",
		UserResolver: &userResolver{users: map[string]contractsauth.Authenticatable{"": user}},
	}
}

func requestWithoutUser(channel string) broadcasting.AuthRequest {
	return broadcasting.AuthRequest{
		ChannelName:  channel,
		SocketID:     "abcd.1234",
		UserResolver: &userResolver{users: map[string]contractsauth.Authenticatable{}},
	}
}

func allowHandler(value any) broadcasting.ChannelHandler {
	return func(contractsauth.Authenticatable, ...any) (any, error) {
		return value, nil
	}
}

func assertAccessDenied(t *testing.T, err error) {
	t.Helper()

	if !errors.Is(err, broadcasting.ErrAccessDenied) {
		t.Fatalf("error = %v, want ErrAccessDenied", err)
	}
}

func (c *fakePusherClient) SocketAuth(string, string) (string, error) {
	c.socketAuthCalls++

	return `{"auth":"abcd:efgh"}`, nil
}

func (c *fakePusherClient) PresenceAuth(_ string, _ string, userID string, userInfo any) (string, error) {
	c.presenceAuthCalls++

	data, err := json.Marshal(map[string]any{
		"auth": "abcd:efgh",
		"channel_data": map[string]any{
			"user_id":   userID,
			"user_info": userInfo,
		},
	})

	if err != nil {
		return "", err
	}

	return string(data), nil
}

func (c *fakePusherClient) Trigger(channels []string, event string, payload map[string]any, params map[string]string) error {
	c.triggers = append(c.triggers, pusherTrigger{
		channels: append([]string(nil), channels...),
		event:    event,
		payload:  cloneMap(payload),
		params:   cloneStringMap(params),
	})

	return nil
}

func (c *fakePusherClient) Settings() broadcasting.PusherSettings {
	return c.settings
}

func (p *fakeRedisPublisher) Publish(channel string, payload []byte) error {
	if p.messages == nil {
		p.messages = make(map[string][]byte)
	}

	p.messages[channel] = append([]byte(nil), payload...)

	return nil
}

func (p *fakeAblyPublisher) Publish(channel string, message broadcasting.AblyMessage) error {
	if p.messages == nil {
		p.messages = make(map[string]broadcasting.AblyMessage)
	}

	p.messages[channel] = message

	return nil
}

func (b *recordingBroadcaster) Broadcast(_ context.Context, channels []string, event string, payload map[string]any) error {
	b.calls = append(b.calls, broadcastCall{
		channels: append([]string(nil), channels...),
		event:    event,
		payload:  cloneMap(payload),
	})

	return nil
}

func cloneMap(input map[string]any) map[string]any {
	output := make(map[string]any, len(input))

	for k, v := range input {
		output[k] = v
	}

	return output
}

func cloneStringMap(input map[string]string) map[string]string {
	output := make(map[string]string, len(input))

	for k, v := range input {
		output[k] = v
	}

	return output
}

func requireEqual(t *testing.T, got, want any) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
}

func requireJSONMap(t *testing.T, raw []byte) map[string]any {
	t.Helper()

	var decoded map[string]any

	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("decode json: %v", err)
	}

	return decoded
}

func modelBinding(prefix string) broadcasting.BindingFunc {
	return func(value string) (any, error) {
		if value == "" {
			return nil, fmt.Errorf("missing value")
		}

		return prefix + "." + value + ".instance", nil
	}
}
