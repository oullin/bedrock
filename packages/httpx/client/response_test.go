package client_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/bedrock/packages/httpx/client"
)

func makeResponse(status int, body string, headers ...map[string]string) *client.Response {
	raw := &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			raw.Header.Set(k, v)
		}
	}

	return client.NewResponse(raw)
}

func TestResponseStatus(t *testing.T) {
	t.Parallel()

	resp := makeResponse(http.StatusOK, "")

	if resp.Status() != 200 {
		t.Fatalf("expected 200, got %d", resp.Status())
	}
}

func TestResponseBody(t *testing.T) {
	t.Parallel()

	resp := makeResponse(http.StatusOK, "hello world")

	if resp.Body() != "hello world" {
		t.Fatalf("expected 'hello world', got %s", resp.Body())
	}
}

func TestResponseJSON(t *testing.T) {
	t.Parallel()

	resp := makeResponse(http.StatusOK, `{"name":"Taylor"}`)

	var data map[string]string

	if err := resp.JSON(&data); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data["name"] != "Taylor" {
		t.Fatalf("expected Taylor, got %s", data["name"])
	}
}

func TestResponseStatusHelpers(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status int
		check  func(*client.Response) bool
		name   string
	}{
		{200, (*client.Response).Ok, "Ok"},
		{201, (*client.Response).Created, "Created"},
		{202, (*client.Response).Accepted, "Accepted"},
		{204, (*client.Response).NoContent, "NoContent"},
		{301, (*client.Response).MovedPermanently, "MovedPermanently"},
		{302, (*client.Response).Found, "Found"},
		{304, (*client.Response).NotModified, "NotModified"},
		{400, (*client.Response).BadRequest, "BadRequest"},
		{401, (*client.Response).Unauthorized, "Unauthorized"},
		{402, (*client.Response).PaymentRequired, "PaymentRequired"},
		{403, (*client.Response).Forbidden, "Forbidden"},
		{404, (*client.Response).NotFound, "NotFound"},
		{408, (*client.Response).RequestTimeout, "RequestTimeout"},
		{409, (*client.Response).Conflict, "Conflict"},
		{422, (*client.Response).UnprocessableEntity, "UnprocessableEntity"},
		{429, (*client.Response).TooManyRequests, "TooManyRequests"},
	}

	for _, tt := range tests {
		resp := makeResponse(tt.status, "")

		if !tt.check(resp) {
			t.Errorf("%s() should return true for status %d", tt.name, tt.status)
		}
	}
}

func TestResponseRangeHelpers(t *testing.T) {
	t.Parallel()

	resp200 := makeResponse(200, "")

	if !resp200.Successful() {
		t.Fatal("200 should be successful")
	}

	resp301 := makeResponse(301, "")

	if !resp301.Redirect() {
		t.Fatal("301 should be redirect")
	}

	resp404 := makeResponse(404, "")

	if !resp404.ClientError() {
		t.Fatal("404 should be client error")
	}

	if !resp404.Failed() {
		t.Fatal("404 should be failed")
	}

	resp500 := makeResponse(500, "")

	if !resp500.ServerError() {
		t.Fatal("500 should be server error")
	}
}

func TestResponseThrow(t *testing.T) {
	t.Parallel()

	resp200 := makeResponse(200, "")

	if resp200.Throw() != nil {
		t.Fatal("200 should not throw")
	}

	resp500 := makeResponse(500, "")

	if resp500.Throw() == nil {
		t.Fatal("500 should throw")
	}
}

func TestResponseThrowIf(t *testing.T) {
	t.Parallel()

	resp500 := makeResponse(500, "")

	if resp500.ThrowIf(false) != nil {
		t.Fatal("should not throw when condition is false")
	}

	if resp500.ThrowIf(true) == nil {
		t.Fatal("should throw when condition is true")
	}
}

func TestResponseHeader(t *testing.T) {
	t.Parallel()

	resp := makeResponse(200, "", map[string]string{"X-Custom": "value"})

	if resp.Header("X-Custom") != "value" {
		t.Fatalf("expected value, got %s", resp.Header("X-Custom"))
	}
}

func TestResponseHeaders(t *testing.T) {
	t.Parallel()

	resp := makeResponse(200, "", map[string]string{"X-A": "1", "X-B": "2"})

	headers := resp.Headers()

	if headers.Get("X-A") != "1" || headers.Get("X-B") != "2" {
		t.Fatal("expected both headers")
	}
}

func TestResponseCookies(t *testing.T) {
	t.Parallel()

	raw := &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader("")),
	}
	raw.Header.Add("Set-Cookie", "session=abc; Path=/")

	resp := client.NewResponse(raw)

	if len(resp.Cookies()) == 0 {
		t.Fatal("expected cookies")
	}
}

func TestResponseRaw(t *testing.T) {
	t.Parallel()

	resp := makeResponse(200, "")

	if resp.Raw() == nil {
		t.Fatal("expected non-nil raw response")
	}
}

func TestResponseBodyReadOnce(t *testing.T) {
	t.Parallel()

	resp := makeResponse(200, "content")

	// Read twice - should return the same cached content.
	if resp.Body() != "content" {
		t.Fatal("first read failed")
	}

	if resp.Body() != "content" {
		t.Fatal("second read should return cached content")
	}
}

func TestResponseReason(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status int
		reason string
	}{
		{200, "OK"},
		{201, "Created"},
		{301, "Moved Permanently"},
		{404, "Not Found"},
		{500, "Internal Server Error"},
	}

	for _, tt := range tests {
		resp := makeResponse(tt.status, "")

		if resp.Reason() != tt.reason {
			t.Errorf("expected reason %q for %d, got %q", tt.reason, tt.status, resp.Reason())
		}
	}
}

func TestResponseCollect(t *testing.T) {
	t.Parallel()

	resp := makeResponse(200, `{"name":"Taylor","age":30}`)
	data := resp.Collect()

	if data == nil {
		t.Fatal("expected non-nil map")
	}

	if data["name"] != "Taylor" {
		t.Fatalf("expected Taylor, got %v", data["name"])
	}
}

func TestResponseCollectWithKey(t *testing.T) {
	t.Parallel()

	resp := makeResponse(200, `{"user":{"name":"Taylor"}}`)
	data := resp.Collect("user")

	if data == nil {
		t.Fatal("expected non-nil map")
	}

	if data["name"] != "Taylor" {
		t.Fatalf("expected Taylor, got %v", data["name"])
	}
}

func TestResponseCollectWithMissingKey(t *testing.T) {
	t.Parallel()

	resp := makeResponse(200, `{"user":{"name":"Taylor"}}`)
	data := resp.Collect("missing")

	if data != nil {
		t.Fatal("expected nil for missing key")
	}
}

func TestResponseCollectInvalidJSON(t *testing.T) {
	t.Parallel()

	resp := makeResponse(200, "not json")
	data := resp.Collect()

	if data != nil {
		t.Fatal("expected nil for invalid JSON")
	}
}

func TestResponseEffectiveUri(t *testing.T) {
	t.Parallel()

	reqURL := "http://example.com/final"

	req, _ := http.NewRequest("GET", reqURL, nil)

	raw := &http.Response{
		StatusCode: 200,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader("")),
		Request:    req,
	}

	resp := client.NewResponse(raw)

	if resp.EffectiveUri() != reqURL {
		t.Fatalf("expected %s, got %s", reqURL, resp.EffectiveUri())
	}
}

func TestResponseEffectiveUriEmpty(t *testing.T) {
	t.Parallel()

	resp := makeResponse(200, "")

	if resp.EffectiveUri() != "" {
		t.Fatal("expected empty string when no request set")
	}
}

func TestResponseThrowUnless(t *testing.T) {
	t.Parallel()

	resp500 := makeResponse(500, "")

	if resp500.ThrowUnless(true) != nil {
		t.Fatal("should not throw when condition is true")
	}

	if resp500.ThrowUnless(false) == nil {
		t.Fatal("should throw when condition is false")
	}

	resp200 := makeResponse(200, "")

	if resp200.ThrowUnless(false) != nil {
		t.Fatal("should not throw for successful response")
	}
}

func TestResponseThrowIfStatus(t *testing.T) {
	t.Parallel()

	resp200 := makeResponse(200, "")

	if resp200.ThrowIfStatus(200) == nil {
		t.Fatal("should throw when status matches, even for 200")
	}

	if resp200.ThrowIfStatus(404) != nil {
		t.Fatal("should not throw when status does not match")
	}
}

func TestResponseThrowUnlessStatus(t *testing.T) {
	t.Parallel()

	resp200 := makeResponse(200, "")

	if resp200.ThrowUnlessStatus(200) != nil {
		t.Fatal("should not throw when status matches")
	}

	if resp200.ThrowUnlessStatus(201) == nil {
		t.Fatal("should throw when status does not match")
	}
}

func TestResponseThrowIfClientError(t *testing.T) {
	t.Parallel()

	resp404 := makeResponse(404, "")

	if resp404.ThrowIfClientError() == nil {
		t.Fatal("should throw for 4xx")
	}

	resp200 := makeResponse(200, "")

	if resp200.ThrowIfClientError() != nil {
		t.Fatal("should not throw for 2xx")
	}

	resp500 := makeResponse(500, "")

	if resp500.ThrowIfClientError() != nil {
		t.Fatal("should not throw for 5xx")
	}
}

func TestResponseThrowIfServerError(t *testing.T) {
	t.Parallel()

	resp500 := makeResponse(500, "")

	if resp500.ThrowIfServerError() == nil {
		t.Fatal("should throw for 5xx")
	}

	resp200 := makeResponse(200, "")

	if resp200.ThrowIfServerError() != nil {
		t.Fatal("should not throw for 2xx")
	}

	resp404 := makeResponse(404, "")

	if resp404.ThrowIfServerError() != nil {
		t.Fatal("should not throw for 4xx")
	}
}

func TestResponseOnError(t *testing.T) {
	t.Parallel()

	called := false
	resp500 := makeResponse(500, "")

	result := resp500.OnError(func(r *client.Response) {
		called = true
	})

	if !called {
		t.Fatal("callback should be called on failure")
	}

	if result != resp500 {
		t.Fatal("should return self for chaining")
	}
}

func TestResponseOnErrorNotCalledOnSuccess(t *testing.T) {
	t.Parallel()

	called := false
	resp200 := makeResponse(200, "")

	resp200.OnError(func(r *client.Response) {
		called = true
	})

	if called {
		t.Fatal("callback should not be called on success")
	}
}

func TestResponseClose(t *testing.T) {
	t.Parallel()

	resp := makeResponse(200, "content")

	if err := resp.Close(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResponseCloseAfterRead(t *testing.T) {
	t.Parallel()

	resp := makeResponse(200, "content")

	_ = resp.Body() // trigger read

	if err := resp.Close(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResponseToException(t *testing.T) {
	t.Parallel()

	resp500 := makeResponse(500, "")

	ex := resp500.ToException()

	if ex == nil {
		t.Fatal("expected non-nil exception for 5xx")
	}

	if ex.Response != resp500 {
		t.Fatal("exception should reference the response")
	}
}

func TestResponseToExceptionOnSuccess(t *testing.T) {
	t.Parallel()

	resp200 := makeResponse(200, "")

	if resp200.ToException() != nil {
		t.Fatal("expected nil for successful response")
	}
}

func TestResponseHandlerStats(t *testing.T) {
	t.Parallel()

	resp := makeResponse(200, "")

	stats := resp.HandlerStats()

	if stats == nil {
		t.Fatal("expected non-nil stats map")
	}

	if len(stats) != 0 {
		t.Fatalf("expected empty stats, got %d entries", len(stats))
	}
}

func TestResponseHandlerStatsWithData(t *testing.T) {
	t.Parallel()

	resp := makeResponse(200, "")
	resp.SetStats(map[string]any{
		"total_ms":   42.0,
		"connect_ms": 5.0,
	})

	stats := resp.HandlerStats()

	if stats["total_ms"] != 42.0 {
		t.Fatalf("expected total_ms=42, got %v", stats["total_ms"])
	}

	if stats["connect_ms"] != 5.0 {
		t.Fatalf("expected connect_ms=5, got %v", stats["connect_ms"])
	}
}
