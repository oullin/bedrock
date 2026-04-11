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
