package httpx_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bedrock/packages/httpx"
)

func TestJsonResponseSend(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	data := map[string]string{"name": "Taylor"}
	resp := httpx.NewJsonResponse(rec, data, http.StatusOK)

	if err := resp.Send(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	if rec.Header().Get("Content-Type") != "application/json" {
		t.Fatal("expected application/json content type")
	}

	body := rec.Body.String()

	if body != `{"name":"Taylor"}` {
		t.Fatalf("unexpected body: %s", body)
	}
}

func TestJsonResponseWithCallback(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	data := map[string]int{"count": 42}
	resp := httpx.NewJsonResponse(rec, data, http.StatusOK).
		WithCallback("handleResponse")

	if err := resp.Send(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec.Header().Get("Content-Type") != "text/javascript; charset=utf-8" {
		t.Fatalf("expected text/javascript, got %s", rec.Header().Get("Content-Type"))
	}

	body := rec.Body.String()

	if !strings.HasPrefix(body, "/**/ handleResponse(") {
		t.Fatalf("expected JSONP callback wrapper, got %s", body)
	}

	if !strings.HasSuffix(body, ");") {
		t.Fatalf("expected JSONP suffix, got %s", body)
	}
}

func TestJsonResponseStatus(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	resp := httpx.NewJsonResponse(rec, nil, http.StatusOK).Status(http.StatusCreated)

	resp.Send()

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", rec.Code)
	}
}

func TestJsonResponseHeaders(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	resp := httpx.NewJsonResponse(rec, "ok", http.StatusOK).
		Header("X-Custom", "value").
		WithHeaders(map[string]string{"X-Another": "yes"})

	resp.Send()

	if rec.Header().Get("X-Custom") != "value" {
		t.Fatal("expected X-Custom header")
	}

	if rec.Header().Get("X-Another") != "yes" {
		t.Fatal("expected X-Another header")
	}
}

func TestJsonResponseCookie(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	resp := httpx.NewJsonResponse(rec, "ok", http.StatusOK).
		Cookie(&http.Cookie{Name: "token", Value: "abc", Path: "/"})

	resp.Send()

	cookies := rec.Result().Cookies()

	if len(cookies) != 1 || cookies[0].Name != "token" {
		t.Fatal("expected cookie to be set")
	}
}

func TestJsonResponseSetData(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	resp := httpx.NewJsonResponse(rec, "original", http.StatusOK)
	resp.SetData("updated")

	resp.Send()

	if rec.Body.String() != `"updated"` {
		t.Fatalf("expected updated data, got %s", rec.Body.String())
	}
}

func TestJsonResponseGetData(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	data := map[string]int{"count": 1}
	resp := httpx.NewJsonResponse(rec, data, http.StatusOK)

	if resp.GetData() == nil {
		t.Fatal("expected non-nil data")
	}
}

func TestJsonResponseHasValidJSON(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()

	validResp := httpx.NewJsonResponse(rec, map[string]string{"ok": "yes"}, 200)

	if !validResp.HasValidJSON() {
		t.Fatal("expected valid JSON")
	}

	invalidResp := httpx.NewJsonResponse(rec, make(chan int), 200)

	if invalidResp.HasValidJSON() {
		t.Fatal("expected invalid JSON for channel type")
	}
}

func TestJsonResponseWithIndent(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	data := map[string]string{"name": "Taylor"}
	opts := httpx.JsonOptions{Indent: true}
	resp := httpx.NewJsonResponse(rec, data, http.StatusOK, opts)

	resp.Send()

	body := rec.Body.String()

	if !strings.Contains(body, "    ") {
		t.Fatal("expected indented JSON output")
	}
}

func TestJsonResponseGetOriginal(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	data := map[string]string{"key": "val"}
	resp := httpx.NewJsonResponse(rec, data, http.StatusOK)

	resp.Send()

	if resp.GetOriginal() == nil {
		t.Fatal("expected original data after send")
	}
}

func TestJsonResponseWithException(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	resp := httpx.NewJsonResponse(rec, nil, http.StatusOK)

	err := errors.New("test")
	resp.WithException(err)

	if resp.GetException() != err {
		t.Fatal("expected stored exception")
	}
}

func TestFromJsonString(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	resp := httpx.FromJsonString(rec, `{"raw":true}`, http.StatusOK)

	resp.Send()

	if rec.Body.String() != `{"raw":true}` {
		t.Fatalf("expected raw JSON string, got %s", rec.Body.String())
	}
}

func TestJsonResponseGetStatusCode(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	resp := httpx.NewJsonResponse(rec, nil, http.StatusAccepted)

	if resp.GetStatusCode() != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", resp.GetStatusCode())
	}
}

func TestJsonResponseEncodingOptions(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	resp := httpx.NewJsonResponse(rec, "test", http.StatusOK)
	resp.SetEncodingOptions(httpx.JsonOptions{Indent: true})

	resp.Send()

	// With indent, a simple string should still be quoted.
	if rec.Body.String() != `"test"` {
		t.Fatalf("unexpected body: %s", rec.Body.String())
	}
}
