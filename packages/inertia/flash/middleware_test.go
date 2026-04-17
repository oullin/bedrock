package flash

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/bedrock/packages/inertia"
)

func TestMiddlewareConsumesFlash(t *testing.T) {
	t.Parallel()

	s := NewCookieStore(WithCookieName("test_flash"))

	var capturedCtx context.Context

	handler := Middleware(s)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedCtx = r.Context()

		w.WriteHeader(http.StatusOK)
	}))

	msg := Message{Kind: "success", Title: "Done", Message: "Created."}
	data, _ := json.Marshal(msg)

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.AddCookie(&http.Cookie{
		Name:  "test_flash",
		Value: url.QueryEscape(string(data)),
	})

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	// Verify the cookie was consumed (deleted).
	cookies := rec.Result().Cookies()
	found := false

	for _, c := range cookies {
		if c.Name == "test_flash" && c.MaxAge == -1 {
			found = true

			break
		}
	}

	if !found {
		t.Error("expected flash cookie to be deleted after middleware consumed it")
	}

	// Verify the flash prop was set under the default "flash" key.
	props := inertia.PropsFromContext(capturedCtx)

	got, ok := props["flash"].(*Message)

	if !ok || got == nil {
		t.Fatal("expected flash prop to be set in context under \"flash\" key")
	}

	if got.Kind != msg.Kind || got.Title != msg.Title || got.Message != msg.Message {
		t.Errorf("flash prop mismatch: got %+v, want %+v", got, msg)
	}
}

func TestMiddlewareNoFlash(t *testing.T) {
	t.Parallel()

	s := NewCookieStore(WithCookieName("test_flash"))
	called := false

	handler := Middleware(s)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !called {
		t.Error("expected handler to be called even without flash")
	}

	cookies := rec.Result().Cookies()

	for _, c := range cookies {
		if c.Name == "test_flash" {
			t.Error("expected no flash cookie manipulation when no flash exists")
		}
	}
}

func TestMiddlewareNilStoreNoOp(t *testing.T) {
	t.Parallel()

	called := false

	handler := Middleware(nil)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if !called {
		t.Error("expected handler to be called with nil store (no-op middleware)")
	}
}

func TestMiddlewareCustomPropKey(t *testing.T) {
	t.Parallel()

	s := NewCookieStore(WithCookieName("test_flash"))

	var capturedCtx context.Context

	handler := Middleware(s, WithPropKey("notification"))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedCtx = r.Context()

		w.WriteHeader(http.StatusOK)
	}))

	msg := Message{Kind: "error", Title: "Oops", Message: "Failed."}
	data, _ := json.Marshal(msg)

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req.AddCookie(&http.Cookie{
		Name:  "test_flash",
		Value: url.QueryEscape(string(data)),
	})

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	props := inertia.PropsFromContext(capturedCtx)

	// The prop must appear under "notification", not the default "flash".
	got, ok := props["notification"].(*Message)

	if !ok || got == nil {
		t.Fatal("expected flash prop under \"notification\" key")
	}

	if got.Kind != msg.Kind || got.Title != msg.Title || got.Message != msg.Message {
		t.Errorf("prop mismatch: got %+v, want %+v", got, msg)
	}

	if _, exists := props["flash"]; exists {
		t.Error("expected no prop under default \"flash\" key when custom key is set")
	}
}
