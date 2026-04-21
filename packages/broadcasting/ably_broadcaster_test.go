package broadcasting_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/broadcasting"
)

func TestAblyBroadcasterAuthParity(t *testing.T) {
	t.Parallel()

	// AblyBroadcasterTest::testAuthCallValidAuthenticationResponseWithPrivateChannelWhenCallbackReturnTrue
	// AblyBroadcasterTest::testAuthThrowAccessDeniedHttpExceptionWithPrivateChannelWhenCallbackReturnFalse
	// AblyBroadcasterTest::testAuthThrowAccessDeniedHttpExceptionWithPrivateChannelWhenRequestUserNotFound
	// AblyBroadcasterTest::testAuthCallValidAuthenticationResponseWithPresenceChannelWhenCallbackReturnAnArray
	// AblyBroadcasterTest::testAuthThrowAccessDeniedHttpExceptionWithPresenceChannelWhenCallbackReturnNull
	// AblyBroadcasterTest::testAuthThrowAccessDeniedHttpExceptionWithPresenceChannelWhenRequestUserNotFound
	publisher := &fakeAblyPublisher{}
	b := broadcasting.NewAblyBroadcaster("abcd:efgh", publisher)
	b.Channel("test", allowHandler(true))

	if _, err := b.Auth(requestWithUser("private-test")); err != nil {
		t.Fatalf("private channel auth returned error: %v", err)
	}

	b = broadcasting.NewAblyBroadcaster("abcd:efgh", publisher)
	b.Channel("test", allowHandler(false))
	_, err := b.Auth(requestWithUser("private-test"))
	assertAccessDenied(t, err)

	b = broadcasting.NewAblyBroadcaster("abcd:efgh", publisher)
	b.Channel("test", allowHandler(true))
	_, err = b.Auth(requestWithoutUser("private-test"))
	assertAccessDenied(t, err)

	b = broadcasting.NewAblyBroadcaster("abcd:efgh", publisher)
	b.Channel("test", allowHandler([]any{1, 2, 3, 4}))

	if _, err := b.Auth(requestWithUser("presence-test")); err != nil {
		t.Fatalf("presence channel auth returned error: %v", err)
	}

	b = broadcasting.NewAblyBroadcaster("abcd:efgh", publisher)
	b.Channel("test", allowHandler(nil))
	_, err = b.Auth(requestWithUser("presence-test"))
	assertAccessDenied(t, err)

	b = broadcasting.NewAblyBroadcaster("abcd:efgh", publisher)
	b.Channel("test", allowHandler([]any{1, 2, 3, 4}))
	_, err = b.Auth(requestWithoutUser("presence-test"))
	assertAccessDenied(t, err)
}

func TestAblyBroadcasterValidAuthenticationResponses(t *testing.T) {
	t.Parallel()

	b := broadcasting.NewAblyBroadcaster("abcd:efgh", &fakeAblyPublisher{})

	private, err := b.ValidAuthenticationResponse(requestWithUser("private-test"), true)

	if err != nil {
		t.Fatalf("private response returned error: %v", err)
	}

	requireEqual(t, private, map[string]any{"auth": "abcd:70d12333dc84a7daade1c59ab2d7fca09a5c6668367ece30cdc2c1ae3d1b1fda"})

	presence, err := b.ValidAuthenticationResponse(requestWithUser("presence-test"), map[string]any{"name": "Taylor"})

	if err != nil {
		t.Fatalf("presence response returned error: %v", err)
	}

	response := presence.(map[string]any)

	if response["auth"] == "" || response["channel_data"] == "" {
		t.Fatalf("presence response missing auth/channel_data: %#v", response)
	}
}

func TestAblyBroadcasterBroadcastFormatsChannels(t *testing.T) {
	t.Parallel()

	publisher := &fakeAblyPublisher{}
	b := broadcasting.NewAblyBroadcaster("abcd:efgh", publisher)

	err := b.Broadcast(context.Background(), []string{"orders", "private-users", "presence-team"}, "OrderUpdated", map[string]any{"id": 1, "socket": "123.456"})

	if err != nil {
		t.Fatalf("Broadcast returned error: %v", err)
	}

	if _, ok := publisher.messages["public:orders"]; !ok {
		t.Fatal("expected public channel publish")
	}

	if _, ok := publisher.messages["private:users"]; !ok {
		t.Fatal("expected private channel publish")
	}

	if _, ok := publisher.messages["presence:team"]; !ok {
		t.Fatal("expected presence channel publish")
	}
}
