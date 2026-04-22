package broadcasting_test

import (
	"context"
	"testing"

	"github.com/bedrock/packages/broadcasting"
)

func TestRedisBroadcasterAuthParity(t *testing.T) {
	t.Parallel()

	// RedisBroadcasterTest::testAuthCallValidAuthenticationResponseWithPrivateChannelWhenCallbackReturnTrue
	// RedisBroadcasterTest::testAuthThrowAccessDeniedHttpExceptionWithPrivateChannelWhenCallbackReturnFalse
	// RedisBroadcasterTest::testAuthThrowAccessDeniedHttpExceptionWithPrivateChannelWhenRequestUserNotFound
	// RedisBroadcasterTest::testAuthCallValidAuthenticationResponseWithPresenceChannelWhenCallbackReturnAnArray
	// RedisBroadcasterTest::testAuthThrowAccessDeniedHttpExceptionWithPresenceChannelWhenCallbackReturnNull
	// RedisBroadcasterTest::testAuthThrowAccessDeniedHttpExceptionWithPresenceChannelWhenRequestUserNotFound
	publisher := &fakeRedisPublisher{}
	b := broadcasting.NewRedisBroadcaster(publisher, "")
	b.Channel("test", allowHandler(true))

	if _, err := b.Auth(requestWithUser("private-test")); err != nil {
		t.Fatalf("private channel auth returned error: %v", err)
	}

	b = broadcasting.NewRedisBroadcaster(publisher, "")
	b.Channel("test", allowHandler(false))
	_, err := b.Auth(requestWithUser("private-test"))
	assertAccessDenied(t, err)

	b = broadcasting.NewRedisBroadcaster(publisher, "")
	b.Channel("test", allowHandler(true))
	_, err = b.Auth(requestWithoutUser("private-test"))
	assertAccessDenied(t, err)

	b = broadcasting.NewRedisBroadcaster(publisher, "")
	b.Channel("test", allowHandler(map[string]any{"a": "b", "c": "d"}))

	if _, err := b.Auth(requestWithUser("presence-test")); err != nil {
		t.Fatalf("presence channel auth returned error: %v", err)
	}

	b = broadcasting.NewRedisBroadcaster(publisher, "")
	b.Channel("test", allowHandler(nil))
	_, err = b.Auth(requestWithUser("presence-test"))
	assertAccessDenied(t, err)

	b = broadcasting.NewRedisBroadcaster(publisher, "")
	b.Channel("test", allowHandler(map[string]any{"a": "b", "c": "d"}))
	_, err = b.Auth(requestWithoutUser("presence-test"))
	assertAccessDenied(t, err)
}

func TestRedisBroadcasterValidAuthenticationResponses(t *testing.T) {
	t.Parallel()

	// RedisBroadcasterTest::testValidAuthenticationResponseWithPrivateChannel
	// RedisBroadcasterTest::testValidAuthenticationResponseWithPresenceChannel
	b := broadcasting.NewRedisBroadcaster(&fakeRedisPublisher{}, "")

	private, err := b.ValidAuthenticationResponse(requestWithUser("private-test"), true)

	if err != nil {
		t.Fatalf("private response returned error: %v", err)
	}

	if private != "true" {
		t.Fatalf("private response = %q, want true", private)
	}

	presence, err := b.ValidAuthenticationResponse(requestWithUser("presence-test"), map[string]any{"a": "b", "c": "d"})

	if err != nil {
		t.Fatalf("presence response returned error: %v", err)
	}

	requireEqual(t, presence, `{"channel_data":{"user_id":"42","user_info":{"a":"b","c":"d"}}}`)
}

func TestRedisBroadcasterBroadcastPublishesPayload(t *testing.T) {
	t.Parallel()

	publisher := &fakeRedisPublisher{}
	b := broadcasting.NewRedisBroadcaster(publisher, "laravel_database_")

	err := b.Broadcast(context.Background(), []string{"orders"}, "OrderUpdated", map[string]any{"id": 1, "socket": "123.456"})

	if err != nil {
		t.Fatalf("Broadcast returned error: %v", err)
	}

	payload := requireJSONMap(t, publisher.messages["laravel_database_orders"])
	requireEqual(t, payload["event"], "OrderUpdated")
	requireEqual(t, payload["socket"], "123.456")
	requireEqual(t, payload["data"], map[string]any{"id": float64(1)})
}
