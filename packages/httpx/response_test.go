package httpx_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/httpx"
)

func TestResponseStatus(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	resp := httpx.NewResponse(rec)
	resp.Status(http.StatusCreated)

	if resp.GetStatusCode() != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.GetStatusCode())
	}

	resp.Send([]byte("created"))

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201 in recorder, got %d", rec.Code)
	}
}

func TestResponseHeaders(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	resp := httpx.NewResponse(rec)
	resp.Header("X-Custom", "value").
		WithHeaders(map[string]string{"X-Another": "another"})

	resp.Send([]byte("ok"))

	if rec.Header().Get("X-Custom") != "value" {
		t.Fatal("expected X-Custom header")
	}

	if rec.Header().Get("X-Another") != "another" {
		t.Fatal("expected X-Another header")
	}
}

func TestResponseWithoutHeader(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	resp := httpx.NewResponse(rec)
	resp.Header("X-Remove", "yes").WithoutHeader("X-Remove")

	resp.Send([]byte("ok"))

	if rec.Header().Get("X-Remove") != "" {
		t.Fatal("expected X-Remove header to be removed")
	}
}

func TestResponseGetHeader(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	resp := httpx.NewResponse(rec)
	resp.Header("X-Test", "hello")

	if resp.GetHeader("X-Test") != "hello" {
		t.Fatal("expected to get buffered header")
	}
}

func TestResponseCookies(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	resp := httpx.NewResponse(rec)
	resp.Cookie(&http.Cookie{Name: "session", Value: "abc", Path: "/"})

	resp.Send([]byte("ok"))

	cookies := rec.Result().Cookies()

	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}

	if cookies[0].Name != "session" || cookies[0].Value != "abc" {
		t.Fatalf("unexpected cookie: %v", cookies[0])
	}
}

func TestResponseWithoutCookie(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	resp := httpx.NewResponse(rec)
	resp.WithoutCookie("session")

	resp.Send([]byte("ok"))

	cookies := rec.Result().Cookies()

	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie deletion, got %d", len(cookies))
	}

	if cookies[0].MaxAge != -1 {
		t.Fatal("expected MaxAge -1 for cookie deletion")
	}
}

func TestResponseSendString(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	resp := httpx.NewResponse(rec)

	resp.SendString("hello world")

	if rec.Body.String() != "hello world" {
		t.Fatalf("expected 'hello world', got %s", rec.Body.String())
	}

	if rec.Header().Get("Content-Type") != "text/plain; charset=utf-8" {
		t.Fatal("expected text/plain content type")
	}
}

func TestResponseJSON(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	resp := httpx.NewResponse(rec)

	data := map[string]string{"name": "Taylor"}
	resp.Status(http.StatusOK).JSON(data)

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

func TestResponseNoContent(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	resp := httpx.NewResponse(rec)

	resp.NoContent()

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}

	if rec.Body.Len() != 0 {
		t.Fatal("expected empty body")
	}
}

func TestResponseNoContentCustomStatus(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	resp := httpx.NewResponse(rec)

	resp.NoContent(http.StatusResetContent)

	if rec.Code != http.StatusResetContent {
		t.Fatalf("expected 205, got %d", rec.Code)
	}
}

func TestResponseWithException(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	resp := httpx.NewResponse(rec)

	err := errors.New("test error")
	resp.WithException(err)

	if resp.GetException() != err {
		t.Fatal("expected exception to be stored")
	}
}

func TestResponseOriginal(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	resp := httpx.NewResponse(rec)

	original := map[string]string{"key": "value"}
	resp.SetOriginal(original)

	if resp.GetOriginal() == nil {
		t.Fatal("expected original to be stored")
	}
}

func TestResponseWriter(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	resp := httpx.NewResponse(rec)

	if resp.Writer() != rec {
		t.Fatal("expected Writer() to return the underlying writer")
	}
}

func TestResponseDefaultStatus(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	resp := httpx.NewResponse(rec)

	if resp.GetStatusCode() != http.StatusOK {
		t.Fatalf("expected default status 200, got %d", resp.GetStatusCode())
	}
}
