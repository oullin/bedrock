package client_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/bedrock/packages/httpx/client"
)

func TestPendingRequestGet(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("expected GET, got %s", r.Method)
		}

		w.Write([]byte("ok"))
	}))

	defer server.Close()

	resp, err := client.NewFactory().PendingRequest().Get(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Body() != "ok" {
		t.Fatalf("expected ok, got %s", resp.Body())
	}
}

func TestPendingRequestGetWithQuery(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.URL.Query().Get("page")))
	}))

	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().Get(server.URL, map[string]string{"page": "2"})

	if resp.Body() != "2" {
		t.Fatalf("expected 2, got %s", resp.Body())
	}
}

func TestPendingRequestPostJSON(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatal("expected JSON content type")
		}

		body, _ := io.ReadAll(r.Body)

		var data map[string]string

		json.Unmarshal(body, &data)

		w.Write([]byte(data["name"]))
	}))

	defer server.Close()

	resp, err := client.NewFactory().PendingRequest().AsJSON().Post(server.URL, map[string]string{"name": "Taylor"})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Body() != "Taylor" {
		t.Fatalf("expected Taylor, got %s", resp.Body())
	}
}

func TestPendingRequestPostForm(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		w.Write([]byte(r.PostForm.Get("email")))
	}))

	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().AsForm().Post(server.URL, map[string]string{"email": "test@test.com"})

	if resp.Body() != "test@test.com" {
		t.Fatalf("expected test@test.com, got %s", resp.Body())
	}
}

func TestPendingRequestWithHeaders(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.Header.Get("X-Custom")))
	}))

	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().
		WithHeader("X-Custom", "hello").
		Get(server.URL)

	if resp.Body() != "hello" {
		t.Fatalf("expected hello, got %s", resp.Body())
	}
}

func TestPendingRequestWithToken(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.Header.Get("Authorization")))
	}))

	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().
		WithToken("abc123").
		Get(server.URL)

	if resp.Body() != "Bearer abc123" {
		t.Fatalf("expected 'Bearer abc123', got %s", resp.Body())
	}
}

func TestPendingRequestTimeout(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.Write([]byte("slow"))
	}))

	defer server.Close()

	_, err := client.NewFactory().PendingRequest().
		Timeout(50 * time.Millisecond).
		Get(server.URL)

	if err == nil {
		t.Fatal("expected timeout error")
	}
}

func TestPendingRequestPut(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("expected PUT, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().Put(server.URL)

	if !resp.Ok() {
		t.Fatal("expected ok")
	}
}

func TestPendingRequestPatch(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Fatalf("expected PATCH, got %s", r.Method)
		}
	}))

	defer server.Close()

	client.NewFactory().PendingRequest().Patch(server.URL)
}

func TestPendingRequestDelete(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Fatalf("expected DELETE, got %s", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
	}))

	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().Delete(server.URL)

	if !resp.NoContent() {
		t.Fatal("expected 204")
	}
}

func TestPendingRequestWithoutRedirecting(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/redirected", http.StatusFound)
	}))

	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().
		WithoutRedirecting().
		Get(server.URL)

	if !resp.Found() {
		t.Fatalf("expected 302, got %d", resp.Status())
	}
}

func TestPendingRequestMiddleware(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.Header.Get("X-Middleware")))
	}))

	defer server.Close()

	mw := client.Middleware(func(req *http.Request, next client.RoundTripFunc) (*http.Response, error) {
		req.Header.Set("X-Middleware", "injected")

		return next(req)
	})

	resp, _ := client.NewFactory().PendingRequest().
		WithMiddleware(mw).
		Get(server.URL)

	if resp.Body() != "injected" {
		t.Fatalf("expected injected, got %s", resp.Body())
	}
}

func TestPendingRequestHead(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodHead {
			t.Fatalf("expected HEAD, got %s", r.Method)
		}

		w.Header().Set("X-Custom", "exists")
	}))

	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().Head(server.URL)

	if resp.Header("X-Custom") != "exists" {
		t.Fatal("expected X-Custom header")
	}
}

func TestPendingRequestWithQueryParameters(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		sort := r.URL.Query().Get("sort")

		w.Write([]byte(page + ":" + sort))
	}))

	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().
		WithQueryParameters(map[string]string{"page": "1", "sort": "name"}).
		Get(server.URL)

	if resp.Body() != "1:name" {
		t.Fatalf("expected 1:name, got %s", resp.Body())
	}
}

func TestPendingRequestWithQueryParametersMergedWithGet(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		filter := r.URL.Query().Get("filter")

		w.Write([]byte(page + ":" + filter))
	}))

	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().
		WithQueryParameters(map[string]string{"page": "1"}).
		Get(server.URL, map[string]string{"filter": "active"})

	if resp.Body() != "1:active" {
		t.Fatalf("expected 1:active, got %s", resp.Body())
	}
}

func TestPendingRequestContentType(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.Header.Get("Content-Type")))
	}))

	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().
		ContentType("text/xml").
		Post(server.URL, "data")

	if resp.Body() != "text/xml" {
		t.Fatalf("expected text/xml, got %s", resp.Body())
	}
}

func TestPendingRequestWithUserAgent(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.Header.Get("User-Agent")))
	}))

	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().
		WithUserAgent("MyApp/1.0").
		Get(server.URL)

	if resp.Body() != "MyApp/1.0" {
		t.Fatalf("expected MyApp/1.0, got %s", resp.Body())
	}
}

func TestPendingRequestWithUrlParameters(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.URL.Path))
	}))

	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().
		WithUrlParameters(map[string]string{"version": "v1", "id": "42"}).
		Get(server.URL + "/api/{version}/users/{id}")

	if resp.Body() != "/api/v1/users/42" {
		t.Fatalf("expected /api/v1/users/42, got %s", resp.Body())
	}
}

func TestPendingRequestReplaceHeaders(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		old := r.Header.Get("X-Old")
		new := r.Header.Get("X-New")

		w.Write([]byte(old + ":" + new))
	}))

	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().
		WithHeader("X-Old", "exists").
		ReplaceHeaders(map[string]string{"X-New": "replaced"}).
		Get(server.URL)

	if resp.Body() != ":replaced" {
		t.Fatalf("expected ':replaced', got %s", resp.Body())
	}
}

func TestPendingRequestSend(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.Method))
	}))

	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().Send("PATCH", server.URL)

	if resp.Body() != "PATCH" {
		t.Fatalf("expected PATCH, got %s", resp.Body())
	}
}

func TestPendingRequestSendWithData(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)

		var data map[string]string

		json.Unmarshal(body, &data)

		w.Write([]byte(data["key"]))
	}))

	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().
		AsJSON().
		Send("POST", server.URL, map[string]string{"key": "value"})

	if resp.Body() != "value" {
		t.Fatalf("expected value, got %s", resp.Body())
	}
}

func TestPendingRequestSink(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("sink content"))
	}))

	defer server.Close()

	var buf bytes.Buffer

	resp, _ := client.NewFactory().PendingRequest().
		Sink(&buf).
		Get(server.URL)

	if buf.String() != "sink content" {
		t.Fatalf("expected sink content in buffer, got %s", buf.String())
	}

	if resp.Body() != "sink content" {
		t.Fatalf("expected sink content in response, got %s", resp.Body())
	}
}

func TestPendingRequestAttach(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(32 << 20)

		file, header, err := r.FormFile("avatar")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		defer file.Close()

		content, _ := io.ReadAll(file)

		w.Write([]byte(header.Filename + ":" + string(content)))
	}))

	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().
		Attach("avatar", strings.NewReader("file content"), "photo.jpg").
		Post(server.URL)

	if resp.Body() != "photo.jpg:file content" {
		t.Fatalf("expected 'photo.jpg:file content', got %s", resp.Body())
	}
}

func TestPendingRequestAttachWithFormFields(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(32 << 20)

		name := r.FormValue("name")

		file, _, err := r.FormFile("doc")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		defer file.Close()

		content, _ := io.ReadAll(file)

		w.Write([]byte(name + ":" + string(content)))
	}))

	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().
		Attach("doc", strings.NewReader("doc content"), "readme.txt").
		Post(server.URL, map[string]string{"name": "Taylor"})

	if resp.Body() != "Taylor:doc content" {
		t.Fatalf("expected 'Taylor:doc content', got %s", resp.Body())
	}
}

func TestPendingRequestWithoutVerifying(t *testing.T) {
	t.Parallel()

	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("secure"))
	}))

	defer server.Close()

	resp, err := client.NewFactory().PendingRequest().
		WithoutVerifying().
		Get(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Body() != "secure" {
		t.Fatalf("expected secure, got %s", resp.Body())
	}
}

func TestPendingRequestBeforeSending(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(r.Header.Get("X-Before")))
	}))

	defer server.Close()

	resp, _ := client.NewFactory().PendingRequest().
		BeforeSending(func(req *http.Request) {
			req.Header.Set("X-Before", "injected")
		}).
		Get(server.URL)

	if resp.Body() != "injected" {
		t.Fatalf("expected injected, got %s", resp.Body())
	}
}

func TestPendingRequestAfterResponse(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}))

	defer server.Close()

	var capturedBody string

	client.NewFactory().PendingRequest().
		AfterResponse(func(resp *client.Response) {
			capturedBody = resp.Body()
		}).
		Get(server.URL)

	if capturedBody != "hello" {
		t.Fatalf("expected hello, got %s", capturedBody)
	}
}

func TestPendingRequestThrowOnFailure(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	defer server.Close()

	_, err := client.NewFactory().PendingRequest().
		Throw().
		Get(server.URL)

	if err == nil {
		t.Fatal("expected error for 500 with Throw()")
	}
}

func TestPendingRequestThrowOnFailureNotOnSuccess(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	defer server.Close()

	_, err := client.NewFactory().PendingRequest().
		Throw().
		Get(server.URL)

	if err != nil {
		t.Fatalf("unexpected error for 200 with Throw(): %v", err)
	}
}

func TestPendingRequestThrowWithCallback(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	defer server.Close()

	var callbackCalled bool

	_, err := client.NewFactory().PendingRequest().
		Throw(func(resp *client.Response, e error) {
			callbackCalled = true
		}).
		Get(server.URL)

	if err == nil {
		t.Fatal("expected error")
	}

	if !callbackCalled {
		t.Fatal("throw callback should have been called")
	}
}

func TestPendingRequestThrowIf(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	defer server.Close()

	_, err := client.NewFactory().PendingRequest().
		ThrowIf(true).
		Get(server.URL)

	if err == nil {
		t.Fatal("expected error when ThrowIf(true) and 500")
	}

	_, err = client.NewFactory().PendingRequest().
		ThrowIf(false).
		Get(server.URL)

	if err != nil {
		t.Fatalf("unexpected error when ThrowIf(false): %v", err)
	}
}

func TestPendingRequestThrowUnless(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	defer server.Close()

	_, err := client.NewFactory().PendingRequest().
		ThrowUnless(false).
		Get(server.URL)

	if err == nil {
		t.Fatal("expected error when ThrowUnless(false) and 500")
	}

	_, err = client.NewFactory().PendingRequest().
		ThrowUnless(true).
		Get(server.URL)

	if err != nil {
		t.Fatalf("unexpected error when ThrowUnless(true): %v", err)
	}
}

func TestPendingRequestStub(t *testing.T) {
	t.Parallel()

	resp, err := client.NewFactory().Fake().PendingRequest().
		Stub(func(req *http.Request) *http.Response {
			return &http.Response{
				StatusCode: 201,
				Header:     make(http.Header),
				Body:       io.NopCloser(strings.NewReader("stubbed")),
			}
		}).
		Get("http://example.com/test")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.Body() != "stubbed" {
		t.Fatalf("expected stubbed, got %s", resp.Body())
	}

	if resp.Status() != 201 {
		t.Fatalf("expected 201, got %d", resp.Status())
	}
}

func TestPendingRequestPreventStrayRequests(t *testing.T) {
	t.Parallel()

	_, err := client.NewFactory().Fake().PendingRequest().
		Stub(func(req *http.Request) *http.Response {
			return nil // does not handle the request
		}).
		PreventStrayRequests().
		Get("http://example.com/unmatched")

	if err == nil {
		t.Fatal("expected stray request error")
	}
}

func TestPendingRequestWithAttributes(t *testing.T) {
	t.Parallel()

	var captured string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))

	defer server.Close()

	client.NewFactory().PendingRequest().
		WithAttributes(map[string]any{"custom_key": "custom_value"}).
		BeforeSending(func(req *http.Request) {
			if val := req.Context().Value("custom_key"); val != nil {
				captured = val.(string)
			}
		}).
		Get(server.URL)

	if captured != "custom_value" {
		t.Fatalf("expected custom_value, got %s", captured)
	}
}

func TestPendingRequestConnectTimeout(t *testing.T) {
	t.Parallel()

	// Use a non-routable address to trigger a connect timeout.
	_, err := client.NewFactory().PendingRequest().
		ConnectTimeout(50 * time.Millisecond).
		Timeout(1 * time.Second).
		Get("http://10.255.255.1:80")

	if err == nil {
		t.Fatal("expected connect timeout error")
	}
}

func TestPendingRequestHandlerStatsIntegration(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	}))

	defer server.Close()

	resp, err := client.NewFactory().PendingRequest().Get(server.URL)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stats := resp.HandlerStats()

	if _, ok := stats["total_ms"]; !ok {
		t.Fatal("expected total_ms in handler stats")
	}

	totalMs := stats["total_ms"].(float64)

	if totalMs < 0 {
		t.Fatalf("total_ms should be non-negative, got %f", totalMs)
	}
}
