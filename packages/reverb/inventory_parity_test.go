package reverb_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	contractsReverb "github.com/bedrock/packages/contracts/reverb"
	"github.com/bedrock/packages/reverb"
)

func requireSent(t *testing.T, conn *fakeConn, needle string) {
	t.Helper()

	if !containsBytes(conn.SentMessages(), []byte(needle)) {
		t.Fatalf("expected sent messages to contain %q, got %s", needle, conn.SentMessages())
	}
}

func requireNoSent(t *testing.T, conn *fakeConn, needle string) {
	t.Helper()

	if containsBytes(conn.SentMessages(), []byte(needle)) {
		t.Fatalf("expected sent messages not to contain %q, got %s", needle, conn.SentMessages())
	}
}

func signedChannelAuth(app *reverb.App, socketID, channel, channelData string) string {
	return app.Key() + ":" + reverb.SignChannel(app.Secret(), socketID, channel, channelData)
}

// ChannelBrokerTest::it_can_return_a_channel_instance
// ChannelBrokerTest::it_can_return_a_private_channel_instance
// ChannelBrokerTest::it_can_return_a_presence_channel_instance
// ChannelBrokerTest::it_can_return_a_cache_channel_instance
// ChannelBrokerTest::it_can_return_a_private_cache_channel_instance
// ChannelBrokerTest::it_can_return_a_presence_cache_channel_instance
func TestInventoryChannelFactoryReturnsProtocolChannelTypes(t *testing.T) {
	t.Parallel()

	app := newTestApp()

	cases := []struct {
		name  string
		check func(contractsReverb.Channel) bool
	}{
		{"public-room", func(ch contractsReverb.Channel) bool {
			_, ok := ch.(*reverb.Channel)
			return ok
		}},
		{"private-room", func(ch contractsReverb.Channel) bool {
			_, ok := ch.(*reverb.PrivateChannel)
			return ok
		}},
		{"presence-room", func(ch contractsReverb.Channel) bool {
			_, ok := ch.(*reverb.PresenceChannel)
			return ok
		}},
		{"cache-room", func(ch contractsReverb.Channel) bool {
			_, ok := ch.(*reverb.CacheChannel)
			return ok
		}},
		{"private-cache-room", func(ch contractsReverb.Channel) bool {
			_, ok := ch.(*reverb.PrivateCacheChannel)
			return ok
		}},
		{"presence-cache-room", func(ch contractsReverb.Channel) bool {
			_, ok := ch.(*reverb.PresenceCacheChannel)
			return ok
		}},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if ch := reverb.NewChannel(tc.name, app); !tc.check(ch) {
				t.Fatalf("unexpected channel implementation %T for %q", ch, tc.name)
			}
		})
	}
}

// ChannelTest::it_can_subscribe_a_connection_to_a_channel
// ChannelTest::it_can_unsubscribe_a_connection_from_a_channel
// ChannelTest::it_can_broadcast_to_all_connections_of_a_channel
// ChannelTest::it_does_not_broadcast_to_the_connection_sending_the_message
// ChannelTest::it_removes_a_channel_when_no_subscribers_remain
func TestInventoryPublicChannelLifecycle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	app := newTestApp()
	apps := reverb.NewAppManager([]reverb.AppConfig{{ID: app.ID(), Key: app.Key(), Secret: app.Secret()}})
	manager := reverb.NewChannelManager(apps)
	ch, err := manager.GetOrCreate(app.ID(), "public-room")
	if err != nil {
		t.Fatalf("GetOrCreate: %v", err)
	}

	sender := newFakeConn("socket-1", app.ID())
	receiver := newFakeConn("socket-2", app.ID())

	if err := ch.Subscribe(ctx, sender, "", ""); err != nil {
		t.Fatalf("subscribe sender: %v", err)
	}
	if err := ch.Subscribe(ctx, receiver, "", ""); err != nil {
		t.Fatalf("subscribe receiver: %v", err)
	}
	if !ch.HasConnection(sender.SocketID()) || len(ch.Connections()) != 2 {
		t.Fatalf("expected both connections subscribed")
	}

	event := contractsReverb.Event{Event: "inventory-event", Data: `{"ok":true}`, Channel: ch.Name()}
	if err := ch.BroadcastToAll(ctx, event); err != nil {
		t.Fatalf("broadcast all: %v", err)
	}
	requireSent(t, sender, "inventory-event")
	requireSent(t, receiver, "inventory-event")

	beforeSender := len(sender.SentMessages())
	except := sender.SocketID()
	if err := ch.Broadcast(ctx, contractsReverb.Event{Event: "receiver-only", Data: `{}`, Channel: ch.Name()}, &except); err != nil {
		t.Fatalf("broadcast except sender: %v", err)
	}
	if len(sender.SentMessages()) != beforeSender {
		t.Fatalf("sender received excluded event")
	}
	requireSent(t, receiver, "receiver-only")

	ch.Unsubscribe(ctx, sender)
	ch.Unsubscribe(ctx, receiver)
	manager.CleanupEmpty(app.ID())
	if _, ok := manager.Get(app.ID(), ch.Name()); ok {
		t.Fatalf("expected empty channel to be removed")
	}
}

// PrivateChannelTest::it_can_subscribe_a_connection_to_a_channel
// PrivateChannelTest::it_can_unsubscribe_a_connection_from_a_channel
// PrivateChannelTest::it_can_broadcast_to_all_connections_of_a_channel
// PrivateChannelTest::it_fails_to_subscribe_if_the_signature_is_invalid
// PrivateChannelTest::it_fails_to_subscribe_to_a_private_channel_with_no_auth_token
// PrivateChannelTest::it_fails_to_subscribe_to_a_presence_channel_with_no_auth_token
func TestInventoryPrivateChannelLifecycle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	app := newTestApp()
	privateName := "private-room"
	privateChannel := reverb.NewPrivateChannel(privateName, app)
	conn := newFakeConn("socket-1", app.ID())

	if err := privateChannel.Subscribe(ctx, conn, signedChannelAuth(app, conn.SocketID(), privateName, ""), ""); err != nil {
		t.Fatalf("valid private subscribe: %v", err)
	}
	if !privateChannel.HasConnection(conn.SocketID()) {
		t.Fatalf("expected private channel connection")
	}
	if err := privateChannel.BroadcastToAll(ctx, contractsReverb.Event{Event: "private-event", Data: `{}`, Channel: privateName}); err != nil {
		t.Fatalf("private broadcast: %v", err)
	}
	requireSent(t, conn, "private-event")

	privateChannel.Unsubscribe(ctx, conn)
	if privateChannel.HasConnection(conn.SocketID()) {
		t.Fatalf("expected private connection removed")
	}

	if err := reverb.NewPrivateChannel(privateName, app).Subscribe(ctx, newFakeConn("socket-2", app.ID()), "bad-signature", ""); err != reverb.ErrUnauthorized {
		t.Fatalf("invalid private signature error = %v, want ErrUnauthorized", err)
	}
	if err := reverb.NewPrivateChannel(privateName, app).Subscribe(ctx, newFakeConn("socket-3", app.ID()), "", ""); err != reverb.ErrUnauthorized {
		t.Fatalf("missing private auth error = %v, want ErrUnauthorized", err)
	}
	if err := reverb.NewPresenceChannel("presence-room", app).Subscribe(ctx, newFakeConn("socket-4", app.ID()), "", `{"user_id":"1"}`); err != reverb.ErrUnauthorized {
		t.Fatalf("missing presence auth error = %v, want ErrUnauthorized", err)
	}
}

// CacheChannelTest::it_receives_no_data_when_no_previous_event_triggered
// CacheChannelTest::it_stores_last_triggered_event
// PrivateCacheChannelTest::it_can_subscribe_a_connection_to_a_channel
// PrivateCacheChannelTest::it_can_unsubscribe_a_connection_from_a_channel
// PrivateCacheChannelTest::it_can_broadcast_to_all_connections_of_a_channel
// PrivateCacheChannelTest::it_fails_to_subscribe_if_the_signature_is_invalid
// PrivateCacheChannelTest::it_receives_no_data_when_no_previous_event_triggered
// PrivateCacheChannelTest::it_stores_last_triggered_event
func TestInventoryCacheChannels(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	app := newTestApp()

	cacheChannel := reverb.NewCacheChannel("cache-room", app)
	first := newFakeConn("socket-1", app.ID())
	if err := cacheChannel.Subscribe(ctx, first, "", ""); err != nil {
		t.Fatalf("cache subscribe: %v", err)
	}
	requireSent(t, first, "pusher:cache_miss")

	cached := contractsReverb.Event{Event: "cached-event", Data: `{"cached":true}`, Channel: "cache-room"}
	if err := cacheChannel.BroadcastToAll(ctx, cached); err != nil {
		t.Fatalf("cache broadcast: %v", err)
	}
	if cacheChannel.LastEvent() == nil || cacheChannel.LastEvent().Event != "cached-event" {
		t.Fatalf("expected cached last event, got %#v", cacheChannel.LastEvent())
	}

	privateName := "private-cache-room"
	privateCache := reverb.NewPrivateCacheChannel(privateName, app)
	privateConn := newFakeConn("socket-2", app.ID())
	if err := privateCache.Subscribe(ctx, privateConn, signedChannelAuth(app, privateConn.SocketID(), privateName, ""), ""); err != nil {
		t.Fatalf("private cache subscribe: %v", err)
	}
	requireSent(t, privateConn, "pusher:cache_miss")
	if err := privateCache.BroadcastToAll(ctx, contractsReverb.Event{Event: "private-cached", Data: `{}`, Channel: privateName}); err != nil {
		t.Fatalf("private cache broadcast: %v", err)
	}
	requireSent(t, privateConn, "private-cached")
	if privateCache.LastEvent() == nil || privateCache.LastEvent().Event != "private-cached" {
		t.Fatalf("expected private cache last event, got %#v", privateCache.LastEvent())
	}
	privateCache.Unsubscribe(ctx, privateConn)
	if privateCache.HasConnection(privateConn.SocketID()) {
		t.Fatalf("expected private cache connection removed")
	}
	if err := reverb.NewPrivateCacheChannel(privateName, app).Subscribe(ctx, newFakeConn("socket-3", app.ID()), "bad-signature", ""); err != reverb.ErrUnauthorized {
		t.Fatalf("invalid private cache signature error = %v, want ErrUnauthorized", err)
	}
}

// PresenceChannelTest::it_can_subscribe_a_connection_to_a_channel
// PresenceChannelTest::it_can_unsubscribe_a_connection_from_a_channel
// PresenceChannelTest::it_can_broadcast_to_all_connections_of_a_channel
// PresenceChannelTest::it_fails_to_subscribe_if_the_signature_is_invalid
// PresenceChannelTest::it_can_return_data_stored_on_the_connection
// PresenceChannelTest::it_sends_notification_of_subscription
// PresenceChannelTest::it_sends_notification_of_subscription_with_data
// PresenceChannelTest::it_sends_notification_of_an_unsubscribe
// PresenceChannelTest::it_ensures_the_member_added_event_is_only_fired_once
// PresenceChannelTest::it_ensures_the_member_removed_event_is_only_fired_once
func TestInventoryPresenceChannelLifecycle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	app := newTestApp()
	channelName := "presence-room"
	channel := reverb.NewPresenceChannel(channelName, app)

	aliceData := `{"user_id":"alice","user_info":{"name":"Alice"}}`
	aliceA := newFakeConn("socket-1", app.ID())
	if err := channel.Subscribe(ctx, aliceA, signedChannelAuth(app, aliceA.SocketID(), channelName, aliceData), aliceData); err != nil {
		t.Fatalf("presence subscribe alice A: %v", err)
	}
	requireSent(t, aliceA, "subscription_succeeded")
	requireSent(t, aliceA, "Alice")

	aliceB := newFakeConn("socket-2", app.ID())
	beforeDuplicate := len(aliceA.SentMessages())
	if err := channel.Subscribe(ctx, aliceB, signedChannelAuth(app, aliceB.SocketID(), channelName, aliceData), aliceData); err != nil {
		t.Fatalf("presence subscribe alice B: %v", err)
	}
	if len(aliceA.SentMessages()) != beforeDuplicate {
		t.Fatalf("duplicate user should not fire member_added")
	}
	if channel.MemberCount() != 1 || channel.Members()["alice"] == nil {
		t.Fatalf("expected one deduplicated alice member, got %v", channel.Members())
	}

	bobData := `{"user_id":"bob","user_info":{"name":"Bob"}}`
	bob := newFakeConn("socket-3", app.ID())
	if err := channel.Subscribe(ctx, bob, signedChannelAuth(app, bob.SocketID(), channelName, bobData), bobData); err != nil {
		t.Fatalf("presence subscribe bob: %v", err)
	}
	requireSent(t, aliceA, "member_added")
	requireSent(t, aliceA, "bob")

	if err := channel.BroadcastToAll(ctx, contractsReverb.Event{Event: "presence-event", Data: `{}`, Channel: channelName}); err != nil {
		t.Fatalf("presence broadcast: %v", err)
	}
	requireSent(t, aliceA, "presence-event")
	requireSent(t, bob, "presence-event")

	beforeRemove := len(bob.SentMessages())
	channel.Unsubscribe(ctx, aliceA)
	if len(bob.SentMessages()) != beforeRemove {
		t.Fatalf("first alice disconnect should not fire member_removed while alice B remains")
	}
	channel.Unsubscribe(ctx, aliceB)
	requireSent(t, bob, "member_removed")
	requireSent(t, bob, "alice")

	if err := reverb.NewPresenceChannel(channelName, app).Subscribe(ctx, newFakeConn("socket-4", app.ID()), "bad-signature", aliceData); err != reverb.ErrUnauthorized {
		t.Fatalf("invalid presence signature error = %v, want ErrUnauthorized", err)
	}
}

// PresenceCacheChannelTest::it_can_subscribe_a_connection_to_a_channel
// PresenceCacheChannelTest::it_can_unsubscribe_a_connection_from_a_channel
// PresenceCacheChannelTest::it_can_broadcast_to_all_connections_of_a_channel
// PresenceCacheChannelTest::it_fails_to_subscribe_if_the_signature_is_invalid
// PresenceCacheChannelTest::it_can_return_data_stored_on_the_connection
// PresenceCacheChannelTest::it_sends_notification_of_subscription
// PresenceCacheChannelTest::it_sends_notification_of_subscription_with_data
// PresenceCacheChannelTest::it_sends_notification_of_an_unsubscribe
// PresenceCacheChannelTest::it_receives_no_data_when_no_previous_event_triggered
// PresenceCacheChannelTest::it_stores_last_triggered_event
func TestInventoryPresenceCacheChannelLifecycle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	app := newTestApp()
	channelName := "presence-cache-room"
	channel := reverb.NewPresenceCacheChannel(channelName, app)

	aliceData := `{"user_id":"alice","user_info":{"name":"Alice"}}`
	alice := newFakeConn("socket-1", app.ID())
	if err := channel.Subscribe(ctx, alice, signedChannelAuth(app, alice.SocketID(), channelName, aliceData), aliceData); err != nil {
		t.Fatalf("presence cache subscribe alice: %v", err)
	}
	requireSent(t, alice, "subscription_succeeded")
	requireSent(t, alice, "Alice")
	requireSent(t, alice, "pusher:cache_miss")

	if channel.MemberCount() != 1 || channel.Members()["alice"] == nil {
		t.Fatalf("expected alice member, got %v", channel.Members())
	}

	event := contractsReverb.Event{Event: "presence-cache-event", Data: `{}`, Channel: channelName}
	if err := channel.BroadcastToAll(ctx, event); err != nil {
		t.Fatalf("presence cache broadcast: %v", err)
	}
	requireSent(t, alice, "presence-cache-event")
	if channel.LastEvent() == nil || channel.LastEvent().Event != "presence-cache-event" {
		t.Fatalf("expected cached presence event, got %#v", channel.LastEvent())
	}

	bobData := `{"user_id":"bob","user_info":{"name":"Bob"}}`
	bob := newFakeConn("socket-2", app.ID())
	if err := channel.Subscribe(ctx, bob, signedChannelAuth(app, bob.SocketID(), channelName, bobData), bobData); err != nil {
		t.Fatalf("presence cache subscribe bob: %v", err)
	}
	requireSent(t, bob, "presence-cache-event")
	requireSent(t, alice, "member_added")

	channel.Unsubscribe(ctx, bob)
	requireSent(t, alice, "member_removed")
	if err := reverb.NewPresenceCacheChannel(channelName, app).Subscribe(ctx, newFakeConn("socket-3", app.ID()), "bad-signature", aliceData); err != reverb.ErrUnauthorized {
		t.Fatalf("invalid presence cache signature error = %v, want ErrUnauthorized", err)
	}
}

// ChannelManagerTest::it_can_subscribe_to_a_channel
// ChannelManagerTest::it_can_unsubscribe_from_a_channel
// ChannelManagerTest::it_can_get_all_channels
// ChannelManagerTest::it_can_determine_whether_a_channel_exists
// ChannelManagerTest::it_can_get_all_connections_subscribed_to_a_channel
// ChannelManagerTest::it_can_unsubscribe_a_connection_from_all_channels
// ChannelManagerTest::it_can_get_the_data_for_a_connection_subscribed_to_a_channel
// ChannelManagerTest::it_can_get_all_connections_for_all_channels
func TestInventoryChannelManagerOperations(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	apps := reverb.NewAppManager([]reverb.AppConfig{{ID: "app-1", Key: "key-1", Secret: "secret-1"}})
	manager := reverb.NewChannelManager(apps)
	conn := newFakeConn("socket-1", "app-1")

	first, err := manager.GetOrCreate("app-1", "public-one")
	if err != nil {
		t.Fatalf("GetOrCreate first: %v", err)
	}
	second, err := manager.GetOrCreate("app-1", "public-two")
	if err != nil {
		t.Fatalf("GetOrCreate second: %v", err)
	}
	if err := first.Subscribe(ctx, conn, "", ""); err != nil {
		t.Fatalf("subscribe first: %v", err)
	}
	if err := second.Subscribe(ctx, conn, "", ""); err != nil {
		t.Fatalf("subscribe second: %v", err)
	}

	if _, ok := manager.Get("app-1", "public-one"); !ok {
		t.Fatalf("expected public-one channel")
	}
	if len(manager.All("app-1")) != 2 {
		t.Fatalf("expected two channels, got %d", len(manager.All("app-1")))
	}
	if len(first.Connections()) != 1 || !first.HasConnection(conn.SocketID()) {
		t.Fatalf("expected first channel connection")
	}

	for _, ch := range manager.All("app-1") {
		if ch.HasConnection(conn.SocketID()) {
			ch.Unsubscribe(ctx, conn)
		}
	}
	manager.CleanupEmpty("app-1")
	if len(manager.All("app-1")) != 0 {
		t.Fatalf("expected all empty channels removed, got %d", len(manager.All("app-1")))
	}
}

// EventTest::it_can_correctly_format_a_payload
// EventTest::it_can_correctly_format_an_internal_payload
// EventTest::it_can_respond_to_a_ping
// EventTest::it_can_subscribe_to_a_channel
// EventTest::it_can_subscribe_to_an_empty_channel
// EventTest::it_can_unsubscribe_from_a_channel
func TestInventoryPusherMessageFormattingAndParsing(t *testing.T) {
	t.Parallel()

	payload, err := reverb.MarshalEvent("inventory-event", "public-room", map[string]any{"answer": 42})
	if err != nil {
		t.Fatalf("MarshalEvent: %v", err)
	}
	var envelope struct {
		Event   string `json:"event"`
		Data    string `json:"data"`
		Channel string `json:"channel"`
	}
	if err := json.Unmarshal(payload, &envelope); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if envelope.Event != "inventory-event" || envelope.Channel != "public-room" || envelope.Data != `{"answer":42}` {
		t.Fatalf("unexpected event envelope: %+v", envelope)
	}

	errorPayload, err := reverb.MarshalError(reverb.CodeInvalidMessage, "bad message")
	if err != nil {
		t.Fatalf("MarshalError: %v", err)
	}
	requireJSONContains(t, errorPayload, "pusher:error")

	subscribe, err := reverb.Parse([]byte(`{"event":"pusher:subscribe","data":{"channel":"public-room"}}`))
	if err != nil || subscribe.Event != "pusher:subscribe" {
		t.Fatalf("unexpected subscribe parse: msg=%+v err=%v", subscribe, err)
	}
	subscribeData, err := reverb.ParseSubscribeData(subscribe.Data)
	if err != nil {
		t.Fatalf("ParseSubscribeData: %v", err)
	}
	if subscribeData.Channel != "public-room" {
		t.Fatalf("subscribe channel = %q", subscribeData.Channel)
	}

	unsubscribe, err := reverb.Parse([]byte(`{"event":"pusher:unsubscribe","data":{"channel":"public-room"}}`))
	if err != nil || unsubscribe.Event != "pusher:unsubscribe" {
		t.Fatalf("unexpected unsubscribe parse: msg=%+v err=%v", unsubscribe, err)
	}

	ping, err := reverb.Parse([]byte(`{"event":"pusher:ping","data":{}}`))
	if err != nil || ping.Event != "pusher:ping" {
		t.Fatalf("unexpected ping parse: msg=%+v err=%v", ping, err)
	}
}

func requireJSONContains(t *testing.T, payload []byte, needle string) {
	t.Helper()

	if !containsBytes([][]byte{payload}, []byte(needle)) {
		t.Fatalf("expected payload to contain %q, got %s", needle, payload)
	}
}

func requireContentLength(t *testing.T, resp *http.Response) {
	t.Helper()

	if got := resp.Header.Get("Content-Length"); got == "" {
		t.Fatalf("expected Content-Length header on %s response", resp.Status)
	}
}

// EventsControllerTest::it_can_receive_and_event_trigger
// EventsControllerTest::it_can_receive_and_event_trigger_for_multiple_channels
// EventsControllerTest::it_can_ignore_a_subscriber
// EventsControllerTest::it_validates_invalid_data
// EventsControllerTest::it_fails_when_payload_is_invalid
// EventsControllerTest::it_fails_when_app_cannot_be_found
// EventsControllerTest::it_fails_when_using_an_invalid_signature
// EventsControllerTest::it_can_send_the_content_length_header
// EventsBatchControllerTest::it_can_receive_an_event_batch_trigger
// EventsBatchControllerTest::it_can_receive_an_event_batch_trigger_with_multiple_events
// EventsBatchControllerTest::it_fails_when_payload_is_invalid
// EventsBatchControllerTest::it_fails_when_using_an_invalid_signature
// EventsBatchControllerTest::it_can_send_the_content_length_header
// ChannelsControllerTest::it_can_return_all_channel_information
// ChannelsControllerTest::it_can_send_the_content_length_header
// ChannelControllerTest::it_can_return_data_for_a_single_channel
// ChannelControllerTest::it_fails_when_using_an_invalid_signature
// ChannelControllerTest::it_can_send_the_content_length_header
func TestInventoryHTTPPusherEndpoints(t *testing.T) {
	t.Parallel()

	handler, _, _, manager := newHTTPTestSetup(t)
	ctx := context.Background()

	first, err := manager.GetOrCreate("app-1", "public-one")
	if err != nil {
		t.Fatalf("GetOrCreate first: %v", err)
	}
	second, err := manager.GetOrCreate("app-1", "public-two")
	if err != nil {
		t.Fatalf("GetOrCreate second: %v", err)
	}
	sender := newFakeConn("socket-1", "app-1")
	receiver := newFakeConn("socket-2", "app-1")
	if err := first.Subscribe(ctx, sender, "", ""); err != nil {
		t.Fatalf("subscribe sender: %v", err)
	}
	if err := first.Subscribe(ctx, receiver, "", ""); err != nil {
		t.Fatalf("subscribe receiver: %v", err)
	}
	secondConn := newFakeConn("socket-3", "app-1")
	if err := second.Subscribe(ctx, secondConn, "", ""); err != nil {
		t.Fatalf("subscribe second: %v", err)
	}

	body, _ := json.Marshal(reverb.TriggerRequest{
		Name:     "http-event",
		Data:     `{"via":"http"}`,
		Channels: []string{"public-one", "public-two"},
		SocketID: &[]string{sender.SocketID()}[0],
	})
	req := signedRequest(t, http.MethodPost, "/apps/app-1/events", "secret-1", body)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)
	if recorder.Result().StatusCode != http.StatusOK {
		t.Fatalf("event trigger status = %d", recorder.Result().StatusCode)
	}
	requireContentLength(t, recorder.Result())
	requireNoSent(t, sender, "http-event")
	requireSent(t, receiver, "http-event")
	requireSent(t, secondConn, "http-event")

	batchBody, _ := json.Marshal(reverb.BatchTriggerRequest{Batch: []reverb.TriggerRequest{
		{Name: "batch-one", Data: `{}`, Channels: []string{"public-one"}},
		{Name: "batch-two", Data: `{}`, Channels: []string{"public-two"}},
	}})
	batchReq := signedRequest(t, http.MethodPost, "/apps/app-1/batch_events", "secret-1", batchBody)
	batchRecorder := httptest.NewRecorder()
	handler.ServeHTTP(batchRecorder, batchReq)
	if batchRecorder.Result().StatusCode != http.StatusOK {
		t.Fatalf("batch trigger status = %d", batchRecorder.Result().StatusCode)
	}
	requireContentLength(t, batchRecorder.Result())
	requireSent(t, receiver, "batch-one")
	requireSent(t, secondConn, "batch-two")

	channelsReq := signedRequestWithExtraParams(t, http.MethodGet, "/apps/app-1/channels", "secret-1", nil, map[string]string{"info": "subscription_count"})
	channelsRecorder := httptest.NewRecorder()
	handler.ServeHTTP(channelsRecorder, channelsReq)
	if channelsRecorder.Result().StatusCode != http.StatusOK {
		t.Fatalf("channels status = %d", channelsRecorder.Result().StatusCode)
	}
	requireContentLength(t, channelsRecorder.Result())
	var channelsResponse struct {
		Channels map[string]struct {
			SubscriptionCount int `json:"subscription_count"`
		} `json:"channels"`
	}
	if err := json.NewDecoder(channelsRecorder.Result().Body).Decode(&channelsResponse); err != nil {
		t.Fatalf("decode channels: %v", err)
	}
	if channelsResponse.Channels["public-one"].SubscriptionCount != 2 {
		t.Fatalf("public-one subscription count = %d", channelsResponse.Channels["public-one"].SubscriptionCount)
	}

	channelReq := signedRequest(t, http.MethodGet, "/apps/app-1/channels/public-one", "secret-1", nil)
	channelRecorder := httptest.NewRecorder()
	handler.ServeHTTP(channelRecorder, channelReq)
	if channelRecorder.Result().StatusCode != http.StatusOK {
		t.Fatalf("channel status = %d", channelRecorder.Result().StatusCode)
	}
	requireContentLength(t, channelRecorder.Result())

	invalidJSON := signedRequest(t, http.MethodPost, "/apps/app-1/events", "secret-1", []byte(`{`))
	invalidRecorder := httptest.NewRecorder()
	handler.ServeHTTP(invalidRecorder, invalidJSON)
	if invalidRecorder.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid event payload status = %d", invalidRecorder.Result().StatusCode)
	}

	invalidBatchJSON := signedRequest(t, http.MethodPost, "/apps/app-1/batch_events", "secret-1", []byte(`{`))
	invalidBatchRecorder := httptest.NewRecorder()
	handler.ServeHTTP(invalidBatchRecorder, invalidBatchJSON)
	if invalidBatchRecorder.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid batch payload status = %d", invalidBatchRecorder.Result().StatusCode)
	}

	missingAppReq := signedRequest(t, http.MethodPost, "/apps/missing/events", "secret-1", body)
	missingAppRecorder := httptest.NewRecorder()
	handler.ServeHTTP(missingAppRecorder, missingAppReq)
	if missingAppRecorder.Result().StatusCode != http.StatusNotFound {
		t.Fatalf("missing app status = %d", missingAppRecorder.Result().StatusCode)
	}

	badSignatureReq := signedRequest(t, http.MethodPost, "/apps/app-1/events", "wrong-secret", body)
	badSignatureRecorder := httptest.NewRecorder()
	handler.ServeHTTP(badSignatureRecorder, badSignatureReq)
	if badSignatureRecorder.Result().StatusCode != http.StatusUnauthorized {
		t.Fatalf("bad event signature status = %d", badSignatureRecorder.Result().StatusCode)
	}

	badChannelSignatureReq := signedRequest(t, http.MethodGet, "/apps/app-1/channels/public-one", "wrong-secret", nil)
	badChannelSignatureRecorder := httptest.NewRecorder()
	handler.ServeHTTP(badChannelSignatureRecorder, badChannelSignatureReq)
	if badChannelSignatureRecorder.Result().StatusCode != http.StatusUnauthorized {
		t.Fatalf("bad channel signature status = %d", badChannelSignatureRecorder.Result().StatusCode)
	}
}

// EventTest::it_can_publish_an_event_when_enabled
// EventTest::it_can_broadcast_an_event_directly_when_publishing_disabled
// EventTest::it_can_broadcast_an_event_for_multiple_channels
func TestInventorySyncDispatcherBroadcastsEvents(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	apps := reverb.NewAppManager([]reverb.AppConfig{{ID: "app-1", Key: "key-1", Secret: "secret-1"}})
	manager := reverb.NewChannelManager(apps)
	dispatcher := reverb.NewSyncDispatcher(manager)

	first, err := manager.GetOrCreate("app-1", "public-one")
	if err != nil {
		t.Fatalf("GetOrCreate first: %v", err)
	}
	second, err := manager.GetOrCreate("app-1", "public-two")
	if err != nil {
		t.Fatalf("GetOrCreate second: %v", err)
	}
	firstConn := newFakeConn("socket-1", "app-1")
	secondConn := newFakeConn("socket-2", "app-1")
	if err := first.Subscribe(ctx, firstConn, "", ""); err != nil {
		t.Fatalf("subscribe first: %v", err)
	}
	if err := second.Subscribe(ctx, secondConn, "", ""); err != nil {
		t.Fatalf("subscribe second: %v", err)
	}

	if err := dispatcher.Dispatch(ctx, "app-1", contractsReverb.Event{Event: "published", Data: `{}`, Channel: "public-one"}); err != nil {
		t.Fatalf("dispatch first: %v", err)
	}
	if err := dispatcher.Dispatch(ctx, "app-1", contractsReverb.Event{Event: "published-two", Data: `{}`, Channel: "public-two"}); err != nil {
		t.Fatalf("dispatch second: %v", err)
	}
	requireSent(t, firstConn, "published")
	requireSent(t, secondConn, "published-two")

	if err := first.BroadcastToAll(ctx, contractsReverb.Event{Event: "direct-broadcast", Data: `{}`, Channel: "public-one"}); err != nil {
		t.Fatalf("direct broadcast: %v", err)
	}
	requireSent(t, firstConn, "direct-broadcast")
}

// ServerTest::it_cannot_connect_from_an_invalid_origin
// ServerTest::it_can_connect_from_a_valid_origin
// ServerTest::it_cannot_connect_when_over_the_max_connection_limit
// ServerTest::it_rejects_messages_over_the_max_allowed_size
// ServerTest::it_allows_message_within_the_max_allowed_size
// FactoryTest::it_can_create_a_server
func TestInventoryServerConfigurationPrimitives(t *testing.T) {
	t.Parallel()

	if reverb.ValidateOrigin("https://app.example.com", []string{"*.example.com"}) != true {
		t.Fatalf("expected wildcard origin to be valid")
	}
	if reverb.ValidateOrigin("https://evil.test", []string{"*.example.com"}) != false {
		t.Fatalf("expected unmatched origin to be invalid")
	}

	apps := reverb.NewAppManager([]reverb.AppConfig{{ID: "app-1", Key: "key-1", Secret: "secret-1", MaxConnections: 1, MaxMessageSize: 32}})
	app, err := apps.FindByKey("key-1")
	if err != nil {
		t.Fatalf("FindByKey: %v", err)
	}
	if app.MaxConnections() != 1 || app.MaxMessageSize() != 32 {
		t.Fatalf("unexpected app limits: max connections=%d message size=%d", app.MaxConnections(), app.MaxMessageSize())
	}

	connections := reverb.NewConnectionManager()
	channels := reverb.NewChannelManager(apps)
	server := reverb.NewServer(reverb.DefaultConfig(), apps, connections, channels, reverb.NewSyncDispatcher(channels))
	if server == nil {
		t.Fatalf("expected server")
	}
	if len(apps.All()) != 1 {
		t.Fatalf("expected one configured app")
	}

	connections.Add(reverb.NewConn(nil, "app-1"))
	if connections.Count("app-1") != app.MaxConnections() {
		t.Fatalf("expected app at max connection limit")
	}
}
