package http_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	bedhttp "github.com/bedrock/packages/anvil/http"
)

// Upstream: testJsonResponsesAreConvertedAndHeadersAreSet
func TestResponseJsonConversion(t *testing.T) {
	t.Parallel()

	resp, err := bedhttp.NewResponse(map[string]string{"name": "John"}, http.StatusOK)
	if err != nil {
		t.Fatalf("NewResponse: %v", err)
	}

	rec := httptest.NewRecorder()
	if err := resp.WriteTo(rec); err != nil {
		t.Fatalf("WriteTo: %v", err)
	}

	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected 'application/json', got %q", ct)
	}

	var parsed map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &parsed); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if parsed["name"] != "John" {
		t.Fatalf("expected 'John', got %q", parsed["name"])
	}
}

// Upstream: testHeader
func TestResponseHeader(t *testing.T) {
	t.Parallel()

	resp, _ := bedhttp.NewResponse("hello", http.StatusOK)
	resp.SetHeader("X-Custom", "value")

	if got := resp.GetHeader("X-Custom"); got != "value" {
		t.Fatalf("expected 'value', got %q", got)
	}
	if got := resp.GetHeader("X-Missing", "default"); got != "default" {
		t.Fatalf("expected 'default', got %q", got)
	}
}

// Upstream: testGetOriginalContent
func TestResponseGetOriginalContent(t *testing.T) {
	t.Parallel()

	resp, _ := bedhttp.NewResponse("hello world", http.StatusOK)
	got, ok := resp.GetOriginalContent().(string)
	if !ok || got != "hello world" {
		t.Fatalf("expected 'hello world', got %v", resp.GetOriginalContent())
	}
}

// Upstream: testGetOriginalContentRetrievesTheFirstOriginalContent
func TestResponseGetOriginalContentRetrievesFirst(t *testing.T) {
	t.Parallel()

	resp, _ := bedhttp.NewResponse("first", http.StatusOK)

	// Overwrite content.
	if err := resp.SetContent("second"); err != nil {
		t.Fatalf("SetContent: %v", err)
	}

	// Original content should be the last set value.
	got, ok := resp.GetOriginalContent().(string)
	if !ok || got != "second" {
		t.Fatalf("expected 'second', got %v", resp.GetOriginalContent())
	}
}

// Upstream: testSetAndRetrieveStatusCode
func TestResponseSetAndRetrieveStatusCode(t *testing.T) {
	t.Parallel()

	resp, _ := bedhttp.NewResponse("", http.StatusOK)
	if got := resp.GetStatusCode(); got != http.StatusOK {
		t.Fatalf("expected 200, got %d", got)
	}

	resp.SetStatusCode(http.StatusNotFound)
	if got := resp.GetStatusCode(); got != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", got)
	}
}

// Upstream: testSetStatusCodeAndRetrieveStatusText
func TestResponseStatusText(t *testing.T) {
	t.Parallel()

	resp, _ := bedhttp.NewResponse("", http.StatusOK)
	if got := resp.StatusText(); got != "OK" {
		t.Fatalf("expected 'OK', got %q", got)
	}

	resp.SetStatusCode(http.StatusNotFound)
	if got := resp.StatusText(); got != "Not Found" {
		t.Fatalf("expected 'Not Found', got %q", got)
	}

	resp.SetStatusCode(http.StatusInternalServerError)
	if got := resp.StatusText(); got != "Internal Server Error" {
		t.Fatalf("expected 'Internal Server Error', got %q", got)
	}
}

// Upstream: testWithHeaders
func TestResponseWithHeaders(t *testing.T) {
	t.Parallel()

	resp, _ := bedhttp.NewResponse("", http.StatusOK)
	resp.WithHeaders(map[string]string{
		"X-Foo": "foo",
		"X-Bar": "bar",
	})

	if got := resp.GetHeader("X-Foo"); got != "foo" {
		t.Fatalf("expected 'foo', got %q", got)
	}
	if got := resp.GetHeader("X-Bar"); got != "bar" {
		t.Fatalf("expected 'bar', got %q", got)
	}
}

// Upstream: testWithoutHeader
func TestResponseWithoutHeader(t *testing.T) {
	t.Parallel()

	resp, _ := bedhttp.NewResponse("", http.StatusOK)
	resp.SetHeader("X-Foo", "foo")
	resp.SetHeader("X-Bar", "bar")
	resp.SetHeader("X-Baz", "baz")

	resp.WithoutHeader("X-Foo", "X-Bar")

	if got := resp.GetHeader("X-Foo"); got != "" {
		t.Fatalf("expected X-Foo to be removed, got %q", got)
	}
	if got := resp.GetHeader("X-Bar"); got != "" {
		t.Fatalf("expected X-Bar to be removed, got %q", got)
	}
	if got := resp.GetHeader("X-Baz"); got != "baz" {
		t.Fatalf("expected X-Baz to remain 'baz', got %q", got)
	}
}
