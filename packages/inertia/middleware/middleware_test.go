package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bedrock/packages/inertia/protocol"
	"github.com/bedrock/packages/inertia/middleware"
)

func TestMiddleware_SetsVaryHeader(t *testing.T) {
	t.Parallel()

	mw := middleware.New(middleware.Config{Version: "v1"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if got := w.Header().Get("Vary"); got != protocol.HeaderInertia {
		t.Errorf("Vary = %q, want %q", got, protocol.HeaderInertia)
	}
}

func TestMiddleware_NonInertiaPassesThrough(t *testing.T) {
	t.Parallel()

	called := false
	mw := middleware.New(middleware.Config{Version: "v1"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if !called {
		t.Error("next handler was not called for non-Inertia request")
	}

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestMiddleware_VersionMatch(t *testing.T) {
	t.Parallel()

	called := false
	mw := middleware.New(middleware.Config{Version: "v1"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderVersion, "v1")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if !called {
		t.Error("next handler was not called when versions match")
	}
}

func TestMiddleware_VersionMismatch(t *testing.T) {
	t.Parallel()

	called := false
	mw := middleware.New(middleware.Config{Version: "v2"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
	}))

	r := httptest.NewRequest(http.MethodGet, "/test", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderVersion, "v1")

	r.RequestURI = "/test"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if called {
		t.Error("next handler should not be called on version mismatch")
	}

	if w.Code != http.StatusConflict {
		t.Errorf("status = %d, want %d", w.Code, http.StatusConflict)
	}

	if loc := w.Header().Get(protocol.HeaderLocation); loc != "/test" {
		t.Errorf("X-Inertia-Location = %q, want %q", loc, "/test")
	}
}

func TestMiddleware_VersionMismatchOnlyGET(t *testing.T) {
	t.Parallel()

	called := false
	mw := middleware.New(middleware.Config{Version: "v2"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodPost, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderVersion, "v1")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if !called {
		t.Error("next handler should be called for POST even on version mismatch")
	}
}

func TestMiddleware_RedirectConversion_PUTTo303(t *testing.T) {
	t.Parallel()

	mw := middleware.New(middleware.Config{Version: "v1"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/target", http.StatusFound)
	}))

	r := httptest.NewRequest(http.MethodPut, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderVersion, "v1")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusSeeOther {
		t.Errorf("status = %d, want %d (303)", w.Code, http.StatusSeeOther)
	}
}

func TestMiddleware_RedirectConversion_PATCHTo303(t *testing.T) {
	t.Parallel()

	mw := middleware.New(middleware.Config{Version: "v1"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/target", http.StatusFound)
	}))

	r := httptest.NewRequest(http.MethodPatch, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderVersion, "v1")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusSeeOther {
		t.Errorf("status = %d, want %d (303)", w.Code, http.StatusSeeOther)
	}
}

func TestMiddleware_RedirectConversion_DELETETo303(t *testing.T) {
	t.Parallel()

	mw := middleware.New(middleware.Config{Version: "v1"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/target", http.StatusFound)
	}))

	r := httptest.NewRequest(http.MethodDelete, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderVersion, "v1")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusSeeOther {
		t.Errorf("status = %d, want %d (303)", w.Code, http.StatusSeeOther)
	}
}

func TestMiddleware_RedirectConversion_GETStays302(t *testing.T) {
	t.Parallel()

	mw := middleware.New(middleware.Config{Version: "v1"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/target", http.StatusFound)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderVersion, "v1")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusFound {
		t.Errorf("status = %d, want %d (302)", w.Code, http.StatusFound)
	}
}

func TestMiddleware_WriteWithoutExplicitWriteHeader(t *testing.T) {
	t.Parallel()

	mw := middleware.New(middleware.Config{Version: "v1"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write without calling WriteHeader first triggers implicit 200.
		w.Write([]byte("hello"))
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderVersion, "v1")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	if w.Body.String() != "hello" {
		t.Errorf("body = %q", w.Body.String())
	}
}

func TestMiddleware_DoubleWriteHeaderIgnored(t *testing.T) {
	t.Parallel()

	mw := middleware.New(middleware.Config{Version: "v1"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.WriteHeader(http.StatusInternalServerError) // Second call ignored.
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderVersion, "v1")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d (first call wins)", w.Code, http.StatusCreated)
	}
}

// --- Version Edge Cases ---

func TestMiddleware_EmptyClientVersionPassesThrough(t *testing.T) {
	t.Parallel()

	called := false
	mw := middleware.New(middleware.Config{Version: "v1"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	// No X-Inertia-Version header set.
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if !called {
		t.Error("handler should be called when client sends no version")
	}

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestMiddleware_EmptyServerVersionWithClientVersion(t *testing.T) {
	t.Parallel()

	mw := middleware.New(middleware.Config{Version: ""})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderVersion, "v1")

	r.RequestURI = "/"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusConflict {
		t.Errorf("status = %d, want %d (client version != empty server version)", w.Code, http.StatusConflict)
	}
}

func TestMiddleware_VersionMismatchPreservesQueryString(t *testing.T) {
	t.Parallel()

	mw := middleware.New(middleware.Config{Version: "v2"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	r := httptest.NewRequest(http.MethodGet, "/users?page=2&sort=name", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderVersion, "v1")

	r.RequestURI = "/users?page=2&sort=name"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if loc := w.Header().Get(protocol.HeaderLocation); loc != "/users?page=2&sort=name" {
		t.Errorf("X-Inertia-Location = %q, want %q", loc, "/users?page=2&sort=name")
	}
}

func TestMiddleware_VersionMismatchOnPOST_NoConflict(t *testing.T) {
	t.Parallel()

	called := false
	mw := middleware.New(middleware.Config{Version: "v2"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodPost, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderVersion, "v1")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if !called {
		t.Error("POST should not trigger version conflict")
	}
}

func TestMiddleware_VersionMismatchOnDELETE_NoConflict(t *testing.T) {
	t.Parallel()

	called := false
	mw := middleware.New(middleware.Config{Version: "v2"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodDelete, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderVersion, "v1")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if !called {
		t.Error("DELETE should not trigger version conflict")
	}
}

// --- Redirect Conversion ---

func TestMiddleware_RedirectConversion_POSTStays302(t *testing.T) {
	t.Parallel()

	mw := middleware.New(middleware.Config{Version: "v1"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/target", http.StatusFound)
	}))

	r := httptest.NewRequest(http.MethodPost, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderVersion, "v1")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusFound {
		t.Errorf("status = %d, want %d (POST 302 not converted)", w.Code, http.StatusFound)
	}
}

func TestMiddleware_301NotConverted(t *testing.T) {
	t.Parallel()

	mw := middleware.New(middleware.Config{Version: "v1"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/target", http.StatusMovedPermanently)
	}))

	r := httptest.NewRequest(http.MethodPut, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderVersion, "v1")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusMovedPermanently {
		t.Errorf("status = %d, want %d (only 302 converted)", w.Code, http.StatusMovedPermanently)
	}
}

func TestMiddleware_307NotConverted(t *testing.T) {
	t.Parallel()

	mw := middleware.New(middleware.Config{Version: "v1"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/target", http.StatusTemporaryRedirect)
	}))

	r := httptest.NewRequest(http.MethodDelete, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderVersion, "v1")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("status = %d, want %d (307 not converted)", w.Code, http.StatusTemporaryRedirect)
	}
}

// --- Status Interceptor ---

func TestMiddleware_VaryHeaderOnInertiaRequest(t *testing.T) {
	t.Parallel()

	mw := middleware.New(middleware.Config{Version: "v1"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderVersion, "v1")

	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if got := w.Header().Get("Vary"); got != protocol.HeaderInertia {
		t.Errorf("Vary = %q, want %q", got, protocol.HeaderInertia)
	}
}

func TestMiddleware_VersionMismatchResponseHasNoBody(t *testing.T) {
	t.Parallel()

	mw := middleware.New(middleware.Config{Version: "v2"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

	r := httptest.NewRequest(http.MethodGet, "/", nil)

	r.Header.Set(protocol.HeaderInertia, "true")
	r.Header.Set(protocol.HeaderVersion, "v1")

	r.RequestURI = "/"
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Body.Len() > 0 {
		t.Errorf("409 body should be empty, got %q", w.Body.String())
	}
}

func TestMiddleware_NonInertia302NotConverted(t *testing.T) {
	t.Parallel()

	mw := middleware.New(middleware.Config{Version: "v1"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/target", http.StatusFound)
	}))

	// Non-Inertia PUT request — no status interceptor.
	r := httptest.NewRequest(http.MethodPut, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, r)

	if w.Code != http.StatusFound {
		t.Errorf("status = %d, want %d (non-Inertia 302 not converted)", w.Code, http.StatusFound)
	}
}

func TestMiddleware_ConcurrentRequests(t *testing.T) {
	t.Parallel()

	mw := middleware.New(middleware.Config{Version: "v1"})

	handler := mw(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	done := make(chan struct{})

	for n := 0; n < 20; n++ {
		go func() {
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			r.Header.Set(protocol.HeaderInertia, "true")
			r.Header.Set(protocol.HeaderVersion, "v1")

			w := httptest.NewRecorder()

			handler.ServeHTTP(w, r)

			done <- struct{}{}
		}()
	}

	for n := 0; n < 20; n++ {
		<-done
	}
}
