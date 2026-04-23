package reverb_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	contractsReverb "github.com/bedrock/packages/contracts/reverb"
	"github.com/bedrock/packages/reverb"
	"nhooyr.io/websocket"
)

type wsEnvelope struct {
	Event   string `json:"event"`
	Data    string `json:"data"`
	Channel string `json:"channel"`
}

func newWebSocketHarness(t *testing.T, cfg reverb.AppConfig) (*httptest.Server, *reverb.App, *reverb.ChannelManager, *reverb.SyncDispatcher) {
	ts, app, _, _, channels, dispatcher := newWebSocketHarnessWithManagers(t, cfg)

	return ts, app, channels, dispatcher
}

func newWebSocketHarnessWithManagers(t *testing.T, cfg reverb.AppConfig) (*httptest.Server, *reverb.App, *reverb.AppManager, *reverb.ConnectionManager, *reverb.ChannelManager, *reverb.SyncDispatcher) {
	t.Helper()

	apps := reverb.NewAppManager([]reverb.AppConfig{cfg})
	app, err := apps.FindByID(cfg.ID)

	if err != nil {
		t.Fatalf("find configured app: %v", err)
	}

	channels := reverb.NewChannelManager(apps)
	conns := reverb.NewConnectionManager()
	dispatcher := reverb.NewSyncDispatcher(channels)
	server := reverb.NewServer(reverb.DefaultConfig(), apps, conns, channels, dispatcher)
	handler := reverb.NewHTTPHandler(apps, conns, channels, dispatcher)

	mux := http.NewServeMux()
	mux.Handle("/app/", server)
	mux.Handle("/apps/", handler)

	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	return ts, app, apps, conns, channels, dispatcher
}

func dialWebSocket(t *testing.T, ts *httptest.Server, path, origin string) (*websocket.Conn, context.Context) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	t.Cleanup(cancel)

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + path
	conn, _, err := websocket.Dial(ctx, wsURL, &websocket.DialOptions{
		HTTPHeader: http.Header{
			"Origin": {origin},
		},
	})

	if err != nil {
		t.Fatalf("dial websocket %s: %v", wsURL, err)
	}

	t.Cleanup(func() {
		_ = conn.Close(websocket.StatusNormalClosure, "")
	})

	return conn, ctx
}

func readEnvelope(t *testing.T, ctx context.Context, conn *websocket.Conn) wsEnvelope {
	t.Helper()

	_, raw, err := conn.Read(ctx)

	if err != nil {
		t.Fatalf("read websocket frame: %v", err)
	}

	var env wsEnvelope

	if err := json.Unmarshal(raw, &env); err != nil {
		t.Fatalf("decode websocket frame: %v", err)
	}

	return env
}

func writeJSONMessage(t *testing.T, ctx context.Context, conn *websocket.Conn, payload any) {
	t.Helper()

	raw, err := json.Marshal(payload)

	if err != nil {
		t.Fatalf("encode websocket frame: %v", err)
	}

	if err := conn.Write(ctx, websocket.MessageText, raw); err != nil {
		t.Fatalf("write websocket frame: %v", err)
	}
}

func writeTextFrame(t *testing.T, ctx context.Context, conn *websocket.Conn, payload string) {
	t.Helper()

	if err := conn.Write(ctx, websocket.MessageText, []byte(payload)); err != nil {
		t.Fatalf("write websocket frame: %v", err)
	}
}

func assertEvent(t *testing.T, env wsEnvelope, event string) {
	t.Helper()

	if env.Event != event {
		t.Fatalf("event = %q, want %q (payload=%+v)", env.Event, event, env)
	}
}

func assertErrorCode(t *testing.T, env wsEnvelope, code string) {
	t.Helper()

	assertEvent(t, env, "pusher:error")

	if !strings.Contains(env.Data, `"code":`+code) {
		t.Fatalf("error payload = %q, want code %s", env.Data, code)
	}
}

func assertNoFrame(t *testing.T, conn *websocket.Conn) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 75*time.Millisecond)

	defer cancel()

	_, _, err := conn.Read(ctx)

	if err == nil {
		t.Fatal("expected no websocket frame")
	}
}

func assertEventsUnordered(t *testing.T, ctx context.Context, conn *websocket.Conn, want ...string) {
	t.Helper()

	if len(want) == 0 {
		t.Fatal("expected at least one event")
	}

	seen := make(map[string]int, len(want))

	for range want {
		env := readEnvelope(t, ctx, conn)
		seen[env.Event]++
	}

	for _, event := range want {
		if seen[event] == 0 {
			t.Fatalf("expected event %q in %v", event, seen)
		}
	}
}

func subscribeFrame(channel string, auth string, channelData string) map[string]any {
	payload := map[string]any{
		"channel": channel,
	}

	if auth != "" {
		payload["auth"] = auth
	}

	if channelData != "" {
		payload["channel_data"] = channelData
	}

	return map[string]any{
		"event": "pusher:subscribe",
		"data":  payload,
	}
}

func readEstablishedSocketID(t *testing.T, ctx context.Context, conn *websocket.Conn) string {
	t.Helper()

	env := readEnvelope(t, ctx, conn)
	assertEvent(t, env, "pusher:connection_established")

	var established reverb.ConnectionEstablishedData

	if err := json.Unmarshal([]byte(env.Data), &established); err != nil {
		t.Fatalf("decode connection established payload: %v", err)
	}

	if established.SocketID == "" {
		t.Fatal("expected non-empty socket id")
	}

	return established.SocketID
}

func waitForChannelRemoval(t *testing.T, channels *reverb.ChannelManager, appID, name string) {
	t.Helper()

	deadline := time.Now().Add(2 * time.Second)

	for time.Now().Before(deadline) {
		if _, ok := channels.Get(appID, name); !ok {
			return
		}

		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("expected channel %q to be removed", name)
}

// ClientEventTest::it_can_forward_a_client_message
// ClientEventTest::it_can_forward_an_unauthenticated_client_message_on_public_channel
// ClientEventTest::it_forwards_a_client_message_for_unauthenticated_client_when_set_to_all
// ClientEventTest::it_does_not_forward_client_message_when_set_to_none
// ClientEventTest::it_does_not_forward_unauthenticated_client_message_when_in_members_mode
// ClientEventTest::it_does_not_forward_a_message_to_itself
// ServerTest::it_allow_receiving_client_event_with_empty_data
func TestInventoryServerClientEventPolicies(t *testing.T) {
	t.Parallel()

	t.Run("all forwards public client events with empty data", func(t *testing.T) {
		t.Parallel()

		ts, app, _, _ := newWebSocketHarness(t, reverb.AppConfig{
			ID:             "app-1",
			Key:            "key-1",
			Secret:         "secret-1",
			AllowedOrigins: []string{"*"},
			ClientEvents:   reverb.ClientEventsConfig{Mode: "all"},
		})

		sender, senderCtx := dialWebSocket(t, ts, "/app/"+app.Key(), "https://app.example.com")
		observer, observerCtx := dialWebSocket(t, ts, "/app/"+app.Key(), "https://app.example.com")
		readEstablishedSocketID(t, senderCtx, sender)
		readEstablishedSocketID(t, observerCtx, observer)

		channelName := "public-client-events"
		writeJSONMessage(t, senderCtx, sender, subscribeFrame(channelName, "", ""))
		writeJSONMessage(t, observerCtx, observer, subscribeFrame(channelName, "", ""))
		assertEvent(t, readEnvelope(t, senderCtx, sender), "pusher_internal:subscription_succeeded")
		assertEvent(t, readEnvelope(t, observerCtx, observer), "pusher_internal:subscription_succeeded")

		writeJSONMessage(t, senderCtx, sender, map[string]any{
			"event":   "client-empty",
			"channel": channelName,
			"data":    map[string]any{},
		})

		assertEvent(t, readEnvelope(t, observerCtx, observer), "client-empty")
		assertNoFrame(t, sender)
	})

	t.Run("none rejects client events", func(t *testing.T) {
		t.Parallel()

		ts, app, _, _ := newWebSocketHarness(t, reverb.AppConfig{
			ID:             "app-1",
			Key:            "key-1",
			Secret:         "secret-1",
			AllowedOrigins: []string{"*"},
			ClientEvents:   reverb.ClientEventsConfig{Mode: "none"},
		})

		sender, senderCtx := dialWebSocket(t, ts, "/app/"+app.Key(), "https://app.example.com")
		observer, observerCtx := dialWebSocket(t, ts, "/app/"+app.Key(), "https://app.example.com")
		readEstablishedSocketID(t, senderCtx, sender)
		readEstablishedSocketID(t, observerCtx, observer)

		channelName := "public-client-events"
		writeJSONMessage(t, senderCtx, sender, subscribeFrame(channelName, "", ""))
		writeJSONMessage(t, observerCtx, observer, subscribeFrame(channelName, "", ""))
		assertEvent(t, readEnvelope(t, senderCtx, sender), "pusher_internal:subscription_succeeded")
		assertEvent(t, readEnvelope(t, observerCtx, observer), "pusher_internal:subscription_succeeded")

		writeJSONMessage(t, senderCtx, sender, map[string]any{
			"event":   "client-empty",
			"channel": channelName,
			"data":    map[string]any{},
		})

		assertErrorCode(t, readEnvelope(t, senderCtx, sender), "4301")
		assertNoFrame(t, observer)
	})

	t.Run("members rejects public channel client events", func(t *testing.T) {
		t.Parallel()

		ts, app, _, _ := newWebSocketHarness(t, reverb.AppConfig{
			ID:             "app-1",
			Key:            "key-1",
			Secret:         "secret-1",
			AllowedOrigins: []string{"*"},
			ClientEvents:   reverb.ClientEventsConfig{Mode: "members"},
		})

		sender, senderCtx := dialWebSocket(t, ts, "/app/"+app.Key(), "https://app.example.com")
		observer, observerCtx := dialWebSocket(t, ts, "/app/"+app.Key(), "https://app.example.com")
		readEstablishedSocketID(t, senderCtx, sender)
		readEstablishedSocketID(t, observerCtx, observer)

		channelName := "public-client-events"
		writeJSONMessage(t, senderCtx, sender, subscribeFrame(channelName, "", ""))
		writeJSONMessage(t, observerCtx, observer, subscribeFrame(channelName, "", ""))
		assertEvent(t, readEnvelope(t, senderCtx, sender), "pusher_internal:subscription_succeeded")
		assertEvent(t, readEnvelope(t, observerCtx, observer), "pusher_internal:subscription_succeeded")

		writeJSONMessage(t, senderCtx, sender, map[string]any{
			"event":   "client-empty",
			"channel": channelName,
			"data":    map[string]any{},
		})

		assertErrorCode(t, readEnvelope(t, senderCtx, sender), "4009")
		assertNoFrame(t, observer)
	})
}

// ClientEventTest::it_fails_on_unsupported_message
// ServerTest::it_sends_an_error_if_something_fails
// ServerTest::it_sends_an_error_if_something_fails_for_event_type
// ServerTest::it_sends_an_error_if_something_fails_for_data_type
// ServerTest::it_sends_an_error_if_something_fails_for_data_channel_type
// ServerTest::it_sends_an_error_if_something_fails_for_data_auth_type
// ServerTest::it_sends_an_error_if_something_fails_for_data_channel_data_type
// ServerTest::it_sends_an_error_if_something_fails_for_channel_type
func TestInventoryServerReportsMalformedMessages(t *testing.T) {
	t.Parallel()

	ts, app, _, _ := newWebSocketHarness(t, reverb.AppConfig{
		ID:             "app-1",
		Key:            "key-1",
		Secret:         "secret-1",
		AllowedOrigins: []string{"*"},
		ClientEvents:   reverb.ClientEventsConfig{Mode: "all"},
	})

	conn, ctx := dialWebSocket(t, ts, "/app/"+app.Key(), "https://app.example.com")
	readEstablishedSocketID(t, ctx, conn)

	cases := []string{
		`{"event":42,"data":{}}`,
		`{"event":"pusher:subscribe","data":42}`,
		`{"event":"pusher:subscribe","data":{"channel":42}}`,
		`{"event":"pusher:subscribe","data":{"channel":"private-room","auth":42}}`,
		`{"event":"pusher:subscribe","data":{"channel":"presence-room","auth":"bad","channel_data":{}}}`,
		`{"event":"client-bad","channel":42,"data":{}}`,
		`{"event":"client-bad","data":{}}`,
	}

	for _, payload := range cases {
		writeTextFrame(t, ctx, conn, payload)
		assertEvent(t, readEnvelope(t, ctx, conn), "pusher:error")
	}
}

// ServerTest::it_cannot_connect_when_over_the_max_connection_limit
// ServerTest::it_limits_the_size_of_messages
// ServerTest::it_rejects_messages_over_the_max_allowed_size
// ServerTest::it_allows_message_within_the_max_allowed_size
func TestInventoryServerEnforcesConnectionAndMessageLimits(t *testing.T) {
	t.Parallel()

	ts, app, _, _ := newWebSocketHarness(t, reverb.AppConfig{
		ID:              "app-1",
		Key:             "key-1",
		Secret:          "secret-1",
		MaxConnections:  1,
		MaxMessageSize:  64,
		AllowedOrigins:  []string{"*"},
		ActivityTimeout: 30,
	})

	first, firstCtx := dialWebSocket(t, ts, "/app/"+app.Key(), "https://app.example.com")
	readEstablishedSocketID(t, firstCtx, first)

	limitCtx, cancelLimit := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelLimit()

	_, resp, err := websocket.Dial(limitCtx, "ws"+strings.TrimPrefix(ts.URL, "http")+"/app/"+app.Key(), &websocket.DialOptions{
		HTTPHeader: http.Header{"Origin": {"https://app.example.com"}},
	})

	if err == nil {
		t.Fatal("expected second websocket dial to fail over max connection limit")
	}

	if resp == nil || resp.StatusCode != http.StatusServiceUnavailable {
		t.Fatalf("connection limit status = %v, want %d", resp, http.StatusServiceUnavailable)
	}

	writeTextFrame(t, firstCtx, first, `{"event":"pusher:ping","data":{}}`)
	assertEvent(t, readEnvelope(t, firstCtx, first), "pusher:pong")

	writeTextFrame(t, firstCtx, first, `{"event":"pusher:ping","data":{"padding":"`+strings.Repeat("x", 96)+`"}}`)

	if _, _, err := first.Read(firstCtx); err == nil {
		t.Fatal("expected oversized websocket message to close the connection")
	}
}

// ServerTest::it_it_can_ping_inactive_subscribers
// ServerTest::it_it_can_disconnect_inactive_subscribers
// PingInactiveConnectionsTest::it_pings_inactive_connections
// PruneStaleConnectionsTest::it_cleans_up_stale_connections
func TestInventoryServerInactiveConnectionJobs(t *testing.T) {
	t.Parallel()

	ts, app, apps, conns, _, _ := newWebSocketHarnessWithManagers(t, reverb.AppConfig{
		ID:              "app-1",
		Key:             "key-1",
		Secret:          "secret-1",
		AllowedOrigins:  []string{"*"},
		ActivityTimeout: 1,
		PingInterval:    1,
	})

	conn, ctx := dialWebSocket(t, ts, "/app/"+app.Key(), "https://app.example.com")
	readEstablishedSocketID(t, ctx, conn)

	time.Sleep(1100 * time.Millisecond)
	reverb.PingInactiveConnections(ctx, conns, apps)
	assertEvent(t, readEnvelope(t, ctx, conn), "pusher:ping")

	reverb.PruneStaleConnections(ctx, conns, apps)

	if _, _, err := conn.Read(ctx); err == nil {
		t.Fatal("expected stale websocket connection to be closed")
	}

	if conns.Count(app.ID()) != 0 {
		t.Fatalf("stale connection count = %d, want 0", conns.Count(app.ID()))
	}
}

// ServerTest::it_can_handle_a_connection
// EventTest::it_can_send_an_acknowledgement
// ServerTest::it_can_handle_a_new_message
// ServerTest::it_can_respond_to_a_ping
// ServerTest::it_can_handle_a_pong_control_frame
// ServerTest::it_can_subscribe_to_a_channel
// ServerTest::it_can_subscribe_to_a_private_channel
// ServerTest::it_can_handle_an_event
// ServerTest::it_can_receive_a_message_broadcast_from_the_server
// ServerTest::it_can_publish_and_subscribe_to_a_triggered_event
// ServerTest::it_can_handle_a_client_whisper
// ServerTest::it_can_publish_and_subscribe_to_a_client_whisper
// ServerTest::it_does_not_forward_a_message_to_itself
// ServerTest::it_can_handle_a_disconnection
func TestInventoryServerPublicWebSocketFlow(t *testing.T) {
	t.Parallel()

	ts, app, channels, dispatcher := newWebSocketHarness(t, reverb.AppConfig{
		ID:             "app-1",
		Key:            "key-1",
		Secret:         "secret-1",
		AllowedOrigins: []string{"*"},
		ClientEvents:   reverb.ClientEventsConfig{Mode: "all"},
	})

	sender, senderCtx := dialWebSocket(t, ts, "/app/"+app.Key(), "https://app.example.com")
	observer, observerCtx := dialWebSocket(t, ts, "/app/"+app.Key(), "https://app.example.com")

	senderSocketID := readEstablishedSocketID(t, senderCtx, sender)
	observerSocketID := readEstablishedSocketID(t, observerCtx, observer)

	if senderSocketID == observerSocketID {
		t.Fatalf("expected unique socket ids, got %q", senderSocketID)
	}

	channelName := "public-room"
	writeJSONMessage(t, senderCtx, sender, subscribeFrame(channelName, "", ""))
	writeJSONMessage(t, observerCtx, observer, subscribeFrame(channelName, "", ""))

	assertEvent(t, readEnvelope(t, senderCtx, sender), "pusher_internal:subscription_succeeded")
	assertEvent(t, readEnvelope(t, observerCtx, observer), "pusher_internal:subscription_succeeded")

	if err := dispatcher.Dispatch(senderCtx, app.ID(), contractsReverb.Event{
		Event:   "server-broadcast",
		Data:    `{"from":"dispatcher"}`,
		Channel: channelName,
	}); err != nil {
		t.Fatalf("dispatch broadcast: %v", err)
	}

	assertEvent(t, readEnvelope(t, senderCtx, sender), "server-broadcast")
	assertEvent(t, readEnvelope(t, observerCtx, observer), "server-broadcast")

	triggerBody, _ := json.Marshal(reverb.TriggerRequest{
		Name:     "triggered-event",
		Data:     `{"from":"http"}`,
		Channels: []string{channelName},
	})
	req := signedRequest(t, http.MethodPost, "/apps/"+app.ID()+"/events", app.Secret(), triggerBody)
	rec := httptest.NewRecorder()
	ts.Config.Handler.ServeHTTP(rec, req)

	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("trigger status = %d", rec.Result().StatusCode)
	}

	assertEvent(t, readEnvelope(t, senderCtx, sender), "triggered-event")
	assertEvent(t, readEnvelope(t, observerCtx, observer), "triggered-event")

	writeJSONMessage(t, senderCtx, sender, map[string]any{
		"event":   "client-whisper",
		"channel": channelName,
		"data":    map[string]any{"message": "hello"},
	})

	assertEvent(t, readEnvelope(t, observerCtx, observer), "client-whisper")

	writeJSONMessage(t, senderCtx, sender, map[string]any{
		"event":   "pusher:ping",
		"channel": channelName,
		"data":    map[string]any{},
	})
	assertEvent(t, readEnvelope(t, senderCtx, sender), "pusher:pong")

	writeJSONMessage(t, senderCtx, sender, map[string]any{
		"event":   "pusher:pong",
		"channel": channelName,
		"data":    map[string]any{},
	})

	assertNoFrame(t, sender)

	_ = sender.Close(websocket.StatusNormalClosure, "done")
	_ = observer.Close(websocket.StatusNormalClosure, "done")

	waitForChannelRemoval(t, channels, app.ID(), channelName)
}

// ServerTest::it_can_handle_connections_to_different_applications
// ServerTest::it_can_subscribe_a_connection_to_multiple_channels
// ServerTest::it_can_subscribe_multiple_connections_to_multiple_channels
func TestInventoryServerMultiAppAndMultiChannelFlow(t *testing.T) {
	t.Parallel()

	app1 := reverb.AppConfig{
		ID:             "app-1",
		Key:            "key-1",
		Secret:         "secret-1",
		AllowedOrigins: []string{"*"},
		ClientEvents:   reverb.ClientEventsConfig{Mode: "all"},
	}
	app2 := reverb.AppConfig{
		ID:             "app-2",
		Key:            "key-2",
		Secret:         "secret-2",
		AllowedOrigins: []string{"*"},
		ClientEvents:   reverb.ClientEventsConfig{Mode: "all"},
	}

	apps := reverb.NewAppManager([]reverb.AppConfig{app1, app2})
	channels := reverb.NewChannelManager(apps)
	conns := reverb.NewConnectionManager()
	dispatcher := reverb.NewSyncDispatcher(channels)
	server := reverb.NewServer(reverb.DefaultConfig(), apps, conns, channels, dispatcher)
	handler := reverb.NewHTTPHandler(apps, conns, channels, dispatcher)

	mux := http.NewServeMux()
	mux.Handle("/app/", server)
	mux.Handle("/apps/", handler)

	ts := httptest.NewServer(mux)
	t.Cleanup(ts.Close)

	alice, aliceCtx := dialWebSocket(t, ts, "/app/"+app1.Key, "https://app.example.com")
	bob, bobCtx := dialWebSocket(t, ts, "/app/"+app1.Key, "https://app.example.com")
	charlie, charlieCtx := dialWebSocket(t, ts, "/app/"+app2.Key, "https://app.example.com")

	aliceSocketID := readEstablishedSocketID(t, aliceCtx, alice)
	bobSocketID := readEstablishedSocketID(t, bobCtx, bob)
	charlieSocketID := readEstablishedSocketID(t, charlieCtx, charlie)

	if aliceSocketID == bobSocketID || aliceSocketID == charlieSocketID || bobSocketID == charlieSocketID {
		t.Fatalf("expected unique socket ids, got %q %q %q", aliceSocketID, bobSocketID, charlieSocketID)
	}

	channelOne := "public-one"
	channelTwo := "public-two"

	for _, tc := range []struct {
		conn *websocket.Conn
		ctx  context.Context
	}{
		{alice, aliceCtx},
		{bob, bobCtx},
	} {
		writeJSONMessage(t, tc.ctx, tc.conn, subscribeFrame(channelOne, "", ""))
		writeJSONMessage(t, tc.ctx, tc.conn, subscribeFrame(channelTwo, "", ""))
		assertEvent(t, readEnvelope(t, tc.ctx, tc.conn), "pusher_internal:subscription_succeeded")
		assertEvent(t, readEnvelope(t, tc.ctx, tc.conn), "pusher_internal:subscription_succeeded")
	}

	writeJSONMessage(t, charlieCtx, charlie, subscribeFrame(channelOne, "", ""))
	assertEvent(t, readEnvelope(t, charlieCtx, charlie), "pusher_internal:subscription_succeeded")

	if err := dispatcher.Dispatch(aliceCtx, app1.ID, contractsReverb.Event{Event: "app1-first", Data: `{}`, Channel: channelOne}); err != nil {
		t.Fatalf("dispatch app1 channelOne: %v", err)
	}

	assertEvent(t, readEnvelope(t, aliceCtx, alice), "app1-first")
	assertEvent(t, readEnvelope(t, bobCtx, bob), "app1-first")

	if err := dispatcher.Dispatch(aliceCtx, app1.ID, contractsReverb.Event{Event: "app1-second", Data: `{}`, Channel: channelTwo}); err != nil {
		t.Fatalf("dispatch app1 channelTwo: %v", err)
	}

	assertEvent(t, readEnvelope(t, aliceCtx, alice), "app1-second")
	assertEvent(t, readEnvelope(t, bobCtx, bob), "app1-second")

	if err := dispatcher.Dispatch(charlieCtx, app2.ID, contractsReverb.Event{Event: "app2-first", Data: `{}`, Channel: channelOne}); err != nil {
		t.Fatalf("dispatch app2 channelOne: %v", err)
	}

	assertEvent(t, readEnvelope(t, charlieCtx, charlie), "app2-first")
	assertNoFrame(t, alice)
	assertNoFrame(t, bob)
}

// ServerTest::it_can_subscribe_to_a_presence_channel
// ServerTest::it_can_notify_other_subscribers_of_a_presence_channel_when_a_new_member_joins
// ServerTest::it_can_notify_other_subscribers_of_a_presence_channel_when_a_member_leaves
// ServerTest::it_subscription_succeeded_event_contains_unique_list_of_users
// ServerTest::it_can_subscribe_to_a_cache_channel
// ServerTest::it_can_receive_a_cache_missed_message_when_joining_a_cache_channel_with_an_empty_cache
// ServerTest::it_can_receive_a_cached_message_when_joining_a_cache_channel
// ServerTest::it_can_subscribe_to_a_private_cache_channel
// ServerTest::it_can_receive_a_cache_missed_message_when_joining_a_private_cache_channel_with_an_empty_cache
// ServerTest::it_can_receive_a_cached_message_when_joining_a_private_cache_channel
// ServerTest::it_can_subscribe_to_a_presence_cache_channel
// ServerTest::it_can_receive_a_cache_missed_message_when_joining_a_presence_cache_channel_with_an_empty_cache
// ServerTest::it_can_receive_a_cached_message_when_joining_a_presence_cache_channel
func TestInventoryServerPresenceAndCacheJoinFlow(t *testing.T) {
	t.Parallel()

	ts, app, channels, dispatcher := newWebSocketHarness(t, reverb.AppConfig{
		ID:             "app-1",
		Key:            "key-1",
		Secret:         "secret-1",
		AllowedOrigins: []string{"*"},
		ClientEvents:   reverb.ClientEventsConfig{Mode: "all"},
	})

	alice, aliceCtx := dialWebSocket(t, ts, "/app/"+app.Key(), "https://app.example.com")
	bob, bobCtx := dialWebSocket(t, ts, "/app/"+app.Key(), "https://app.example.com")
	aliceSocketID := readEstablishedSocketID(t, aliceCtx, alice)
	bobSocketID := readEstablishedSocketID(t, bobCtx, bob)

	presenceName := "presence-room"
	aliceData := `{"user_id":"alice","user_info":{"name":"Alice"}}`
	bobData := `{"user_id":"bob","user_info":{"name":"Bob"}}`

	writeJSONMessage(t, aliceCtx, alice, subscribeFrame(
		presenceName,
		signedChannelAuth(app, aliceSocketID, presenceName, aliceData),
		aliceData,
	))
	env := readEnvelope(t, aliceCtx, alice)
	assertEvent(t, env, "pusher_internal:subscription_succeeded")

	if !strings.Contains(env.Data, `"count":1`) {
		t.Fatalf("alice presence payload = %q", env.Data)
	}

	writeJSONMessage(t, bobCtx, bob, subscribeFrame(
		presenceName,
		signedChannelAuth(app, bobSocketID, presenceName, bobData),
		bobData,
	))
	assertEvent(t, readEnvelope(t, aliceCtx, alice), "pusher_internal:member_added")
	env = readEnvelope(t, bobCtx, bob)
	assertEvent(t, env, "pusher_internal:subscription_succeeded")

	if !strings.Contains(env.Data, `"count":2`) || !strings.Contains(env.Data, `"Alice"`) || !strings.Contains(env.Data, `"Bob"`) {
		t.Fatalf("bob presence payload = %q", env.Data)
	}

	_ = bob.Close(websocket.StatusNormalClosure, "done")
	assertEvent(t, readEnvelope(t, aliceCtx, alice), "pusher_internal:member_removed")

	cacheName := "cache-room"
	writeJSONMessage(t, aliceCtx, alice, subscribeFrame(cacheName, "", ""))
	assertEventsUnordered(t, aliceCtx, alice, "pusher_internal:subscription_succeeded", "pusher:cache_miss")

	if err := dispatcher.Dispatch(aliceCtx, app.ID(), contractsReverb.Event{
		Event:   "cache-event",
		Data:    `{"cached":true}`,
		Channel: cacheName,
	}); err != nil {
		t.Fatalf("dispatch cache event: %v", err)
	}

	assertEvent(t, readEnvelope(t, aliceCtx, alice), "cache-event")

	cacheObserver, cacheObserverCtx := dialWebSocket(t, ts, "/app/"+app.Key(), "https://app.example.com")
	cacheObserverSocketID := readEstablishedSocketID(t, cacheObserverCtx, cacheObserver)
	writeJSONMessage(t, cacheObserverCtx, cacheObserver, subscribeFrame(cacheName, "", ""))
	assertEventsUnordered(t, cacheObserverCtx, cacheObserver, "pusher_internal:subscription_succeeded", "cache-event")

	privateCacheName := "private-cache-room"
	writeJSONMessage(t, aliceCtx, alice, subscribeFrame(
		privateCacheName,
		signedChannelAuth(app, aliceSocketID, privateCacheName, ""),
		"",
	))
	assertEventsUnordered(t, aliceCtx, alice, "pusher_internal:subscription_succeeded", "pusher:cache_miss")

	if err := dispatcher.Dispatch(aliceCtx, app.ID(), contractsReverb.Event{
		Event:   "private-cache-event",
		Data:    `{"cached":true}`,
		Channel: privateCacheName,
	}); err != nil {
		t.Fatalf("dispatch private cache event: %v", err)
	}

	assertEvent(t, readEnvelope(t, aliceCtx, alice), "private-cache-event")

	writeJSONMessage(t, cacheObserverCtx, cacheObserver, subscribeFrame(
		privateCacheName,
		signedChannelAuth(app, cacheObserverSocketID, privateCacheName, ""),
		"",
	))
	assertEventsUnordered(t, cacheObserverCtx, cacheObserver, "pusher_internal:subscription_succeeded", "private-cache-event")

	presenceCacheName := "presence-cache-room"
	alicePresenceData := `{"user_id":"alice","user_info":{"name":"Alice"}}`
	writeJSONMessage(t, aliceCtx, alice, subscribeFrame(
		presenceCacheName,
		signedChannelAuth(app, aliceSocketID, presenceCacheName, alicePresenceData),
		alicePresenceData,
	))
	assertEventsUnordered(t, aliceCtx, alice, "pusher_internal:subscription_succeeded", "pusher:cache_miss")

	if err := dispatcher.Dispatch(aliceCtx, app.ID(), contractsReverb.Event{
		Event:   "presence-cache-event",
		Data:    `{"cached":true}`,
		Channel: presenceCacheName,
	}); err != nil {
		t.Fatalf("dispatch presence cache event: %v", err)
	}

	assertEvent(t, readEnvelope(t, aliceCtx, alice), "presence-cache-event")

	writeJSONMessage(t, cacheObserverCtx, cacheObserver, subscribeFrame(
		presenceCacheName,
		signedChannelAuth(app, cacheObserverSocketID, presenceCacheName, bobData),
		bobData,
	))
	assertEvent(t, readEnvelope(t, aliceCtx, alice), "pusher_internal:member_added")
	assertEventsUnordered(t, cacheObserverCtx, cacheObserver, "pusher_internal:subscription_succeeded", "presence-cache-event")

	_ = cacheObserver.Close(websocket.StatusNormalClosure, "done")
	assertEvent(t, readEnvelope(t, aliceCtx, alice), "pusher_internal:member_removed")

	_ = alice.Close(websocket.StatusNormalClosure, "done")

	waitForChannelRemoval(t, channels, app.ID(), presenceName)
	waitForChannelRemoval(t, channels, app.ID(), cacheName)
	waitForChannelRemoval(t, channels, app.ID(), privateCacheName)
	waitForChannelRemoval(t, channels, app.ID(), presenceCacheName)
}

// ServerTest::it_fails_to_connect_when_an_invalid_application_is_provided
// ServerTest::it_fails_to_subscribe_to_a_private_channel_with_no_auth_token
// ServerTest::it_fails_to_subscribe_to_a_presence_channel_with_no_auth_token
// ServerTest::it_fails_to_subscribe_to_a_private_channel_with_invalid_auth_signature
// ServerTest::it_fails_to_subscribe_to_a_presence_channel_with_invalid_auth_signature
// ServerTest::it_fails_to_subscribe_to_a_private_cache_channel_with_invalid_auth_signature
// ServerTest::it_fails_to_subscribe_to_a_presence_cache_channel_with_invalid_auth_signature
func TestInventoryServerRejectsInvalidApplicationsAndAuth(t *testing.T) {
	t.Parallel()

	ts, app, _, _ := newWebSocketHarness(t, reverb.AppConfig{
		ID:             "app-1",
		Key:            "key-1",
		Secret:         "secret-1",
		AllowedOrigins: []string{"*"},
	})

	badReq := httptest.NewRequest(http.MethodGet, "/app/missing", nil)
	badReq.Header.Set("Origin", "https://app.example.com")
	badRec := httptest.NewRecorder()
	ts.Config.Handler.ServeHTTP(badRec, badReq)

	if badRec.Result().StatusCode != http.StatusNotFound {
		t.Fatalf("invalid app status = %d", badRec.Result().StatusCode)
	}

	conn, ctx := dialWebSocket(t, ts, "/app/"+app.Key(), "https://app.example.com")
	socketID := readEstablishedSocketID(t, ctx, conn)

	readAndAssertError := func(name string, payload map[string]any) {
		t.Helper()

		writeJSONMessage(t, ctx, conn, payload)
		env := readEnvelope(t, ctx, conn)
		assertEvent(t, env, "pusher:error")

		if !strings.Contains(env.Data, `"code":4009`) {
			t.Fatalf("%s error payload = %q", name, env.Data)
		}
	}

	readAndAssertError("private missing auth", subscribeFrame("private-room", "", ""))
	readAndAssertError("presence missing auth", subscribeFrame("presence-room", "", `{"user_id":"u1"}`))
	readAndAssertError("private invalid auth", subscribeFrame("private-room", "bad-signature", ""))
	readAndAssertError("presence invalid auth", subscribeFrame("presence-room", "bad-signature", `{"user_id":"u1"}`))
	readAndAssertError("private cache invalid auth", subscribeFrame("private-cache-room", "bad-signature", ""))
	readAndAssertError("presence cache invalid auth", subscribeFrame("presence-cache-room", "bad-signature", `{"user_id":"u1"}`))

	if socketID == "" {
		t.Fatal("expected socket id")
	}
}

// ServerTest::it_can_handle_a_ping_control_frame
func TestInventoryServerHandlesPingControlFrame(t *testing.T) {
	t.Parallel()

	ts, app, _, _ := newWebSocketHarness(t, reverb.AppConfig{
		ID:             "app-1",
		Key:            "key-1",
		Secret:         "secret-1",
		AllowedOrigins: []string{"*"},
	})

	conn, ctx := dialWebSocket(t, ts, "/app/"+app.Key(), "https://app.example.com")
	readEstablishedSocketID(t, ctx, conn)

	go func() {
		_, _, _ = conn.Read(context.Background())
	}()

	if err := conn.Ping(ctx); err != nil {
		t.Fatalf("ping control frame: %v", err)
	}
}
