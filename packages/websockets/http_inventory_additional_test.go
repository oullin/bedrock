package reverb_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/bedrock/packages/websockets"
)

func signedRequestWithExtraParams(t *testing.T, method, path, secret string, body []byte, extra map[string]string) *http.Request {
	t.Helper()

	params := map[string]string{
		"auth_key":       "key-1",
		"auth_timestamp": "1234567890",
		"auth_version":   "1.0",
	}
	for k, v := range extra {
		params[k] = v
	}

	sig := websockets.SignHTTPRequest(secret, method, path, params)
	params["auth_signature"] = sig

	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}

	reader := bytes.NewReader(body)
	req := httptest.NewRequest(method, path+"?"+values.Encode(), reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req
}

func responseBodyString(t *testing.T, resp *http.Response) string {
	t.Helper()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}

	return string(body)
}

// HealthCheckControllerTest::it_can_respond_to_a_health_check_request
// ConnectionsControllerTest::it_can_return_a_connection_count
// ConnectionsControllerTest::it_can_return_the_correct_connection_count_when_subscribed_to_multiple_channels
// EventsControllerTest::it_can_verify_signature_when_using_a_custom_server_path
// EventsControllerTest::it_can_send_the_content_length_header
func TestInventoryHTTPHealthConnectionsAndCustomPath(t *testing.T) {
	t.Parallel()

	handler, _, conns, mgr := newHTTPTestSetup(t)

	healthReq := httptest.NewRequest(http.MethodGet, "/up", nil)
	healthRec := httptest.NewRecorder()
	handler.ServeHTTP(healthRec, healthReq)
	if healthRec.Result().StatusCode != http.StatusOK {
		t.Fatalf("health status = %d", healthRec.Result().StatusCode)
	}
	if got := responseBodyString(t, healthRec.Result()); got != `{"health":"OK"}` {
		t.Fatalf("health body = %s", got)
	}
	if got := healthRec.Result().Header.Get("Content-Length"); got != "15" {
		t.Fatalf("health content-length = %s", got)
	}

	conns.Add(websockets.NewConn(nil, "app-1"))
	conns.Add(websockets.NewConn(nil, "app-1"))
	countReq := signedRequest(t, http.MethodGet, "/apps/app-1/connections", "secret-1", nil)
	countRec := httptest.NewRecorder()
	handler.ServeHTTP(countRec, countReq)
	if countRec.Result().StatusCode != http.StatusOK {
		t.Fatalf("connections status = %d", countRec.Result().StatusCode)
	}
	if got := responseBodyString(t, countRec.Result()); got != `{"connections":2}` {
		t.Fatalf("connections body = %s", got)
	}
	if got := countRec.Result().Header.Get("Content-Length"); got != "17" {
		t.Fatalf("connections content-length = %s", got)
	}

	ch, err := mgr.GetOrCreate("app-1", "public-prefix")
	if err != nil {
		t.Fatalf("GetOrCreate: %v", err)
	}
	receiver := newFakeConn("socket-1", "app-1")
	if err := ch.Subscribe(context.Background(), receiver, "", ""); err != nil {
		t.Fatalf("subscribe: %v", err)
	}

	body := []byte(`{"name":"custom-path-event","data":"{}","channels":["public-prefix"]}`)
	params := map[string]string{
		"auth_key":       "key-1",
		"auth_timestamp": "1234567890",
		"auth_version":   "1.0",
	}
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	values.Set("auth_signature", websockets.SignHTTPRequest("secret-1", http.MethodPost, "/apps/app-1/events", params))
	req := httptest.NewRequest(http.MethodPost, "/ws/apps/app-1/events?"+values.Encode(), bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("custom path status = %d", rec.Result().StatusCode)
	}
	if got := responseBodyString(t, rec.Result()); got != `{}` {
		t.Fatalf("custom path body = %s", got)
	}
	if got := rec.Result().Header.Get("Content-Length"); got != "2" {
		t.Fatalf("custom path content-length = %s", got)
	}
	if !containsBytes(receiver.SentMessages(), []byte("custom-path-event")) {
		t.Fatalf("receiver did not get custom-path-event")
	}
}

// ChannelUsersControllerTest::it_returns_the_user_data
// ChannelUsersControllerTest::it_returns_the_unique_user_data
// ChannelUsersControllerTest::it_can_send_the_content_length_header
func TestInventoryHTTPChannelUsers(t *testing.T) {
	t.Parallel()

	handler, apps, _, mgr := newHTTPTestSetup(t)
	app, err := apps.FindByID("app-1")
	if err != nil {
		t.Fatalf("FindByID: %v", err)
	}

	channel, err := mgr.GetOrCreate(app.ID(), "presence-test-channel")
	if err != nil {
		t.Fatalf("GetOrCreate: %v", err)
	}

	first := newFakeConn("socket-1", app.ID())
	second := newFakeConn("socket-2", app.ID())
	third := newFakeConn("socket-3", app.ID())

	firstData := `{"user_id":1,"user_info":{"name":"Taylor"}}`
	secondData := `{"user_id":2,"user_info":{"name":"Joe"}}`
	thirdData := `{"user_id":3,"user_info":{"name":"Jess"}}`

	if err := channel.Subscribe(context.Background(), first, signedChannelAuth(app, first.SocketID(), "presence-test-channel", firstData), firstData); err != nil {
		t.Fatalf("subscribe first: %v", err)
	}
	if err := channel.Subscribe(context.Background(), second, signedChannelAuth(app, second.SocketID(), "presence-test-channel", secondData), secondData); err != nil {
		t.Fatalf("subscribe second: %v", err)
	}
	if err := channel.Subscribe(context.Background(), third, signedChannelAuth(app, third.SocketID(), "presence-test-channel", thirdData), thirdData); err != nil {
		t.Fatalf("subscribe third: %v", err)
	}

	req := signedRequest(t, http.MethodGet, "/apps/app-1/channels/presence-test-channel/users", "secret-1", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Result().StatusCode != http.StatusOK {
		t.Fatalf("users status = %d", rec.Result().StatusCode)
	}
	if got := responseBodyString(t, rec.Result()); got != `{"users":[{"id":1},{"id":2},{"id":3}]}` {
		t.Fatalf("users body = %s", got)
	}
	if got := rec.Result().Header.Get("Content-Length"); got != "38" {
		t.Fatalf("users content-length = %s", got)
	}

	dupHandler, dupApps, _, dupMgr := newHTTPTestSetup(t)
	dupApp, err := dupApps.FindByID("app-1")
	if err != nil {
		t.Fatalf("FindByID duplicate: %v", err)
	}
	dupChannel, err := dupMgr.GetOrCreate(dupApp.ID(), "presence-test-channel")
	if err != nil {
		t.Fatalf("GetOrCreate duplicate: %v", err)
	}
	dupOne := newFakeConn("socket-4", dupApp.ID())
	dupTwo := newFakeConn("socket-5", dupApp.ID())
	dupThree := newFakeConn("socket-6", dupApp.ID())
	dupA := `{"user_id":3,"user_info":{"name":"Taylor"}}`
	dupB := `{"user_id":2,"user_info":{"name":"Joe"}}`
	dupC := `{"user_id":3,"user_info":{"name":"Jess"}}`
	if err := dupChannel.Subscribe(context.Background(), dupOne, signedChannelAuth(dupApp, dupOne.SocketID(), "presence-test-channel", dupA), dupA); err != nil {
		t.Fatalf("subscribe dup one: %v", err)
	}
	if err := dupChannel.Subscribe(context.Background(), dupTwo, signedChannelAuth(dupApp, dupTwo.SocketID(), "presence-test-channel", dupB), dupB); err != nil {
		t.Fatalf("subscribe dup two: %v", err)
	}
	if err := dupChannel.Subscribe(context.Background(), dupThree, signedChannelAuth(dupApp, dupThree.SocketID(), "presence-test-channel", dupC), dupC); err != nil {
		t.Fatalf("subscribe dup three: %v", err)
	}

	dupReq := signedRequest(t, http.MethodGet, "/apps/app-1/channels/presence-test-channel/users", "secret-1", nil)
	dupRec := httptest.NewRecorder()
	dupHandler.ServeHTTP(dupRec, dupReq)
	if got := responseBodyString(t, dupRec.Result()); got != `{"users":[{"id":3},{"id":2}]}` {
		t.Fatalf("unique users body = %s", got)
	}
}

// EventsControllerTest::it_can_return_user_counts_when_requested
// EventsControllerTest::it_can_return_subscription_counts_when_requested
// EventsControllerTest::it_can_receive_and_event_trigger
// EventsControllerTest::it_can_receive_and_event_trigger_for_multiple_channels
// EventsBatchControllerTest::it_can_receive_an_event_batch_trigger
// EventsBatchControllerTest::it_can_receive_an_event_batch_trigger_with_multiple_events
func TestInventoryHTTPEventInfoResponses(t *testing.T) {
	t.Parallel()

	handler, _, _, mgr := newHTTPTestSetup(t)
	ctx := httptest.NewRequest(http.MethodPost, "/", nil).Context()

	presence, err := mgr.GetOrCreate("app-1", "presence-test-channel-one")
	if err != nil {
		t.Fatalf("GetOrCreate presence: %v", err)
	}
	public, err := mgr.GetOrCreate("app-1", "test-channel-two")
	if err != nil {
		t.Fatalf("GetOrCreate public: %v", err)
	}

	presenceConn := newFakeConn("socket-1", "app-1")
	publicConn := newFakeConn("socket-2", "app-1")

	presenceData := `{"user_id":1,"user_info":{"name":"Taylor"}}`
	app, err := websockets.NewAppManager([]websockets.AppConfig{{ID: "app-1", Key: "key-1", Secret: "secret-1"}}).FindByID("app-1")
	if err != nil {
		t.Fatalf("FindByID app: %v", err)
	}

	if err := presence.Subscribe(ctx, presenceConn, signedChannelAuth(app, presenceConn.SocketID(), "presence-test-channel-one", presenceData), presenceData); err != nil {
		t.Fatalf("subscribe presence: %v", err)
	}
	if err := public.Subscribe(ctx, publicConn, "", ""); err != nil {
		t.Fatalf("subscribe public: %v", err)
	}

	body := []byte(`{"name":"NewEvent","channels":["presence-test-channel-one","test-channel-two"],"data":"{\"some\":\"data\"}"}`)
	userCountReq := signedRequestWithExtraParams(t, http.MethodPost, "/apps/app-1/events", "secret-1", body, map[string]string{"info": "user_count"})
	userCountRec := httptest.NewRecorder()
	handler.ServeHTTP(userCountRec, userCountReq)
	if got := responseBodyString(t, userCountRec.Result()); got != `{"channels":{"presence-test-channel-one":{"user_count":1},"test-channel-two":{}}}` {
		t.Fatalf("user_count body = %s", got)
	}

	subCountReq := signedRequestWithExtraParams(t, http.MethodPost, "/apps/app-1/events", "secret-1", body, map[string]string{"info": "subscription_count"})
	subCountRec := httptest.NewRecorder()
	handler.ServeHTTP(subCountRec, subCountReq)
	if got := responseBodyString(t, subCountRec.Result()); got != `{"channels":{"presence-test-channel-one":{},"test-channel-two":{"subscription_count":1}}}` {
		t.Fatalf("subscription_count body = %s", got)
	}

	batchBody := []byte(`{"batch":[{"name":"NewEvent","channels":["presence-test-channel-one"],"data":"{}"},{"name":"NewEvent","channels":["test-channel-two"],"data":"{}"}]}`)
	batchReq := signedRequest(t, http.MethodPost, "/apps/app-1/batch_events", "secret-1", batchBody)
	batchRec := httptest.NewRecorder()
	handler.ServeHTTP(batchRec, batchReq)
	if got := responseBodyString(t, batchRec.Result()); got != `{}` {
		t.Fatalf("batch body = %s", got)
	}
}
