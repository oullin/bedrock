package websockets_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/websockets"
)

// signedRequest builds a signed HTTP request for the WebSockets HTTP API.
// It computes the HMAC signature over method + path + sorted query params
// and attaches it as auth_signature.
func signedRequest(t *testing.T, method, path, secret string, body []byte) *http.Request {
	t.Helper()

	params := map[string]string{
		"auth_key":       "key-1",
		"auth_timestamp": "1234567890",
		"auth_version":   "1.0",
	}

	sig := websockets.SignHTTPRequest(secret, method, path, params)

	// Build query string
	var qparts []string

	for k, v := range params {
		qparts = append(qparts, fmt.Sprintf("%s=%s", k, v))
	}

	qparts = append(qparts, "auth_signature="+sig)
	queryString := strings.Join(qparts, "&")

	fullURL := path + "?" + queryString

	var bodyReader *bytes.Reader

	if body != nil {
		bodyReader = bytes.NewReader(body)
	} else {
		bodyReader = bytes.NewReader([]byte{})
	}

	req := httptest.NewRequest(method, fullURL, bodyReader)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req
}

func newHTTPTestSetup(t *testing.T) (*websockets.HTTPHandler, *websockets.AppManager, *websockets.ConnectionManager, *websockets.ChannelManager) {
	t.Helper()

	apps := websockets.NewAppManager([]websockets.AppConfig{
		{ID: "app-1", Key: "key-1", Secret: "secret-1"},
	})
	conns := websockets.NewConnectionManager()
	mgr := websockets.NewChannelManager(apps)
	dispatcher := websockets.NewSyncDispatcher(mgr)
	handler := websockets.NewHTTPHandler(apps, conns, mgr, dispatcher)

	return handler, apps, conns, mgr
}

func TestHTTPHandler_Trigger_Success(t *testing.T) {
	t.Parallel()

	handler, _, _, mgr := newHTTPTestSetup(t)
	ctx := httptest.NewRequest(http.MethodPost, "/", nil).Context()

	// Create a channel and subscribe a fakeConn
	ch, err := mgr.GetOrCreate("app-1", "public-test")

	if err != nil {
		t.Fatalf("GetOrCreate: %v", err)
	}

	conn := newFakeConn("sock-1", "app-1")

	if err := ch.Subscribe(ctx, conn, "", ""); err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	msgsBeforeTrigger := len(conn.SentMessages())

	body, _ := json.Marshal(websockets.TriggerRequest{
		Name:     "triggered-event",
		Data:     `{"hello":"world"}`,
		Channels: []string{"public-test"},
	})

	path := "/apps/app-1/events"
	req := signedRequest(t, http.MethodPost, path, "secret-1", body)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	msgs := conn.SentMessages()

	if len(msgs) <= msgsBeforeTrigger {
		t.Error("expected fakeConn to receive the triggered event")
	}

	if !containsBytes(msgs[msgsBeforeTrigger:], []byte("triggered-event")) {
		t.Errorf("expected triggered-event in messages, got: %s", msgs[msgsBeforeTrigger:])
	}
}

func TestHTTPHandler_Trigger_InvalidSignature(t *testing.T) {
	t.Parallel()

	handler, _, _, _ := newHTTPTestSetup(t)

	body, _ := json.Marshal(websockets.TriggerRequest{
		Name:     "test-event",
		Data:     `{}`,
		Channels: []string{"public-test"},
	})

	path := "/apps/app-1/events"
	// Use wrong secret to produce invalid signature
	req := signedRequest(t, http.MethodPost, path, "wrong-secret", body)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}

func TestHTTPHandler_Trigger_AppNotFound(t *testing.T) {
	t.Parallel()

	handler, _, _, _ := newHTTPTestSetup(t)

	body, _ := json.Marshal(websockets.TriggerRequest{
		Name:     "test-event",
		Data:     `{}`,
		Channels: []string{"public-test"},
	})

	path := "/apps/nonexistent/events"
	req := signedRequest(t, http.MethodPost, path, "secret-1", body)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}

func TestHTTPHandler_BatchTrigger_Success(t *testing.T) {
	t.Parallel()

	handler, _, _, mgr := newHTTPTestSetup(t)
	ctx := httptest.NewRequest(http.MethodPost, "/", nil).Context()

	// Set up two channels
	ch1, err := mgr.GetOrCreate("app-1", "public-ch1")

	if err != nil {
		t.Fatalf("GetOrCreate ch1: %v", err)
	}

	ch2, err := mgr.GetOrCreate("app-1", "public-ch2")

	if err != nil {
		t.Fatalf("GetOrCreate ch2: %v", err)
	}

	conn1 := newFakeConn("sock-1", "app-1")
	conn2 := newFakeConn("sock-2", "app-1")

	if err := ch1.Subscribe(ctx, conn1, "", ""); err != nil {
		t.Fatalf("Subscribe conn1: %v", err)
	}

	if err := ch2.Subscribe(ctx, conn2, "", ""); err != nil {
		t.Fatalf("Subscribe conn2: %v", err)
	}

	msgs1Before := len(conn1.SentMessages())
	msgs2Before := len(conn2.SentMessages())

	batchBody, _ := json.Marshal(websockets.BatchTriggerRequest{
		Batch: []websockets.TriggerRequest{
			{Name: "event-for-ch1", Data: `{"a":1}`, Channels: []string{"public-ch1"}},
			{Name: "event-for-ch2", Data: `{"b":2}`, Channels: []string{"public-ch2"}},
		},
	})

	path := "/apps/app-1/batch_events"
	req := signedRequest(t, http.MethodPost, path, "secret-1", batchBody)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	msgs1 := conn1.SentMessages()
	msgs2 := conn2.SentMessages()

	if len(msgs1) <= msgs1Before {
		t.Error("expected conn1 to receive event-for-ch1")
	}

	if len(msgs2) <= msgs2Before {
		t.Error("expected conn2 to receive event-for-ch2")
	}

	if !containsBytes(msgs1[msgs1Before:], []byte("event-for-ch1")) {
		t.Errorf("expected event-for-ch1 in conn1 messages, got: %s", msgs1[msgs1Before:])
	}

	if !containsBytes(msgs2[msgs2Before:], []byte("event-for-ch2")) {
		t.Errorf("expected event-for-ch2 in conn2 messages, got: %s", msgs2[msgs2Before:])
	}
}

func TestHTTPHandler_GetChannels(t *testing.T) {
	t.Parallel()

	handler, _, _, mgr := newHTTPTestSetup(t)
	ctx := httptest.NewRequest(http.MethodGet, "/", nil).Context()

	// Create a channel with a subscriber
	ch, err := mgr.GetOrCreate("app-1", "public-test")

	if err != nil {
		t.Fatalf("GetOrCreate: %v", err)
	}

	conn := newFakeConn("sock-1", "app-1")

	if err := ch.Subscribe(ctx, conn, "", ""); err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	path := "/apps/app-1/channels"
	req := signedRequest(t, http.MethodGet, path, "secret-1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]any

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	channels, ok := result["channels"]

	if !ok {
		t.Error("expected 'channels' key in response JSON")
	}

	_ = channels
}

func TestHTTPHandler_GetChannel_Found(t *testing.T) {
	t.Parallel()

	handler, _, _, mgr := newHTTPTestSetup(t)
	ctx := httptest.NewRequest(http.MethodGet, "/", nil).Context()

	ch, err := mgr.GetOrCreate("app-1", "public-test")

	if err != nil {
		t.Fatalf("GetOrCreate: %v", err)
	}

	conn := newFakeConn("sock-1", "app-1")

	if err := ch.Subscribe(ctx, conn, "", ""); err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	path := "/apps/app-1/channels/public-test"
	req := signedRequest(t, http.MethodGet, path, "secret-1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestHTTPHandler_GetChannel_NotFound(t *testing.T) {
	t.Parallel()

	handler, _, _, _ := newHTTPTestSetup(t)

	path := "/apps/app-1/channels/nonexistent"
	req := signedRequest(t, http.MethodGet, path, "secret-1", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", resp.StatusCode)
	}
}
