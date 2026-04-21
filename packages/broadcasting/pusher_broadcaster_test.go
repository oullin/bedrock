package broadcasting_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/broadcasting"
)

func TestPusherBroadcasterAuthParity(t *testing.T) {
	t.Parallel()

	// PusherBroadcasterTest::testAuthCallValidAuthenticationResponseWithPrivateChannelWhenCallbackReturnTrue
	// PusherBroadcasterTest::testAuthThrowAccessDeniedHttpExceptionWithPrivateChannelWhenCallbackReturnFalse
	// PusherBroadcasterTest::testAuthThrowAccessDeniedHttpExceptionWithPrivateChannelWhenRequestUserNotFound
	// PusherBroadcasterTest::testAuthCallValidAuthenticationResponseWithPresenceChannelWhenCallbackReturnAnArray
	// PusherBroadcasterTest::testAuthThrowAccessDeniedHttpExceptionWithPresenceChannelWhenCallbackReturnNull
	// PusherBroadcasterTest::testAuthThrowAccessDeniedHttpExceptionWithPresenceChannelWhenRequestUserNotFound
	client := &fakePusherClient{}
	b := broadcasting.NewPusherBroadcaster(client)
	b.Channel("test", allowHandler(true))

	if _, err := b.Auth(requestWithUser("private-test")); err != nil {
		t.Fatalf("private channel auth returned error: %v", err)
	}

	b = broadcasting.NewPusherBroadcaster(client)
	b.Channel("test", allowHandler(false))
	_, err := b.Auth(requestWithUser("private-test"))
	assertAccessDenied(t, err)

	b = broadcasting.NewPusherBroadcaster(client)
	b.Channel("test", allowHandler(true))
	_, err = b.Auth(requestWithoutUser("private-test"))
	assertAccessDenied(t, err)

	b = broadcasting.NewPusherBroadcaster(client)
	b.Channel("test", allowHandler([]any{1, 2, 3, 4}))

	if _, err := b.Auth(requestWithUser("presence-test")); err != nil {
		t.Fatalf("presence channel auth returned error: %v", err)
	}

	b = broadcasting.NewPusherBroadcaster(client)
	b.Channel("test", allowHandler(nil))
	_, err = b.Auth(requestWithUser("presence-test"))
	assertAccessDenied(t, err)

	b = broadcasting.NewPusherBroadcaster(client)
	b.Channel("test", allowHandler([]any{1, 2, 3, 4}))
	_, err = b.Auth(requestWithoutUser("presence-test"))
	assertAccessDenied(t, err)
}

func TestPusherBroadcasterValidAuthenticationResponses(t *testing.T) {
	t.Parallel()

	// PusherBroadcasterTest::testValidAuthenticationResponseCallPusherSocketAuthMethodWithPrivateChannel
	// PusherBroadcasterTest::testValidAuthenticationResponseCallPusherPresenceAuthMethodWithPresenceChannel
	client := &fakePusherClient{}
	b := broadcasting.NewPusherBroadcaster(client)

	private, err := b.ValidAuthenticationResponse(requestWithUser("private-test"), true)

	if err != nil {
		t.Fatalf("private response returned error: %v", err)
	}

	requireEqual(t, private, map[string]any{"auth": "abcd:efgh"})

	if client.socketAuthCalls != 1 {
		t.Fatalf("socket auth calls = %d, want 1", client.socketAuthCalls)
	}

	presence, err := b.ValidAuthenticationResponse(requestWithUser("presence-test"), []any{1, 2, 3, 4})

	if err != nil {
		t.Fatalf("presence response returned error: %v", err)
	}

	requireEqual(t, presence, map[string]any{
		"auth": "abcd:efgh",
		"channel_data": map[string]any{
			"user_id":   "42",
			"user_info": []any{float64(1), float64(2), float64(3), float64(4)},
		},
	})

	if client.presenceAuthCalls != 1 {
		t.Fatalf("presence auth calls = %d, want 1", client.presenceAuthCalls)
	}
}

func TestPusherBroadcasterUserAuthentication(t *testing.T) {
	t.Parallel()

	// PusherBroadcasterTest::testUserAuthenticationForPusher
	client := &fakePusherClient{settings: broadcasting.PusherSettings{
		AuthKey: "278d425bdf160c739803",
		Secret:  "7ad3773142a6692b25b8",
	}}
	b := broadcasting.NewPusherBroadcaster(client)
	b.ResolveAuthenticatedUserUsing(func(broadcasting.AuthRequest) any {
		return map[string]any{"id": "12345"}
	})

	response, err := b.ResolveAuthenticatedUser(broadcasting.AuthRequest{SocketID: "1234.1234"})

	if err != nil {
		t.Fatalf("ResolveAuthenticatedUser returned error: %v", err)
	}

	requireEqual(t, response, map[string]any{
		"auth":      "278d425bdf160c739803:4708d583dada6a56435fb8bc611c77c359a31eebde13337c16ab43aa6de336ba",
		"user_data": `{"id":"12345"}`,
	})
}

func TestPusherBroadcasterBroadcastPassesSocketAsParameter(t *testing.T) {
	t.Parallel()

	client := &fakePusherClient{}
	b := broadcasting.NewPusherBroadcaster(client)

	err := b.Broadcast(context.Background(), []string{"orders"}, "OrderUpdated", map[string]any{
		"id":     1,
		"socket": "123.456",
	})

	if err != nil {
		t.Fatalf("Broadcast returned error: %v", err)
	}

	requireEqual(t, client.triggers[0].payload, map[string]any{"id": 1})
	requireEqual(t, client.triggers[0].params, map[string]string{"socket_id": "123.456"})
}
